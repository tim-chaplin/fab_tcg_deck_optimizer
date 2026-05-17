package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestWarmongersRecital_NoAttackReturnsZero: no qualifying next attack card → +N rider fizzles.
func TestWarmongersRecital_NoAttackReturnsZero(t *testing.T) {
	ge := gameengine.New()
	for _, c := range []card.Card{
		cards.WarmongersRecitalRed{}, cards.WarmongersRecitalYellow{}, cards.WarmongersRecitalBlue{},
	} {
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: c})
		if ge.Value() != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), ge.Value())
		}
	}
}

// TestWarmongersRecital_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestWarmongersRecital_NonAttackInRemainingFizzles(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.GenericAction()}}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.WarmongersRecitalRed{}})
	if ge.Value() != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", ge.Value())
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
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
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
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{deckTop}).Build()}
	ge.SetCardsRemaining([]*card.CardState{targetState})
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.WarmongersRecitalRed{}})
	if len(targetState.OnHit) != 1 {
		t.Fatalf("OnHit not registered: len=%d", len(targetState.OnHit))
	}

	// Simulate the chain step'ge graveyard deposit, then fire the OnHit.
	ge.AddToGraveyard(target)
	preLog := len(ge.LogEntries())
	h := &targetState.OnHit[0]
	h.Fire(ge, ge.Logger(), targetState, h)

	if g := ge.Graveyard(); len(g) != 0 {
		t.Errorf("Graveyard after recycle = %v, want empty (target pulled out)", g)
	}
	if got := ge.Deck().Size(); got != 2 {
		t.Errorf("Deck size after recycle = %d, want 2 (target appended to bottom)", got)
	}
	if top := ge.Deck().PeekTop(); top != card.Card(deckTop) {
		t.Errorf("Deck top after recycle = %v, want %v (target went to bottom, deckTop unchanged)", top, deckTop)
	}
	// Rider line attributes the recycle to the buffed attack, not Warmonger'ge Recital.
	added := ge.LogEntries()[preLog:]
	if len(added) != 1 || added[0].Source != target.DisplayName() {
		t.Errorf("rider log = %+v, want one entry sourced to %q", added, target.DisplayName())
	}
}
