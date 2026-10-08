package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) registerHandlers() {
	b.Session.AddHandler(handleInteraction)
}

func handleInteraction(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	if i.ApplicationCommandData().Name != "ping" {
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "pong",
		},
	})

	if err != nil {
		fmt.Println("Error responding to interaction:", err)
	}
}