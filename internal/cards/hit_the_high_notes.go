// Hit the High Notes — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you've played or created an aura this turn, this gets +2{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (HitTheHighNotesRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(g)
}

func (HitTheHighNotesYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(g)
}

func (HitTheHighNotesBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(g)
}
func hitTheHighNotesBonus(g card.GameEngine) int {
	if g.AuraCreated() {
		return 2
	}
	return 0
}
