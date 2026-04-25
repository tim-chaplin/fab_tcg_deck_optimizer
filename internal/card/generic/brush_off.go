// Brush Off — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "The next time you would be dealt 3 or less damage this turn, prevent it."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var brushOffTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type BrushOffRed struct{}

func (BrushOffRed) ID() card.ID                               { return card.BrushOffRed }
func (BrushOffRed) Name() string                              { return "Brush Off (Red)" }
func (BrushOffRed) Cost(*card.TurnState) int                  { return 0 }
func (BrushOffRed) Pitch() int                                { return 1 }
func (BrushOffRed) Attack() int                               { return 0 }
func (BrushOffRed) Defense() int                              { return 0 }
func (BrushOffRed) Types() card.TypeSet                       { return brushOffTypes }
func (BrushOffRed) GoAgain() bool                             { return false }
// not implemented: Instant prevent 1 damage; gated on aura/item with counter
func (BrushOffRed) NotImplemented()                           {}
func (BrushOffRed) Play(*card.TurnState, *card.CardState) int { return 0 }

type BrushOffYellow struct{}

func (BrushOffYellow) ID() card.ID                               { return card.BrushOffYellow }
func (BrushOffYellow) Name() string                              { return "Brush Off (Yellow)" }
func (BrushOffYellow) Cost(*card.TurnState) int                  { return 0 }
func (BrushOffYellow) Pitch() int                                { return 2 }
func (BrushOffYellow) Attack() int                               { return 0 }
func (BrushOffYellow) Defense() int                              { return 0 }
func (BrushOffYellow) Types() card.TypeSet                       { return brushOffTypes }
func (BrushOffYellow) GoAgain() bool                             { return false }
// not implemented: Instant prevent 1 damage; gated on aura/item with counter
func (BrushOffYellow) NotImplemented()                           {}
func (BrushOffYellow) Play(*card.TurnState, *card.CardState) int { return 0 }

type BrushOffBlue struct{}

func (BrushOffBlue) ID() card.ID                               { return card.BrushOffBlue }
func (BrushOffBlue) Name() string                              { return "Brush Off (Blue)" }
func (BrushOffBlue) Cost(*card.TurnState) int                  { return 0 }
func (BrushOffBlue) Pitch() int                                { return 3 }
func (BrushOffBlue) Attack() int                               { return 0 }
func (BrushOffBlue) Defense() int                              { return 0 }
func (BrushOffBlue) Types() card.TypeSet                       { return brushOffTypes }
func (BrushOffBlue) GoAgain() bool                             { return false }
// not implemented: Instant prevent 1 damage; gated on aura/item with counter
func (BrushOffBlue) NotImplemented()                           {}
func (BrushOffBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
