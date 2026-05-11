// Talisman of Balance — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** At the beginning of your end phase, if you have less cards in arsenal than an
// opposing hero, destroy Talisman of Balance and put the top card of your deck into an empty
// arsenal zone you control."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: end-phase arsenal-fill from top of deck if behind on arsenal count

func (TalismanOfBalanceBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
}
