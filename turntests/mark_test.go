package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Pursue to the Edge of Oblivion's on-hit rider marks the opposing hero.
func TestPursueToTheEdgeOfOblivion_MarksOpponentOnHit(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.PursueToTheEdgeOfOblivionRed{}}
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand)
	if !state.OpponentMarked {
		t.Fatalf("OpponentMarked = false after Pursue resolved, want true")
	}
}

// Tests that with the opponent already marked, Pursue's attack strips the prior mark and
// its on-hit rider reapplies it — net end-of-chain mark stays on.
func TestPursueToTheEdgeOfOblivion_PreservesMarkWhenOpponentAlreadyMarked(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.PursueToTheEdgeOfOblivionRed{}}
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, gameengine.GameStateBuilder().SetOpponentMarked(true).Build(), hand)
	if !state.OpponentMarked {
		t.Fatalf("OpponentMarked = false after Pursue resolved against pre-marked opponent, want true")
	}
}

// Tests that Outed against an unmarked opponent gets no marked-defender bonus.
func TestOuted_NoBonusWhenOpponentUnmarked(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.OutedRed{}}
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand)
	if state.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Outed printed 3{p}, no marked-defender bonus)", state.Value)
	}
}

// Tests that Outed against a pre-marked opponent self-buffs +1{p} for the marked-defender
// rider.
func TestOuted_BonusWhenOpponentAlreadyMarked(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.OutedRed{}}
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, gameengine.GameStateBuilder().SetOpponentMarked(true).Build(), hand)
	if state.Value != 4 {
		t.Fatalf("Value = %d, want 4 (Outed printed 3{p} + 1 marked-defender)", state.Value)
	}
}
