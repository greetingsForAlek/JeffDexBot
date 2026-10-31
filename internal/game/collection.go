package game

import (
	"sync"

	"github.com/greetingsForAlek/JeffDexBot/internal/models"
)

// collections and collectionsMu
var (
	collections = make(map[string]map[int]models.Character)
	collectionsMu sync.Mutex
)

// add to collection function
func AddToCollection (
	userID string,
	character models.Character,
) bool {
	// lock
	collectionsMu.Lock()
	defer collectionsMu.Unlock()

	// if the userID does not have a collection
	if collections[userID] == nil {
		// make one
		collections[userID] = make(map[int]models.Character)
	}

	// check if the characterID exists within the user's collection
	if _, exists := collections[userID][character.ID]; exists {
		return false // return
	}

	collections[userID][character.ID] = character // set the character
	return true // and retrun
}

// get collection function
func GetCollection(userID string) []models.Character {
	// lock
	collectionsMu.Lock()
	defer collectionsMu.Unlock()

	// player collection and characters slice
	playerCollection := collections[userID]
	characters := make([]models.Character, 0, len(playerCollection))

	// loop through characters in the player collection
	for _, character := range playerCollection {
		characters = append(characters, character) // append
	}

	return characters // return
}