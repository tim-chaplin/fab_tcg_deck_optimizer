// Hit the High Notes — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you've played or created an aura this turn, this gets +2{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (HitTheHighNotesRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(s)
}

func (HitTheHighNotesYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(s)
}

func (HitTheHighNotesBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	self.BonusAttack += hitTheHighNotesBonus(s)
}
func hitTheHighNotesBonus(s sim.GameEngine) int {
	if s.HasPlayedOrCreatedAura() {
		return 2
	}
	return 0
}
