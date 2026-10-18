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
		Rarity: "Common",
	},
	{
		ID: 3,
		Name: "Mr No Internet",
		ImageURL: "https://res.cloudinary.com/cwohba23/image/upload/v1786985444/MNI.png",
		Rarity: "Common",
	},
	{
		ID: 4,
		Name: "Mr Squirt",
		ImageURL: "https://res.cloudinary.com/cwohba23/image/upload/v1786985457/MrSquirt.png",
		Rarity: "Common",
	},
	{
		ID: 5,
		Name: "Bald Boi",
		ImageURL: "https://res.cloudinary.com/cwohba23/image/upload/v1786985468/BaldBoi.png",
		Rarity: "Common",
	},
	{
		ID: 6,
		Name: "Mr No Name",
		ImageURL: "https://res.cloudinary.com/cwohba23/image/upload/v1786986655/MrNoName.png",
		Rarity: "Common",
	},
	{
		ID: 7,
		Name: "Mr Paper",
		ImageURL: "https://res.cloudinary.com/cwohba23/image/upload/v1786986695/MrPaper.png",
		Rarity: "Common",
	},
	{
		ID: 8,
		Name: "Mrs Paper",
		ImageURL: "https://res.cloudinary.com/cwohba23/image/upload/v1786986747/MrsPaper.png",
		Rarity: "Common",
	},
	{
		ID: 9,
		Name: "Joe",
		ImageURL: "https://res.cloudinary.com/cwohba23/image/upload/v1786986664/Joe.png",
		Rarity: "Common",
	},
	{
		ID: 10,
		Name: "Top Hat Tom",
		ImageURL: "https://res.cloudinary.com/cwohba23/image/upload/v1786986756/TopHatTom.png",
		Rarity: "Common",
	},
}

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