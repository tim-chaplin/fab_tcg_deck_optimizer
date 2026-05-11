// Gravekeeping — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4, Blue 3. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks a hero, you may banish a card from their graveyard."
//
// The on-attack opponent-graveyard banish is not modelled: opponent-graveyard hate is a
// sideboard-time consideration against graveyard-recursion matchups, not a property of the
// card-vs-deck Value the simulator optimises.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func gravekeepingPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
}

func (GravekeepingRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	gravekeepingPlay(s, l, self)
}

func (GravekeepingYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	gravekeepingPlay(s, l, self)
}

func (GravekeepingBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	gravekeepingPlay(s, l, self)
}
