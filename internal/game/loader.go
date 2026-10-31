package game

import (
	"encoding/json"
	"fmt"
	"os"
	"net/url"

	"github.com/greetingsForAlek/JeffDexBot/internal/models"
)

// load characters function
func LoadCharacters(path string) ([]models.Character, error) {
	// get the data from the json file path & handle errors
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read character file: %w", err)
	}

	// loaded characters
	var loadedCharacters []models.Character

	// decode the data and check for erros
	if err := json.Unmarshal(data, &loadedCharacters); err != nil {
		return nil, fmt.Errorf("parsing character file: %w", err)
	}

	// check the length and return if too short
	if len(loadedCharacters) == 0 {
		return nil, fmt.Errorf("character file contains no characters")
	}

	// valid rarities
	validRarities := map[string]bool{
		"Common": true,
		"Uncommon": true,
		"Rare": true,
		"Epic": true,
		"Legendary": true,
	}

	// seenIDs map
	seenIDs := make(map[int]bool)

	// loop through loaded characters
	for i, character := range loadedCharacters {
		// check the id and return error if needed
		if character.ID <= 0 {
			return nil, fmt.Errorf(
				"character at index %d has an invalid id: %d",
				i, 
				character.ID,
			)
		}

		// check if the characterID is already in seenIDs and return error if needed
		if seenIDs[character.ID] {
			return nil, fmt.Errorf(
				"duplicate character ID: %d",
				character.ID,
			)
		}
		seenIDs[character.ID] = true // then set the id to true in seenIDs

		// check if the character's name is an empty string and error if it is
		if character.Name == "" {
			return nil, fmt.Errorf(
				"character with ID %d has no name",
				character.ID,
			)
		}

		// parse the image URLs and return an error if needed
		parsedURL, err := url.ParseRequestURI(character.ImageURL)
		if err != nil || parsedURL.Scheme != "https" || parsedURL.Host == "" {
			return nil, fmt.Errorf(
				"character %q has no valid image URL: %q (must be valid HTTPS URL)",
				character.Name,
				character.ImageURL,
			)
		}

		// then check if the rarity is valid and also error if needed
		if !validRarities[character.Rarity] {
			return nil, fmt.Errorf(
				"character %q has invalid rarity %q",
				character.Name,
				character.Rarity,
			)
		}
	}

	return loadedCharacters, nil // and then return
}