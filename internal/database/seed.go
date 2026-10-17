package database

import (
	"fmt"

	"github.com/greetingsForAlek/JeffDexBot/internal/game"
)

func (d *Database) SeedCharacters() error {
	characters := game.AllCharacters()

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