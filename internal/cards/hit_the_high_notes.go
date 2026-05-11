// Hit the High Notes — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you've played or created an aura this turn, this gets +2{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (HitTheHighNotesRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(s)
}

func (HitTheHighNotesYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(s)
}

func (HitTheHighNotesBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(s)
}
func hitTheHighNotesBonus(s card.GameEngine) int {
	if s.HasPlayedOrCreatedAura() {
		return 2
	}
	return 0
}
