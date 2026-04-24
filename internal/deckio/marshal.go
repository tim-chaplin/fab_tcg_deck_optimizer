package deckio

// Runtime Deck → JSON encoding: Marshal is the public entry point; toJSON / statsToJSON /
// perCardToJSON / bestTurnToJSON walk the deck, flatten interface values to names, and sort
// the outputs so diffs across runs stay stable.

import (
	"encoding/json"
	"sort"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
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
		Runs:        s.Runs,
		Hands:       s.Hands,
		TotalValue:  s.TotalValue,
		Avg:         s.Mean(),
		FirstCycle:  s.FirstCycle,
		SecondCycle: s.SecondCycle,
		Best:        bestTurnToJSON(s.Best),
		PerCard:     perCardToJSON(s.PerCard),
		Histogram:   s.Histogram,
	}
}

// perCardToJSON flattens the card.ID-keyed map into a slice sorted by Avg descending, total
// appearances descending, then card name — so the JSON output is stable and the best-performing
// cards surface at the top.
func perCardToJSON(m map[card.ID]deck.CardPlayStats) []CardPlayStatsJSON {
	if len(m) == 0 {
		return nil
	}
	out := make([]CardPlayStatsJSON, 0, len(m))
	for id, s := range m {
		out = append(out, CardPlayStatsJSON{
			Card:              cards.Get(id).Name(),
			Plays:             s.Plays,
			Pitches:           s.Pitches,
			TotalContribution: s.TotalContribution,
			Avg:               s.Avg(),
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Avg != out[j].Avg {
			return out[i].Avg > out[j].Avg
		}
		ni, nj := out[i].Plays+out[i].Pitches, out[j].Plays+out[j].Pitches
		if ni != nj {
			return ni > nj
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
	// the single source of truth. Weapon names are extracted from AttackChain. Arsenal-in
	// entries are emitted separately in ArsenalIn so reload can re-append them with
	// FromArsenal=true, which keeps the "(from arsenal)" tag on the printout.
	var handNames, roles []string
	var contribs []float64
	var arsenalIn *ArsenalInJSON
	for _, a := range b.Summary.BestLine {
		if a.FromArsenal {
			arsenalIn = &ArsenalInJSON{
				Card:         a.Card.Name(),
				Role:         a.Role.String(),
				Contribution: a.Contribution,
			}
			continue
		}
		handNames = append(handNames, a.Card.Name())
		roles = append(roles, a.Role.String())
		contribs = append(contribs, a.Contribution)
	}
	var weaponNames []string
	var chain []AttackChainEntryJSON
	for _, e := range b.Summary.AttackChain {
		if w, ok := e.Card.(weapon.Weapon); ok {
			weaponNames = append(weaponNames, w.Name())
		}
		chain = append(chain, AttackChainEntryJSON{
			Card:              e.Card.Name(),
			Damage:            e.Damage,
			TriggerDamage:     e.TriggerDamage,
			AuraTriggerDamage: e.AuraTriggerDamage,
		})
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
		Contributions:        contribs,
		Weapons:              weaponNames,
		Chain:                chain,
		StartOfTurnAuras:     startOfTurnAuras,
		ArsenalIn:            arsenalIn,
		TriggersFromLastTurn: triggers,
		Value:                b.Summary.Value,
		StartingRunechants:   b.StartingRunechants,
	}
}
