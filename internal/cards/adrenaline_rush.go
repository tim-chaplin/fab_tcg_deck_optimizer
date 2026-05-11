// Adrenaline Rush — Generic Action - Attack. Cost 2. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When you play this, if you have less {h} than an opposing hero, this gets +3{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// adrenalineRushBonus returns the +3{p} rider when the current hero opts into
// LowerHealthWanter, else 0.
func adrenalineRushBonus() int {
	if sim.HeroWantsLowerHealth() {
		return 3
	}
	return 0
}

func (AdrenalineRushRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	self.BonusAttack += adrenalineRushBonus()
}

func (AdrenalineRushYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	self.BonusAttack += adrenalineRushBonus()
}

func (AdrenalineRushBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	self.BonusAttack += adrenalineRushBonus()
}
