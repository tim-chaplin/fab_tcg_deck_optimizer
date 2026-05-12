// Eirina's Prayer — Generic Instant. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Reveal the top card of your deck. Prevent the next X arcane damage that would be dealt to
// your hero this turn, where X is 6 minus the pitch value of the card revealed this way."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// not implemented: Instant prevent X arcane to your hero; X scaled by revealed top-card pitch

func (EirinasPrayerRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {}

// not implemented: Instant prevent X arcane to your hero; X scaled by revealed top-card pitch

func (EirinasPrayerYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {}

// not implemented: Instant prevent X arcane to your hero; X scaled by revealed top-card pitch

func (EirinasPrayerBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {}
