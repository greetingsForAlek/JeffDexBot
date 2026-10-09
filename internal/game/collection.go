package game

import (
	"sync"

	"github.com/greetingsForAlek/JeffDexBot/internal/models"
)

var (
	collections = make(map[string]map[int]models.Character)
	collectionsMu sync.Mutex
)

func AddToCollection (
	userID string,
	character models.Character,
) bool {
	collectionsMu.Lock()
	defer collectionsMu.Unlock()

	if collections[userID] == nil {
		collections[userID] = make(map[int]models.Character)
	}

	if _, exists := collections[userID][character.ID]; exists {
		return false
	}

	collections[userID][character.ID] = character
	return true
}

func GetCollection(userID string) []models.Character {
	collectionsMu.Lock()
	defer collectionsMu.Unlock()

	playerCollection := collections[userID]
	characters := make([]models.Character, 0, len(playerCollection))

	for _, character := range playerCollection {
		characters = append(characters, character)
	}

	return characters
}