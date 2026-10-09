package game

import (
	"math/rand/v2"

	"github.com/greetingsForAlek/JeffDexBot/internal/models"
)

var characters = []models.Character {
	{
		ID: 1,
		Name: "Jeff",
		ImageURL: "https://res.cloudinary.com/cwohba23/image/upload/v1786985403/Jeff.png",
	},
	{
		ID: 2,
		Name: "Rob",
		ImageURL: "https://res.cloudinary.com/cwohba23/image/upload/v1786985433/Rob.png",
	},
}

func Random() models.Character {
	return characters[rand.IntN(len(characters))]
}