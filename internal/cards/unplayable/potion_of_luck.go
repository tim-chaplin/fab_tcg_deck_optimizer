// Potion of Luck — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy Potion of Luck: Shuffle your hand and arsenal into your deck then
// draw that many cards."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (PotionOfLuckBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {}
