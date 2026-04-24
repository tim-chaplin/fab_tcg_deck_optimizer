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

func fromJSON(dj *DeckJSON) (*deck.Deck, error) {
	h, ok := hero.ByName(dj.Hero)
	if !ok {
		return nil, fmt.Errorf("deckio: unknown hero %q", dj.Hero)
	}
	weapons := make([]weapon.Weapon, len(dj.Weapons))
	for i, name := range dj.Weapons {
		w, ok := cards.WeaponByName(name)
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
	// Sideboard and Equipment are name-only lists — the optimizer doesn't read them and the
	// registry isn't consulted (so the user can list equipment pieces or any other items
	// the sim doesn't model). Copy the names through verbatim.
	if len(dj.Sideboard) > 0 {
		d.Sideboard = append([]string(nil), dj.Sideboard...)
	}
	if len(dj.Equipment) > 0 {
		d.Equipment = append([]string(nil), dj.Equipment...)
	}
	d.Stats = deck.Stats{
		Runs:        dj.Stats.Runs,
		Hands:       dj.Stats.Hands,
		TotalValue:  dj.Stats.TotalValue,
		FirstCycle:  dj.Stats.FirstCycle,
		SecondCycle: dj.Stats.SecondCycle,
		Best:        best,
		PerCard:     perCard,
		Histogram:   dj.Stats.Histogram,
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
	// Size the rebuilt line to include the arsenal-in entry when present so the "(from
	// arsenal)" tag survives the round-trip. The arsenal-in entry goes at the tail,
	// matching bestUncached's convention (hand cards at indices [0,n); arsenal-in at n).
	lineLen := len(bj.Hand)
	if bj.ArsenalIn != nil {
		lineLen++
	}
	line := make([]hand.CardAssignment, lineLen)
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
	if bj.ArsenalIn != nil {
		ac, err := lookupCardByName(bj.ArsenalIn.Card)
		if err != nil {
			return deck.BestTurn{}, fmt.Errorf("deckio: unknown arsenal_in card %q", bj.ArsenalIn.Card)
		}
		r, err := roleFromString(bj.ArsenalIn.Role)
		if err != nil {
			return deck.BestTurn{}, err
		}
		line[len(bj.Hand)] = hand.CardAssignment{
			Card:         ac,
			Role:         r,
			Contribution: bj.ArsenalIn.Contribution,
			FromArsenal:  true,
		}
	}
	chain, err := rebuildAttackChain(bj, line)
	if err != nil {
		return deck.BestTurn{}, err
	}
	var startOfTurnAuras []card.Card
	if len(bj.StartOfTurnAuras) > 0 {
		startOfTurnAuras = make([]card.Card, len(bj.StartOfTurnAuras))
		for i, name := range bj.StartOfTurnAuras {
			c, err := lookupCardByName(name)
			if err != nil {
				return deck.BestTurn{}, fmt.Errorf("deckio: unknown start_of_turn_aura %q", name)
			}
			startOfTurnAuras[i] = c
		}
	}
	var triggers []hand.TriggerContribution
	if len(bj.TriggersFromLastTurn) > 0 {
		triggers = make([]hand.TriggerContribution, len(bj.TriggersFromLastTurn))
		for i, t := range bj.TriggersFromLastTurn {
			c, err := lookupCardByName(t.Card)
			if err != nil {
				return deck.BestTurn{}, fmt.Errorf("deckio: unknown triggers_from_last_turn card %q", t.Card)
			}
			entry := hand.TriggerContribution{Card: c, Damage: t.Damage}
			if t.Revealed != "" {
				rc, err := lookupCardByName(t.Revealed)
				if err != nil {
					return deck.BestTurn{}, fmt.Errorf("deckio: unknown triggers_from_last_turn revealed %q", t.Revealed)
				}
				entry.Revealed = rc
			}
			triggers[i] = entry
		}
	}
	return deck.BestTurn{
		Summary: hand.TurnSummary{
			BestLine:             line,
			AttackChain:          chain,
			Value:                bj.Value,
			StartOfTurnAuras:     startOfTurnAuras,
			TriggersFromLastTurn: triggers,
		},
		StartingRunechants: bj.StartingRunechants,
	}, nil
}

// rebuildAttackChain reconstructs TurnSummary.AttackChain from the JSON form. When the file has
// an explicit Chain array we use it: it carries true play order plus per-step damage,
// hero-trigger damage, and aura-trigger damage, which FormatBestTurn needs to render "+N"
// contribution labels. Files without a Chain field fall back to a best-effort rebuild
// (hand-order Attack-role cards then weapons) so they still load, though damage labels will
// all read "+0".
func rebuildAttackChain(bj BestTurnJSON, line []hand.CardAssignment) ([]hand.AttackChainEntry, error) {
	if len(bj.Chain) > 0 {
		chain := make([]hand.AttackChainEntry, len(bj.Chain))
		for i, e := range bj.Chain {
			c, err := lookupCardByName(e.Card)
			if err != nil {
				return nil, fmt.Errorf("%v in attack chain", err)
			}
			chain[i] = hand.AttackChainEntry{
				Card:              c,
				Damage:            e.Damage,
				TriggerDamage:     e.TriggerDamage,
				AuraTriggerDamage: e.AuraTriggerDamage,
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
		if w, ok := cards.WeaponByName(name); ok {
			chain = append(chain, hand.AttackChainEntry{Card: w})
		}
	}
	return chain, nil
}

// lookupCardByName resolves a card name from the JSON form to either a registered card or a
// known weapon. Returns an error on unknown names so a corrupted file doesn't silently produce
// nil entries that'd crash FormatBestTurn. Callers wrap the bare error with the field
// context (attack chain, start-of-turn aura, etc.) since those strings differ.
func lookupCardByName(name string) (card.Card, error) {
	if w, ok := cards.WeaponByName(name); ok {
		return w, nil
	}
	id, ok := cards.ByName(name)
	if !ok {
		return nil, fmt.Errorf("deckio: unknown card/weapon %q", name)
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
