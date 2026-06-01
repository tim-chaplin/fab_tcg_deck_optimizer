package textio

import "github.com/tim-chaplin/fab-deck-optimizer/internal/deck"

// On-disk JSON shapes for a Deck and its accumulated Stats. Every field here trades a
// runtime interface value for a display-name string so files are human-readable and don't
// depend on card-registry indexing.

// DeckJSON is the on-disk shape of a Deck with its Stats. Sideboard and Equipment are
// user-managed name lists the simulator never reads; they round-trip but don't participate
// in scoring. Both omit when empty so files without them stay diff-clean.
type DeckJSON struct {
	Hero      string          `json:"hero"`
	Weapons   []string        `json:"weapons"`
	Cards     []string        `json:"cards"`
	Sideboard []string        `json:"sideboard,omitempty"`
	Equipment []string        `json:"equipment,omitempty"`
	Pitch     PitchCountsJSON `json:"pitch"`
	Stats     StatsJSON       `json:"stats"`
}

// PitchCountsJSON reports how many cards of each pitch colour are in the deck. Derived from
// Cards on marshal, ignored on unmarshal — kept in the file for human readability.
type PitchCountsJSON struct {
	Red    int `json:"red"`
	Yellow int `json:"yellow"`
	Blue   int `json:"blue"`
}

// StatsJSON mirrors deck.Stats with card references flattened to names.
type StatsJSON struct {
	// Matchup is the opponent profile these stats were measured under. Pointer + omitempty so an
	// unevaluated deck (Runs == 0) omits it rather than claim a 0/0 matchup it never faced.
	Matchup *MatchupJSON `json:"matchup,omitempty"`
	// Avg is TotalValue/Hands, emitted for human readability. Unmarshal ignores it and
	// rederives via DeckStats.Mean(); the canonical state is (Runs, Hands, TotalValue).
	Avg         float64         `json:"avg"`
	Runs        int             `json:"runs"`
	Hands       int             `json:"hands"`
	TotalValue  float64         `json:"total_value"`
	FirstCycle  deck.CycleStats `json:"first_cycle"`
	SecondCycle deck.CycleStats `json:"second_cycle"`
	Best        BestTurnJSON    `json:"best"`
	// Histogram counts hands seen at each Value. encoding/json writes int-keyed maps with
	// the int formatted as a string ("7": 42), which round-trips fine since the field type
	// is map[int]int. Omitted when empty.
	Histogram map[int]int `json:"histogram,omitempty"`
}

// BestTurnJSON is the on-disk shape of deck.BestTurn — just the headline Value.
type BestTurnJSON struct {
	Value int `json:"value"`
}

// MatchupJSON is the on-disk shape of the opponent profile a deck's stats were measured under.
type MatchupJSON struct {
	IncomingPhysicalDamage int `json:"incoming_physical_damage"`
	IncomingArcaneDamage   int `json:"incoming_arcane_damage"`
}
