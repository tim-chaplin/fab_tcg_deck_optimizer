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

// heroesByName is the fixed registry used to resolve Hero names during deserialization. Add new
// heroes here as they're implemented.
var heroesByName = map[string]hero.Hero{
	(hero.Viserai{}).Name(): hero.Viserai{},
}

// DeckJSON is the on-disk shape of a Deck with its Stats.
type DeckJSON struct {
	Hero    string    `json:"hero"`
	Weapons []string  `json:"weapons"`
	Cards   []string  `json:"cards"`
	Stats   StatsJSON `json:"stats"`
}

// StatsJSON mirrors deck.Stats with card references flattened to names.
type StatsJSON struct {
	Runs        int            `json:"runs"`
	Hands       int            `json:"hands"`
	TotalValue  float64        `json:"total_value"`
	FirstCycle  deck.CycleStats `json:"first_cycle"`
	SecondCycle deck.CycleStats `json:"second_cycle"`
	Best        BestHandJSON   `json:"best"`
}

// BestHandJSON is the JSON form of deck.BestHand: card names and role names instead of interface
// values.
type BestHandJSON struct {
	Hand    []string `json:"hand"`
	Roles   []string `json:"roles"`
	Weapons []string `json:"weapons"`
	Value   int      `json:"value"`
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
	for i, c := range d.Cards {
		cardNames[i] = c.Name()
	}
	sort.Strings(cardNames)
	return &DeckJSON{
		Hero:    d.Hero.Name(),
		Weapons: weapons,
		Cards:   cardNames,
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
		Best:        bestHandToJSON(s.Best),
	}
}

func bestHandToJSON(b deck.BestHand) BestHandJSON {
	if b.Hand == nil {
		return BestHandJSON{}
	}
	handNames := make([]string, len(b.Hand))
	for i, c := range b.Hand {
		handNames[i] = c.Name()
	}
	roles := make([]string, len(b.Play.Roles))
	for i, r := range b.Play.Roles {
		roles[i] = r.String()
	}
	return BestHandJSON{
		Hand:    handNames,
		Roles:   roles,
		Weapons: append([]string(nil), b.Play.Weapons...),
		Value:   b.Play.Value,
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
	best, err := bestHandFromJSON(dj.Stats.Best)
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
	}
	return d, nil
}

func bestHandFromJSON(bj BestHandJSON) (deck.BestHand, error) {
	if len(bj.Hand) == 0 {
		return deck.BestHand{}, nil
	}
	if len(bj.Roles) != len(bj.Hand) {
		return deck.BestHand{}, fmt.Errorf("deckio: best hand has %d cards but %d roles", len(bj.Hand), len(bj.Roles))
	}
	cs := make([]card.Card, len(bj.Hand))
	for i, name := range bj.Hand {
		id, ok := cards.ByName(name)
		if !ok {
			return deck.BestHand{}, fmt.Errorf("deckio: unknown card %q in best hand", name)
		}
		cs[i] = cards.Get(id)
	}
	roles := make([]hand.Role, len(bj.Roles))
	for i, name := range bj.Roles {
		r, err := roleFromString(name)
		if err != nil {
			return deck.BestHand{}, err
		}
		roles[i] = r
	}
	return deck.BestHand{
		Hand: cs,
		Play: hand.Play{
			Roles:   roles,
			Weapons: append([]string(nil), bj.Weapons...),
			Value:   bj.Value,
		},
	}, nil
}

func roleFromString(s string) (hand.Role, error) {
	switch s {
	case "PITCH":
		return hand.Pitch, nil
	case "ATTACK":
		return hand.Attack, nil
	case "DEFEND":
		return hand.Defend, nil
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
