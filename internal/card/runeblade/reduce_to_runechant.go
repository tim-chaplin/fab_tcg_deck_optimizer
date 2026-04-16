// Reduce to Runechant — Runeblade Defense Reaction. Cost 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "Reduce to Runechant costs {r} less to play for each Runechant you control. Create a
// Runechant token."
// Simplification: assume we always have (at least) 1 Runechant, so the cost is effectively 0.
// The created Runechant token isn't tracked across turns.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var reduceToRunechantTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeDefenseReaction)

type ReduceToRunechantRed struct{}

func (ReduceToRunechantRed) ID() card.ID                 { return card.ReduceToRunechantRed }
func (ReduceToRunechantRed) Name() string             { return "Reduce to Runechant (Red)" }
func (ReduceToRunechantRed) Cost() int                { return 0 }
func (ReduceToRunechantRed) Pitch() int               { return 1 }
func (ReduceToRunechantRed) Attack() int              { return 0 }
func (ReduceToRunechantRed) Defense() int             { return 4 }
func (ReduceToRunechantRed) Types() card.TypeSet      { return reduceToRunechantTypes }
func (ReduceToRunechantRed) GoAgain() bool            { return false }
func (ReduceToRunechantRed) Play(*card.TurnState) int { return 0 }

type ReduceToRunechantYellow struct{}

func (ReduceToRunechantYellow) ID() card.ID                 { return card.ReduceToRunechantYellow }
func (ReduceToRunechantYellow) Name() string             { return "Reduce to Runechant (Yellow)" }
func (ReduceToRunechantYellow) Cost() int                { return 0 }
func (ReduceToRunechantYellow) Pitch() int               { return 2 }
func (ReduceToRunechantYellow) Attack() int              { return 0 }
func (ReduceToRunechantYellow) Defense() int             { return 3 }
func (ReduceToRunechantYellow) Types() card.TypeSet      { return reduceToRunechantTypes }
func (ReduceToRunechantYellow) GoAgain() bool            { return false }
func (ReduceToRunechantYellow) Play(*card.TurnState) int { return 0 }

type ReduceToRunechantBlue struct{}

func (ReduceToRunechantBlue) ID() card.ID                 { return card.ReduceToRunechantBlue }
func (ReduceToRunechantBlue) Name() string             { return "Reduce to Runechant (Blue)" }
func (ReduceToRunechantBlue) Cost() int                { return 0 }
func (ReduceToRunechantBlue) Pitch() int               { return 3 }
func (ReduceToRunechantBlue) Attack() int              { return 0 }
func (ReduceToRunechantBlue) Defense() int             { return 2 }
func (ReduceToRunechantBlue) Types() card.TypeSet      { return reduceToRunechantTypes }
func (ReduceToRunechantBlue) GoAgain() bool            { return false }
func (ReduceToRunechantBlue) Play(*card.TurnState) int { return 0 }
