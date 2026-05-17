// Copper token: cost {4}, draw a card, destroy one Copper token. Carries Go again.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// copperTokenName is the canonical display name. The engine matches by CardName when
// bumping an existing entry's Count or reading a count. External callers reach it via
// CopperToken{}.Name() rather than importing the const directly.
const copperTokenName = "Copper"

// copperTypes is the Types bitmask CopperToken reports. Cached as a package-level var so
// Types() returns without recomputing the OR per call.
var copperTypes = card.NewTypeSet(card.TypeGeneric, card.TypeItem)

// copperConsumer is the narrow engine surface CopperToken.Play needs to destroy the
// consumed token entry. *gameengine.GameEngine satisfies it structurally.
type copperConsumer interface {
	ConsumeItemByName(name string, n int)
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
func (CopperToken) Types(card.GameEngine) card.TypeSet { return copperTypes }
func (CopperToken) GoAgain(card.GameEngine) bool       { return true }

func (CopperToken) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return ge.CopperCount() > 0
}

func (CopperToken) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.(copperConsumer).ConsumeItemByName(copperTokenName, 1)
	ge.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 copper to draw a card", 0)
}
