package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestArcaneCussing_BlockCoversIncomingReturnsN confirms the aura's value is N when the
// partition's block total meets or exceeds incoming damage — we don't take damage, the aura
// survives to pay out later.
func TestArcaneCussing_BlockCoversIncomingReturnsN(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{ArcaneCussingRed{}, 3},
		{ArcaneCussingYellow{}, 2},
		{ArcaneCussingBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{IncomingDamage: 3, BlockTotal: 3}
		if got := tc.c.Play(&s, &card.CardState{}); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d (block == incoming)", tc.c.Name(), got, tc.n)
		}
	}
}

// TestArcaneCussing_OverBlockReturnsN pins that BlockTotal is uncapped — over-blocking still
// counts as covering incoming.
func TestArcaneCussing_OverBlockReturnsN(t *testing.T) {
	s := card.TurnState{IncomingDamage: 3, BlockTotal: 7}
	if got := (ArcaneCussingRed{}).Play(&s, &card.CardState{}); got != 3 {
		t.Errorf("Play() = %d, want 3 (over-block still covers)", got)
	}
}

// TestArcaneCussing_BlockShortReturnsZero confirms the aura collapses to 0 when incoming damage
// gets through — we take damage, aura dies without pay-out, no same-turn attack to save it.
func TestArcaneCussing_BlockShortReturnsZero(t *testing.T) {
	cases := []card.Card{
		ArcaneCussingRed{},
		ArcaneCussingYellow{},
		ArcaneCussingBlue{},
	}
	for _, c := range cases {
		s := card.TurnState{IncomingDamage: 3, BlockTotal: 2}
		if got := c.Play(&s, &card.CardState{}); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (block < incoming, no same-turn pop)", c.Name(), got)
		}
	}
}

// TestArcaneCussing_SameTurnPopBySalientAttack: even if we're taking damage, a later attack
// with a likely-to-hit power pops the aura this turn for its full N.
func TestArcaneCussing_SameTurnPopBySalientAttack(t *testing.T) {
	s := card.TurnState{
		IncomingDamage: 3,
		BlockTotal:     0,
		CardsRemaining: []*card.CardState{{Card: stubAttackWithPower{power: 4}}},
	}
	if got := (ArcaneCussingRed{}).Play(&s, &card.CardState{}); got != 3 {
		t.Errorf("Play() = %d, want 3 (Attack=4 likely to hit, pops Cussing same turn)", got)
	}
}

// TestArcaneCussing_SameTurnPopByWeaponSwing: Cussing's "deal damage" trigger fires off weapon
// swings as well as attack actions, so a following weapon whose swing amount is likely to hit
// pops the aura. Weapon stub has Attack=0, so the weapon swing alone can't pop it — we use a
// carryover Runechant to satisfy the likely-to-hit check.
func TestArcaneCussing_SameTurnPopByWeaponSwing(t *testing.T) {
	s := card.TurnState{
		IncomingDamage: 3,
		BlockTotal:     0,
		Runechants:     1,
		CardsRemaining: []*card.CardState{{Card: stubRunebladeWeapon{}}},
	}
	if got := (ArcaneCussingRed{}).Play(&s, &card.CardState{}); got != 3 {
		t.Errorf("Play() = %d, want 3 (1 Runechant fires with weapon, likely to hit)", got)
	}
}

// TestArcaneCussing_SameTurnPopByRunechantAlone: even an attack whose Attack value is a
// multiple of 3 (blockable) pops the aura if a single Runechant fires alongside it.
func TestArcaneCussing_SameTurnPopByRunechantAlone(t *testing.T) {
	s := card.TurnState{
		IncomingDamage: 3,
		BlockTotal:     0,
		Runechants:     1,
		CardsRemaining: []*card.CardState{{Card: stubAttackWithPower{power: 6}}},
	}
	if got := (ArcaneCussingRed{}).Play(&s, &card.CardState{}); got != 3 {
		t.Errorf("Play() = %d, want 3 (Attack=6 blockable, but 1 Runechant likely to slip through)", got)
	}
}

// TestArcaneCussing_BlockableAttackNoRunechantReturnsZero: a following attack whose power is a
// multiple of 3 (blockable) and no Runechants firing can't pop Cussing — and we're taking
// damage, so value collapses to 0.
func TestArcaneCussing_BlockableAttackNoRunechantReturnsZero(t *testing.T) {
	s := card.TurnState{
		IncomingDamage: 3,
		BlockTotal:     0,
		CardsRemaining: []*card.CardState{{Card: stubAttackWithPower{power: 6}}},
	}
	if got := (ArcaneCussingRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Errorf("Play() = %d, want 0 (Attack=6 blockable, no Runechants, taking damage)", got)
	}
}
