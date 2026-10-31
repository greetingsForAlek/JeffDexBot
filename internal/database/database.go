// this file contains all of the database stuff
// be warned, there are a lot of SQL queries.

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

// database struct
type Database struct {
	DB *sql.DB // sql database
}

// new database function
func New() (*Database, error) {
	dbPath := os.Getenv("DATABASE_PATH") // get the database path from the .env
	// if there is no provided db path, we're probably in testing and not deployment
	// so we need to set it to this path
	if dbPath == "" {
		dbPath = "./character-bot.db"
	}

	// check if there's errors while making the dbPath
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	// open the dbPath with sqlite
	db, err := sql.Open("sqlite", dbPath)
	// error if needed
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// check for errors while pinging the database
	if err := db.Ping(); err != nil {
		db.Close() // close the db
		return nil, fmt.Errorf("connect to database: %w", err) // error message
	}

	return &Database{DB: db}, nil // return the new database and nil, which means there are no further errors
}

// close database function
func (d *Database) Close() error {
	// literally just return d.DB.Close()
	return d.DB.Close()
}

// init function
func (d *Database) Init() error {
	// query the database with this big sql query
	// the query will create all of the tables in the database:
	// characters, collections, user_balances, daily_shop & daily_shop_purchases
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

	// now we need to execute the query
	_, err := d.DB.Exec(query)
	// and error if there was any execution errors
	if err != nil {
		return fmt.Errorf("create database tables: %w", err)
	}

	// also need to ensure that the rarity column exists
	if err := d.ensureRarityColumn(); err != nil {
		return err
	}

	return nil // return nil because there are no further errors
}

// add character function
func (d *Database) AddCharacter (
	id int,
	name string,
	imageURL string,
	rarity string,
) error {
	// first we write the query
	// it will basically update the characters table to add the new character
	query := `
	INSERT INTO characters (id, name, image_url, rarity)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		image_url = excluded.image_url,
		rarity = excluded.rarity;
	`

	// execute the query and catch any errors
	_, err := d.DB.Exec(query, id, name, imageURL, rarity)
	if err != nil {
		return fmt.Errorf("add character: %w", err)
	}

	return nil // return nil, as there are no further errors
}

//  get characters function
func (d *Database) GetCharacters() ([]models.Character, error) {
	// write the query
	// this will select the character by using the provided details
	query := `
	SELECT id, name, image_url, rarity
	FROM characters
	ORDER BY id;
	`

	// execute the query and catch any errors
	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query characters: %w", err)
	}
	defer rows.Close()

	// characters slice
	var characters []models.Character

	// loop through the rows
	for rows.Next() {
		var character models.Character // character variable

		// scan the rows and catch any errors
		err := rows.Scan(
			&character.ID,
			&character.Name,
			&character.ImageURL,
			&character.Rarity,
		)
		if err != nil {
			return nil, fmt.Errorf("scan character: %w", err)
		}

		characters = append(characters, character) // append the new character to the characters slice
	}

	// error if there is one
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate characters: %w", err)
	}

	return characters, nil // return the characters slice and nil indicating no errors
}

// add to collection function
func (d *Database) AddToCollection(
	userID string,
	characterID int,
) (bool, error) {
	// write the query
	// this will add a character's id to the user's collection table
	query := `
	INSERT INTO collections (user_id, character_id)
	VALUES (?, ?)
	ON CONFLICT(user_id, character_id) DO NOTHING;
	`

	// execute the query and catch errors
	result, err := d.DB.Exec(query, userID, characterID)
	if err != nil {
		return false, fmt.Errorf("add to collection: %w", err)
	}

	// get the affected rows and catch any errors
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("check collection insert: %w", err)
	}

	return rowsAffected > 0, nil // return a boolean which is true if there are 1 or more rows affected and also return nil as there are no errors
}

// Get collection function
func (d *Database) GetCollection (
	userID string,
) ([]models.Character, error) {
	// write the query
	// this will fetch the user's collection
	query := `
		SELECT c.id, c.name, c.image_url
		FROM characters c
		INNER JOIN collections col
			ON col.character_id = c.id
		WHERE col.user_id  = ?
		ORDER BY col.collected_at, c.id;
	`

	// execute the query and catch any errors
	rows, err := d.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("get collection: %w", err)
	}
	defer rows.Close()

	var characters []models.Character // characters slice

	// loop through rows
	for rows.Next() {
		var character models.Character // character variable

		// scan the rows and catch any errors
		if err := rows.Scan(
			&character.ID,
			&character.Name,
			&character.ImageURL,
		); err != nil {
			return nil, fmt.Errorf("scan collection character: %w", err) // return any errors
		}

		characters = append(characters, character) // append the character to the slice
	}

	// throw any errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate collection: %w", err)
	}

	return characters, nil // return the slice and nil meaning there are no errors
}

// ensure rarity column function
// the whole purpose of this is to ensure that the character's rarity column exists
func (d *Database) ensureRarityColumn() error {
	// execute a small query and return any errors
	rows, err := d.DB.Query("PRAGMA table_info(characters)")
	if err != nil {
		return fmt.Errorf("inspect characters table: %w", err)
	}
	defer rows.Close()

	// loop through the rows
	for rows.Next() {
		// define all of these variables
		var (
			cid int
			name string
			columnType string
			notNull int
			defaultValue sql.NullString
			primaryKey int
		)

		// scan the rows to set the variables and return any errors
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

		// check if rarity shows up and return nil
		if name == "rarity" {
			return nil
		}
	}

	// error if needed
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("read character columns: %w", err)
	}
	
	// more errors if needed
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close character columns: %w", err)
	}

	// execute a small query to alter the characters table
	// this will make sure every character has a rarity
	_, err = d.DB.Exec(`
		ALTER TABLE characters
		ADD COLUMN rarity TEXT NOT NULL DEFAULT 'Common';
	`)

	// error if there is any
	if err != nil {
		return fmt.Errorf("add rarity column: %w", err)
	}

	return nil // return nil because there are no more errors
}

// get coins function
func (d *Database) GetCoins(userID string) (int, error) {
	var coins int // coins

	// query the rows to get the coins from the user balance
	err := d.DB.QueryRow(`
	SELECT coins
	FROM user_balances
	WHERE user_id = ?;
	`, userID).Scan(&coins) // scan to set the coins

	// error if there are no rows
	if err == sql.ErrNoRows {
		return 0, nil
	}

	// error if anything else went wrong
	if err != nil {
		return 0, fmt.Errorf("get user coins: %w", err)
	}

	return coins, nil // return the coins and nil because no errors
}

// add coins function
func (d *Database) AddCoins(userID string, amount int) error {
	// check if provided amount is less than zero and throw an error
	if amount < 0 {
		return fmt.Errorf("coin amount cannot be negative")
	}

	// execute a query
	// this will insert coins into the user's balance
	_, err := d.DB.Exec(`
	INSERT INTO user_balances (user_id, coins)
	VALUES (?, ?)
	ON CONFLICT(user_id) DO UPDATE SET
		coins = coins + excluded.coins;
	`, userID, amount)

	// error if there is one
	if err != nil {
		return fmt.Errorf("add user coins: %w", err)
	}

	return nil // return nil for no errors
}

// get collected character ids function
// basically gets the ids of all the characters that user has collected
func (d *Database) GetCollectedCharacterIDs(userID string) (map[int]bool, error) {
	// write a query
	// this will select all of the character ids in a users collection
	rows, err := d.DB.Query(`
	SELECT character_id
	FROM collections
	WHERE user_id = ?
	`, userID)

	// error if needed
	if err != nil {
		return nil, fmt.Errorf("get collected character IDs: %w", err)
	}
	defer rows.Close()

	collected := make(map[int]bool) // collected variable

	// loop through the rows
	for rows.Next() {
		var characterID int // character id

		// catch errors while scanning the rows
		if err := rows.Scan(&characterID); err != nil {
			return nil, fmt.Errorf("scan collected character ID: %w", err)
		}

		collected[characterID] = true // set the characterID's collected status to true
	}

	// error if needed
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read collected character IDs: %w", err)
	}

	return collected, nil // return
}

// get daily shop function
func (d *Database) GetDailyShop() ([]models.Character, map[int]int, map[int]bool, error) {
	today := time.Now().Format("2006-01-02") // today (current time)

	var count int // count variable
	// query rows and catch errors
	err := d.DB.QueryRow(`
		SELECT COUNT(*)
		FROM daily_shop
		WHERE shop_date = ?;
	`, today).Scan(&count)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("check daily shop: %w", err)
	}

	// check if count is 0
	if count == 0 {
		// do another query and catch errors
		rows, err := d.DB.Query(`
			SELECT id, name, image_url, rarity
			FROM characters
			WHERE rarity IN ('Rare', 'Epic', 'Legendary');
		`)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("get shop candidates: %w", err)
		}

		// candidates slice
		var candidates []models.Character
		// loop through rows
		for rows.Next() {
			var c models.Character // character variable
			if err := rows.Scan(&c.ID, &c.Name, &c.ImageURL, &c.Rarity); err != nil { // scan the rows and catch errors
				rows.Close()
				return nil, nil, nil, fmt.Errorf("scan shop candidate: %w", err)
			}
			// append the character to the candidates slice
			candidates = append(candidates, c)
		}
		// error if there is one
		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, nil, nil, fmt.Errorf("read shop candidates: %w", err)
		}
		rows.Close()

		// select random candidates
		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})

		// if there are over five candidates, only pick five
		if len(candidates) > 5 {
			candidates = candidates[:5]
		}

		// begin the database and catch errors
		tx, err := d.DB.Begin()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("begin shop generation: %w", err)
		}

		// loop through the candidates and select prices
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

			// execute a query and catch errors
			_, err := tx.Exec(`
				INSERT OR IGNORE INTO daily_shop (character_id, shop_date, price)
				VALUES (?, ?, ?);
			`, c.ID, today, price)
			if err != nil {
				tx.Rollback()
				return nil, nil, nil, fmt.Errorf("save shop stock: %w", err)
			}
		}

		// commit and catch errors
		if err := tx.Commit(); err != nil {
			return nil, nil, nil, fmt.Errorf("commit shop stock: %w", err)
		}
	}

	// do another query
	rows, err := d.DB.Query(`
		SELECT c.id, c.name, c.image_url, c.rarity, s.price
		FROM daily_shop s
		JOIN characters c ON c.id = s.character_id
		WHERE s.shop_date = ?
		ORDER BY s.price, c.name;
	`, today)
	// catch any errors
	if err != nil {
		return nil, nil, nil, fmt.Errorf("load daily shop: %w", err)
	}
	defer rows.Close()

	var stock []models.Character // stock slice
	prices := make(map[int]int) // prices map

	// loop through rows
	for rows.Next() {
		// character and price variables
		var c models.Character
		var price int
		// scan the rows and catch errors
		if err := rows.Scan(&c.ID, &c.Name, &c.ImageURL, &c.Rarity, &price); err != nil {
			return nil, nil, nil, fmt.Errorf("scan daily shop: %w", err)
		}
		stock = append(stock, c) // append the character to the stock
		prices[c.ID] = price // set the character price
	}
	// catch errors
	if err := rows.Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("read daily shop: %w", err)
	}

	return stock, prices, nil, nil // return
}

// buy daily shop character function
func (d *Database) BuyDailyShopCharacter(
	userID string,
	characterID int,
) (string, int, error) {
	today := time.Now().Format("2006-01-02") // get the current time

	// begin the db; catch errors
	tx, err := d.DB.Begin()
	if err != nil {
		return "", 0, fmt.Errorf("begin purchase: %w", err)
	}
	defer tx.Rollback()

	// get the price
	var price int
	// do a query to select the character's price
	err = tx.QueryRow(`
		SELECT price
		FROM daily_shop
		WHERE character_id = ? AND shop_date = ?;
	`, characterID, today).Scan(&price)

	// no rows error
	if err == sql.ErrNoRows {
		return "not_available", 0, nil
	}

	// other errors
	if err != nil {
		return "", 0, fmt.Errorf("find shop item: %w", err)
	}

	// do another query
	// this will insert the character_id into daily_shop_purchases
	result, err := tx.Exec(`
		INSERT INTO daily_shop_purchases (user_id, character_id, shop_date)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, character_id, shop_date) DO NOTHING;
	`, userID, characterID, today)
	// catch errors
	if err != nil {
		return "", 0, fmt.Errorf("record shop purchase: %w", err)
	}

	// get the affected rows and catch errors
	inserted, err := result.RowsAffected()
	if err != nil {
		return "", 0, fmt.Errorf("check purchase record: %w", err)
	}
	// error if inserted is 0
	if inserted == 0 {
		return "already_bought", 0, nil
	}

	// execute yet another query to update user balances
	// also catch any errors
	result, err = tx.Exec(`
		UPDATE user_balances
		SET coins = coins - ?
		WHERE user_id = ? AND coins >= ?;
	`, price, userID, price)
	if err != nil {
		return "", 0, fmt.Errorf("deduct purchase price: %w", err)
	}

	// get the rows affected from the last query and check for errors
	updated, err := result.RowsAffected()
	if err != nil {
		return "", 0, fmt.Errorf("check balance deduction: %w", err)
	}

	// if updated is 0 throw an error
	if updated == 0 {
		return "insufficient_funds", 0, nil
	}

	// execute ANOTHER QUERY and check for MORE ERRORS
	_, err = tx.Exec(`
		INSERT INTO collections (user_id, character_id)
		VALUES (?, ?)
		ON CONFLICT (user_id, character_id) DO NOTHING;
	`, userID, characterID)

	if err != nil {
		return "", 0, fmt.Errorf("add purchased character: %w", err)
	}

	// commit and check for errors
	if err := tx.Commit(); err != nil {
		return "", 0, fmt.Errorf("commit purchase: %w", err)
	}

	return "success", price, nil // and finally return
}