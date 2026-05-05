package e2etest

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Pursue to the Edge of Oblivion's on-hit rider marks the opposing hero.
func TestPursueToTheEdgeOfOblivion_MarksOpponentOnHit(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.PursueToTheEdgeOfOblivionRed{}}
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
	if !state.OpponentMarked {
		t.Fatalf("OpponentMarked = false after Pursue resolved, want true")
	}
}

// Tests that with the opponent already marked, Pursue's attack strips the prior mark and
// its on-hit rider reapplies it — net end-of-chain mark stays on.
func TestPursueToTheEdgeOfOblivion_PreservesMarkWhenOpponentAlreadyMarked(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.PursueToTheEdgeOfOblivionRed{}}
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{OpponentMarked: true}, hand)
	if !state.OpponentMarked {
		t.Fatalf("OpponentMarked = false after Pursue resolved against pre-marked opponent, want true")
	}
}

// Tests that Outed against an unmarked opponent gets no marked-defender bonus.
func TestOuted_NoBonusWhenOpponentUnmarked(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.OutedRed{}}
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
	if state.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Outed printed 3{p}, no marked-defender bonus)", state.Value)
	}
}

// Tests that Outed against a pre-marked opponent self-buffs +1{p} for the marked-defender
// rider.
func TestOuted_BonusWhenOpponentAlreadyMarked(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.OutedRed{}}
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{OpponentMarked: true}, hand)
	if state.Value != 4 {
		t.Fatalf("Value = %d, want 4 (Outed printed 3{p} + 1 marked-defender)", state.Value)
	}
}
