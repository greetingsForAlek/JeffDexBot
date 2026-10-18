package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/greetingsForAlek/JeffDexBot/internal/game"
)

func (b *Bot) registerHandlers() {
	b.Session.AddHandler(b.handleInteraction)
	b.Session.AddHandler(b.handleMessage)
}

func (b *Bot) handleInteraction(
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

		collected, err := b.DB.GetCollectedCharacterIDs(userID)
		if err != nil {
			fmt.Println("error getting collected character IDs:", err)
			respondEphemeral(s, i, "Sorry, something went wrong loading your collection progress. Please create a support ticket if needed.")
			return
		}

		var description strings.Builder
		totalCollected := 0

		for _, character := range characters {
			if collected[character.ID] {
				fmt.Fprintf(
					&description,
					"✅**%s** - %s\n",
					character.Name,
					character.Rarity,
				)
				totalCollected++
			} else {
				fmt.Fprintf(
					&description,
					"❓ **???** - %s\n",
					character.Rarity,
				)
			}
		}

		embed := &discordgo.MessageEmbed {
			Title: "📖 Character Index",
			Description: fmt.Sprintf(
				"Collected **%d/%d** characters\n\n%s",
				totalCollected,
				len(characters),
				description.String(),
			),
			Color: 0x5865F2,
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("error responding to index command:", err)
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

	hasRound, correct, character := game.Guess(
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

	var message string

	if added {
		message = "🎉" + m.Author.Username + "guessed correctly and collected **" + character.Name + "**!"
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

	_, err = s.ChannelMessageSend(m.ChannelID, message)
	if err != nil {
		fmt.Println("Error sending winner message: ", err)
	}
}