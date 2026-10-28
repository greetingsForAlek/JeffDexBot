package game

import (
	"math/rand/v2"
	"strings"
	"sync"

	"github.com/greetingsForAlek/JeffDexBot/internal/models"
)

type Round struct {
	Character models.Character
	Active bool
	Shiny bool
}

var (
	rounds = make(map[string]*Round)
	roundsMu sync.Mutex
)

func StartRound(channelID string) models.Character {
	roundsMu.Lock()
	defer roundsMu.Unlock()

	if round, exists := rounds[channelID]; exists && round.Active {
		return round.Character
	}

	character := Random()
	shiny := rand.IntN(100) < 15

	rounds[channelID] = &Round {
		Character: character,
		Active: true,
		Shiny: shiny,
	}

	return character
}

func Guess(channelID, guess string) (bool, bool, models.Character, bool) {
	roundsMu.Lock()
	defer roundsMu.Unlock()

	round, exists := rounds[channelID]
	if !exists || !round.Active {
		return false, false, models.Character{}, false
	}

	if strings.EqualFold(
		strings.TrimSpace(guess),
		strings.TrimSpace(round.Character.Name),
	) {
		round.Active = false
		return true, true, round.Character, round.Shiny
	}

	return true, false, models.Character{}, false
}

func HasActiveRound(channelID string) bool {
	roundsMu.Lock()
	defer roundsMu.Unlock()

	round, exists := rounds[channelID]
	return exists && round.Active
}