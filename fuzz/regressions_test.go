package fuzz

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Pinned fuzz seed: a Viserai hand with a modal blocker funded across two pitches must
// produce the same Stats under the eval cache as without it. Without the bestWinner
// pmask tracking the WINNING defenseDealt, the cached and uncached evaluators diverge.
func TestEvalCacheEquivalence_Seed3307073355315735355(t *testing.T) {
	const (
		deckSize  = 40
		maxCopies = 2
		incoming  = 7
		shuffles  = 16
	)
	setupRNG := rand.New(rand.NewSource(3307073355315735355))
	baseline := deck.Random(heroes.Viserai, deckSize, maxCopies, setupRNG, registry.Registry{})
	cached := sim.NewEvaluator().Evaluate(baseline.Copy(), shuffles, sim.Matchup{IncomingPhysicalDamage: incoming}, rand.New(rand.NewSource(99)))
	uncached := sim.NewEvaluatorWithoutCache().Evaluate(baseline.Copy(), shuffles, sim.Matchup{IncomingPhysicalDamage: incoming}, rand.New(rand.NewSource(99)))
	if cached.Hands != uncached.Hands || cached.TotalValue != uncached.TotalValue {
		t.Errorf("eval cache divergence: cached=(hands=%d val=%.0f) uncached=(hands=%d val=%.0f)",
			cached.Hands, cached.TotalValue, uncached.Hands, uncached.TotalValue)
	}
}

// Pinned fuzz seed: a turn with two consecutive aura-banishing DRs must not double-banish
// the same aura when drGraveyard is recomputed between DRs. Surfaces as a card-conservation
// panic from playOneTurn under verifyTurnInvariants.
func TestRunDefense_DRGraveyardCarriesBetweenDRs_Seed6556227072949927836(t *testing.T) {
	const (
		deckSize  = 40
		maxCopies = 2
		incoming  = 7
		shuffles  = 20
	)
	setupRNG := rand.New(rand.NewSource(6556227072949927836))
	baseline := deck.Random(heroes.Viserai, deckSize, maxCopies, setupRNG, registry.Registry{})
	sim.NewEvaluator().Evaluate(baseline.Copy(), shuffles, sim.Matchup{IncomingPhysicalDamage: incoming}, rand.New(rand.NewSource(99)))
	sim.NewEvaluatorWithoutCache().Evaluate(baseline.Copy(), shuffles, sim.Matchup{IncomingPhysicalDamage: incoming}, rand.New(rand.NewSource(99)))
}
