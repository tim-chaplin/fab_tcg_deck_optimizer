// Life for a Life — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this is played, if you have less {h} than an opposing hero, it gets **go again**.
// When this hits, gain 1{h}."
//
// The on-hit 1{h} gain is modelled as +1 damage-equivalent (1 health saved ≈ 1 damage), gated
// on sim.LikelyToHit. The "less {h} than an opposing hero" clause is modelled as a hero
// attribute — go again fires for heroes that implement sim.LowerHealthWanter and never fires
// otherwise.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var lifeForALifeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// lifeForALifeHealValue is the damage-equivalent credited when the on-hit 1{h} gain fires.
const lifeForALifeHealValue = 1

type LifeForALifeRed struct{}

func (LifeForALifeRed) ID() ids.CardID          { return ids.LifeForALifeRed }
func (LifeForALifeRed) Name() string            { return "Life for a Life" }
func (LifeForALifeRed) Cost(*sim.TurnState) int { return 1 }
func (LifeForALifeRed) Pitch() int              { return 1 }
func (LifeForALifeRed) Attack() int             { return 4 }
func (LifeForALifeRed) Defense() int            { return 2 }
func (LifeForALifeRed) Types() card.TypeSet     { return lifeForALifeTypes }
func (LifeForALifeRed) GoAgain() bool           { return sim.HeroWantsLowerHealth() }
func (LifeForALifeRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	s.ApplyAndLogRiderOnHit(self, lifeForALifeHealValue, "On-hit gained 1 health")
}

type LifeForALifeYellow struct{}

func (LifeForALifeYellow) ID() ids.CardID          { return ids.LifeForALifeYellow }
func (LifeForALifeYellow) Name() string            { return "Life for a Life" }
func (LifeForALifeYellow) Cost(*sim.TurnState) int { return 1 }
func (LifeForALifeYellow) Pitch() int              { return 2 }
func (LifeForALifeYellow) Attack() int             { return 3 }
func (LifeForALifeYellow) Defense() int            { return 2 }
func (LifeForALifeYellow) Types() card.TypeSet     { return lifeForALifeTypes }
func (LifeForALifeYellow) GoAgain() bool           { return sim.HeroWantsLowerHealth() }
func (LifeForALifeYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	s.ApplyAndLogRiderOnHit(self, lifeForALifeHealValue, "On-hit gained 1 health")
}

type LifeForALifeBlue struct{}

func (LifeForALifeBlue) ID() ids.CardID          { return ids.LifeForALifeBlue }
func (LifeForALifeBlue) Name() string            { return "Life for a Life" }
func (LifeForALifeBlue) Cost(*sim.TurnState) int { return 1 }
func (LifeForALifeBlue) Pitch() int              { return 3 }
func (LifeForALifeBlue) Attack() int             { return 2 }
func (LifeForALifeBlue) Defense() int            { return 2 }
func (LifeForALifeBlue) Types() card.TypeSet     { return lifeForALifeTypes }
func (LifeForALifeBlue) GoAgain() bool           { return sim.HeroWantsLowerHealth() }
func (LifeForALifeBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	s.ApplyAndLogRiderOnHit(self, lifeForALifeHealValue, "On-hit gained 1 health")
}
