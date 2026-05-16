// Silver token: cost {3}, draw a card, destroy one Silver token. Carries Go again.

package token

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

const silverTokenName = "Silver"

// NewSilver returns a fresh Silver token item at count n.
func NewSilver(n int) *Item {
	return NewItem(silverTokenName, ids.SilverTokenID, SilverToken{}, n)
}

// silverTypes is the Types bitmask SilverToken reports. Cached as a package-level var so
// Types() returns without recomputing the OR per call.
var silverTypes = card.NewTypeSet(card.TypeGeneric, card.TypeItem)

// silverConsumer is the narrow engine surface SilverToken.Play needs to destroy the
// consumed token entry. *gameengine.GameEngine satisfies it structurally.
type silverConsumer interface {
	ConsumeItemByName(name string, n int)
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
func (SilverToken) Types(card.GameEngine) card.TypeSet { return silverTypes }
func (SilverToken) GoAgain(card.GameEngine) bool       { return true }

func (SilverToken) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return ge.SilverCount() > 0
}

func (SilverToken) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.(silverConsumer).ConsumeItemByName(silverTokenName, 1)
	ge.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 silver to draw a card", 0)
}
