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
	h := []card.Card{fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	got := Best(stubHero{}, nil, h, 4, nil)
	if got.Value != 6 {
		t.Fatalf("want value 6, got %d", got.Value)
	}
}

func TestBest_AllBlueHand(t *testing.T) {
	// Best: pitch 1 blue (3 res), attack with 2 blues (cost 2, dealt 2), defend with 1 blue (prevented
	// 3). Value = 5.
	h := []card.Card{fake.BlueAttack{}, fake.BlueAttack{}, fake.BlueAttack{}, fake.BlueAttack{}}
	got := Best(stubHero{}, nil, h, 4, nil)
	if got.Value != 5 {
		t.Fatalf("want value 5, got %d", got.Value)
	}
}

func TestBest_MixedHand(t *testing.T) {
	// Best: pitch 1 blue (3 res), attack with 2 reds (cost 2, dealt 6), defend with 1 blue (prevented
	// 3). Value = 9.
	h := []card.Card{fake.BlueAttack{}, fake.BlueAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	got := Best(stubHero{}, nil, h, 4, nil)
	if got.Value != 9 {
		t.Fatalf("want value 9, got %d", got.Value)
	}
}

func TestBest_DefenseCappedAtIncoming(t *testing.T) {
	// Best: pitch 1 blue, attack with 2 blues (dealt 2), defend with 1 blue (prevented capped at
	// incoming=2). Value = 4.
	h := []card.Card{fake.BlueAttack{}, fake.BlueAttack{}, fake.BlueAttack{}, fake.BlueAttack{}}
	got := Best(stubHero{}, nil, h, 2, nil)
	if got.Value != 4 {
		t.Fatalf("want value 4, got %d", got.Value)
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
	got := Best(hero.Viserai{}, nil, h, 4, nil)
	if got.Value != 15 {
		t.Fatalf("want value 15, got %d (roles=[%s])",
			got.Value, FormatRoles(h, got.Roles))
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
	got := Best(hero.Viserai{}, weapons, h, 0, nil)
	if got.Value != 8 {
		t.Fatalf("want value 8, got %d (roles=[%s])",
			got.Value, FormatRoles(h, got.Roles))
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
	got := Best(hero.Viserai{}, weapons, h, 0, nil)
	if got.Value != 11 {
		t.Fatalf("want value 11, got %d (roles=[%s])",
			got.Value, FormatRoles(h, got.Roles))
	}
}

func TestBest_ViseraiOathBlueHocusRedMalefic(t *testing.T) {
	// Pitch Blue Hocus Pocus (3 res). Play Red Malefic (3 dmg, go again). Play Red Oath (+1
	// Runechant, peeks ahead and sees the Blade swing = +3 bonus, +1 Viserai Runechant from prior
	// non-attack action = 5). Swing Reaping Blade (cost 1, 3 dmg). Value = 3 + 5 + 3 = 11.
	h := []card.Card{
		runeblade.HocusPocusBlue{},
		runeblade.OathOfTheArknightRed{},
		runeblade.MaleficIncantationRed{},
	}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(hero.Viserai{}, weapons, h, 0, nil)
	if got.Value != 11 {
		t.Fatalf("want value 11, got %d (roles=[%s])",
			got.Value, FormatRoles(h, got.Roles))
	}
}

func TestBest_RunicReapingPrefersAttackPitch(t *testing.T) {
	// Pitching the Blue Hocus Pocus (attack-typed, pitch 3) pays for Runic Reaping + Shrill AND
	// satisfies Runic Reaping's pitched-attack rider. Pitching the Blue Malefic Aura instead would
	// lose the rider. Blue Malefic (1 arcane + 1 Viserai runechant = 2) → Runic Reaping (3 + 1
	// rider + 1 Viserai runechant = 5) → Shrill (4 base + 3 aura-created bonus = 7). Value = 2 + 5
	// + 7 = 14.
	h := []card.Card{
		runeblade.HocusPocusBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.RunicReapingRed{},
		runeblade.ShrillOfSkullformRed{},
	}
	got := Best(hero.Viserai{}, nil, h, 0, nil)
	if got.Value != 14 {
		t.Fatalf("want value 14, got %d (roles=[%s])",
			got.Value, FormatRoles(h, got.Roles))
	}
}

func TestBest_ViseraiMauvrionGrantsGoAgainToShrill(t *testing.T) {
	// Pitch Blue Hocus Pocus (3 res). Play Blue Malefic (1 arcane, go again). Play Red Mauvrion
	// Skies (0 cost, go again; grants go-again to the next Runeblade attack action card = Shrill,
	// and emits 3 runechants). Play Red Shrill (cost 2, 4 base + 3 aura-created bonus = 7; chains
	// thanks to Mauvrion's grant). Swing Reaping Blade (cost 1, 3 dmg). Viserai fires +1 on
	// Mauvrion (prior Malefic is a non-attack action) and +1 on Shrill (priors include non-attack
	// actions). Value = 1 + 3 + 7 + 3 + 2 = 16.
	h := []card.Card{
		runeblade.HocusPocusBlue{},
		runeblade.MaleficIncantationBlue{},
		runeblade.MauvrionSkiesRed{},
		runeblade.ShrillOfSkullformRed{},
	}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(hero.Viserai{}, weapons, h, 0, nil)
	if got.Value != 16 {
		t.Fatalf("want value 16, got %d (roles=[%s])",
			got.Value, FormatRoles(h, got.Roles))
	}
}

func TestCanAfford(t *testing.T) {
	// fake.Red: Cost 1, Pitch 1. fake.Blue: Cost 1, Pitch 3.
	cases := []struct {
		name      string
		pitched   []card.Card
		attackers []card.Card
		want      bool
	}{
		{"empty/empty is trivially affordable", nil, nil, true},
		{"zero pitch covers zero cost", nil, nil, true},
		{"1 Red pitched covers 1 Red attacker (1 == 1)", []card.Card{fake.RedAttack{}}, []card.Card{fake.RedAttack{}}, true},
		{"1 Red pitched can't cover 2 Red attackers (1 < 2)", []card.Card{fake.RedAttack{}}, []card.Card{fake.RedAttack{}, fake.RedAttack{}}, false},
		{"1 Blue pitched covers 3 Red attackers (3 >= 3)", []card.Card{fake.BlueAttack{}}, []card.Card{fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}}, true},
		{"1 Blue pitched can't cover 4 Red attackers (3 < 4)", []card.Card{fake.BlueAttack{}}, []card.Card{fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}}, false},
		{"attackers with 0 cost are always affordable", nil, []card.Card{runeblade.MauvrionSkiesRed{}}, true},
		{"excess resources are fine", []card.Card{fake.BlueAttack{}, fake.BlueAttack{}}, []card.Card{fake.RedAttack{}}, true},
	}
	for _, tc := range cases {
		if got := canAfford(tc.pitched, tc.attackers); got != tc.want {
			t.Errorf("%s: canAfford() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestIsLegalOrder_MauvrionCantSaveShrillWhenRuneragerIsAhead(t *testing.T) {
	// Mauvrion's grant lands on the first matching Runeblade attack action card in CardsRemaining.
	// In the ordering Mauvrion → Runerager → Shrill → weapon, Runerager is that first match, so
	// Shrill never gets the grant. Shrill has no printed go-again, so the Shrill → weapon chain
	// must break — isLegalOrder rejects the ordering.
	order := []card.Card{
		runeblade.MauvrionSkiesRed{},
		runeblade.RuneragerSwarmRed{},
		runeblade.ShrillOfSkullformRed{},
		weapon.ReapingBlade{},
	}
	n := len(order)
	pcBuf := make([]card.PlayedCard, n)
	ptrBuf := make([]*card.PlayedCard, n)
	cpBuf := make([]card.Card, 0, n)
	state := &card.TurnState{}
	if _, legal := playSequence(hero.Viserai{}, nil, nil, order, pcBuf, ptrBuf, cpBuf, state); legal {
		t.Fatalf("ordering %v should be illegal (Shrill has no go-again and Mauvrion granted Runerager instead)",
			cardNames(order))
	}
}

func cardNames(cs []card.Card) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.Name()
	}
	return out
}

func TestBest_ViseraiMauvrionChainsShrillIntoRuneragerIntoWeapon(t *testing.T) {
	// Pitch Blue Hocus → Mauvrion → Shrill → Runerager → Reaping Blade. Value = 3 + 7 + 3 + 3 + 2
	// Viserai runechants = 18.
	h := []card.Card{
		runeblade.HocusPocusBlue{},
		runeblade.MauvrionSkiesRed{},
		runeblade.RuneragerSwarmRed{},
		runeblade.ShrillOfSkullformRed{},
	}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(hero.Viserai{}, weapons, h, 0, nil)
	if got.Value != 18 {
		t.Fatalf("want value 18, got %d (roles=[%s])",
			got.Value, FormatRoles(h, got.Roles))
	}
}

// grantAll is a test-only attacker that sets GrantedGoAgain=true on every PlayedCard remaining in
// CardsRemaining. Used with grantSpy to detect cross-permutation PlayedCard wrapper leakage.
type grantAll struct{}

func (grantAll) Name() string           { return "grantAll" }
func (grantAll) Cost() int              { return 0 }
func (grantAll) Pitch() int              { return 0 }
func (grantAll) Attack() int            { return 0 }
func (grantAll) Defense() int           { return 0 }
func (grantAll) Types() map[string]bool { return map[string]bool{"Runeblade": true, "Action": true, "Attack": true} }
func (grantAll) GoAgain() bool          { return true }
func (grantAll) Play(s *card.TurnState) int {
	for _, pc := range s.CardsRemaining {
		pc.GrantedGoAgain = true
	}
	return 0
}

// grantSpy is a test-only attacker that, when it plays FIRST in a permutation, records whether
// any PlayedCard in CardsRemaining already has GrantedGoAgain=true. With per-permutation fresh
// wrappers, that should never happen (no prior card in this permutation has run yet). If wrappers
// leak across permutations, a grant applied by a previous permutation's grantAll will still be
// visible here — tripping the spy.
type grantSpy struct{ saw *bool }

func (grantSpy) Name() string             { return "grantSpy" }
func (grantSpy) Cost() int                { return 0 }
func (grantSpy) Pitch() int               { return 0 }
func (grantSpy) Attack() int              { return 0 }
func (grantSpy) Defense() int             { return 0 }
func (grantSpy) Types() map[string]bool   { return map[string]bool{"Runeblade": true, "Action": true, "Attack": true} }
func (grantSpy) GoAgain() bool            { return true }
func (g grantSpy) Play(s *card.TurnState) int {
	if len(s.CardsPlayed) != 0 {
		return 0
	}
	for _, pc := range s.CardsRemaining {
		if pc.GrantedGoAgain {
			*g.saw = true
		}
	}
	return 0
}

func TestBestAttackDamage_PlayedCardGrantsDontLeakAcrossPermutations(t *testing.T) {
	// The permutation loop in bestAttackDamage must allocate fresh *PlayedCard wrappers per
	// permutation so a grant applied by one permutation's Play() can't bleed into a later
	// permutation's legality/effect checks.
	//
	// Setup: attackers = [grantAll, grantSpy, grantAll]. The permute order emits grantAll-first
	// permutations before the first grantSpy-first permutation. Each grantAll-first permutation
	// mutates GrantedGoAgain on the wrappers for the cards behind it. When grantSpy later plays
	// FIRST in its own permutation, its CardsRemaining contains the other two cards' wrappers —
	// which must be fresh (GrantedGoAgain=false), since no card has played yet in this permutation.
	// If the wrappers were reused across permutations the spy would see leaked grants and trip.
	var sawLeak bool
	attackers := []card.Card{grantAll{}, grantSpy{saw: &sawLeak}, grantAll{}}
	_ = bestAttackDamage(stubHero{}, attackers, nil, nil, newAttackBufs(0, len(attackers)))
	if sawLeak {
		t.Fatalf("PlayedCard wrapper state leaked across permutations: grantSpy saw a pre-existing GrantedGoAgain when playing first")
	}
}

func TestBest_RespectsResourceConstraint(t *testing.T) {
	// Best: pitch 2 reds (2 res) to attack with 2 reds (cost 2, dealt 6). Value = 6. Resources must
	// cover costs.
	h := []card.Card{fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	got := Best(stubHero{}, nil, h, 0, nil)
	if got.Value != 6 {
		t.Fatalf("want value 6, got %d", got.Value)
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
