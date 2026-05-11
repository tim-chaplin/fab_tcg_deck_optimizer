// Sigil of Deadwood — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2, Go again.
// Only printed in Blue.
// Text: "Go again. At the beginning of your action phase, destroy this. When this leaves the
// arena, create a Runechant token."
//
// Handler creates 1 Runechant next turn.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (SigilOfDeadwoodBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.AddStartOfTurnAura(self, sigilOfDeadwoodAuraHandler, 1)
}

// sigilOfDeadwoodAuraHandler creates 1 runechant on the next-turn fire and destroys the
// aura. Top-level so the Aura.Handler assignment doesn't allocate a closure.
func sigilOfDeadwoodAuraHandler(s card.GameEngine, l card.Logger, a card.Aura) {
	s.CreateRunechants(1)
	l.AppendPostTrigger(a.SelfName(), "Created a runechant", 1)
	a.Destroy(true)
}
