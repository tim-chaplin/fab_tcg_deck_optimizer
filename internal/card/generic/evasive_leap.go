// Evasive Leap — Generic Defense Reaction. Cost 0.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 3, Yellow 2, Blue 1.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

type EvasiveLeapRed struct{}

func (EvasiveLeapRed) ID() card.ID                 { return card.EvasiveLeapRed }
func (EvasiveLeapRed) Name() string             { return "Evasive Leap (Red)" }
func (EvasiveLeapRed) Cost(*card.TurnState) int                { return 0 }
func (EvasiveLeapRed) Pitch() int               { return 1 }
func (EvasiveLeapRed) Attack() int              { return 0 }
func (EvasiveLeapRed) Defense() int             { return 3 }
func (EvasiveLeapRed) Types() card.TypeSet      { return defenseReactionTypes }
func (EvasiveLeapRed) GoAgain() bool            { return false }
func (EvasiveLeapRed) Play(*card.TurnState) int { return 0 }

type EvasiveLeapYellow struct{}

func (EvasiveLeapYellow) ID() card.ID                 { return card.EvasiveLeapYellow }
func (EvasiveLeapYellow) Name() string             { return "Evasive Leap (Yellow)" }
func (EvasiveLeapYellow) Cost(*card.TurnState) int                { return 0 }
func (EvasiveLeapYellow) Pitch() int               { return 2 }
func (EvasiveLeapYellow) Attack() int              { return 0 }
func (EvasiveLeapYellow) Defense() int             { return 2 }
func (EvasiveLeapYellow) Types() card.TypeSet      { return defenseReactionTypes }
func (EvasiveLeapYellow) GoAgain() bool            { return false }
func (EvasiveLeapYellow) Play(*card.TurnState) int { return 0 }

type EvasiveLeapBlue struct{}

func (EvasiveLeapBlue) ID() card.ID                 { return card.EvasiveLeapBlue }
func (EvasiveLeapBlue) Name() string             { return "Evasive Leap (Blue)" }
func (EvasiveLeapBlue) Cost(*card.TurnState) int                { return 0 }
func (EvasiveLeapBlue) Pitch() int               { return 3 }
func (EvasiveLeapBlue) Attack() int              { return 0 }
func (EvasiveLeapBlue) Defense() int             { return 1 }
func (EvasiveLeapBlue) Types() card.TypeSet      { return defenseReactionTypes }
func (EvasiveLeapBlue) GoAgain() bool            { return false }
func (EvasiveLeapBlue) Play(*card.TurnState) int { return 0 }
