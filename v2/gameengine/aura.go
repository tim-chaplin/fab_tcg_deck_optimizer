package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Card-facing aura creation methods on *GameEngine. Cards call these via the
// v2/card.GameEngine interface; the methods construct the concrete *aura.Aura directly,
// wrapping the typed card.AuraHandler into the leaf package's type-erased Handler shape.

// CreateStartOfTurnAura registers a triggertype.StartOfTurn aura: the handler fires at
// the start of each subsequent turn.
func (ge *GameEngine) CreateStartOfTurnAura(self *card.CardState, handler card.AuraHandler, count int) {
	ge.CreateAura(aura.NewFromCard(self.Card, triggertype.StartOfTurn, WrapAuraHandler(handler), count, false))
}

// CreateOncePerTurnAttackActionAura registers a triggertype.AttackAction aura with the
// OncePerTurn gate set — fires at most once per turn regardless of how many attack
// actions resolve.
func (ge *GameEngine) CreateOncePerTurnAttackActionAura(self *card.CardState, handler card.AuraHandler, count int) {
	ge.CreateAura(aura.NewFromCard(self.Card, triggertype.AttackAction, WrapAuraHandler(handler), count, true))
}
