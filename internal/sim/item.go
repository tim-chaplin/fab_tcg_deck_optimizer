package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Item is the sim concrete impl of gameengine.Item. Sim's token-item factories
// (NewGoldItem, NewSilverItem, NewCopperItem) build one via NewTokenItem and seed the
// engine's item list.
type Item struct {
	tokenName string // "Gold" / "Silver" / "Copper" — only token items today
	tokenID   ids.CardID
	ability   card.Card
	count     int
}

// NewTokenItem builds a token item with the supplied name, identifier, activated-ability
// card, and initial count.
func NewTokenItem(name string, tokenID ids.CardID, ability card.Card, count int) *Item {
	return &Item{
		tokenName: name,
		tokenID:   tokenID,
		ability:   ability,
		count:     count,
	}
}

// itemSliceAsEngine converts []*Item to []gameengine.Item for engine-API call sites.
func itemSliceAsEngine(src []*Item) []gameengine.Item {
	if len(src) == 0 {
		return nil
	}
	out := make([]gameengine.Item, len(src))
	for i, it := range src {
		out[i] = it
	}
	return out
}

// itemSliceFromEngine type-asserts engine-returned []gameengine.Item back to []*Item.
func itemSliceFromEngine(src []gameengine.Item) []*Item {
	if len(src) == 0 {
		return nil
	}
	out := make([]*Item, len(src))
	for i, it := range src {
		out[i] = it.(*Item)
	}
	return out
}

func (i *Item) CardName() string   { return i.tokenName }
func (i *Item) CardID() ids.CardID { return i.tokenID }
func (i *Item) Count() int         { return i.count }
func (i *Item) SetCount(n int)     { i.count = n }
func (i *Item) Ability() card.Card { return i.ability }

func (i *Item) Clone() gameengine.Item {
	out := *i
	return &out
}
