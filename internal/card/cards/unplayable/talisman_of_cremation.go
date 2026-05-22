// Talisman of Cremation — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** When you play a card from your banished zone, destroy Talisman of Cremation
// and name a card. Banish all cards with the chosen name from each opposing hero's graveyard."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (TalismanOfCremationBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
