// Life for a Life — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this is played, if you have less {h} than an opposing hero, it gets **go again**.
// When this hits, gain 1{h}."
//
// 1{h} gain is modelled as +1 damage-equivalent. The "less {h}" go-again clause routes
// through g.HeroWantsLowerHealth — fires for heroes implementing card.LowerHealthWanter.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// lifeForALifeHealValue is the damage-equivalent credited when the on-hit 1{h} gain fires.
const lifeForALifeHealValue = 1

// lifeForALifeOnHit fires the printed "When this hits, gain 1{h}" rider. Top-level so
// registration stays alloc-free.
func lifeForALifeOnHit(g card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	g.AddValue(lifeForALifeHealValue)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit gained 1 health", lifeForALifeHealValue)
}

// lifeForALifeGoAgain returns the printed "less {h}" go-again rider when the active hero
// opts into LowerHealthWanter. nil-g (Opt-time / cardmeta lookups before a hero is set)
// reads as false — printed default.
func lifeForALifeGoAgain(g card.GameEngine) bool {
	if g == nil {
		return false
	}
	return g.HeroWantsLowerHealth()
}

func (LifeForALifeRed) GoAgain(g card.GameEngine) bool { return lifeForALifeGoAgain(g) }
func (LifeForALifeRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(lifeForALifeOnHit)
}

func (LifeForALifeYellow) GoAgain(g card.GameEngine) bool { return lifeForALifeGoAgain(g) }
func (LifeForALifeYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(lifeForALifeOnHit)
}

func (LifeForALifeBlue) GoAgain(g card.GameEngine) bool { return lifeForALifeGoAgain(g) }
func (LifeForALifeBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(lifeForALifeOnHit)
}
