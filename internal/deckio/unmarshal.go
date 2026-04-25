package deckio

// JSON → runtime Deck decoding: Unmarshal is the public entry point; fromJSON /
// perCardFromJSON / bestTurnFromJSON walk the decoded form, resolve every name through the
// card / weapon / hero registries, and reassemble the Deck. Unknown names fail loudly so a
// corrupted file doesn't produce silent nil-entry crashes downstream.

import (
	"encoding/json"
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
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
		PerCard:         perCard,
		PerCardMarginal: perCardMarginal,
		Histogram:       dj.Stats.Histogram,
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

func perCardMarginalFromJSON(entries []CardMarginalStatsJSON) (map[card.ID]deck.CardMarginalStats, error) {
	if len(entries) == 0 {
		return nil, nil
	}
	out := make(map[card.ID]deck.CardMarginalStats, len(entries))
	for _, e := range entries {
		id, ok := cards.ByName(e.Card)
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
		if w, ok := weapon.ByName(name); ok {
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
	if w, ok := weapon.ByName(name); ok {
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
