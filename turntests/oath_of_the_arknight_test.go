package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

func TestOathOfTheArknight_NoRemainingCards(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.OathOfTheArknightRed{}})
	if got := ge.Value(); got != 1 {
		t.Errorf("Play() = %d, want 1 (Runechant only, no attack to buff)", got)
	}
}

func TestOathOfTheArknight_RunebladeAttackInRemaining(t *testing.T) {
	// Oath always creates a Runechant (+1 damage credited to Oath itself). The +N{p} buff
	// rides on the target'ge BonusAttack — so Play returns just the Runechant value, and the
	// target'ge BonusAttack picks up +N.
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
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != 1 {
			t.Errorf("%s: Play() = %d, want 1 (Runechant only; +N rides on target'ge BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.bonus {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.bonus)
		}
	}
}

func TestOathOfTheArknight_WeaponCountsAsAttack(t *testing.T) {
	target := &card.CardState{Card: testutils.RunebladeWeapon{}}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.OathOfTheArknightRed{}})
	if got := ge.Value(); got != 1 {
		t.Errorf("Play() = %d, want 1 (Runechant only; +3 rides on weapon'ge BonusAttack)", got)
	}
	if target.BonusAttack != 3 {
		t.Errorf("weapon BonusAttack = %d, want 3", target.BonusAttack)
	}
}

func TestOathOfTheArknight_NonRunebladeAttackDoesNotQualify(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.NonRunebladeAttack{}}}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.OathOfTheArknightRed{}})
	if got := ge.Value(); got != 1 {
		t.Errorf("Play() = %d, want 1 (non-Runeblade attack shouldn't trigger bonus)", got)
	}
}

func TestOathOfTheArknight_RunebladeNonAttackDoesNotQualify(t *testing.T) {
	// Read the Runes is Runeblade + Action but NOT Attack or Weapon.
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.NonAttack{}}}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.OathOfTheArknightRed{}})
	if got := ge.Value(); got != 1 {
		t.Errorf("Play() = %d, want 1 (non-attack Runeblade card shouldn't trigger bonus)", got)
	}
}
