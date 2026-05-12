// Wounded Bull — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When you play this, if you have less {h} than an opposing hero, this gains +1{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// woundedBullBonus returns the +1{p} power buff when the current hero opts into
// LowerHealthWanter, else 0.
func woundedBullBonus(s card.GameEngine) int {
	if s.HeroWantsLowerHealth() {
		return 1
	}
	return 0
}

func (WoundedBullRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += woundedBullBonus(s)
}

func (WoundedBullYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += woundedBullBonus(s)
}

func (WoundedBullBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += woundedBullBonus(s)
}
