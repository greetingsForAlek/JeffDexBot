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

	characters, err := game.LoadCharacters("data/characters.json")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded characters from JSON!\n")

	db, err := database.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Init(); err != nil {
		panic(err)
	}

	fmt.Println("Database Initialized.")

	if err := db.SeedCharacters(characters); err != nil {
		panic(err)
	}

	fmt.Println("Characters loaded into database.")

	game.SetCharacters(characters)

	dropChannelID := os.Getenv("DROP_CHANNEL_ID")

	if dropChannelID == "" {
		panic("DROP_CHANNEL_ID is not set.")
	}

	b, err := bot.New(token, db, dropChannelID)
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