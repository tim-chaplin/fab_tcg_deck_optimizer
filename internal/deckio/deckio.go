// Package deckio serializes Deck values (plus their accumulated Stats) to and from JSON. Cards,
// weapons, and heroes are referenced by name; deserialization looks names up in package cards and
// in the hero registry below.
package deckio

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// heroesByName resolves Hero names during deserialization. Add new heroes here as implemented.
var heroesByName = map[string]hero.Hero{
	(hero.Viserai{}).Name(): hero.Viserai{},
}

// DeckJSON is the on-disk shape of a Deck with its Stats.
type DeckJSON struct {
	Hero    string           `json:"hero"`
	Weapons []string         `json:"weapons"`
	Cards   []string         `json:"cards"`
	Pitch   PitchCountsJSON  `json:"pitch"`
	Stats   StatsJSON        `json:"stats"`
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
	Runs        int                    `json:"runs"`
	Hands       int                    `json:"hands"`
	TotalValue  float64                `json:"total_value"`
	FirstCycle  deck.CycleStats        `json:"first_cycle"`
	SecondCycle deck.CycleStats        `json:"second_cycle"`
	Best        BestTurnJSON           `json:"best"`
	PerCard     []CardPlayStatsJSON    `json:"per_card,omitempty"`
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

// BestTurnJSON is the JSON form of deck.BestTurn: card names and role names instead of interface
// values. Contributions parallels Hand/Roles and carries CardAssignment.Contribution for each
// hand slot so loaded decks render per-card damage/block/pitch credit instead of "+0". Chain is
// the ordered attack sequence — cards and weapons in play order with their per-step damage —
// so FormatBestTurn can reconstruct the same layout the live sim produced. Both fields are
// omitempty so files produced before this change still load (missing contributions default to
// 0, missing chain falls back to a plausible hand-order reconstruction).
type BestTurnJSON struct {
	Hand               []string               `json:"hand"`
	Roles              []string               `json:"roles"`
	Contributions      []float64              `json:"contributions,omitempty"`
	Weapons            []string               `json:"weapons"`
	Chain              []AttackChainEntryJSON `json:"chain,omitempty"`
	Value              int                    `json:"value"`
	StartingRunechants int                    `json:"starting_runechants"`
}

// AttackChainEntryJSON serialises one attack step (card or weapon) along with the damage it
// dealt in the sim's winning chain. TriggerDamage is the hero's OnCardPlayed contribution for
// that step (e.g. a Runechant Viserai fires as the card resolves) — omitted when zero.
type AttackChainEntryJSON struct {
	Card          string  `json:"card"`
	Damage        float64 `json:"damage"`
	TriggerDamage float64 `json:"trigger_damage,omitempty"`
}

// Marshal returns the JSON encoding of `d` (indented) with card/weapon/hero names in place of
// interface values.
func Marshal(d *deck.Deck) ([]byte, error) {
	return json.MarshalIndent(toJSON(d), "", "  ")
}

// Unmarshal decodes JSON produced by Marshal into a *deck.Deck. Returns an error if any card,
// weapon, or hero name isn't recognized.
func Unmarshal(data []byte) (*deck.Deck, error) {
	var dj DeckJSON
	if err := json.Unmarshal(data, &dj); err != nil {
		return nil, err
	}
	return fromJSON(&dj)
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
		Hero:    d.Hero.Name(),
		Weapons: weapons,
		Cards:   cardNames,
		Pitch:   pitchCounts,
		Stats:   statsToJSON(d.Stats),
	}
}

func statsToJSON(s deck.Stats) StatsJSON {
	return StatsJSON{
		Runs:        s.Runs,
		Hands:       s.Hands,
		TotalValue:  s.TotalValue,
		FirstCycle:  s.FirstCycle,
		SecondCycle: s.SecondCycle,
		Best:        bestTurnToJSON(s.Best),
		PerCard:     perCardToJSON(s.PerCard),
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
	// Serialise hand cards only (arsenal-in entries belong to a previous turn's hand). JSON
	// stays with parallel name + role arrays for human readability / backward compatibility;
	// the in-memory BestLine is still the single source of truth. Weapon names get extracted
	// from the AttackChain since TurnSummary no longer carries them separately.
	var handNames, roles []string
	var contribs []float64
	for _, a := range b.Summary.BestLine {
		if a.FromArsenal {
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
			Card:          e.Card.Name(),
			Damage:        e.Damage,
			TriggerDamage: e.TriggerDamage,
		})
	}
	return BestTurnJSON{
		Hand:               handNames,
		Roles:              roles,
		Contributions:      contribs,
		Weapons:            weaponNames,
		Chain:              chain,
		Value:              b.Summary.Value,
		StartingRunechants: b.StartingRunechants,
	}
}

func fromJSON(dj *DeckJSON) (*deck.Deck, error) {
	h, ok := heroesByName[dj.Hero]
	if !ok {
		return nil, fmt.Errorf("deckio: unknown hero %q", dj.Hero)
	}
	weaponReg := weaponsByName()
	weapons := make([]weapon.Weapon, len(dj.Weapons))
	for i, name := range dj.Weapons {
		w, ok := weaponReg[name]
		if !ok {
			return nil, fmt.Errorf("deckio: unknown weapon %q", name)
		}
		weapons[i] = w
	}
	cs := make([]card.Card, len(dj.Cards))
	for i, name := range dj.Cards {
		id, ok := cards.ByName(name)
		if !ok {
			return nil, fmt.Errorf("deckio: unknown card %q", name)
		}
		cs[i] = cards.Get(id)
	}
	best, err := bestTurnFromJSON(dj.Stats.Best)
	if err != nil {
		return nil, err
	}
	perCard, err := perCardFromJSON(dj.Stats.PerCard)
	if err != nil {
		return nil, err
	}
	d := deck.New(h, weapons, cs)
	d.Stats = deck.Stats{
		Runs:        dj.Stats.Runs,
		Hands:       dj.Stats.Hands,
		TotalValue:  dj.Stats.TotalValue,
		FirstCycle:  dj.Stats.FirstCycle,
		SecondCycle: dj.Stats.SecondCycle,
		Best:        best,
		PerCard:     perCard,
	}
	return d, nil
}

func perCardFromJSON(entries []CardPlayStatsJSON) (map[card.ID]deck.CardPlayStats, error) {
	if len(entries) == 0 {
		return nil, nil
	}
	out := make(map[card.ID]deck.CardPlayStats, len(entries))
	for _, e := range entries {
		id, ok := cards.ByName(e.Card)
		if !ok {
			return nil, fmt.Errorf("deckio: unknown card %q in per_card stats", e.Card)
		}
		out[id] = deck.CardPlayStats{
			Plays:             e.Plays,
			Pitches:           e.Pitches,
			TotalContribution: e.TotalContribution,
		}
	}
	return out, nil
}

func bestTurnFromJSON(bj BestTurnJSON) (deck.BestTurn, error) {
	if len(bj.Hand) == 0 {
		return deck.BestTurn{}, nil
	}
	if len(bj.Roles) != len(bj.Hand) {
		return deck.BestTurn{}, fmt.Errorf("deckio: best turn has %d cards but %d roles", len(bj.Hand), len(bj.Roles))
	}
	if len(bj.Contributions) != 0 && len(bj.Contributions) != len(bj.Hand) {
		return deck.BestTurn{}, fmt.Errorf("deckio: best turn has %d cards but %d contributions", len(bj.Hand), len(bj.Contributions))
	}
	line := make([]hand.CardAssignment, len(bj.Hand))
	for i, name := range bj.Hand {
		id, ok := cards.ByName(name)
		if !ok {
			return deck.BestTurn{}, fmt.Errorf("deckio: unknown card %q in best turn", name)
		}
		r, err := roleFromString(bj.Roles[i])
		if err != nil {
			return deck.BestTurn{}, err
		}
		ca := hand.CardAssignment{Card: cards.Get(id), Role: r}
		if len(bj.Contributions) > 0 {
			ca.Contribution = bj.Contributions[i]
		}
		line[i] = ca
	}
	chain, err := rebuildAttackChain(bj, line)
	if err != nil {
		return deck.BestTurn{}, err
	}
	return deck.BestTurn{
		Summary: hand.TurnSummary{
			BestLine:    line,
			AttackChain: chain,
			Value:       bj.Value,
		},
		StartingRunechants: bj.StartingRunechants,
	}, nil
}

// rebuildAttackChain reconstructs TurnSummary.AttackChain from the JSON form. When the file has
// an explicit Chain array we use it: it carries true play order plus per-step damage and
// hero-trigger damage, which FormatBestTurn needs to render "+N" contribution labels. Files
// without a Chain field fall back to a best-effort rebuild (hand-order Attack-role cards then
// weapons) so they still load, though damage labels will all read "+0".
func rebuildAttackChain(bj BestTurnJSON, line []hand.CardAssignment) ([]hand.AttackChainEntry, error) {
	weaponReg := weaponsByName()
	if len(bj.Chain) > 0 {
		chain := make([]hand.AttackChainEntry, len(bj.Chain))
		for i, e := range bj.Chain {
			c, err := lookupChainCard(e.Card, weaponReg)
			if err != nil {
				return nil, err
			}
			chain[i] = hand.AttackChainEntry{
				Card:          c,
				Damage:        e.Damage,
				TriggerDamage: e.TriggerDamage,
			}
		}
		return chain, nil
	}
	var chain []hand.AttackChainEntry
	for _, a := range line {
		if a.Role == hand.Attack {
			chain = append(chain, hand.AttackChainEntry{Card: a.Card, Damage: a.Contribution})
		}
	}
	for _, name := range bj.Weapons {
		if w, ok := weaponReg[name]; ok {
			chain = append(chain, hand.AttackChainEntry{Card: w})
		}
	}
	return chain, nil
}

// lookupChainCard resolves an attack-chain entry's name to either a registered card or a known
// weapon. Returns an error on unknown names so a corrupted file doesn't silently produce nil
// entries that'd crash FormatBestTurn.
func lookupChainCard(name string, weaponReg map[string]weapon.Weapon) (card.Card, error) {
	if w, ok := weaponReg[name]; ok {
		return w, nil
	}
	id, ok := cards.ByName(name)
	if !ok {
		return nil, fmt.Errorf("deckio: unknown card/weapon %q in attack chain", name)
	}
	return cards.Get(id), nil
}

func roleFromString(s string) (hand.Role, error) {
	switch s {
	case "PITCH":
		return hand.Pitch, nil
	case "ATTACK":
		return hand.Attack, nil
	case "DEFEND":
		return hand.Defend, nil
	case "HELD":
		return hand.Held, nil
	case "ARSENAL":
		return hand.Arsenal, nil
	}
	return 0, fmt.Errorf("deckio: unknown role %q", s)
}

func weaponsByName() map[string]weapon.Weapon {
	m := make(map[string]weapon.Weapon, len(cards.AllWeapons))
	for _, w := range cards.AllWeapons {
		m[w.Name()] = w
	}
	return m
}
