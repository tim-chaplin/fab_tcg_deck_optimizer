// Talisman of Cremation — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** When you play a card from your banished zone, destroy Talisman of Cremation
// and name a card. Banish all cards with the chosen name from each opposing hero's graveyard."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: self-destroys on play-from-banished → banish a named card from opposing
// graveyards

func (TalismanOfCremationBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
}
