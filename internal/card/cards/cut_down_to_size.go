// Cut Down to Size — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, if they have 4 or more cards in hand, they discard a card."
//
// On-hit rider intentionally not modelled: an optimal defender drops one card to block,
// falls under the 4+ threshold, and dodges the discard — crediting the rider would
// double-count the block already captured by IncomingPhysicalDamage / BlockTotal.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (CutDownToSizeRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (CutDownToSizeYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (CutDownToSizeBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
