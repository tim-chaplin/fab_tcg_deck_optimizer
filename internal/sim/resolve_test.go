package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// Each card-type branch of ResolveAttackStep gets one assertion so the standard
// "card resolves and the sim credits / caps / logs" mechanic is pinned in exactly
// one place. Per-card tests don't repeat this — they exercise card-specific
// behaviour (riders, conditional self-buffs) and assume the standard mechanic.

// fakeAttack is a vanilla attack-action card with printed power 3 and an empty Play.
type fakeAttack struct{}

func (fakeAttack) ID() ids.CardID      { return ids.InvalidCard }
func (fakeAttack) Name() string        { return "fakeAttack" }
func (fakeAttack) DisplayName() string { return "fakeAttack" }
func (fakeAttack) Cost() int           { return 0 }
func (fakeAttack) Pitch() int          { return 0 }
func (fakeAttack) Attack() int         { return 3 }
func (fakeAttack) Defense() int        { return 0 }
func (fakeAttack) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (fakeAttack) GoAgain(card.GameEngine) bool                                 { return false }
func (fakeAttack) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

// fakeDR is a vanilla defense-reaction card with printed defense 4.
type fakeDR struct{}

func (fakeDR) ID() ids.CardID      { return ids.InvalidCard }
func (fakeDR) Name() string        { return "fakeDR" }
func (fakeDR) DisplayName() string { return "fakeDR" }
func (fakeDR) Cost() int           { return 0 }
func (fakeDR) Pitch() int          { return 0 }
func (fakeDR) Attack() int         { return 0 }
func (fakeDR) Defense() int        { return 4 }
func (fakeDR) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeDefenseReaction)
}
func (fakeDR) GoAgain(card.GameEngine) bool                                 { return false }
func (fakeDR) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

// fakeNonAttack is a non-attack action that flips a flag from Play (used to
// confirm the Play body still runs even though the sim contributes no n).
type fakeNonAttack struct{ played *bool }

func (fakeNonAttack) ID() ids.CardID      { return ids.InvalidCard }
func (fakeNonAttack) Name() string        { return "fakeNonAttack" }
func (fakeNonAttack) DisplayName() string { return "fakeNonAttack" }
func (fakeNonAttack) Cost() int           { return 0 }
func (fakeNonAttack) Pitch() int          { return 0 }
func (fakeNonAttack) Attack() int         { return 0 }
func (fakeNonAttack) Defense() int        { return 0 }
func (fakeNonAttack) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction)
}
func (fakeNonAttack) GoAgain(card.GameEngine) bool { return false }
func (n fakeNonAttack) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	*n.played = true
}

// fakeSelfBuff is an attack action whose Play body flips BonusAttack so the
// sim's post-Play EffectiveAttack reads the buffed value. Pins the contract
// that ResolveAttackStep computes n AFTER Play returns.
type fakeSelfBuff struct{}

func (fakeSelfBuff) ID() ids.CardID      { return ids.InvalidCard }
func (fakeSelfBuff) Name() string        { return "fakeSelfBuff" }
func (fakeSelfBuff) DisplayName() string { return "fakeSelfBuff" }
func (fakeSelfBuff) Cost() int           { return 0 }
func (fakeSelfBuff) Pitch() int          { return 0 }
func (fakeSelfBuff) Attack() int         { return 2 }
func (fakeSelfBuff) Defense() int        { return 0 }
func (fakeSelfBuff) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (fakeSelfBuff) GoAgain(card.GameEngine) bool { return false }
func (fakeSelfBuff) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += 1
}

func TestResolveAttackStep_AttackCreditsEffectiveAttack(t *testing.T) {
	ge := gameengine.New()
	pc := &card.CardState{Card: fakeAttack{}}
	ge.ResolveAttackStep(ge.Logger(), pc)
	if ge.Value() != 3 {
		t.Errorf("Value = %d, want 3 (printed attack)", ge.Value())
	}
}

func TestResolveAttackStep_DefenseReactionCapsToIncomingDamage(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetIncomingDamage(2).Build()}
	pc := &card.CardState{Card: fakeDR{}}
	ge.ResolveAttackStep(ge.Logger(), pc)
	if ge.Value() != 2 {
		t.Errorf("Value = %d, want 2 (capped at IncomingDamage)", ge.Value())
	}
	if ge.RemainingPhysicalDamage() != 0 {
		t.Errorf("RemainingPhysicalDamage = %d, want 0 (capped block banked against it)", ge.RemainingPhysicalDamage())
	}
}

func TestResolveAttackStep_DefenseReactionUncappedWhenIncomingExceedsDefense(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetIncomingDamage(10).Build()}
	pc := &card.CardState{Card: fakeDR{}}
	ge.ResolveAttackStep(ge.Logger(), pc)
	if ge.Value() != 4 {
		t.Errorf("Value = %d, want 4 (printed defense, uncapped)", ge.Value())
	}
	if ge.RemainingPhysicalDamage() != 6 {
		t.Errorf("RemainingPhysicalDamage = %d, want 6 (10 - 4)", ge.RemainingPhysicalDamage())
	}
}

func TestResolveAttackStep_NonAttackContributesZero(t *testing.T) {
	played := false
	ge := gameengine.New()
	pc := &card.CardState{Card: fakeNonAttack{played: &played}}
	ge.ResolveAttackStep(ge.Logger(), pc)
	if !played {
		t.Error("non-attack Play body did not run")
	}
	if ge.Value() != 0 {
		t.Errorf("Value = %d, want 0", ge.Value())
	}
}

func TestResolveAttackStep_SelfBuffInPlayAppliesBeforeCredit(t *testing.T) {
	ge := gameengine.New()
	pc := &card.CardState{Card: fakeSelfBuff{}}
	ge.ResolveAttackStep(ge.Logger(), pc)
	if ge.Value() != 3 {
		t.Errorf("Value = %d, want 3 (printed 2 + Play'ge +1 BonusAttack)", ge.Value())
	}
}
