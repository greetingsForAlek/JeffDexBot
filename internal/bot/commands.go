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

	hasRound, correct, answer := game.Guess(
		m.ChannelID,
		m.Content,
	)

	if !hasRound {
		return
	}

	if correct {
		_, err := s.ChannelMessageSend(
			m.ChannelID,
			"🎉" + m.Author.Username + " guessed correctly! The character was **" + answer + "**",
		)

		if err != nil {
			fmt.Println("Error sending winner message: ", err)
		}
	}
}