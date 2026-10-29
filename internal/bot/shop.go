package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleShop(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	stock, prices, _, err := b.DB.GetDailyShop()
	if err != nil {
		fmt.Println("Error loading daily shop:", err)
		respondEphemeral(s, i, "Sorry, I couldn't load today's shop")
		return
	}

	if len(stock) == 0 {
		respondEphemeral(
			s,
			i,
			"There are no Rare, Epic or Legendary characters available in today's shop",
		)
		return
	}

	var description strings.Builder

	for _, character := range stock {
		description.WriteString(
			fmt.Sprintf(
				"**%s** - %s - 🪙 %d coins\n",
				character.Name,
				character.Rarity,
				prices[character.ID],
			),
		)
	}

	embed := &discordgo.MessageEmbed{
		Title: "🛒 Daily Character Shop",
		Description: description.String(),
		Color: 0xF1C40F,
	}

	err = s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		},
	)

	if err != nil {
		fmt.Println("Error responding to shop command:", err)
	}
}