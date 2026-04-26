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
		weapons[i] = w.Name()
	}
	cardNames := make([]string, len(d.Cards))
	var pitchCounts PitchCountsJSON
	for i, c := range d.Cards {
		cardNames[i] = c.Name()
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
			Card:         cards.Get(id).Name(),
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

func bestTurnToJSON(b deck.BestTurn) BestTurnJSON {
	if len(b.Summary.BestLine) == 0 {
		return BestTurnJSON{}
	}
	// JSON uses parallel name + role arrays for human readability; the in-memory BestLine is
	// the single source of truth. Arsenal-in entries are emitted separately in ArsenalIn so
	// reload can re-append them with FromArsenal=true, which keeps the "(from arsenal)" tag
	// on the printout.
	var handNames, roles []string
	var arsenalIn *ArsenalInJSON
	for _, a := range b.Summary.BestLine {
		if a.FromArsenal {
			arsenalIn = &ArsenalInJSON{
				Card: a.Card.Name(),
				Role: a.Role.String(),
			}
			continue
		}
		handNames = append(handNames, a.Card.Name())
		roles = append(roles, a.Role.String())
	}
	var weaponNames []string
	if len(b.Summary.SwungWeapons) > 0 {
		weaponNames = append([]string(nil), b.Summary.SwungWeapons...)
	}
	var startOfTurnAuras []string
	if len(b.Summary.StartOfTurnAuras) > 0 {
		startOfTurnAuras = make([]string, len(b.Summary.StartOfTurnAuras))
		for i, a := range b.Summary.StartOfTurnAuras {
			startOfTurnAuras[i] = a.Name()
		}
	}
	var triggers []TriggerContributionJSON
	for _, t := range b.Summary.TriggersFromLastTurn {
		entry := TriggerContributionJSON{Card: t.Card.Name(), Damage: t.Damage}
		if t.Revealed != nil {
			entry.Revealed = t.Revealed.Name()
		}
		triggers = append(triggers, entry)
	}
	return BestTurnJSON{
		Hand:                 handNames,
		Roles:                roles,
		Weapons:              weaponNames,
		StartOfTurnAuras:     startOfTurnAuras,
		ArsenalIn:            arsenalIn,
		TriggersFromLastTurn: triggers,
		Value:                b.Summary.Value,
		StartingRunechants:   b.StartingRunechants,
	}
}
