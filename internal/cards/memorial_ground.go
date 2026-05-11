// Memorial Ground — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Put target attack action card with cost 2 or less from your graveyard on top of your
// deck."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Cost predicate reads s so variable-cost targets are gated on their current cost.
func memorialGroundPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	if _, ok := s.RecycleFromGraveyardToTop(func(c sim.Card) bool {
		return c.Types().IsAttackAction() && c.Cost(s) <= 2
	}); ok {
		self.LogRider(l, 0, "Recycled an attack action card to top of deck")
	}
	self.Log(l, 0)
}

func (MemorialGroundRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	memorialGroundPlay(s, l, self)
}

func (MemorialGroundYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	memorialGroundPlay(s, l, self)
}

func (MemorialGroundBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	memorialGroundPlay(s, l, self)
}
