// High Striker — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 2.
//
// Text: "The next time an attack you control hits this turn, create 6 Copper tokens. **Go again**"
//
// Simplification: Copper-token economy isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var highStrikerTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type HighStrikerRed struct{}

func (HighStrikerRed) ID() card.ID                 { return card.HighStrikerRed }
func (HighStrikerRed) Name() string                { return "High Striker (Red)" }
func (HighStrikerRed) Cost(*card.TurnState) int                   { return 0 }
func (HighStrikerRed) Pitch() int                  { return 1 }
func (HighStrikerRed) Attack() int                 { return 0 }
func (HighStrikerRed) Defense() int                { return 2 }
func (HighStrikerRed) Types() card.TypeSet         { return highStrikerTypes }
func (HighStrikerRed) GoAgain() bool               { return true }
func (HighStrikerRed) Play(s *card.TurnState, _ *card.CardState) int { return 0 }

type HighStrikerYellow struct{}

func (HighStrikerYellow) ID() card.ID                 { return card.HighStrikerYellow }
func (HighStrikerYellow) Name() string                { return "High Striker (Yellow)" }
func (HighStrikerYellow) Cost(*card.TurnState) int                   { return 0 }
func (HighStrikerYellow) Pitch() int                  { return 2 }
func (HighStrikerYellow) Attack() int                 { return 0 }
func (HighStrikerYellow) Defense() int                { return 2 }
func (HighStrikerYellow) Types() card.TypeSet         { return highStrikerTypes }
func (HighStrikerYellow) GoAgain() bool               { return true }
func (HighStrikerYellow) Play(s *card.TurnState, _ *card.CardState) int { return 0 }

type HighStrikerBlue struct{}

func (HighStrikerBlue) ID() card.ID                 { return card.HighStrikerBlue }
func (HighStrikerBlue) Name() string                { return "High Striker (Blue)" }
func (HighStrikerBlue) Cost(*card.TurnState) int                   { return 0 }
func (HighStrikerBlue) Pitch() int                  { return 3 }
func (HighStrikerBlue) Attack() int                 { return 0 }
func (HighStrikerBlue) Defense() int                { return 2 }
func (HighStrikerBlue) Types() card.TypeSet         { return highStrikerTypes }
func (HighStrikerBlue) GoAgain() bool               { return true }
func (HighStrikerBlue) Play(s *card.TurnState, _ *card.CardState) int { return 0 }
