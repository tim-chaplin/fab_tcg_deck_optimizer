// Gold token: cost {2}, draw a card, destroy one Gold token. Carries Go again.
//
// The Item-side bookkeeping (count, name, ID) lives on token.Item; this file owns the
// activated-ability card the chain runner enqueues as a 1-AP playable while at least one
// Gold token is in play.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/token"
)

// goldTokenName is the canonical display name. The engine matches by CardName when
// bumping an existing entry's Count or reading a count.
const goldTokenName = "Gold"

// NewGold returns a fresh Gold token item at count n.
func NewGold(n int) *token.Item {
	return token.NewItem(goldTokenName, ids.GoldTokenID, GoldToken{}, n)
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
func (GoldToken) Types(card.GameEngine) card.TypeSet { return tokenAbilityTypes }
func (GoldToken) GoAgain(card.GameEngine) bool       { return true }

func (GoldToken) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return ge.GoldCount() > 0
}

func (GoldToken) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.(tokenAbilityConsumer).ConsumeItemByName(goldTokenName, 1)
	ge.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 gold to draw a card", 0)
}
