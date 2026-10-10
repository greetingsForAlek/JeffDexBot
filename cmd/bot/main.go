package main

import (
	"os"
	"fmt"
	
	"github.com/joho/godotenv"
	"github.com/greetingsForAlek/JeffDexBot/internal/bot"
	"github.com/greetingsForAlek/JeffDexBot/internal/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file.")
	}

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

	b, err := bot.New(token)
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