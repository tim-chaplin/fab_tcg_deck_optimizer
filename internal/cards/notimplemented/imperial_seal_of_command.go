// Imperial Seal of Command — Generic Action - Item. Cost 0. Printed pitch variants: Red 1.
//
// Text: "**Legendary** **Action** - Destroy this: Defense reaction cards can't be played this turn.
// If you are Royal, the next time you hit a hero this turn, destroy all cards in their arsenal.
// **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: activated 'no DR this turn' + Royal-only arsenal-wipe on hit

func (ImperialSealOfCommandRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
}
