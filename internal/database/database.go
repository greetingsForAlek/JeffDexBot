package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/greetingsForAlek/JeffDexBot/internal/models"
	_ "modernc.org/sqlite"
)

type Database struct {
	DB *sql.DB
}

func New() (*Database, error) {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./character-bot.db"
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}

func (d *Database) Init() error {
	query := `
	CREATE TABLE IF NOT EXISTS characters (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		image_url TEXT NOT NULL,
		rarity TEXT NOT NULL DEFAULT 'Common'
	);

	CREATE TABLE IF NOT EXISTS collections (
		user_id TEXT NOT NULL,
		character_id INTEGER NOT NULL,
		collected_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

		PRIMARY KEY (user_id, character_id),

		FOREIGN KEY (character_id)
			REFERENCES characters(id)
			ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS user_balances (
		user_id PRIMARY KEY,
		coins INTEGER DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS daily_shop (
		character_id INTEGER NOT NULL,
		shop_date TEXT NOT NULL,
		price INTEGER NOT NULL,

		PRIMARY KEY (character_id, shop_date),

		FOREIGN KEY (character_id)
			REFERENCES characters(id)
			ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS daily_shop_purchases (
		user_id TEXT NOT NULL,
		character_id INTEGER NOT NULL,
		shop_date TEXT NOT NULL,

		PRIMARY KEY (user_id, character_id, shop_date),

		FOREIGN KEY (character_id)
			REFERENCES characters(id)
			ON DELETE CASCADE
	);
	`

	_, err := d.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("create database tables: %w", err)
	}

	if err := d.ensureRarityColumn(); err != nil {
		return err
	}

	return nil
}

func (d *Database) AddCharacter (
	id int,
	name string,
	imageURL string,
	rarity string,
) error {
	query := `
	INSERT INTO characters (id, name, image_url, rarity)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		image_url = excluded.image_url,
		rarity = excluded.rarity;
	`

	_, err := d.DB.Exec(query, id, name, imageURL, rarity)
	if err != nil {
		return fmt.Errorf("add character: %w", err)
	}

	return nil
}

func (d *Database) GetCharacters() ([]models.Character, error) {
	query := `
	SELECT id, name, image_url, rarity
	FROM characters
	ORDER BY id;
	`

	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query characters: %w", err)
	}
	defer rows.Close()

	var characters []models.Character

	for rows.Next() {
		var character models.Character

		err := rows.Scan(
			&character.ID,
			&character.Name,
			&character.ImageURL,
			&character.Rarity,
		)
		if err != nil {
			return nil, fmt.Errorf("scan character: %w", err)
		}

		characters = append(characters, character)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate characters: %w", err)
	}

	return characters, nil
}

func (d *Database) AddToCollection(
	userID string,
	characterID int,
) (bool, error) {
	query := `
	INSERT INTO collections (user_id, character_id)
	VALUES (?, ?)
	ON CONFLICT(user_id, character_id) DO NOTHING;
	`

	result, err := d.DB.Exec(query, userID, characterID)
	if err != nil {
		return false, fmt.Errorf("add to collection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("check collection insert: %w", err)
	}

	return rowsAffected > 0, nil
}

func (d *Database) GetCollection (
	userID string,
) ([]models.Character, error) {
	query := `
		SELECT c.id, c.name, c.image_url
		FROM characters c
		INNER JOIN collections col
			ON col.character_id = c.id
		WHERE col.user_id  = ?
		ORDER BY col.collected_at, c.id;
	`

	rows, err := d.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("get collection: %w", err)
	}
	defer rows.Close()

	var characters []models.Character

	for rows.Next() {
		var character models.Character

		if err := rows.Scan(
			&character.ID,
			&character.Name,
			&character.ImageURL,
		); err != nil {
			return nil, fmt.Errorf("scan collection character: %w", err)
		}

		characters = append(characters, character)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate collection: %w", err)
	}

	return characters, nil
}

func (d *Database) ensureRarityColumn() error {
	rows, err := d.DB.Query("PRAGMA table_info(characters)")
	if err != nil {
		return fmt.Errorf("inspect characters table: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cid int
			name string
			columnType string
			notNull int
			defaultValue sql.NullString
			primaryKey int
		)

		if err := rows.Scan(
			&cid,
			&name,
			&columnType,
			&notNull,
			&defaultValue,
			&primaryKey,
		); err != nil {
			rows.Close()
			return fmt.Errorf("inspect character column: %w", err)
		}

		if name == "rarity" {
			return nil
		}
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("read character columns: %w", err)
	}
	
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close character columns: %w", err)
	}

	_, err = d.DB.Exec(`
		ALTER TABLE characters
		ADD COLUMN rarity TEXT NOT NULL DEFAULT 'Common';
	`)

	if err != nil {
		return fmt.Errorf("add rarity column: %w", err)
	}

	return nil
}

func (d *Database) GetCoins(userID string) (int, error) {
	var coins int

	err := d.DB.QueryRow(`
	SELECT coins
	FROM user_balances
	WHERE user_id = ?;
	`, userID).Scan(&coins)

	if err == sql.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return 0, fmt.Errorf("get user coins: %w", err)
	}

	return coins, nil
}

func (d *Database) AddCoins(userID string, amount int) error {
	if amount < 0 {
		return fmt.Errorf("coin amount cannot be negative")
	}

	_, err := d.DB.Exec(`
	INSERT INTO user_balances (user_id, coins)
	VALUES (?, ?)
	ON CONFLICT(user_id) DO UPDATE SET
		coins = coins + excluded.coins;
	`, userID, amount)

	if err != nil {
		return fmt.Errorf("add user coins: %w", err)
	}

	return nil
}

func (d *Database) GetCollectedCharacterIDs(userID string) (map[int]bool, error) {
	rows, err := d.DB.Query(`
	SELECT character_id
	FROM collections
	WHERE user_id = ?
	`, userID)

	if err != nil {
		return nil, fmt.Errorf("get collected character IDs: %w", err)
	}
	defer rows.Close()

	collected := make(map[int]bool)

	for rows.Next() {
		var characterID int

		if err := rows.Scan(&characterID); err != nil {
			return nil, fmt.Errorf("scan collected character ID: %w", err)
		}

		collected[characterID] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read collected character IDs: %w", err)
	}

	return collected, nil
}

func (d *Database) GetDailyShop() ([]models.Character, map[int]int, map[int]bool, error) {
	today := time.Now().Format("2006-01-02")

	var count int
	err := d.DB.QueryRow(`
		SELECT COUNT(*)
		FROM daily_shop
		WHERE shop_date = ?;
	`, today).Scan(&count)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("check daily shop: %w", err)
	}

	if count == 0 {
		rows, err := d.DB.Query(`
			SELECT id, name, image_url, rarity
			FROM characters
			WHERE rarity IN ('Rare', 'Epic', 'Legendary');
		`)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("get shop candidates: %w", err)
		}

		var candidates []models.Character
		for rows.Next() {
			var c models.Character
			if err := rows.Scan(&c.ID, &c.Name, &c.ImageURL, &c.Rarity); err != nil {
				rows.Close()
				return nil, nil, nil, fmt.Errorf("scan shop candidate: %w", err)
			}
			candidates = append(candidates, c)
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, nil, nil, fmt.Errorf("read shop candidates: %w", err)
		}
		rows.Close()

		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})

		if len(candidates) > 5 {
			candidates = candidates[:5]
		}

		tx, err := d.DB.Begin()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("begin shop generation: %w", err)
		}

		for _, c := range candidates {
			price := 0
			switch c.Rarity {
			case "Rare":
				price = 250
			case "Epic":
				price = 500
			case "Legendary":
				price = 1000
			}

			_, err := tx.Exec(`
				INSERT OR IGNORE INTO daily_shop (character_id, shop_date, price)
				VALUES (?, ?, ?);
			`, c.ID, today, price)
			if err != nil {
				tx.Rollback()
				return nil, nil, nil, fmt.Errorf("save shop stock: %w", err)
			}
		}

		if err := tx.Commit(); err != nil {
			return nil, nil, nil, fmt.Errorf("commit shop stock: %w", err)
		}
	}

	rows, err := d.DB.Query(`
		SELECT c.id, c.name, c.image_url, c.rarity, s.price
		FROM daily_shop s
		JOIN characters c ON c.id = s.character_id
		WHERE s.shop_date = ?
		ORDER BY s.price, c.name;
	`, today)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("load daily shop: %w", err)
	}
	defer rows.Close()

	var stock []models.Character
	prices := make(map[int]int)

	for rows.Next() {
		var c models.Character
		var price int
		if err := rows.Scan(&c.ID, &c.Name, &c.ImageURL, &c.Rarity, &price); err != nil {
			return nil, nil, nil, fmt.Errorf("scan daily shop: %w", err)
		}
		stock = append(stock, c)
		prices[c.ID] = price
	}
	if err := rows.Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("read daily shop: %w", err)
	}

	return stock, prices, nil, nil
}