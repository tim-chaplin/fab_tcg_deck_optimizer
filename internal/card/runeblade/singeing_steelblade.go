// Singeing Steelblade — Runeblade Action - Attack. Cost 1, Defense 3, Arcane 1.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When you attack with Singeing Steelblade, deal 1 arcane damage to target hero."
//
// Simplification: arcane damage counts as regular damage. Play returns power + 1.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var singeingSteelbladeTypes = map[string]bool{"Runeblade": true, "Action": true, "Attack": true}

type SingeingSteelbladeRed struct{}

func (SingeingSteelbladeRed) Name() string               { return "Singeing Steelblade (Red)" }
func (SingeingSteelbladeRed) Cost() int                  { return 1 }
func (SingeingSteelbladeRed) Pitch() int                 { return 1 }
func (SingeingSteelbladeRed) Attack() int                { return 4 }
func (SingeingSteelbladeRed) Defense() int               { return 3 }
func (SingeingSteelbladeRed) Types() map[string]bool     { return singeingSteelbladeTypes }
func (SingeingSteelbladeRed) GoAgain() bool              { return false }
func (c SingeingSteelbladeRed) Play(*card.TurnState) int { return c.Attack() + 1 }

type SingeingSteelbladeYellow struct{}

func (SingeingSteelbladeYellow) Name() string               { return "Singeing Steelblade (Yellow)" }
func (SingeingSteelbladeYellow) Cost() int                  { return 1 }
func (SingeingSteelbladeYellow) Pitch() int                 { return 2 }
func (SingeingSteelbladeYellow) Attack() int                { return 3 }
func (SingeingSteelbladeYellow) Defense() int               { return 3 }
func (SingeingSteelbladeYellow) Types() map[string]bool     { return singeingSteelbladeTypes }
func (SingeingSteelbladeYellow) GoAgain() bool              { return false }
func (c SingeingSteelbladeYellow) Play(*card.TurnState) int { return c.Attack() + 1 }

type SingeingSteelbladeBlue struct{}

func (SingeingSteelbladeBlue) Name() string               { return "Singeing Steelblade (Blue)" }
func (SingeingSteelbladeBlue) Cost() int                  { return 1 }
func (SingeingSteelbladeBlue) Pitch() int                 { return 3 }
func (SingeingSteelbladeBlue) Attack() int                { return 2 }
func (SingeingSteelbladeBlue) Defense() int               { return 3 }
func (SingeingSteelbladeBlue) Types() map[string]bool     { return singeingSteelbladeTypes }
func (SingeingSteelbladeBlue) GoAgain() bool              { return false }
func (c SingeingSteelbladeBlue) Play(*card.TurnState) int { return c.Attack() + 1 }
