package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Card-facing aura creation methods on *GameEngine. Cards call these via the
// v2/card.GameEngine interface; the methods delegate to the registered BuildCardAura
// factory so the concrete Aura type lives outside gameengine.

// CreateStartOfTurnAura registers a triggertype.StartOfTurn aura: the handler fires at
// the start of each subsequent turn. The handler signature matches card.GameEngine's
// inline declaration; BuildCardAura's slot accepts aura.Handler structurally.
func (ge *GameEngine) CreateStartOfTurnAura(self *card.CardState, handler func(card.GameEngine, card.Logger, card.Aura), count int) {
	ge.CreateAura(BuildCardAura(self, triggertype.StartOfTurn, handler, count, false))
}

// CreateOncePerTurnAttackActionAura registers a triggertype.AttackAction aura with the
// OncePerTurn gate set — fires at most once per turn regardless of how many attack
// actions resolve.
func (ge *GameEngine) CreateOncePerTurnAttackActionAura(self *card.CardState, handler func(card.GameEngine, card.Logger, card.Aura), count int) {
	ge.CreateAura(BuildCardAura(self, triggertype.AttackAction, handler, count, true))
}
