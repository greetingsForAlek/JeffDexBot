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
		message = "🎉" + m.Author.Username + " guessed correctly and collected **" + character.Name + "**!"
	} else {
		message = "🎉" + m.Author.Username + " guessed correctly! They already own **" + character.Name + "**, so no duplicate was added."
	}

	_, err = s.ChannelMessageSend(m.ChannelID, message)
	if err != nil {
		fmt.Println("Error sending winner message: ", err)
	}
}