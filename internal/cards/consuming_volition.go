// Consuming Volition — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you've dealt arcane damage this turn, this gets 'When this hits a hero, they discard
// a card.'"
//
// "This hits" reads only this card's own damage; co-firing runechants don't satisfy it.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var consumingVolitionTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// consumingVolitionApplyRider registers the on-hit discard rider when ArcaneDamageDealt is
// set.
func consumingVolitionApplyRider(s *sim.TurnState, self *sim.CardState) {
	if !s.ArcaneDamageDealt {
		return
	}
	self.OnHit = append(self.OnHit, sim.OnHitHandler{Fire: consumingVolitionOnHit})
}

// consumingVolitionOnHit fires the conditional "When this hits a hero, they discard a card"
// rider. Top-level so registration stays alloc-free.
func consumingVolitionOnHit(s *sim.TurnState, self *sim.CardState, _ *sim.OnHitHandler) {
	s.AddValue(sim.DiscardValue)
	s.LogRider(self, sim.DiscardValue, "On-hit discarded a card")
}

type ConsumingVolitionRed struct{}

func (ConsumingVolitionRed) ID() ids.CardID          { return ids.ConsumingVolitionRed }
func (ConsumingVolitionRed) Name() string            { return "Consuming Volition" }
func (ConsumingVolitionRed) Cost(*sim.TurnState) int { return 1 }
func (ConsumingVolitionRed) Pitch() int              { return 1 }
func (ConsumingVolitionRed) Attack() int             { return 4 }
func (ConsumingVolitionRed) Defense() int            { return 3 }
func (ConsumingVolitionRed) Types() card.TypeSet     { return consumingVolitionTypes }
func (ConsumingVolitionRed) GoAgain() bool           { return false }
func (ConsumingVolitionRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	consumingVolitionApplyRider(s, self)
}

type ConsumingVolitionYellow struct{}

func (ConsumingVolitionYellow) ID() ids.CardID          { return ids.ConsumingVolitionYellow }
func (ConsumingVolitionYellow) Name() string            { return "Consuming Volition" }
func (ConsumingVolitionYellow) Cost(*sim.TurnState) int { return 1 }
func (ConsumingVolitionYellow) Pitch() int              { return 2 }
func (ConsumingVolitionYellow) Attack() int             { return 3 }
func (ConsumingVolitionYellow) Defense() int            { return 3 }
func (ConsumingVolitionYellow) Types() card.TypeSet     { return consumingVolitionTypes }
func (ConsumingVolitionYellow) GoAgain() bool           { return false }
func (ConsumingVolitionYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	consumingVolitionApplyRider(s, self)
}

type ConsumingVolitionBlue struct{}

func (ConsumingVolitionBlue) ID() ids.CardID          { return ids.ConsumingVolitionBlue }
func (ConsumingVolitionBlue) Name() string            { return "Consuming Volition" }
func (ConsumingVolitionBlue) Cost(*sim.TurnState) int { return 1 }
func (ConsumingVolitionBlue) Pitch() int              { return 3 }
func (ConsumingVolitionBlue) Attack() int             { return 2 }
func (ConsumingVolitionBlue) Defense() int            { return 3 }
func (ConsumingVolitionBlue) Types() card.TypeSet     { return consumingVolitionTypes }
func (ConsumingVolitionBlue) GoAgain() bool           { return false }
func (ConsumingVolitionBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	consumingVolitionApplyRider(s, self)
}
