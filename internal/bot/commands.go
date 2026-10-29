package bot

import (
	"fmt"
	"strconv"
	"strings"
	"sort"

	"github.com/bwmarrin/discordgo"
	"github.com/greetingsForAlek/JeffDexBot/internal/game"
	"github.com/greetingsForAlek/JeffDexBot/internal/models"
)

func (b *Bot) registerHandlers() {
	b.Session.AddHandler(b.handleInteraction)
	b.Session.AddHandler(b.handleShopButton)
	b.Session.AddHandler(b.handleMessage)
}

func (b *Bot) handleInteraction(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	if i.Type == discordgo.InteractionMessageComponent {
		b.handleIndexButton(s, i)
		return
	}

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
							Description: "Rarity: **" + character.Rarity + "**",
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
	
	case "collection":
		if i.Member == nil  || i.Member.User == nil {
			respond(s, i, "This command must be used in a server.")
			return
		}

		characters, err := b.DB.GetCollection(i.Member.User.ID)
		if err != nil {
			fmt.Println("Error loading collection: ", err)
			respond(s, i, "Sorry, I couldn't load your collection. Feel free to create a support ticket if needed.")
			return
		}

		if len(characters) == 0 {
			respond(
				s,
				i,
				"You haven't collected any characters yet! Use `/guess` to get started playing 🎮!",
			)
			return
		}

		var description strings.Builder

		for _, character := range characters {
			description.WriteString("• **")
			description.WriteString(character.Name)
			description.WriteString("**\n")
		}

		err = s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Embeds: []*discordgo.MessageEmbed {
						{
							Title: "Your Character Collection",
							Description: description.String(),
						},
					},
				},
			},
		)

		if err != nil {
			fmt.Println("Error responding to collection command:", err)
		}
	case "wallet":
		coins, err := b.DB.GetCoins(i.Member.User.ID)
		if err != nil {
			fmt.Println("error getting coin balance:", err)

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Sorry, something went wrong retrieving your balance. Please create a support ticket if needed.",
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				fmt.Println("error responding to balance command:", err)
			}
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("🪙 You currently have: **%d coins**!", coins),
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("error responding to balance command:", err)
		}
	
	case "index":
		userID := i.Member.User.ID

		characters, err := b.DB.GetCharacters()
		if err != nil {
			fmt.Println("error getting character index:", err)
			respondEphemeral(s, i, "Sorry, something went wrong loading the character index. Please create a support ticket if needed.")
			return
		}

		rarityOrder := map[string]int {
			"Common": 0,
			"Uncommon": 1,
			"Rare": 2,
			"Epic": 3,
			"Legendary": 4,
		}

		sort.Slice(characters, func(i, j int) bool {
			rarityI := rarityOrder[characters[i].Rarity]
			rarityJ := rarityOrder[characters[j].Rarity]

			if rarityI != rarityJ {
				return rarityI < rarityJ
			}

			return strings.ToLower(characters[i].Name) < strings.ToLower(characters[j].Name)
		})

		collected, err := b.DB.GetCollectedCharacterIDs(userID)
		if err != nil {
			fmt.Println("error getting collected character IDs:", err)
			respondEphemeral(s, i, "Sorry, something went wrong loading your collection progress. Please create a support ticket if needed.")
			return
		}

		totalCollected := 0
		for _, character := range characters {
			if collected[character.ID] {
				totalCollected++
			}
		}

		const pageSize = 10
		totalPages := (len(characters) + pageSize - 1) / pageSize

		embed := buildIndexEmbed(characters, collected, 0, totalCollected)

		var components []discordgo.MessageComponent

		if totalPages > 1 {
			components = []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "Previous",
							Style: discordgo.SecondaryButton,
							CustomID: "index_page_0",
							Disabled: true,
						},
						discordgo.Button{
							Label: "Next",
							Style: discordgo.PrimaryButton,
							CustomID: "index_page_1",
						},
					},
				},
			}
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
				Components: components,
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})

		if err != nil {
			fmt.Println("error responding to index command: ", err)
		}
	
	case "shop":
		b.handleShop(s, i)
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

func respondEphemeral(
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
				Flags: discordgo.MessageFlagsEphemeral,
			},
		},
	)

	if err != nil {
		fmt.Println("error sending ephemeral response:", err)
	}
}

func (b *Bot) handleMessage(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	if m.Author == nil || m.Author.Bot {
		return
	}

	hasRound, correct, character, shiny := game.Guess(
		m.ChannelID,
		m.Content,
	)

	if !hasRound || !correct {
		return
	}

	added, err := b.DB.AddToCollection(
		m.Author.ID,
		character.ID,
	)

	if err != nil {
		fmt.Println("error saving collection: ", err)
		_, _ = s.ChannelMessageSend(
			m.ChannelID,
			"Sorry, something went wrong saving your character. Please try again later or feel free to create a support ticket.",
		)
		return
	}

	shinyReward := 0

	if shiny {
		shinyReward = game.ShinyReward(character.Rarity)

		if shinyReward > 0 {
			err := b.DB.AddCoins(m.Author.ID, shinyReward)
			if err != nil {
				fmt.Println("error awarding shiny bonus: ", err)
				_, _ = s.ChannelMessageSend(
					m.ChannelID,
					"⚠️ You collected the character, but something went wrong awarding your shiny bonus. Please create a support ticket.",
				)
				return
			}
		}
	}

	var message string

	if added {
		message = "🎉" + m.Author.Username + " guessed correctly and collected **" + character.Name + "**!"
	} else {
		reward := game.DuplicateReward(character.Rarity)

		if reward > 0 {
			err := b.DB.AddCoins(m.Author.ID, reward)

			if err != nil {
				fmt.Println("error awarding duplicate coins: ", err)
				_, _ = s.ChannelMessageSend(
					m.ChannelID,
					"Your duplicate was detected, but something went wrong awarding your coins. Please create a support ticket if needed.",
				)
				return
			}
		}

		message = "🪙 " + m.Author.Username + " guessed correctly! They already own **" + character.Name + "**, so they earned **" + fmt.Sprint(reward) + " coins** for the duplicate!"
	}

	if shiny {
		message += fmt.Sprintf(
			"\n✨ **SHINY CHARACTER!** You earned an additonal **%d coins**! ✨",
			shinyReward,
		)
	}

	_, err = s.ChannelMessageSend(m.ChannelID, message)
	if err != nil {
		fmt.Println("Error sending winner message: ", err)
	}
}

func buildIndexEmbed(
	characters []models.Character,
	collected map[int]bool,
	page int,
	totalCollected int,
) *discordgo.MessageEmbed {
	const pageSize = 10

	start := page * pageSize
	end := start + pageSize

	if end > len(characters) {
		end = len(characters)
	}

	var description strings.Builder

	for _, character := range characters[start:end] {
		if collected[character.ID] {
			fmt.Fprintf(
				&description,
				"✅ **%s** - %s\n",
				character.Name,
				character.Rarity,
			)
		} else {
			fmt.Fprintf(
				&description,
				"❓ **???** - %s\n",
				character.Rarity,
			)
		}
	}

	totalPages := (len(characters) + pageSize - 1) / pageSize

	return &discordgo.MessageEmbed{
		Title: "📖 Character Index",
		Description: fmt.Sprintf(
			"Collected **%d/%d** characters\n\n%s",
			totalCollected,
			len(characters),
			description.String(),
		),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page %d/%d", page+1, totalPages),
		},
		Color: 0x5865F2,
	}
}

func (b *Bot) handleIndexButton(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
    	Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	if err != nil {
    	fmt.Println("Error acknowledging index button:", err)
    	return
	}

	data := i.MessageComponentData()

	if !strings.HasPrefix(data.CustomID, "index_page_") {
		return
	}

	page, err := strconv.Atoi(strings.TrimPrefix(data.CustomID, "index_page_"))
	if err != nil {
		respondComponentError(s, i, "Invalid index page.")
		return
	}

	if i.Member == nil || i.Member.User == nil {
		respondComponentError(s, i, "The index can only be used in a server")
		return
	}

	characters, err := b.DB.GetCharacters()
	if err != nil {
		fmt.Println("Error loading character index:", err)
		respondComponentError(s, i, "Sorry, I couldn't load the character index.")
		return
	}

	collected, err := b.DB.GetCollectedCharacterIDs(i.Member.User.ID)
	if err != nil {
		fmt.Println("Error loading collected character IDs:", err)
		respondComponentError(s, i, "Sorry, I couldn't load your collection progress.")
		return
	}

	totalCollected := 0
	for _, character := range characters {
		if collected[character.ID] {
			totalCollected++
		}
	}

	const pageSize = 10
	totalPages := (len(characters) + pageSize - 1) / pageSize

	if totalPages == 0 {
		respondComponentError(s, i, "There are no characters in the index.")
		return
	}

	if page >= totalPages {
		page = totalPages - 1
	}

	var components []discordgo.MessageComponent

	if totalPages > 1 {
		components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Previous",
						Style: discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("index_page_%d", page-1),
						Disabled: page == 0,
					},
					discordgo.Button{
						Label: "Next",
						Style: discordgo.PrimaryButton,
						CustomID: fmt.Sprintf("index_page_%d", page+1),
						Disabled: page == totalPages - 1,
					},
				},
			},
		}
	}

	_, err = s.InteractionResponseEdit(
    	i.Interaction,
    	&discordgo.WebhookEdit{
        	Embeds: &[]*discordgo.MessageEmbed{
            	buildIndexEmbed(characters, collected, page, totalCollected),
        	},
        	Components: &components,
    	},
	)

	if err != nil {
    	fmt.Println("Error updating character index:", err)
	}
}

func respondComponentError(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	message string,
) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		fmt.Println("Error responding to index button:", err)
	}
}