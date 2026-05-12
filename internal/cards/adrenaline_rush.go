// Adrenaline Rush — Generic Action - Attack. Cost 2. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When you play this, if you have less {h} than an opposing hero, this gets +3{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// adrenalineRushBonus returns the +3{p} rider when the current hero opts into
// LowerHealthWanter, else 0.
func adrenalineRushBonus(s card.GameEngine) int {
	if s.HeroWantsLowerHealth() {
		return 3
	}
	return 0
}

func (AdrenalineRushRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += adrenalineRushBonus(s)
}

func (AdrenalineRushYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += adrenalineRushBonus(s)
}

func (AdrenalineRushBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += adrenalineRushBonus(s)
}
