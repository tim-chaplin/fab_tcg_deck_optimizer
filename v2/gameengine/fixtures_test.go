package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Shared in-package test stubs. Inline rather than reaching for internal/testutils — using
// testutils from v2/gameengine tests would import v2/card transitively, which is fine here
// but the project prefers no testutils dependency at all so v2/card / v2/token / v2/trigger
// tests (which would cycle through testutils) stay consistent with v2/gameengine.

// stubCard is a minimal card.Card with configurable attack — exercises the fields each
// test reads, zero values for the rest.
type stubCard struct {
	name   string
	attack int
}

func (c stubCard) ID() ids.CardID                                   { return ids.InvalidCard }
func (c stubCard) Name() string                                     { return c.name }
func (c stubCard) DisplayName() string                              { return c.name }
func (stubCard) Cost(card.GameEngine) int                           { return 0 }
func (stubCard) Pitch() int                                         { return 0 }
func (c stubCard) Attack() int                                      { return c.attack }
func (stubCard) Defense() int                                       { return 0 }
func (stubCard) Types(card.GameEngine) card.TypeSet                 { return 0 }
func (stubCard) GoAgain(card.GameEngine) bool                       { return false }
func (stubCard) Play(card.GameEngine, card.Logger, *card.CardState) {}

// dominatingStub is a card.Card with configurable attack AND the card.Dominator marker —
// exercises the printed-Dominate branch of LikelyToHit.
type dominatingStub struct {
	name   string
	attack int
}

func (c dominatingStub) ID() ids.CardID                                   { return ids.InvalidCard }
func (c dominatingStub) Name() string                                     { return c.name }
func (c dominatingStub) DisplayName() string                              { return c.name }
func (dominatingStub) Cost(card.GameEngine) int                           { return 0 }
func (dominatingStub) Pitch() int                                         { return 0 }
func (c dominatingStub) Attack() int                                      { return c.attack }
func (dominatingStub) Defense() int                                       { return 0 }
func (dominatingStub) Types(card.GameEngine) card.TypeSet                 { return 0 }
func (dominatingStub) GoAgain(card.GameEngine) bool                       { return false }
func (dominatingStub) Play(card.GameEngine, card.Logger, *card.CardState) {}
func (dominatingStub) Dominate()                                          {}
