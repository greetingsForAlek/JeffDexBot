package game

import (
	"math/rand/v2"
	"strings"
	"sync"

	"github.com/greetingsForAlek/JeffDexBot/internal/models"
)

// round struct
type Round struct {
	Character models.Character
	Active bool
	Shiny bool
}

// rounds and roundsMu
var (
	rounds = make(map[string]*Round)
	roundsMu sync.Mutex
)

// start round function
func StartRound(channelID string) models.Character {
	// lock
	roundsMu.Lock()
	defer roundsMu.Unlock()

	character := Random() // random character
	shiny := rand.IntN(100) < 15 // shiny

	// start a new round
	rounds[channelID] = &Round {
		Character: character,
		Active: true,
		Shiny: shiny,
	}

	return character // and return the character
}

// guess function
func Guess(channelID, guess string) (bool, bool, models.Character, bool) {
	// lock
	roundsMu.Lock()
	defer roundsMu.Unlock()

	// check if round exists and is active
	round, exists := rounds[channelID]
	if !exists || !round.Active {
		return false, false, models.Character{}, false
	}

	// check if the user guessed correctly
	if strings.EqualFold(
		strings.TrimSpace(guess),
		strings.TrimSpace(round.Character.Name),
	) {
		// make the round inactive
		round.Active = false
		return true, true, round.Character, round.Shiny
	}

	return true, false, models.Character{}, false // and return
}

// has active round function
// basically for checking if there is currently an active round
func HasActiveRound(channelID string) bool {
	// lock
	roundsMu.Lock()
	defer roundsMu.Unlock()

	// check if the round is currently active
	round, exists := rounds[channelID]
	return exists && round.Active
}