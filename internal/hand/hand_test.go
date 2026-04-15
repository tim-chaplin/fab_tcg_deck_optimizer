package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// stubHero is a no-op Hero used by tests that want to measure raw hand
// value without any hero-ability contribution.
type stubHero struct{}

func (stubHero) Name() string                          { return "stubHero" }
func (stubHero) Health() int                           { return 20 }
func (stubHero) Intelligence() int                     { return 4 }
func (stubHero) Types() map[string]bool                { return map[string]bool{} }
func (stubHero) OnCardPlayed(card.Card, *card.TurnState) int { return 0 }

func TestBest_AllRedHand(t *testing.T) {
	// Best: pitch 2 reds (2 res) to attack with the other 2 (cost 2, dealt 6). Value = 6.
	h := []card.Card{fake.Red{}, fake.Red{}, fake.Red{}, fake.Red{}}
	got := Best(stubHero{}, nil, h, 4)
	if got.Value() != 6 {
		t.Fatalf("want value 6, got %d (dealt=%d prevented=%d)", got.Value(), got.Dealt, got.Prevented)
	}
}

func TestBest_AllBlueHand(t *testing.T) {
	// Best: pitch 1 blue (3 res), attack with 2 blues (cost 2, dealt 2), defend with 1 blue (prevented
	// 3). Value = 5.
	h := []card.Card{fake.Blue{}, fake.Blue{}, fake.Blue{}, fake.Blue{}}
	got := Best(stubHero{}, nil, h, 4)
	if got.Value() != 5 {
		t.Fatalf("want value 5, got %d (dealt=%d prevented=%d)", got.Value(), got.Dealt, got.Prevented)
	}
}

func TestBest_MixedHand(t *testing.T) {
	// Best: pitch 1 blue (3 res), attack with 2 reds (cost 2, dealt 6), defend with 1 blue (prevented
	// 3). Value = 9.
	h := []card.Card{fake.Blue{}, fake.Blue{}, fake.Red{}, fake.Red{}}
	got := Best(stubHero{}, nil, h, 4)
	if got.Value() != 9 {
		t.Fatalf("want value 9, got %d (dealt=%d prevented=%d)", got.Value(), got.Dealt, got.Prevented)
	}
}

func TestBest_DefenseCappedAtIncoming(t *testing.T) {
	// Best: pitch 1 blue, attack with 2 blues (dealt 2), defend with 1 blue (prevented capped at
	// incoming=2). Value = 4.
	h := []card.Card{fake.Blue{}, fake.Blue{}, fake.Blue{}, fake.Blue{}}
	got := Best(stubHero{}, nil, h, 2)
	if got.Value() != 4 {
		t.Fatalf("want value 4, got %d (dealt=%d prevented=%d)", got.Value(), got.Dealt, got.Prevented)
	}
}

func TestBest_ViseraiMaleficShrillCombo(t *testing.T) {
	// Hero = Viserai. Best line: pitch the Blue Malefic, then play both Red Maleficas and the Red
	// Shrill. Value = 15.
	h := []card.Card{
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationRed{},
		runeblade.MaleficIncantationRed{},
		runeblade.ShrillOfSkullformRed{},
	}
	got := Best(hero.Viserai{}, nil, h, 4)
	if got.Value() != 15 {
		t.Fatalf("want value 15, got %d (dealt=%d prevented=%d roles=%v)",
			got.Value(), got.Dealt, got.Prevented, got.Roles)
	}
}

func TestBest_ViseraiReapingBladeBlueMalefics(t *testing.T) {
	// Pitch 1 Blue Malefic (3 res), play the other 3 Blue Malefics (Runechants from Viserai on #2
	// and #3), then swing Reaping Blade (cost 1, 3 dmg). Value = 3 + 2 + 3 = 8.
	h := []card.Card{
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
	}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(hero.Viserai{}, weapons, h, 0)
	if got.Value() != 8 {
		t.Fatalf("want value 8, got %d (dealt=%d prevented=%d roles=%v)",
			got.Value(), got.Dealt, got.Prevented, got.Roles)
	}
}

func TestBest_ViseraiReapingBladeMaleficsPlusShrill(t *testing.T) {
	// Pitch 1 Blue Malefic (3 res), play 2 Blue Malefics (2 dmg + 1 Runechant), then Red Shrill
	// (cost 2, 4+3 aura bonus + 1 Runechant = 8). Reaping Blade stays holstered — Shrill has no
	// Go again, so nothing can follow it. Value = 2 + 1 + 8 = 11.
	h := []card.Card{
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.ShrillOfSkullformRed{},
	}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(hero.Viserai{}, weapons, h, 0)
	if got.Value() != 11 {
		t.Fatalf("want value 11, got %d (dealt=%d prevented=%d roles=%v)",
			got.Value(), got.Dealt, got.Prevented, got.Roles)
	}
}

func TestBest_RespectsResourceConstraint(t *testing.T) {
	// Best: pitch 2 reds (2 res) to attack with 2 reds (cost 2, dealt 6). Value = 6. Resources must
	// cover costs.
	h := []card.Card{fake.Red{}, fake.Red{}, fake.Red{}, fake.Red{}}
	got := Best(stubHero{}, nil, h, 0)
	if got.Value() != 6 {
		t.Fatalf("want value 6, got %d", got.Value())
	}
	var res, cost int
	for i, c := range h {
		switch got.Roles[i] {
		case Pitch:
			res += c.Pitch()
		case Attack:
			cost += c.Cost()
		}
	}
	if res < cost {
		t.Fatalf("invalid play: resources %d < costs %d", res, cost)
	}
}
