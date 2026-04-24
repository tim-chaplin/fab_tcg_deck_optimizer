package deck

// Single-slot mutation generation for the iterate-mode hill climb: every alternative weapon
// loadout plus every (remove one, add one) card swap the deck admits. Ordering is deterministic
// so first-improvement search explores low-value removal slots first.

import (
	"fmt"
	"sort"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Mutation is one candidate single-slot change: the mutated Deck plus a human-readable summary
// (e.g. "swapped Aether Slash (Red) for Arcanic Spike (Red)"). Consumers use Deck to evaluate
// and Description for logging.
type Mutation struct {
	Deck        *Deck
	Description string
}

// AllMutations returns every single-card mutation of d in a deterministic order: first every
// alternative weapon loadout (sorted by loadout key), then every (removeID, addID) pair where
// one copy of removeID is dropped and one copy of addID is added. removeID must be in the deck;
// addID's post-mutation count must not exceed maxCopies. Pairs with removeID == addID are
// skipped.
//
// Card-mutation ordering: the outer loop iterates uniqueIDs by ascending per-card average
// contribution (d.Stats.PerCard[id].Avg()), so low-value cards get swap candidates tried first.
// Cards without stats tie at 0 and fall back to card.ID. The inner loop iterates the addID pool
// by card.ID. Favouring low-value removal slots surfaces useful swaps early for a
// first-improvement hill climb.
//
// Single-card swaps (not paired swaps) let the hill climber reach decks with odd per-card counts
// (e.g. 1× X + 3× Y at maxCopies=3).
//
// legal filters the addition pool: only accepted IDs become swap-in candidates, so format-banned
// cards can't be introduced. Removal targets aren't filtered — a deck that entered the climb
// holding a banned card can still have it swapped out. Pass nil to skip filtering.
//
// Returned decks have zero Stats and share no backing slices with d or each other.
func AllMutations(d *Deck, maxCopies int, legal func(card.Card) bool) []Mutation {
	out := weaponLoadoutMutations(d)
	out = append(out, cardSwapMutations(d, maxCopies, legal)...)
	return out
}

// weaponLoadoutMutations emits one Mutation per distinct weapon loadout that isn't the current
// one. Loadouts are canonicalised by weaponKey (names sorted) and processed in key order so
// the output is deterministic regardless of map-iteration randomness.
func weaponLoadoutMutations(d *Deck) []Mutation {
	loadouts := weaponLoadouts(cards.AllWeapons)
	currentKey := weaponKey(d.Weapons)
	type keyedLoadout struct {
		key     string
		weapons []weapon.Weapon
	}
	sortedLoadouts := make([]keyedLoadout, 0, len(loadouts))
	for _, l := range loadouts {
		sortedLoadouts = append(sortedLoadouts, keyedLoadout{key: weaponKey(l), weapons: l})
	}
	sort.Slice(sortedLoadouts, func(i, j int) bool { return sortedLoadouts[i].key < sortedLoadouts[j].key })
	var out []Mutation
	for _, l := range sortedLoadouts {
		if l.key == currentKey {
			continue
		}
		newCards := make([]card.Card, len(d.Cards))
		copy(newCards, d.Cards)
		nd := New(d.Hero, l.weapons, newCards)
		nd.Sideboard = d.Sideboard
		nd.Equipment = d.Equipment
		out = append(out, Mutation{
			Deck:        nd,
			Description: fmt.Sprintf("swapped weapons from %s to %s", loadoutLabel(d.Weapons), loadoutLabel(l.weapons)),
		})
	}
	return out
}

// cardSwapMutations emits every single-card remove+add mutation the deck admits. Remove targets
// iterate in ascending per-card avg contribution so the hill climb spends its budget on the
// currently-worst cards first; with no Stats yet the tiebreak falls through to stable card.ID
// order. Add candidates skip no-ops (same ID) and entries already at maxCopies.
func cardSwapMutations(d *Deck, maxCopies int, legal func(card.Card) bool) []Mutation {
	counts := map[card.ID]int{}
	for _, c := range d.Cards {
		counts[c.ID()]++
	}
	uniqueIDs := make([]card.ID, 0, len(counts))
	for id := range counts {
		uniqueIDs = append(uniqueIDs, id)
	}
	sort.Slice(uniqueIDs, func(i, j int) bool {
		ai := d.Stats.PerCard[uniqueIDs[i]].Avg()
		aj := d.Stats.PerCard[uniqueIDs[j]].Avg()
		if ai != aj {
			return ai < aj
		}
		return uniqueIDs[i] < uniqueIDs[j]
	})

	// legalPool returns IDs in ascending order (cards.Deckable() iterates byID).
	pool := legalPool(legal)

	var out []Mutation
	for _, removeID := range uniqueIDs {
		removed := cards.Get(removeID)
		for _, addID := range pool {
			if addID == removeID {
				continue // no-op: remove one and add one of the same card.
			}
			if counts[addID] >= maxCopies {
				continue // at max copies.
			}
			replacement := cards.Get(addID)
			newCards := make([]card.Card, 0, len(d.Cards))
			removed1 := false
			for _, c := range d.Cards {
				if !removed1 && c.ID() == removeID {
					removed1 = true
					continue
				}
				newCards = append(newCards, c)
			}
			newCards = append(newCards, replacement)
			nd := New(d.Hero, d.Weapons, newCards)
			nd.Sideboard = d.Sideboard
			nd.Equipment = d.Equipment
			out = append(out, Mutation{
				Deck:        nd,
				Description: fmt.Sprintf("-1 %s, +1 %s", removed.Name(), replacement.Name()),
			})
		}
	}
	return out
}
