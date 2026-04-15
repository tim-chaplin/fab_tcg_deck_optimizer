// Rune Flash — Runeblade Action - Attack. Defense 3. Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "Rune Flash costs {r} less to play for each Runechant you control."
// Simplification: assume enough Runechants to fully discount, so the effective cost is 0.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var runeFlashTypes = map[string]bool{"Runeblade": true, "Action": true, "Attack": true}

type RuneFlashRed struct{}

func (RuneFlashRed) Name() string               { return "Rune Flash (Red)" }
func (RuneFlashRed) Cost() int                  { return 0 }
func (RuneFlashRed) Pitch() int                 { return 1 }
func (RuneFlashRed) Attack() int                { return 4 }
func (RuneFlashRed) Defense() int               { return 3 }
func (RuneFlashRed) Types() map[string]bool     { return runeFlashTypes }
func (RuneFlashRed) GoAgain() bool              { return true }
func (c RuneFlashRed) Play(*card.TurnState) int { return c.Attack() }

type RuneFlashYellow struct{}

func (RuneFlashYellow) Name() string               { return "Rune Flash (Yellow)" }
func (RuneFlashYellow) Cost() int                  { return 0 }
func (RuneFlashYellow) Pitch() int                 { return 2 }
func (RuneFlashYellow) Attack() int                { return 3 }
func (RuneFlashYellow) Defense() int               { return 3 }
func (RuneFlashYellow) Types() map[string]bool     { return runeFlashTypes }
func (RuneFlashYellow) GoAgain() bool              { return true }
func (c RuneFlashYellow) Play(*card.TurnState) int { return c.Attack() }

type RuneFlashBlue struct{}

func (RuneFlashBlue) Name() string               { return "Rune Flash (Blue)" }
func (RuneFlashBlue) Cost() int                  { return 0 }
func (RuneFlashBlue) Pitch() int                 { return 3 }
func (RuneFlashBlue) Attack() int                { return 2 }
func (RuneFlashBlue) Defense() int               { return 3 }
func (RuneFlashBlue) Types() map[string]bool     { return runeFlashTypes }
func (RuneFlashBlue) GoAgain() bool              { return true }
func (c RuneFlashBlue) Play(*card.TurnState) int { return c.Attack() }
