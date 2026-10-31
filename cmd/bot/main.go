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
	// load the env file
	_ = godotenv.Load()

	// get the discord token
	token := os.Getenv("DISCORD_TOKEN")

	if token == "" {
		panic("DISCORD_TOKEN is not set.") // error if discord token is not set
	}

	// load the characters from the json file
	characters, err := game.LoadCharacters("data/characters.json")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded characters from JSON!\n")

	// New Database
	db, err := database.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Init(); err != nil { // initialise the database
		panic(err)
	}

	fmt.Println("Database Initialized.")

	// seed the characters
	if err := db.SeedCharacters(characters); err != nil {
		panic(err)
	}

	fmt.Println("Characters loaded into database.")

	// set characters
	game.SetCharacters(characters)

	dropChannelID := os.Getenv("DROP_CHANNEL_ID") // get the drop channel id for daily drops

	if dropChannelID == "" {
		panic("DROP_CHANNEL_ID is not set.") // error if the drop channel id is not set
	}

	// New bot
	b, err := bot.New(token, db, dropChannelID)
	if err != nil {
		panic(err)
	}

	// Start the bot
	err = b.Start()
	if err != nil {
		panic(err)
	}

	defer b.Stop()

	select {}
}