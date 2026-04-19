// Runeblood Incantation — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. Runeblood Incantation enters the arena with N verse counters on it. At the
// beginning of your action phase, remove a verse counter. If you do, create a Runechant token.
// Otherwise, destroy Runeblood Incantation." (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: the aura's total payout is N Runechants if it ticks down all N counters, so
// we credit N as future-turn damage. Any turn we take damage interrupts the tick-down (we die
// faster, fewer turns to milk it, and in practice the aura gets answered), so when the current
// partition doesn't block all incoming damage we collapse the value to 0 — same fragile-aura
// model as Arcane Cussing.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var runebloodIncantationTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type RunebloodIncantationRed struct{}

func (RunebloodIncantationRed) ID() card.ID                 { return card.RunebloodIncantationRed }
func (RunebloodIncantationRed) Name() string                { return "Runeblood Incantation (Red)" }
func (RunebloodIncantationRed) Cost() int                   { return 1 }
func (RunebloodIncantationRed) Pitch() int                  { return 1 }
func (RunebloodIncantationRed) Attack() int                 { return 0 }
func (RunebloodIncantationRed) Defense() int                { return 2 }
func (RunebloodIncantationRed) Types() card.TypeSet         { return runebloodIncantationTypes }
func (RunebloodIncantationRed) GoAgain() bool               { return true }
func (RunebloodIncantationRed) Play(s *card.TurnState) int  { return auraSurvivalValue(s, 3) }

type RunebloodIncantationYellow struct{}

func (RunebloodIncantationYellow) ID() card.ID                 { return card.RunebloodIncantationYellow }
func (RunebloodIncantationYellow) Name() string                { return "Runeblood Incantation (Yellow)" }
func (RunebloodIncantationYellow) Cost() int                   { return 1 }
func (RunebloodIncantationYellow) Pitch() int                  { return 2 }
func (RunebloodIncantationYellow) Attack() int                 { return 0 }
func (RunebloodIncantationYellow) Defense() int                { return 2 }
func (RunebloodIncantationYellow) Types() card.TypeSet         { return runebloodIncantationTypes }
func (RunebloodIncantationYellow) GoAgain() bool               { return true }
func (RunebloodIncantationYellow) Play(s *card.TurnState) int  { return auraSurvivalValue(s, 2) }

type RunebloodIncantationBlue struct{}

func (RunebloodIncantationBlue) ID() card.ID                 { return card.RunebloodIncantationBlue }
func (RunebloodIncantationBlue) Name() string                { return "Runeblood Incantation (Blue)" }
func (RunebloodIncantationBlue) Cost() int                   { return 1 }
func (RunebloodIncantationBlue) Pitch() int                  { return 3 }
func (RunebloodIncantationBlue) Attack() int                 { return 0 }
func (RunebloodIncantationBlue) Defense() int                { return 2 }
func (RunebloodIncantationBlue) Types() card.TypeSet         { return runebloodIncantationTypes }
func (RunebloodIncantationBlue) GoAgain() bool               { return true }
func (RunebloodIncantationBlue) Play(s *card.TurnState) int  { return auraSurvivalValue(s, 1) }
