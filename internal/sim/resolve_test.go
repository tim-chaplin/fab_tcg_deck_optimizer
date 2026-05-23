package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// Each card-type branch of ResolveChainStep gets one assertion so the standard
// "card resolves and the sim credits / caps / logs" mechanic is pinned in exactly
// one place. Per-card tests don't repeat this — they exercise card-specific
// behaviour (riders, conditional self-buffs) and assume the standard mechanic.

// fakeChainAttack is a vanilla attack-action card with printed power 3 and an empty Play.
type fakeChainAttack struct{}

func (fakeChainAttack) ID() ids.CardID           { return ids.InvalidCard }
func (fakeChainAttack) Name() string             { return "fakeChainAttack" }
func (fakeChainAttack) DisplayName() string      { return "fakeChainAttack" }
func (fakeChainAttack) Cost(card.GameEngine) int { return 0 }
func (fakeChainAttack) Pitch() int               { return 0 }
func (fakeChainAttack) Attack() int              { return 3 }
func (fakeChainAttack) Defense() int             { return 0 }
func (fakeChainAttack) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (fakeChainAttack) GoAgain(card.GameEngine) bool                                 { return false }
func (fakeChainAttack) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

// fakeDR is a vanilla defense-reaction card with printed defense 4.
type fakeDR struct{}

func (fakeDR) ID() ids.CardID           { return ids.InvalidCard }
func (fakeDR) Name() string             { return "fakeDR" }
func (fakeDR) DisplayName() string      { return "fakeDR" }
func (fakeDR) Cost(card.GameEngine) int { return 0 }
func (fakeDR) Pitch() int               { return 0 }
func (fakeDR) Attack() int              { return 0 }
func (fakeDR) Defense() int             { return 4 }
func (fakeDR) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeDefenseReaction)
}
func (fakeDR) GoAgain(card.GameEngine) bool                                 { return false }
func (fakeDR) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

// fakeNonAttack is a non-attack action that flips a flag from Play (used to
// confirm the Play body still runs even though the sim contributes no n).
type fakeNonAttack struct{ played *bool }

func (fakeNonAttack) ID() ids.CardID           { return ids.InvalidCard }
func (fakeNonAttack) Name() string             { return "fakeNonAttack" }
func (fakeNonAttack) DisplayName() string      { return "fakeNonAttack" }
func (fakeNonAttack) Cost(card.GameEngine) int { return 0 }
func (fakeNonAttack) Pitch() int               { return 0 }
func (fakeNonAttack) Attack() int              { return 0 }
func (fakeNonAttack) Defense() int             { return 0 }
func (fakeNonAttack) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction)
}
func (fakeNonAttack) GoAgain(card.GameEngine) bool { return false }
func (n fakeNonAttack) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	*n.played = true
}

// fakeSelfBuff is an attack action whose Play body flips BonusAttack so the
// sim's post-Play EffectiveAttack reads the buffed value. Pins the contract
// that ResolveChainStep computes n AFTER Play returns.
type fakeSelfBuff struct{}

func (fakeSelfBuff) ID() ids.CardID           { return ids.InvalidCard }
func (fakeSelfBuff) Name() string             { return "fakeSelfBuff" }
func (fakeSelfBuff) DisplayName() string      { return "fakeSelfBuff" }
func (fakeSelfBuff) Cost(card.GameEngine) int { return 0 }
func (fakeSelfBuff) Pitch() int               { return 0 }
func (fakeSelfBuff) Attack() int              { return 2 }
func (fakeSelfBuff) Defense() int             { return 0 }
func (fakeSelfBuff) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (fakeSelfBuff) GoAgain(card.GameEngine) bool { return false }
func (fakeSelfBuff) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.BonusAttack += 1
}

func TestResolveChainStep_AttackCreditsEffectiveAttack(t *testing.T) {
	ge := gameengine.New()
	pc := &card.CardState{Card: fakeChainAttack{}}
	ge.ResolveChainStep(ge.Logger(), pc)
	if ge.Value() != 3 {
		t.Errorf("Value = %d, want 3 (printed attack)", ge.Value())
	}
}

func TestResolveChainStep_DefenseReactionCapsToIncomingDamage(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetIncomingDamage(2).Build()}
	pc := &card.CardState{Card: fakeDR{}}
	ge.ResolveChainStep(ge.Logger(), pc)
	if ge.Value() != 2 {
		t.Errorf("Value = %d, want 2 (capped at IncomingDamage)", ge.Value())
	}
	if ge.RemainingUnblockedDamage() != 0 {
		t.Errorf("RemainingUnblockedDamage = %d, want 0 (capped block banked against it)", ge.RemainingUnblockedDamage())
	}
}

func TestResolveChainStep_DefenseReactionUncappedWhenIncomingExceedsDefense(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetIncomingDamage(10).Build()}
	pc := &card.CardState{Card: fakeDR{}}
	ge.ResolveChainStep(ge.Logger(), pc)
	if ge.Value() != 4 {
		t.Errorf("Value = %d, want 4 (printed defense, uncapped)", ge.Value())
	}
	if ge.RemainingUnblockedDamage() != 6 {
		t.Errorf("RemainingUnblockedDamage = %d, want 6 (10 - 4)", ge.RemainingUnblockedDamage())
	}
}

func TestResolveChainStep_NonAttackContributesZero(t *testing.T) {
	played := false
	ge := gameengine.New()
	pc := &card.CardState{Card: fakeNonAttack{played: &played}}
	ge.ResolveChainStep(ge.Logger(), pc)
	if !played {
		t.Error("non-attack Play body did not run")
	}
	if ge.Value() != 0 {
		t.Errorf("Value = %d, want 0", ge.Value())
	}
}

func TestResolveChainStep_SelfBuffInPlayAppliesBeforeCredit(t *testing.T) {
	ge := gameengine.New()
	pc := &card.CardState{Card: fakeSelfBuff{}}
	ge.ResolveChainStep(ge.Logger(), pc)
	if ge.Value() != 3 {
		t.Errorf("Value = %d, want 3 (printed 2 + Play'ge +1 BonusAttack)", ge.Value())
	}
}
