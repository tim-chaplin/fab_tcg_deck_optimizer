package e2etest

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests Wage Gold's wager via the chain runner. Prior Gold seeds the wager; the
// blue pitch funds the {3} cost. Red (7 power) is in the LikelyToHit window so the
// wager survives — Gold carries to next turn and Value is just the printed power.
// Blue (5 power) is fully blocked in our model, so the wager-loss path fires:
// the Gold goes to opponent and the AddOpponentValue(1) credit nets out one off
// the printed power.
func TestWageGold_WagerOutcomesPerPitch(t *testing.T) {
	cases := []struct {
		name      string
		card      sim.Card
		wantValue int
		wantGold  int
	}{
		{"Red hits, wager survives", cards.WageGoldRed{}, 7, 1},
		{"Blue blocked, wager lost", cards.WageGoldBlue{}, 4, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			deck := fillerDeck()
			hand := []sim.Card{tc.card, testutils.BluePitch{}}
			priorItems := []sim.Item{sim.NewGoldItem(1)}
			got := sim.BestWithTriggers(heroes.Viserai{}, nil, hand, sim.Matchup{IncomingDamage: 0}, deck, nil, nil, priorItems)
			if got.Value != tc.wantValue {
				t.Fatalf("Value = %d, want %d\nBestLine: %s", got.Value, tc.wantValue, formatBestLine(got.BestLine))
			}
			if g := got.State.Gold(); g != tc.wantGold {
				t.Fatalf("end-of-turn Gold = %d, want %d", g, tc.wantGold)
			}
		})
	}
}

// Tests that with no prior Gold, Wage Gold doesn't engage the wager at all — Yellow
// (6 power, fully blocked in our model) just credits its printed power without a
// loss penalty. Confirms the Play-time gate keeps us out of the wager when there's
// nothing to set aside.
func TestWageGold_NoWagerWhenNoGoldHeld(t *testing.T) {
	deck := fillerDeck()
	hand := []sim.Card{cards.WageGoldYellow{}, testutils.BluePitch{}}
	got := sim.BestWithTriggers(heroes.Viserai{}, nil, hand, sim.Matchup{IncomingDamage: 0}, deck, nil, nil, nil)
	if got.Value != 6 {
		t.Fatalf("Value = %d, want 6 (no wager → no penalty even on blocked attack)\nBestLine: %s",
			got.Value, formatBestLine(got.BestLine))
	}
	if g := got.State.Gold(); g != 0 {
		t.Fatalf("end-of-turn Gold = %d, want 0", g)
	}
}
