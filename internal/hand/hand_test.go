package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestBest_AllRedHand(t *testing.T) {
	// Best: pitch 2 reds (2 res) to attack with the other 2 (cost 2, dealt 6). Value = 6.
	h := []card.Card{card.TestCardRed, card.TestCardRed, card.TestCardRed, card.TestCardRed}
	got := Best(h, 4)
	if got.Value() != 6 {
		t.Fatalf("want value 6, got %d (dealt=%d prevented=%d)", got.Value(), got.Dealt, got.Prevented)
	}
}

func TestBest_AllBlueHand(t *testing.T) {
	// Best: pitch 1 blue (3 res), attack with 2 blues (cost 2, dealt 2), defend with 1 blue (prevented 3). Value = 5.
	h := []card.Card{card.TestCardBlue, card.TestCardBlue, card.TestCardBlue, card.TestCardBlue}
	got := Best(h, 4)
	if got.Value() != 5 {
		t.Fatalf("want value 5, got %d (dealt=%d prevented=%d)", got.Value(), got.Dealt, got.Prevented)
	}
}

func TestBest_MixedHand(t *testing.T) {
	// Best: pitch 1 blue (3 res), attack with 2 reds (cost 2, dealt 6), defend with 1 blue (prevented 3). Value = 9.
	h := []card.Card{card.TestCardBlue, card.TestCardBlue, card.TestCardRed, card.TestCardRed}
	got := Best(h, 4)
	if got.Value() != 9 {
		t.Fatalf("want value 9, got %d (dealt=%d prevented=%d)", got.Value(), got.Dealt, got.Prevented)
	}
}

func TestBest_DefenseCappedAtIncoming(t *testing.T) {
	// Best: pitch 1 blue, attack with 2 blues (dealt 2), defend with 1 blue (prevented capped at incoming=2). Value = 4.
	h := []card.Card{card.TestCardBlue, card.TestCardBlue, card.TestCardBlue, card.TestCardBlue}
	got := Best(h, 2)
	if got.Value() != 4 {
		t.Fatalf("want value 4, got %d (dealt=%d prevented=%d)", got.Value(), got.Dealt, got.Prevented)
	}
}

func TestBest_RespectsResourceConstraint(t *testing.T) {
	// Best: pitch 2 reds (2 res) to attack with 2 reds (cost 2, dealt 6). Value = 6. Resources must cover costs.
	h := []card.Card{card.TestCardRed, card.TestCardRed, card.TestCardRed, card.TestCardRed}
	got := Best(h, 0)
	if got.Value() != 6 {
		t.Fatalf("want value 6, got %d", got.Value())
	}
	var res, cost int
	for i, c := range h {
		switch got.Roles[i] {
		case Pitch:
			res += c.Pitch
		case Attack:
			cost += c.Cost
		}
	}
	if res < cost {
		t.Fatalf("invalid play: resources %d < costs %d", res, cost)
	}
}
