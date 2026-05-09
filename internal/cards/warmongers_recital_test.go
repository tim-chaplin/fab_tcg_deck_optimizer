package cards_test

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestWarmongersRecital_NoAttackReturnsZero: no qualifying next attack card → +N rider fizzles.
func TestWarmongersRecital_NoAttackReturnsZero(t *testing.T) {
	s := sim.TurnState{}
	for _, c := range []sim.Card{
		cards.WarmongersRecitalRed{}, cards.WarmongersRecitalYellow{}, cards.WarmongersRecitalBlue{},
	} {
		c.Play(&s, &sim.CardState{Card: c})
		if s.Value != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), s.Value)
		}
	}
}

// TestWarmongersRecital_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestWarmongersRecital_NonAttackInRemainingFizzles(t *testing.T) {
	s := sim.TurnState{CardsRemaining: []*sim.CardState{{Card: testutils.GenericAction()}}}
	(cards.WarmongersRecitalRed{}).Play(&s, &sim.CardState{Card: cards.WarmongersRecitalRed{}})
	if s.Value != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", s.Value)
	}
}

// TestWarmongersRecital_NextAttackReceivesBonusAndOnHit: first attack-action target gets the
// per-variant +N{p} bonus AND an OnHit handler appended for the recycle rider.
func TestWarmongersRecital_NextAttackReceivesBonusAndOnHit(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{cards.WarmongersRecitalRed{}, 3},
		{cards.WarmongersRecitalYellow{}, 2},
		{cards.WarmongersRecitalBlue{}, 1},
	}
	for _, tc := range cases {
		target := &sim.CardState{Card: testutils.GenericAttack(0, 0)}
		s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
		if len(target.OnHit) != 1 {
			t.Errorf("%s: target OnHit len = %d, want 1 (recycle handler)", tc.c.Name(), len(target.OnHit))
		}
	}
}

// TestWarmongersRecital_OnHitFireRecyclesTargetFromGraveyardToDeckBottom: firing the OnHit
// handler pulls target from graveyard and appends it to the bottom of the deck.
func TestWarmongersRecital_OnHitFireRecyclesTargetFromGraveyardToDeckBottom(t *testing.T) {
	target := testutils.GenericAttack(0, 5)
	targetState := &sim.CardState{Card: target}
	deckTop := testutils.GenericAttack(1, 7)
	s := sim.NewTurnState([]sim.Card{deckTop}, nil)
	s.CardsRemaining = []*sim.CardState{targetState}
	(cards.WarmongersRecitalRed{}).Play(s, &sim.CardState{Card: cards.WarmongersRecitalRed{}})
	if len(targetState.OnHit) != 1 {
		t.Fatalf("OnHit not registered: len=%d", len(targetState.OnHit))
	}

	// Simulate the chain step's graveyard deposit, then fire the OnHit.
	s.AddToGraveyard(target)
	preLog := len(s.LogEntries())
	h := &targetState.OnHit[0]
	h.Fire(s, targetState, h)

	if g := s.Graveyard(); len(g) != 0 {
		t.Errorf("Graveyard after recycle = %v, want empty (target pulled out)", g)
	}
	if d := s.Deck(); len(d) != 2 || d[0] != deckTop || d[1] != target {
		t.Errorf("Deck after recycle = %v, want [deckTop, target] (target appended to bottom)", d)
	}
	// Rider line attributes the recycle to the buffed attack, not Warmonger's Recital.
	added := s.LogEntries()[preLog:]
	if len(added) != 1 || added[0].Source != target.DisplayName() {
		t.Errorf("rider log = %+v, want one entry sourced to %q", added, target.DisplayName())
	}
}
