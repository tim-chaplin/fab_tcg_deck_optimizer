package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// featherfootDeck is five filler attacks — enough to deal a full hand next turn. This
// turn draws no cards, so deck contents don't affect the line under test.
func featherfootDeck() *deck.Deck {
	filler := testutils.FakeRedAttack().
		WithCost(1).
		WithPower(3)
	return deck.New(heroes.Viserai, nil, []deck.Card{filler, filler, filler, filler, filler})
}

// Tests that Talisman of Featherfoot, in play at start of turn, fires when a +1{p}
// attack reaction buffs the player's attack: the talisman is destroyed and the buffed
// attack gains go again, banking the action point that funds a follow-up attack from
// hand.
func TestTalismanOfFeatherfoot_PlusOneBuffGrantsGoAgainAndSelfDestructs(t *testing.T) {
	// Arcanic Crackle Red: 0-cost 3-power attack action + 1 arcane.
	// Lunging Press Blue: 0-cost +1{p} AR.
	// FakeRedAttack power 1 cost 0: the follow-up that consumes the granted AP.
	hand := []card.Card{
		cards.ArcanicCrackleRed{},
		cards.LungingPressBlue{},
		testutils.FakeRedAttack().
			WithCost(0).
			WithPower(1),
	}
	state := gameengine.GameStateBuilder().
		SetIncomingPhysicalDamage(0).
		CreateItemFromCard(cards.TalismanOfFeatherfootYellow{}).
		Build()

	summary := sim.EvalOneTurnForTesting(featherfootDeck(), state, hand)

	// Expect Crackle 3 + LP +1 buff + Crackle arcane 1 + follow-up 1-power = 6.
	if summary.Value != 6 {
		t.Fatalf("Value = %d, want 6 (Crackle 3 + LP +1 + arcane 1 + Featherfoot-funded follow-up 1)", summary.Value)
	}
	if n := len(summary.State.Items()); n != 0 {
		t.Fatalf("Items() = %d after turn, want 0 (Featherfoot must self-destruct on +1 fire)", n)
	}
}

// Control: without Talisman of Featherfoot in play, the same hand has no AP for the
// follow-up after Crackle + LP, so the 1-power filler stays in hand.
func TestTalismanOfFeatherfoot_PlusOneBuffWithoutTalismanNoFollowUp(t *testing.T) {
	hand := []card.Card{
		cards.ArcanicCrackleRed{},
		cards.LungingPressBlue{},
		testutils.FakeRedAttack().
			WithCost(0).
			WithPower(1),
	}

	summary := sim.EvalOneTurnForTesting(featherfootDeck(), gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build(), hand)

	// Expect Crackle 3 + LP +1 buff + arcane 1 = 5 (no follow-up: only 1 AP available).
	if summary.Value != 5 {
		t.Fatalf("Value = %d, want 5 (Crackle 3 + LP +1 + arcane 1; no AP for follow-up)", summary.Value)
	}
}

// Tests that Featherfoot grants go-again to a no-go-again attack via the deferred
// go-again check in finalizeActiveAttack: hand is two no-go-again attacks + a +1
// attack reaction. The AR buffs the first attack, Featherfoot fires and flips
// GrantedGoAgain, the deferred check (after all ARs targeting that attack resolve)
// banks the AP, and the second attack plays. Whole hand played.
func TestTalismanOfFeatherfoot_GoAgainGrantPaysForSecondAttack(t *testing.T) {
	attacker1 := testutils.FakeRedAttack().
		WithCost(0).
		WithPower(3).
		WithName("Attacker1")
	attacker2 := testutils.FakeRedAttack().
		WithCost(0).
		WithPower(2).
		WithName("Attacker2")
	hand := []card.Card{attacker1, cards.LungingPressBlue{}, attacker2}
	state := gameengine.GameStateBuilder().
		SetIncomingPhysicalDamage(0).
		CreateItemFromCard(cards.TalismanOfFeatherfootYellow{}).
		Build()

	summary := sim.EvalOneTurnForTesting(featherfootDeck(), state, hand)

	// Whole hand resolves: Attacker1 3 + LP +1 buff = 4, then Attacker2 2 = 6 total.
	if summary.Value != 6 {
		t.Fatalf("Value = %d, want 6 (Attacker1 3 + LP +1 + Attacker2 2; Featherfoot's go-again grant funds the second attack)", summary.Value)
	}
	if n := len(summary.State.Items()); n != 0 {
		t.Fatalf("Items() = %d after turn, want 0 (Featherfoot must self-destruct on +1 fire)", n)
	}
}

// Counterfactual: same hand without Talisman of Featherfoot in play. The first attack
// has no go-again, the AR is free, but the second attack has no AP — the optimiser can
// only play one attack + AR (or drop the AR for the second attack).
func TestTalismanOfFeatherfoot_NoTalismanOnlyOneAttackPlays(t *testing.T) {
	attacker1 := testutils.FakeRedAttack().
		WithCost(0).
		WithPower(3).
		WithName("Attacker1")
	attacker2 := testutils.FakeRedAttack().
		WithCost(0).
		WithPower(2).
		WithName("Attacker2")
	hand := []card.Card{attacker1, cards.LungingPressBlue{}, attacker2}

	summary := sim.EvalOneTurnForTesting(featherfootDeck(), gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build(), hand)

	// Best line: Attacker1 3 + LP +1 buff = 4. Attacker2 can't be paid for; AP exhausted.
	if summary.Value != 4 {
		t.Fatalf("Value = %d, want 4 (Attacker1 3 + LP +1; no Featherfoot means no AP for the second attack)", summary.Value)
	}
}

// Tests that a +2{p} buff does NOT fire Featherfoot. Razor Reflex Yellow mode 1 grants
// +2{p} to a cost-≤1 attack action card — exactly the "exactly +1" gate's negative
// case. The talisman stays in play, no extra AP is banked, and the 1-power follow-up
// stays in hand.
func TestTalismanOfFeatherfoot_PlusTwoBuffDoesNotFire(t *testing.T) {
	// Arcanic Crackle Red pitches for 1 to fund Razor Reflex Yellow's 1-resource cost.
	// Second Crackle is the AR target; Razor Reflex Yellow buffs +2{p} in mode 1.
	hand := []card.Card{
		cards.ArcanicCrackleRed{},
		cards.RazorReflexYellow{},
		cards.ArcanicCrackleRed{},
		testutils.FakeRedAttack().
			WithCost(0).
			WithPower(1),
	}
	state := gameengine.GameStateBuilder().
		SetIncomingPhysicalDamage(0).
		CreateItemFromCard(cards.TalismanOfFeatherfootYellow{}).
		Build()

	summary := sim.EvalOneTurnForTesting(featherfootDeck(), state, hand)

	if n := len(summary.State.Items()); n != 1 {
		t.Fatalf("Items() = %d after turn, want 1 (Featherfoot must NOT fire on +2 buff)", n)
	}
}
