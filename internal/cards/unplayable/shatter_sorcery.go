// Shatter Sorcery — Generic Instant. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "Choose 1 or both; - Destroy target aura permanent with Sigil in its name.
// - Prevent the next 1 arcane damage that would be dealt to target hero this turn."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (ShatterSorceryBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.Log(self, 0)
}
