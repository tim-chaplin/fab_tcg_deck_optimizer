package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// TurnSummary is the sim-runtime working type returned by playOneTurn: winning card-role
// assignments plus a *GameState reflecting the start-of-next-turn boundary (pitched
// recycled to deck bottom, next-turn hand drawn, accrued Value carried). Pure data, no
// rules engine. deck.BestTurn persists only the durable fields.
type TurnSummary struct {
	BestLine       []card.CardAssignment
	SwungWeapons   []string
	Value          int
	State          *gameengine.GameState
	IncomingDamage int
	Cacheable      bool
}

// ArsenalIn returns the assignment for the card that started the turn in the arsenal.
func (t TurnSummary) ArsenalIn() (card.CardAssignment, bool) {
	for _, a := range t.BestLine {
		if a.FromArsenal {
			return a, true
		}
	}
	return card.CardAssignment{}, false
}
