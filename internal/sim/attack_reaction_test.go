package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// fakeAR is a minimal Card + AttackReaction stub. The unit tests exercise
// CardState.GrantAttackReactionBuff's bookkeeping.
type fakeAR struct{}

func (fakeAR) ID() ids.CardID      { return ids.InvalidCard }
func (fakeAR) Name() string        { return "fakeAR" }
func (fakeAR) DisplayName() string { return "fakeAR [B]" }
func (fakeAR) Cost() int           { return 0 }
func (fakeAR) Pitch() int          { return 3 }
func (fakeAR) Attack() int         { return 0 }
func (fakeAR) Defense() int        { return 0 }
func (fakeAR) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)
}
func (fakeAR) GoAgain(card.GameEngine) bool                          { return false }
func (fakeAR) ARTargetAllowed(card.GameEngine, *card.CardState, int8) bool { return true }
func (fakeAR) Play(card.GameEngine, card.Logger, *card.CardState)    {}

// fakeBaseAttack is a Generic Action - Attack target candidate.
type fakeBaseAttack struct{}

func (fakeBaseAttack) ID() ids.CardID      { return ids.InvalidCard }
func (fakeBaseAttack) Name() string        { return "fakeBaseAttack" }
func (fakeBaseAttack) DisplayName() string { return "fakeBaseAttack [R]" }
func (fakeBaseAttack) Cost() int           { return 0 }
func (fakeBaseAttack) Pitch() int          { return 1 }
func (fakeBaseAttack) Attack() int         { return 1 }
func (fakeBaseAttack) Defense() int        { return 0 }
func (fakeBaseAttack) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (fakeBaseAttack) GoAgain(card.GameEngine) bool                       { return true }
func (fakeBaseAttack) Play(card.GameEngine, card.Logger, *card.CardState) {}

// Tests that GrantAttackReactionBuff is a no-op when no target is set.
func TestGrantAttackReactionBuff_NoTargetIsNoOp(t *testing.T) {
	ge := gameengine.New()
	(&card.CardState{Card: fakeAR{}}).GrantAttackReactionBuff(ge, ge.Logger(), 5)
	if ge.Value() != 0 {
		t.Errorf("Value = %d, want 0", ge.Value())
	}
}

// Tests that GrantAttackReactionBuff buffs BonusAttack and credits Value.
func TestGrantAttackReactionBuff_AppliesBuffAndCreditsValue(t *testing.T) {
	target := &card.CardState{Card: fakeBaseAttack{}}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetAttackReactionTarget(target).Build()}
	(&card.CardState{Card: fakeAR{}}).GrantAttackReactionBuff(ge, ge.Logger(), 3)
	if target.BonusAttack != 3 {
		t.Errorf("target BonusAttack = %d, want 3", target.BonusAttack)
	}
	if ge.Value() != 3 {
		t.Errorf("Value = %d, want 3", ge.Value())
	}
}
