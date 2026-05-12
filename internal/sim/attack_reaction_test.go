package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// stubAR is a minimal Card + AttackReaction stub. The unit tests exercise
// CardState.GrantAttackReactionBuff's bookkeeping; ARTargetAllowed is consulted by the chain runner.
type stubAR struct{}

func (stubAR) ID() ids.CardID           { return ids.InvalidCard }
func (stubAR) Name() string             { return "stubAR" }
func (stubAR) DisplayName() string      { return "stubAR [B]" }
func (stubAR) Cost(card.GameEngine) int { return 0 }
func (stubAR) Pitch() int               { return 3 }
func (stubAR) Attack() int              { return 0 }
func (stubAR) Defense() int             { return 0 }
func (stubAR) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)
}
func (stubAR) GoAgain(card.GameEngine) bool                          { return false }
func (stubAR) ARTargetAllowed(card.GameEngine, card.Card, int8) bool { return true }
func (stubAR) Play(card.GameEngine, card.Logger, *card.CardState)    {}

// stubAttack is a Generic Action - Attack target candidate.
type stubAttack struct{}

func (stubAttack) ID() ids.CardID           { return ids.InvalidCard }
func (stubAttack) Name() string             { return "stubAttack" }
func (stubAttack) DisplayName() string      { return "stubAttack [R]" }
func (stubAttack) Cost(card.GameEngine) int { return 0 }
func (stubAttack) Pitch() int               { return 1 }
func (stubAttack) Attack() int              { return 1 }
func (stubAttack) Defense() int             { return 0 }
func (stubAttack) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (stubAttack) GoAgain(card.GameEngine) bool                       { return true }
func (stubAttack) Play(card.GameEngine, card.Logger, *card.CardState) {}

// Tests that GrantAttackReactionBuff is a no-op when no target is set.
func TestGrantAttackReactionBuff_NoTargetIsNoOp(t *testing.T) {
	s := TurnState{}
	(&card.CardState{Card: stubAR{}}).GrantAttackReactionBuff(&s, s.Logger(), 5)
	if s.Value() != 0 {
		t.Errorf("Value = %d, want 0", s.Value())
	}
}

// Tests that GrantAttackReactionBuff buffs BonusAttack, credits Value, and amends the
// target's chain-step log delta.
func TestGrantAttackReactionBuff_AppliesBuffAndCreditsValue(t *testing.T) {
	target := &card.CardState{Card: stubAttack{}}
	s := NewTurnStateFromSpec(TurnStateSpec{AttackReactionTarget: target})
	s.logger.AppendChainStep("stubAttack: ATTACK", 1)
	(&card.CardState{Card: stubAR{}}).GrantAttackReactionBuff(&s, s.Logger(), 3)
	if target.BonusAttack != 3 {
		t.Errorf("target BonusAttack = %d, want 3", target.BonusAttack)
	}
	if s.Value() != 3 {
		t.Errorf("Value = %d, want 3", s.Value())
	}
	if got := s.LogEntries()[0].N; got != 4 {
		t.Errorf("amended chain-step N = %d, want 4", got)
	}
}

// Tests that AmendLastChainStepN skips non-chain-step entries to find the most recent
// chain-step.
func TestAmendLastChainStepN_SkipsNonChainEntries(t *testing.T) {
	s := NewTurnStateFromSpec(TurnStateSpec{})
	s.logger.AppendChainStep("first", 2)
	s.logger.AppendPostTrigger("first", "rider", 0)
	s.Logger().AmendLastChainStepN(5)
	entries := s.LogEntries()
	if got := entries[0].N; got != 7 {
		t.Errorf("first chain-step N = %d, want 7", got)
	}
	if got := entries[1].N; got != 0 {
		t.Errorf("post-trigger N = %d, want 0 (untouched)", got)
	}
}
