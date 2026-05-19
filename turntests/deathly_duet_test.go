package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

func TestDeathlyDuet_BaseDamage(t *testing.T) {
	// Nothing attributed → just printed power.
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.DeathlyDuetRed{}, 4},
		{cards.DeathlyDuetYellow{}, 3},
		{cards.DeathlyDuetBlue{}, 2},
	}
	for _, tc := range cases {
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

func TestDeathlyDuet_AttackAttributedAddsPower(t *testing.T) {
	// Attack attributed → +2{p}.
	ge := gameengine.New()
	pc := &card.CardState{
		Card:          cards.DeathlyDuetRed{},
		PitchedToPlay: []card.Card{testutils.RunebladeAttack{}},
	}
	ge.ResolveChainStep(ge.Logger(), pc)
	if got := ge.Value(); got != 6 {
		t.Errorf("Deathly Duet Red with attack attributed: Play() = %d, want 6", got)
	}
}

func TestDeathlyDuet_NonAttackActionAttributedCreatesRunechants(t *testing.T) {
	// Non-attack action attributed → 2 Runechant tokens enter play, credited +1 each at creation.
	// Play returns base + 2 (Deathly Duet Red base 4 + 2 token credits = 6). state.Runechants=2
	// for downstream consume bookkeeping.
	ge := gameengine.New()
	pc := &card.CardState{
		Card:          cards.DeathlyDuetRed{},
		PitchedToPlay: []card.Card{testutils.NonAttack{}},
	}
	ge.ResolveChainStep(ge.Logger(), pc)
	if got := ge.Value(); got != 6 {
		t.Errorf("Deathly Duet Red with non-attack attributed: Play() = %d, want 6 (base 4 + 2 token credits)", got)
	}
	if ge.RunechantCount() != 2 {
		t.Errorf("Runechants = %d, want 2", ge.RunechantCount())
	}
	if !ge.AuraCreated() {
		t.Errorf("AuraCreated should be set when Runechants are created")
	}
}

func TestDeathlyDuet_BothBranchesFire(t *testing.T) {
	// Both an attack AND a non-attack action attributed → both riders fire: +2 power bonus,
	// plus 2 Runechants credited +1 each at creation. Play returns base 4 + 2 power + 2 = 8.
	ge := gameengine.New()
	pc := &card.CardState{
		Card:          cards.DeathlyDuetRed{},
		PitchedToPlay: []card.Card{testutils.RunebladeAttack{}, testutils.NonAttack{}},
	}
	ge.ResolveChainStep(ge.Logger(), pc)
	if got := ge.Value(); got != 8 {
		t.Errorf("Deathly Duet Red with both attributed: Play() = %d, want 8 (base 4 + 2 power + 2 token credits)", got)
	}
	if ge.RunechantCount() != 2 {
		t.Errorf("Runechants = %d, want 2", ge.RunechantCount())
	}
}
