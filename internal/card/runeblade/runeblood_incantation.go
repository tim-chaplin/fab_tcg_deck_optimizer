// Runeblood Incantation — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. Runeblood Incantation enters the arena with N verse counters on it. At the
// beginning of your action phase, remove a verse counter. If you do, create a Runechant token.
// Otherwise, destroy Runeblood Incantation." (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: the first counter ticks via PlayNextTurn (creating one Runechant at the start
// of the turn after Runeblood is played); the remaining N-1 are credited as flat future-turn
// damage at play time. The aura then destroys itself on that same PlayNextTurn — a coarse
// under-model of Red/Yellow, which would normally linger several turns.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var runebloodIncantationTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

// runebloodPlay credits N-1 of the future Runechant ticks as flat damage; the first ticks on
// PlayNextTurn. AuraCreated is flipped so same-turn aura-readers see it.
func runebloodPlay(s *card.TurnState, n int) int {
	s.AuraCreated = true
	return n - 1
}

// runebloodNextTurn destroys self at the start of the next turn and creates one Runechant
// token as its leave-arena payoff.
func runebloodNextTurn(s *card.TurnState, self card.Card) card.DelayedPlayResult {
	s.AddToGraveyard(self)
	return card.DelayedPlayResult{Damage: s.CreateRunechant()}
}

type RunebloodIncantationRed struct{}

func (RunebloodIncantationRed) ID() card.ID                { return card.RunebloodIncantationRed }
func (RunebloodIncantationRed) Name() string               { return "Runeblood Incantation (Red)" }
func (RunebloodIncantationRed) Cost(*card.TurnState) int   { return 1 }
func (RunebloodIncantationRed) Pitch() int                 { return 1 }
func (RunebloodIncantationRed) Attack() int                { return 0 }
func (RunebloodIncantationRed) Defense() int               { return 2 }
func (RunebloodIncantationRed) Types() card.TypeSet        { return runebloodIncantationTypes }
func (RunebloodIncantationRed) GoAgain() bool              { return true }
func (RunebloodIncantationRed) Play(s *card.TurnState) int { return runebloodPlay(s, 3) }
func (c RunebloodIncantationRed) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	return runebloodNextTurn(s, c)
}

type RunebloodIncantationYellow struct{}

func (RunebloodIncantationYellow) ID() card.ID                { return card.RunebloodIncantationYellow }
func (RunebloodIncantationYellow) Name() string               { return "Runeblood Incantation (Yellow)" }
func (RunebloodIncantationYellow) Cost(*card.TurnState) int   { return 1 }
func (RunebloodIncantationYellow) Pitch() int                 { return 2 }
func (RunebloodIncantationYellow) Attack() int                { return 0 }
func (RunebloodIncantationYellow) Defense() int               { return 2 }
func (RunebloodIncantationYellow) Types() card.TypeSet        { return runebloodIncantationTypes }
func (RunebloodIncantationYellow) GoAgain() bool              { return true }
func (RunebloodIncantationYellow) Play(s *card.TurnState) int { return runebloodPlay(s, 2) }
func (c RunebloodIncantationYellow) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	return runebloodNextTurn(s, c)
}

type RunebloodIncantationBlue struct{}

func (RunebloodIncantationBlue) ID() card.ID                { return card.RunebloodIncantationBlue }
func (RunebloodIncantationBlue) Name() string               { return "Runeblood Incantation (Blue)" }
func (RunebloodIncantationBlue) Cost(*card.TurnState) int   { return 1 }
func (RunebloodIncantationBlue) Pitch() int                 { return 3 }
func (RunebloodIncantationBlue) Attack() int                { return 0 }
func (RunebloodIncantationBlue) Defense() int               { return 2 }
func (RunebloodIncantationBlue) Types() card.TypeSet        { return runebloodIncantationTypes }
func (RunebloodIncantationBlue) GoAgain() bool              { return true }
func (RunebloodIncantationBlue) Play(s *card.TurnState) int { return runebloodPlay(s, 1) }
func (c RunebloodIncantationBlue) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	return runebloodNextTurn(s, c)
}
