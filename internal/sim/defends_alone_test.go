package sim_test

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// aloneBonusBlocker is a 2-defense plain blocker with a configurable DefendsAloneBonus.
type aloneBonusBlocker struct{ bonus int }

func (aloneBonusBlocker) ID() ids.CardID      { return ids.InvalidCard }
func (aloneBonusBlocker) Name() string        { return "aloneBonusBlocker" }
func (aloneBonusBlocker) Cost(*TurnState) int { return 0 }
func (aloneBonusBlocker) Pitch() int          { return 0 }
func (aloneBonusBlocker) Attack() int         { return 0 }
func (aloneBonusBlocker) Defense() int        { return 2 }
func (aloneBonusBlocker) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (aloneBonusBlocker) GoAgain() bool               { return false }
func (b aloneBonusBlocker) DefendsAloneBonus() int    { return b.bonus }
func (aloneBonusBlocker) Play(*TurnState, *CardState) {}

// plainBlocker is a 2-defense plain blocker with no bonus markers.
type plainBlocker struct{}

func (plainBlocker) ID() ids.CardID      { return ids.InvalidCard }
func (plainBlocker) Name() string        { return "plainBlocker" }
func (plainBlocker) Cost(*TurnState) int { return 0 }
func (plainBlocker) Pitch() int          { return 0 }
func (plainBlocker) Attack() int         { return 0 }
func (plainBlocker) Defense() int        { return 2 }
func (plainBlocker) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (plainBlocker) GoAgain() bool               { return false }
func (plainBlocker) Play(*TurnState, *CardState) {}

// togetherBonusBlocker is a 2-defense plain blocker with a configurable DefendsTogetherBonus.
type togetherBonusBlocker struct{ bonus int }

func (togetherBonusBlocker) ID() ids.CardID      { return ids.InvalidCard }
func (togetherBonusBlocker) Name() string        { return "togetherBonusBlocker" }
func (togetherBonusBlocker) Cost(*TurnState) int { return 0 }
func (togetherBonusBlocker) Pitch() int          { return 0 }
func (togetherBonusBlocker) Attack() int         { return 0 }
func (togetherBonusBlocker) Defense() int        { return 2 }
func (togetherBonusBlocker) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (togetherBonusBlocker) GoAgain() bool               { return false }
func (b togetherBonusBlocker) DefendsTogetherBonus() int { return b.bonus }
func (togetherBonusBlocker) Play(*TurnState, *CardState) {}

// blockOneDR is a 1-defense Defense Reaction.
type blockOneDR struct{}

func (blockOneDR) ID() ids.CardID      { return ids.InvalidCard }
func (blockOneDR) Name() string        { return "blockOneDR" }
func (blockOneDR) Cost(*TurnState) int { return 0 }
func (blockOneDR) Pitch() int          { return 0 }
func (blockOneDR) Attack() int         { return 0 }
func (blockOneDR) Defense() int        { return 1 }
func (blockOneDR) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeDefenseReaction)
}
func (blockOneDR) GoAgain() bool { return false }
func (blockOneDR) Play(s *TurnState, self *CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}

// Tests that DefendsAloneBonus fires for a solo plain blocker.
func TestDefendsAloneBonus_FiresWhenSoloPlainBlocker(t *testing.T) {
	bufs := NewAttackBufs(2, 0, nil)
	got, _ := DefendersDamage(
		[]Card{aloneBonusBlocker{bonus: 1}},
		nil, nil,
		bufs.State(),
		bufs.DefenseGravScratch(),
		bufs.DRCardStateScratch(),
		5, -1,
	)
	if got != 3 {
		t.Errorf("solo blocker = %d, want 3 (Defense 2 + alone bonus 1)", got)
	}
}

// Tests that a Defense Reaction blocking alongside doesn't cancel DefendsAloneBonus.
func TestDefendsAloneBonus_DRDoesNotCancel(t *testing.T) {
	bufs := NewAttackBufs(2, 0, nil)
	got, _ := DefendersDamage(
		[]Card{aloneBonusBlocker{bonus: 1}, blockOneDR{}},
		nil, nil,
		bufs.State(),
		bufs.DefenseGravScratch(),
		bufs.DRCardStateScratch(),
		5, -1,
	)
	if got != 4 {
		t.Errorf("plain + DR = %d, want 4 (DR 1 + plain 2 + alone bonus 1)", got)
	}
}

// Tests that a second plain blocker cancels DefendsAloneBonus.
func TestDefendsAloneBonus_SecondPlainBlockerCancels(t *testing.T) {
	bufs := NewAttackBufs(2, 0, nil)
	got, _ := DefendersDamage(
		[]Card{aloneBonusBlocker{bonus: 1}, plainBlocker{}},
		nil, nil,
		bufs.State(),
		bufs.DefenseGravScratch(),
		bufs.DRCardStateScratch(),
		5, -1,
	)
	if got != 4 {
		t.Errorf("two plain blockers = %d, want 4 (no alone bonus, 2+2)", got)
	}
}

// Tests that DefendsAloneBonus is capped by remaining incoming damage.
func TestDefendsAloneBonus_CappedByRemainingDamage(t *testing.T) {
	bufs := NewAttackBufs(2, 0, nil)
	got, _ := DefendersDamage(
		[]Card{aloneBonusBlocker{bonus: 1}},
		nil, nil,
		bufs.State(),
		bufs.DefenseGravScratch(),
		bufs.DRCardStateScratch(),
		2, -1,
	)
	if got != 2 {
		t.Errorf("over-blocked solo blocker = %d, want 2 (capped at incoming)", got)
	}
}

// Tests that DefendsTogetherBonus does not fire when a card with the marker blocks alone.
func TestDefendsTogetherBonus_SuppressedWhenSolo(t *testing.T) {
	bufs := NewAttackBufs(2, 0, nil)
	got, _ := DefendersDamage(
		[]Card{togetherBonusBlocker{bonus: 1}},
		nil, nil,
		bufs.State(),
		bufs.DefenseGravScratch(),
		bufs.DRCardStateScratch(),
		5, -1,
	)
	if got != 2 {
		t.Errorf("solo together-marker blocker = %d, want 2 (no together bonus)", got)
	}
}

// Tests that a Defense Reaction alongside doesn't satisfy DefendsTogetherBonus.
func TestDefendsTogetherBonus_DRDoesNotSatisfy(t *testing.T) {
	bufs := NewAttackBufs(2, 0, nil)
	got, _ := DefendersDamage(
		[]Card{togetherBonusBlocker{bonus: 1}, blockOneDR{}},
		nil, nil,
		bufs.State(),
		bufs.DefenseGravScratch(),
		bufs.DRCardStateScratch(),
		5, -1,
	)
	if got != 3 {
		t.Errorf("together + DR = %d, want 3 (DR 1 + plain 2; DR doesn't satisfy together)", got)
	}
}

// Tests that a second plain blocker satisfies DefendsTogetherBonus.
func TestDefendsTogetherBonus_FiresWithSecondPlainBlocker(t *testing.T) {
	bufs := NewAttackBufs(2, 0, nil)
	got, _ := DefendersDamage(
		[]Card{togetherBonusBlocker{bonus: 1}, plainBlocker{}},
		nil, nil,
		bufs.State(),
		bufs.DefenseGravScratch(),
		bufs.DRCardStateScratch(),
		5, -1,
	)
	if got != 5 {
		t.Errorf("two plain blockers = %d, want 5 (2 + 2 + together bonus 1)", got)
	}
}
