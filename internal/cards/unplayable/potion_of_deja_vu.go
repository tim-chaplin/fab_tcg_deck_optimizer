// Potion of Déjà Vu — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy Potion of Déjà Vu: Put all cards from your pitch zone on top of your
// deck in any order."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (PotionOfDejaVuBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
