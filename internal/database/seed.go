package database

import (
	"fmt"

	"github.com/greetingsForAlek/JeffDexBot/internal/models"
)

func (d *Database) SeedCharacters(characters []models.Character) error {
	for _, character := range characters {
		err := d.AddCharacter(
			character.ID,
			character.Name,
			character.ImageURL,
			character.Rarity,
		)

		if err != nil {
			return fmt.Errorf(
				"seed character: %q: %w",
				character.Name,
				err,
			)
		}
	}

	return nil
}