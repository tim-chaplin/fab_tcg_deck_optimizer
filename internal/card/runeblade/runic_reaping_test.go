package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// stubRunebladeAttack is a minimal Runeblade Action-Attack used to satisfy Runic Reaping's
// lookahead for a qualifying next attack.
type stubRunebladeAttack struct{}

func (stubRunebladeAttack) Name() string  { return "StubRunebladeAttack" }
func (stubRunebladeAttack) Cost() int     { return 0 }
func (stubRunebladeAttack) Pitch() int    { return 0 }
func (stubRunebladeAttack) Attack() int   { return 0 }
func (stubRunebladeAttack) Defense() int  { return 0 }
func (stubRunebladeAttack) Types() map[string]bool {
	return map[string]bool{"Runeblade": true, "Action": true, "Attack": true}
}
func (stubRunebladeAttack) GoAgain() bool           { return true }
func (stubRunebladeAttack) Play(*card.TurnState) int { return 0 }

// stubRunebladeWeapon is a Runeblade weapon — not an attack action card.
type stubRunebladeWeapon struct{}

func (stubRunebladeWeapon) Name() string  { return "StubRunebladeWeapon" }
func (stubRunebladeWeapon) Cost() int     { return 0 }
func (stubRunebladeWeapon) Pitch() int    { return 0 }
func (stubRunebladeWeapon) Attack() int   { return 0 }
func (stubRunebladeWeapon) Defense() int  { return 0 }
func (stubRunebladeWeapon) Types() map[string]bool {
	return map[string]bool{"Runeblade": true, "Weapon": true}
}
func (stubRunebladeWeapon) GoAgain() bool            { return false }
func (stubRunebladeWeapon) Play(*card.TurnState) int { return 0 }

// stubNonAttackPitch is a non-attack card used to confirm the pitched-attack rider does NOT fire.
type stubNonAttackPitch struct{}

func (stubNonAttackPitch) Name() string                { return "StubNonAttackPitch" }
func (stubNonAttackPitch) Cost() int                   { return 0 }
func (stubNonAttackPitch) Pitch() int                  { return 0 }
func (stubNonAttackPitch) Attack() int                 { return 0 }
func (stubNonAttackPitch) Defense() int                { return 0 }
func (stubNonAttackPitch) Types() map[string]bool      { return map[string]bool{"Action": true} }
func (stubNonAttackPitch) GoAgain() bool               { return false }
func (stubNonAttackPitch) Play(*card.TurnState) int    { return 0 }

func TestRunicReaping_NoNextAttackReturnsZero(t *testing.T) {
	// No attack action following → no bonus at all, and AuraCreated must remain false.
	s := card.TurnState{Pitched: []card.Card{stubRunebladeAttack{}}}
	if got := (RunicReapingRed{}).Play(&s); got != 0 {
		t.Fatalf("want 0 when no next attack, got %d", got)
	}
	if s.AuraCreated {
		t.Fatalf("AuraCreated should stay false when no bonus fires")
	}
}

func TestRunicReaping_WeaponNextDoesNotQualify(t *testing.T) {
	// A Runeblade weapon swing later in the turn is not an attack action card, so the rider doesn't
	// trigger.
	s := card.TurnState{CardsRemaining: []card.Card{stubRunebladeWeapon{}}}
	if got := (RunicReapingRed{}).Play(&s); got != 0 {
		t.Fatalf("want 0 with weapon-only next, got %d", got)
	}
}

func TestRunicReaping_NextAttackNoPitchedAttack(t *testing.T) {
	// Next attack exists, but nothing attack-typed was pitched → just N runechants. Each variant
	// contributes its printed count.
	cases := []struct {
		c    card.Card
		want int
	}{
		{RunicReapingRed{}, 3},
		{RunicReapingYellow{}, 2},
		{RunicReapingBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{
			CardsRemaining: []card.Card{stubRunebladeAttack{}},
			Pitched:        []card.Card{stubNonAttackPitch{}},
		}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should be set when bonus fires", tc.c.Name())
		}
	}
}

func TestRunicReaping_NextAttackWithPitchedAttack(t *testing.T) {
	// Next attack exists AND an attack card was pitched → N+1 (the +1{p} rider stacks on the
	// runechant count).
	cases := []struct {
		c    card.Card
		want int
	}{
		{RunicReapingRed{}, 4},
		{RunicReapingYellow{}, 3},
		{RunicReapingBlue{}, 2},
	}
	for _, tc := range cases {
		s := card.TurnState{
			CardsRemaining: []card.Card{stubRunebladeAttack{}},
			Pitched:        []card.Card{stubRunebladeAttack{}},
		}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
