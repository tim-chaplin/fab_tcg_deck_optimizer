// Strategic Planning — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Put an action card with cost 2 or less from a graveyard on the bottom of its owner's deck.
// At the beginning of the end phase, draw a card. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// drawOneAtEndOfTurn is the end-of-turn TriggerHandler that fires Strategic Planning's
// deferred draw.
func drawOneAtEndOfTurn(g card.GameEngine, l card.Logger, _ card.Trigger) {
	g.DrawOne()
}

func strategicPlanningPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	if _, ok := g.RecycleFromGraveyardToBottom(func(c card.Card) bool {
		return c.Types(nil).Has(card.TypeAction) && c.Cost(g) <= 2
	}); ok {
		l.AppendPostTrigger(self.Card.DisplayName(), "Recycled an action card to bottom of deck", 0)
	}
	g.AddEndOfTurnTrigger(self, drawOneAtEndOfTurn)
	l.AppendPostTrigger(self.Card.DisplayName(), "End-phase draw queued", 0)
}

func (StrategicPlanningRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	strategicPlanningPlay(g, l, self)
}

func (StrategicPlanningYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	strategicPlanningPlay(g, l, self)
}

func (StrategicPlanningBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	strategicPlanningPlay(g, l, self)
}
