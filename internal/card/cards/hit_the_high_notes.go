// Hit the High Notes — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you've played or created an aura this turn, this gets +2{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (HitTheHighNotesRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(ge)
}

func (HitTheHighNotesYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(ge)
}

func (HitTheHighNotesBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(ge)
}
func hitTheHighNotesBonus(ge card.GameEngine) int {
	if ge.AuraCreated() {
		return 2
	}
	return 0
}
