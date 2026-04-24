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
	Avg         float64             `json:"avg"`
	Runs        int                 `json:"runs"`
	Hands       int                 `json:"hands"`
	TotalValue  float64             `json:"total_value"`
	FirstCycle  deck.CycleStats     `json:"first_cycle"`
	SecondCycle deck.CycleStats     `json:"second_cycle"`
	Best        BestTurnJSON        `json:"best"`
	PerCard     []CardPlayStatsJSON `json:"per_card,omitempty"`
	// Histogram counts hands seen at each Value. encoding/json writes int-keyed maps with the
	// int formatted as a string ("7": 42), which round-trips fine since we declare the field
	// as map[int]int. Omitted when empty so old files stay valid.
	Histogram map[int]int `json:"histogram,omitempty"`
}

// CardPlayStatsJSON is the JSON form of deck.CardPlayStats keyed by card name. Avg is included
// even though it's derivable from the other fields — it's what a human reader actually wants
// when skimming the file.
type CardPlayStatsJSON struct {
	Card              string  `json:"card"`
	Plays             int     `json:"plays"`
	Pitches           int     `json:"pitches"`
	TotalContribution float64 `json:"total_contribution"`
	Avg               float64 `json:"avg"`
}

// BestTurnJSON is the JSON form of deck.BestTurn: card names and role names instead of
// interface values. Contributions parallels Hand/Roles and carries
// CardAssignment.Contribution for each hand slot. Chain is the ordered attack sequence —
// cards and weapons in play order with their per-step damage. StartOfTurnAuras mirrors
// hand.TurnSummary.StartOfTurnAuras as a list of card names; ArsenalIn carries the card
// that started the turn in the arsenal slot (distinct from Hand, which is the dealt hand
// only). TriggersFromLastTurn records carryover AuraTrigger fires so the "from previous
// turn" printout round-trips across a reload. Omitempty-omitted fields fall back to
// defaults (contributions = 0, chain rebuilt in hand order, no prior auras, no arsenal-in,
// no carryover trigger lines).
type BestTurnJSON struct {
	Hand                 []string                  `json:"hand"`
	Roles                []string                  `json:"roles"`
	Contributions        []float64                 `json:"contributions,omitempty"`
	Weapons              []string                  `json:"weapons"`
	Chain                []AttackChainEntryJSON    `json:"chain,omitempty"`
	StartOfTurnAuras     []string                  `json:"start_of_turn_auras,omitempty"`
	ArsenalIn            *ArsenalInJSON            `json:"arsenal_in,omitempty"`
	TriggersFromLastTurn []TriggerContributionJSON `json:"triggers_from_last_turn,omitempty"`
	Value                int                       `json:"value"`
	StartingRunechants   int                       `json:"starting_runechants"`
}

// ArsenalInJSON carries the arsenal-in card's role-assigned entry for the best turn so a
// reloaded deck can re-render the "(from arsenal)" tag in the play order. Hand serialises
// only dealt-hand entries; the arsenal-in card lives here separately because it belongs to
// a previous turn and isn't part of the dealt hand the reader sees in the Card list.
type ArsenalInJSON struct {
	Card         string  `json:"card"`
	Role         string  `json:"role"`
	Contribution float64 `json:"contribution,omitempty"`
}

// TriggerContributionJSON is the serialised form of hand.TriggerContribution — one
// carryover AuraTrigger fire at the top of the saved turn. Damage / Revealed are both
// omitempty so a zero-damage non-reveal entry would still round-trip as just the aura
// name (shouldn't happen in practice since FormatBestTurn drops those lines, but the
// shape stays lossless).
type TriggerContributionJSON struct {
	Card     string `json:"card"`
	Damage   int    `json:"damage,omitempty"`
	Revealed string `json:"revealed,omitempty"`
}

// AttackChainEntryJSON serialises one attack step (card or weapon) with the damage it dealt
// in the sim's winning chain. TriggerDamage is the hero's OnCardPlayed contribution for that
// step; AuraTriggerDamage is the mid-chain AuraTrigger contribution (e.g. a prior-turn Malefic
// Incantation firing on this attack). Both are omitted when zero.
type AttackChainEntryJSON struct {
	Card              string  `json:"card"`
	Damage            float64 `json:"damage"`
	TriggerDamage     float64 `json:"trigger_damage,omitempty"`
	AuraTriggerDamage float64 `json:"aura_trigger_damage,omitempty"`
}
