package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// Each card-type branch of ResolveChainStep gets one assertion so the standard
// "card resolves and the sim credits / caps / logs" mechanic is pinned in exactly
// one place. Per-card tests don't repeat this — they exercise card-specific
// behaviour (riders, conditional self-buffs) and assume the standard mechanic.

// attackStub is a vanilla attack-action card with printed power 3 and an empty Play.
type attackStub struct{}

func (attackStub) ID() ids.CardID      { return ids.InvalidCard }
func (attackStub) Name() string        { return "attackStub" }
func (attackStub) DisplayName() string { return "attackStub" }
func (attackStub) Cost(*TurnState) int { return 0 }
func (attackStub) Pitch() int          { return 0 }
func (attackStub) Attack() int         { return 3 }
func (attackStub) Defense() int        { return 0 }
func (attackStub) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (attackStub) GoAgain() bool                                { return false }
func (attackStub) Play(s *TurnState, l Logger, self *CardState) {}

// drStub is a vanilla defense-reaction card with printed defense 4.
type drStub struct{}

func (drStub) ID() ids.CardID      { return ids.InvalidCard }
func (drStub) Name() string        { return "drStub" }
func (drStub) DisplayName() string { return "drStub" }
func (drStub) Cost(*TurnState) int { return 0 }
func (drStub) Pitch() int          { return 0 }
func (drStub) Attack() int         { return 0 }
func (drStub) Defense() int        { return 4 }
func (drStub) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeDefenseReaction)
}
func (drStub) GoAgain() bool                                { return false }
func (drStub) Play(s *TurnState, l Logger, self *CardState) {}

// nonAttackStub is a non-attack action that flips a flag from Play (used to
// confirm the Play body still runs even though the sim contributes no n).
type nonAttackStub struct{ played *bool }

func (nonAttackStub) ID() ids.CardID                                 { return ids.InvalidCard }
func (nonAttackStub) Name() string                                   { return "nonAttackStub" }
func (nonAttackStub) DisplayName() string                            { return "nonAttackStub" }
func (nonAttackStub) Cost(*TurnState) int                            { return 0 }
func (nonAttackStub) Pitch() int                                     { return 0 }
func (nonAttackStub) Attack() int                                    { return 0 }
func (nonAttackStub) Defense() int                                   { return 0 }
func (nonAttackStub) Types() card.TypeSet                            { return card.NewTypeSet(card.TypeGeneric, card.TypeAction) }
func (nonAttackStub) GoAgain() bool                                  { return false }
func (n nonAttackStub) Play(s *TurnState, l Logger, self *CardState) { *n.played = true }

// selfBuffStub is an attack action whose Play body flips BonusAttack so the
// sim's post-Play EffectiveAttack reads the buffed value. Pins the contract
// that ResolveChainStep computes n AFTER Play returns.
type selfBuffStub struct{}

func (selfBuffStub) ID() ids.CardID      { return ids.InvalidCard }
func (selfBuffStub) Name() string        { return "selfBuffStub" }
func (selfBuffStub) DisplayName() string { return "selfBuffStub" }
func (selfBuffStub) Cost(*TurnState) int { return 0 }
func (selfBuffStub) Pitch() int          { return 0 }
func (selfBuffStub) Attack() int         { return 2 }
func (selfBuffStub) Defense() int        { return 0 }
func (selfBuffStub) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (selfBuffStub) GoAgain() bool { return false }
func (selfBuffStub) Play(s *TurnState, l Logger, self *CardState) {
	self.BonusAttack += 1
}

func TestResolveChainStep_AttackCreditsEffectiveAttack(t *testing.T) {
	s := &TurnState{logger: turnlogger.New()}
	self := &CardState{Card: attackStub{}}
	ResolveChainStep(s, s.logger, self)
	if s.Value != 3 {
		t.Errorf("Value = %d, want 3 (printed attack)", s.Value)
	}
	if got := s.LogEntries(); len(got) != 1 || got[0].N != 3 {
		t.Errorf("log = %v, want one chain-step entry with N=3", got)
	}
}

func TestResolveChainStep_DefenseReactionCapsToIncomingDamage(t *testing.T) {
	s := &TurnState{logger: turnlogger.New(), IncomingDamage: 2}
	self := &CardState{Card: drStub{}}
	ResolveChainStep(s, s.logger, self)
	if s.Value != 2 {
		t.Errorf("Value = %d, want 2 (capped at IncomingDamage)", s.Value)
	}
	if s.IncomingDamage != 0 {
		t.Errorf("IncomingDamage = %d, want 0 (decremented by capped block)", s.IncomingDamage)
	}
	if got := s.LogEntries(); len(got) != 1 || got[0].N != 2 {
		t.Errorf("log = %v, want one chain-step entry with N=2", got)
	}
}

func TestResolveChainStep_DefenseReactionUncappedWhenIncomingExceedsDefense(t *testing.T) {
	s := &TurnState{logger: turnlogger.New(), IncomingDamage: 10}
	self := &CardState{Card: drStub{}}
	ResolveChainStep(s, s.logger, self)
	if s.Value != 4 {
		t.Errorf("Value = %d, want 4 (printed defense, uncapped)", s.Value)
	}
	if s.IncomingDamage != 6 {
		t.Errorf("IncomingDamage = %d, want 6 (10 - 4)", s.IncomingDamage)
	}
}

func TestResolveChainStep_NonAttackContributesZero(t *testing.T) {
	played := false
	s := &TurnState{logger: turnlogger.New()}
	self := &CardState{Card: nonAttackStub{played: &played}}
	ResolveChainStep(s, s.logger, self)
	if !played {
		t.Error("non-attack Play body did not run")
	}
	if s.Value != 0 {
		t.Errorf("Value = %d, want 0", s.Value)
	}
	if got := s.LogEntries(); len(got) != 1 || got[0].N != 0 {
		t.Errorf("log = %v, want one chain-step entry with N=0", got)
	}
}

func TestResolveChainStep_SelfBuffInPlayAppliesBeforeCredit(t *testing.T) {
	s := &TurnState{logger: turnlogger.New()}
	self := &CardState{Card: selfBuffStub{}}
	ResolveChainStep(s, s.logger, self)
	if s.Value != 3 {
		t.Errorf("Value = %d, want 3 (printed 2 + Play's +1 BonusAttack)", s.Value)
	}
	if got := s.LogEntries(); len(got) != 1 || got[0].N != 3 {
		t.Errorf("log = %v, want one chain-step entry with N=3", got)
	}
}
