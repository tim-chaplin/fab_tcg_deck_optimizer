package registry

// Persistent card ranking for the anneal search. Ranking holds a total order over the legal card
// pool, best (index 0) to worst, so the search can try mutations involving promising cards first.
// PromoteAll lifts a whole deck's cards to the front each time a swap improves the deck, making the
// order a recency list — the cards of recently-good decks lead, untried chaff trails — so loading a
// tuned deck bootstraps the order from it on the first improvement. The order persists to disk
// between runs; a fresh ranking just seeds the pool in registration order and lets promotions
// settle it.

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/format"
)

// DefaultRankingPath is where the anneal search persists its card ranking, relative to the working
// directory fabsim runs from (the repo root). Gitignored — it's evolving run state, not source.
const DefaultRankingPath = "internal/registry/card_ranking.json"

// Ranking is a total order over the legal card pool, best (index 0) to worst. Safe for concurrent
// PromoteAll / Rank.
type Ranking struct {
	mu    sync.Mutex
	order []CardID       // best → worst
	pos   map[CardID]int // CardID → index in order
}

// NewRanking seeds a ranking from the legal pool in its registration order — the arbitrary initial
// ranking that Promote settles over time.
func NewRanking() *Ranking {
	return rankingFromNames(nil)
}

// rankingFromNames builds a ranking from saved DisplayNames (best → worst): names no longer in the
// pool are dropped, and any pool cards the list omits (new / first-seen cards) are appended at the
// bottom with an arbitrary initial rank.
func rankingFromNames(names []string) *Ranking {
	r := &Ranking{}
	seen := make(map[CardID]bool)
	for _, name := range names {
		id, ok := CardByName(name)
		if !ok || seen[id] {
			continue
		}
		r.order = append(r.order, id)
		seen[id] = true
	}
	// Hero-agnostic: the ranking is a global recency order keyed by CardID, so it spans the
	// whole format-legal pool. Cards illegal for a given run's hero simply never surface as
	// mutations (buildLegalByID filters per hero), so their presence here is harmless.
	for _, c := range legalCardsForFormat(format.SilverAge) {
		if id := c.ID(); !seen[id] {
			r.order = append(r.order, id)
			seen[id] = true
		}
	}
	r.reindex()
	return r
}

// reindex rebuilds pos from order. Callers hold mu (or hold no concurrent access, as in construction).
func (r *Ranking) reindex() {
	r.pos = make(map[CardID]int, len(r.order))
	for i, id := range r.order {
		r.pos[id] = i
	}
}

// Rank returns the card's position, 0 = best. A card not in the ranking sorts last.
func (r *Ranking) Rank(id CardID) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p, ok := r.pos[id]; ok {
		return p
	}
	return len(r.order)
}

// PromoteAll lifts every listed card to the front of the order — a stable partition that keeps the
// listed cards in their current relative order and the rest in theirs. Used to promote a whole
// deck's cards when a swap improved it, so the deck (staples included) leads the ranking. Cards not
// in the ranking are ignored; an empty or all-unknown list is a no-op.
func (r *Ranking) PromoteAll(ids []CardID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	promote := make(map[CardID]bool, len(ids))
	for _, id := range ids {
		if _, ok := r.pos[id]; ok {
			promote[id] = true
		}
	}
	if len(promote) == 0 {
		return
	}
	front := make([]CardID, 0, len(promote))
	back := make([]CardID, 0, len(r.order)-len(promote))
	for _, id := range r.order { // walking r.order preserves each group's current relative order
		if promote[id] {
			front = append(front, id)
		} else {
			back = append(back, id)
		}
	}
	r.order = append(front, back...)
	r.reindex()
}

// Save writes the ranking to path as an indented JSON array of DisplayNames, best → worst.
func (r *Ranking) Save(path string) error {
	r.mu.Lock()
	names := make([]string, len(r.order))
	for i, id := range r.order {
		names[i] = GetCard(id).DisplayName()
	}
	r.mu.Unlock()
	data, err := json.MarshalIndent(names, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadRanking reads a ranking from path (a JSON array of DisplayNames, best → worst). A missing file
// yields a fresh pool-order ranking; pool cards absent from the file are seeded at the bottom.
func LoadRanking(path string) (*Ranking, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return NewRanking(), nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	if err := json.Unmarshal(data, &names); err != nil {
		return nil, err
	}
	return rankingFromNames(names), nil
}
