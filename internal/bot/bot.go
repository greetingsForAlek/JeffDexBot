package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
}

func New(token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	b := &Bot {
		Session: session,
	}

	b.registerHandlers()

	return b, nil
}

func (b *Bot) Start() error {
	err := b.Session.Open()
	if err != nil {
		return err
	}

	commands := []*discordgo.ApplicationCommand {
		{
			Name: "ping",
			Description: "Replies with pong.",
		},
		{
			Name: "guess",
			Description: "Guess the character shown in the image to add it to your collection!",
		},
	}

	for _, command := range commands {
		_, err := b.Session.ApplicationCommandCreate(
			b.Session.State.User.ID,
			"",
			command,
		)

		if err != nil {
			return err
		}

		fmt.Println("Registered /" + command.Name)
	}

	fmt.Println("Bot is online!")
	fmt.Println("Registered /ping!")

	return nil
}

func (b *Bot) Stop() error {
	return b.Session.Close()
}