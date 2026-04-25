package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// grantBonusAttack is a test-only non-attack action card that scans CardsRemaining and adds n
// to BonusAttack on the first attack action card it finds. Mirrors the production shape used
// by Come to Fight / Minnowism / Captain's Call once they migrate to the BonusAttack path:
// the grant lives on the target's CardState rather than being returned from the granter's
// own Play, so the buff is attributed to the attack being buffed and feeds EffectiveAttack
// for any "if this hits" rider on that target.
type grantBonusAttack struct{ n int }

func (grantBonusAttack) ID() card.ID              { return card.Invalid }
func (grantBonusAttack) Name() string             { return "grantBonusAttack" }
func (grantBonusAttack) Cost(*card.TurnState) int { return 0 }
func (grantBonusAttack) Pitch() int               { return 0 }
func (grantBonusAttack) Attack() int              { return 0 }
func (grantBonusAttack) Defense() int             { return 0 }
func (grantBonusAttack) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction)
}
func (grantBonusAttack) GoAgain() bool { return true }
func (g grantBonusAttack) Play(s *card.TurnState, _ *card.CardState) int {
	for _, pc := range s.CardsRemaining {
		if pc.Card.Types().IsAttackAction() {
			pc.BonusAttack += g.n
			return 0
		}
	}
	return 0
}

// grantBonusAttackWeapon scans CardsRemaining for the first weapon swing (TypeWeapon, no
// TypeAction) and adds n to its BonusAttack. Mirrors the production shape of Brandish's
// "next weapon attack +1{p}" rider — the target is a weapon, not an attack action.
type grantBonusAttackWeapon struct{ n int }

func (grantBonusAttackWeapon) ID() card.ID              { return card.Invalid }
func (grantBonusAttackWeapon) Name() string             { return "grantBonusAttackWeapon" }
func (grantBonusAttackWeapon) Cost(*card.TurnState) int { return 0 }
func (grantBonusAttackWeapon) Pitch() int               { return 0 }
func (grantBonusAttackWeapon) Attack() int              { return 0 }
func (grantBonusAttackWeapon) Defense() int             { return 0 }
func (grantBonusAttackWeapon) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction)
}
func (grantBonusAttackWeapon) GoAgain() bool { return true }
func (g grantBonusAttackWeapon) Play(s *card.TurnState, _ *card.CardState) int {
	for _, pc := range s.CardsRemaining {
		if pc.Card.Types().Has(card.TypeWeapon) {
			pc.BonusAttack += g.n
			return 0
		}
	}
	return 0
}

// TestPlaySequence_BonusAttackAppliedToTargetDamage pins the core wiring: a granter scheduled
// before an attack action sets BonusAttack on the target's CardState; playSequence folds the
// buff into damage at the target's Play step rather than the granter's, so the chain total
// reflects printed-attack + bonus.
func TestPlaySequence_BonusAttackAppliedToTargetDamage(t *testing.T) {
	order := []card.Card{grantBonusAttack{n: 3}, fake.RedAttack{}}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 10, 0, len(order))
	dmg, _, _, legal := ctx.playSequence(order, nil, nil, nil)
	if !legal {
		t.Fatalf("playSequence returned legal=false; expected granter→RedAttack to chain via go-again")
	}
	// Granter (cost 0, attack 0, go again) → RedAttack (cost 1, printed power 3, bonus +3).
	// Total: 0 + (3 + 3) = 6.
	if dmg != 6 {
		t.Fatalf("dmg = %d, want 6 (RedAttack 3 + granted bonus 3)", dmg)
	}
}

// TestPlaySequence_BonusAttackCreditedToTargetSlot pins per-card attribution: the +N lands in
// the target's perCardOut slot, not the granter's, so chain-display callers see the buff
// credited to the attack receiving it.
func TestPlaySequence_BonusAttackCreditedToTargetSlot(t *testing.T) {
	order := []card.Card{grantBonusAttack{n: 3}, fake.RedAttack{}}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 10, 0, len(order))
	perCard := make([]float64, len(order))
	dmg, _, _, legal := ctx.playSequence(order, perCard, nil, nil)
	if !legal {
		t.Fatalf("playSequence returned legal=false")
	}
	if dmg != 6 {
		t.Fatalf("dmg = %d, want 6", dmg)
	}
	if perCard[0] != 0 {
		t.Errorf("granter perCardOut = %.1f, want 0 (granter returns 0; the +3 belongs to the target)", perCard[0])
	}
	if perCard[1] != 6 {
		t.Errorf("RedAttack perCardOut = %.1f, want 6 (printed 3 + bonus 3)", perCard[1])
	}
}

// TestPlaySequence_BonusAttackNoTargetFizzles pins the no-target case: a granter alone
// scans CardsRemaining, finds no attack action, and contributes nothing — the BonusAttack
// state simply stays 0.
func TestPlaySequence_BonusAttackNoTargetFizzles(t *testing.T) {
	order := []card.Card{grantBonusAttack{n: 3}}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 10, 0, len(order))
	dmg, _, _, legal := ctx.playSequence(order, nil, nil, nil)
	if !legal {
		t.Fatalf("playSequence returned legal=false; expected granter alone to be a legal 1-card chain")
	}
	if dmg != 0 {
		t.Fatalf("dmg = %d, want 0 (granter has no damage and no target to buff)", dmg)
	}
}

// TestPlaySequence_BonusAttackStacksAcrossGranters pins that two granters in front of the
// same target both write to BonusAttack; the field accumulates rather than overwriting.
func TestPlaySequence_BonusAttackStacksAcrossGranters(t *testing.T) {
	order := []card.Card{grantBonusAttack{n: 3}, grantBonusAttack{n: 2}, fake.RedAttack{}}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 10, 0, len(order))
	dmg, _, _, legal := ctx.playSequence(order, nil, nil, nil)
	if !legal {
		t.Fatalf("playSequence returned legal=false; expected two granters→RedAttack to chain via go-again")
	}
	// Granter +3 → granter +2 → RedAttack (cost 1, printed power 3, bonus 3+2=5). Total 0+0+8 = 8.
	if dmg != 8 {
		t.Fatalf("dmg = %d, want 8 (RedAttack 3 + stacked grants 5)", dmg)
	}
}

// TestPlaySequence_BonusAttackAppliesToWeapon pins that BonusAttack works on weapon swings,
// not just attack action cards. Brandish, Razor Reflex's sword/dagger branch, Thrust, and
// Visit the Blacksmith all target weapon attacks.
func TestPlaySequence_BonusAttackAppliesToWeapon(t *testing.T) {
	order := []card.Card{grantBonusAttackWeapon{n: 2}, weapon.ReapingBlade{}}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 10, 0, len(order))
	perCard := make([]float64, len(order))
	dmg, _, _, legal := ctx.playSequence(order, perCard, nil, nil)
	if !legal {
		t.Fatalf("playSequence returned legal=false; expected granter→weapon swing to chain via go-again")
	}
	// Granter (cost 0, returns 0) → Reaping Blade (cost 1, printed power 3, bonus +2 = 5).
	// Total: 0 + 5 = 5.
	if dmg != 5 {
		t.Fatalf("dmg = %d, want 5 (Reaping Blade 3 + granted bonus 2); perCard=%v", dmg, perCard)
	}
	if perCard[1] != 5 {
		t.Errorf("Reaping Blade perCardOut = %.1f, want 5 (printed 3 + bonus 2)", perCard[1])
	}
}

// TestPlaySequence_BonusAttackNegativeClampsAtZero pins the FaB attack-power floor: a
// negative grant (defender-side -N{p} debuff like Drag Down's printed text) reduces the
// target attack's contribution but never drives it below 0. A 1-power attack with a -3
// grant deals 0, not -2 — the chain total is unchanged below the floor.
func TestPlaySequence_BonusAttackNegativeClampsAtZero(t *testing.T) {
	order := []card.Card{grantBonusAttack{n: -3}, fake.BlueAttack{}}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 10, 0, len(order))
	perCard := make([]float64, len(order))
	dmg, _, _, legal := ctx.playSequence(order, perCard, nil, nil)
	if !legal {
		t.Fatalf("playSequence returned legal=false")
	}
	// Granter (returns 0, no own attack) → BlueAttack (printed power 1, bonus -3 →
	// pre-clamp -2, post-clamp 0). Total 0+0 = 0.
	if dmg != 0 {
		t.Fatalf("dmg = %d, want 0 (1-power attack with -3 bonus floors at 0)", dmg)
	}
	if perCard[1] != 0 {
		t.Errorf("BlueAttack perCardOut = %.1f, want 0 (clamped from -2)", perCard[1])
	}
}

// TestPlaySequence_BonusAttackNegativePartialReduction pins the in-range case: a negative
// grant that doesn't drive the target below 0 reduces the contribution by the full bonus,
// no clamp.
func TestPlaySequence_BonusAttackNegativePartialReduction(t *testing.T) {
	order := []card.Card{grantBonusAttack{n: -2}, fake.RedAttack{}}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 10, 0, len(order))
	dmg, _, _, legal := ctx.playSequence(order, nil, nil, nil)
	if !legal {
		t.Fatalf("playSequence returned legal=false")
	}
	// Granter → RedAttack (printed power 3, bonus -2 → 1). Total 1.
	if dmg != 1 {
		t.Fatalf("dmg = %d, want 1 (RedAttack 3 - debuff 2)", dmg)
	}
}

// TestPlaySequence_BonusAttackNoAttackTargetFizzles pins the granter-side scan: if no attack
// action follows the granter, the rider has nowhere to land and total damage stays 0.
func TestPlaySequence_BonusAttackNoAttackTargetFizzles(t *testing.T) {
	order := []card.Card{grantBonusAttack{n: 3}, grantBonusAttack{n: 2}}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 10, 0, len(order))
	dmg, _, _, legal := ctx.playSequence(order, nil, nil, nil)
	if !legal {
		t.Fatalf("playSequence returned legal=false")
	}
	if dmg != 0 {
		t.Fatalf("dmg = %d, want 0 (no attack actions present; both granters fizzle)", dmg)
	}
}

// TestPlaySequence_BonusAttackPerPermutationReset pins the per-permutation reset contract.
// playSequence rebuilds CardState wrappers fresh per call, but inside one call the
// re-entrant playSequenceWithMeta must zero BonusAttack before reading the chain — otherwise
// a wrapper carried in via pcBuf could leak from a previous run. We verify by running the
// same hand twice through one playSequence (which re-enters playSequenceWithMeta): each run
// must start with BonusAttack = 0 and the totals must match.
func TestPlaySequence_BonusAttackPerPermutationReset(t *testing.T) {
	order := []card.Card{grantBonusAttack{n: 3}, fake.RedAttack{}}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 10, 0, len(order))
	first, _, _, _ := ctx.playSequence(order, nil, nil, nil)
	second, _, _, _ := ctx.playSequence(order, nil, nil, nil)
	if first != 6 || second != 6 {
		t.Fatalf("non-deterministic damage across reuses: first=%d, second=%d, want both=6", first, second)
	}
}
