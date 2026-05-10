// Potion of Seeing — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy Potion of Seeing: Look at target hero's hand."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (PotionOfSeeingBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { l.Log(self, 0) }
