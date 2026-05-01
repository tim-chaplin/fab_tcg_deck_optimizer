// Runerager Swarm — Runeblade Action - Attack. Cost 0, Defense 3.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "If you've played or created an aura this turn, this gets go again."
//
// Go again is conditional on the aura clause, not a printed keyword (docs/dev-standards.md
// covers the conditional grant wiring).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var runeragerSwarmTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type RuneragerSwarmRed struct{}

func (RuneragerSwarmRed) ID() ids.CardID          { return ids.RuneragerSwarmRed }
func (RuneragerSwarmRed) Name() string            { return "Runerager Swarm" }
func (RuneragerSwarmRed) Cost(*sim.TurnState) int { return 0 }
func (RuneragerSwarmRed) Pitch() int              { return 1 }
func (RuneragerSwarmRed) Attack() int             { return 3 }
func (RuneragerSwarmRed) Defense() int            { return 3 }
func (RuneragerSwarmRed) Types() card.TypeSet     { return runeragerSwarmTypes }
func (RuneragerSwarmRed) GoAgain() bool           { return false }
func (RuneragerSwarmRed) Play(s *sim.TurnState, self *sim.CardState) {
	runeragerSwarmPlay(s, self)
}

type RuneragerSwarmYellow struct{}

func (RuneragerSwarmYellow) ID() ids.CardID          { return ids.RuneragerSwarmYellow }
func (RuneragerSwarmYellow) Name() string            { return "Runerager Swarm" }
func (RuneragerSwarmYellow) Cost(*sim.TurnState) int { return 0 }
func (RuneragerSwarmYellow) Pitch() int              { return 2 }
func (RuneragerSwarmYellow) Attack() int             { return 2 }
func (RuneragerSwarmYellow) Defense() int            { return 3 }
func (RuneragerSwarmYellow) Types() card.TypeSet     { return runeragerSwarmTypes }
func (RuneragerSwarmYellow) GoAgain() bool           { return false }
func (RuneragerSwarmYellow) Play(s *sim.TurnState, self *sim.CardState) {
	runeragerSwarmPlay(s, self)
}

type RuneragerSwarmBlue struct{}

func (RuneragerSwarmBlue) ID() ids.CardID          { return ids.RuneragerSwarmBlue }
func (RuneragerSwarmBlue) Name() string            { return "Runerager Swarm" }
func (RuneragerSwarmBlue) Cost(*sim.TurnState) int { return 0 }
func (RuneragerSwarmBlue) Pitch() int              { return 3 }
func (RuneragerSwarmBlue) Attack() int             { return 1 }
func (RuneragerSwarmBlue) Defense() int            { return 3 }
func (RuneragerSwarmBlue) Types() card.TypeSet     { return runeragerSwarmTypes }
func (RuneragerSwarmBlue) GoAgain() bool           { return false }
func (RuneragerSwarmBlue) Play(s *sim.TurnState, self *sim.CardState) {
	runeragerSwarmPlay(s, self)
}
func runeragerSwarmPlay(s *sim.TurnState, self *sim.CardState) {
	if s.HasPlayedOrCreatedAura() {
		self.GrantedGoAgain = true
	}
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}

func (RuneragerSwarmRed) ConditionalGoAgain()    {}
func (RuneragerSwarmYellow) ConditionalGoAgain() {}
func (RuneragerSwarmBlue) ConditionalGoAgain()   {}
