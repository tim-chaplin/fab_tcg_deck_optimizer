// Blow for a Blow — Generic Action - Attack. Cost 2, Pitch 1, Power 4, Defense 2. Only printed in
// Red.
//
// Text: "When this is played, if you have less {h} than an opposing hero, it gets **go again**.
// When this hits, deal 1 damage to any target."
//
// On-hit 1 damage is modelled as +1 damage-equivalent. The "less {h}" go-again clause routes
// through g.HeroWantsLowerHealth — fires for heroes implementing card.LowerHealthWanter.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// blowForABlowPingValue is the damage-equivalent credited when the on-hit 1-damage rider fires.
const blowForABlowPingValue = 1

// BlowForABlowRed.GoAgain returns the printed "less {h}" go-again rider when the active
// hero opts into LowerHealthWanter. nil-g (Opt-time / cardmeta lookups before a hero is
// set) reads as false — printed default.
func (BlowForABlowRed) GoAgain(g card.GameEngine) bool {
	if g == nil {
		return false
	}
	return g.HeroWantsLowerHealth()
}
func (BlowForABlowRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(blowForABlowOnHit)
}

// blowForABlowOnHit fires the printed "When this hits, deal 1 damage" rider. Top-level so
// registration doesn't allocate a closure on the hot anneal path.
func blowForABlowOnHit(g card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	g.AddValue(blowForABlowPingValue)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit dealt 1 damage", blowForABlowPingValue)
}
