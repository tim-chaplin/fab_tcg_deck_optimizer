// Aether Slash — Runeblade Action - Attack. Cost 1, Defense 3, Arcane 1.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When Aether Slash attacks, if a 'non-attack' action card was pitched to play it, deal 1
// arcane damage to any target."
//
// Reads self.PitchedToPlay to gate the +1 arcane rider.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (AetherSlashRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	aetherSlashApplyRider(ge, l, self)
}

func (AetherSlashYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	aetherSlashApplyRider(ge, l, self)
}

func (AetherSlashBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	aetherSlashApplyRider(ge, l, self)
}

// aetherSlashApplyRider deals 1 arcane when a non-attack action is among the pitched cards.
func aetherSlashApplyRider(ge card.GameEngine, l card.Logger, self *card.CardState) {
	for _, p := range self.PitchedToPlay {
		if p.Types(nil).IsNonAttackAction() {
			ge.DealArcaneDamage(l, self.Card.DisplayName(), 1)
			return
		}
	}
}
