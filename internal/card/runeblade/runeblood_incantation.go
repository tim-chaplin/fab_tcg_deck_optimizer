// Runeblood Incantation — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. Runeblood Incantation enters the arena with N verse counters on it. At the
// beginning of your action phase, remove a verse counter. If you do, create a Runechant token.
// Otherwise, destroy Runeblood Incantation." (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: 1 Runechant is delayed into next turn's carryover (the counter that'll tick
// at next turn's action phase start); the remaining N-1 are credited as flat future-turn damage
// without tracking the tokens. Avoids over-crediting same-turn state so variable-cost
// cards can't use all N runes for a full discount in the turn Runeblood was played.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var runebloodIncantationTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type RunebloodIncantationRed struct{}

func (RunebloodIncantationRed) ID() card.ID                 { return card.RunebloodIncantationRed }
func (RunebloodIncantationRed) Name() string              { return "Runeblood Incantation (Red)" }
func (RunebloodIncantationRed) Cost(*card.TurnState) int                 { return 1 }
func (RunebloodIncantationRed) Pitch() int                { return 1 }
func (RunebloodIncantationRed) Attack() int               { return 0 }
func (RunebloodIncantationRed) Defense() int              { return 2 }
func (RunebloodIncantationRed) Types() card.TypeSet       { return runebloodIncantationTypes }
func (RunebloodIncantationRed) GoAgain() bool             { return true }
func (RunebloodIncantationRed) Play(s *card.TurnState, _ *card.CardState) int  { return runebloodPlay(s, 3) }

type RunebloodIncantationYellow struct{}

func (RunebloodIncantationYellow) ID() card.ID                 { return card.RunebloodIncantationYellow }
func (RunebloodIncantationYellow) Name() string             { return "Runeblood Incantation (Yellow)" }
func (RunebloodIncantationYellow) Cost(*card.TurnState) int                { return 1 }
func (RunebloodIncantationYellow) Pitch() int               { return 2 }
func (RunebloodIncantationYellow) Attack() int              { return 0 }
func (RunebloodIncantationYellow) Defense() int             { return 2 }
func (RunebloodIncantationYellow) Types() card.TypeSet      { return runebloodIncantationTypes }
func (RunebloodIncantationYellow) GoAgain() bool            { return true }
func (RunebloodIncantationYellow) Play(s *card.TurnState, _ *card.CardState) int { return runebloodPlay(s, 2) }

type RunebloodIncantationBlue struct{}

func (RunebloodIncantationBlue) ID() card.ID                 { return card.RunebloodIncantationBlue }
func (RunebloodIncantationBlue) Name() string             { return "Runeblood Incantation (Blue)" }
func (RunebloodIncantationBlue) Cost(*card.TurnState) int                { return 1 }
func (RunebloodIncantationBlue) Pitch() int               { return 3 }
func (RunebloodIncantationBlue) Attack() int              { return 0 }
func (RunebloodIncantationBlue) Defense() int             { return 2 }
func (RunebloodIncantationBlue) Types() card.TypeSet      { return runebloodIncantationTypes }
func (RunebloodIncantationBlue) GoAgain() bool            { return true }
func (RunebloodIncantationBlue) Play(s *card.TurnState, _ *card.CardState) int { return runebloodPlay(s, 1) }

// runebloodPlay delays 1 Runechant to next turn (the counter that'll tick first) and credits
// the remaining n-1 as untracked flat damage. At n=1 it's just the single delayed token.
func runebloodPlay(s *card.TurnState, n int) int {
	return s.DelayRunechants(1) + (n - 1)
}
