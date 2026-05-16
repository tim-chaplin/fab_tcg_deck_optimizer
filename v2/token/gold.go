// Gold token: cost {2}, draw a card, destroy one Gold token. Carries Go again.

package token

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// goldTokenName is the canonical display name. The engine matches by CardName when
// bumping an existing entry's Count or reading a count.
const goldTokenName = "Gold"

// NewGold returns a fresh Gold token item at count n.
func NewGold(n int) *Item {
	return NewItem(goldTokenName, ids.GoldTokenID, GoldToken{}, n)
}

// goldTypes is the Types bitmask GoldToken reports. Cached as a package-level var so
// Types() returns without recomputing the OR per call.
var goldTypes = card.NewTypeSet(card.TypeGeneric, card.TypeItem)

// goldConsumer is the narrow engine surface GoldToken.Play needs to destroy the consumed
// token entry. *gameengine.GameEngine satisfies it structurally.
type goldConsumer interface {
	ConsumeItemByName(name string, n int)
}

// GoldToken is the activated-ability card: cost {2}, draw a card, destroy one Gold token.
type GoldToken struct{}

func (GoldToken) ID() ids.CardID                     { return ids.GoldTokenAbilityID }
func (GoldToken) Name() string                       { return goldTokenName }
func (GoldToken) DisplayName() string                { return goldTokenName }
func (GoldToken) Cost(card.GameEngine) int           { return 2 }
func (GoldToken) Pitch() int                         { return 0 }
func (GoldToken) Attack() int                        { return 0 }
func (GoldToken) Defense() int                       { return 0 }
func (GoldToken) Types(card.GameEngine) card.TypeSet { return goldTypes }
func (GoldToken) GoAgain(card.GameEngine) bool       { return true }

func (GoldToken) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return ge.GoldCount() > 0
}

func (GoldToken) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.(goldConsumer).ConsumeItemByName(goldTokenName, 1)
	ge.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 gold to draw a card", 0)
}
