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

func (SigilOfDeadwoodBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateStartOfTurnAura(self, sigilOfDeadwoodAuraHandler, 1)
}

// sigilOfDeadwoodAuraHandler creates 1 runechant on the next-turn fire and destroys the aura.
func sigilOfDeadwoodAuraHandler(ge card.GameEngine, l card.Logger, a card.Aura) {
	ge.CreateRunechants(1)
	l.AppendPostTrigger(a.CardName(), "Created a runechant", 1)
	a.Destroy(true)
}
