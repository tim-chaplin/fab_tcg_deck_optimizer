// Tremor of íArathael — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If a card has been put into your banished zone this turn, Tremor of íArathael gains
// +2{p}."
//
// Snapshot at Play: only banishes earlier in the chain trigger the +2{p}; later banishes
// don't retroactively buff this attack (matches the past-tense "has been put").

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func tremorOfIArathaelPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if ge.CardBanished() {
		self.BonusAttack += 2
	}
}

func (TremorOfIArathaelRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	tremorOfIArathaelPlay(ge, l, self)
}

func (TremorOfIArathaelYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	tremorOfIArathaelPlay(ge, l, self)
}

func (TremorOfIArathaelBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	tremorOfIArathaelPlay(ge, l, self)
}
