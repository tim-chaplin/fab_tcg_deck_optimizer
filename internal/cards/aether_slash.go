// Aether Slash — Runeblade Action - Attack. Cost 1, Defense 3, Arcane 1.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When Aether Slash attacks, if a 'non-attack' action card was pitched to play it, deal 1
// arcane damage to any target."
//
// Reads self.PitchedToPlay (the cards the chain runner attributed to funding THIS copy's
// cost) to gate the +1 arcane rider.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (AetherSlashRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	aetherSlashApplyRider(g, l, self)
}

func (AetherSlashYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	aetherSlashApplyRider(g, l, self)
}

func (AetherSlashBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	aetherSlashApplyRider(g, l, self)
}

// aetherSlashApplyRider deals 1 arcane and emits the rider sub-line when a non-attack action
// is among the pitched cards the runner attributed to paying for this Aether Slash.
func aetherSlashApplyRider(g card.GameEngine, l card.Logger, self *card.CardState) {
	for _, p := range self.PitchedToPlay {
		if p.Types(nil).IsNonAttackAction() {
			g.DealArcaneDamage(l, self.Card.DisplayName(), 1)
			return
		}
	}
}
