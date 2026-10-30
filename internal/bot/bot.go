package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/greetingsForAlek/JeffDexBot/internal/database"
)

// bot struct
type Bot struct {
	Session *discordgo.Session // discordgo session
	DB *database.Database // database so the bot can load stuff
	DropChannelID string // drop channel for daily drops
	dropSchedulerStop chan struct{} // dropSchedulerStop for daily drops
}

// new bot function
func New(token string, db *database.Database, dropChannelID string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token) // new discordgo sessio 
	if err != nil {
		return nil, err
	}

	session.Identify.Intents |= discordgo.IntentGuildMessages |
		discordgo.IntentMessageContent // get the intents

	// set up the bot
	b := &Bot {
		Session: session,
		DB: db,
		DropChannelID: dropChannelID,
	}

	b.registerHandlers() // register handlers

	return b, nil
}

// Start the bot
func (b *Bot) Start() error {
	// Check for the dropChannelID
	if b.DropChannelID == "" {
		return fmt.Errorf("DROP_CHANNEL_ID is not set") // error if not found
	}

	// Open a new session
	err := b.Session.Open()
	if err != nil {
		return err
	}

	// commands
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
		{
			Name: "index",
			Description: "View your character collection progress",
		},
		{
			Name: "shop",
			Description: "View today's character shop",
		},
	}

	// filter through the commands
	for _, command := range commands {
		// create every command
		_, err := b.Session.ApplicationCommandCreate(
			b.Session.State.User.ID,
			"",
			command,
		)

		if err != nil {
			return err // return error if there is one
		}

		fmt.Println("Registered /" + command.Name) // print if registry succeeds
	}

	// register the commands
	registered, err := b.Session.ApplicationCommands(
    	b.Session.State.User.ID,
    	"",
	)
	if err != nil {
    	return fmt.Errorf("list registered commands: %w", err) // error if there is one
	}

	// loop through and print if the registry is ok
	for _, command := range registered {
    	fmt.Println("Discord command:", command.Name)
	}

	fmt.Println("Bot is online!") // print message to confirm the bot is online

	b.dropSchedulerStop = b.startDropScheduler() // start the drop scheduler

	return nil // return nil because there are no further errors
}

// stop function
func (b *Bot) Stop() error {
	if b.dropSchedulerStop != nil { // close if there is no dropSchedulerStop
		close(b.dropSchedulerStop)
	}

	return b.Session.Close()
}