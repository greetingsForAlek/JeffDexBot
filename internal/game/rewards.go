package game

func DuplicateReward(rarity string) int {
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

func ShinyReward(rarity string) int {
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