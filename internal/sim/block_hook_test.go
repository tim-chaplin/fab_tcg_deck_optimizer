package sim_test

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// soloBlocker fires +1{d} when no other plain block shares the defenders slot — the
// "defends alone" pattern (BattlefrontBastion).
type soloBlocker struct{}

func (soloBlocker) ID() ids.CardID      { return ids.InvalidCard }
func (soloBlocker) Name() string        { return "soloBlocker" }
func (soloBlocker) Cost(*TurnState) int { return 0 }
func (soloBlocker) Pitch() int          { return 0 }
func (soloBlocker) Attack() int         { return 0 }
func (soloBlocker) Defense() int        { return 2 }
func (soloBlocker) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (soloBlocker) GoAgain() bool               { return false }
func (soloBlocker) Play(*TurnState, *CardState) {}
func (soloBlocker) Block(s *TurnState, self *CardState) {
	plainCount := 0
	for _, d := range s.Defenders {
		if d.Types().IsDefenseReaction() {
			continue
		}
		plainCount++
		if plainCount > 1 {
			return
		}
	}
	self.BonusDefense += 1
}

// togetherBlocker fires +1{d} when at least one OTHER plain block shares the slot —
// the "defends together" pattern (RightBehindYou).
type togetherBlocker struct{}

func (togetherBlocker) ID() ids.CardID      { return ids.InvalidCard }
func (togetherBlocker) Name() string        { return "togetherBlocker" }
func (togetherBlocker) Cost(*TurnState) int { return 0 }
func (togetherBlocker) Pitch() int          { return 0 }
func (togetherBlocker) Attack() int         { return 0 }
func (togetherBlocker) Defense() int        { return 2 }
func (togetherBlocker) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (togetherBlocker) GoAgain() bool               { return false }
func (togetherBlocker) Play(*TurnState, *CardState) {}
func (togetherBlocker) Block(s *TurnState, self *CardState) {
	plainCount := 0
	for _, d := range s.Defenders {
		if d.Types().IsDefenseReaction() {
			continue
		}
		plainCount++
		if plainCount >= 2 {
			self.BonusDefense += 1
			return
		}
	}
}

// plainBlocker is a vanilla 2-defense plain blocker with no Block hook.
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

// blockOneDR is a 1-defense Defense Reaction. Used to verify DRs alongside a plain blocker
// don't satisfy "another plain block" for either Block-time pattern.
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

// Tests that defendersDamage folds the alone-style Block hook's BonusDefense into the
// lone plain blocker's contribution.
func TestBlockHook_AloneFiresWhenSoloPlainBlocker(t *testing.T) {
	bufs := NewAttackBufs(2, 0, nil)
	got, _ := DefendersDamage(
		[]Card{soloBlocker{}},
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

// Tests that a Defense Reaction blocking alongside doesn't satisfy "another plain block."
func TestBlockHook_AloneStandsAlongsideDR(t *testing.T) {
	bufs := NewAttackBufs(2, 0, nil)
	got, _ := DefendersDamage(
		[]Card{soloBlocker{}, blockOneDR{}},
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

// Tests that a second plain blocker cancels the alone-style bonus.
func TestBlockHook_AloneCancelledBySecondPlainBlocker(t *testing.T) {
	bufs := NewAttackBufs(2, 0, nil)
	got, _ := DefendersDamage(
		[]Card{soloBlocker{}, plainBlocker{}},
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

// Tests that the alone-style bonus is capped by remaining incoming damage.
func TestBlockHook_AloneCappedByRemainingDamage(t *testing.T) {
	bufs := NewAttackBufs(2, 0, nil)
	got, _ := DefendersDamage(
		[]Card{soloBlocker{}},
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

// Tests that the together-style hook stays silent when its card blocks alone.
func TestBlockHook_TogetherSuppressedWhenSolo(t *testing.T) {
	bufs := NewAttackBufs(2, 0, nil)
	got, _ := DefendersDamage(
		[]Card{togetherBlocker{}},
		nil, nil,
		bufs.State(),
		bufs.DefenseGravScratch(),
		bufs.DRCardStateScratch(),
		5, -1,
	)
	if got != 2 {
		t.Errorf("solo together-blocker = %d, want 2 (no together bonus)", got)
	}
}

// Tests that a Defense Reaction alongside doesn't satisfy "another plain block" for the
// together-style hook.
func TestBlockHook_TogetherDRDoesNotSatisfy(t *testing.T) {
	bufs := NewAttackBufs(2, 0, nil)
	got, _ := DefendersDamage(
		[]Card{togetherBlocker{}, blockOneDR{}},
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

// Tests that a second plain blocker satisfies the together-style hook.
func TestBlockHook_TogetherFiresWithSecondPlainBlocker(t *testing.T) {
	bufs := NewAttackBufs(2, 0, nil)
	got, _ := DefendersDamage(
		[]Card{togetherBlocker{}, plainBlocker{}},
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
