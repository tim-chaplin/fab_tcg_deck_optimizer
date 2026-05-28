package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/item"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// CreateItem registers a card-sourced triggered item.
func (gs *GameState) CreateItem(source card.Card, tt triggertype.Type, handler func(card.GameEngine, card.Logger, card.Item, card.FireContext), oncePerTurn bool, filter func(card.TypeSet) bool) {
	gs.items = append(gs.items, item.NewFromCard(source, tt, handler, oncePerTurn, filter))
}

// CreateItemWithAbility registers a card-sourced in-play permanent carrying an activated
// ability. See card.GameEngine for the wmask same-turn-delay contract.
func (gs *GameState) CreateItemWithAbility(source card.Card, ability card.Card) {
	gs.items = append(gs.items, item.NewFromCardWithAbility(source, ability))
}
