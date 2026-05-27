package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Tests that EBT's PlayPrecondition fails when no damage has been dealt this turn.
func TestEvenBiggerThanThat_PlayPreconditionFailsWithoutPriorDamage(t *testing.T) {
	for _, c := range []card.Card{
		cards.EvenBiggerThanThatRed{}, cards.EvenBiggerThanThatYellow{}, cards.EvenBiggerThanThatBlue{},
	} {
		ge := gameengine.New()
		gate := c.(card.PlayPrecondition).PlayPrecondition(ge, &card.CardState{Card: c})
		if gate {
			t.Errorf("%s: PlayPrecondition = true with DamageDealt=0, want false", c.Name())
		}
	}
}

// Tests that a landing power-4 attack credits 1 unblocked damage (blocks-in-multiples-
// of-3 model: opponent blocks 3, leaks 1) and flips HitThisTurn.
func TestEvenBiggerThanThat_LandingAttackCreditsUnblockedDamage(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{testutils.FakeRedAttack().WithPower(4)}
	summary := sim.EvalOneTurnForTesting(d, nil, hand)
	if got := summary.State.DamageDealt(); got != 1 {
		t.Errorf("DamageDealt = %d, want 1 (power-4 attack, opponent blocks 3)\nBestLine: %s",
			got, formatBestLine(summary.BestLine))
	}
	if !summary.State.HitThisTurn() {
		t.Errorf("HitThisTurn = false, want true (physical hit landed)")
	}
}

// Tests EBT's PlayPrecondition gate flips once a physical attack has landed this turn.
func TestEvenBiggerThanThat_PlayPreconditionPassesAfterHit(t *testing.T) {
	ge := gameengine.New()
	ge.SetHitThisTurn(true)
	if !(cards.EvenBiggerThanThatRed{}).PlayPrecondition(ge, &card.CardState{}) {
		t.Errorf("PlayPrecondition = false after SetHitThisTurn(true), want true")
	}
}

// Tests that DamageDealt > 0 alone (arcane only, no physical hit) does not satisfy the
// EBT precondition — HitThisTurn must be set by a physical hit.
func TestEvenBiggerThanThat_PlayPreconditionFailsWithArcaneOnly(t *testing.T) {
	ge := gameengine.New()
	ge.AddDamageDealt(1) // simulates arcane-only credit
	if (cards.EvenBiggerThanThatRed{}).PlayPrecondition(ge, &card.CardState{}) {
		t.Errorf("PlayPrecondition = true with arcane-only damage, want false (arcane doesn't 'hit')")
	}
}

// Tests EBT's hit branch: deck top has power > DamageDealt → Quicken created, card drawn.
func TestEvenBiggerThanThat_HighDeckTopCreatesQuickenAndDraws(t *testing.T) {
	highPower := testutils.FakeRedAttack().
		WithPower(4).
		WithName("HighPower")
	lowPower := testutils.FakeRedAttack().
		WithPower(2).
		WithName("LowPower")
	d := deck.New(testutils.Hero{Intel: 4}, nil, []deck.Card{highPower, lowPower, lowPower})
	state := gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetDeck(d).
		Build()
	state.AddDamageDealt(3)
	ge := state.Engine()
	handBefore := ge.HandSize()
	(cards.EvenBiggerThanThatRed{}).Play(ge, ge.Logger(), &card.CardState{Card: cards.EvenBiggerThanThatRed{}})

	if ge.QuickenCount() != 1 {
		t.Errorf("QuickenCount = %d, want 1 (top power 4 > dealt 3 should mint Quicken)", ge.QuickenCount())
	}
	if ge.HandSize() != handBefore+1 {
		t.Errorf("HandSize = %d, want %d (top-rider should draw a card)", ge.HandSize(), handBefore+1)
	}
}

// Tests EBT's miss branch: deck top has power <= DamageDealt → no Quicken, no draw.
func TestEvenBiggerThanThat_LowDeckTopFizzles(t *testing.T) {
	lowPower := testutils.FakeRedAttack().
		WithPower(2).
		WithName("LowPower")
	d := deck.New(testutils.Hero{Intel: 4}, nil, []deck.Card{lowPower, lowPower, lowPower})
	state := gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetDeck(d).
		Build()
	state.AddDamageDealt(5)
	ge := state.Engine()
	handBefore := ge.HandSize()
	(cards.EvenBiggerThanThatRed{}).Play(ge, ge.Logger(), &card.CardState{Card: cards.EvenBiggerThanThatRed{}})

	if ge.QuickenCount() != 0 {
		t.Errorf("QuickenCount = %d, want 0 (top power 2 <= dealt 5 should fizzle)", ge.QuickenCount())
	}
	if ge.HandSize() != handBefore {
		t.Errorf("HandSize = %d, want %d (fizzle should not draw)", ge.HandSize(), handBefore)
	}
}

// Tests that a Quicken token grants Go again to the firing attack's CardState and
// consumes one charge per fire.
func TestQuicken_GrantsGoAgainAndConsumesPerFire(t *testing.T) {
	ge := gameengine.New()
	ge.CreateQuicken(2)
	if got := ge.QuickenCount(); got != 2 {
		t.Fatalf("QuickenCount after CreateQuicken(2) = %d, want 2", got)
	}

	first := &card.CardState{Card: testutils.FakeRedAttack()}
	ge.FireTriggers(triggertype.CardOrAbility, first)
	if !first.GrantedGoAgain {
		t.Errorf("first attack GrantedGoAgain = false, want true (Quicken should grant)")
	}
	if got := ge.QuickenCount(); got != 1 {
		t.Errorf("QuickenCount after first fire = %d, want 1 (one charge consumed)", got)
	}

	second := &card.CardState{Card: testutils.FakeRedAttack()}
	ge.FireTriggers(triggertype.CardOrAbility, second)
	if !second.GrantedGoAgain {
		t.Errorf("second attack GrantedGoAgain = false, want true")
	}
	if got := ge.QuickenCount(); got != 0 {
		t.Errorf("QuickenCount after second fire = %d, want 0 (last charge consumed)", got)
	}

	third := &card.CardState{Card: testutils.FakeRedAttack()}
	ge.FireTriggers(triggertype.CardOrAbility, third)
	if third.GrantedGoAgain {
		t.Errorf("third attack GrantedGoAgain = true, want false (Quicken exhausted)")
	}
}

// Tests that Quicken's IsAttack filter skips non-attack CardOrAbility fires.
func TestQuicken_SkipsNonAttackFires(t *testing.T) {
	ge := gameengine.New()
	ge.CreateQuicken(1)
	action := &card.CardState{Card: testutils.FakeRedAction()}
	ge.FireTriggers(triggertype.CardOrAbility, action)
	if action.GrantedGoAgain {
		t.Errorf("non-attack action GrantedGoAgain = true, want false")
	}
	if got := ge.QuickenCount(); got != 1 {
		t.Errorf("QuickenCount = %d, want 1 (filter should preserve the charge)", got)
	}
}
