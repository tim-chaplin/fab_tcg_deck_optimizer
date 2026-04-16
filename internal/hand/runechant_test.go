package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// TestPlaySequence_DiscountRejectsInsufficientBudget verifies that a DiscountPerRunechant card
// fails its per-play cost check when the chain budget can't cover the effective cost.
func TestPlaySequence_DiscountRejectsInsufficientBudget(t *testing.T) {
	order := []card.Card{runeblade.AmplifyTheArknightRed{}} // PrintedCost 3, Cost() 0
	pcBuf := make([]card.PlayedCard, len(order))
	ptrBuf := make([]*card.PlayedCard, len(order))
	cpBuf := make([]card.Card, 0, len(order))
	state := &card.TurnState{}
	// Chain budget 0, carryover 0 → effective cost = 3 - 0 = 3 > 0, chain illegal.
	dmg, leftover, _, legal := playSequence(hero.Viserai{}, nil, nil, order, pcBuf, ptrBuf, cpBuf, state, 0, 0)
	if legal {
		t.Fatalf("expected illegal chain, got legal (dmg=%d, leftover=%d)", dmg, leftover)
	}
}

// TestPlaySequence_DiscountAffordableWithBudget shows the same card becomes legal once the
// budget covers its printed cost.
func TestPlaySequence_DiscountAffordableWithBudget(t *testing.T) {
	order := []card.Card{runeblade.AmplifyTheArknightRed{}}
	pcBuf := make([]card.PlayedCard, len(order))
	ptrBuf := make([]*card.PlayedCard, len(order))
	cpBuf := make([]card.Card, 0, len(order))
	state := &card.TurnState{}
	// Chain budget 3, carryover 0 → effective cost 3, budget just covers it. Amplify's Attack(6)
	// is the only damage; no runechants to consume.
	dmg, leftover, _, legal := playSequence(hero.Viserai{}, nil, nil, order, pcBuf, ptrBuf, cpBuf, state, 3, 0)
	if !legal {
		t.Fatalf("expected legal chain")
	}
	if dmg != 6 {
		t.Errorf("dmg = %d, want 6", dmg)
	}
	if leftover != 0 {
		t.Errorf("leftover = %d, want 0", leftover)
	}
}

// TestPlaySequence_DiscountUsesCarryoverRunechants shows the discount applies from carryover
// tokens — no chain budget needed when there are enough runechants already in play.
func TestPlaySequence_DiscountUsesCarryoverRunechants(t *testing.T) {
	order := []card.Card{runeblade.AmplifyTheArknightRed{}}
	pcBuf := make([]card.PlayedCard, len(order))
	ptrBuf := make([]*card.PlayedCard, len(order))
	cpBuf := make([]card.Card, 0, len(order))
	state := &card.TurnState{}
	// Chain budget 0, carryover 3 → effective cost 3-3 = 0, legal. Damage is just Amplify's
	// Attack(); the consumed carryover tokens aren't re-credited (they were credited on the
	// previous turn when they were created).
	dmg, leftover, _, legal := playSequence(hero.Viserai{}, nil, nil, order, pcBuf, ptrBuf, cpBuf, state, 0, 3)
	if !legal {
		t.Fatalf("expected legal chain")
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
	pcBuf := make([]card.PlayedCard, len(order))
	ptrBuf := make([]*card.PlayedCard, len(order))
	cpBuf := make([]card.Card, 0, len(order))
	state := &card.TurnState{}
	dmg, leftover, _, legal := playSequence(hero.Viserai{}, nil, nil, order, pcBuf, ptrBuf, cpBuf, state, 0, 0)
	if !legal {
		t.Fatalf("expected legal chain")
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
	got := Best(hero.Viserai{}, nil, h, 0, nil, 0)
	if got.Value != 4 {
		t.Errorf("Value = %d, want 4 (3 Read tokens + 1 Viserai token)", got.Value)
	}
	if got.LeftoverRunechants != 4 {
		t.Errorf("LeftoverRunechants = %d, want 4", got.LeftoverRunechants)
	}
}

// TestBest_MauvrionReadWithCarryover is the same hand with 1 runechant carried in from the
// previous turn. The hand still creates 4 tokens this turn, and the 1 carryover token doesn't
// get consumed (no attack in the chain), so leftover = 5.
func TestBest_MauvrionReadWithCarryover(t *testing.T) {
	h := []card.Card{runeblade.MauvrionSkiesRed{}, runeblade.ReadTheRunesRed{}}
	got := Best(hero.Viserai{}, nil, h, 0, nil, 1)
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
	got := Best(hero.Viserai{}, weapons, h, 0, nil, 1)
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
	got := Best(stubHero{}, nil, h, 0, nil, 0)
	if got.Value != 3 {
		t.Errorf("Value = %d, want 3", got.Value)
	}
	if got.LeftoverRunechants != 3 {
		t.Errorf("LeftoverRunechants = %d, want 3 (0 live + 3 delayed from Blessing)",
			got.LeftoverRunechants)
	}
}

// TestBest_ReduceToRunechantRejectsInsufficientBudget verifies that Reduce to Runechant's
// per-runechant discount is re-priced against the attack line's leftoverRunechants. The hand has
// three Reduces plus an Aether Slash; before the fix, the solver could block with all three
// Reduces at Cost()=0 each while Aether attacked, scoring 13 (5 damage + 8 prevented). With the
// discount correctly modeled, Aether's attack consumes the carryover (none here) and leaves 0
// runechants, so each Reduce really costs 1 — the 3-Reduce full-block partition is rejected.
// Best feasible partition: 2 Reduce block + 1 Reduce pitch + Aether pitch = prevent 8 + 2 damage
// credits from the two Reduce reactions creating runechants = 10.
func TestBest_ReduceToRunechantRejectsInsufficientBudget(t *testing.T) {
	h := []card.Card{
		runeblade.ReduceToRunechantRed{},
		runeblade.ReduceToRunechantRed{},
		runeblade.ReduceToRunechantRed{},
		runeblade.AetherSlashRed{},
	}
	got := Best(hero.Viserai{}, nil, h, 12, nil, 0)
	if got.Value != 10 {
		t.Errorf("Value = %d, want 10 (3-Reduce block rejected without runechants; 2-Reduce block with each creating a runechant reaches 10). roles=[%s]",
			got.Value, FormatRoles(h, got.Roles))
	}
}

// TestBest_ReduceToRunechantDiscountedByCarryover pins the other side: when the previous turn
// left 3 runechants behind and no attack this turn consumes them, Reduce's effective cost is
// max(1-3, 0) = 0, and a full three-Reduce block becomes affordable again. Optimal: Aether pitch,
// all three Reduces block → prevent 12 + 3 damage credits (one per Reduce's Runechant creation).
func TestBest_ReduceToRunechantDiscountedByCarryover(t *testing.T) {
	h := []card.Card{
		runeblade.ReduceToRunechantRed{},
		runeblade.ReduceToRunechantRed{},
		runeblade.ReduceToRunechantRed{},
		runeblade.AetherSlashRed{},
	}
	got := Best(hero.Viserai{}, nil, h, 12, nil, 3)
	if got.Value != 15 {
		t.Errorf("Value = %d, want 15 (12 prevented + 3 Runechant-creation credits)", got.Value)
	}
}

// TestBest_CarryoverFeedsDiscount verifies end-to-end: a hand containing a discount attacker is
// playable when the previous turn left enough runechants behind.
func TestBest_CarryoverFeedsDiscount(t *testing.T) {
	// Single Amplify the Arknight (Red): Cost()=0, PrintedCost=3, Attack()=6. With no pitch,
	// chain budget is 0. Without any runechants, effective cost 3 exceeds the budget — so
	// attacking is illegal and Value should be 0.
	h := []card.Card{runeblade.AmplifyTheArknightRed{}}
	got := Best(stubHero{}, nil, h, 0, nil, 0)
	if got.Value != 0 {
		t.Errorf("no carryover: Value = %d, want 0 (discount insufficient without runechants)", got.Value)
	}
	// With 3 runechants carried in, the discount fully covers the cost. Value is just the
	// Attack() power — consumed carryover runechants aren't re-credited.
	got = Best(stubHero{}, nil, h, 0, nil, 3)
	if got.Value != 6 {
		t.Errorf("carryover=3: Value = %d, want 6 (Attack only; carryover tokens don't re-credit)", got.Value)
	}
	if got.LeftoverRunechants != 0 {
		t.Errorf("carryover=3: LeftoverRunechants = %d, want 0 (attack consumes tokens)", got.LeftoverRunechants)
	}
}

