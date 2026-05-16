package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Card-facing aura creation methods on *GameEngine. Cards call these via the
// v2/card.GameEngine interface; the methods construct the concrete *aura.Aura directly.
//
// This file is also the junction that proves v2/card.GameEngine's inline handler type
// (func(GameEngine, Logger, Aura)) agrees with v2/aura.Handler. v2/card and v2/aura
// can't import each other (v2/aura imports v2/card; the reverse would cycle), so they
// each declare the shape independently. The aura.NewFromCard call below assigns the
// inline-typed handler to aura.Handler — if the two signatures ever drift, this file
// stops compiling. The same pattern applies to v2/trigger via engine.go's AddHitTrigger
// and AddEndOfTurnTrigger.

// CreateStartOfTurnAura registers a triggertype.StartOfTurn aura: the handler fires at
// the start of each subsequent turn. The handler signature matches card.GameEngine's
// inline declaration; aura.NewFromCard's fire parameter accepts it as aura.Handler
// (same underlying func type).
func (ge *GameEngine) CreateStartOfTurnAura(self *card.CardState, handler func(card.GameEngine, card.Logger, card.Aura), count int) {
	ge.CreateAura(aura.NewFromCard(self.Card, triggertype.StartOfTurn, handler, count, false))
}

// CreateOncePerTurnAttackActionAura registers a triggertype.AttackAction aura with the
// OncePerTurn gate set — fires at most once per turn regardless of how many attack
// actions resolve.
func (ge *GameEngine) CreateOncePerTurnAttackActionAura(self *card.CardState, handler func(card.GameEngine, card.Logger, card.Aura), count int) {
	ge.CreateAura(aura.NewFromCard(self.Card, triggertype.AttackAction, handler, count, true))
}
