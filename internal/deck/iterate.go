package deck

// Parallel iterate-mode round runner: IterateParallel fans out candidate Mutations across
// workers, two-phase evaluates each (shallow screen then deep confirm), and applies the
// Metropolis acceptance rule so simulated-annealing temperature widens the acceptable set.

import (
	"context"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/klauspost/cpuid/v2"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
)

// IterateParallel runs one iterate-mode round. Workers share a queue; each goroutine does
// the shallow screen and, if shallow clears the effective threshold, the deep-shuffles
// confirmation for the same mutation. The first worker to land an acceptable mutation wins
// (cancellation stops the others). Parallelising deep confirms keeps rounds with noisy
// shallow screens bounded by max(shallow wall, deeps/workers × deep wall).
//
// Annealing: at temperature == 0 only strict improvements clearing the minImprovement
// margin are accepted (classical hill climb with a noise floor). At temperature > 0 worse
// mutations are also accepted with probability exp((deepAvg - baseline) / temperature) — a
// Metropolis-style SA gate that bypasses the minImprovement margin entirely (so the SA walk
// retains its escape-local-maxima behaviour even when the floor is non-zero). The shallow
// pre-screen widens proportionally: at T==0 it cuts at baseline + minImprovement; at T>0 it
// cuts at baseline - 3·T (covers ~95% of probabilistic acceptances; exp(-3) ≈ 0.05).
//
// minImprovement is the noise floor on strict improvements: a mutation must lift deepAvg by
// more than this amount above bestAvg to be accepted at T==0. Prevents infinite loops where
// repeated near-zero "wins" (within shuffle noise) keep accepting indefinitely. Pass 0 for
// the original strict-greater behaviour.
//
// Mutations are pulled FIFO so the earliest-position-wins heuristic of serial iterate
// generally holds, but a worker locked on a deep confirm at position 20 doesn't block
// position 25 — a later-position mutation can occasionally win if its deep confirm
// finishes first.
//
// bestAvg is the current deck's avg (SA "current state", not the all-time best). seed is
// a base; worker w uses (seed + w) for shallow and a derived stream for deep + acceptance
// rolls. shallowCompleted / deepsCompleted are optional live-progress counters.
//
// Returns (acceptedDeck, acceptedAvg, acceptedIndex, true) on first acceptance, or
// (nil, bestAvg, -1, false) if nothing cleared the gate or ctx was cancelled.
func IterateParallel(
	ctx context.Context,
	mutations []Mutation,
	bestAvg float64,
	temperature float64,
	minImprovement float64,
	shallowShuffles, deepShuffles, incoming, numWorkers int,
	seed int64,
	shallowCompleted *atomic.Int64,
	deepsCompleted *atomic.Int64,
) (*Deck, float64, int, bool) {
	if numWorkers <= 0 {
		numWorkers = defaultWorkers()
	}
	if len(mutations) == 0 {
		return nil, bestAvg, -1, false
	}

	// Shallow threshold mirrors the deep acceptance gate's reach: at T==0 require a clear of
	// the noise floor (bestAvg + minImprovement); at T>0 widen to bestAvg - 3·T so
	// probabilistically-acceptable mutations still clear the pre-screen. The probabilistic
	// gate ignores minImprovement so the shallow gate ignores it too at T>0.
	shallowThreshold := bestAvg + minImprovement
	if temperature > 0 {
		shallowThreshold = bestAvg - 3*temperature
	}

	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	type improvement struct {
		idx  int
		avg  float64
		deck *Deck
	}
	improvementCh := make(chan improvement, numWorkers)

	jobs := make(chan int, len(mutations))
	for i := range mutations {
		jobs <- i
	}
	close(jobs)

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerIdx int) {
			defer wg.Done()
			ev := hand.NewEvaluator()
			shallowRng := rand.New(rand.NewSource(seed + int64(workerIdx)))
			// Derive an independent deep stream so the two phases don't share rng state. The
			// acceptance-roll rng shares the deep stream — the deep eval has already happened
			// by the time the roll runs, so no cross-influence on the deep result.
			deepRng := rand.New(rand.NewSource(seed ^ (int64(workerIdx)+1)*int64(0x9e3779b9)))
			for i := range jobs {
				if innerCtx.Err() != nil {
					return
				}
				mut := mutations[i]
				d := New(mut.Deck.Hero, mut.Deck.Weapons, mut.Deck.Cards)
				shallowAvg := d.EvaluateWith(shallowShuffles, incoming, shallowRng, ev).Mean()
				if shallowCompleted != nil {
					shallowCompleted.Add(1)
				}
				if shallowAvg <= shallowThreshold {
					continue
				}
				if innerCtx.Err() != nil {
					return
				}
				// Fresh Deck for the deep pass so d.Stats from the shallow run doesn't leak in.
				dd := New(mut.Deck.Hero, mut.Deck.Weapons, mut.Deck.Cards)
				deepAvg := dd.EvaluateWith(deepShuffles, incoming, deepRng, ev).Mean()
				if deepsCompleted != nil {
					deepsCompleted.Add(1)
				}
				if !acceptMutation(deepAvg, bestAvg, temperature, minImprovement, deepRng) {
					continue
				}
				select {
				case improvementCh <- improvement{idx: i, avg: deepAvg, deck: dd}:
				default:
					// Another worker already filled the buffer; drop silently.
				}
				cancel()
				return
			}
		}(w)
	}

	workersDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(workersDone)
	}()

	select {
	case imp := <-improvementCh:
		<-workersDone
		return imp.deck, imp.avg, imp.idx, true
	case <-workersDone:
		// A last-moment acceptance may have landed just before all senders returned.
		select {
		case imp := <-improvementCh:
			return imp.deck, imp.avg, imp.idx, true
		default:
		}
		return nil, bestAvg, -1, false
	}
}

// defaultWorkers returns the worker count when callers pass numWorkers<=0. The workload is
// purely CPU-bound, so SMT siblings fight for cache and execution units rather than adding
// throughput: capping at physical cores outperforms defaulting to GOMAXPROCS by ~20% on a
// typical consumer CPU. Still clamped by GOMAXPROCS so a lower user/cgroup override wins.
func defaultWorkers() int {
	maxProcs := runtime.GOMAXPROCS(0)
	physical := cpuid.CPU.PhysicalCores
	if physical <= 0 || physical > maxProcs {
		return maxProcs
	}
	return physical
}

// acceptMutation implements the Metropolis acceptance rule with a noise-floor guard. Strict
// improvements that clear the minImprovement margin (deepAvg > bestAvg + minImprovement)
// always pass. Worse-or-marginal mutations pass with probability exp((deepAvg - bestAvg) /
// T) when T > 0; at T == 0 they're rejected, recovering the classical hill-climb behaviour.
//
// minImprovement guards against infinite loops where shuffle noise lets repeated near-zero
// "wins" keep accepting indefinitely. The probabilistic gate intentionally ignores it so SA
// can still walk through ties / shallow dips to escape local maxima — without that bypass,
// raising the floor would shrink the explorable region of the SA walk in proportion.
func acceptMutation(deepAvg, bestAvg, temperature, minImprovement float64, rng *rand.Rand) bool {
	if deepAvg > bestAvg+minImprovement {
		return true
	}
	if temperature <= 0 {
		return false
	}
	prob := math.Exp((deepAvg - bestAvg) / temperature)
	return rng.Float64() < prob
}
