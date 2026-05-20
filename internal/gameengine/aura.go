package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Card-facing aura creation methods on *GameEngine. Cards call these via the
// internal/card.GameEngine interface; the methods construct the concrete *aura.Aura directly.
//
// Compile-time bridge: assigning the inline handler type to aura.Handler below makes the
// two independently-declared func types check-equal — drift makes this file stop compiling.

// CreateStartOfTurnAura registers a triggertype.StartOfTurn aura: the handler fires at
// the start of each subsequent turn.
func (ge *GameEngine) CreateStartOfTurnAura(pc *card.CardState, handler func(card.GameEngine, card.Logger, card.Aura), count int) {
	ge.CreateAura(aura.NewFromCard(pc.Card, triggertype.StartOfTurn, handler, count, false, nil))
}

// CreateOncePerTurnAttackActionAura registers a triggertype.CardOrAbility aura filtered to
// attack-action cards, with the OncePerTurn gate set — fires at most once per turn
// regardless of how many attack actions resolve.
func (ge *GameEngine) CreateOncePerTurnAttackActionAura(pc *card.CardState, handler func(card.GameEngine, card.Logger, card.Aura), count int) {
	ge.CreateAura(aura.NewFromCard(pc.Card, triggertype.CardOrAbility, handler, count, true, card.TypeSet.IsAttackAction))
}

// CreateHitOrDamageTakenAura registers an aura that fires when an attack hits
// (triggertype.Hit) or when the defense phase ends with damage unblocked
// (triggertype.DamageTaken) — the "destroyed when you deal or are dealt damage" aura
// shape. filter narrows the Hit side to a card-type predicate (nil = any hit); it never
// gates DamageTaken, which has no triggering card.
func (ge *GameEngine) CreateHitOrDamageTakenAura(pc *card.CardState, handler func(card.GameEngine, card.Logger, card.Aura), count int, filter func(card.TypeSet) bool) {
	ge.CreateAura(aura.NewFromCard(pc.Card, triggertype.Hit|triggertype.DamageTaken, handler, count, false, filter))
}
