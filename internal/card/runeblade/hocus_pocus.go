// Hocus Pocus — Runeblade Action - Attack. Cost 0, Defense 3.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "When this attacks, create a Runechant token."
//
// Simplification: each Runechant = +1 future damage. Play returns power + 1.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var hocusPocusTypes = map[string]bool{"Runeblade": true, "Action": true, "Attack": true}

type HocusPocusRed struct{}

func (HocusPocusRed) Name() string               { return "Hocus Pocus (Red)" }
func (HocusPocusRed) Cost() int                  { return 0 }
func (HocusPocusRed) Pitch() int                 { return 1 }
func (HocusPocusRed) Attack() int                { return 3 }
func (HocusPocusRed) Defense() int               { return 3 }
func (HocusPocusRed) Types() map[string]bool     { return hocusPocusTypes }
func (HocusPocusRed) GoAgain() bool              { return false }
func (c HocusPocusRed) Play(*card.TurnState) int { return c.Attack() + 1 }

type HocusPocusYellow struct{}

func (HocusPocusYellow) Name() string               { return "Hocus Pocus (Yellow)" }
func (HocusPocusYellow) Cost() int                  { return 0 }
func (HocusPocusYellow) Pitch() int                 { return 2 }
func (HocusPocusYellow) Attack() int                { return 2 }
func (HocusPocusYellow) Defense() int               { return 3 }
func (HocusPocusYellow) Types() map[string]bool     { return hocusPocusTypes }
func (HocusPocusYellow) GoAgain() bool              { return false }
func (c HocusPocusYellow) Play(*card.TurnState) int { return c.Attack() + 1 }

type HocusPocusBlue struct{}

func (HocusPocusBlue) Name() string               { return "Hocus Pocus (Blue)" }
func (HocusPocusBlue) Cost() int                  { return 0 }
func (HocusPocusBlue) Pitch() int                 { return 3 }
func (HocusPocusBlue) Attack() int                { return 1 }
func (HocusPocusBlue) Defense() int               { return 3 }
func (HocusPocusBlue) Types() map[string]bool     { return hocusPocusTypes }
func (HocusPocusBlue) GoAgain() bool              { return false }
func (c HocusPocusBlue) Play(*card.TurnState) int { return c.Attack() + 1 }
