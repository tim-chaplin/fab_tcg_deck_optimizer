// Consuming Volition — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you've dealt arcane damage this turn, this gets 'When this hits a hero, they discard
// a card.'"
//
// "This hits" reads only this card's own damage; co-firing runechants don't satisfy it.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// consumingVolitionApplyRider registers the on-hit discard rider; the ArcaneDamageDealt
// gate runs at hit time so a Runechant firing on this same attack can satisfy it.
func consumingVolitionApplyRider(_ *sim.TurnState, l sim.Logger, self *sim.CardState) {
	self.RegisterOnHit(consumingVolitionOnHit)
}

// consumingVolitionOnHit fires the "When this hits a hero, they discard a card" rider
// when ArcaneDamageDealt is set. Top-level so registration stays alloc-free.
func consumingVolitionOnHit(s *sim.TurnState, l sim.Logger, self *sim.CardState, _ *sim.OnHitHandler) {
	if !s.ArcaneDamageDealt {
		return
	}
	s.AddValue(sim.DiscardValue)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit discarded a card", sim.DiscardValue)
}

func (ConsumingVolitionRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	consumingVolitionApplyRider(s, l, self)
}

func (ConsumingVolitionYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	consumingVolitionApplyRider(s, l, self)
}

func (ConsumingVolitionBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	consumingVolitionApplyRider(s, l, self)
}
