package database

import (
	"fmt"

	"github.com/greetingsForAlek/JeffDexBot/internal/models"
)

// seed characters function
func (d *Database) SeedCharacters(characters []models.Character) error {
	// loop through characters
	for _, character := range characters {
		// add characters
		err := d.AddCharacter(
			character.ID,
			character.Name,
			character.ImageURL,
			character.Rarity,
		)

		// check for errors
		if err != nil {
			// throw an error
			return fmt.Errorf(
				"seed character: %q: %w",
				character.Name,
				err,
			)
		}
	}

	return nil // and return nil for no errors
}