package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon/weapons"
)

// TestGraveyard_AttackChainAppends: every attacker in the chain lands in state.Graveyard, in
// play order. Confirms the solver actually populates the list as cards resolve.
func TestGraveyard_AttackChainAppends(t *testing.T) {
	order := []card.Card{testutils.FakeRedAttack().WithGoAgain(), testutils.FakeRedAttack().WithGoAgain(), testutils.FakeRedAttack().WithGoAgain()}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 1_000_000, 0, len(order))
	if _, _, _, legal := ctx.PlaySequence(order); !legal {
		t.Fatalf("playSequence rejected the chain")
	}
	got := ctx.PermEngine().Graveyard()
	if len(got) != len(order) {
		t.Fatalf("graveyard len = %d, want %d", len(got), len(order))
	}
	for i := range order {
		if got[i].ID() != order[i].ID() {
			t.Errorf("graveyard[%d] = %s, want %s", i, got[i].Name(), order[i].Name())
		}
	}
}

// TestGraveyard_WeaponSwingDoesNotEnterGraveyard: a weapon ability resolves in the attack
// chain but stays equipped (PersistsInPlay via TypeWeapon); the adjacent action-attack
// still lands in the graveyard.
func TestGraveyard_WeaponSwingDoesNotEnterGraveyard(t *testing.T) {
	attack := testutils.FakeRedAttack().WithGoAgain()
	swing := weapons.ReapingBlade{}.Ability()
	order := []card.Card{attack, swing}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 1_000_000, 0, len(order))
	if _, _, _, legal := ctx.PlaySequence(order); !legal {
		t.Fatalf("playSequence rejected attack → weapon")
	}
	got := ctx.PermEngine().Graveyard()
	if len(got) != 1 {
		t.Fatalf("graveyard len = %d, want 1 (attack only, weapon doesn't enter)", len(got))
	}
	if got[0].ID() != attack.ID() {
		t.Errorf("graveyard[0] = %s, want %s", got[0].Name(), attack.Name())
	}
}

// gravSpyDR is a test-only Defense Reaction whose Play captures a snapshot of gs.Graveyard so
// tests can assert the solver seeded it with the expected defenders before the DR resolved.
type gravSpyDR struct{ saw *[]card.Card }

func (gravSpyDR) ID() ids.CardID           { return ids.InvalidCard }
func (gravSpyDR) Name() string             { return "gravSpyDR" }
func (gravSpyDR) DisplayName() string      { return "gravSpyDR" }
func (gravSpyDR) Cost(card.GameEngine) int { return 0 }
func (gravSpyDR) Pitch() int               { return 0 }
func (gravSpyDR) Attack() int              { return 0 }
func (gravSpyDR) Defense() int             { return 1 }
func (gravSpyDR) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeDefenseReaction)
}
func (gravSpyDR) GoAgain(card.GameEngine) bool { return false }
func (gs gravSpyDR) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	*gs.saw = append((*gs.saw)[:0], ge.(*gameengine.GameEngine).Graveyard()...)
}

// auraDefender is a test-only card whose type line is Aura — a persistent type that normally
// stays in the arena until a destroy condition fires. The test uses it as a plain blocker to
// verify that cards used for raw block value land in the graveyard via the defense-phase
// seeding regardless of their type mask.
type auraDefender struct{}

func (auraDefender) ID() ids.CardID                                     { return ids.InvalidCard }
func (auraDefender) Name() string                                       { return "auraDefender" }
func (auraDefender) DisplayName() string                                { return "auraDefender" }
func (auraDefender) Cost(card.GameEngine) int                           { return 0 }
func (auraDefender) Pitch() int                                         { return 0 }
func (auraDefender) Attack() int                                        { return 0 }
func (auraDefender) Defense() int                                       { return 3 }
func (auraDefender) Types(card.GameEngine) card.TypeSet                 { return card.NewTypeSet(card.TypeAura) }
func (auraDefender) GoAgain(card.GameEngine) bool                       { return false }
func (auraDefender) Play(card.GameEngine, card.Logger, *card.CardState) {}

// Tests that a plain blocker enters the graveyard regardless of PersistsInPlay — a paired
// DR snapshotting state.Graveyard sees the blocker.
func TestGraveyard_PlainBlockEntersGraveyardRegardlessOfType(t *testing.T) {
	blocker := auraDefender{}
	if !blocker.Types(nil).PersistsInPlay() {
		t.Fatal("auraDefender's type mask should set PersistsInPlay; otherwise the test " +
			"isn't isolating the plain-block path")
	}
	var saw []card.Card
	dr := gravSpyDR{saw: &saw}
	bufs := NewAttackBufs(2, 0, nil)
	_, _ = DefendersDamage(
		[]card.Card{blocker, dr},
		nil, nil,
		bufs.State(),
		bufs.DefenseGravScratch(),
		bufs.DRCardStateScratch(),
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

// Tests that Graveyard resets between back-to-back playSequence calls (no double-up).
func TestGraveyard_PermutationReset(t *testing.T) {
	first := []card.Card{testutils.FakeRedAttack().WithGoAgain(), testutils.FakeRedAttack().WithGoAgain(), testutils.FakeRedAttack().WithGoAgain()}
	second := []card.Card{testutils.FakeRedAttack().WithGoAgain()}
	ctx := NewSequenceContextForTest(testutils.Hero{Intel: 4}, nil, nil, 1_000_000, 0, len(first))

	if _, _, _, legal := ctx.PlaySequence(first); !legal {
		t.Fatalf("first playSequence rejected")
	}
	if got := len(ctx.PermEngine().Graveyard()); got != len(first) {
		t.Fatalf("after first run, graveyard len = %d, want %d", got, len(first))
	}

	if _, _, _, legal := ctx.PlaySequence(second); !legal {
		t.Fatalf("second playSequence rejected")
	}
	if got := len(ctx.PermEngine().Graveyard()); got != len(second) {
		t.Fatalf("after second run, graveyard len = %d, want %d (leaked from first?)",
			got, len(second))
	}
}
