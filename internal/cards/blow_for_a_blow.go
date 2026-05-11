// Blow for a Blow — Generic Action - Attack. Cost 2, Pitch 1, Power 4, Defense 2. Only printed in
// Red.
//
// Text: "When this is played, if you have less {h} than an opposing hero, it gets **go again**.
// When this hits, deal 1 damage to any target."
//
// On-hit 1 damage is modelled as +1 damage-equivalent. The "less {h}" go-again clause routes
// through sim.HeroWantsLowerHealth — fires for heroes implementing sim.LowerHealthWanter.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// blowForABlowPingValue is the damage-equivalent credited when the on-hit 1-damage rider fires.
const blowForABlowPingValue = 1

func (BlowForABlowRed) GoAgain() bool { return sim.HeroWantsLowerHealth() }
func (BlowForABlowRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(blowForABlowOnHit)
}

// blowForABlowOnHit fires the printed "When this hits, deal 1 damage" rider. Top-level so
// registration doesn't allocate a closure on the hot anneal path.
func blowForABlowOnHit(s card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	s.AddValue(blowForABlowPingValue)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit dealt 1 damage", blowForABlowPingValue)
}
