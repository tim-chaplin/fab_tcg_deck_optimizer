package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// Shared in-package test stubs. Inline rather than reaching for internal/testutils — using
// testutils from internal/gameengine tests would import internal/card transitively, which is fine here
// but the project prefers no testutils dependency at all so internal/card / internal/token / internal/trigger
// tests (which would cycle through testutils) stay consistent with internal/gameengine.

// fakeCard is a minimal card.Card with configurable id and attack — exercises the fields
// each test reads, zero values for the rest. The zero-value id is ids.InvalidCard.
type fakeCard struct {
	name   string
	attack int
	id     ids.CardID
}

func (c fakeCard) ID() ids.CardID                                   { return c.id }
func (c fakeCard) Name() string                                     { return c.name }
func (c fakeCard) DisplayName() string                              { return c.name }
func (fakeCard) Cost() int                                          { return 0 }
func (fakeCard) Pitch() int                                         { return 0 }
func (c fakeCard) Attack() int                                      { return c.attack }
func (fakeCard) Defense() int                                       { return 0 }
func (fakeCard) Types(card.GameEngine) card.TypeSet                 { return 0 }
func (fakeCard) GoAgain(card.GameEngine) bool                       { return false }
func (fakeCard) Play(card.GameEngine, card.Logger, *card.CardState) {}

// dominatingFake is a fakeCard plus the card.Dominator marker — exercises the printed-
// Dominate branch of LikelyToHit.
type dominatingFake struct {
	fakeCard
}

func (dominatingFake) Dominate() {}
