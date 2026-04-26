package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// TestGraveyard_AttackChainAppends: every attacker in the chain lands in state.Graveyard, in
// play order. Confirms the solver actually populates the list as cards resolve.
func TestGraveyard_AttackChainAppends(t *testing.T) {
	order := []card.Card{fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 1_000_000, 0, len(order))
	if _, _, _, legal := ctx.playSequence(order); !legal {
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

// TestGraveyard_WeaponSwingDoesNotEnterGraveyard: weapon swings resolve in the attack chain but
// don't hit the graveyard — they stay equipped. The action-attack sitting next to them still
// lands in the graveyard.
func TestGraveyard_WeaponSwingDoesNotEnterGraveyard(t *testing.T) {
	attack := fake.RedAttack{}
	swing := weapon.ReapingBlade{}
	order := []card.Card{attack, swing}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 1_000_000, 0, len(order))
	if _, _, _, legal := ctx.playSequence(order); !legal {
		t.Fatalf("playSequence rejected attack → weapon")
	}
	got := ctx.bufs.state.Graveyard
	if len(got) != 1 {
		t.Fatalf("graveyard len = %d, want 1 (attack only, weapon doesn't enter)", len(got))
	}
	if got[0].ID() != attack.ID() {
		t.Errorf("graveyard[0] = %s, want %s", got[0].Name(), attack.Name())
	}
}

// gravSpyDR is a test-only Defense Reaction whose Play captures a snapshot of s.Graveyard so
// tests can assert the solver seeded it with the expected defenders before the DR resolved.
type gravSpyDR struct{ saw *[]card.Card }

func (gravSpyDR) ID() card.ID              { return card.Invalid }
func (gravSpyDR) Name() string             { return "gravSpyDR" }
func (gravSpyDR) Cost(*card.TurnState) int { return 0 }
func (gravSpyDR) Pitch() int               { return 0 }
func (gravSpyDR) Attack() int              { return 0 }
func (gravSpyDR) Defense() int             { return 1 }
func (gravSpyDR) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeDefenseReaction)
}
func (gravSpyDR) GoAgain() bool { return false }
func (g gravSpyDR) Play(s *card.TurnState, self *card.CardState) {
	*g.saw = append((*g.saw)[:0], s.Graveyard...)
	s.LogPlay(self)
}

// auraDefender is a test-only card whose type line is Aura — a persistent type that normally
// stays in the arena until a destroy condition fires. The test uses it as a plain blocker to
// verify that cards used for raw block value land in the graveyard via the defense-phase
// seeding regardless of their type mask.
type auraDefender struct{}

func (auraDefender) ID() card.ID                           { return card.Invalid }
func (auraDefender) Name() string                          { return "auraDefender" }
func (auraDefender) Cost(*card.TurnState) int              { return 0 }
func (auraDefender) Pitch() int                            { return 0 }
func (auraDefender) Attack() int                           { return 0 }
func (auraDefender) Defense() int                          { return 3 }
func (auraDefender) Types() card.TypeSet                   { return card.NewTypeSet(card.TypeAura) }
func (auraDefender) GoAgain() bool                         { return false }
func (auraDefender) Play(*card.TurnState, *card.CardState) {}

// TestGraveyard_PlainBlockEntersGraveyardRegardlessOfType: a defender whose type mask
// normally keeps it in play still lands in the graveyard the instant it's used to block.
// The test pairs an aura-typed plain blocker with a DR whose Play snapshots state.Graveyard
// — confirming the DR sees the plain blocker in the graveyard alongside itself.
func TestGraveyard_PlainBlockEntersGraveyardRegardlessOfType(t *testing.T) {
	blocker := auraDefender{}
	if !blocker.Types().PersistsInPlay() {
		t.Fatal("auraDefender's type mask should set PersistsInPlay; otherwise the test " +
			"isn't isolating the plain-block path")
	}
	var saw []card.Card
	dr := gravSpyDR{saw: &saw}
	bufs := newAttackBufs(2, 0, nil)
	_, _ = defendersDamage(
		[]card.Card{blocker, dr},
		nil, nil,
		bufs.state,
		bufs.defenseGravScratch,
		&bufs.drCardStateScratch,
		0, -1,
	)
	foundBlocker := false
	for _, c := range saw {
		if _, ok := c.(auraDefender); ok {
			foundBlocker = true
			break
		}
	}
	if !foundBlocker {
		t.Errorf("plain-blocked auraDefender missing from DR's view of state.Graveyard: %v", saw)
	}
}

// TestGraveyard_PermutationReset: running playSequence twice must reset Graveyard between
// calls. Without the reset, the second call's list would double-up. A changing chain length
// between runs makes the leak obvious — the second graveyard's length should match the second
// order.
func TestGraveyard_PermutationReset(t *testing.T) {
	first := []card.Card{fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	second := []card.Card{fake.RedAttack{}}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 1_000_000, 0, len(first))

	if _, _, _, legal := ctx.playSequence(first); !legal {
		t.Fatalf("first playSequence rejected")
	}
	if got := len(ctx.bufs.state.Graveyard); got != len(first) {
		t.Fatalf("after first run, graveyard len = %d, want %d", got, len(first))
	}

	if _, _, _, legal := ctx.playSequence(second); !legal {
		t.Fatalf("second playSequence rejected")
	}
	if got := len(ctx.bufs.state.Graveyard); got != len(second) {
		t.Fatalf("after second run, graveyard len = %d, want %d (leaked from first?)",
			got, len(second))
	}
}
