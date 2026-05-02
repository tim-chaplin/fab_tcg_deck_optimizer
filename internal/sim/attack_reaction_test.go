package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// stubAR is a minimal Card + AttackReaction stub. The unit tests exercise
// GrantAttackReactionBuff's bookkeeping; ARTargetAllowed is consulted by the chain runner.
type stubAR struct{}

func (stubAR) ID() ids.CardID      { return ids.InvalidCard }
func (stubAR) Name() string        { return "stubAR" }
func (stubAR) Cost(*TurnState) int { return 0 }
func (stubAR) Pitch() int          { return 3 }
func (stubAR) Attack() int         { return 0 }
func (stubAR) Defense() int        { return 0 }
func (stubAR) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)
}
func (stubAR) GoAgain() bool                       { return false }
func (stubAR) ARTargetAllowed(c Card, _ int8) bool { return true }
func (stubAR) Play(*TurnState, *CardState)         {}

// stubAttack is a Generic Action - Attack target candidate.
type stubAttack struct{}

func (stubAttack) ID() ids.CardID      { return ids.InvalidCard }
func (stubAttack) Name() string        { return "stubAttack" }
func (stubAttack) Cost(*TurnState) int { return 0 }
func (stubAttack) Pitch() int          { return 1 }
func (stubAttack) Attack() int         { return 1 }
func (stubAttack) Defense() int        { return 0 }
func (stubAttack) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (stubAttack) GoAgain() bool               { return true }
func (stubAttack) Play(*TurnState, *CardState) {}

// Tests that GrantAttackReactionBuff is a no-op when no target is set.
func TestGrantAttackReactionBuff_NoTargetIsNoOp(t *testing.T) {
	s := TurnState{}
	GrantAttackReactionBuff(&s, &CardState{Card: stubAR{}}, 5)
	if s.Value != 0 {
		t.Errorf("Value = %d, want 0", s.Value)
	}
}

// Tests that GrantAttackReactionBuff buffs BonusAttack, credits Value, and amends the
// target's chain-step log delta.
func TestGrantAttackReactionBuff_AppliesBuffAndCreditsValue(t *testing.T) {
	target := &CardState{Card: stubAttack{}}
	s := TurnState{attackReactionTarget: target}
	s.turnLog = append(s.turnLog, LogEntry{Kind: LogEntryChainStep, Text: "stubAttack: ATTACK", N: 1})
	GrantAttackReactionBuff(&s, &CardState{Card: stubAR{}}, 3)
	if target.BonusAttack != 3 {
		t.Errorf("target BonusAttack = %d, want 3", target.BonusAttack)
	}
	if s.Value != 3 {
		t.Errorf("Value = %d, want 3", s.Value)
	}
	if got := s.turnLog[0].N; got != 4 {
		t.Errorf("amended chain-step N = %d, want 4", got)
	}
}

// Tests that AmendLastChainStepN skips non-chain-step entries to find the most recent
// chain-step.
func TestAmendLastChainStepN_SkipsNonChainEntries(t *testing.T) {
	s := TurnState{}
	s.turnLog = append(s.turnLog,
		LogEntry{Kind: LogEntryChainStep, Text: "first", N: 2},
		LogEntry{Kind: LogEntryPostTrigger, Source: "first", Text: "rider", N: 0},
	)
	s.AmendLastChainStepN(5)
	if got := s.turnLog[0].N; got != 7 {
		t.Errorf("first chain-step N = %d, want 7", got)
	}
	if got := s.turnLog[1].N; got != 0 {
		t.Errorf("post-trigger N = %d, want 0 (untouched)", got)
	}
}
