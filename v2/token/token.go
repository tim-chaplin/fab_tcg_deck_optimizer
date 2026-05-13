// Package token owns the concrete Item type — an in-play permanent with an activated
// ability (currently FaB's three item-flavored tokens: Gold, Silver, Copper). Item itself
// is generic; concrete activated-ability cards (GoldToken / SilverToken / CopperToken)
// and their NewGold / NewSilver / NewCopper factories live in internal/cards, where the
// rest of FaB's card implementations sit.
//
// The package defines its own narrow interfaces and does NOT import v2/card — the stored
// ability is typed as `any` so consumers (the chain runner) assert it back to a richer
// card type when invoking.
package token

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// Item is the concrete entry the engine stores in its persistent item list. Today every
// item is a token; the ability is the activated card the chain runner enqueues as a
// 1-AP playable while the token is in play.
type Item struct {
	tokenName string
	tokenID   ids.CardID
	ability   any
	count     int
}

// NewItem builds a token item with the supplied name, identifier, activated-ability card,
// and initial count. Production callers reach for cards.NewGold / cards.NewSilver /
// cards.NewCopper which wire the appropriate ability card.
func NewItem(name string, tokenID ids.CardID, ability any, count int) *Item {
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
// richer card type (typically v2/card.Card) when invoking.
func (i *Item) Ability() any { return i.ability }

// Copy returns a deep copy boxed as any so the gameengine.Item interface declaration can
// avoid referencing its own type — letting concrete impls satisfy without importing.
func (i *Item) Copy() any {
	out := *i
	return &out
}
