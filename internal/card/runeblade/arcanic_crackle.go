// Arcanic Crackle — Runeblade Action - Attack. Cost 0, Defense 3, Arcane 1.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "Deal 1 arcane damage to target hero."
//
// Simplification: arcane damage counts as regular damage. Play returns power + 1.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var arcanicCrackleTypes = map[string]bool{"Runeblade": true, "Action": true, "Attack": true}

type ArcanicCrackleRed struct{}

func (ArcanicCrackleRed) Name() string                 { return "Arcanic Crackle (Red)" }
func (ArcanicCrackleRed) Cost() int                    { return 0 }
func (ArcanicCrackleRed) Pitch() int                   { return 1 }
func (ArcanicCrackleRed) Attack() int                  { return 3 }
func (ArcanicCrackleRed) Defense() int                 { return 3 }
func (ArcanicCrackleRed) Types() map[string]bool       { return arcanicCrackleTypes }
func (ArcanicCrackleRed) GoAgain() bool                { return false }
func (c ArcanicCrackleRed) Play(*card.TurnState) int   { return c.Attack() + 1 }

type ArcanicCrackleYellow struct{}

func (ArcanicCrackleYellow) Name() string               { return "Arcanic Crackle (Yellow)" }
func (ArcanicCrackleYellow) Cost() int                  { return 0 }
func (ArcanicCrackleYellow) Pitch() int                 { return 2 }
func (ArcanicCrackleYellow) Attack() int                { return 2 }
func (ArcanicCrackleYellow) Defense() int               { return 3 }
func (ArcanicCrackleYellow) Types() map[string]bool     { return arcanicCrackleTypes }
func (ArcanicCrackleYellow) GoAgain() bool              { return false }
func (c ArcanicCrackleYellow) Play(*card.TurnState) int { return c.Attack() + 1 }

type ArcanicCrackleBlue struct{}

func (ArcanicCrackleBlue) Name() string               { return "Arcanic Crackle (Blue)" }
func (ArcanicCrackleBlue) Cost() int                  { return 0 }
func (ArcanicCrackleBlue) Pitch() int                 { return 3 }
func (ArcanicCrackleBlue) Attack() int                { return 1 }
func (ArcanicCrackleBlue) Defense() int               { return 3 }
func (ArcanicCrackleBlue) Types() map[string]bool     { return arcanicCrackleTypes }
func (ArcanicCrackleBlue) GoAgain() bool              { return false }
func (c ArcanicCrackleBlue) Play(*card.TurnState) int { return c.Attack() + 1 }
