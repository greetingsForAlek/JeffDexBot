package game

import (
	"math/rand/v2"

	"github.com/greetingsForAlek/JeffDexBot/internal/models"
)

// rarity weight struct
type rarityWeight struct {
	Rarity string
	Weight int
}

// the rarity weights
// basically the chance of the character appearing
var rarityWeights = []rarityWeight {
	{Rarity: "Common", Weight: 50},
	{Rarity: "Uncommon", Weight: 25},
	{Rarity: "Rare", Weight: 15},
	{Rarity: "Epic", Weight: 7},
	{Rarity: "Legendary", Weight: 3},
}

var characters []models.Character // characters slice

// random function to pick a random character
func Random() models.Character {
	// if the slice length is 0
	if len(characters) == 0 {
		return models.Character{} // return an empty character
	}

	// map of available characters
	available := make(map[string][]models.Character)

	// loop through provided characters
	for _, character := range characters {
		// append
		available[character.Rarity] = append(
			available[character.Rarity],
			character,
		)
	}

	totalWeight := 0 // get the total weight

	// loop through rarity weights
	for _, rarity := range rarityWeights {
		if len(available[rarity.Rarity]) > 0 {
			totalWeight += rarity.Weight // increment the total weight
		}
	}

	// pick at random from the total weight
	roll := rand.IntN(totalWeight)
	
	// loop through the weights
	for _, rarity := range rarityWeights {
		// get the matching weight
		matching := available[rarity.Rarity]

		// if matching is 0 then continue
		if len(matching) == 0 {
			continue
		}

		// if the roll is less than the weight then return
		if roll < rarity.Weight {
			return matching[rand.IntN(len(matching))]
		}

		// take away the weight from the roll
		roll -= rarity.Weight
	}

	// return
	return characters[rand.IntN(len(characters))]
}

// set the characters function
func SetCharacters(newCharacters []models.Character) {
	// basically just append the new characters to characters
	characters = append([]models.Character(nil), newCharacters...)
}

// all characters function
func AllCharacters() []models.Character {
	// append
	return append([]models.Character(nil), characters...)
}