// Sigil of Fyendal — Generic Action - Aura. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "**Go again** At the beginning of your action phase, destroy this. When this leaves the
// arena, gain 1{h}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

func (SigilOfFyendalBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateAura(self, triggertype.StartOfTurn, selfDestructAuraHandler, 1, false, nil)
}

// OnLeavesArena runs the "when this leaves the arena" clause: gain 1{h}, valued 1-to-1 with
// damage.
func (c SigilOfFyendalBlue) OnLeavesArena(g card.GameEngine, l card.Logger) {
	g.AddValue(1)
	l.AppendPostTrigger(c.DisplayName(), "Gained 1 health", 1)
}
