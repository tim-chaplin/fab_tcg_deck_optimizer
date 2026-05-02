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
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var blowForABlowTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// blowForABlowPingValue is the damage-equivalent credited when the on-hit 1-damage rider fires.
const blowForABlowPingValue = 1

type BlowForABlowRed struct{}

func (BlowForABlowRed) ID() ids.CardID          { return ids.BlowForABlowRed }
func (BlowForABlowRed) Name() string            { return "Blow for a Blow" }
func (BlowForABlowRed) Cost(*sim.TurnState) int { return 2 }
func (BlowForABlowRed) Pitch() int              { return 1 }
func (BlowForABlowRed) Attack() int             { return 4 }
func (BlowForABlowRed) Defense() int            { return 2 }
func (BlowForABlowRed) Types() card.TypeSet     { return blowForABlowTypes }
func (BlowForABlowRed) GoAgain() bool           { return sim.HeroWantsLowerHealth() }
func (BlowForABlowRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	self.OnHit = append(self.OnHit, sim.OnHitHandler{Fire: blowForABlowOnHit})
}

// blowForABlowOnHit fires the printed "When this hits, deal 1 damage" rider. Top-level so
// registration doesn't allocate a closure on the hot anneal path.
func blowForABlowOnHit(s *sim.TurnState, self *sim.CardState, _ *sim.OnHitHandler) {
	s.AddValue(blowForABlowPingValue)
	s.LogRider(self, blowForABlowPingValue, "On-hit dealt 1 damage")
}
