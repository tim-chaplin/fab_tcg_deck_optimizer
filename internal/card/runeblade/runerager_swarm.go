// Runerager Swarm — Runeblade Action - Attack. Cost 0, Defense 3.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "If you've played or created an aura this turn, this gets go again."
//
// Go again is CONDITIONAL — it's not a printed keyword but a text-granted effect. Play sets
// self.GrantedGoAgain when the aura condition is met so the chain-legality check sees the
// grant.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var runeragerSwarmTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type RuneragerSwarmRed struct{}

func (RuneragerSwarmRed) ID() card.ID                 { return card.RuneragerSwarmRed }
func (RuneragerSwarmRed) Name() string             { return "Runerager Swarm (Red)" }
func (RuneragerSwarmRed) Cost(*card.TurnState) int                { return 0 }
func (RuneragerSwarmRed) Pitch() int               { return 1 }
func (RuneragerSwarmRed) Attack() int              { return 3 }
func (RuneragerSwarmRed) Defense() int             { return 3 }
func (RuneragerSwarmRed) Types() card.TypeSet      { return runeragerSwarmTypes }
func (RuneragerSwarmRed) GoAgain() bool            { return false }
func (c RuneragerSwarmRed) Play(s *card.TurnState, self *card.CardState) int {
	return runeragerSwarmPlay(c.Attack(), s, self)
}

type RuneragerSwarmYellow struct{}

func (RuneragerSwarmYellow) ID() card.ID                 { return card.RuneragerSwarmYellow }
func (RuneragerSwarmYellow) Name() string             { return "Runerager Swarm (Yellow)" }
func (RuneragerSwarmYellow) Cost(*card.TurnState) int                { return 0 }
func (RuneragerSwarmYellow) Pitch() int               { return 2 }
func (RuneragerSwarmYellow) Attack() int              { return 2 }
func (RuneragerSwarmYellow) Defense() int             { return 3 }
func (RuneragerSwarmYellow) Types() card.TypeSet      { return runeragerSwarmTypes }
func (RuneragerSwarmYellow) GoAgain() bool            { return false }
func (c RuneragerSwarmYellow) Play(s *card.TurnState, self *card.CardState) int {
	return runeragerSwarmPlay(c.Attack(), s, self)
}

type RuneragerSwarmBlue struct{}

func (RuneragerSwarmBlue) ID() card.ID                 { return card.RuneragerSwarmBlue }
func (RuneragerSwarmBlue) Name() string             { return "Runerager Swarm (Blue)" }
func (RuneragerSwarmBlue) Cost(*card.TurnState) int                { return 0 }
func (RuneragerSwarmBlue) Pitch() int               { return 3 }
func (RuneragerSwarmBlue) Attack() int              { return 1 }
func (RuneragerSwarmBlue) Defense() int             { return 3 }
func (RuneragerSwarmBlue) Types() card.TypeSet      { return runeragerSwarmTypes }
func (RuneragerSwarmBlue) GoAgain() bool            { return false }
func (c RuneragerSwarmBlue) Play(s *card.TurnState, self *card.CardState) int {
	return runeragerSwarmPlay(c.Attack(), s, self)
}

func runeragerSwarmPlay(base int, s *card.TurnState, self *card.CardState) int {
	if s.HasAuraInPlay() {
		self.GrantedGoAgain = true
	}
	return base
}
