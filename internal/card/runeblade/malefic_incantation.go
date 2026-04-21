// Malefic Incantation — Runeblade Action - Aura. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "This enters the arena with N verse counters. When it has none, destroy it. Once per turn,
// when you play an attack action card, remove a verse counter from this. If you do, create a
// Runechant token." (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: credit n-1 flat damage on Play for the later-turn verse-counter ticks, and
// model the first tick via card.DelayedPlay — PlayNextTurn creates 1 live Runechant on next
// turn's starting state and destroys the aura. The "once per turn, when you play an attack"
// condition is approximated by always firing once at next turn's upkeep.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var maleficTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type MaleficIncantationRed struct{}

func (MaleficIncantationRed) ID() card.ID              { return card.MaleficIncantationRed }
func (MaleficIncantationRed) Name() string             { return "Malefic Incantation (Red)" }
func (MaleficIncantationRed) Cost(*card.TurnState) int { return 0 }
func (MaleficIncantationRed) Pitch() int               { return 1 }
func (MaleficIncantationRed) Attack() int              { return 0 }
func (MaleficIncantationRed) Defense() int             { return 2 }
func (MaleficIncantationRed) Types() card.TypeSet      { return maleficTypes }
func (MaleficIncantationRed) GoAgain() bool            { return true }
func (c MaleficIncantationRed) Play(s *card.TurnState, _ *card.CardState) int {
	return maleficPlay(s, 3)
}
func (c MaleficIncantationRed) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	return maleficPlayNextTurn(s, c)
}

type MaleficIncantationYellow struct{}

func (MaleficIncantationYellow) ID() card.ID              { return card.MaleficIncantationYellow }
func (MaleficIncantationYellow) Name() string             { return "Malefic Incantation (Yellow)" }
func (MaleficIncantationYellow) Cost(*card.TurnState) int { return 0 }
func (MaleficIncantationYellow) Pitch() int               { return 2 }
func (MaleficIncantationYellow) Attack() int              { return 0 }
func (MaleficIncantationYellow) Defense() int             { return 2 }
func (MaleficIncantationYellow) Types() card.TypeSet      { return maleficTypes }
func (MaleficIncantationYellow) GoAgain() bool            { return true }
func (c MaleficIncantationYellow) Play(s *card.TurnState, _ *card.CardState) int {
	return maleficPlay(s, 2)
}
func (c MaleficIncantationYellow) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	return maleficPlayNextTurn(s, c)
}

type MaleficIncantationBlue struct{}

func (MaleficIncantationBlue) ID() card.ID              { return card.MaleficIncantationBlue }
func (MaleficIncantationBlue) Name() string             { return "Malefic Incantation (Blue)" }
func (MaleficIncantationBlue) Cost(*card.TurnState) int { return 0 }
func (MaleficIncantationBlue) Pitch() int               { return 3 }
func (MaleficIncantationBlue) Attack() int              { return 0 }
func (MaleficIncantationBlue) Defense() int             { return 2 }
func (MaleficIncantationBlue) Types() card.TypeSet      { return maleficTypes }
func (MaleficIncantationBlue) GoAgain() bool            { return true }
func (c MaleficIncantationBlue) Play(s *card.TurnState, _ *card.CardState) int {
	return maleficPlay(s, 1)
}
func (c MaleficIncantationBlue) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	return maleficPlayNextTurn(s, c)
}

// maleficPlay flips AuraCreated for same-turn aura-readers and credits n-1 flat damage for
// the future-turn verse-counter ticks that aren't separately modelled. The first tick's rune
// is created in PlayNextTurn.
func maleficPlay(s *card.TurnState, n int) int {
	s.AuraCreated = true
	return n - 1
}

// maleficPlayNextTurn fires the first verse-counter tick at the start of the next turn:
// destroy the aura and create 1 live Runechant on the new turn's starting state.
func maleficPlayNextTurn(s *card.TurnState, self card.Card) card.DelayedPlayResult {
	s.AddToGraveyard(self)
	return card.DelayedPlayResult{Damage: s.CreateRunechants(1)}
}
