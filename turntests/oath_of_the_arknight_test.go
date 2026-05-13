package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

func TestOathOfTheArknight_NoRemainingCards(t *testing.T) {
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.OathOfTheArknightRed{}})
	if got := s.Value(); got != 1 {
		t.Errorf("Play() = %d, want 1 (Runechant only, no attack to buff)", got)
	}
}

func TestOathOfTheArknight_RunebladeAttackInRemaining(t *testing.T) {
	// Oath always creates a Runechant (+1 damage credited to Oath itself). The +N{p} buff
	// rides on the target's BonusAttack — so Play returns just the Runechant value, and the
	// target's BonusAttack picks up +N.
	cases := []struct {
		c     card.Card
		bonus int
	}{
		{cards.OathOfTheArknightRed{}, 3},
		{cards.OathOfTheArknightYellow{}, 2},
		{cards.OathOfTheArknightBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.RunebladeAttack{}}
		s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != 1 {
			t.Errorf("%s: Play() = %d, want 1 (Runechant only; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.bonus {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.bonus)
		}
	}
}

func TestOathOfTheArknight_WeaponCountsAsAttack(t *testing.T) {
	target := &card.CardState{Card: testutils.RunebladeWeapon{}}
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.OathOfTheArknightRed{}})
	if got := s.Value(); got != 1 {
		t.Errorf("Play() = %d, want 1 (Runechant only; +3 rides on weapon's BonusAttack)", got)
	}
	if target.BonusAttack != 3 {
		t.Errorf("weapon BonusAttack = %d, want 3", target.BonusAttack)
	}
}

func TestOathOfTheArknight_NonRunebladeAttackDoesNotQualify(t *testing.T) {
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.NonRunebladeAttack{}}}).Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.OathOfTheArknightRed{}})
	if got := s.Value(); got != 1 {
		t.Errorf("Play() = %d, want 1 (non-Runeblade attack shouldn't trigger bonus)", got)
	}
}

func TestOathOfTheArknight_RunebladeNonAttackDoesNotQualify(t *testing.T) {
	// Read the Runes is Runeblade + Action but NOT Attack or Weapon.
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.NonAttack{}}}).Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.OathOfTheArknightRed{}})
	if got := s.Value(); got != 1 {
		t.Errorf("Play() = %d, want 1 (non-attack Runeblade card shouldn't trigger bonus)", got)
	}
}
