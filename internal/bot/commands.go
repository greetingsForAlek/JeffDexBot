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

// register handlers
func (b *Bot) registerHandlers() {
	b.Session.AddHandler(b.handleInteraction) // interaction handler
	b.Session.AddHandler(b.handleShopButton) // shop button handler
	b.Session.AddHandler(b.handleWhoIsThisButton) // "who is this" button handler
	b.Session.AddHandler(b.handleMessage) // handle message
}

// handle interaction function
func (b *Bot) handleInteraction(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	// interaction message component
	if i.Type == discordgo.InteractionMessageComponent {
		b.handleIndexButton(s, i)
		return
	}

	// switch through the application command data
	switch i.ApplicationCommandData().Name {
	case "ping": // ping command (testing)
		respond(s, i, "pong")

	case "guess": // guess command
		character := game.StartRound(i.ChannelID) // start round

		// components
		components := []discordgo.MessageComponent {
			// action row
			discordgo.ActionsRow{ 
				Components: []discordgo.MessageComponent {
					// button - "who is this" button
					// reveals the character's name
					discordgo.Button {
						Label: "Who is this?",
						Style: discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("who_is_this_%s", character.Name),
					},
				},
			},
		}

		// interaction respond - the guess message
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
					Components: components,
				},
			},
		)

		// check for errors
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
	
	case "collection": // collection command
		if i.Member == nil  || i.Member.User == nil { // check if command is being used outside of a server
			respond(s, i, "This command must be used in a server.")
			return
		}

		// get the collection from the database
		characters, err := b.DB.GetCollection(i.Member.User.ID)
		if err != nil {
			fmt.Println("Error loading collection: ", err)
			respond(s, i, "Sorry, I couldn't load your collection. Feel free to create a support ticket if needed.")
			return
		}

		// check the length of the characters list
		if len(characters) == 0 {
			respond(
				s,
				i,
				"You haven't collected any characters yet! Use `/guess` to get started playing 🎮!",
			) // respond if there is no characters in the list
			return
		}

		var description strings.Builder // string builder

		// build the message
		for _, character := range characters {
			description.WriteString("• **")
			description.WriteString(character.Name)
			description.WriteString("**\n")
		}

		// respond to the interaction
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

		// print errors if there is any
		if err != nil {
			fmt.Println("Error responding to collection command:", err)
		}
	case "wallet": // wallet command
		coins, err := b.DB.GetCoins(i.Member.User.ID) // get the coins from the database
		if err != nil { // check for errors
			fmt.Println("error getting coin balance:", err) // print error if there is one

			// respond to interaction
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Sorry, something went wrong retrieving your balance. Please create a support ticket if needed.",
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})

			// print error if there is any
			if err != nil {
				fmt.Println("error responding to balance command:", err)
			}
		}

		// interaction respond
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("🪙 You currently have: **%d coins**!", coins),
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})

		// error if there is an error
		if err != nil {
			fmt.Println("error responding to balance command:", err)
		}
	
	case "index": // index command
		userID := i.Member.User.ID // get user id

		characters, err := b.DB.GetCharacters() // get characters from database
		if err != nil { // error if there is one
			fmt.Println("error getting character index:", err)
			respondEphemeral(s, i, "Sorry, something went wrong loading the character index. Please create a support ticket if needed.") // ephemeral respond
			return
		}

		// rarity order
		rarityOrder := map[string]int {
			"Common": 0,
			"Uncommon": 1,
			"Rare": 2,
			"Epic": 3,
			"Legendary": 4,
		}

		// create the characters slice
		sort.Slice(characters, func(i, j int) bool {
			rarityI := rarityOrder[characters[i].Rarity]
			rarityJ := rarityOrder[characters[j].Rarity]

			if rarityI != rarityJ {
				return rarityI < rarityJ
			}

			return strings.ToLower(characters[i].Name) < strings.ToLower(characters[j].Name)
		})

		// get the collected character ids
		collected, err := b.DB.GetCollectedCharacterIDs(userID)
		if err != nil { // error if we have to
			fmt.Println("error getting collected character IDs:", err)
			respondEphemeral(s, i, "Sorry, something went wrong loading your collection progress. Please create a support ticket if needed.")
			return
		}

		// total collected
		totalCollected := 0
		for _, character := range characters { // loop through characters
			if collected[character.ID] { // check if the character has an id
				totalCollected++ // increment the total
			}
		}

		const pageSize = 10 // page size
		totalPages := (len(characters) + pageSize - 1) / pageSize // total pages

		embed := buildIndexEmbed(characters, collected, 0, totalCollected) // build the index embed

		var components []discordgo.MessageComponent // message components

		// check total pages
		if totalPages > 1 {
			components = []discordgo.MessageComponent{ // build the components
				discordgo.ActionsRow{ // action row
					Components: []discordgo.MessageComponent{
						discordgo.Button{ // previous page button
							Label: "Previous",
							Style: discordgo.SecondaryButton,
							CustomID: "index_page_0",
							Disabled: true,
						},
						discordgo.Button{ // next page button
							Label: "Next",
							Style: discordgo.PrimaryButton,
							CustomID: "index_page_1",
						},
					},
				},
			}
		}

		// interaction respond
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
				Components: components,
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})

		// error if there is one
		if err != nil {
			fmt.Println("error responding to index command: ", err)
		}
	
	case "shop": // shop command
		b.handleShop(s, i) // handle the shop
	}
}

// respond function
func respond(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	message string,
) {
	// simply do an interaction respond
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message, // with the provided message as the content
			},
		},
	)

	// and error if there is an error
	if err != nil {
		fmt.Println("Error responding to interaction: ", err)
	}
}


// ephemeral response (only one user can see it)
func respondEphemeral(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	message string,
) {
	// simple interaction response
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message, // with provided message as content
				Flags: discordgo.MessageFlagsEphemeral, // and ephemeral flags
			},
		},
	)

	// and also error if needed
	if err != nil {
		fmt.Println("error sending ephemeral response:", err)
	}
}

// message handler function
func (b *Bot) handleMessage(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	// check if the message has an author or if it's a bot and return
	if m.Author == nil || m.Author.Bot {
		return
	}

	// get the values
	hasRound, correct, character, shiny := game.Guess(
		m.ChannelID,
		m.Content,
	)

	// check if it has 'hasRound' and if it's correct
	if !hasRound || !correct {
		return
	}

	// add the character to the user's collection
	added, err := b.DB.AddToCollection(
		m.Author.ID,
		character.ID,
	)

	// error if needed
	if err != nil {
		fmt.Println("error saving collection: ", err)
		_, _ = s.ChannelMessageSend(
			m.ChannelID,
			"Sorry, something went wrong saving your character. Please try again later or feel free to create a support ticket.",
		)
		return
	}

	shinyReward := 0 // shiny reward

	// check if the character is shiny
	if shiny {
		shinyReward = game.ShinyReward(character.Rarity) // check the rarity for the reward using the ShinyReward function

		// if the shiny reward is not 0
		if shinyReward > 0 {
			err := b.DB.AddCoins(m.Author.ID, shinyReward) // add the coins
			if err != nil { // error if needed
				fmt.Println("error awarding shiny bonus: ", err) // print error
				_, _ = s.ChannelMessageSend(
					m.ChannelID,
					"⚠️ You collected the character, but something went wrong awarding your shiny bonus. Please create a support ticket.",
				)
				return
			}
		}
	}

	var message string // message string


	if added {
		message = "🎉" + m.Author.Username + " guessed correctly and collected **" + character.Name + "**!" // build winner message
	} else {
		reward := game.DuplicateReward(character.Rarity) // set the duplicate reward

		// if there is a reward
		if reward > 0 {
			err := b.DB.AddCoins(m.Author.ID, reward) // add the coins to the user's balance

			// error if needed
			if err != nil {
				fmt.Println("error awarding duplicate coins: ", err)
				_, _ = s.ChannelMessageSend(
					m.ChannelID,
					"Your duplicate was detected, but something went wrong awarding your coins. Please create a support ticket if needed.",
				)
				return
			}
		}

		message = "🪙 " + m.Author.Username + " guessed correctly! They already own **" + character.Name + "**, so they earned **" + fmt.Sprint(reward) + " coins** for the duplicate!" // duplicate winner message
	}

	// if the character is shiny
	if shiny {
		// build add the shiny text onto the main winner message
		message += fmt.Sprintf(
			"\n✨ **SHINY CHARACTER!** You earned an additonal **%d coins**! ✨",
			shinyReward,
		)
	}

	// send the message into the channel where the command was run
	_, err = s.ChannelMessageSend(m.ChannelID, message)
	// print any errors
	if err != nil {
		fmt.Println("Error sending winner message: ", err)
	}
}

// build index embed function
func buildIndexEmbed(
	characters []models.Character,
	collected map[int]bool,
	page int,
	totalCollected int,
) *discordgo.MessageEmbed {
	// page size
	const pageSize = 10

	// start and end variables
	start := page * pageSize
	end := start + pageSize

	// set the end depending on the length of the character list
	if end > len(characters) {
		end = len(characters)
	}

	var description strings.Builder // string builder

	// loop through the characters
	for _, character := range characters[start:end] {
		// set text to collected text
		if collected[character.ID] {
			fmt.Fprintf(
				&description,
				"✅ **%s** - %s\n",
				character.Name,
				character.Rarity,
			)
		} else {
			// set text to not collected text
			fmt.Fprintf(
				&description,
				"❓ **???** - %s\n",
				character.Rarity,
			)
		}
	}

	totalPages := (len(characters) + pageSize - 1) / pageSize // set the total pages

	// return the final embed
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

// handle index button click function
func (b *Bot) handleIndexButton(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	// interaction respond
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
    	Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	// error if needed
	if err != nil {
    	fmt.Println("Error acknowledging index button:", err)
    	return
	}

	data := i.MessageComponentData() // get the message component data

	// check if the data's customID has the prefix we're looking for
	if !strings.HasPrefix(data.CustomID, "index_page_") {
		return
	}

	// extract the character's id from the custom ID
	page, err := strconv.Atoi(strings.TrimPrefix(data.CustomID, "index_page_"))
	if err != nil {
		respondComponentError(s, i, "Invalid index page.")
		return
	}

	// check if the command is being used outside a server
	if i.Member == nil || i.Member.User == nil {
		respondComponentError(s, i, "The index can only be used in a server")
		return
	}

	// Get the characters from the database
	characters, err := b.DB.GetCharacters()
	// error if needed
	if err != nil {
		fmt.Println("Error loading character index:", err)
		respondComponentError(s, i, "Sorry, I couldn't load the character index.")
		return
	}

	// get the ids of the collected characters
	collected, err := b.DB.GetCollectedCharacterIDs(i.Member.User.ID)
	// error if needed
	if err != nil {
		fmt.Println("Error loading collected character IDs:", err)
		respondComponentError(s, i, "Sorry, I couldn't load your collection progress.")
		return
	}

	// total number of collected characters
	totalCollected := 0
	// loop through the characters
	for _, character := range characters {
		// check if the character's id is in the collected list
		if collected[character.ID] {
			totalCollected++ // increment the total collected characters
		}
	}

	const pageSize = 10 // page size
	totalPages := (len(characters) + pageSize - 1) / pageSize // find out the total pages

	// error if the total amount pages is zero
	if totalPages == 0 {
		respondComponentError(s, i, "There are no characters in the index.")
		return
	}

	// check if the page is more or equal to the total
	if page >= totalPages {
		page = totalPages - 1 // make it one less
	}

	var components []discordgo.MessageComponent // components

	// if the total pages is more than one
	if totalPages > 1 {
		// set the components to a new discordgo message component slice
		components = []discordgo.MessageComponent{
			// actions row
			discordgo.ActionsRow{
				// message components
				Components: []discordgo.MessageComponent{
					// previous button
					discordgo.Button{
						Label: "Previous",
						Style: discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("index_page_%d", page-1),
						Disabled: page == 0,
					},
					// next button
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

	// interaction response edit; edit the previous index message to change the page to the new one decided by the user.
	_, err = s.InteractionResponseEdit(
    	i.Interaction,
    	&discordgo.WebhookEdit{
        	Embeds: &[]*discordgo.MessageEmbed{
            	buildIndexEmbed(characters, collected, page, totalCollected),
        	},
        	Components: &components,
    	},
	)

	// error if needed
	if err != nil {
    	fmt.Println("Error updating character index:", err)
	}
}

// component error function
func respondComponentError(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	message string,
) {
	// interaction respond
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags: discordgo.MessageFlagsEphemeral, // with ephemeral flags
		},
	})

	// error if needed
	if err != nil {
		fmt.Println("Error responding to index button:", err)
	}
}

// handle the "who is this?" button
func (b *Bot) handleWhoIsThisButton(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	// check the type of i
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	// set the customID to the one from the message component data
	customID := i.MessageComponentData().CustomID
	// check if the customID has the desired prefix
	if !strings.HasPrefix(customID, "who_is_this_") {
		return
	}

	// extract the character name from the custom ID
	characterName := strings.TrimPrefix(customID, "who_is_this_")

	// build the message
	message := fmt.Sprintf("Not sure who this character is? This is **%s**!", characterName)

	// ephemeral response with the message
	respondEphemeral(s, i, message)
}