package deck

// Card-pair mutations: swap a synergy pair of cards in or out of the deck as a single atomic
// step so the hill climb can discover combinations whose halves are individually weaker than
// other candidates and would never be added by the regular single-slot generator.
//
// A pair is two CardGroups (variant lists). The generator enumerates the cross-product of
// (deck-index pair) × (firstVariant, secondVariant) — every position pair in the deck times
// every variant combination — so duplicate cards at distinct positions can both be removed
// in one mutation (e.g. a deck of [HocusPocusBlue, HocusPocusBlue] swapping both copies for
// a Sun Kiss / Moon Wish pair). Single-slot remains the primary mutation source; pair
// mutations add the orthogonal "atomic 2-for-2 swap" the single-slot generator can't express.
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
	moonWishGroup     = CardGroup{card.MoonWishRed, card.MoonWishYellow, card.MoonWishBlue}
	sunKissGroup      = CardGroup{card.SunKissRed, card.SunKissYellow, card.SunKissBlue}
	belittleGroup     = CardGroup{card.BelittleRed, card.BelittleYellow, card.BelittleBlue}
	minnowismGroup    = CardGroup{card.MinnowismRed, card.MinnowismYellow, card.MinnowismBlue}
	nimblismGroup     = CardGroup{card.NimblismRed, card.NimblismYellow, card.NimblismBlue}
	sloggismGroup     = CardGroup{card.SloggismRed, card.SloggismYellow, card.SloggismBlue}
	nimbleStrikeGroup = CardGroup{
		card.NimbleStrikeRed, card.NimbleStrikeYellow, card.NimbleStrikeBlue,
	}
	regurgitatingSlogGroup = CardGroup{
		card.RegurgitatingSlogRed, card.RegurgitatingSlogYellow, card.RegurgitatingSlogBlue,
	}
	amuletOfHavencallGroup = CardGroup{card.AmuletOfHavencallBlue}
	rallyTheRearguardGroup = CardGroup{
		card.RallyTheRearguardRed, card.RallyTheRearguardYellow, card.RallyTheRearguardBlue,
	}
)

// cardPairs is the registry of synergy pairs the anneal mutation generator considers as
// units — every pair where one card's printed text mentions the other by name. Order matters
// for deterministic mutation output: pairs are emitted in registry order, with cross-product
// (firstVariant, secondVariant) enumeration in group-slice order and (i, j) deck-index
// iteration in ascending order.
//
// Pair-variant combos whose halves carry card.NotImplemented are skipped at enumeration
// time (see pairAddAllowed), so a pair stays valid to register even when one half is still
// a stub — the combo activates once both halves are fully implemented.
var cardPairs = []CardPair{
	{First: moonWishGroup, Second: sunKissGroup},
	{First: belittleGroup, Second: minnowismGroup},
	{First: nimbleStrikeGroup, Second: nimblismGroup},
	{First: regurgitatingSlogGroup, Second: sloggismGroup},
	{First: amuletOfHavencallGroup, Second: rallyTheRearguardGroup},
}

// pairDedupeKey is the canonical (sorted removed IDs, sorted add IDs) tuple identifying a
// pair-mutation candidate by content; see pairSwapMutations for the iteration shape.
type pairDedupeKey struct {
	rmA, rmB   card.ID
	addA, addB card.ID
}

// pairSwapMutations emits paired add mutations for every entry in cardPairs by taking the
// cross-product of every (i, j) deck-index pair (i < j) with every (firstVariant,
// secondVariant) combo from the pair's two groups. Each emitted mutation removes the cards
// at positions i and j from a fresh copy of d.Cards and appends the chosen pair variants.
//
// Index-based iteration is what makes "remove both copies of a card that appears twice"
// reachable: with a [HocusPocusBlue, HocusPocusBlue] deck, the (0, 1) index pair removes
// both copies in one mutation. A unique-ID iteration would skip this case (only one unique
// ID, no inter-ID pair).
//
// Overlap suppression: when a removed card's ID equals one of the pair add IDs, the
// mutation reduces to a single-slot swap (the matching pair member's count is unchanged
// after -1 +1). Single-slot already covers that, so we skip those combos to keep the pair
// generator strictly orthogonal — and correctness-wise, this filter is also what guarantees
// pair mutations never produce a no-op (multiset unchanged) deck.
//
// Dedupe: duplicate cards at distinct indices generate the same result deck. We track
// emitted (sorted-removed-IDs, sorted-add-IDs) tuples in pairDedupeKey form and drop
// repeats.
//
// Pair-variant adds run through pairAddAllowed: NotImplemented printings are filtered
// unconditionally and the caller's legal predicate is honoured per-variant. Removal targets
// aren't filtered — same convention as singleSwapMutations: a deck that arrived holding a
// banned card can still have it removed.
//
// maxCopies enforcement is NOT applied here; AllMutations runs filterMaxCopiesViolations on
// the combined output so single-slot and pair candidates share one cap-checking pass.
//
// Returned decks have zero Stats and share no backing slices with d or each other.
func pairSwapMutations(d *Deck, legal func(card.Card) bool) []Mutation {
	if len(cardPairs) == 0 || len(d.Cards) < 2 {
		return nil
	}
	seen := map[pairDedupeKey]bool{}
	var out []Mutation
	for _, pair := range cardPairs {
		for _, firstID := range pair.First {
			first := cards.Get(firstID)
			if !pairAddAllowed(first, legal) {
				continue
			}
			for _, secondID := range pair.Second {
				second := cards.Get(secondID)
				if !pairAddAllowed(second, legal) {
					continue
				}
				addA, addB := sortedIDPair(firstID, secondID)
				for i := 0; i < len(d.Cards); i++ {
					for j := i + 1; j < len(d.Cards); j++ {
						idI, idJ := d.Cards[i].ID(), d.Cards[j].ID()
						if idI == firstID || idI == secondID ||
							idJ == firstID || idJ == secondID {
							continue
						}
						rmA, rmB := sortedIDPair(idI, idJ)
						key := pairDedupeKey{rmA, rmB, addA, addB}
						if seen[key] {
							continue
						}
						seen[key] = true
						newCards := pairSwapByIndex(d.Cards, i, j, first, second)
						nd := New(d.Hero, d.Weapons, newCards)
						nd.Sideboard = d.Sideboard
						nd.Equipment = d.Equipment
						out = append(out, Mutation{
							Deck: nd,
							Description: fmt.Sprintf("-1 %s, -1 %s, +1 %s, +1 %s",
								d.Cards[i].Name(), d.Cards[j].Name(),
								first.Name(), second.Name()),
						})
					}
				}
			}
		}
	}
	return out
}

// pairAddAllowed reports whether c is eligible as a pair-mutation add target. Rejects
// anything carrying card.NotImplemented (the simulator can't model the rider, so introducing
// the printing into a deck would corrupt the search) and applies the caller's legal filter
// when present. legal=nil keeps every implemented card eligible.
func pairAddAllowed(c card.Card, legal func(card.Card) bool) bool {
	if _, unimplemented := c.(card.NotImplemented); unimplemented {
		return false
	}
	if legal != nil && !legal(c) {
		return false
	}
	return true
}

// sortedIDPair returns (a, b) sorted ascending so callers can build canonical
// order-independent keys.
func sortedIDPair(a, b card.ID) (card.ID, card.ID) {
	if b < a {
		return b, a
	}
	return a, b
}

// pairSwapByIndex returns a fresh slice equal to src with positions i and j removed and
// first and second appended. i and j must be distinct and in range; callers guarantee this
// via i < j enumeration over a sized loop.
func pairSwapByIndex(src []card.Card, i, j int, first, second card.Card) []card.Card {
	out := make([]card.Card, 0, len(src))
	for k, c := range src {
		if k == i || k == j {
			continue
		}
		out = append(out, c)
	}
	out = append(out, first, second)
	return out
}

// cardMultisetKey returns a comparable string summarising a card slice's ID histogram. Two
// slices with the same IDs in different orders produce equal keys. Tests use it to assert
// pair mutations never produce a deck whose composition equals the source.
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
