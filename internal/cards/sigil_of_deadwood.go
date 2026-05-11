// Sigil of Deadwood — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2, Go again.
// Only printed in Blue.
// Text: "Go again. At the beginning of your action phase, destroy this. When this leaves the
// arena, create a Runechant token."
//
// Handler creates 1 Runechant next turn.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (c SigilOfDeadwoodBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.AddAura(sim.Aura{
		Trigger: sim.Trigger{TriggerType: sim.TriggerStartOfTurn, Handler: sigilOfDeadwoodAuraHandler},
		Self:    sim.CardOrTokenType{Card: c},
		Count:   1,
	})
}

// sigilOfDeadwoodAuraHandler creates 1 runechant on the next-turn fire and destroys the
// aura. Top-level so the Aura.Handler assignment doesn't allocate a closure.
func sigilOfDeadwoodAuraHandler(s *sim.TurnState, l card.Logger, _ *sim.Trigger, a *sim.Aura) {
	name := a.Self.DisplayName()
	s.CreateRunechants(1)
	l.AppendPostTrigger(name, "Created a runechant", 1)
	s.DestroyAura(a, true)
}
