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

	bot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ApplicationCommandData().Name == "ping" {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Content: "pong",
				},
			})

			if err != nil {
				fmt.Println("Error responding to interaction:", err)
			}
		}
	})

	err = bot.Open()
	if err != nil {
		panic(err)
	}

	command := &discordgo.ApplicationCommand {
		Name: "ping",
		Description: "Replies with pong.",
	}

	_, err = bot.ApplicationCommandCreate(bot.State.User.ID, "", command)
	if err != nil {
		panic(err)
	}

	fmt.Println("Bot is online!")

	select {}
}