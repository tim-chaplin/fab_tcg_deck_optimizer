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

func (fakeAR) ID() ids.CardID           { return ids.InvalidCard }
func (fakeAR) Name() string             { return "fakeAR" }
func (fakeAR) DisplayName() string      { return "fakeAR [B]" }
func (fakeAR) Cost(card.GameEngine) int { return 0 }
func (fakeAR) Pitch() int               { return 3 }
func (fakeAR) Attack() int              { return 0 }
func (fakeAR) Defense() int             { return 0 }
func (fakeAR) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)
}
func (fakeAR) GoAgain(card.GameEngine) bool                          { return false }
func (fakeAR) ARTargetAllowed(card.GameEngine, card.Card, int8) bool { return true }
func (fakeAR) Play(card.GameEngine, card.Logger, *card.CardState)    {}

// fakeAttack is a Generic Action - Attack target candidate.
type fakeAttack struct{}

func (fakeAttack) ID() ids.CardID           { return ids.InvalidCard }
func (fakeAttack) Name() string             { return "fakeAttack" }
func (fakeAttack) DisplayName() string      { return "fakeAttack [R]" }
func (fakeAttack) Cost(card.GameEngine) int { return 0 }
func (fakeAttack) Pitch() int               { return 1 }
func (fakeAttack) Attack() int              { return 1 }
func (fakeAttack) Defense() int             { return 0 }
func (fakeAttack) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (fakeAttack) GoAgain(card.GameEngine) bool                       { return true }
func (fakeAttack) Play(card.GameEngine, card.Logger, *card.CardState) {}

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
	target := &card.CardState{Card: fakeAttack{}}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetAttackReactionTarget(target).Build()}
	(&card.CardState{Card: fakeAR{}}).GrantAttackReactionBuff(ge, ge.Logger(), 3)
	if target.BonusAttack != 3 {
		t.Errorf("target BonusAttack = %d, want 3", target.BonusAttack)
	}
	if ge.Value() != 3 {
		t.Errorf("Value = %d, want 3", ge.Value())
	}
}
