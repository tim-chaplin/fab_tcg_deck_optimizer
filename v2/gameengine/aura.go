package gameengine

import "github.com/tim-chaplin/fab-deck-optimizer/v2/card"

// Card-facing aura creation methods on *GameEngine. Cards call these via the
// v2/card.GameEngine interface; the methods delegate to the registered BuildCardAura
// factory so the concrete Aura type lives outside gameengine.

// CreateStartOfTurnAura registers a TriggerStartOfTurn aura: the handler fires at the
// start of each subsequent turn.
func (g *GameEngine) CreateStartOfTurnAura(self *card.CardState, handler card.AuraHandler, count int) {
	g.CreateAura(BuildCardAura(self, TriggerStartOfTurn, handler, count, false))
}

// CreateOncePerTurnAttackActionAura registers a TriggerAttackAction aura with the
// OncePerTurn gate set — fires at most once per turn regardless of how many attack
// actions resolve.
func (g *GameEngine) CreateOncePerTurnAttackActionAura(self *card.CardState, handler card.AuraHandler, count int) {
	g.CreateAura(BuildCardAura(self, TriggerAttackAction, handler, count, true))
}
