package bot

import (
	"fmt"
	"math/rand/v2"

	"github.com/bwmarrin/discordgo"
	"github.com/greetingsForAlek/JeffDexBot/internal/game"
)

var dropMessages = []string {
	"who is this rando?",
	"wait a minute, who *ARE* you?",
	"suspect has entered the room.",
	"now's your chance.",
	"alright, who invited [name]?",
	"wait a sec i think i know this guy",
	"hello, happy halloween",
	"here i got you a gift",
	"don't tell anyone, but this is... uh hold on i forgot",
	"a wild [name] has appeared!",
	"i wonder if this one is shiny?",
	"did you know that these characters have a chance of being shiny?",
	"duplicates and shinies award coins. the more you know.",
	"it's good to check your coins balance now and again. run `/wallet` to do so.",
	"want to see your collected characters? run `/index`!",
	"you know you don't have to wait for these drops? run `/guess` any time you want!",
	"what's in the shop today?",
	"only characters that are rare and above can appear in the shop.",
	"well, well, well if it isn't...",
	"look who finally decided to show up.",
}

func (b *Bot) sendDrop(s *discordgo.Session, channelID string) {
	if game.HasActiveRound(channelID) {
		_, err := s.ChannelMessageSend(
			channelID,
			"The previous character escaped! A new drop is arriving...",
		)
		if err != nil {
			fmt.Println("Error announcing expired character:", err)
		}
	}

	character := game.StartRound(channelID)

	message := dropMessages[rand.IntN(len(dropMessages))]

	_, err := s.ChannelMessageSendComplex(
		channelID,
		&discordgo.MessageSend{
			Content: message,
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
	)

	if err != nil {
		fmt.Println("Error sending character drop:", err)
	}
}