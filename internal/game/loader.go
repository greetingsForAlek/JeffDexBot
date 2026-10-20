package game

import (
	"encoding/json"
	"fmt"
	"os"

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

	return loadedCharacters, nil
}