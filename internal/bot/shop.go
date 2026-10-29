package bot

import (
	"fmt"
	"strconv"
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

	var (
		description strings.Builder
		components []discordgo.MessageComponent
	)

	for _, character := range stock {
		description.WriteString(
			fmt.Sprintf(
				"**%s** - %s - 🪙 %d coins\n",
				character.Name,
				character.Rarity,
				prices[character.ID],
			),
		)

		components = append(
			components, 
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Buy " + character.Name,
						Style: discordgo.PrimaryButton,
						CustomID: fmt.Sprintf("shop_buy_%d", character.ID),
					},
				},
			},
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
				Components: components,
			},
		},
	)

	if err != nil {
		fmt.Println("Error responding to shop command:", err)
	}
}

func (b *Bot) handleShopButton(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	customID := i.MessageComponentData().CustomID
	if !strings.HasPrefix(customID, "shop_buy_") {
		return
	}

	characterID, err := strconv.Atoi(
		strings.TrimPrefix(customID, "shop_buy_"),
	)
	if err != nil {
		respondEphemeral(s, i, "Invalid shop item.")
		return
	}

	userID := i.Member.User.ID
	if i.User != nil {
		userID = i.User.ID
	}

	status, price, err := b.DB.BuyDailyShopCharacter(
		userID,
		characterID,
	)
	if err != nil {
		fmt.Println("Error purchasing shop character:", err)
		respondEphemeral(s, i, "Something went wrong with your purchase. Please try again later.")
		return
	}

	switch status {
	case "success":
		respondEphemeral(
			s,
			i,
			fmt.Sprintf("🎉 Purchase successful! You bought the character for 🪙 %d coins.", price),
		)
	case "not_available":
		respondEphemeral(s, i, "That character isn't available in today's shop.")
	case "already_bought":
		respondEphemeral(s, i, "You've already purchased that shop listing today.")
	case "insufficient_funds":
		respondEphemeral(s, i, "You don't have enough coins to buy that character!")
	default:
		respondEphemeral(s, i, "Your purchase couldn't be completed.")
	}
}
