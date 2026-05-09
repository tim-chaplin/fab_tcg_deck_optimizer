// Destructive Tendencies — Generic Instant. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "Choose 1 or both; - Remove all counters from target item token. - Remove all
// counters from target aura token."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (DestructiveTendenciesBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.Log(self, 0)
}
