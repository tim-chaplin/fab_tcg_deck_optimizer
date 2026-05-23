// Arcane Cussing — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. When you deal or are dealt damage, destroy this. When this leaves the arena
// during your turn, create N Runechants." (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

func (ArcaneCussingRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateAura(self, triggertype.Hit|triggertype.DamageTaken, selfDestructAuraHandler, 3, false, nil)
}

func (ArcaneCussingYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateAura(self, triggertype.Hit|triggertype.DamageTaken, selfDestructAuraHandler, 2, false, nil)
}

func (ArcaneCussingBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateAura(self, triggertype.Hit|triggertype.DamageTaken, selfDestructAuraHandler, 1, false, nil)
}

// OnLeavesArena runs the "when this leaves the arena during your turn" clause.
func (c ArcaneCussingRed) OnLeavesArena(g card.GameEngine, l card.Logger) {
	arcaneCussingLeavesArena(g, l, c.DisplayName(), 3)
}

func (c ArcaneCussingYellow) OnLeavesArena(g card.GameEngine, l card.Logger) {
	arcaneCussingLeavesArena(g, l, c.DisplayName(), 2)
}

func (c ArcaneCussingBlue) OnLeavesArena(g card.GameEngine, l card.Logger) {
	arcaneCussingLeavesArena(g, l, c.DisplayName(), 1)
}

// arcaneCussingLeavesArena creates n Runechants when Arcane Cussing leaves the arena during
// our turn — destroyed by an attack we land. Leaving off our turn (the defense phase,
// destroyed by taking damage) doesn't satisfy the printed "during your turn" clause.
func arcaneCussingLeavesArena(g card.GameEngine, l card.Logger, name string, n int) {
	if !g.IsMyTurn() {
		return
	}
	g.CreateRunechants(n)
	l.AppendPostTriggerf(name, n, "Created %d runechants", n)
}
