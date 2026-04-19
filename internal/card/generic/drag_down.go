// Drag Down — Generic Defense Reaction. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 0.
//
// Text: "When this defends an attack, it gets -3{p}."
//
// Simplification: The -3{p} attacker debuff isn't modelled (solver doesn't expose defender-side
// power reductions).
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

type DragDownRed struct{}

func (DragDownRed) ID() card.ID                 { return card.DragDownRed }
func (DragDownRed) Name() string                { return "Drag Down (Red)" }
func (DragDownRed) Cost(*card.TurnState) int                   { return 0 }
func (DragDownRed) Pitch() int                  { return 1 }
func (DragDownRed) Attack() int                 { return 0 }
func (DragDownRed) Defense() int                { return 0 }
func (DragDownRed) Types() card.TypeSet         { return defenseReactionTypes }
func (DragDownRed) GoAgain() bool               { return false }
func (DragDownRed) Play(s *card.TurnState) int { return 0 }

type DragDownYellow struct{}

func (DragDownYellow) ID() card.ID                 { return card.DragDownYellow }
func (DragDownYellow) Name() string                { return "Drag Down (Yellow)" }
func (DragDownYellow) Cost(*card.TurnState) int                   { return 0 }
func (DragDownYellow) Pitch() int                  { return 2 }
func (DragDownYellow) Attack() int                 { return 0 }
func (DragDownYellow) Defense() int                { return 0 }
func (DragDownYellow) Types() card.TypeSet         { return defenseReactionTypes }
func (DragDownYellow) GoAgain() bool               { return false }
func (DragDownYellow) Play(s *card.TurnState) int { return 0 }

type DragDownBlue struct{}

func (DragDownBlue) ID() card.ID                 { return card.DragDownBlue }
func (DragDownBlue) Name() string                { return "Drag Down (Blue)" }
func (DragDownBlue) Cost(*card.TurnState) int                   { return 0 }
func (DragDownBlue) Pitch() int                  { return 3 }
func (DragDownBlue) Attack() int                 { return 0 }
func (DragDownBlue) Defense() int                { return 0 }
func (DragDownBlue) Types() card.TypeSet         { return defenseReactionTypes }
func (DragDownBlue) GoAgain() bool               { return false }
func (DragDownBlue) Play(s *card.TurnState) int { return 0 }
