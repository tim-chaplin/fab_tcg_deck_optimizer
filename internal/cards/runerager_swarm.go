// Runerager Swarm — Runeblade Action - Attack. Cost 0, Defense 3.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "If you've played or created an aura this turn, this gets go again."
//
// Go again is conditional on the aura clause, not a printed keyword (docs/dev-standards.md
// covers the conditional grant wiring).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (RuneragerSwarmRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	runeragerSwarmPlay(s, l, self)
}

func (RuneragerSwarmYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	runeragerSwarmPlay(s, l, self)
}

func (RuneragerSwarmBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	runeragerSwarmPlay(s, l, self)
}
func runeragerSwarmPlay(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	if s.HasPlayedOrCreatedAura() {
		self.GrantedGoAgain = true
	}
}
