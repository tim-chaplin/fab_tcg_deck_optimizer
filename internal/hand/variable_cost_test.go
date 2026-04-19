package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// newSequenceContextForTest builds a sequenceContext wired to a fresh attackBufs sized for the
// given chain length. Tests use this instead of hand-rolling the context fields so the common
// shape is centralised.
func newSequenceContextForTest(h hero.Hero, pitched, deck []card.Card, resourceBudget, runechantCarryover, chainLen int) *sequenceContext {
	return &sequenceContext{
		hero:               h,
		pitched:            pitched,
		deck:               deck,
		bufs:               newAttackBufs(chainLen, 0, nil),
		resourceBudget:     resourceBudget,
		runechantCarryover: runechantCarryover,
	}
}

// TestPlaySequence_SetsArcaneDamageDealtWhenRunechantsFire pins the sim-wide contract: when an
// attack or weapon plays with Runechants live, playSequence flips ArcaneDamageDealt on before
// calling the card's Play (so same-hand triggers reading the flag see it). A chain with no
// runechants leaves the flag false. Uses the fake generic attack (package fake) so the test
// observes the playSequence pre-Play hook, not a card's own flag-setting inside its Play.
func TestPlaySequence_SetsArcaneDamageDealtWhenRunechantsFire(t *testing.T) {
	order := []card.Card{fake.RedAttack{}}

	// No runechants → flag stays false.
	ctx := newSequenceContextForTest(hero.Viserai{}, nil, nil, 10, 0, len(order))
	_, _, _, _ = ctx.playSequence(order, nil, nil)
	if ctx.bufs.state.ArcaneDamageDealt {
		t.Errorf("no runechants carried over; expected ArcaneDamageDealt=false, got true")
	}

	// Carryover runechant → fires on the attack → flag set.
	ctx = newSequenceContextForTest(hero.Viserai{}, nil, nil, 10, 1, len(order))
	_, _, _, _ = ctx.playSequence(order, nil, nil)
	if !ctx.bufs.state.ArcaneDamageDealt {
		t.Errorf("runechant carryover fired on attack; expected ArcaneDamageDealt=true, got false")
	}
}

// TestPlaySequence_DiscountRejectsInsufficientBudget verifies that a variable-cost card
// fails its per-play cost check when the sequence's resource budget can't cover the effective
// cost.
func TestPlaySequence_DiscountRejectsInsufficientBudget(t *testing.T) {
	order := []card.Card{runeblade.AmplifyTheArknightRed{}} // printed cost 3, MinCost 0
	ctx := newSequenceContextForTest(hero.Viserai{}, nil, nil, 0, 0, len(order))
	// Resource budget 0, carryover 0 → effective cost = 3 - 0 = 3 > 0, sequence illegal.
	dmg, leftover, _, legal := ctx.playSequence(order, nil, nil)
	if legal {
		t.Fatalf("expected illegal sequence, got legal (dmg=%d, leftover=%d)", dmg, leftover)
	}
}

// TestPlaySequence_DiscountAffordableWithBudget shows the same card becomes legal once the
// budget covers its printed cost.
func TestPlaySequence_DiscountAffordableWithBudget(t *testing.T) {
	order := []card.Card{runeblade.AmplifyTheArknightRed{}}
	ctx := newSequenceContextForTest(hero.Viserai{}, nil, nil, 3, 0, len(order))
	// Resource budget 3, carryover 0 → effective cost 3, budget just covers it. Amplify's
	// Attack(6) is the only damage; no runechants to consume.
	dmg, leftover, _, legal := ctx.playSequence(order, nil, nil)
	if !legal {
		t.Fatalf("expected legal sequence")
	}
	if dmg != 6 {
		t.Errorf("dmg = %d, want 6", dmg)
	}
	if leftover != 0 {
		t.Errorf("leftover = %d, want 0", leftover)
	}
}

// TestPlaySequence_DiscountUsesCarryoverRunechants shows the discount applies from carryover
// tokens — no resource budget needed when there are enough runechants already in play.
func TestPlaySequence_DiscountUsesCarryoverRunechants(t *testing.T) {
	order := []card.Card{runeblade.AmplifyTheArknightRed{}}
	ctx := newSequenceContextForTest(hero.Viserai{}, nil, nil, 0, 3, len(order))
	// Resource budget 0, carryover 3 → effective cost 3-3 = 0, legal. Damage is just Amplify's
	// Attack(); the consumed carryover tokens aren't re-credited (they were credited on the
	// previous turn when they were created).
	dmg, leftover, _, legal := ctx.playSequence(order, nil, nil)
	if !legal {
		t.Fatalf("expected legal sequence")
	}
	if dmg != 6 {
		t.Errorf("dmg = %d, want 6 (Attack only; consumed carryover isn't re-credited)", dmg)
	}
	if leftover != 0 {
		t.Errorf("leftover = %d, want 0 (attack consumes all runechants)", leftover)
	}
}

// TestPlaySequence_LeftoverFromNonAttackAction confirms that runechants created by a non-attack
// action with no following attack persist as leftover, and that their creation credits damage.
func TestPlaySequence_LeftoverFromNonAttackAction(t *testing.T) {
	order := []card.Card{runeblade.ReadTheRunesRed{}} // creates 3 runechants, not an attack
	ctx := newSequenceContextForTest(hero.Viserai{}, nil, nil, 0, 0, len(order))
	dmg, leftover, _, legal := ctx.playSequence(order, nil, nil)
	if !legal {
		t.Fatalf("expected legal sequence")
	}
	if dmg != 3 {
		t.Errorf("dmg = %d, want 3 (3 tokens created, each credited +1)", dmg)
	}
	if leftover != 3 {
		t.Errorf("leftover = %d, want 3", leftover)
	}
}

// TestBest_MauvrionReadNoCarryover exercises carryover bookkeeping end-to-end. Hand is Red
// Mauvrion Skies + Red Read the Runes with Viserai and no starting runechants. Optimal line:
// attack with Mauvrion then Read the Runes — Mauvrion's rider doesn't match (Read isn't an
// attack action), so Mauvrion contributes 0 tokens; Read then creates 3 tokens, and Viserai
// fires on Read (prior Mauvrion is a non-attack action) for +1 more. Total tokens created = 4,
// Value = 4 (each token credited +1 at creation), no attack consumes them → leftover = 4.
func TestBest_MauvrionReadNoCarryover(t *testing.T) {
	h := []card.Card{runeblade.MauvrionSkiesRed{}, runeblade.ReadTheRunesRed{}}
	got := Best(hero.Viserai{}, nil, h, 0, nil, 0, nil)
	if got.Value != 4 {
		t.Errorf("Value = %d, want 4 (3 Read tokens + 1 Viserai token)", got.Value)
	}
	if got.LeftoverRunechants != 4 {
		t.Errorf("LeftoverRunechants = %d, want 4 (non-attack action; no consumption)",
			got.LeftoverRunechants)
	}
}

// TestBest_MauvrionReadWithCarryover is the same hand with 1 runechant carried in from the
// previous turn. The hand still creates 4 tokens this turn, and the 1 carryover token doesn't
// get consumed (no attack in the chain), so leftover = 5.
func TestBest_MauvrionReadWithCarryover(t *testing.T) {
	h := []card.Card{runeblade.MauvrionSkiesRed{}, runeblade.ReadTheRunesRed{}}
	got := Best(hero.Viserai{}, nil, h, 0, nil, 1, nil)
	if got.LeftoverRunechants != 5 {
		t.Errorf("LeftoverRunechants = %d, want 5 (1 carryover + 4 created)", got.LeftoverRunechants)
	}
}

// TestBest_AetherSlashAloneConsumesCarryover covers the attack-consumes case. Hand is a single
// Red Aether Slash with Reaping Blade equipped and 1 runechant carried in. Pitching Aether Slash
// (pitch 1) and swinging the weapon (cost 1, attack 3) is the only legal line. The weapon's
// attack consumes the 1 carryover token without re-crediting damage (the token was credited on
// the turn it was created), so Value = 3 and leftover = 0.
func TestBest_AetherSlashAloneConsumesCarryover(t *testing.T) {
	h := []card.Card{runeblade.AetherSlashRed{}}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(hero.Viserai{}, weapons, h, 0, nil, 1, nil)
	if got.Value != 3 {
		t.Errorf("Value = %d, want 3 (Reaping Blade attack; carryover consumed without credit)", got.Value)
	}
	if got.LeftoverRunechants != 0 {
		t.Errorf("LeftoverRunechants = %d, want 0 (weapon swing consumed the carryover)", got.LeftoverRunechants)
	}
}

// TestBest_BlessingOfOccultTokensOnlyAppearNextTurn shows Blessing of Occult's tokens sit in
// DelayedRunechants and don't interact with same-turn attacks or discount checks: they only
// surface as LeftoverRunechants into the next turn.
func TestBest_BlessingOfOccultTokensOnlyAppearNextTurn(t *testing.T) {
	// Red Malefic (cost 0, pitch 1, flat N=3 with no follow-up) + Red Blessing (cost 1, pitch 1,
	// Play returns 3 via DelayRunechants(3)). Two Value=3 partitions tie (pitch Malefic / attack
	// Blessing, or vice versa); the solver's leftover tie-break picks the Blessing chain because
	// its 3 delayed tokens carry over.
	h := []card.Card{
		runeblade.MaleficIncantationRed{},
		runeblade.BlessingOfOccultRed{},
	}
	got := Best(stubHero{}, nil, h, 0, nil, 0, nil)
	if got.Value != 3 {
		t.Errorf("Value = %d, want 3", got.Value)
	}
	if got.LeftoverRunechants != 3 {
		t.Errorf("LeftoverRunechants = %d, want 3 (0 live + 3 delayed from Blessing)",
			got.LeftoverRunechants)
	}
}

// TestBest_ReduceToRunechantAffordableWithCarryover: a solo Reduce in hand with one Runechant
// already in play can defend — the single carryover discounts PrintedCost 1 down to 0, so the
// partition is affordable with no pitch. Value = 4 prevented + 1 from the token Reduce creates.
func TestBest_ReduceToRunechantAffordableWithCarryover(t *testing.T) {
	h := []card.Card{runeblade.ReduceToRunechantRed{}}
	got := Best(hero.Viserai{}, nil, h, 4, nil, 1, nil)
	if got.Value != 5 {
		t.Errorf("Value = %d, want 5 (Reduce defends at cost 0 thanks to 1 carryover Runechant)", got.Value)
	}
}

// TestBest_ReduceToRunechantUnaffordableWithoutCarryover: the same solo Reduce with zero
// Runechants in play can't be played at all — effective cost is 1 and there's no pitch to cover
// it. The Defend partition is rejected and the best feasible line is pitching Reduce (value 0).
func TestBest_ReduceToRunechantUnaffordableWithoutCarryover(t *testing.T) {
	h := []card.Card{runeblade.ReduceToRunechantRed{}}
	got := Best(hero.Viserai{}, nil, h, 4, nil, 0, nil)
	if got.Value != 0 {
		t.Errorf("Value = %d, want 0 (Reduce can't pay its cost without Runechants or pitch)", got.Value)
	}
}

// TestBest_DiscountAttackerPaysByPitchWithoutCarryover: a variable-cost attack can be
// played by pitching for the full printed cost when no Runechants are available. Amplify
// (PrintedCost 3, Attack 6) + a pitch-3 card with zero carryover should land for 6.
func TestBest_DiscountAttackerPaysByPitchWithoutCarryover(t *testing.T) {
	h := []card.Card{runeblade.AmplifyTheArknightRed{}, fake.BlueAttack{}}
	got := Best(stubHero{}, nil, h, 0, nil, 0 /* carryover */, nil)
	if got.Value != 6 {
		t.Errorf("Value = %d, want 6", got.Value)
	}
}

// TestBest_DiscountAttackerPaysByPartialCarryoverAndTightPitch: Runechants cover part of the
// printed cost, and a tight pitch covers the remainder. Amplify (PrintedCost 3, Attack 6) with
// 2 carryover Runechants (effective cost 1) and a fake pitch-1 card should land for 6.
func TestBest_DiscountAttackerPaysByPartialCarryoverAndTightPitch(t *testing.T) {
	h := []card.Card{runeblade.AmplifyTheArknightRed{}, fake.RedAttack{}}
	got := Best(stubHero{}, nil, h, 0, nil, 2 /* carryover */, nil)
	if got.Value != 6 {
		t.Errorf("Value = %d, want 6", got.Value)
	}
}

// TestBest_DiscountDefenderPaysByPitchWithoutCarryover: a variable-cost defense
// reaction can be played by pitching for the full printed cost when no Runechants are
// available. Reduce (PrintedCost 1, Defense 4, creates one Runechant) + a fake pitch-1 card,
// zero carryover, against 4 incoming should land for 5 (4 prevented + 1 for the created token).
func TestBest_DiscountDefenderPaysByPitchWithoutCarryover(t *testing.T) {
	h := []card.Card{runeblade.ReduceToRunechantRed{}, fake.RedAttack{}}
	got := Best(stubHero{}, nil, h, 4, nil, 0 /* carryover */, nil)
	if got.Value != 5 {
		t.Errorf("Value = %d, want 5", got.Value)
	}
}

// TestBest_CarryoverFeedsDiscount verifies end-to-end: a hand containing a discount attacker is
// playable when the previous turn left enough runechants behind.
func TestBest_CarryoverFeedsDiscount(t *testing.T) {
	// Single Amplify the Arknight (Red): printed cost 3, MinCost 0, Attack 6. With no pitch,
	// resource budget is 0. Without any runechants, effective cost 3 exceeds the budget — so
	// attacking is illegal and Value should be 0.
	h := []card.Card{runeblade.AmplifyTheArknightRed{}}
	got := Best(stubHero{}, nil, h, 0, nil, 0, nil)
	if got.Value != 0 {
		t.Errorf("no carryover: Value = %d, want 0 (discount insufficient without runechants)", got.Value)
	}
	// With 3 runechants carried in, the discount fully covers the cost. Value is just the
	// Attack() power — consumed carryover runechants aren't re-credited.
	got = Best(stubHero{}, nil, h, 0, nil, 3, nil)
	if got.Value != 6 {
		t.Errorf("carryover=3: Value = %d, want 6 (Attack only; carryover tokens don't re-credit)", got.Value)
	}
	if got.LeftoverRunechants != 0 {
		t.Errorf("carryover=3: LeftoverRunechants = %d, want 0 (attack consumes tokens)", got.LeftoverRunechants)
	}
}
