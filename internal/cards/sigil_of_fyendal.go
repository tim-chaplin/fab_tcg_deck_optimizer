// Sigil of Fyendal — Generic Action - Aura. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "**Go again** At the beginning of your action phase, destroy this. When this leaves the
// arena, gain 1{h}."
//
// Handler credits +1 next turn for the 1{h} gain (valued 1-to-1 with damage).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (SigilOfFyendalBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.CreateStartOfTurnAura(self, sigilOfFyendalAuraHandler, 1)
}

// sigilOfFyendalAuraHandler credits the +1 health (valued 1-to-1 with damage) next turn
// and destroys the aura. Top-level so Aura.Handler doesn't allocate a closure.
func sigilOfFyendalAuraHandler(g card.GameEngine, l card.Logger, a card.Aura) {
	g.AddValue(1)
	l.AppendPostTrigger(a.CardName(), "Gained 1 health", 1)
	a.Destroy(true)
}
