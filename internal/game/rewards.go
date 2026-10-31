package game

// duplicate reward function
// stores all the duplicate rewards
func DuplicateReward(rarity string) int {
	// switch the rarity and return the coins that will be awarded for a duplicate of that rarity
	switch rarity {
	case "Common":
		return 10
	case "Uncommon":
		return 20
	case "Rare":
		return 50
	case "Epic":
		return 100
	case "Legendary":
		return 250
	default:
		return 0
	}
}

// shiny reward function
func ShinyReward(rarity string) int {
	// basically the same as the duplicate reward function
	// switch the rarity and return the amount of coins that will be returned as a shiny reward for that rarity
	switch rarity{
	case "Common":
		return 25
	case "Uncommon":
		return 50
	case "Rare":
		return 100
	case "Epic":
		return 250
	case "Legendary":
		return 500
	default:
		return 0
	}
}