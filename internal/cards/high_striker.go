// High Striker — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 2.
//
// Text: "The next time an attack you control hits this turn, create 6 Copper tokens. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var highStrikerTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type HighStrikerRed struct{}

func (HighStrikerRed) ID() ids.CardID          { return ids.HighStrikerRed }
func (HighStrikerRed) Name() string            { return "High Striker" }
func (HighStrikerRed) Cost(*sim.TurnState) int { return 0 }
func (HighStrikerRed) Pitch() int              { return 1 }
func (HighStrikerRed) Attack() int             { return 0 }
func (HighStrikerRed) Defense() int            { return 2 }
func (HighStrikerRed) Types() card.TypeSet     { return highStrikerTypes }
func (HighStrikerRed) GoAgain() bool           { return true }

// not implemented: copper tokens
func (HighStrikerRed) NotImplemented()                            {}
func (HighStrikerRed) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }

type HighStrikerYellow struct{}

func (HighStrikerYellow) ID() ids.CardID          { return ids.HighStrikerYellow }
func (HighStrikerYellow) Name() string            { return "High Striker" }
func (HighStrikerYellow) Cost(*sim.TurnState) int { return 0 }
func (HighStrikerYellow) Pitch() int              { return 2 }
func (HighStrikerYellow) Attack() int             { return 0 }
func (HighStrikerYellow) Defense() int            { return 2 }
func (HighStrikerYellow) Types() card.TypeSet     { return highStrikerTypes }
func (HighStrikerYellow) GoAgain() bool           { return true }

// not implemented: copper tokens
func (HighStrikerYellow) NotImplemented()                            {}
func (HighStrikerYellow) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }

type HighStrikerBlue struct{}

func (HighStrikerBlue) ID() ids.CardID          { return ids.HighStrikerBlue }
func (HighStrikerBlue) Name() string            { return "High Striker" }
func (HighStrikerBlue) Cost(*sim.TurnState) int { return 0 }
func (HighStrikerBlue) Pitch() int              { return 3 }
func (HighStrikerBlue) Attack() int             { return 0 }
func (HighStrikerBlue) Defense() int            { return 2 }
func (HighStrikerBlue) Types() card.TypeSet     { return highStrikerTypes }
func (HighStrikerBlue) GoAgain() bool           { return true }

// not implemented: copper tokens
func (HighStrikerBlue) NotImplemented()                            {}
func (HighStrikerBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
