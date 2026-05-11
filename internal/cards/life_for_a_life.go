// Life for a Life — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this is played, if you have less {h} than an opposing hero, it gets **go again**.
// When this hits, gain 1{h}."
//
// 1{h} gain is modelled as +1 damage-equivalent. The "less {h}" go-again clause routes
// through sim.HeroWantsLowerHealth — fires for heroes implementing sim.LowerHealthWanter.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// lifeForALifeHealValue is the damage-equivalent credited when the on-hit 1{h} gain fires.
const lifeForALifeHealValue = 1

// lifeForALifeOnHit fires the printed "When this hits, gain 1{h}" rider. Top-level so
// registration stays alloc-free.
func lifeForALifeOnHit(s card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	s.AddValue(lifeForALifeHealValue)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit gained 1 health", lifeForALifeHealValue)
}

func (LifeForALifeRed) GoAgain() bool { return sim.HeroWantsLowerHealth() }
func (LifeForALifeRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(lifeForALifeOnHit)
}

func (LifeForALifeYellow) GoAgain() bool { return sim.HeroWantsLowerHealth() }
func (LifeForALifeYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(lifeForALifeOnHit)
}

func (LifeForALifeBlue) GoAgain() bool { return sim.HeroWantsLowerHealth() }
func (LifeForALifeBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(lifeForALifeOnHit)
}
