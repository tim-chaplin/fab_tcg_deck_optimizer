// Out Muscle — Generic Action - Attack. Cost 3. Printed power: Red 6, Yellow 5, Blue 4. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "While Out Muscle isn't defended by a card with equal or greater {p}, it has **go again**."
//
// Simplification: Defended-by-equal-or-greater-power gate isn't modelled; the printed Go again
// keyword is kept.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var outMuscleTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type OutMuscleRed struct{}

func (OutMuscleRed) ID() card.ID                 { return card.OutMuscleRed }
func (OutMuscleRed) Name() string                { return "Out Muscle (Red)" }
func (OutMuscleRed) Cost(*card.TurnState) int                   { return 3 }
func (OutMuscleRed) Pitch() int                  { return 1 }
func (OutMuscleRed) Attack() int                 { return 6 }
func (OutMuscleRed) Defense() int                { return 2 }
func (OutMuscleRed) Types() card.TypeSet         { return outMuscleTypes }
func (OutMuscleRed) GoAgain() bool               { return true }
func (c OutMuscleRed) Play(s *card.TurnState) int { return c.Attack() }

type OutMuscleYellow struct{}

func (OutMuscleYellow) ID() card.ID                 { return card.OutMuscleYellow }
func (OutMuscleYellow) Name() string                { return "Out Muscle (Yellow)" }
func (OutMuscleYellow) Cost(*card.TurnState) int                   { return 3 }
func (OutMuscleYellow) Pitch() int                  { return 2 }
func (OutMuscleYellow) Attack() int                 { return 5 }
func (OutMuscleYellow) Defense() int                { return 2 }
func (OutMuscleYellow) Types() card.TypeSet         { return outMuscleTypes }
func (OutMuscleYellow) GoAgain() bool               { return true }
func (c OutMuscleYellow) Play(s *card.TurnState) int { return c.Attack() }

type OutMuscleBlue struct{}

func (OutMuscleBlue) ID() card.ID                 { return card.OutMuscleBlue }
func (OutMuscleBlue) Name() string                { return "Out Muscle (Blue)" }
func (OutMuscleBlue) Cost(*card.TurnState) int                   { return 3 }
func (OutMuscleBlue) Pitch() int                  { return 3 }
func (OutMuscleBlue) Attack() int                 { return 4 }
func (OutMuscleBlue) Defense() int                { return 2 }
func (OutMuscleBlue) Types() card.TypeSet         { return outMuscleTypes }
func (OutMuscleBlue) GoAgain() bool               { return true }
func (c OutMuscleBlue) Play(s *card.TurnState) int { return c.Attack() }
