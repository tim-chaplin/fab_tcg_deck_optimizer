// Shrill of Skullform — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you have played or created an aura this turn, Shrill of Skullform gains +3{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (ShrillOfSkullformRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	shrillPlay(g, l, self)
}

func (ShrillOfSkullformYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	shrillPlay(g, l, self)
}

func (ShrillOfSkullformBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	shrillPlay(g, l, self)
}

// shrillPlay routes the +3{p} aura-in-play buff through self.BonusAttack so EffectiveAttack
// and LikelyToHit see the buffed power, then emits the chain step at the buffed value. No
// rider sub-line — this is a power buff, not a separable effect.
func shrillPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	if g.AuraCreated() {
		self.BonusAttack += 3
	}
}
