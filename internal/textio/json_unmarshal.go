package textio

// JSON → runtime Deck decoding. UnmarshalDeck is the public entry point; fromJSON /
// perCardMarginalFromJSON / bestTurnFromJSON resolve every name through the card / weapon /
// hero registries. Unknown names fail loudly so a corrupted file doesn't produce silent
// nil-entry crashes downstream.

import (
	"encoding/json"
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
)

// UnmarshalDeck decodes JSON produced by MarshalDeck into a *deck.Deck and its Stats.
// Returns an error if any card, weapon, or hero name isn't recognised.
func UnmarshalDeck(data []byte) (*deck.Deck, deck.Stats, error) {
	var dj DeckJSON
	if err := json.Unmarshal(data, &dj); err != nil {
		return nil, deck.Stats{}, err
	}
	return fromJSON(&dj)
}

func fromJSON(dj *DeckJSON) (*deck.Deck, deck.Stats, error) {
	h, ok := registry.HeroByName(dj.Hero)
	if !ok {
		return nil, deck.Stats{}, fmt.Errorf("deckio: unknown hero %q", dj.Hero)
	}
	weapons := make([]deck.Weapon, len(dj.Weapons))
	for i, name := range dj.Weapons {
		w, ok := registry.WeaponByName(name)
		if !ok {
			return nil, deck.Stats{}, fmt.Errorf("deckio: unknown weapon %q", name)
		}
		weapons[i] = w
	}
	cs := make([]deck.Card, len(dj.Cards))
	for i, name := range dj.Cards {
		id, ok := registry.CardByName(name)
		if !ok {
			return nil, deck.Stats{}, fmt.Errorf("deckio: unknown card %q", name)
		}
		cs[i] = registry.GetCard(id)
	}
	best, err := bestTurnFromJSON(dj.Stats.Best)
	if err != nil {
		return nil, deck.Stats{}, err
	}
	perCardMarginal, err := perCardMarginalFromJSON(dj.Stats.PerCardMarginal)
	if err != nil {
		return nil, deck.Stats{}, err
	}
	d := deck.New(h, weapons, cs)
	// Sideboard and Equipment are name-only lists — the registry isn't consulted, so the
	// user can list equipment pieces or any other items the sim doesn't model.
	if len(dj.Sideboard) > 0 {
		d.Sideboard = append([]string(nil), dj.Sideboard...)
	}
	if len(dj.Equipment) > 0 {
		d.Equipment = append([]string(nil), dj.Equipment...)
	}
	stats := deck.Stats{
		Runs:            dj.Stats.Runs,
		Hands:           dj.Stats.Hands,
		TotalValue:      dj.Stats.TotalValue,
		FirstCycle:      dj.Stats.FirstCycle,
		SecondCycle:     dj.Stats.SecondCycle,
		Best:            best,
		PerCardMarginal: perCardMarginal,
		Histogram:       dj.Stats.Histogram,
	}
	return d, stats, nil
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

// bestTurnFromJSON restores the headline Value, returning the zero BestTurn when none was
// recorded.
func bestTurnFromJSON(bj BestTurnJSON) (deck.BestTurn, error) {
	if bj.Value == 0 {
		return deck.BestTurn{}, nil
	}
	return deck.BestTurn{Value: bj.Value}, nil
}
