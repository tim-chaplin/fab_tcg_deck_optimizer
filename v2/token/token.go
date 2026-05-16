// Package token owns the factories for FaB's five built-in tokens: the item tokens
// Gold / Silver / Copper, and the aura-flavored tokens Runechant / Ponder. Each factory
// returns the concrete value the engine stores — *item.Item for the item tokens, *aura.Aura
// for the aura tokens — wiring in the per-token name, identifier, and (for auras) the
// trigger type + fire handler.
//
// The activated-ability card types backing the item tokens (cards.GoldToken et al.) live
// in internal/cards alongside other card implementations. The aura-token fire handlers
// and their narrow consumer interfaces live in runechant.go / ponder.go alongside the
// token-specific consts.
package token

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/item"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// NewGold returns a fresh Gold token item at count n.
func NewGold(n int) *item.Item {
	return item.NewFromToken(cards.GoldTokenName, ids.GoldTokenID, cards.GoldToken{}, n)
}

// NewSilver returns a fresh Silver token item at count n.
func NewSilver(n int) *item.Item {
	return item.NewFromToken(cards.SilverTokenName, ids.SilverTokenID, cards.SilverToken{}, n)
}

// NewCopper returns a fresh Copper token item at count n.
func NewCopper(n int) *item.Item {
	return item.NewFromToken(cards.CopperTokenName, ids.CopperTokenID, cards.CopperToken{}, n)
}

// NewRunechant returns a fresh Runechant token aura at count n. Fires per attack
// (triggertype.Attack); the fire handler lives in runechant.go.
func NewRunechant(n int) *aura.Aura {
	return aura.NewFromToken(runechantTokenName, ids.RunechantTokenID, triggertype.Attack, runechantFire, n)
}

// NewPonder returns a fresh Ponder token aura at count n. Fires at end-of-turn
// (triggertype.EndOfTurn); the fire handler lives in ponder.go.
func NewPonder(n int) *aura.Aura {
	return aura.NewFromToken(ponderTokenName, ids.PonderTokenID, triggertype.EndOfTurn, ponderFire, n)
}
