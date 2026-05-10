// Clarity Potion — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy Clarity Potion: **Opt 2**"

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (ClarityPotionBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { l.Log(self, 0) }
