package gameengine

import "github.com/tim-chaplin/fab-deck-optimizer/v2/card"

// Card-facing trigger registration on *GameEngine. Cards call these via
// v2/card.GameEngine; the methods delegate to the registered BuildCardTrigger factory so
// the concrete Trigger type lives outside gameengine.

// AddHitTrigger registers a one-shot TriggerHit listener. filter narrows the qualifying
// hits to a card-type predicate; nil = any hit qualifies.
func (g *GameEngine) AddHitTrigger(self *card.CardState, handler card.TriggerHandler, filter func(card.TypeSet) bool) {
	g.CreateTrigger(BuildCardTrigger(self, TriggerHit, handler, filter))
}

// AddEndOfTurnTrigger registers a one-shot TriggerEndOfTurn listener — fires after the
// chain finishes resolving but before the carry-state snapshot.
func (g *GameEngine) AddEndOfTurnTrigger(self *card.CardState, handler card.TriggerHandler) {
	g.CreateTrigger(BuildCardTrigger(self, TriggerEndOfTurn, handler, nil))
}
