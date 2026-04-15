// Runeblood Incantation — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. Runeblood Incantation enters the arena with N verse counters on it. At the
// beginning of your action phase, remove a verse counter. If you do, create a Runechant token.
// Otherwise, destroy Runeblood Incantation." (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: assume every verse counter eventually produces a Runechant. Play returns N
// (Red=3, Yellow=2, Blue=1).
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var runebloodIncantationTypes = map[string]bool{"Runeblade": true, "Action": true, "Aura": true}

type RunebloodIncantationRed struct{}

func (RunebloodIncantationRed) Name() string              { return "Runeblood Incantation (Red)" }
func (RunebloodIncantationRed) Cost() int                 { return 1 }
func (RunebloodIncantationRed) Pitch() int                { return 1 }
func (RunebloodIncantationRed) Attack() int               { return 0 }
func (RunebloodIncantationRed) Defense() int              { return 2 }
func (RunebloodIncantationRed) Types() map[string]bool    { return runebloodIncantationTypes }
func (RunebloodIncantationRed) GoAgain() bool             { return true }
func (RunebloodIncantationRed) Play(*card.TurnState) int  { return 3 }

type RunebloodIncantationYellow struct{}

func (RunebloodIncantationYellow) Name() string             { return "Runeblood Incantation (Yellow)" }
func (RunebloodIncantationYellow) Cost() int                { return 1 }
func (RunebloodIncantationYellow) Pitch() int               { return 2 }
func (RunebloodIncantationYellow) Attack() int              { return 0 }
func (RunebloodIncantationYellow) Defense() int             { return 2 }
func (RunebloodIncantationYellow) Types() map[string]bool   { return runebloodIncantationTypes }
func (RunebloodIncantationYellow) GoAgain() bool            { return true }
func (RunebloodIncantationYellow) Play(*card.TurnState) int { return 2 }

type RunebloodIncantationBlue struct{}

func (RunebloodIncantationBlue) Name() string             { return "Runeblood Incantation (Blue)" }
func (RunebloodIncantationBlue) Cost() int                { return 1 }
func (RunebloodIncantationBlue) Pitch() int               { return 3 }
func (RunebloodIncantationBlue) Attack() int              { return 0 }
func (RunebloodIncantationBlue) Defense() int             { return 2 }
func (RunebloodIncantationBlue) Types() map[string]bool   { return runebloodIncantationTypes }
func (RunebloodIncantationBlue) GoAgain() bool            { return true }
func (RunebloodIncantationBlue) Play(*card.TurnState) int { return 1 }
