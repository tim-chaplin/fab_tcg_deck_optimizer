package textio

// Runtime Deck → JSON encoding. MarshalDeck is the public entry point; toJSON / statsToJSON
// / bestTurnToJSON walk the deck, flatten interface values to names, and sort outputs so
// diffs across runs stay stable.

import (
	"encoding/json"
	"sort"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
)

// MarshalDeck returns the indented JSON encoding of d and stats, with card / weapon / hero
// names in place of interface values.
func MarshalDeck(d *deck.Deck, stats deck.Stats) ([]byte, error) {
	return json.MarshalIndent(toJSON(d, stats), "", "  ")
}

func toJSON(d *deck.Deck, stats deck.Stats) *DeckJSON {
	weapons := make([]string, len(d.Weapons))
	for i, w := range d.Weapons {
		weapons[i] = w.Name()
	}
	cardNames := d.DisplayNames()
	sort.Strings(cardNames)
	red, yellow, blue := d.PitchCounts()
	return &DeckJSON{
		Hero:      d.Hero.(hero.Hero).Name(),
		Weapons:   weapons,
		Cards:     cardNames,
		Sideboard: sortedStrings(d.Sideboard),
		Equipment: sortedStrings(d.Equipment),
		Pitch:     PitchCountsJSON{Red: red, Yellow: yellow, Blue: blue},
		Stats:     statsToJSON(stats),
	}
}

// sortedStrings returns a sorted copy of ss, nil on empty input so omitempty elides the
// field entirely.
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

// perCardMarginalToJSON flattens the ids.CardID-keyed marginal-stats map into a slice
// sorted by Marginal descending, then by card name — matching the on-screen card-value
// table's order.
func perCardMarginalToJSON(m map[ids.CardID]deck.CardMarginalStats) []CardMarginalStatsJSON {
	if len(m) == 0 {
		return nil
	}
	out := make([]CardMarginalStatsJSON, 0, len(m))
	for id, s := range m {
		out = append(out, CardMarginalStatsJSON{
			Card:         registry.GetCard(id).DisplayName(),
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

// bestTurnToJSON serialises the headline Value, returning the zero value on no recorded
// best so omitempty drops the field.
func bestTurnToJSON(b deck.BestTurn) BestTurnJSON {
	if len(b.BestLine) == 0 {
		return BestTurnJSON{}
	}
	return BestTurnJSON{Value: b.Value}
}
