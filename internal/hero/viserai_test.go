package hero

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// stubRuneAttack is a minimal Runeblade attack-action card.
type stubRuneAttack struct{}

func (stubRuneAttack) ID() card.ID             { return card.Invalid }
func (stubRuneAttack) Name() string                { return "StubRuneAttack" }
func (stubRuneAttack) Cost() int                   { return 0 }
func (stubRuneAttack) Pitch() int                  { return 0 }
func (stubRuneAttack) Attack() int                 { return 0 }
func (stubRuneAttack) Defense() int                { return 0 }
func (stubRuneAttack) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (stubRuneAttack) GoAgain() bool            { return true }
func (stubRuneAttack) Play(*card.TurnState) int { return 0 }

// stubRuneAura is a minimal Runeblade non-attack action (an Aura).
type stubRuneAura struct{}

func (stubRuneAura) ID() card.ID             { return card.Invalid }
func (stubRuneAura) Name() string                { return "StubRuneAura" }
func (stubRuneAura) Cost() int                   { return 0 }
func (stubRuneAura) Pitch() int                  { return 0 }
func (stubRuneAura) Attack() int                 { return 0 }
func (stubRuneAura) Defense() int                { return 0 }
func (stubRuneAura) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)
}
func (stubRuneAura) GoAgain() bool            { return true }
func (stubRuneAura) Play(*card.TurnState) int { return 0 }

// stubNonRuneblade is an Action-Attack with no Runeblade type — should never trigger Viserai.
type stubNonRuneblade struct{}

func (stubNonRuneblade) ID() card.ID             { return card.Invalid }
func (stubNonRuneblade) Name() string             { return "StubGeneric" }
func (stubNonRuneblade) Cost() int                { return 0 }
func (stubNonRuneblade) Pitch() int               { return 0 }
func (stubNonRuneblade) Attack() int              { return 0 }
func (stubNonRuneblade) Defense() int             { return 0 }
func (stubNonRuneblade) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (stubNonRuneblade) GoAgain() bool            { return true }
func (stubNonRuneblade) Play(*card.TurnState) int { return 0 }

func TestViserai_RunebladeAfterNonAttackActionTriggers(t *testing.T) {
	// Non-attack action played first, then a Runeblade attack. The second play's OnCardPlayed
	// sees the prior non-attack action and creates a Runechant token on state (consumed
	// downstream by the triggering attack's resolution). OnCardPlayed itself returns 0 damage.
	s := card.TurnState{CardsPlayed: []card.Card{stubRuneAura{}}}
	if got := (Viserai{}).OnCardPlayed(stubRuneAttack{}, &s); got != 0 {
		t.Fatalf("expected 0 damage from OnCardPlayed, got %d", got)
	}
	if s.Runechants != 1 {
		t.Fatalf("expected 1 Runechant on state, got %d", s.Runechants)
	}
}

func TestViserai_NoPriorNonAttackAction(t *testing.T) {
	// Runeblade card, but the only prior play was an attack — no trigger.
	s := card.TurnState{CardsPlayed: []card.Card{stubRuneAttack{}}}
	if got := (Viserai{}).OnCardPlayed(stubRuneAttack{}, &s); got != 0 {
		t.Fatalf("expected 0 (no non-attack action in CardsPlayed), got %d", got)
	}
}

func TestViserai_PlayedCardNotRuneblade(t *testing.T) {
	// Played card isn't Runeblade — Viserai's ability doesn't trigger even if a non-attack action was
	// played earlier.
	s := card.TurnState{CardsPlayed: []card.Card{stubRuneAura{}}}
	if got := (Viserai{}).OnCardPlayed(stubNonRuneblade{}, &s); got != 0 {
		t.Fatalf("expected 0 (non-Runeblade played), got %d", got)
	}
}

// stubRuneWeapon is a Runeblade weapon — tagged with Types["Weapon"] so Viserai should NOT trigger
// when it swings.
type stubRuneWeapon struct{}

func (stubRuneWeapon) ID() card.ID             { return card.Invalid }
func (stubRuneWeapon) Name() string             { return "StubRuneWeapon" }
func (stubRuneWeapon) Cost() int                { return 0 }
func (stubRuneWeapon) Pitch() int               { return 0 }
func (stubRuneWeapon) Attack() int              { return 0 }
func (stubRuneWeapon) Defense() int             { return 0 }
func (stubRuneWeapon) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon)
}
func (stubRuneWeapon) GoAgain() bool            { return true }
func (stubRuneWeapon) Play(*card.TurnState) int { return 0 }

func TestViserai_WeaponSwingDoesNotTrigger(t *testing.T) {
	// Even with a prior non-attack action in CardsPlayed, swinging a Runeblade weapon isn't "playing a
	// card" and must not trigger.
	s := card.TurnState{CardsPlayed: []card.Card{stubRuneAura{}}}
	if got := (Viserai{}).OnCardPlayed(stubRuneWeapon{}, &s); got != 0 {
		t.Fatalf("expected 0 for weapon swing, got %d", got)
	}
}

func TestViserai_EmptyTurn(t *testing.T) {
	// First card of the turn: no prior plays, nothing to trigger on.
	var s card.TurnState
	if got := (Viserai{}).OnCardPlayed(stubRuneAura{}, &s); got != 0 {
		t.Fatalf("expected 0 on empty turn, got %d", got)
	}
}
