// Reek of Corruption — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you have played or created an aura this turn, Reek of Corruption gains 'When this
// hits a hero, they discard a card.'"
//
// Simplification: assume both the aura condition and the on-hit trigger are always satisfied.
// Discard is valued at +3 (matching Consuming Volition's discard valuation).
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var reekOfCorruptionTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type ReekOfCorruptionRed struct{}

func (ReekOfCorruptionRed) Name() string               { return "Reek of Corruption (Red)" }
func (ReekOfCorruptionRed) Cost() int                  { return 2 }
func (ReekOfCorruptionRed) Pitch() int                 { return 1 }
func (ReekOfCorruptionRed) Attack() int                { return 4 }
func (ReekOfCorruptionRed) Defense() int               { return 3 }
func (ReekOfCorruptionRed) Types() card.TypeSet        { return reekOfCorruptionTypes }
func (ReekOfCorruptionRed) GoAgain() bool              { return false }
func (c ReekOfCorruptionRed) Play(*card.TurnState) int { return c.Attack() + 3 }

type ReekOfCorruptionYellow struct{}

func (ReekOfCorruptionYellow) Name() string               { return "Reek of Corruption (Yellow)" }
func (ReekOfCorruptionYellow) Cost() int                  { return 2 }
func (ReekOfCorruptionYellow) Pitch() int                 { return 2 }
func (ReekOfCorruptionYellow) Attack() int                { return 3 }
func (ReekOfCorruptionYellow) Defense() int               { return 3 }
func (ReekOfCorruptionYellow) Types() card.TypeSet        { return reekOfCorruptionTypes }
func (ReekOfCorruptionYellow) GoAgain() bool              { return false }
func (c ReekOfCorruptionYellow) Play(*card.TurnState) int { return c.Attack() + 3 }

type ReekOfCorruptionBlue struct{}

func (ReekOfCorruptionBlue) Name() string               { return "Reek of Corruption (Blue)" }
func (ReekOfCorruptionBlue) Cost() int                  { return 2 }
func (ReekOfCorruptionBlue) Pitch() int                 { return 3 }
func (ReekOfCorruptionBlue) Attack() int                { return 2 }
func (ReekOfCorruptionBlue) Defense() int               { return 3 }
func (ReekOfCorruptionBlue) Types() card.TypeSet        { return reekOfCorruptionTypes }
func (ReekOfCorruptionBlue) GoAgain() bool              { return false }
func (c ReekOfCorruptionBlue) Play(*card.TurnState) int { return c.Attack() + 3 }
