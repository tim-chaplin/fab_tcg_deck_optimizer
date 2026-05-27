package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
)

// Tests that Flock's PlayPrecondition fails when no card in hand has cost ≤ 1.
func TestFlockOfTheFeatherWalkers_PrecondFailsWithNoLowCostHand(t *testing.T) {
	for _, c := range []card.Card{
		cards.FlockOfTheFeatherWalkersRed{}, cards.FlockOfTheFeatherWalkersYellow{}, cards.FlockOfTheFeatherWalkersBlue{},
	} {
		state := gameengine.GameStateBuilder().Build()
		state.SetHand([]card.Card{testutils.FakeRedAttack().WithCost(2)})
		if c.(card.PlayPrecondition).PlayPrecondition(state.Engine(), &card.CardState{Card: c}) {
			t.Errorf("%s: PlayPrecondition = true with only a cost-2 card in hand, want false", c.Name())
		}
	}
}

// Tests that Flock's PlayPrecondition passes when a cost-≤-1 card is in hand.
func TestFlockOfTheFeatherWalkers_PrecondPassesWithLowCostHand(t *testing.T) {
	for _, cost := range []int{0, 1} {
		state := gameengine.GameStateBuilder().Build()
		state.SetHand([]card.Card{testutils.FakeRedAttack().WithCost(cost)})
		if !(cards.FlockOfTheFeatherWalkersRed{}).PlayPrecondition(state.Engine(), &card.CardState{}) {
			t.Errorf("PlayPrecondition = false with a cost-%d card in hand, want true", cost)
		}
	}
}

// Tests that Flock's Play creates exactly one Quicken token.
func TestFlockOfTheFeatherWalkers_PlayCreatesQuicken(t *testing.T) {
	for _, c := range []card.Card{
		cards.FlockOfTheFeatherWalkersRed{}, cards.FlockOfTheFeatherWalkersYellow{}, cards.FlockOfTheFeatherWalkersBlue{},
	} {
		ge := gameengine.New()
		c.Play(ge, ge.Logger(), &card.CardState{Card: c})
		if got := ge.QuickenCount(); got != 1 {
			t.Errorf("%s: QuickenCount after Play = %d, want 1", c.Name(), got)
		}
	}
}

// End-to-end: a pitch source pays Flock's cost while a cost-1 card stays in hand to satisfy
// the reveal, and the post-attack-turn state carries a Quicken token.
func TestFlockOfTheFeatherWalkers_PlaysAndMintsQuickenInAttackTurn(t *testing.T) {
	pitchSource := testutils.FakeRedAttack().
		WithCost(0).
		WithName("PitchSource")
	revealTarget := testutils.FakeRedAttack().
		WithCost(1).
		WithName("RevealTarget")
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.FlockOfTheFeatherWalkersRed{}, pitchSource, revealTarget}
	summary := sim.EvalOneTurnForTesting(d, nil, hand)

	if got := summary.State.QuickenCount(); got != 1 {
		t.Errorf("QuickenCount = %d, want 1 (Flock should mint Quicken on play)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}

// Tests that a printed-cost-2 VariableCost card (Drawn to the Dark Dimension) qualifies as a
// reveal target when its EffectiveCost is discounted to ≤ 1 by Runechants on board. The
// predicate reads EffectiveCost so the gate sees the discount.
func TestFlockOfTheFeatherWalkers_PrecondAcceptsDiscountedVariableCost(t *testing.T) {
	noRune := gameengine.GameStateBuilder().Build()
	noRune.SetHand([]card.Card{cards.DrawnToTheDarkDimensionRed{}})
	if (cards.FlockOfTheFeatherWalkersRed{}).PlayPrecondition(noRune.Engine(), &card.CardState{}) {
		t.Errorf("PlayPrecondition = true with un-discounted Drawn (effective cost 2), want false")
	}

	oneRune := gameengine.GameStateBuilder().
		AddAura(token.NewRunechant(1)).
		Build()
	oneRune.SetHand([]card.Card{cards.DrawnToTheDarkDimensionRed{}})
	if !(cards.FlockOfTheFeatherWalkersRed{}).PlayPrecondition(oneRune.Engine(), &card.CardState{}) {
		t.Errorf("PlayPrecondition = false with discounted Drawn (effective cost 1), want true")
	}
}

// End-to-end: a card can't be both played AND used to satisfy Flock's reveal cost on the
// same turn. Hand: blue pitch + a 0-cost go-again 3-power red attack + Flock R. Pitching
// blue funds Flock's resource cost; the red attack is the only cost-≤-1 reveal target, so
// playing it would empty the hand of eligible reveals and block Flock. Best line: hold
// the attack as the reveal target and play Flock for 5.
func TestFlockOfTheFeatherWalkers_RevealTargetCannotAlsoBePlayed(t *testing.T) {
	bluePitch := testutils.FakeBlueAttack().
		WithName("BluePitch")
	redAttack := testutils.FakeRedAttack().
		WithCost(0).
		WithPower(3).
		WithGoAgain().
		WithName("RedAttack")
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{bluePitch, redAttack, cards.FlockOfTheFeatherWalkersRed{}}
	summary := sim.EvalOneTurnForTesting(d, nil, hand)

	if summary.Value != 5 {
		t.Errorf("Value = %d, want 5 (Flock plays solo for printed 5; the red attack is the reveal target and can't also be played)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
}

// End-to-end: no cost-≤-1 card available besides Flock itself, so the precondition blocks the
// play and no Quicken is minted.
func TestFlockOfTheFeatherWalkers_NoEligibleRevealBlocksPlay(t *testing.T) {
	highCost := testutils.FakeRedAttack().
		WithCost(2).
		WithName("HighCost")
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.FlockOfTheFeatherWalkersRed{}, highCost}
	summary := sim.EvalOneTurnForTesting(d, nil, hand)

	if got := summary.State.QuickenCount(); got != 0 {
		t.Errorf("QuickenCount = %d, want 0 (precondition should block Flock and prevent Quicken)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
}
