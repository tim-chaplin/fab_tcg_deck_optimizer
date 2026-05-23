package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// CreateAura is the card-facing aura registration method on *GameEngine. Cards call it
// via the internal/card.GameEngine interface; the method constructs the concrete
// *aura.Aura and appends it via AppendAura. Source is taken from pc.Card; pass
// triggertype bits OR'd together for multi-event auras (e.g.
// triggertype.Hit|triggertype.DamageTaken). filter narrows the firing site to a
// card-type predicate (nil = any); it is ignored on turn-boundary and DamageTaken
// events that have no triggering card.
//
// Compile-time bridge: the inline handler type assigns cleanly to aura.Handler — drift
// between the two declarations would break this call.
func (ge *GameEngine) CreateAura(pc *card.CardState, tt triggertype.Type, handler func(card.GameEngine, card.Logger, card.Aura), count int, oncePerTurn bool, filter func(card.TypeSet) bool) {
	ge.AppendAura(aura.NewFromCard(pc.Card, tt, handler, count, oncePerTurn, filter))
}
