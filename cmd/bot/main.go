package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file!")
	}
	
	token := os.Getenv("DISCORD_TOKEN")

	if token == "" {
		panic("Discord Token is not set!")
	}

	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	err = bot.Open()
	if err != nil {
		panic(err)
	}

	fmt.Println("Bot is online!")

	select {}
}