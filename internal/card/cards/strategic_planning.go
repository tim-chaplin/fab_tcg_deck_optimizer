// Strategic Planning — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Put an action card with cost 2 or less from a graveyard on the bottom of its owner's deck.
// At the beginning of the end phase, draw a card. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// drawOneAtEndOfTurn is the end-of-turn TriggerHandler that fires Strategic Planning's
// deferred draw.
func drawOneAtEndOfTurn(ge card.GameEngine, l card.Logger, _ card.EphemeralTrigger, _ triggertype.Type) {
	ge.DrawOne()
}

func strategicPlanningPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if _, ok := ge.RecycleFromGraveyardToBottom(func(c card.Card) bool {
		return c.Types(nil).Has(card.TypeAction) && c.Cost() <= 2
	}); ok {
		l.AppendPostTrigger(self.Card.DisplayName(), "Recycled an action card to bottom of deck", 0)
	}
	ge.CreateTrigger(self.Card, triggertype.EndOfTurn, drawOneAtEndOfTurn, nil)
	l.AppendPostTrigger(self.Card.DisplayName(), "End-phase draw queued", 0)
}

func (StrategicPlanningRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	strategicPlanningPlay(ge, l, self)
}

func (StrategicPlanningYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	strategicPlanningPlay(ge, l, self)
}

func (StrategicPlanningBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	strategicPlanningPlay(ge, l, self)
}
