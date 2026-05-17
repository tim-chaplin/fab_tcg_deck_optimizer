// Blow for a Blow — Generic Action - Attack. Cost 2, Pitch 1, Power 4, Defense 2. Only printed in
// Red.
//
// Text: "When this is played, if you have less {h} than an opposing hero, it gets **go again**.
// When this hits, deal 1 damage to any target."
//
// On-hit 1 damage is modelled as +1 damage-equivalent. The "less {h}" go-again clause routes
// through ge.HeroWantsLowerHealth — fires for heroes implementing card.LowerHealthWanter.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// blowForABlowPingValue is the damage-equivalent credited when the on-hit 1-damage rider fires.
const blowForABlowPingValue = 1

// BlowForABlowRed.GoAgain returns true when the active hero opts into LowerHealthWanter.
// nil-g reads as false (the printed default).
func (BlowForABlowRed) GoAgain(ge card.GameEngine) bool {
	if ge == nil {
		return false
	}
	return ge.HeroWantsLowerHealth()
}
func (BlowForABlowRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(blowForABlowOnHit)
}

// blowForABlowOnHit fires the printed "When this hits, deal 1 damage" rider.
func blowForABlowOnHit(ge card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	ge.AddValue(blowForABlowPingValue)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit dealt 1 damage", blowForABlowPingValue)
}
