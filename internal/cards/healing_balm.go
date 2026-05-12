// Healing Balm — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2. Printed health-gain: Red 3{h}, Yellow 2{h}, Blue 1{h}.
//
// Text: "Gain N{h}." (N is the printed variant value above.)
//
// Modelling: health is valued 1-to-1 with damage.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// healingBalmPlay credits N{h} as a sub-line under self.
func healingBalmPlay(g card.GameEngine, l card.Logger, self *card.CardState, heal int) {
	g.AddValue(heal)
	l.AppendPostTriggerf(self.Card.DisplayName(), heal, "Gained %d health", heal)
}

func (HealingBalmRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	healingBalmPlay(g, l, self, 3)
}

func (HealingBalmYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	healingBalmPlay(g, l, self, 2)
}

func (HealingBalmBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	healingBalmPlay(g, l, self, 1)
}
