// Malefic Incantation — Runeblade Action - Aura. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "This enters the arena with N verse counters. When it has none, destroy it. Once per turn,
// when you play an attack action card, remove a verse counter from this. If you do, create a
// Runechant token." (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: if at least one attack action card follows Malefic in this turn's chain, the
// "once per turn" trigger fires — create a live Runechant now. For Blue (N=1) the counter hits
// zero, so the aura is destroyed immediately and lands in this turn's graveyard. For Red/Yellow
// the remaining n-1 future-turn ticks are credited as flat damage and the aura heads to next
// turn's graveyard via card.DelayedPlay. Without a follow-up attack action the trigger can't
// fire; credit flat n for the full run of future ticks and still graveyard next turn.
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
func (MaleficIncantationRed) AddsFutureValue()         {}
func (c MaleficIncantationRed) Play(s *card.TurnState, _ *card.CardState) int {
	return maleficPlay(s, 3, c)
}
func (c MaleficIncantationRed) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	s.AddToGraveyard(c)
	return card.DelayedPlayResult{}
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
func (MaleficIncantationYellow) AddsFutureValue()         {}
func (c MaleficIncantationYellow) Play(s *card.TurnState, _ *card.CardState) int {
	return maleficPlay(s, 2, c)
}
func (c MaleficIncantationYellow) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	s.AddToGraveyard(c)
	return card.DelayedPlayResult{}
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
func (MaleficIncantationBlue) AddsFutureValue()         {}
func (c MaleficIncantationBlue) Play(s *card.TurnState, _ *card.CardState) int {
	return maleficPlay(s, 1, c)
}
func (c MaleficIncantationBlue) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	s.AddToGraveyard(c)
	return card.DelayedPlayResult{}
}

// maleficPlay flips AuraCreated and handles the same-turn tick. With a follow-up attack action
// in the chain, create a live Runechant now. At n=1 the verse counter hits zero so the aura
// lands in this turn's graveyard immediately; at n>1 the remaining n-1 future-turn ticks are
// credited as flat damage (the aura stays in play and heads to next turn's graveyard via
// PlayNextTurn). Without a follow-up attack, no same-turn rune; credit flat n and let
// PlayNextTurn graveyard the aura next turn.
func maleficPlay(s *card.TurnState, n int, self card.Card) int {
	s.AuraCreated = true
	if !followUpAttackAction(s.CardsRemaining) {
		return n
	}
	runes := s.CreateRunechants(1)
	if n == 1 {
		s.AddToGraveyard(self)
		return runes
	}
	return runes + (n - 1)
}

// followUpAttackAction reports whether any CardState in remaining is an attack action card
// (TypeAttack excludes weapons, which carry TypeWeapon on the type line instead).
func followUpAttackAction(remaining []*card.CardState) bool {
	for _, pc := range remaining {
		if pc.Card.Types().Has(card.TypeAttack) {
			return true
		}
	}
	return false
}
