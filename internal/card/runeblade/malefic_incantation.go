// Malefic Incantation — Runeblade Action - Aura. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "This enters the arena with N verse counters. When it has none, destroy it. Once per turn,
// when you play an attack action card, remove a verse counter from this. If you do, create a
// Runechant token." (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: if any attack (card or weapon) follows Malefic this turn, one counter ticks
// off and a Runechant fires — the aura is treated as destroyed this turn and goes to the
// graveyard. Otherwise the aura lingers and is destroyed at the start of the next turn via
// PlayNextTurn. This under-models Red/Yellow (they can legitimately stay in play across multiple
// turns before all counters are removed), but it matches the solver's coarse granularity and
// ensures Weeping Battleground can't retroactively banish the aura in the same turn's defense
// phase (DRs resolve before attacks, so CardsRemaining hasn't been consulted yet). The remaining
// N-1 counters are credited up-front as flat damage so the card's full damage potential is
// captured.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var maleficTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

// maleficPlay fires a Runechant when a future attack will pop a counter this turn, and moves
// the aura to the graveyard in that case. Remaining n-1 counters are credited as flat damage.
// Without a follow-up attack, the aura lingers for PlayNextTurn to destroy.
func maleficPlay(s *card.TurnState, n int) int {
	s.AuraCreated = true
	if hasFutureAttack(s) {
		if s.Self != nil {
			s.AddToGraveyard(s.Self.Card)
		}
		return s.CreateRunechant() + (n - 1)
	}
	return n - 1
}

// maleficNextTurn destroys a lingering Malefic at the start of the next turn — the solver's
// simplification treats any aura that survived the play-turn as destroyed here so it doesn't
// stick around forever (the real card keeps ticking one counter per turn).
func maleficNextTurn(s *card.TurnState, self card.Card) card.DelayedPlayResult {
	s.AddToGraveyard(self)
	return card.DelayedPlayResult{}
}

// hasFutureAttack reports whether any attack card or weapon swing remains in CardsRemaining
// after the current play position.
func hasFutureAttack(s *card.TurnState) bool {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if t.Has(card.TypeAttack) || t.Has(card.TypeWeapon) {
			return true
		}
	}
	return false
}

type MaleficIncantationRed struct{}

func (MaleficIncantationRed) ID() card.ID                 { return card.MaleficIncantationRed }
func (MaleficIncantationRed) Name() string                { return "Malefic Incantation (Red)" }
func (MaleficIncantationRed) Cost(*card.TurnState) int    { return 0 }
func (MaleficIncantationRed) Pitch() int                  { return 1 }
func (MaleficIncantationRed) Attack() int                 { return 0 }
func (MaleficIncantationRed) Defense() int                { return 2 }
func (MaleficIncantationRed) Types() card.TypeSet         { return maleficTypes }
func (MaleficIncantationRed) GoAgain() bool               { return true }
func (MaleficIncantationRed) Play(s *card.TurnState) int  { return maleficPlay(s, 3) }
func (c MaleficIncantationRed) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	return maleficNextTurn(s, c)
}

type MaleficIncantationYellow struct{}

func (MaleficIncantationYellow) ID() card.ID                 { return card.MaleficIncantationYellow }
func (MaleficIncantationYellow) Name() string                { return "Malefic Incantation (Yellow)" }
func (MaleficIncantationYellow) Cost(*card.TurnState) int    { return 0 }
func (MaleficIncantationYellow) Pitch() int                  { return 2 }
func (MaleficIncantationYellow) Attack() int                 { return 0 }
func (MaleficIncantationYellow) Defense() int                { return 2 }
func (MaleficIncantationYellow) Types() card.TypeSet         { return maleficTypes }
func (MaleficIncantationYellow) GoAgain() bool               { return true }
func (MaleficIncantationYellow) Play(s *card.TurnState) int  { return maleficPlay(s, 2) }
func (c MaleficIncantationYellow) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	return maleficNextTurn(s, c)
}

type MaleficIncantationBlue struct{}

func (MaleficIncantationBlue) ID() card.ID                 { return card.MaleficIncantationBlue }
func (MaleficIncantationBlue) Name() string                { return "Malefic Incantation (Blue)" }
func (MaleficIncantationBlue) Cost(*card.TurnState) int    { return 0 }
func (MaleficIncantationBlue) Pitch() int                  { return 3 }
func (MaleficIncantationBlue) Attack() int                 { return 0 }
func (MaleficIncantationBlue) Defense() int                { return 2 }
func (MaleficIncantationBlue) Types() card.TypeSet         { return maleficTypes }
func (MaleficIncantationBlue) GoAgain() bool               { return true }
func (MaleficIncantationBlue) Play(s *card.TurnState) int  { return maleficPlay(s, 1) }
func (c MaleficIncantationBlue) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	return maleficNextTurn(s, c)
}
