package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/greetingsForAlek/JeffDexBot/internal/game"
)

func (b *Bot) registerHandlers() {
	b.Session.AddHandler(handleInteraction)
}

func handleInteraction(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	switch i.ApplicationCommandData().Name {
	case "ping":
		respond(s, i, "pong")

	case "guess":
		character := game.Random()

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