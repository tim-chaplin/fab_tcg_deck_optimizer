// Sigil of Cycles — Generic Action - Aura. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "**Go again** At the beginning of your action phase, destroy this. When this leaves
// the arena, discard a card then draw a card."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

func (SigilOfCyclesBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateAura(self.Card, triggertype.StartOfTurn, selfDestructAuraHandler, 1, false, nil)
}

// OnLeavesArena runs the "when this leaves the arena" clause: discard a card, then draw a
// card. An empty hand drops the rider.
func (c SigilOfCyclesBlue) OnLeavesArena(g card.GameEngine, l card.Logger) {
	discarded, ok := g.Discard()
	if !ok {
		return
	}
	g.DrawOne()
	l.AppendPostTriggerf(c.DisplayName(), 0, "Discarded %s and drew", discarded.DisplayName())
}
