package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// CreateAura implements card.GameEngine.CreateAura.
func (ge *GameEngine) CreateAura(pc *card.CardState, tt triggertype.Type, handler func(card.GameEngine, card.Logger, card.Aura), count int, oncePerTurn bool, filter func(card.TypeSet) bool) {
	ge.AppendAura(aura.NewFromCard(pc.Card, tt, handler, count, oncePerTurn, filter))
}
