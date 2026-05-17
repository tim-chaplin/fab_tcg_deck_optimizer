// Arcanic Spike — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 5, Yellow 4, Blue 3.
// Text: "If you've dealt arcane damage this turn, this gets +2{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// arcaneDamageBonus is the +2{p} gained when the "dealt arcane damage this turn" clause is live.
const arcaneDamageBonus = 2

// arcanicSpikeBonus returns the +2{p} power buff when ArcaneDamageDealt is set, else 0.
func arcanicSpikeBonus(ge card.GameEngine) int {
	if ge != nil && ge.ArcaneDamageDealt() {
		return arcaneDamageBonus
	}
	return 0
}

func (ArcanicSpikeRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += arcanicSpikeBonus(ge)
}

func (ArcanicSpikeYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += arcanicSpikeBonus(ge)
}

func (ArcanicSpikeBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += arcanicSpikeBonus(ge)
}
