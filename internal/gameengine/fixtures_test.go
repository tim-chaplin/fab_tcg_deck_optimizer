package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// Shared in-package test stubs. Inline rather than reaching for internal/testutils — using
// testutils from internal/gameengine tests would import internal/card transitively, which is fine here
// but the project prefers no testutils dependency at all so internal/card / internal/token / internal/trigger
// tests (which would cycle through testutils) stay consistent with internal/gameengine.

// stubCard is a minimal card.Card with configurable id and attack — exercises the fields
// each test reads, zero values for the rest. The zero-value id is ids.InvalidCard.
type stubCard struct {
	name   string
	attack int
	id     ids.CardID
}

func (c stubCard) ID() ids.CardID                                   { return c.id }
func (c stubCard) Name() string                                     { return c.name }
func (c stubCard) DisplayName() string                              { return c.name }
func (stubCard) Cost(card.GameEngine) int                           { return 0 }
func (stubCard) Pitch() int                                         { return 0 }
func (c stubCard) Attack() int                                      { return c.attack }
func (stubCard) Defense() int                                       { return 0 }
func (stubCard) Types(card.GameEngine) card.TypeSet                 { return 0 }
func (stubCard) GoAgain(card.GameEngine) bool                       { return false }
func (stubCard) Play(card.GameEngine, card.Logger, *card.CardState) {}

// dominatingStub is a stubCard plus the card.Dominator marker — exercises the printed-
// Dominate branch of LikelyToHit.
type dominatingStub struct {
	stubCard
}

func (dominatingStub) Dominate() {}
