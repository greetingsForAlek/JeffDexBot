package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/greetingsForAlek/JeffDexBot/internal/game"
)

func (b *Bot) registerHandlers() {
	b.Session.AddHandler(handleInteraction)
	b.Session.AddHandler(handleMessage)
}

func handleInteraction(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	switch i.ApplicationCommandData().Name {
	case "ping":
		respond(s, i, "pong")

	case "guess":
		character := game.StartRound(i.ChannelID)

		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed {
						{
							Title: "Who is this character?",
							Description: "Take a guess!",
							Image: &discordgo.MessageEmbedImage{
								URL: character.ImageURL,
							},
						},
					},
				},
			},
		)

		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
	}
}

func respond(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	message string,
) {
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		},
	)

	if err != nil {
		fmt.Println("Error responding to interaction: ", err)
	}
}

func handleMessage(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	if m.Author == nil || m.Author.Bot {
		return
	}

	hasRound, correct, character := game.Guess(
		m.ChannelID,
		m.Content,
	)

	if !hasRound || !correct {
		return
	}

	added := game.AddToCollection(m.Author.ID, character)

	var message string

	if added {
		message = "🎉 " + m.Author.Username + " guessed correctly and collected ** " + character.Name + "**!"
	} else {
		message = "🎉 " + m.Author.Username + " guessed correctly! But they already own **" + character.Name + "**, so no duplicate was added."
	}

	_, err := s.ChannelMessageSend(m.ChannelID, message)
	if err != nil {
		fmt.Println("Error sending winner message:", err)
	}
}