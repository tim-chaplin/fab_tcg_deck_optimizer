package deck

// Card-pair mutations: swap a synergy pair of cards in or out of the deck as a single atomic
// step so the hill climb can discover combinations whose halves are individually weaker than
// other candidates and would never be added by the regular single-slot generator.
//
// A pair is two CardGroups (variant lists). The generator enumerates every (firstVariant,
// secondVariant) cross-product so the climber can try Red/Red, Red/Yellow, Red/Blue, … and
// let the regular per-pitch single-slot mutations re-tune from there. Single-slot remains the
// primary mutation source; pair mutations add the orthogonal "atomic 2-for-2 swap" the
// single-slot generator can't express.
//
// Sun Kiss / Moon Wish is the pilot pairing: the synergy reads any Moon Wish printing in
// CardsPlayed by name prefix, so any (Moon Wish variant, Sun Kiss variant) combination is a
// legal pair entry; we register both card groups and let the cross-product enumeration cover
// all 9 variant pairings.

import (
	"fmt"
	"sort"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
)

// CardGroup is a set of card IDs that share a printed name (i.e. all pitch variants of one
// card). Pair mutations enumerate cross-products across two groups so every variant
// combination becomes its own mutation candidate.
type CardGroup []card.ID

// CardPair is two card groups whose members an anneal mutation should add together because
// either alone has weak realised value but the combination unlocks a hidden rider.
type CardPair struct {
	First  CardGroup
	Second CardGroup
}

// Card-group definitions for registered pairs. Splitting them out as named vars keeps the
// cardPairs registry compact and lets future pairs reuse a group on both sides if needed
// (e.g. a self-pair like "two copies of X").
var (
	moonWishGroup = CardGroup{card.MoonWishRed, card.MoonWishYellow, card.MoonWishBlue}
	sunKissGroup  = CardGroup{card.SunKissRed, card.SunKissYellow, card.SunKissBlue}
)

// cardPairs is the registry of synergy pairs the anneal mutation generator considers as
// units. Order matters for deterministic mutation output: pairs are emitted in registry
// order, with cross-product (firstVariant, secondVariant) enumeration in group-slice order.
var cardPairs = []CardPair{
	{First: moonWishGroup, Second: sunKissGroup},
}

// cardPairTopK caps how many low-avg removal slots each pair-variant combo tries. The pair
// generator emits one mutation per (i, j) drawn from the K lowest-avg unique IDs in the
// deck, so total candidates per (firstVariant, secondVariant) combo is K*(K-1)/2. K=5 gives
// 10 candidates per combo — enough to surface meaningful "drop the deck's two worst slots"
// options without bloating round size: a Moon Wish × Sun Kiss pair with all 9 variant combos
// then contributes ≤90 candidates per round, dwarfed by the thousands of single-slot
// mutations.
const cardPairTopK = 5

// cardPairMutations emits paired add mutations for every entry in cardPairs. For each pair
// the generator iterates (firstVariant, secondVariant) cross-products from the two groups
// and, for each combo, emits one mutation per (i, j) pair drawn from the cardPairTopK
// lowest-avg unique IDs in the deck. Each emitted mutation removes one copy of each removal
// target and adds one copy of each pair variant.
//
// Per-variant maxCopies cap: a variant whose count would exceed maxCopies after the +1 add
// is skipped. A deck saturated with Red on both halves yields Yellow/Blue cross-add
// candidates only.
//
// Overlap suppression: when a removal ID matches one of the pair's add IDs, the mutation
// reduces to a single-slot swap (the matching pair member's count is unchanged net of the
// removal). The single-slot generator already covers that, so we skip those combos to keep
// the pair generator strictly orthogonal.
//
// legal filters BOTH pair-variant adds: a combo where either variant is rejected by legal
// (e.g. a banned printing) is skipped. Removal targets aren't filtered — same convention as
// cardSwapMutations: a deck that arrived holding a banned card can still have it removed.
//
// Returned decks have zero Stats and share no backing slices with d or each other. No-op
// mutations (resulting deck has the same card multiset as the source) are filtered out as a
// defensive guard; the overlap suppression above ensures none should arise.
func cardPairMutations(d *Deck, maxCopies int, legal func(card.Card) bool) []Mutation {
	if len(cardPairs) == 0 {
		return nil
	}
	counts := map[card.ID]int{}
	for _, c := range d.Cards {
		counts[c.ID()]++
	}
	uniqueIDs := lowestAvgUniqueIDs(d, counts, cardPairTopK)
	if len(uniqueIDs) < 2 {
		return nil // need at least two distinct removal targets to emit a 2-for-2 swap.
	}

	srcKey := cardMultisetKey(d.Cards)
	var out []Mutation
	for _, pair := range cardPairs {
		out = append(out, mutationsForPair(d, pair, counts, uniqueIDs, srcKey, maxCopies, legal)...)
	}
	return out
}

// mutationsForPair enumerates the (firstVariant, secondVariant) × (removeI, removeJ) candidates
// for one pair entry. Split out so cardPairMutations stays a simple per-pair fan-out and the
// per-pair loops are testable in isolation.
func mutationsForPair(d *Deck, pair CardPair, counts map[card.ID]int, uniqueIDs []card.ID,
	srcKey string, maxCopies int, legal func(card.Card) bool) []Mutation {
	var out []Mutation
	for _, firstID := range pair.First {
		first := cards.Get(firstID)
		if legal != nil && !legal(first) {
			continue
		}
		if counts[firstID]+1 > maxCopies {
			continue
		}
		for _, secondID := range pair.Second {
			second := cards.Get(secondID)
			if legal != nil && !legal(second) {
				continue
			}
			if counts[secondID]+1 > maxCopies {
				continue
			}
			for i := 0; i < len(uniqueIDs); i++ {
				for j := i + 1; j < len(uniqueIDs); j++ {
					removeI, removeJ := uniqueIDs[i], uniqueIDs[j]
					if removeI == firstID || removeI == secondID ||
						removeJ == firstID || removeJ == secondID {
						continue // single-slot generator covers the overlap case.
					}
					newCards := pairSwapDeck(d.Cards, removeI, removeJ, first, second)
					if cardMultisetKey(newCards) == srcKey {
						continue // defensive no-op guard.
					}
					nd := New(d.Hero, d.Weapons, newCards)
					nd.Sideboard = d.Sideboard
					nd.Equipment = d.Equipment
					out = append(out, Mutation{
						Deck: nd,
						Description: fmt.Sprintf("-1 %s, -1 %s, +1 %s, +1 %s",
							cards.Get(removeI).Name(), cards.Get(removeJ).Name(),
							first.Name(), second.Name()),
					})
				}
			}
		}
	}
	return out
}

// lowestAvgUniqueIDs returns up to k unique card IDs from d's card list, ordered by ascending
// per-card avg contribution (weakest first). Ties fall through to ascending card.ID — same
// tiebreak rule as cardSwapMutations so removal-slot ordering stays consistent across the
// single-slot and pair generators.
func lowestAvgUniqueIDs(d *Deck, counts map[card.ID]int, k int) []card.ID {
	ids := make([]card.ID, 0, len(counts))
	for id := range counts {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		ai := d.Stats.PerCard[ids[i]].Avg()
		aj := d.Stats.PerCard[ids[j]].Avg()
		if ai != aj {
			return ai < aj
		}
		return ids[i] < ids[j]
	})
	if len(ids) > k {
		ids = ids[:k]
	}
	return ids
}

// pairSwapDeck builds the post-mutation card slice: copies d.Cards with one removeI removed,
// one removeJ removed, then first and second appended. Removal is one-shot per ID so a deck
// holding multiple copies keeps the rest. The returned slice shares no backing storage with
// d.
func pairSwapDeck(src []card.Card, removeI, removeJ card.ID, first, second card.Card) []card.Card {
	out := make([]card.Card, 0, len(src))
	removedI, removedJ := false, false
	for _, c := range src {
		if !removedI && c.ID() == removeI {
			removedI = true
			continue
		}
		if !removedJ && c.ID() == removeJ {
			removedJ = true
			continue
		}
		out = append(out, c)
	}
	out = append(out, first, second)
	return out
}

// cardMultisetKey returns a comparable string summarising a card slice's ID histogram. Two
// slices with the same IDs in different orders produce equal keys; used by the no-op guard
// to detect mutations whose result equals the source by composition.
func cardMultisetKey(cs []card.Card) string {
	counts := map[card.ID]int{}
	for _, c := range cs {
		counts[c.ID()]++
	}
	ids := make([]int, 0, len(counts))
	for id := range counts {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)
	var b []byte
	for _, id := range ids {
		b = append(b, fmt.Sprintf("%d:%d,", id, counts[card.ID(id)])...)
	}
	return string(b)
}
