package deck

// Single-slot mutation generation for the iterate-mode hill climb: every alternative weapon
// loadout plus every (remove one, add one) card swap the deck admits. Ordering is by
// ascending ids.CardID for stability — no value-based bias. The anneal driver shuffles the
// returned slice each round so exploration order is unbiased.

import (
	"fmt"
	"sort"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Mutation is one candidate single-slot change: the mutated Deck plus a human-readable summary
// (e.g. "swapped Aether Slash [R] for Arcanic Spike [R]"). Consumers use Deck to evaluate
// and Description for logging.
type Mutation struct {
	Deck        *Deck
	Description string
}

// AllMutations returns every single-card mutation of d in a deterministic order: first every
// alternative weapon loadout (sorted by loadout key), then every (removeID, addID) pair where
// one copy of removeID is dropped and one copy of addID is added, then the synergy-pair
// "swap two for two" mutations from pairSwapMutations. removeID must be in the deck. Pairs
// with removeID == addID are skipped.
//
// Card-mutation ordering is by ascending ids.CardID for stability — no value-based bias. The
// anneal driver shuffles the returned slice each round so neither the first-found classical
// climb nor the probabilistic SA gate disproportionately samples the head of the slice.
//
// Single-card swaps (not paired swaps) let the hill climber reach decks with odd per-card counts
// (e.g. 1× X + 3× Y at maxCopies=3). The pair-swap layer is the orthogonal escape hatch for
// synergies whose halves are individually weaker than competitors and would never enter the
// deck via single-slot mutations alone — see cardPairs in card_pairs.go.
//
// legal filters the addition pool: only accepted IDs become swap-in candidates, so format-banned
// cards can't be introduced. Removal targets aren't filtered — a deck that entered the climb
// holding a banned card can still have it swapped out. Pass nil to skip filtering.
//
// maxCopies is enforced by filterMaxCopiesViolations as a final post-pass over the combined
// candidate list — both single-slot and pair generators emit cap-blind candidates and the
// shared filter strips any whose result deck exceeds the per-printing limit.
//
// Returned decks have zero Stats and share no backing slices with d or each other.
func AllMutations(d *Deck, maxCopies int, legal func(card.Card) bool) []Mutation {
	out := weaponLoadoutMutations(d)
	out = append(out, singleSwapMutations(d, legal)...)
	out = append(out, pairSwapMutations(d, legal)...)
	return filterMaxCopiesViolations(out, maxCopies)
}

// weaponLoadoutMutations emits one Mutation per distinct weapon loadout that isn't the current
// one. Loadouts are canonicalised by weaponKey (names sorted) and processed in key order so
// the output is deterministic regardless of map-iteration randomness.
func weaponLoadoutMutations(d *Deck) []Mutation {
	loadouts := weaponLoadouts(weapon.All)
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

// singleSwapMutations emits every single-card remove+add mutation the deck admits. Remove
// targets iterate in ascending ids.CardID for stability (no value-based bias; the anneal driver
// shuffles afterward). Add candidates skip no-ops (same ID); the maxCopies cap is enforced
// by filterMaxCopiesViolations downstream so this generator stays cap-blind.
func singleSwapMutations(d *Deck, legal func(card.Card) bool) []Mutation {
	uniqueIDs := sortedDeckIDs(d.Cards)

	// legalPool returns IDs in ascending order (registry.DeckableCards() iterates byID).
	pool := legalPool(legal)

	var out []Mutation
	for _, removeID := range uniqueIDs {
		removed := registry.GetCard(removeID)
		for _, addID := range pool {
			if addID == removeID {
				continue // no-op: remove one and add one of the same card.
			}
			replacement := registry.GetCard(addID)
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

// sortedDeckIDs returns every distinct card ID appearing in cs, sorted ascending. Used as
// the removal-target ordering for singleSwapMutations; the order is purely for stability since
// the anneal driver shuffles the final mutation slice.
func sortedDeckIDs(cs []card.Card) []ids.CardID {
	seen := map[ids.CardID]bool{}
	ids := make([]ids.CardID, 0, len(cs))
	for _, c := range cs {
		id := c.ID()
		if seen[id] {
			continue
		}
		seen[id] = true
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}

// filterMaxCopiesViolations returns a fresh slice holding the subset of muts whose
// post-mutation deck respects the per-printing maxCopies cap. Centralising the cap check
// here keeps the per-mutation generators free to enumerate cap-blind candidates; the shared
// post-pass guarantees no downstream consumer ever sees a candidate that violates the
// construction limit.
//
// Source decks that themselves violate maxCopies (e.g. a hand-curated deck loaded from
// disk) flow through unchanged on weapon-only mutations; only mutations that grow a
// violation strictly worse get filtered.
//
// The returned slice does not share storage with the input — callers can keep the original
// muts slice intact for diagnostics if they want.
func filterMaxCopiesViolations(muts []Mutation, maxCopies int) []Mutation {
	out := make([]Mutation, 0, len(muts))
	for _, m := range muts {
		if respectsMaxCopies(m.Deck.Cards, maxCopies) {
			out = append(out, m)
		}
	}
	return out
}

// respectsMaxCopies reports whether every distinct ID in cs appears at most maxCopies times.
// Returns false at the first overshoot so a single hot card short-circuits the count.
func respectsMaxCopies(cs []card.Card, maxCopies int) bool {
	counts := map[ids.CardID]int{}
	for _, c := range cs {
		counts[c.ID()]++
		if counts[c.ID()] > maxCopies {
			return false
		}
	}
	return true
}
