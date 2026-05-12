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
func NewTokenItem(name string, tokenID ids.CardID, ability card.Card, count int) gameengine.Item {
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
func (i *Item) Ability() card.Card { return i.ability }

func (i *Item) Clone() gameengine.Item {
	out := *i
	return &out
}
