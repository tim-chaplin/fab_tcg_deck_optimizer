// Rise Above — Generic Defense Reaction. Cost 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "You may put a card from your hand on top of your deck rather than pay Rise Above's {r}
// cost."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// not implemented: hand-as-cost alt cost not modelled; card fails when printed cost can't be paid

func (RiseAboveRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}

// not implemented: hand-as-cost alt cost not modelled; card fails when printed cost can't be paid

func (RiseAboveYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}

// not implemented: hand-as-cost alt cost not modelled; card fails when printed cost can't be paid

func (RiseAboveBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}
