package deckio

// JSON → runtime Deck decoding: Unmarshal is the public entry point; fromJSON /
// perCardFromJSON / bestTurnFromJSON walk the decoded form, resolve every name through the
// card / weapon / hero registries, and reassemble the Deck. Unknown names fail loudly so a
// corrupted file doesn't produce silent nil-entry crashes downstream.

import (
	"encoding/json"
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Unmarshal decodes JSON produced by Marshal into a *deck.Deck. Returns an error if any card,
// weapon, or hero name isn't recognized.
func Unmarshal(data []byte) (*deck.Deck, error) {
	var dj DeckJSON
	if err := json.Unmarshal(data, &dj); err != nil {
		return nil, err
	}
	return fromJSON(&dj)
}

func fromJSON(dj *DeckJSON) (*deck.Deck, error) {
	h, ok := hero.ByName(dj.Hero)
	if !ok {
		return nil, fmt.Errorf("deckio: unknown hero %q", dj.Hero)
	}
	weapons := make([]weapon.Weapon, len(dj.Weapons))
	for i, name := range dj.Weapons {
		w, ok := weapon.ByName(name)
		if !ok {
			return nil, fmt.Errorf("deckio: unknown weapon %q", name)
		}
		weapons[i] = w
	}
	cs := make([]card.Card, len(dj.Cards))
	for i, name := range dj.Cards {
		id, ok := registry.CardByName(name)
		if !ok {
			return nil, fmt.Errorf("deckio: unknown card %q", name)
		}
		cs[i] = registry.GetCard(id)
	}
	best, err := bestTurnFromJSON(dj.Stats.Best)
	if err != nil {
		return nil, err
	}
	perCardMarginal, err := perCardMarginalFromJSON(dj.Stats.PerCardMarginal)
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
		Runs:            dj.Stats.Runs,
		Hands:           dj.Stats.Hands,
		TotalValue:      dj.Stats.TotalValue,
		FirstCycle:      dj.Stats.FirstCycle,
		SecondCycle:     dj.Stats.SecondCycle,
		Best:            best,
		PerCardMarginal: perCardMarginal,
		Histogram:       dj.Stats.Histogram,
	}
	return d, nil
}

func perCardMarginalFromJSON(entries []CardMarginalStatsJSON) (map[ids.CardID]deck.CardMarginalStats, error) {
	if len(entries) == 0 {
		return nil, nil
	}
	out := make(map[ids.CardID]deck.CardMarginalStats, len(entries))
	for _, e := range entries {
		id, ok := registry.CardByName(e.Card)
		if !ok {
			return nil, fmt.Errorf("deckio: unknown card %q in per_card_marginal stats", e.Card)
		}
		out[id] = deck.CardMarginalStats{
			PresentTotal: e.PresentTotal,
			PresentHands: e.PresentHands,
			AbsentTotal:  e.AbsentTotal,
			AbsentHands:  e.AbsentHands,
		}
	}
	return out, nil
}

// bestTurnFromJSON restores the structured TurnLog plus the headline Value /
// StartingRunechants ints. The structured TurnSummary fields (BestLine, SwungWeapons,
// StartOfTurnAuras, TriggersFromLastTurn, State) aren't reconstructed — fabsim's print path
// renders Log via hand.FormatTurnLog. Returns a zero BestTurn when the JSON has no log.
func bestTurnFromJSON(bj BestTurnJSON) (deck.BestTurn, error) {
	if bj.Log.IsEmpty() {
		return deck.BestTurn{}, nil
	}
	return deck.BestTurn{
		Summary:            hand.TurnSummary{Value: bj.Value},
		StartingRunechants: bj.StartingRunechants,
		Log:                bj.Log,
	}, nil
}
