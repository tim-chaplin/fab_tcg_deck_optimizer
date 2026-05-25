package sim

// Reads back the gameengine.Pool's HighWaterMark after a representative workload — used
// to size a future fixed-capacity Pool. Output-only; no assertions.

import (
	"math/rand"
	"testing"
)

// Runs viserai_v4 through Evaluate and logs the per-worker pool's peak in-flight count.
// Skipped when the saved deck is absent. With FreeAll-at-end-of-shuffle wired into
// runOneShuffle, HighWaterMark reports the max Gets-minus-Puts seen in any single
// shuffle — the bound a fixed-capacity Pool would need.
func TestPool_HighWaterAfterRealDeckEval(t *testing.T) {
	const (
		incoming = 7
		shuffles = 200
	)
	loaded := loadRealDeck(t)
	if loaded == nil {
		t.Skip("mydecks/viserai_v4.json not found — need the saved deck to measure the realistic pool peak")
	}
	ev := NewEvaluator()
	rng := rand.New(rand.NewSource(42))
	ev.Evaluate(loaded.Copy(), shuffles, Matchup{IncomingDamage: incoming}, rng)
	if ev.cachedBufs == nil {
		t.Fatalf("Evaluator never built its attackBufs scratch — no pool to inspect")
	}
	t.Logf("statePool high-water mark across %d shuffles = %d", shuffles, ev.cachedBufs.statePool.HighWaterMark())
}
