package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// stubCard is a minimal Card implementation for exercising TurnState helpers. Tests only care
// about identity and the Types / Attack / GoAgain hooks some helpers probe; everything else
// returns a zero value. id defaults to Invalid (0); tests that reach into ID-keyed caches
// should set it to a distinct value to avoid colliding on slot 0.
type stubCard struct {
	id      ids.CardID
	name    string
	types   card.TypeSet
	attack  int
	goAgain bool
}

func (c stubCard) ID() ids.CardID            { return c.id }
func (c stubCard) Name() string              { return c.name }
func (stubCard) Cost(*TurnState) int         { return 0 }
func (stubCard) Pitch() int                  { return 0 }
func (c stubCard) Attack() int               { return c.attack }
func (stubCard) Defense() int                { return 0 }
func (c stubCard) Types() card.TypeSet       { return c.types }
func (c stubCard) GoAgain() bool             { return c.goAgain }
func (stubCard) Play(*TurnState, *CardState) {}

// dominatingStubCard is a stubCard that implements the Dominator marker — exercises the
// printed-Dominate branch of EffectiveDominate / HasDominate.
type dominatingStubCard struct {
	stubCard
}

func (dominatingStubCard) Dominate() {}

// notImplementedStubCard is a stubCard that implements the NotImplemented marker — exercises
// the type assertion the deck legal-pool filter keys on.
type notImplementedStubCard struct {
	stubCard
}

func (notImplementedStubCard) NotImplemented() {}

// unplayableStubCard is a stubCard that implements the Unplayable marker — exercises the
// second pool-exclusion path the deck legal-pool filter keys on.
type unplayableStubCard struct {
	stubCard
}

func (unplayableStubCard) Unplayable() {}
