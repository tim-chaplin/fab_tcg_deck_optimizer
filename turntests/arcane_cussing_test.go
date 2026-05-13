package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestArcaneCussing_BlockCoversIncomingReturnsN confirms the aura's value is N when the
// partition's block total meets or exceeds incoming damage — we don't take damage, the aura
// survives to pay out later.
func TestArcaneCussing_BlockCoversIncomingReturnsN(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{cards.ArcaneCussingRed{}, 3},
		{cards.ArcaneCussingYellow{}, 2},
		{cards.ArcaneCussingBlue{}, 1},
	}
	for _, tc := range cases {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
			SetIncomingDamage(3).
			SetBlockTotal(3).
			Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d (block == incoming)", tc.c.Name(), got, tc.n)
		}
	}
}

// TestArcaneCussing_OverBlockReturnsN pins that BlockTotal is uncapped — over-blocking still
// counts as covering incoming.
func TestArcaneCussing_OverBlockReturnsN(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetIncomingDamage(3).
		SetBlockTotal(7).
		Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.ArcaneCussingRed{}})
	if got := ge.Value(); got != 3 {
		t.Errorf("Play() = %d, want 3 (over-block still covers)", got)
	}
}

// TestArcaneCussing_BlockShortReturnsZero confirms the aura collapses to 0 when incoming damage
// gets through — we take damage, aura dies without pay-out, no same-turn attack to save it.
func TestArcaneCussing_BlockShortReturnsZero(t *testing.T) {
	cases := []card.Card{
		cards.ArcaneCussingRed{},
		cards.ArcaneCussingYellow{},
		cards.ArcaneCussingBlue{},
	}
	for _, c := range cases {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
			SetIncomingDamage(3).
			SetBlockTotal(2).
			Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (block < incoming, no same-turn pop)", c.Name(), got)
		}
	}
}

// TestArcaneCussing_SameTurnPopBySalientAttack: even if we're taking damage, a later attack
// with a likely-to-hit power pops the aura this turn for its full N.
func TestArcaneCussing_SameTurnPopBySalientAttack(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetIncomingDamage(3).
		SetBlockTotal(0).
		SetCardsRemaining([]*card.CardState{{Card: testutils.AttackWithPower{Power: 4}}}).
		Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.ArcaneCussingRed{}})
	if got := ge.Value(); got != 3 {
		t.Errorf("Play() = %d, want 3 (Attack=4 likely to hit, pops Cussing same turn)", got)
	}
}

// Tests that Cussing's pop trigger fires off a likely-to-hit weapon swing.
func TestArcaneCussing_SameTurnPopByWeaponSwing(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetIncomingDamage(3).
		SetBlockTotal(0).
		SetCardsRemaining([]*card.CardState{{Card: testutils.RunebladeWeapon{}}}).
		Build()}
	ge.CreateAura(sim.NewRunechantAura(1))
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.ArcaneCussingRed{}})
	if got := ge.Value(); got != 3 {
		t.Errorf("Play() = %d, want 3 (1 Runechant fires with weapon, likely to hit)", got)
	}
}

// TestArcaneCussing_SameTurnPopByRunechantAlone: even an attack whose Attack value is a
// multiple of 3 (blockable) pops the aura if a single Runechant fires alongside it.
func TestArcaneCussing_SameTurnPopByRunechantAlone(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetIncomingDamage(3).
		SetBlockTotal(0).
		SetCardsRemaining([]*card.CardState{{Card: testutils.AttackWithPower{Power: 6}}}).
		Build()}
	ge.CreateAura(sim.NewRunechantAura(1))
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.ArcaneCussingRed{}})
	if got := ge.Value(); got != 3 {
		t.Errorf("Play() = %d, want 3 (Attack=6 blockable, but 1 Runechant likely to slip through)", got)
	}
}

// TestArcaneCussing_BlockableAttackNoRunechantReturnsZero: a following attack whose power is a
// multiple of 3 (blockable) and no Runechants firing can't pop Cussing — and we're taking
// damage, so value collapses to 0.
func TestArcaneCussing_BlockableAttackNoRunechantReturnsZero(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetIncomingDamage(3).
		SetBlockTotal(0).
		SetCardsRemaining([]*card.CardState{{Card: testutils.AttackWithPower{Power: 6}}}).
		Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.ArcaneCussingRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (Attack=6 blockable, no Runechants, taking damage)", got)
	}
}
