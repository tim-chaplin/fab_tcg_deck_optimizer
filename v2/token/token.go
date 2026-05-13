// Package token owns the concrete Item type (in-play permanents with activated abilities)
// and FaB's three item-flavored tokens: Gold, Silver, Copper. Each token is an Item paired
// with an activated-ability card the chain runner enqueues as a 1-AP playable while the
// token is in play.
//
// FaB's aura-flavored tokens (Runechant, Ponder) live in v2/aura because they share the
// concrete Aura type with card-backed auras.
//
// Satisfies gameengine.Item structurally without importing gameengine; Copy returns any,
// and the ability Play methods type-assert the engine to a local consumer interface.
package token

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Item is the concrete entry the engine stores in its persistent item list. Today every
// item is a token; the constructors below pair the count with the activated-ability card
// the chain runner enqueues each turn.
type Item struct {
	tokenName string
	tokenID   ids.CardID
	ability   card.Card
	count     int
}

// NewItem builds a token item with the supplied name, identifier, activated-ability card,
// and initial count. Production callers reach for NewGold / NewSilver / NewCopper.
func NewItem(name string, tokenID ids.CardID, ability card.Card, count int) *Item {
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

// Copy returns a deep copy boxed as any so the gameengine.Item interface declaration can
// avoid referencing its own type — letting concrete impls satisfy without importing.
func (i *Item) Copy() any {
	out := *i
	return &out
}

// Token display names — the engine matches by CardName when bumping an existing entry's
// Count or reading a count.
const (
	NameGold   = "Gold"
	NameSilver = "Silver"
	NameCopper = "Copper"
)

// NewGold / NewSilver / NewCopper return fresh token items at count n.
func NewGold(n int) *Item   { return NewItem(NameGold, ids.GoldTokenID, GoldAbility{}, n) }
func NewSilver(n int) *Item { return NewItem(NameSilver, ids.SilverTokenID, SilverAbility{}, n) }
func NewCopper(n int) *Item { return NewItem(NameCopper, ids.CopperTokenID, CopperAbility{}, n) }

// abilityConsumer is the slice of engine surface token activated abilities need.
// *gameengine.GameEngine satisfies it structurally.
type abilityConsumer interface {
	ConsumeItemByName(name string, n int)
}

var abilityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeItem)

// GoldAbility: cost {2}, draw a card, destroy one Gold token. Carries Go again.
type GoldAbility struct{}

func (GoldAbility) ID() ids.CardID                     { return ids.GoldTokenAbilityID }
func (GoldAbility) Name() string                       { return NameGold }
func (GoldAbility) DisplayName() string                { return NameGold }
func (GoldAbility) Cost(card.GameEngine) int           { return 2 }
func (GoldAbility) Pitch() int                         { return 0 }
func (GoldAbility) Attack() int                        { return 0 }
func (GoldAbility) Defense() int                       { return 0 }
func (GoldAbility) Types(card.GameEngine) card.TypeSet { return abilityTypes }
func (GoldAbility) GoAgain(card.GameEngine) bool       { return true }

func (GoldAbility) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return ge.GoldCount() > 0
}

func (GoldAbility) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.(abilityConsumer).ConsumeItemByName(NameGold, 1)
	ge.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 gold to draw a card", 0)
}

// SilverAbility: cost {3}, draw a card, destroy one Silver token.
type SilverAbility struct{}

func (SilverAbility) ID() ids.CardID                     { return ids.SilverTokenAbilityID }
func (SilverAbility) Name() string                       { return NameSilver }
func (SilverAbility) DisplayName() string                { return NameSilver }
func (SilverAbility) Cost(card.GameEngine) int           { return 3 }
func (SilverAbility) Pitch() int                         { return 0 }
func (SilverAbility) Attack() int                        { return 0 }
func (SilverAbility) Defense() int                       { return 0 }
func (SilverAbility) Types(card.GameEngine) card.TypeSet { return abilityTypes }
func (SilverAbility) GoAgain(card.GameEngine) bool       { return true }

func (SilverAbility) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return ge.SilverCount() > 0
}

func (SilverAbility) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.(abilityConsumer).ConsumeItemByName(NameSilver, 1)
	ge.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 silver to draw a card", 0)
}

// CopperAbility: cost {4}, draw a card, destroy one Copper token.
type CopperAbility struct{}

func (CopperAbility) ID() ids.CardID                     { return ids.CopperTokenAbilityID }
func (CopperAbility) Name() string                       { return NameCopper }
func (CopperAbility) DisplayName() string                { return NameCopper }
func (CopperAbility) Cost(card.GameEngine) int           { return 4 }
func (CopperAbility) Pitch() int                         { return 0 }
func (CopperAbility) Attack() int                        { return 0 }
func (CopperAbility) Defense() int                       { return 0 }
func (CopperAbility) Types(card.GameEngine) card.TypeSet { return abilityTypes }
func (CopperAbility) GoAgain(card.GameEngine) bool       { return true }

func (CopperAbility) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return ge.CopperCount() > 0
}

func (CopperAbility) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.(abilityConsumer).ConsumeItemByName(NameCopper, 1)
	ge.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Spent 1 copper to draw a card", 0)
}
