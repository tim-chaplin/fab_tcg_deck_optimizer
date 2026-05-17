// Healing Balm — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2. Printed health-gain: Red 3{h}, Yellow 2{h}, Blue 1{h}.
//
// Text: "Gain N{h}." (N is the printed variant value above.)
//
// Modelling: health is valued 1-to-1 with damage.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// healingBalmPlay credits N{h} as a sub-line under self.
func healingBalmPlay(ge card.GameEngine, l card.Logger, self *card.CardState, heal int) {
	ge.AddValue(heal)
	l.AppendPostTriggerf(self.Card.DisplayName(), heal, "Gained %d health", heal)
}

func (HealingBalmRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	healingBalmPlay(ge, l, self, 3)
}

func (HealingBalmYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	healingBalmPlay(ge, l, self, 2)
}

func (HealingBalmBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	healingBalmPlay(ge, l, self, 1)
}
