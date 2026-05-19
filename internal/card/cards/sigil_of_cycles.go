// Sigil of Cycles — Generic Action - Aura. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "**Go again** At the beginning of your action phase, destroy this. When this leaves
// the arena, discard a card then draw a card."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (SigilOfCyclesBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateStartOfTurnAura(self, sigilOfCyclesAuraHandler, 1)
}

func sigilOfCyclesAuraHandler(ge card.GameEngine, l card.Logger, a card.Aura) {
	a.Destroy(true)
	discarded, ok := ge.Discard()
	if !ok {
		return
	}
	ge.DrawOne()
	l.AppendPostTriggerf(a.CardName(), 0, "Discarded %s and drew", discarded.DisplayName())
}
