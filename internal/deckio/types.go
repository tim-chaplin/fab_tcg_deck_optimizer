package deckio

// On-disk JSON shapes for a Deck and its accumulated Stats. Every field here trades a runtime
// interface value for a display-name string so files are human-readable and don't depend on
// card-registry indexing. Marshal / Unmarshal in their own files convert between these and
// the runtime Deck / Stats.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
)

// DeckJSON is the on-disk shape of a Deck with its Stats. Sideboard and Equipment are
// user-managed parallel card lists that the simulator never reads — both round-trip through
// Marshal / Unmarshal but don't participate in scoring, mutations, or NotImplemented
// sanitization. Each is omitted from the JSON when empty so existing files stay untouched.
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
// Cards on marshal and ignored on unmarshal (kept in the file purely for human readability).
type PitchCountsJSON struct {
	Red    int `json:"red"`
	Yellow int `json:"yellow"`
	Blue   int `json:"blue"`
}

// StatsJSON mirrors deck.Stats with card references flattened to names.
type StatsJSON struct {
	// Avg is TotalValue/Hands, emitted for human readability when skimming the JSON. Loaders
	// ignore it — Unmarshal rederives via Stats.Mean() so the canonical state is always
	// (Runs, Hands, TotalValue). Kept first so it's the first number a human sees.
	Avg             float64                 `json:"avg"`
	Runs            int                     `json:"runs"`
	Hands           int                     `json:"hands"`
	TotalValue      float64                 `json:"total_value"`
	FirstCycle      deck.CycleStats         `json:"first_cycle"`
	SecondCycle     deck.CycleStats         `json:"second_cycle"`
	Best            BestTurnJSON            `json:"best"`
	PerCardMarginal []CardMarginalStatsJSON `json:"per_card_marginal,omitempty"`
	// Histogram counts hands seen at each Value. encoding/json writes int-keyed maps with the
	// int formatted as a string ("7": 42), which round-trips fine since we declare the field
	// as map[int]int. Omitted when empty so old files stay valid.
	Histogram map[int]int `json:"histogram,omitempty"`
}

// CardMarginalStatsJSON is the JSON form of deck.CardMarginalStats keyed by card name.
// Marginal (PresentMean - AbsentMean) is the actionable smell-test signal a human reader
// scans for, so it's included alongside the raw with/without sums even though it's
// derivable.
type CardMarginalStatsJSON struct {
	Card         string  `json:"card"`
	PresentTotal float64 `json:"present_total"`
	PresentHands int     `json:"present_hands"`
	AbsentTotal  float64 `json:"absent_total"`
	AbsentHands  int     `json:"absent_hands"`
	Marginal     float64 `json:"marginal"`
}

// BestTurnJSON is the on-disk shape of deck.BestTurn — just the rendered printout lines.
// Marshal serialises deck.BestTurn.Lines verbatim; Unmarshal restores them into the loaded
// deck so fabsim's print path can dump them directly without reconstructing a TurnSummary.
// Value and StartingRunechants ride along for diff-friendly searchability of the file but
// don't drive the printout — Lines already includes the "Best turn played (value N):"
// header.
type BestTurnJSON struct {
	Value              int      `json:"value"`
	StartingRunechants int      `json:"starting_runechants,omitempty"`
	Lines              []string `json:"lines,omitempty"`
}
