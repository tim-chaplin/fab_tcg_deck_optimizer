// Malefic Incantation — Runeblade Action - Aura. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "This enters the arena with N verse counters. When it has none, destroy it. Once per turn,
// when you play an attack action card, remove a verse counter from this. If you do, create a
// Runechant token." (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: if at least one attack action card follows Malefic in this turn's chain, the
// "once per turn" trigger fires — create a live Runechant now and credit n-1 flat damage for
// the remaining verse counters that'll tick on future turns. Otherwise no same-turn tick;
// credit flat n for the full set of future ticks without tracking the aura's persistence.
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
func (MaleficIncantationRed) Play(s *card.TurnState, _ *card.CardState) int {
	return maleficPlay(s, 3)
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
func (MaleficIncantationYellow) Play(s *card.TurnState, _ *card.CardState) int {
	return maleficPlay(s, 2)
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
func (MaleficIncantationBlue) Play(s *card.TurnState, _ *card.CardState) int {
	return maleficPlay(s, 1)
}

// maleficPlay flips AuraCreated for same-turn aura-readers and ticks a verse counter if any
// attack action card follows in this turn's chain — creating a live Runechant that can feed
// later attacks, and crediting n-1 flat damage for the remaining future-turn ticks. Without a
// follow-up attack action, the tick doesn't fire; credit flat n for the full run of future
// ticks instead.
func maleficPlay(s *card.TurnState, n int) int {
	s.AuraCreated = true
	if followUpAttackAction(s.CardsRemaining) {
		return s.CreateRunechants(1) + (n - 1)
	}
	return n
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
