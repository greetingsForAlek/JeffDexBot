package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/greetingsForAlek/JeffDexBot/internal/database"
)

type Bot struct {
	Session *discordgo.Session
	DB *database.Database
}

func New(token string, db *database.Database) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	session.Identify.Intents |= discordgo.IntentGuildMessages |
		discordgo.IntentMessageContent

	b := &Bot {
		Session: session,
		DB: db,
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
		{
			Name: "collection",
			Description: "View your collected characters!",
		},
		{
			Name: "wallet",
			Description: "Check your coins wallet!",
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

	registered, err := b.Session.ApplicationCommands(
    	b.Session.State.User.ID,
    	"",
	)
	if err != nil {
    	return fmt.Errorf("list registered commands: %w", err)
	}

	for _, command := range registered {
    	fmt.Println("Discord command:", command.Name)
	}

	fmt.Println("Bot is online!")

	return nil
}

func (b *Bot) Stop() error {
	return b.Session.Close()
}