package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestWarmongersRecital_NoAttackReturnsZero: no qualifying next attack card → +N rider fizzles.
func TestWarmongersRecital_NoAttackReturnsZero(t *testing.T) {
	s := &gameengine.GameEngine{GameState: gameengine.NewState()}
	for _, c := range []card.Card{
		cards.WarmongersRecitalRed{}, cards.WarmongersRecitalYellow{}, cards.WarmongersRecitalBlue{},
	} {
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: c})
		if s.Value() != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), s.Value())
		}
	}
}

// TestWarmongersRecital_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestWarmongersRecital_NonAttackInRemainingFizzles(t *testing.T) {
	s := gameengine.NewFromSpec(gameengine.Spec{CardsRemaining: []*card.CardState{{Card: testutils.GenericAction()}}})
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.WarmongersRecitalRed{}})
	if s.Value() != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", s.Value())
	}
}

// TestWarmongersRecital_NextAttackReceivesBonusAndOnHit: first attack-action target gets the
// per-variant +N{p} bonus AND an OnHit handler appended for the recycle rider.
func TestWarmongersRecital_NextAttackReceivesBonusAndOnHit(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.WarmongersRecitalRed{}, 3},
		{cards.WarmongersRecitalYellow{}, 2},
		{cards.WarmongersRecitalBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.GenericAttack(0, 0)}
		s := gameengine.NewFromSpec(gameengine.Spec{CardsRemaining: []*card.CardState{target}})
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: tc.c})
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
	targetState := &card.CardState{Card: target}
	deckTop := testutils.GenericAttack(1, 7)
	s := gameengine.NewFromCards([]card.Card{deckTop}, nil)
	s.SetCardsRemaining([]*card.CardState{targetState})
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.WarmongersRecitalRed{}})
	if len(targetState.OnHit) != 1 {
		t.Fatalf("OnHit not registered: len=%d", len(targetState.OnHit))
	}

	// Simulate the chain step's graveyard deposit, then fire the OnHit.
	s.AddToGraveyard(target)
	preLog := len(s.LogEntries())
	h := &targetState.OnHit[0]
	h.Fire(s, s.Logger(), targetState, h)

	if g := s.Graveyard(); len(g) != 0 {
		t.Errorf("Graveyard after recycle = %v, want empty (target pulled out)", g)
	}
	if got := s.Deck().Size(); got != 2 {
		t.Errorf("Deck size after recycle = %d, want 2 (target appended to bottom)", got)
	}
	if top := s.Deck().PeekTop(); top != card.Card(deckTop) {
		t.Errorf("Deck top after recycle = %v, want %v (target went to bottom, deckTop unchanged)", top, deckTop)
	}
	// Rider line attributes the recycle to the buffed attack, not Warmonger's Recital.
	added := s.LogEntries()[preLog:]
	if len(added) != 1 || added[0].Source != target.DisplayName() {
		t.Errorf("rider log = %+v, want one entry sourced to %q", added, target.DisplayName())
	}
}
