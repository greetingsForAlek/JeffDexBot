package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)


// handle shop function
func (b *Bot) handleShop(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	stock, prices, _, err := b.DB.GetDailyShop() // get the daily shop's stocks and prices from the database
	// error if required
	if err != nil {
		fmt.Println("Error loading daily shop:", err)
		respondEphemeral(s, i, "Sorry, I couldn't load today's shop")
		return
	}

	// if there are no characters in stock
	if len(stock) == 0 {
		// ephemeral response to let the user know
		respondEphemeral(
			s,
			i,
			"There are no Rare, Epic or Legendary characters available in today's shop",
		)
		return
	}

	// description and components variables
	var (
		description strings.Builder
		components []discordgo.MessageComponent
	)

	// loop through character in stock
	for _, character := range stock {
		// build the string for each of them
		description.WriteString(
			fmt.Sprintf(
				"**%s** - %s - 🪙 %d coins\n",
				character.Name,
				character.Rarity,
				prices[character.ID],
			),
		)

		// append a button for each character
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

	// build the embed
	embed := &discordgo.MessageEmbed{
		Title: "🛒 Daily Character Shop",
		Description: description.String(),
		Color: 0xF1C40F,
	}

	// respond to the interaction
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

	// error if needed
	if err != nil {
		fmt.Println("Error responding to shop command:", err)
	}
}

// handle shop button
func (b *Bot) handleShopButton(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	// check the type of the component
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	// get the customID
	customID := i.MessageComponentData().CustomID
	// check if it has the desired prefix
	if !strings.HasPrefix(customID, "shop_buy_") {
		return
	}

	// extract the character id from the custom id
	characterID, err := strconv.Atoi(
		strings.TrimPrefix(customID, "shop_buy_"),
	)
	// error if required to
	if err != nil {
		respondEphemeral(s, i, "Invalid shop item.")
		return
	}

	// get the user id
	userID := i.Member.User.ID
	// set userID if no user
	if i.User != nil {
		userID = i.User.ID
	}

	// buy daily shop character
	status, price, err := b.DB.BuyDailyShopCharacter(
		userID,
		characterID,
	)
	// error if needed
	if err != nil {
		fmt.Println("Error purchasing shop character:", err)
		respondEphemeral(s, i, "Something went wrong with your purchase. Please try again later.")
		return
	}

	// switch the status
	switch status {
	case "success": // success
		respondEphemeral( // success response
			s,
			i,
			fmt.Sprintf("🎉 Purchase successful! You bought the character for 🪙 %d coins.", price),
		)
	case "not_available": // not available
		respondEphemeral(s, i, "That character isn't available in today's shop.") // not available response
	case "already_bought": // already bought
		respondEphemeral(s, i, "You've already purchased that shop listing today.") // already bought response
	case "insufficient_funds": // insufficient funds
		respondEphemeral(s, i, "You don't have enough coins to buy that character!") // insufficient funds response
	default: // default
		respondEphemeral(s, i, "Your purchase couldn't be completed.") // default response
	}
}
