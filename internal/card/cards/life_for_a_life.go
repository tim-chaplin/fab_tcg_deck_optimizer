// Life for a Life — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this is played, if you have less {h} than an opposing hero, it gets **go again**.
// When this hits, gain 1{h}."
//
// 1{h} gain is modelled as +1 damage-equivalent. The "less {h}" go-again clause routes
// through ge.HeroWantsLowerHealth — fires for heroes implementing card.LowerHealthWanter.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// lifeForALifeHealValue is the damage-equivalent credited when the on-hit 1{h} gain fires.
const lifeForALifeHealValue = 1

// lifeForALifeOnHit fires the printed "When this hits, gain 1{h}" rider.
func lifeForALifeOnHit(ge card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	ge.AddValue(lifeForALifeHealValue)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit gained 1 health", lifeForALifeHealValue)
}

// lifeForALifeGoAgain returns true when the active hero opts into LowerHealthWanter. nil-g
// reads as false (the printed default).
func lifeForALifeGoAgain(ge card.GameEngine) bool {
	if ge == nil {
		return false
	}
	return ge.HeroWantsLowerHealth()
}

func (LifeForALifeRed) GoAgain(ge card.GameEngine) bool { return lifeForALifeGoAgain(ge) }
func (LifeForALifeRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(lifeForALifeOnHit)
}

func (LifeForALifeYellow) GoAgain(ge card.GameEngine) bool { return lifeForALifeGoAgain(ge) }
func (LifeForALifeYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(lifeForALifeOnHit)
}

func (LifeForALifeBlue) GoAgain(ge card.GameEngine) bool { return lifeForALifeGoAgain(ge) }
func (LifeForALifeBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(lifeForALifeOnHit)
}
