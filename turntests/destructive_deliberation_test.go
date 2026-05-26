package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

func TestDestructiveDeliberation_PlayCreditsAttack(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.DestructiveDeliberationRed{}, 5},
		{cards.DestructiveDeliberationYellow{}, 4},
		{cards.DestructiveDeliberationBlue{}, 3},
	}
	for _, tc := range cases {
		ge := gameengine.New()
		pc := &card.CardState{Card: tc.c}
		ge.ResolveAttackStep(ge.Logger(), pc)
		if ge.Value() != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.c.Name(), ge.Value(), tc.want)
		}
		if len(pc.OnHit) != 1 {
			t.Errorf("%s: OnHit = %d, want 1 (Ponder rider)", tc.c.Name(), len(pc.OnHit))
		}
	}
}

func TestDestructiveDeliberation_OnHitCreatesPonder(t *testing.T) {
	for _, c := range []card.Card{
		cards.DestructiveDeliberationRed{},
		cards.DestructiveDeliberationYellow{},
		cards.DestructiveDeliberationBlue{},
	} {
		ge := gameengine.New()
		pc := &card.CardState{Card: c}
		ge.ResolveAttackStep(ge.Logger(), pc)
		pc.OnHit[0].Fire(ge, ge.Logger(), pc, &pc.OnHit[0])
		if got := ge.PonderCount(); got != 1 {
			t.Errorf("%s: Ponders = %d, want 1", c.Name(), got)
		}
		if !ge.AuraCreated() {
			t.Errorf("%s: AuraCreated = false, want true", c.Name())
		}
	}
}
