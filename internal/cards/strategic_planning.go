// Strategic Planning — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Put an action card with cost 2 or less from a graveyard on the bottom of its owner's deck.
// At the beginning of the end phase, draw a card. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// drawOneAtEndOfTurn is the end-of-turn TriggerHandler that fires Strategic Planning's
// deferred draw. Top-level so the registration stays alloc-free.
func drawOneAtEndOfTurn(s *sim.TurnState, l sim.Logger, _ *sim.Trigger, _ *sim.Aura) {
	s.DrawOne()
}

func strategicPlanningPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	if _, ok := s.RecycleFromGraveyardToBottom(func(c sim.Card) bool {
		return c.Types().Has(card.TypeAction) && c.Cost(s) <= 2
	}); ok {
		l.AppendPostTrigger(self.Card.DisplayName(), "Recycled an action card to bottom of deck", 0)
	}
	s.AddTrigger(sim.Trigger{
		Source:      self.Card,
		TriggerType: sim.TriggerEndOfTurn,
		Handler:     drawOneAtEndOfTurn,
	})
	l.AppendPostTrigger(self.Card.DisplayName(), "End-phase draw queued", 0)
}

func (StrategicPlanningRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	strategicPlanningPlay(s, l, self)
}

func (StrategicPlanningYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	strategicPlanningPlay(s, l, self)
}

func (StrategicPlanningBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	strategicPlanningPlay(s, l, self)
}
