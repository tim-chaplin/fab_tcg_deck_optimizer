package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
)

// grantBonusDamage is a test-only non-attack action card that scans CardsRemaining and adds n
// to BonusDamage on the first attack action card it finds. Mirrors the production shape used
// by Come to Fight / Minnowism / Captain's Call once they migrate to the BonusDamage path:
// the grant lives on the target's CardState rather than being returned from the granter's
// own Play, so the buff is attributed to the attack being buffed and feeds EffectiveAttack
// for any "if this hits" rider on that target.
type grantBonusDamage struct{ n int }

func (grantBonusDamage) ID() card.ID              { return card.Invalid }
func (grantBonusDamage) Name() string             { return "grantBonusDamage" }
func (grantBonusDamage) Cost(*card.TurnState) int { return 0 }
func (grantBonusDamage) Pitch() int               { return 0 }
func (grantBonusDamage) Attack() int              { return 0 }
func (grantBonusDamage) Defense() int             { return 0 }
func (grantBonusDamage) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction)
}
func (grantBonusDamage) GoAgain() bool { return true }
func (g grantBonusDamage) Play(s *card.TurnState, _ *card.CardState) int {
	for _, pc := range s.CardsRemaining {
		if pc.Card.Types().IsAttackAction() {
			pc.BonusDamage += g.n
			return 0
		}
	}
	return 0
}

// grantBonusDamageUngated is grantBonusDamage's misbehaving cousin: it writes BonusDamage to
// every entry in CardsRemaining, attack action or not. Used to exercise the solver-side
// `isAttackAction` gate that protects against a buggy grantor leaking damage into a non-attack
// target's slot.
type grantBonusDamageUngated struct{ n int }

func (grantBonusDamageUngated) ID() card.ID              { return card.Invalid }
func (grantBonusDamageUngated) Name() string             { return "grantBonusDamageUngated" }
func (grantBonusDamageUngated) Cost(*card.TurnState) int { return 0 }
func (grantBonusDamageUngated) Pitch() int               { return 0 }
func (grantBonusDamageUngated) Attack() int              { return 0 }
func (grantBonusDamageUngated) Defense() int             { return 0 }
func (grantBonusDamageUngated) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction)
}
func (grantBonusDamageUngated) GoAgain() bool { return true }
func (g grantBonusDamageUngated) Play(s *card.TurnState, _ *card.CardState) int {
	for _, pc := range s.CardsRemaining {
		pc.BonusDamage += g.n
	}
	return 0
}

// TestPlaySequence_BonusDamageAppliedToTargetDamage pins the core wiring: a granter scheduled
// before an attack action sets BonusDamage on the target's CardState; playSequence folds the
// buff into damage at the target's Play step rather than the granter's, so the chain total
// reflects printed-attack + bonus.
func TestPlaySequence_BonusDamageAppliedToTargetDamage(t *testing.T) {
	order := []card.Card{grantBonusDamage{n: 3}, fake.RedAttack{}}
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

// TestPlaySequence_BonusDamageCreditedToTargetSlot pins per-card attribution: the +N lands in
// the target's perCardOut slot, not the granter's, so chain-display callers see the buff
// credited to the attack receiving it.
func TestPlaySequence_BonusDamageCreditedToTargetSlot(t *testing.T) {
	order := []card.Card{grantBonusDamage{n: 3}, fake.RedAttack{}}
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

// TestPlaySequence_BonusDamageNoTargetFizzles pins the no-target case: a granter alone
// scans CardsRemaining, finds no attack action, and contributes nothing — the BonusDamage
// state simply stays 0.
func TestPlaySequence_BonusDamageNoTargetFizzles(t *testing.T) {
	order := []card.Card{grantBonusDamage{n: 3}}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 10, 0, len(order))
	dmg, _, _, legal := ctx.playSequence(order, nil, nil, nil)
	if !legal {
		t.Fatalf("playSequence returned legal=false; expected granter alone to be a legal 1-card chain")
	}
	if dmg != 0 {
		t.Fatalf("dmg = %d, want 0 (granter has no damage and no target to buff)", dmg)
	}
}

// TestPlaySequence_BonusDamageStacksAcrossGranters pins that two granters in front of the
// same target both write to BonusDamage; the field accumulates rather than overwriting.
func TestPlaySequence_BonusDamageStacksAcrossGranters(t *testing.T) {
	order := []card.Card{grantBonusDamage{n: 3}, grantBonusDamage{n: 2}, fake.RedAttack{}}
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

// TestPlaySequence_BonusDamageSolverGateOnNonAttack pins the solver-side `isAttackAction`
// gate in playSequenceWithMeta. Uses an ungated granter that writes BonusDamage onto every
// CardsRemaining entry — including a non-attack action card scheduled after it. The solver
// must NOT fold that bonus into the non-attack card's perCardOut, only the attack action's
// slot. Without the gate, the non-attack would over-credit by 5 and the chain total would
// rise from 8 to 13.
func TestPlaySequence_BonusDamageSolverGateOnNonAttack(t *testing.T) {
	order := []card.Card{
		grantBonusDamageUngated{n: 5},
		grantBonusDamage{n: 0}, // non-attack action target — solver must skip BonusDamage on it
		fake.RedAttack{},
	}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 10, 0, len(order))
	perCard := make([]float64, len(order))
	dmg, _, _, legal := ctx.playSequence(order, perCard, nil, nil)
	if !legal {
		t.Fatalf("playSequence returned legal=false")
	}
	// Ungated granter writes BonusDamage += 5 to BOTH the non-attack granter (index 1) and
	// the RedAttack (index 2). Only the attack action has the bonus folded into damage:
	//   index 0: 0 (granter's own Play return)
	//   index 1: 0 (non-attack action — gate skips its BonusDamage)
	//   index 2: 3 (printed) + 5 (bonus) = 8
	// Total: 8.
	if dmg != 8 {
		t.Fatalf("dmg = %d, want 8 (only the attack action's BonusDamage applies); perCard=%v", dmg, perCard)
	}
	if perCard[1] != 0 {
		t.Errorf("non-attack target perCardOut = %.1f, want 0 (solver-side isAttackAction gate skipped it)", perCard[1])
	}
	if perCard[2] != 8 {
		t.Errorf("RedAttack perCardOut = %.1f, want 8 (printed 3 + bonus 5)", perCard[2])
	}
}

// TestPlaySequence_BonusDamageNegativeClampsAtZero pins the FaB attack-power floor: a
// negative grant (defender-side -N{p} debuff like Drag Down's printed text) reduces the
// target attack's contribution but never drives it below 0. A 1-power attack with a -3
// grant deals 0, not -2 — the chain total is unchanged below the floor.
func TestPlaySequence_BonusDamageNegativeClampsAtZero(t *testing.T) {
	order := []card.Card{grantBonusDamage{n: -3}, fake.BlueAttack{}}
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

// TestPlaySequence_BonusDamageNegativePartialReduction pins the in-range case: a negative
// grant that doesn't drive the target below 0 reduces the contribution by the full bonus,
// no clamp.
func TestPlaySequence_BonusDamageNegativePartialReduction(t *testing.T) {
	order := []card.Card{grantBonusDamage{n: -2}, fake.RedAttack{}}
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

// TestPlaySequence_BonusDamageNoAttackTargetFizzles pins the granter-side scan: if no attack
// action follows the granter, the rider has nowhere to land and total damage stays 0.
func TestPlaySequence_BonusDamageNoAttackTargetFizzles(t *testing.T) {
	order := []card.Card{grantBonusDamage{n: 3}, grantBonusDamage{n: 2}}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 10, 0, len(order))
	dmg, _, _, legal := ctx.playSequence(order, nil, nil, nil)
	if !legal {
		t.Fatalf("playSequence returned legal=false")
	}
	if dmg != 0 {
		t.Fatalf("dmg = %d, want 0 (no attack actions present; both granters fizzle)", dmg)
	}
}

// TestPlaySequence_BonusDamagePerPermutationReset pins the per-permutation reset contract.
// playSequence rebuilds CardState wrappers fresh per call, but inside one call the
// re-entrant playSequenceWithMeta must zero BonusDamage before reading the chain — otherwise
// a wrapper carried in via pcBuf could leak from a previous run. We verify by running the
// same hand twice through one playSequence (which re-enters playSequenceWithMeta): each run
// must start with BonusDamage = 0 and the totals must match.
func TestPlaySequence_BonusDamagePerPermutationReset(t *testing.T) {
	order := []card.Card{grantBonusDamage{n: 3}, fake.RedAttack{}}
	ctx := newSequenceContextForTest(stubHero, nil, nil, 10, 0, len(order))
	first, _, _, _ := ctx.playSequence(order, nil, nil, nil)
	second, _, _, _ := ctx.playSequence(order, nil, nil, nil)
	if first != 6 || second != 6 {
		t.Fatalf("non-deterministic damage across reuses: first=%d, second=%d, want both=6", first, second)
	}
}
