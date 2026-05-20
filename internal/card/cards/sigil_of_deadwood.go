// Sigil of Deadwood — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2, Go again.
// Only printed in Blue.
// Text: "Go again. At the beginning of your action phase, destroy this. When this leaves the
// arena, create a Runechant token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (SigilOfDeadwoodBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateStartOfTurnAura(self, selfDestructAuraHandler, 1)
}

// OnLeavesArena runs the "when this leaves the arena" clause: create 1 Runechant.
func (c SigilOfDeadwoodBlue) OnLeavesArena(g card.GameEngine, l card.Logger) {
	g.CreateRunechants(1)
	l.AppendPostTrigger(c.DisplayName(), "Created a runechant", 1)
}
