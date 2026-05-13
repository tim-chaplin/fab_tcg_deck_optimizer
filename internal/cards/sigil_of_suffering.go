// Sigil of Suffering — Runeblade Defense Reaction. Cost 0, Arcane 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 3, Yellow 2, Blue 1.
// Text: "Deal 1 arcane damage to the attacking hero. If you have dealt arcane damage this turn,
// Sigil of Suffering gains +1{d}."
//
// The Sigil's own printed-1 arcane satisfies the conditional via LikelyDamageHits(1, false),
// so the +1{d} is credited whenever there's IncomingDamage left to absorb it.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func sigilOfSufferingPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if ge.ArcaneDamageDealt() || ge.LikelyDamageHits(1, false) {
		self.BonusDefense++
	}
	ge.DealArcaneDamage(l, self.Card.DisplayName(), 1)
}

func (SigilOfSufferingRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	sigilOfSufferingPlay(ge, l, self)
}

func (SigilOfSufferingYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	sigilOfSufferingPlay(ge, l, self)
}

func (SigilOfSufferingBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	sigilOfSufferingPlay(ge, l, self)
}
