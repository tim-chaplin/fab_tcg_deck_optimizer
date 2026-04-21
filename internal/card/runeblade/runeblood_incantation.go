// Runeblood Incantation — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. Runeblood Incantation enters the arena with N verse counters on it. At the
// beginning of your action phase, remove a verse counter. If you do, create a Runechant token.
// Otherwise, destroy Runeblood Incantation." (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: credit n-1 flat damage on Play for the later-turn rune ticks, and model the
// first tick via card.DelayedPlay — PlayNextTurn creates 1 live Runechant on next turn's
// starting state and destroys the aura. Avoids over-crediting same-turn state so
// variable-cost cards can't use any of Runeblood's runes for a discount on the turn it was
// played.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var runebloodIncantationTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type RunebloodIncantationRed struct{}

func (RunebloodIncantationRed) ID() card.ID              { return card.RunebloodIncantationRed }
func (RunebloodIncantationRed) Name() string             { return "Runeblood Incantation (Red)" }
func (RunebloodIncantationRed) Cost(*card.TurnState) int { return 1 }
func (RunebloodIncantationRed) Pitch() int               { return 1 }
func (RunebloodIncantationRed) Attack() int              { return 0 }
func (RunebloodIncantationRed) Defense() int             { return 2 }
func (RunebloodIncantationRed) Types() card.TypeSet      { return runebloodIncantationTypes }
func (RunebloodIncantationRed) GoAgain() bool            { return true }
func (c RunebloodIncantationRed) Play(s *card.TurnState, _ *card.CardState) int {
	return runebloodPlay(s, 3)
}
func (c RunebloodIncantationRed) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	return runebloodPlayNextTurn(s, c)
}

type RunebloodIncantationYellow struct{}

func (RunebloodIncantationYellow) ID() card.ID              { return card.RunebloodIncantationYellow }
func (RunebloodIncantationYellow) Name() string             { return "Runeblood Incantation (Yellow)" }
func (RunebloodIncantationYellow) Cost(*card.TurnState) int { return 1 }
func (RunebloodIncantationYellow) Pitch() int               { return 2 }
func (RunebloodIncantationYellow) Attack() int              { return 0 }
func (RunebloodIncantationYellow) Defense() int             { return 2 }
func (RunebloodIncantationYellow) Types() card.TypeSet      { return runebloodIncantationTypes }
func (RunebloodIncantationYellow) GoAgain() bool            { return true }
func (c RunebloodIncantationYellow) Play(s *card.TurnState, _ *card.CardState) int {
	return runebloodPlay(s, 2)
}
func (c RunebloodIncantationYellow) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	return runebloodPlayNextTurn(s, c)
}

type RunebloodIncantationBlue struct{}

func (RunebloodIncantationBlue) ID() card.ID              { return card.RunebloodIncantationBlue }
func (RunebloodIncantationBlue) Name() string             { return "Runeblood Incantation (Blue)" }
func (RunebloodIncantationBlue) Cost(*card.TurnState) int { return 1 }
func (RunebloodIncantationBlue) Pitch() int               { return 3 }
func (RunebloodIncantationBlue) Attack() int              { return 0 }
func (RunebloodIncantationBlue) Defense() int             { return 2 }
func (RunebloodIncantationBlue) Types() card.TypeSet      { return runebloodIncantationTypes }
func (RunebloodIncantationBlue) GoAgain() bool            { return true }
func (c RunebloodIncantationBlue) Play(s *card.TurnState, _ *card.CardState) int {
	return runebloodPlay(s, 1)
}
func (c RunebloodIncantationBlue) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	return runebloodPlayNextTurn(s, c)
}

// runebloodPlay flips AuraCreated for same-turn aura-readers and credits n-1 flat damage for
// the future-turn verse-counter ticks that aren't separately modelled. The first tick's rune
// is created in PlayNextTurn.
func runebloodPlay(s *card.TurnState, n int) int {
	s.AuraCreated = true
	return n - 1
}

// runebloodPlayNextTurn fires the first verse-counter tick at the start of the next turn:
// destroy the aura and create 1 live Runechant on the new turn's starting state.
func runebloodPlayNextTurn(s *card.TurnState, self card.Card) card.DelayedPlayResult {
	s.AddToGraveyard(self)
	return card.DelayedPlayResult{Damage: s.CreateRunechants(1)}
}
