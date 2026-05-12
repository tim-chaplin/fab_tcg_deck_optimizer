// Runerager Swarm — Runeblade Action - Attack. Cost 0, Defense 3.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "If you've played or created an aura this turn, this gets go again."
//
// Go again is conditional on the aura clause, not a printed keyword (docs/dev-standards.md
// covers the conditional grant wiring).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (RuneragerSwarmRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	runeragerSwarmPlay(g, l, self)
}

func (RuneragerSwarmYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	runeragerSwarmPlay(g, l, self)
}

func (RuneragerSwarmBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	runeragerSwarmPlay(g, l, self)
}
func runeragerSwarmPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	if g.AuraCreated() {
		self.GrantedGoAgain = true
	}
}
