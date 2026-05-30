package registry

// Persistent card ranking for the anneal search. Ranking holds a total order over the legal card
// pool, best (index 0) to worst, so the search can try mutations involving promising cards first.
// RecordResult bubbles a winning card above a losing one — one comparison at a time — so genuinely
// good cards rise and chaff sinks the more the search runs. The order persists to disk between
// runs; a fresh ranking just seeds the pool in registration order and lets the results settle it.

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// DefaultRankingPath is where the anneal search persists its card ranking, relative to the working
// directory fabsim runs from (the repo root). Gitignored — it's evolving run state, not source.
const DefaultRankingPath = "internal/registry/card_ranking.json"

// Ranking is a total order over the legal card pool, best (index 0) to worst. Safe for concurrent
// RecordResult / Rank.
type Ranking struct {
	mu    sync.Mutex
	order []CardID       // best → worst
	pos   map[CardID]int // CardID → index in order
}

// NewRanking seeds a ranking from the legal pool in its registration order — the arbitrary initial
// ranking that RecordResult settles over time.
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
	for _, c := range (Registry{}).LegalCards() {
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

// RecordResult registers that winner beat loser in a head-to-head. If the winner currently ranks
// below the loser they exchange positions, so the winner rises and the loser sinks one step; if the
// winner already ranks above the loser nothing changes. Unknown cards are ignored.
func (r *Ranking) RecordResult(winner, loser CardID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	wi, wok := r.pos[winner]
	li, lok := r.pos[loser]
	if !wok || !lok || wi <= li {
		return
	}
	r.order[wi], r.order[li] = r.order[li], r.order[wi]
	r.pos[winner], r.pos[loser] = li, wi
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
