package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon/weapons"
)

// Tests that playSequence flips ArcaneDamageDealt before calling Play when Runechants are
// live, and leaves it false otherwise.
func TestPlaySequence_SetsArcaneDamageDealtWhenRunechantsFire(t *testing.T) {
	order := []card.Card{testutils.FakeRedAttack()}

	// No runechants → flag stays false.
	ctx := NewSequenceContextForTest(heroes.Viserai, nil, nil, 10, 0, len(order))
	_, _, _, _ = ctx.PlaySequence(order)
	if ctx.PermEngine().ArcaneDamageDealt() {
		t.Errorf("no runechants carried over; expected ArcaneDamageDealt=false, got true")
	}

	// Carryover runechant → fires on the attack → flag set.
	ctx = NewSequenceContextForTest(heroes.Viserai, nil, nil, 10, 1, len(order))
	_, _, _, _ = ctx.PlaySequence(order)
	if !ctx.PermEngine().ArcaneDamageDealt() {
		t.Errorf("runechant carryover fired on attack; expected ArcaneDamageDealt=true, got false")
	}
}

// TestPlaySequence_DiscountRejectsInsufficientBudget verifies that a variable-cost card
// fails its per-play cost check when the sequence's resource budget can't cover the effective
// cost.
func TestPlaySequence_DiscountRejectsInsufficientBudget(t *testing.T) {
	order := []card.Card{cards.AmplifyTheArknightRed{}} // printed cost 3, MinCost 0
	ctx := NewSequenceContextForTest(heroes.Viserai, nil, nil, 0, 0, len(order))
	// Resource budget 0, carryover 0 → effective cost = 3 - 0 = 3 > 0, sequence illegal.
	dmg, leftover, _, legal := ctx.PlaySequence(order)
	if legal {
		t.Fatalf("expected illegal sequence, got legal (dmg=%d, leftover=%d)", dmg, leftover)
	}
}

// TestPlaySequence_DiscountAffordableWithBudget shows the same card becomes legal once the
// budget covers its printed cost.
func TestPlaySequence_DiscountAffordableWithBudget(t *testing.T) {
	order := []card.Card{cards.AmplifyTheArknightRed{}}
	ctx := NewSequenceContextForTest(heroes.Viserai, nil, nil, 3, 0, len(order))
	// Resource budget 3, carryover 0 → effective cost 3, budget just covers it. Amplify's
	// Attack(6) is the only damage; no runechants to consume.
	dmg, leftover, _, legal := ctx.PlaySequence(order)
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
	order := []card.Card{cards.AmplifyTheArknightRed{}}
	ctx := NewSequenceContextForTest(heroes.Viserai, nil, nil, 0, 3, len(order))
	// Resource budget 0, carryover 3 → effective cost 3-3 = 0, legal. Damage is just Amplify's
	// Attack(); the consumed carryover tokens aren't re-credited (they were credited on the
	// previous turn when they were created).
	dmg, leftover, _, legal := ctx.PlaySequence(order)
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
	order := []card.Card{cards.ReadTheRunesRed{}} // creates 3 runechants, not an attack
	ctx := NewSequenceContextForTest(heroes.Viserai, nil, nil, 0, 0, len(order))
	dmg, leftover, _, legal := ctx.PlaySequence(order)
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

// stateWithRunechants returns a *GameState seeded with hero h and n Runechants on the
// aura list — wrapper for tests that exercise Best with non-zero carryover. n<=0 returns
// the bare state (no runechants).
func stateWithRunechants(h hero.Hero, n int) *gameengine.GameState {
	b := gameengine.GameStateBuilder().SetHero(h)
	if n > 0 {
		b.AddAura(token.NewRunechant(n))
	}
	return b.Build()
}

// Tests carryover bookkeeping end-to-end with no starting runechants — every created token
// is credited once and persists as leftover.
func TestBest_MauvrionReadNoCarryover(t *testing.T) {
	h := []card.Card{cards.MauvrionSkiesRed{}, cards.ReadTheRunesRed{}}
	got := Best(nil, h, nil, gameengine.GameStateBuilder().SetHero(heroes.Viserai).Build())
	if got.Value != 4 {
		t.Errorf("Value = %d, want 4 (3 Read tokens + 1 Viserai token)", got.Value)
	}
	if got.State.RunechantCount() != 4 {
		t.Errorf("leftover Runechants = %d, want 4 (non-attack action; no consumption)",
			got.State.RunechantCount())
	}
}

// TestBest_MauvrionReadWithCarryover is the same hand with 1 runechant carried in from the
// previous turn. The hand still creates 4 tokens this turn, and the 1 carryover token doesn't
// get consumed (no attack in the chain), so leftover = 5.
func TestBest_MauvrionReadWithCarryover(t *testing.T) {
	h := []card.Card{cards.MauvrionSkiesRed{}, cards.ReadTheRunesRed{}}
	got := Best(nil, h, nil, stateWithRunechants(heroes.Viserai, 1))
	if got.State.RunechantCount() != 5 {
		t.Errorf("leftover Runechants = %d, want 5 (1 carryover + 4 created)", got.State.RunechantCount())
	}
}

// Tests that an attack consumes a carryover runechant without re-crediting damage.
func TestBest_AetherSlashAloneConsumesCarryover(t *testing.T) {
	h := []card.Card{cards.AetherSlashRed{}}
	weapons := []weapon.Weapon{weapons.ReapingBlade{}}
	got := Best(weapons, h, nil, stateWithRunechants(heroes.Viserai, 1))
	if got.Value != 3 {
		t.Errorf("Value = %d, want 3 (Reaping Blade attack; carryover consumed without credit)", got.Value)
	}
	if got.State.RunechantCount() != 0 {
		t.Errorf("leftover Runechants = %d, want 0 (weapon swing consumed the carryover)", got.State.RunechantCount())
	}
}

// Tests that Blessing's start-of-turn-deferred runes don't appear in the same-turn chain.
func TestBest_BlessingOfOccultTokensDoNotAffectSameTurnChain(t *testing.T) {
	h := []card.Card{
		cards.MaleficIncantationRed{},
		cards.BlessingOfOccultRed{},
	}
	got := Best(nil, h, nil, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	if got.Value != 0 {
		t.Errorf("Value = %d, want 0 (Malefic needs an attack action to fire; Blessing is deferred)", got.Value)
	}
	if got.State.RunechantCount() != 0 {
		t.Errorf("leftover Runechants = %d, want 0 (Blessing's runes materialise via start-of-turn trigger, not carryover)",
			got.State.RunechantCount())
	}
}

// Tests that a single carryover runechant covers Reduce's printed cost so it can defend with
// no pitch.
func TestBest_ReduceToRunechantAffordableWithCarryover(t *testing.T) {
	h := []card.Card{cards.ReduceToRunechantRed{}}
	state := stateWithRunechants(heroes.Viserai, 1)
	state.SetIncomingDamage(4)
	got := Best(nil, h, nil, state)
	if got.Value != 5 {
		t.Errorf("Value = %d, want 5 (Reduce defends at cost 0 thanks to 1 carryover Runechant)", got.Value)
	}
}

// Tests that solo Reduce with no carryover and no pitch source is unplayable, falling back to
// a pitch (Value 0).
func TestBest_ReduceToRunechantUnaffordableWithoutCarryover(t *testing.T) {
	h := []card.Card{cards.ReduceToRunechantRed{}}
	got := Best(nil, h, nil, gameengine.GameStateBuilder().
		SetHero(heroes.Viserai).
		SetIncomingDamage(4).
		Build())
	if got.Value != 0 {
		t.Errorf("Value = %d, want 0 (Reduce can't pay its cost without Runechants or pitch)", got.Value)
	}
}

// Tests that a variable-cost attack pays its full printed cost by pitch when no runechants
// are available.
func TestBest_DiscountAttackerPaysByPitchWithoutCarryover(t *testing.T) {
	h := []card.Card{cards.AmplifyTheArknightRed{}, testutils.FakeBlueAttack()}
	got := Best(nil, h, nil, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	if got.Value != 6 {
		t.Errorf("Value = %d, want 6", got.Value)
	}
}

// Tests that runechants cover part of the printed cost and a tight pitch covers the rest.
func TestBest_DiscountAttackerPaysByPartialCarryoverAndTightPitch(t *testing.T) {
	h := []card.Card{cards.AmplifyTheArknightRed{}, testutils.FakeRedAttack()}
	got := Best(nil, h, nil, stateWithRunechants(testutils.Hero{Intel: 4}, 2))
	if got.Value != 6 {
		t.Errorf("Value = %d, want 6", got.Value)
	}
}

// Tests that a variable-cost defense reaction pays its full printed cost by pitch when no
// runechants are available.
func TestBest_DiscountDefenderPaysByPitchWithoutCarryover(t *testing.T) {
	h := []card.Card{cards.ReduceToRunechantRed{}, testutils.FakeRedAttack()}
	got := Best(nil, h, nil, gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetIncomingDamage(4).
		Build())
	if got.Value != 5 {
		t.Errorf("Value = %d, want 5", got.Value)
	}
}

// Tests end-to-end that a hand containing a discount attacker is playable iff carryover
// runechants cover the printed cost.
func TestBest_CarryoverFeedsDiscount(t *testing.T) {
	h := []card.Card{cards.AmplifyTheArknightRed{}}
	got := Best(nil, h, nil, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	if got.Value != 0 {
		t.Errorf("no carryover: Value = %d, want 0 (discount insufficient without runechants)", got.Value)
	}
	// With 3 runechants carried in, the discount fully covers the cost. Value is just the
	// Attack() power — consumed carryover runechants aren't re-credited.
	got = Best(nil, h, nil, stateWithRunechants(testutils.Hero{Intel: 4}, 3))
	if got.Value != 6 {
		t.Errorf("carryover=3: Value = %d, want 6 (Attack only; carryover tokens don't re-credit)", got.Value)
	}
	if got.State.RunechantCount() != 0 {
		t.Errorf("carryover=3: leftover Runechants = %d, want 0 (attack consumes tokens)", got.State.RunechantCount())
	}
}
