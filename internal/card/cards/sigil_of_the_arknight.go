// Sigil of the Arknight — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2. Go again.
// Only printed in Blue.
//
// Text: "At the beginning of your action phase, destroy this. When this leaves the arena,
// reveal the top card of your deck. If it's an attack action card, put it into your hand."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (SigilOfTheArknightBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateStartOfTurnAura(self, selfDestructAuraHandler, 1)
}

// OnLeavesArena runs the "when this leaves the arena" clause: reveal the top card of the
// deck and, if it's an attack action card, draw it. Logs the outcome — hit or whiff — so
// the random reveal shows in the printout; an empty deck is the silent edge case.
func (c SigilOfTheArknightBlue) OnLeavesArena(g card.GameEngine, l card.Logger) {
	top, ok := g.PeekDeck()
	if !ok {
		return
	}
	name := c.DisplayName()
	if top.Types(nil).IsAttackAction() {
		g.DrawOne()
		l.AppendPostTriggerf(name, 0, "%s drew %s into hand", name, top.DisplayName())
		return
	}
	l.AppendPostTriggerf(name, 0, "%s revealed %s but didn't draw it", name, top.DisplayName())
}
