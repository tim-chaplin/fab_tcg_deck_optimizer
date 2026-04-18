package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestFlyingHigh_NoAttackReturnsZero covers the miss branch: with nothing attack-typed in
// CardsRemaining the grant fizzles and Play returns 0.
func TestFlyingHigh_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	if got := (FlyingHighRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestFlyingHigh_NonAttackInRemainingFizzles confirms a non-attack action in CardsRemaining is
// skipped by the attack-action predicate.
func TestFlyingHigh_NonAttackInRemainingFizzles(t *testing.T) {
	skipped := &card.PlayedCard{Card: stubGenericAction()}
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{skipped}}
	if got := (FlyingHighRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
	if skipped.GrantedGoAgain {
		t.Error("non-attack target should not be granted go again")
	}
}

// TestFlyingHigh_NextAttackRedGrantsGoAgainAndBonus exercises the "target is red" branch: pitch 1
// target gets go again granted and Play returns the +1{p} bonus.
func TestFlyingHigh_NextAttackRedGrantsGoAgainAndBonus(t *testing.T) {
	target := &card.PlayedCard{Card: stubGenericAttackPitch(0, 0, 1)}
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{target}}
	if got := (FlyingHighRed{}).Play(&s); got != 1 {
		t.Errorf("Play() = %d, want 1 (red target → +1{p})", got)
	}
	if !target.GrantedGoAgain {
		t.Error("target GrantedGoAgain = false, want true")
	}
}

// TestFlyingHigh_NextAttackNonRedGrantsGoAgainOnly: yellow (pitch 2) and blue (pitch 3) targets
// still get go again, but the +1{p} rider doesn't fire.
func TestFlyingHigh_NextAttackNonRedGrantsGoAgainOnly(t *testing.T) {
	for _, pitch := range []int{2, 3} {
		target := &card.PlayedCard{Card: stubGenericAttackPitch(0, 0, pitch)}
		s := card.TurnState{CardsRemaining: []*card.PlayedCard{target}}
		if got := (FlyingHighRed{}).Play(&s); got != 0 {
			t.Errorf("pitch %d: Play() = %d, want 0", pitch, got)
		}
		if !target.GrantedGoAgain {
			t.Errorf("pitch %d: target GrantedGoAgain = false, want true", pitch)
		}
	}
}
