package deck

// Card-pair mutations: swap a synergy pair of cards in or out of the deck as a single atomic
// step so the hill climb can discover combinations whose halves are individually weaker than
// other candidates and would never be added by the regular single-slot generator.
//
// Sun Kiss / Moon Wish is the pilot pairing: Sun Kiss alone is a 3{h}-gain card identical in
// solo value to Healing Balm; Moon Wish alone is a vanilla cost-2 attack with both printed
// riders dropped. Either solo would be cut in seconds by a single-slot mutation, so the hill
// climb never gets the chance to keep both around long enough to evaluate the synergy.
// cardPairMutations forces the question by emitting "-X -Y, +A +B" candidates whenever a pair
// is fully absent.

import (
	"fmt"
	"sort"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
)

// CardPair is two card printings that an anneal mutation must add or remove together because
// either alone has no realised value but the combination unlocks a hidden rider. First and
// Second are concrete printing IDs; we add the printings the registry names and let the
// regular per-pitch single-slot mutations re-tune colours afterwards.
type CardPair struct {
	First  card.ID
	Second card.ID
}

// cardPairs is the registry of synergy pairs the anneal mutation generator considers as
// units. Order matters for deterministic mutation output: pairs are emitted in registry
// order. Variant choice favours the Red printing (lowest pitch cost) so the pair lands at
// minimum opportunity cost; the hill climb's normal per-pitch swaps can re-tune from there.
var cardPairs = []CardPair{
	{First: card.MoonWishRed, Second: card.SunKissRed},
}

// cardPairTopK caps how many low-avg removal slots each absent pair tries. The pair generator
// emits one mutation per (i, j) pair drawn from the K lowest-avg unique IDs in the deck, so
// total candidates per absent pair is K*(K-1)/2. K=5 gives 10 candidates per pair — enough to
// surface meaningful "remove the deck's two worst slots" options without bloating round size.
const cardPairTopK = 5

// cardPairMutations emits paired add/remove mutations for every entry in cardPairs whose two
// halves are both absent from d.Cards. For each absent pair the generator removes one copy of
// each of two distinct low-avg cards (drawn from the cardPairTopK lowest-avg unique IDs) and
// adds one copy of each pair member.
//
// "Both halves absent" is the only triggering case: when one half is already in deck the
// regular single-slot generator can already propose adding the other. Pair mutations are
// strictly the escape hatch for the "neither in deck" hill-climb cul-de-sac.
//
// legal filters BOTH pair members: a pair where either half is rejected by legal (e.g. a
// banned card) is skipped entirely. Removal targets aren't filtered — same convention as
// cardSwapMutations: a deck that arrived holding a banned card can still have it removed.
//
// Returned decks have zero Stats and share no backing slices with d or each other. No-op
// mutations (resulting deck has the same card multiset as the source) are filtered out as a
// defensive guard; the design above ensures none should arise, but the check is cheap.
func cardPairMutations(d *Deck, maxCopies int, legal func(card.Card) bool) []Mutation {
	if len(cardPairs) == 0 {
		return nil
	}
	counts := map[card.ID]int{}
	for _, c := range d.Cards {
		counts[c.ID()]++
	}
	uniqueIDs := lowestAvgUniqueIDs(d, counts, cardPairTopK)

	srcKey := cardMultisetKey(d.Cards)
	var out []Mutation
	for _, pair := range cardPairs {
		if counts[pair.First] > 0 || counts[pair.Second] > 0 {
			continue // single-slot generator handles half-present states.
		}
		first := cards.Get(pair.First)
		second := cards.Get(pair.Second)
		if legal != nil && (!legal(first) || !legal(second)) {
			continue
		}
		// Both halves enter at +1 copy each; the cap check guards against a future registry
		// where maxCopies < 1 (defensive — Random already panics in that case).
		if counts[pair.First]+1 > maxCopies || counts[pair.Second]+1 > maxCopies {
			continue
		}
		for i := 0; i < len(uniqueIDs); i++ {
			for j := i + 1; j < len(uniqueIDs); j++ {
				removeI, removeJ := uniqueIDs[i], uniqueIDs[j]
				newCards := pairSwapDeck(d.Cards, removeI, removeJ, first, second)
				if cardMultisetKey(newCards) == srcKey {
					continue // defensive no-op guard; should never trigger by construction.
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
