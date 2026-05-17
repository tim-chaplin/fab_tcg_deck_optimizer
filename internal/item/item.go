// Package item owns the concrete Item type — an in-play permanent with an activated
// ability. The chain runner stores Items on the engine alongside auras; the ability slot
// is typed as `any` so item carries no internal/card dependency. Callers type-assert the ability
// to their richer card type when invoking.
//
// Constructors are typed by source: NewFromToken wraps the built-in FaB tokens (Gold /
// Silver / Copper, sourced from internal/cards via internal/token). A NewFromCard sibling will
// be added when card-source items land (not yet implemented).
package item

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// Item is the concrete entry the engine stores in its persistent item list. Today every
// item is a token; the ability is the activated card the chain runner enqueues as a
// 1-AP playable while the item is in play.
type Item struct {
	tokenName string
	tokenID   ids.CardID
	ability   any
	count     int
}

// NewFromToken builds a token-sourced item with the supplied name, identifier,
// activated-ability card, and initial count. Production callers reach for internal/token's
// NewGold / NewSilver / NewCopper rather than calling this directly.
func NewFromToken(name string, tokenID ids.CardID, ability any, count int) *Item {
	return &Item{
		tokenName: name,
		tokenID:   tokenID,
		ability:   ability,
		count:     count,
	}
}

func (i *Item) CardName() string   { return i.tokenName }
func (i *Item) CardID() ids.CardID { return i.tokenID }
func (i *Item) Count() int         { return i.count }
func (i *Item) SetCount(n int)     { i.count = n }

// Ability returns the activated-ability card boxed as any. Callers type-assert to their
// richer card type (typically internal/card.Card) when invoking.
func (i *Item) Ability() any { return i.ability }

// Copy returns a deep copy boxed as any so the gameengine.Item interface declaration can
// avoid referencing its own type — letting concrete impls satisfy without importing.
func (i *Item) Copy() any {
	out := *i
	return &out
}
