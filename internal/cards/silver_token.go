// Silver token: cost {3}, draw a card, destroy one Silver token. Carries Go again.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/token"
)

const silverTokenName = "Silver"

// NewSilver returns a fresh Silver token item at count n.
func NewSilver(n int) *token.Item {
	return token.NewItem(silverTokenName, ids.SilverTokenID, SilverToken{}, n)
}

// SilverToken is the activated-ability card: cost {3}, draw a card, destroy one Silver.
type SilverToken struct{}

func (SilverToken) ID() ids.CardID                     { return ids.SilverTokenAbilityID }
func (SilverToken) Name() string                       { return silverTokenName }
func (SilverToken) DisplayName() string                { return silverTokenName }
func (SilverToken) Cost(card.GameEngine) int           { return 3 }
func (SilverToken) Pitch() int                         { return 0 }
func (SilverToken) Attack() int                        { return 0 }
func (SilverToken) Defense() int                       { return 0 }
func (SilverToken) Types(card.GameEngine) card.TypeSet { return tokenAbilityTypes }
func (SilverToken) GoAgain(card.GameEngine) bool       { return true }

func (SilverToken) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return ge.SilverCount() > 0
}

func (SilverToken) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.(tokenAbilityConsumer).ConsumeItemByName(silverTokenName, 1)
	ge.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 silver to draw a card", 0)
}
