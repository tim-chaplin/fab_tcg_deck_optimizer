package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
)

// TestGraveyard_AttackChainAppends: every attacker in the chain lands in state.Graveyard, in
// play order. Confirms the solver actually populates the list as cards resolve.
func TestGraveyard_AttackChainAppends(t *testing.T) {
	order := []card.Card{fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	ctx := newSequenceContextForTest(stubHero{}, nil, nil, 1_000_000, 0, len(order))
	if _, _, _, legal := ctx.playSequence(order, nil, nil); !legal {
		t.Fatalf("playSequence rejected the chain")
	}
	got := ctx.bufs.state.Graveyard
	if len(got) != len(order) {
		t.Fatalf("graveyard len = %d, want %d", len(got), len(order))
	}
	for i := range order {
		if got[i].ID() != order[i].ID() {
			t.Errorf("graveyard[%d] = %s, want %s", i, got[i].Name(), order[i].Name())
		}
	}
}

// TestGraveyard_PermutationReset: running playSequence twice must reset Graveyard between calls.
// Without reset, the second call's list would double-up. A changing chain length between runs
// makes the leak obvious — the second graveyard's length should match the second order.
func TestGraveyard_PermutationReset(t *testing.T) {
	first := []card.Card{fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	second := []card.Card{fake.RedAttack{}}
	ctx := newSequenceContextForTest(stubHero{}, nil, nil, 1_000_000, 0, len(first))

	if _, _, _, legal := ctx.playSequence(first, nil, nil); !legal {
		t.Fatalf("first playSequence rejected")
	}
	if got := len(ctx.bufs.state.Graveyard); got != len(first) {
		t.Fatalf("after first run, graveyard len = %d, want %d", got, len(first))
	}

	if _, _, _, legal := ctx.playSequence(second, nil, nil); !legal {
		t.Fatalf("second playSequence rejected")
	}
	if got := len(ctx.bufs.state.Graveyard); got != len(second) {
		t.Fatalf("after second run, graveyard len = %d, want %d (leaked from first?)",
			got, len(second))
	}
}

// TestBest_WeepingBattlegroundFindsPlainBlockedAura: with Sigil of Silphidae (an Aura with
// Defense 3) and Weeping Battleground (DR, Defense 3) as the hand, the optimal line is to plain-
// block with Sigil and play Weeping Battleground as a DR. Weeping Battleground scans
// state.Graveyard, finds the Sigil aura, banishes it, and deals 1 arcane. Against 6 incoming,
// prevented = 6 and DR damage = 1 → Value = 7.
func TestBest_WeepingBattlegroundFindsPlainBlockedAura(t *testing.T) {
	h := []card.Card{runeblade.SigilOfSilphidaeBlue{}, runeblade.WeepingBattlegroundRed{}}
	got := Best(stubHero{}, nil, h, 6, nil, 0, nil)
	if got.Value != 7 {
		t.Fatalf("want value 7 (6 prevented + 1 from Weeping Battleground banish), got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBest_WeepingBattlegroundFizzlesWithoutAura: Weeping Battleground paired with a non-aura
// defender (Red attack used as a plain block) scans an aura-free graveyard and fizzles. Against
// 3 incoming, prevented = 3 + 1 (Red Defense) = 4 (capped at 3), DR damage = 0 → Value = 3.
func TestBest_WeepingBattlegroundFizzlesWithoutAura(t *testing.T) {
	h := []card.Card{fake.RedAttack{}, runeblade.WeepingBattlegroundRed{}}
	got := Best(stubHero{}, nil, h, 3, nil, 0, nil)
	if got.Value != 3 {
		t.Fatalf("want value 3 (3 prevented + 0 DR damage, no aura to banish), got %d (roles=[%s])",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBest_WeepingBattlegroundIsolatedAcrossCalls: repeated Best calls with the same hand return
// the same Value. Exercises the solver's state buffers for residual mutations from the previous
// call — a leaky graveyard could e.g. keep the Sigil banished into the next Best's defense
// window.
func TestBest_WeepingBattlegroundIsolatedAcrossCalls(t *testing.T) {
	h := []card.Card{runeblade.SigilOfSilphidaeBlue{}, runeblade.WeepingBattlegroundRed{}}
	first := Best(stubHero{}, nil, h, 6, nil, 0, nil)
	second := Best(stubHero{}, nil, h, 6, nil, 0, nil)
	if first.Value != second.Value {
		t.Fatalf("Best repeatability broken: first=%d second=%d", first.Value, second.Value)
	}
}
