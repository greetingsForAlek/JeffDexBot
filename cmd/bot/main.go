package main

import (
	"os"
	
	"github.com/joho/godotenv"
	"github.com/greetingsForAlek/JeffDexBot/internal/bot"
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