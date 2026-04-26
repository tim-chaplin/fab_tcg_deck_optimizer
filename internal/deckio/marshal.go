package deckio

// Runtime Deck → JSON encoding: Marshal is the public entry point; toJSON / statsToJSON /
// bestTurnToJSON walk the deck, flatten interface values to names, and sort the outputs so
// diffs across runs stay stable.

import (
	"encoding/json"
	"sort"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
)

// Marshal returns the JSON encoding of `d` (indented) with card/weapon/hero names in place of
// interface values.
func Marshal(d *deck.Deck) ([]byte, error) {
	return json.MarshalIndent(toJSON(d), "", "  ")
}

func toJSON(d *deck.Deck) *DeckJSON {
	weapons := make([]string, len(d.Weapons))
	for i, w := range d.Weapons {
		weapons[i] = card.DisplayName(w)
	}
	cardNames := make([]string, len(d.Cards))
	var pitchCounts PitchCountsJSON
	for i, c := range d.Cards {
		cardNames[i] = card.DisplayName(c)
		switch c.Pitch() {
		case 1:
			pitchCounts.Red++
		case 2:
			pitchCounts.Yellow++
		case 3:
			pitchCounts.Blue++
		}
	}
	sort.Strings(cardNames)
	return &DeckJSON{
		Hero:      d.Hero.Name(),
		Weapons:   weapons,
		Cards:     cardNames,
		Sideboard: sortedStrings(d.Sideboard),
		Equipment: sortedStrings(d.Equipment),
		Pitch:     pitchCounts,
		Stats:     statsToJSON(d.Stats),
	}
}

// sortedStrings returns a sorted copy of ss. Nil on empty input so omitempty can elide the
// field entirely. Used for name-only parallel lists (Sideboard, Equipment) that don't
// participate in simulation.
func sortedStrings(ss []string) []string {
	if len(ss) == 0 {
		return nil
	}
	out := append([]string(nil), ss...)
	sort.Strings(out)
	return out
}

func statsToJSON(s deck.Stats) StatsJSON {
	return StatsJSON{
		Runs:            s.Runs,
		Hands:           s.Hands,
		TotalValue:      s.TotalValue,
		Avg:             s.Mean(),
		FirstCycle:      s.FirstCycle,
		SecondCycle:     s.SecondCycle,
		Best:            bestTurnToJSON(s.Best),
		PerCardMarginal: perCardMarginalToJSON(s.PerCardMarginal),
		Histogram:       s.Histogram,
	}
}

// perCardMarginalToJSON flattens the card.ID-keyed marginal-stats map into a slice sorted
// by Marginal descending, then by card name — matching the on-screen card-value table's
// order so the JSON and the printout read in lockstep.
func perCardMarginalToJSON(m map[card.ID]deck.CardMarginalStats) []CardMarginalStatsJSON {
	if len(m) == 0 {
		return nil
	}
	out := make([]CardMarginalStatsJSON, 0, len(m))
	for id, s := range m {
		out = append(out, CardMarginalStatsJSON{
			Card:         card.DisplayName(cards.Get(id)),
			PresentTotal: s.PresentTotal,
			PresentHands: s.PresentHands,
			AbsentTotal:  s.AbsentTotal,
			AbsentHands:  s.AbsentHands,
			Marginal:     s.Marginal(),
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Marginal != out[j].Marginal {
			return out[i].Marginal > out[j].Marginal
		}
		return out[i].Card < out[j].Card
	})
	return out
}

// bestTurnToJSON serialises deck.BestTurn.Log directly. The structured TurnSummary stays in
// memory for the live computation but never crosses the JSON boundary — the structured Log
// is the single source of truth on disk and feeds the formatter at print time.
func bestTurnToJSON(b deck.BestTurn) BestTurnJSON {
	if b.Log.IsEmpty() {
		return BestTurnJSON{}
	}
	return BestTurnJSON{
		Value:              b.Summary.Value,
		StartingRunechants: b.StartingRunechants,
		Log:                b.Log,
	}
}
