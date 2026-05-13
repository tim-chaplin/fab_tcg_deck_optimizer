// Copper token: cost {4}, draw a card, destroy one Copper token. Carries Go again.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/token"
)

const copperTokenName = "Copper"

// NewCopper returns a fresh Copper token item at count n.
func NewCopper(n int) *token.Item {
	return token.NewItem(copperTokenName, ids.CopperTokenID, CopperToken{}, n)
}

// CopperToken is the activated-ability card: cost {4}, draw a card, destroy one Copper.
type CopperToken struct{}

func (CopperToken) ID() ids.CardID                     { return ids.CopperTokenAbilityID }
func (CopperToken) Name() string                       { return copperTokenName }
func (CopperToken) DisplayName() string                { return copperTokenName }
func (CopperToken) Cost(card.GameEngine) int           { return 4 }
func (CopperToken) Pitch() int                         { return 0 }
func (CopperToken) Attack() int                        { return 0 }
func (CopperToken) Defense() int                       { return 0 }
func (CopperToken) Types(card.GameEngine) card.TypeSet { return tokenAbilityTypes }
func (CopperToken) GoAgain(card.GameEngine) bool       { return true }

func (CopperToken) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return ge.CopperCount() > 0
}

func (CopperToken) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.(tokenAbilityConsumer).ConsumeItemByName(copperTokenName, 1)
	ge.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 copper to draw a card", 0)
}
