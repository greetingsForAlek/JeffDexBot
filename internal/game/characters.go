package game

import (
	"math/rand/v2"

	"github.com/greetingsForAlek/JeffDexBot/internal/models"
)

type rarityWeight struct {
	Rarity string
	Weight int
}

var rarityWeights = []rarityWeight {
	{Rarity: "Common", Weight: 50},
	{Rarity: "Uncommon", Weight: 25},
	{Rarity: "Rare", Weight: 15},
	{Rarity: "Epic", Weight: 7},
	{Rarity: "Legendary", Weight: 3},
}

var characters []models.Character

func Random() models.Character {
	if len(characters) == 0 {
		return models.Character{}
	}

	available := make(map[string][]models.Character)

	for _, character := range characters {
		available[character.Rarity] = append(
			available[character.Rarity],
			character,
		)
	}

	totalWeight := 0

	for _, rarity := range rarityWeights {
		if len(available[rarity.Rarity]) > 0 {
			totalWeight += rarity.Weight
		}
	}

	roll := rand.IntN(totalWeight)
	
	for _, rarity := range rarityWeights {
		matching := available[rarity.Rarity]

		if len(matching) == 0 {
			continue
		}

		if roll < rarity.Weight {
			return matching[rand.IntN(len(matching))]
		}

		roll -= rarity.Weight
	}

	return characters[rand.IntN(len(characters))]
}

func SetCharacters(newCharacters []models.Character) {
	characters = append([]models.Character(nil), newCharacters...)
}

func AllCharacters() []models.Character {
	return append([]models.Character(nil), characters...)
}