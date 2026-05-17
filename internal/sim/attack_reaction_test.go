package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// stubAR is a minimal Card + AttackReaction stub. The unit tests exercise
// CardState.GrantAttackReactionBuff's bookkeeping.
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
	ge := gameengine.New()
	(&card.CardState{Card: stubAR{}}).GrantAttackReactionBuff(ge, ge.Logger(), 5)
	if ge.Value() != 0 {
		t.Errorf("Value = %d, want 0", ge.Value())
	}
}

// Tests that GrantAttackReactionBuff buffs BonusAttack and credits Value.
func TestGrantAttackReactionBuff_AppliesBuffAndCreditsValue(t *testing.T) {
	target := &card.CardState{Card: stubAttack{}}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetAttackReactionTarget(target).Build()}
	(&card.CardState{Card: stubAR{}}).GrantAttackReactionBuff(ge, ge.Logger(), 3)
	if target.BonusAttack != 3 {
		t.Errorf("target BonusAttack = %d, want 3", target.BonusAttack)
	}
	if ge.Value() != 3 {
		t.Errorf("Value = %d, want 3", ge.Value())
	}
}
