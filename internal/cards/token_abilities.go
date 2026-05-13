// FaB's three item-flavored tokens: Gold, Silver, Copper. Each token is a token.Item
// paired with an activated-ability card the chain runner enqueues as a 1-AP playable while
// the token is in play.
//
// The Item type lives in v2/token (which is card-free); the abilities and the
// NewGold / NewSilver / NewCopper factories live here so the v2/token package doesn't
// transitively import v2/card.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/token"
)

// Token display names — the engine matches by CardName when bumping an existing entry's
// Count or reading a count.
const (
	tokenNameGold   = "Gold"
	tokenNameSilver = "Silver"
	tokenNameCopper = "Copper"
)

// NewGold / NewSilver / NewCopper return fresh token items at count n.
func NewGold(n int) *token.Item {
	return token.NewItem(tokenNameGold, ids.GoldTokenID, GoldAbility{}, n)
}
func NewSilver(n int) *token.Item {
	return token.NewItem(tokenNameSilver, ids.SilverTokenID, SilverAbility{}, n)
}
func NewCopper(n int) *token.Item {
	return token.NewItem(tokenNameCopper, ids.CopperTokenID, CopperAbility{}, n)
}

// abilityConsumer is the slice of engine surface token activated abilities need.
// *gameengine.GameEngine satisfies it structurally.
type abilityConsumer interface {
	ConsumeItemByName(name string, n int)
}

var tokenAbilityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeItem)

// GoldAbility: cost {2}, draw a card, destroy one Gold token. Carries Go again.
type GoldAbility struct{}

func (GoldAbility) ID() ids.CardID                     { return ids.GoldTokenAbilityID }
func (GoldAbility) Name() string                       { return tokenNameGold }
func (GoldAbility) DisplayName() string                { return tokenNameGold }
func (GoldAbility) Cost(card.GameEngine) int           { return 2 }
func (GoldAbility) Pitch() int                         { return 0 }
func (GoldAbility) Attack() int                        { return 0 }
func (GoldAbility) Defense() int                       { return 0 }
func (GoldAbility) Types(card.GameEngine) card.TypeSet { return tokenAbilityTypes }
func (GoldAbility) GoAgain(card.GameEngine) bool       { return true }

func (GoldAbility) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return ge.GoldCount() > 0
}

func (GoldAbility) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.(abilityConsumer).ConsumeItemByName(tokenNameGold, 1)
	ge.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 gold to draw a card", 0)
}

// SilverAbility: cost {3}, draw a card, destroy one Silver token.
type SilverAbility struct{}

func (SilverAbility) ID() ids.CardID                     { return ids.SilverTokenAbilityID }
func (SilverAbility) Name() string                       { return tokenNameSilver }
func (SilverAbility) DisplayName() string                { return tokenNameSilver }
func (SilverAbility) Cost(card.GameEngine) int           { return 3 }
func (SilverAbility) Pitch() int                         { return 0 }
func (SilverAbility) Attack() int                        { return 0 }
func (SilverAbility) Defense() int                       { return 0 }
func (SilverAbility) Types(card.GameEngine) card.TypeSet { return tokenAbilityTypes }
func (SilverAbility) GoAgain(card.GameEngine) bool       { return true }

func (SilverAbility) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return ge.SilverCount() > 0
}

func (SilverAbility) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.(abilityConsumer).ConsumeItemByName(tokenNameSilver, 1)
	ge.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 silver to draw a card", 0)
}

// CopperAbility: cost {4}, draw a card, destroy one Copper token.
type CopperAbility struct{}

func (CopperAbility) ID() ids.CardID                     { return ids.CopperTokenAbilityID }
func (CopperAbility) Name() string                       { return tokenNameCopper }
func (CopperAbility) DisplayName() string                { return tokenNameCopper }
func (CopperAbility) Cost(card.GameEngine) int           { return 4 }
func (CopperAbility) Pitch() int                         { return 0 }
func (CopperAbility) Attack() int                        { return 0 }
func (CopperAbility) Defense() int                       { return 0 }
func (CopperAbility) Types(card.GameEngine) card.TypeSet { return tokenAbilityTypes }
func (CopperAbility) GoAgain(card.GameEngine) bool       { return true }

func (CopperAbility) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return ge.CopperCount() > 0
}

func (CopperAbility) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.(abilityConsumer).ConsumeItemByName(tokenNameCopper, 1)
	ge.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 copper to draw a card", 0)
}
