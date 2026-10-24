package game

import (
	"encoding/json"
	"fmt"
	"os"
	"net/url"

	"github.com/greetingsForAlek/JeffDexBot/internal/models"
)

func LoadCharacters(path string) ([]models.Character, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read character file: %w", err)
	}

	var loadedCharacters []models.Character

	if err := json.Unmarshal(data, &loadedCharacters); err != nil {
		return nil, fmt.Errorf("parsing character file: %w", err)
	}

	if len(loadedCharacters) == 0 {
		return nil, fmt.Errorf("character file contains no characters")
	}

	validRarities := map[string]bool{
		"Common": true,
		"Uncommon": true,
		"Rare": true,
		"Epic": true,
		"Legendary": true,
	}

	seenIDs := make(map[int]bool)

	for i, character := range loadedCharacters {
		if character.ID <= 0 {
			return nil, fmt.Errorf(
				"character at index %d has an invalid id: %d",
				i, 
				character.ID,
			)
		}

		if seenIDs[character.ID] {
			return nil, fmt.Errorf(
				"duplicate character ID: %d",
				character.ID,
			)
		}
		seenIDs[character.ID] = true

		if character.Name == "" {
			return nil, fmt.Errorf(
				"character with ID %d has no name",
				character.ID,
			)
		}

		parsedURL, err := url.ParseRequestURI(character.ImageURL)
		if err != nil || parsedURL.Scheme != "https" || parsedURL.Host == "" {
			return nil, fmt.Errorf(
				"character %q has no valid image URL: %q (must be valid HTTPS URL)",
				character.Name,
				character.ImageURL,
			)
		}

		if !validRarities[character.Rarity] {
			return nil, fmt.Errorf(
				"character %q has invalid rarity %q",
				character.Name,
				character.Rarity,
			)
		}
	}

	return loadedCharacters, nil
}