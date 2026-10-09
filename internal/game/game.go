package game

import (
	"strings"
	"sync"

	"github.com/greetingsForAlek/JeffDexBot/internal/models"
)

type Round struct {
	Character models.Character
	Active bool
}

var (
	rounds = make(map[string]*Round)
	roundsMu sync.Mutex
)

func StartRound(channelID string) models.Character {
	roundsMu.Lock()
	defer roundsMu.Unlock()

	character := Random()

	rounds[channelID] = &Round {
		Character: character,
		Active: true,
	}

	return character
}

func Guess(channelID, guess string) (bool, bool, models.Character) {
	roundsMu.Lock()
	defer roundsMu.Unlock()

	round, exists := rounds[channelID]
	if !exists || !round.Active {
		return false, false, models.Character{}
	}

	if strings.EqualFold(
		strings.TrimSpace(guess),
		strings.TrimSpace(round.Character.Name),
	) {
		round.Active = false
		return true, true, round.Character
	}

	return true, false, models.Character{}
}