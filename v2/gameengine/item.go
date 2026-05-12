package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// itemEntry is the engine-owned Item impl. Sim's token-item factories (NewGoldItem,
// NewSilverItem, NewCopperItem) build one via NewTokenItem and seed the engine's item
// list.
type itemEntry struct {
	tokenName string // "Gold" / "Silver" / "Copper" — only token items today
	tokenID   ids.CardID
	ability   card.Card
	count     int
}

// NewTokenItem builds a token item with the supplied name, identifier, activated-ability
// card, and initial count.
func NewTokenItem(name string, tokenID ids.CardID, ability card.Card, count int) Item {
	return &itemEntry{
		tokenName: name,
		tokenID:   tokenID,
		ability:   ability,
		count:     count,
	}
}

func (i *itemEntry) CardName() string   { return i.tokenName }
func (i *itemEntry) CardID() ids.CardID { return i.tokenID }
func (i *itemEntry) Count() int         { return i.count }
func (i *itemEntry) SetCount(n int)     { i.count = n }
func (i *itemEntry) Ability() card.Card { return i.ability }

func (i *itemEntry) Clone() Item {
	out := *i
	return &out
}
