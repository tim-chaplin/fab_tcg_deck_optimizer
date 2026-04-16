package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
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
	dmg, leftover, legal := playSequence(hero.Viserai{}, nil, nil, order, pcBuf, ptrBuf, cpBuf, state, 0, 0)
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
	dmg, leftover, legal := playSequence(hero.Viserai{}, nil, nil, order, pcBuf, ptrBuf, cpBuf, state, 3, 0)
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
// tokens — no budget needed when there are enough runechants already in play.
func TestPlaySequence_DiscountUsesCarryoverRunechants(t *testing.T) {
	order := []card.Card{runeblade.AmplifyTheArknightRed{}}
	pcBuf := make([]card.PlayedCard, len(order))
	ptrBuf := make([]*card.PlayedCard, len(order))
	cpBuf := make([]card.Card, 0, len(order))
	state := &card.TurnState{}
	// Chain budget 0, carryover 3 → effective cost 3-3 = 0, legal. Amplify's attack consumes
	// the 3 runechants for +3 damage, on top of its Attack() of 6.
	dmg, leftover, legal := playSequence(hero.Viserai{}, nil, nil, order, pcBuf, ptrBuf, cpBuf, state, 0, 3)
	if !legal {
		t.Fatalf("expected legal chain")
	}
	if dmg != 9 {
		t.Errorf("dmg = %d, want 9 (6 power + 3 consumed runechants)", dmg)
	}
	if leftover != 0 {
		t.Errorf("leftover = %d, want 0 (attack consumes all runechants)", leftover)
	}
}

// TestPlaySequence_LeftoverFromNonAttackAction confirms that runechants created by a non-attack
// action with no following attack persist as leftover.
func TestPlaySequence_LeftoverFromNonAttackAction(t *testing.T) {
	order := []card.Card{runeblade.ReadTheRunesRed{}} // creates 3 runechants, not an attack
	pcBuf := make([]card.PlayedCard, len(order))
	ptrBuf := make([]*card.PlayedCard, len(order))
	cpBuf := make([]card.Card, 0, len(order))
	state := &card.TurnState{}
	dmg, leftover, legal := playSequence(hero.Viserai{}, nil, nil, order, pcBuf, ptrBuf, cpBuf, state, 0, 0)
	if !legal {
		t.Fatalf("expected legal chain")
	}
	if dmg != 0 {
		t.Errorf("dmg = %d, want 0 (Read the Runes isn't an attack)", dmg)
	}
	if leftover != 3 {
		t.Errorf("leftover = %d, want 3", leftover)
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
	// With 3 runechants carried in, the discount fully covers the cost. Amplify attacks for
	// Attack()(6) + consumed tokens(3) = 9.
	got = Best(stubHero{}, nil, h, 0, nil, 3)
	if got.Value != 9 {
		t.Errorf("carryover=3: Value = %d, want 9", got.Value)
	}
	if got.LeftoverRunechants != 0 {
		t.Errorf("carryover=3: LeftoverRunechants = %d, want 0 (attack consumes tokens)", got.LeftoverRunechants)
	}
}

// TestBest_LeftoverIntoNextTurn demonstrates that a hand whose optimal play leaves runechants in
// play (e.g. non-attack-action chain) returns them via Play.LeftoverRunechants for the next turn.
func TestBest_LeftoverIntoNextTurn(t *testing.T) {
	// Read the Runes (Red) alone: Pitch 1, creates 3 Runechants, not an attack. Best play is
	// ATTACK (play the card); no damage is dealt but 3 tokens carry into next turn.
	h := []card.Card{runeblade.ReadTheRunesRed{}}
	got := Best(stubHero{}, nil, h, 0, nil, 0)
	if got.LeftoverRunechants != 3 {
		t.Errorf("LeftoverRunechants = %d, want 3", got.LeftoverRunechants)
	}
}

