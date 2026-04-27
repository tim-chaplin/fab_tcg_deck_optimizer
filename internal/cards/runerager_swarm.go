// Runerager Swarm — Runeblade Action - Attack. Cost 0, Defense 3.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "If you've played or created an aura this turn, this gets go again."
//
// Go again is conditional on the aura clause, not a printed keyword (docs/dev-standards.md
// covers the conditional grant wiring).

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var runeragerSwarmTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type RuneragerSwarmRed struct{}

func (RuneragerSwarmRed) ID() card.ID              { return card.RuneragerSwarmRed }
func (RuneragerSwarmRed) Name() string             { return "Runerager Swarm" }
func (RuneragerSwarmRed) Cost(*card.TurnState) int { return 0 }
func (RuneragerSwarmRed) Pitch() int               { return 1 }
func (RuneragerSwarmRed) Attack() int              { return 3 }
func (RuneragerSwarmRed) Defense() int             { return 3 }
func (RuneragerSwarmRed) Types() card.TypeSet      { return runeragerSwarmTypes }
func (RuneragerSwarmRed) GoAgain() bool            { return false }
func (RuneragerSwarmRed) Play(s *card.TurnState, self *card.CardState) {
	runeragerSwarmPlay(s, self)
}

type RuneragerSwarmYellow struct{}

func (RuneragerSwarmYellow) ID() card.ID              { return card.RuneragerSwarmYellow }
func (RuneragerSwarmYellow) Name() string             { return "Runerager Swarm" }
func (RuneragerSwarmYellow) Cost(*card.TurnState) int { return 0 }
func (RuneragerSwarmYellow) Pitch() int               { return 2 }
func (RuneragerSwarmYellow) Attack() int              { return 2 }
func (RuneragerSwarmYellow) Defense() int             { return 3 }
func (RuneragerSwarmYellow) Types() card.TypeSet      { return runeragerSwarmTypes }
func (RuneragerSwarmYellow) GoAgain() bool            { return false }
func (RuneragerSwarmYellow) Play(s *card.TurnState, self *card.CardState) {
	runeragerSwarmPlay(s, self)
}

type RuneragerSwarmBlue struct{}

func (RuneragerSwarmBlue) ID() card.ID              { return card.RuneragerSwarmBlue }
func (RuneragerSwarmBlue) Name() string             { return "Runerager Swarm" }
func (RuneragerSwarmBlue) Cost(*card.TurnState) int { return 0 }
func (RuneragerSwarmBlue) Pitch() int               { return 3 }
func (RuneragerSwarmBlue) Attack() int              { return 1 }
func (RuneragerSwarmBlue) Defense() int             { return 3 }
func (RuneragerSwarmBlue) Types() card.TypeSet      { return runeragerSwarmTypes }
func (RuneragerSwarmBlue) GoAgain() bool            { return false }
func (RuneragerSwarmBlue) Play(s *card.TurnState, self *card.CardState) {
	runeragerSwarmPlay(s, self)
}
func runeragerSwarmPlay(s *card.TurnState, self *card.CardState) {
	if s.HasPlayedOrCreatedAura() {
		self.GrantedGoAgain = true
	}
	s.ApplyAndLogEffectiveAttack(self)
}
