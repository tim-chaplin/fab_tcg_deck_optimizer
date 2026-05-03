// Sigil of Deadwood — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2, Go again.
// Only printed in Blue.
// Text: "Go again. At the beginning of your action phase, destroy this. When this leaves the
// arena, create a Runechant token."
//
// Handler creates 1 Runechant next turn.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var sigilOfDeadwoodTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type SigilOfDeadwoodBlue struct{}

func (SigilOfDeadwoodBlue) ID() ids.CardID          { return ids.SigilOfDeadwoodBlue }
func (SigilOfDeadwoodBlue) Name() string            { return "Sigil of Deadwood" }
func (SigilOfDeadwoodBlue) Cost(*sim.TurnState) int { return 0 }
func (SigilOfDeadwoodBlue) Pitch() int              { return 3 }
func (SigilOfDeadwoodBlue) Attack() int             { return 0 }
func (SigilOfDeadwoodBlue) Defense() int            { return 2 }
func (SigilOfDeadwoodBlue) Types() card.TypeSet     { return sigilOfDeadwoodTypes }
func (SigilOfDeadwoodBlue) GoAgain() bool           { return true }
func (SigilOfDeadwoodBlue) AddsFutureValue()        {}
func (c SigilOfDeadwoodBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.AddAura(sim.Aura{
		Self:        c,
		TriggerType: sim.TriggerStartOfTurn,
		Count:       1,
		Handler: func(s *sim.TurnState, t *sim.Aura) int {
			n := s.CreateRunechants(1)
			s.DestroyAura(t)
			return n
		},
		LogText: sim.DisplayName(c) + ": Created a runechant",
	})
	s.Log(self, 0)
}
