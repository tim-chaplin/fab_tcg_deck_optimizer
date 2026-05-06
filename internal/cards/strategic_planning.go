// Strategic Planning — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Put an action card with cost 2 or less from a graveyard on the bottom of its owner's deck.
// At the beginning of the end phase, draw a card. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var strategicPlanningTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// drawOneAtEndOfTurn is the end-of-turn TriggerHandler that fires Strategic Planning's
// deferred draw. Top-level so the registration stays alloc-free.
func drawOneAtEndOfTurn(s *sim.TurnState, _ *sim.Trigger) {
	s.DrawOne()
}

func strategicPlanningPlay(s *sim.TurnState, self *sim.CardState) {
	if _, ok := s.RecycleFromGraveyardToBottom(func(c sim.Card) bool {
		return c.Types().Has(card.TypeAction) && c.Cost(s) <= 2
	}); ok {
		s.LogRider(self, 0, "Recycled an action card to bottom of deck")
	}
	s.AddTrigger(sim.Trigger{
		Source:      self.Card,
		TriggerType: sim.TriggerEndOfTurn,
		Handler:     drawOneAtEndOfTurn,
	})
	s.LogRider(self, 0, "End-phase draw queued")
	s.Log(self, 0)
}

type StrategicPlanningRed struct{}

func (StrategicPlanningRed) ID() ids.CardID          { return ids.StrategicPlanningRed }
func (StrategicPlanningRed) Name() string            { return "Strategic Planning" }
func (StrategicPlanningRed) Cost(*sim.TurnState) int { return 1 }
func (StrategicPlanningRed) Pitch() int              { return 1 }
func (StrategicPlanningRed) Attack() int             { return 0 }
func (StrategicPlanningRed) Defense() int            { return 2 }
func (StrategicPlanningRed) Types() card.TypeSet     { return strategicPlanningTypes }
func (StrategicPlanningRed) GoAgain() bool           { return true }
func (StrategicPlanningRed) Play(s *sim.TurnState, self *sim.CardState) {
	strategicPlanningPlay(s, self)
}

type StrategicPlanningYellow struct{}

func (StrategicPlanningYellow) ID() ids.CardID          { return ids.StrategicPlanningYellow }
func (StrategicPlanningYellow) Name() string            { return "Strategic Planning" }
func (StrategicPlanningYellow) Cost(*sim.TurnState) int { return 1 }
func (StrategicPlanningYellow) Pitch() int              { return 2 }
func (StrategicPlanningYellow) Attack() int             { return 0 }
func (StrategicPlanningYellow) Defense() int            { return 2 }
func (StrategicPlanningYellow) Types() card.TypeSet     { return strategicPlanningTypes }
func (StrategicPlanningYellow) GoAgain() bool           { return true }
func (StrategicPlanningYellow) Play(s *sim.TurnState, self *sim.CardState) {
	strategicPlanningPlay(s, self)
}

type StrategicPlanningBlue struct{}

func (StrategicPlanningBlue) ID() ids.CardID          { return ids.StrategicPlanningBlue }
func (StrategicPlanningBlue) Name() string            { return "Strategic Planning" }
func (StrategicPlanningBlue) Cost(*sim.TurnState) int { return 1 }
func (StrategicPlanningBlue) Pitch() int              { return 3 }
func (StrategicPlanningBlue) Attack() int             { return 0 }
func (StrategicPlanningBlue) Defense() int            { return 2 }
func (StrategicPlanningBlue) Types() card.TypeSet     { return strategicPlanningTypes }
func (StrategicPlanningBlue) GoAgain() bool           { return true }
func (StrategicPlanningBlue) Play(s *sim.TurnState, self *sim.CardState) {
	strategicPlanningPlay(s, self)
}
