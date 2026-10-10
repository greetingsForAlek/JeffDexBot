package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

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
		image_url TEXT NOT NULL
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
	`

	_, err := d.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("create database tables: %w", err)
	}

	return nil
}