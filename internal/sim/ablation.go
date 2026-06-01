package sim

// Coupled per-card ablation: how much mean hand value a deck loses when each card is removed.
// Kept separate from the per-shuffle eval in deck_eval.go because the orchestration — shuffle
// once, then replay that one order for the full deck and every card-removed copy — is its own
// concern. Sharing the shuffle is what makes the deltas tight: the cards a removal leaves
// untouched see identical draws, so the common-random-numbers noise cancels in the difference.

import (
	"math/rand"
	"sync"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// PerCardAblation returns, per unique card in d, the mean hand value the deck loses without it:
// full-deck avg minus the avg of the deck with every copy of that card removed. Each shuffle is
// played by the full deck and every card-removed copy on the SAME shuffled order, so each delta
// carries the coupled (common-random-numbers) variance reduction — far fewer shuffles resolve the
// small per-card contributions than independent re-sims would. Runs ~(unique cards + 1) plays per
// shuffle, parallel across shuffles over ev.numWorkers, sharing ev.cache.
func (ev *Evaluator) PerCardAblation(d *deck.Deck, runs int, mp Matchup, rng *rand.Rand) map[ids.CardID]float64 {
	handSize := d.Hero.(hero.Hero).Intelligence()
	deckSize := d.Size()
	if handSize <= 0 || deckSize < handSize || runs <= 0 {
		return nil
	}
	uniqueIDs, _ := d.UniqueIDs()
	baseHandsPerCycle := deckSize / handSize
	// Score only cards whose removal still leaves a playable deck (>= one hand); a card whose
	// removal would shrink the deck below handSize can't be ablated, so omit it. Each ablated deck
	// is shorter, so it cycles in fewer hands — its own length sets the fair per-hand baseline.
	type ablationTarget struct {
		id            ids.CardID
		handsPerCycle int
	}
	targets := make([]ablationTarget, 0, len(uniqueIDs))
	for _, id := range uniqueIDs {
		if ablatedSize := deckSize - d.Count(id); ablatedSize >= handSize {
			targets = append(targets, ablationTarget{id: id, handsPerCycle: ablatedSize / handSize})
		}
	}

	numWorkers := ev.numWorkers
	if numWorkers < 1 {
		numWorkers = 1
	}

	// One worker's slice of the run: the full deck's stats plus per-target ablation stats
	// (index-aligned to targets), each accumulated across the shuffles that worker drew.
	type ablationAccum struct {
		base    deck.Stats
		ablated []deck.Stats
	}
	results := make(chan ablationAccum, numWorkers)

	runsPerWorker := runs / numWorkers
	extras := runs % numWorkers
	spawned := 0
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		myRuns := runsPerWorker
		if w < extras {
			myRuns++
		}
		if myRuns == 0 {
			continue
		}
		spawned++
		seed := rng.Int63()
		wg.Add(1)
		go func(seed int64, myRuns int) {
			defer wg.Done()
			workerEv := &Evaluator{cache: ev.cache, statePool: gameengine.NewPrewarmedPool()}
			workerRNG := rand.New(rand.NewSource(seed))
			scratch := newShuffleScratch(handSize)
			a := ablationAccum{ablated: make([]deck.Stats, len(targets))}
			for r := 0; r < myRuns; r++ {
				full := d.Copy()
				full.Shuffle(workerRNG)
				// Play copies derived from this one shuffled order: the full deck, then each
				// card-removed variant (Without preserves the surviving cards' order).
				playOrderedDeck(full.Copy(), scratch, &a.base, workerEv, mp, baseHandsPerCycle, handSize)
				for i, tgt := range targets {
					playOrderedDeck(full.Without(tgt.id), scratch, &a.ablated[i], workerEv, mp, tgt.handsPerCycle, handSize)
				}
			}
			results <- a
		}(seed, myRuns)
	}
	wg.Wait()

	var base deck.Stats
	ablated := make([]deck.Stats, len(targets))
	for i := 0; i < spawned; i++ {
		a := <-results
		mergeStatsInto(&base, &a.base)
		for j := range ablated {
			mergeStatsInto(&ablated[j], &a.ablated[j])
		}
	}

	baseAvg := base.Mean()
	deltas := make(map[ids.CardID]float64, len(targets))
	for i, tgt := range targets {
		deltas[tgt.id] = baseAvg - ablated[i].Mean()
	}
	return deltas
}
