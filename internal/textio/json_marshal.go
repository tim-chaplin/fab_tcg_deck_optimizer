package textio

// Runtime Deck → JSON encoding. MarshalDeck is the public entry point; toJSON / statsToJSON
// / bestTurnToJSON walk the deck, flatten interface values to names, and sort outputs so
// diffs across runs stay stable.

import (
	"encoding/json"
	"sort"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
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
	f := d.Format
	if f == nil {
		f = format.SilverAge
	}
	return &DeckJSON{
		Hero:      d.Hero.(hero.Hero).Name(),
		Format:    f.Name(),
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
		Matchup:     matchupToJSON(s),
		Runs:        s.Runs,
		Hands:       s.Hands,
		TotalValue:  s.TotalValue,
		Avg:         s.Mean(),
		FirstCycle:  s.FirstCycle,
		SecondCycle: s.SecondCycle,
		Best:        bestTurnToJSON(s.Best),
		Histogram:   s.Histogram,
	}
}

// matchupToJSON returns the matchup only for evaluated stats (Runs > 0); an unevaluated deck
// has no matchup, so nil drops the field.
func matchupToJSON(s deck.Stats) *MatchupJSON {
	if s.Runs == 0 {
		return nil
	}
	return &MatchupJSON{
		IncomingPhysicalDamage: s.IncomingPhysicalDamage,
		IncomingArcaneDamage:   s.IncomingArcaneDamage,
	}
}

// bestTurnToJSON serialises the headline Value, returning the zero value on no recorded
// best so omitempty drops the field.
func bestTurnToJSON(b deck.BestTurn) BestTurnJSON {
	if len(b.BestLine) == 0 {
		return BestTurnJSON{}
	}
	return BestTurnJSON{Value: b.Value}
}
