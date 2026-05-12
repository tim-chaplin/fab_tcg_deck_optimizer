// Tremor of íArathael — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If a card has been put into your banished zone this turn, Tremor of íArathael gains
// +2{p}."
//
// Snapshot at Play: only banishes earlier in the chain trigger the +2{p}; later banishes
// don't retroactively buff this attack (matches the past-tense "has been put"). Reads the
// CardBanished flag rather than len(g.Banished()) so prior-turn entries that carry
// through the snapshot path don't pollute the check.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func tremorOfIArathaelPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	if g.CardBanished() {
		self.BonusAttack += 2
	}
}

func (TremorOfIArathaelRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	tremorOfIArathaelPlay(g, l, self)
}

func (TremorOfIArathaelYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	tremorOfIArathaelPlay(g, l, self)
}

func (TremorOfIArathaelBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	tremorOfIArathaelPlay(g, l, self)
}
