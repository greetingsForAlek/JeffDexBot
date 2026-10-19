package main

import (
	"fmt"
	"os"

	"github.com/greetingsForAlek/JeffDexBot/internal/bot"
	"github.com/greetingsForAlek/JeffDexBot/internal/database"
	"github.com/greetingsForAlek/JeffDexBot/internal/game"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_TOKEN")

	if token == "" {
		panic("DISCORD_TOKEN is not set.")
	}

	db, err := database.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Init(); err != nil {
		panic(err)
	}

	fmt.Println("Database Initialized.")

	if err := db.SeedCharacters(); err != nil {
		panic(err)
	}

	fmt.Println("Characters loaded into database.")

	characters, err := db.GetCharacters()
	if err != nil {
		panic(err)
	}

	game.SetCharacters(characters)

	fmt.Printf("Loaded %d characters from SQLite!\n", len(characters))

	b, err := bot.New(token, db)
	if err != nil {
		panic(err)
	}

	err = b.Start()
	if err != nil {
		panic(err)
	}

	defer b.Stop()

	select {}
}