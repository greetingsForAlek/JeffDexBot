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

var characters = []models.Character {
	{
		ID: 1,
		Name: "Jeff",
		ImageURL: "https://res.cloudinary.com/cwohba23/image/upload/v1786985403/Jeff.png",
		Rarity: "Common",
	},
	{
		ID: 2,
		Name: "Rob",
		ImageURL: "https://res.cloudinary.com/cwohba23/image/upload/v1786985433/Rob.png",
		Rarity: "Rare",
	},
	{
		ID: 3,
		Name: "Mr No Internet",
		ImageURL: "https://res.cloudinary.com/cwohba23/image/upload/v1786985444/MNI.png",
		Rarity: "Legendary",
	},
}

func Random() models.Character {
	return characters[rand.IntN(len(characters))]
}

func SetCharacters(newCharacters []models.Character) {
	characters = append([]models.Character(nil), newCharacters...)
}

func AllCharacters() []models.Character {
	return append([]models.Character(nil), characters...)
}