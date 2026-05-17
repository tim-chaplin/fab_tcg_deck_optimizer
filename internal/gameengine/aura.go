package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Card-facing aura creation methods on *GameEngine. Cards call these via the
// internal/card.GameEngine interface; the methods construct the concrete *aura.Aura directly.
//
// This file is also the junction that proves internal/card.GameEngine's inline handler type
// (func(GameEngine, Logger, Aura)) agrees with internal/aura.Handler. internal/card and internal/aura
// can't import each other (internal/aura imports internal/card; the reverse would cycle), so they
// each declare the shape independently. The aura.NewFromCard call below assigns the
// inline-typed handler to aura.Handler — if the two signatures ever drift, this file
// stops compiling. The same pattern applies to internal/trigger via engine.go's AddHitTrigger
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
