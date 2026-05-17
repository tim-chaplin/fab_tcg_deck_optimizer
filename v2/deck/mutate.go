package deck

// Mutation enumeration for the anneal-mode hill climb. Every alternative weapon loadout,
// every (remove one, add one) single-card swap, and every paired (-1/-1, +1/+1) synergy
// swap a deck admits. Ordering is by ascending ids.CardID for stability — no value-based
// bias. The anneal driver shuffles the returned slice each round so exploration order is
// unbiased.
//
// CARD PAIRS:
// A pair is two CardGroups (variant lists). The pair generator enumerates the cross-product
// of (deck-index pair) × (firstVariant, secondVariant) — every position pair in the deck
// times every variant combination — so duplicate cards at distinct positions can both be
// removed in one mutation (e.g. a deck of [HocusPocusBlue, HocusPocusBlue] swapping both
// copies for a Sun Kiss / Moon Wish pair). Single-slot remains the primary mutation source;
// pair mutations add the orthogonal "atomic 2-for-2 swap" the single-slot generator can't
// express.
//
// Sun Kiss / Moon Wish is the pilot pairing: the synergy reads any Moon Wish printing in
// CardsPlayed by name prefix, so any (Moon Wish variant, Sun Kiss variant) combination is
// a legal pair entry; we register both card groups and let the cross-product enumeration
// cover all 9 variant pairings.

import (
	"fmt"
	"sort"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// Mutation is one candidate single-slot or pair change: the mutated Deck plus a
// human-readable summary (e.g. "swapped Aether Slash [R] for Arcanic Spike [R]"). Consumers
// use Deck to evaluate and Description for logging.
type Mutation struct {
	Deck        *Deck
	Description string
}

// AllMutations returns every single-card and pair mutation of d in a deterministic order:
// first every alternative weapon loadout (sorted by loadout key), then every (removeID,
// addID) pair where one copy of removeID is dropped and one copy of addID is added, then
// the synergy-pair "swap two for two" mutations. removeID must be in the deck. Pairs with
// removeID == addID are skipped.
//
// Card-mutation ordering is by ascending ids.CardID for stability — no value-based bias.
// The anneal driver shuffles the returned slice each round so neither the first-found
// classical climb nor the probabilistic SA gate disproportionately samples the head of the
// slice.
//
// Single-card swaps (not paired swaps) let the hill climber reach decks with odd per-card
// counts (e.g. 1× X + 3× Y at maxCopies=3). The pair-swap layer is the orthogonal escape
// hatch for synergies whose halves are individually weaker than competitors and would never
// enter the deck via single-slot mutations alone — see CardPairs.
//
// reg supplies the legal card and weapon rosters; legal optionally filters reg's cards
// (format ban check). Removal targets aren't filtered — a deck that entered the climb
// holding a banned card can still have it swapped out. legal=nil disables the format
// filter.
//
// maxCopies is enforced by filterMaxCopiesViolations as a final post-pass over the combined
// candidate list — both single-slot and pair generators emit cap-blind candidates and the
// shared filter strips any whose result deck exceeds the per-printing limit.
//
// Returned decks share no backing slices with d or each other.
func AllMutations(d *Deck, maxCopies int, reg Registry) []Mutation {
	pool := buildLegalByID(reg)
	out := weaponLoadoutMutations(d, reg)
	out = append(out, singleSwapMutations(d, pool)...)
	out = append(out, pairSwapMutations(d, CardPairs, pool)...)
	return filterMaxCopiesViolations(out, maxCopies)
}

// buildLegalByID materialises the registry's legal pool as an ID→Card map plus the IDs
// in ascending order for stable iteration. The map gives the mutation generators an O(1)
// "is this swap-in eligible" check without re-scanning the registry on every candidate.
func buildLegalByID(reg Registry) legalCardPool {
	cards := reg.LegalCards()
	pool := legalCardPool{
		byID: make(map[ids.CardID]Card, len(cards)),
		ids:  make([]ids.CardID, 0, len(cards)),
	}
	for _, c := range cards {
		pool.byID[c.ID()] = c
		pool.ids = append(pool.ids, c.ID())
	}
	return pool
}

// legalCardPool is the materialised view of the registry's legal pool the per-mutation
// generators iterate over. ids is the ascending-ID list for stable enumeration; byID
// resolves an ID to its Card for the swap-in side of every emitted mutation.
type legalCardPool struct {
	byID map[ids.CardID]Card
	ids  []ids.CardID
}

// weaponLoadoutMutations emits one Mutation per distinct weapon loadout that isn't the
// current one. Loadouts are canonicalised by weaponKey (names sorted) and processed in key
// order so the output is deterministic regardless of map-iteration randomness.
func weaponLoadoutMutations(d *Deck, reg Registry) []Mutation {
	loadouts := weaponLoadouts(reg.LegalWeapons())
	currentKey := weaponKey(d.Weapons)
	type keyedLoadout struct {
		key     string
		weapons []Weapon
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
		newCards := append([]Card(nil), d.cards...)
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
// targets iterate in ascending ids.CardID for stability (no value-based bias; the anneal
// driver shuffles afterward). Add candidates skip no-ops (same ID); the maxCopies cap is
// enforced by filterMaxCopiesViolations downstream so this generator stays cap-blind.
func singleSwapMutations(d *Deck, pool legalCardPool) []Mutation {
	uniqueRemovals := uniqueDeckCards(d)
	var out []Mutation
	for _, removed := range uniqueRemovals {
		removeID := removed.ID()
		for _, addID := range pool.ids {
			if addID == removeID {
				continue // no-op: remove one and add one of the same card.
			}
			replacement := pool.byID[addID]
			newCards := make([]Card, 0, d.Size())
			removed1 := false
			for _, c := range d.cards {
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
				Description: fmt.Sprintf("-1 %s, +1 %s", removed.DisplayName(), replacement.DisplayName()),
			})
		}
	}
	return out
}

// uniqueDeckCards returns one representative Card per distinct card ID in d, sorted by
// ascending ID. The ordering is for output stability; the anneal driver shuffles the final
// mutation slice so search bias doesn't enter through the iteration order.
func uniqueDeckCards(d *Deck) []Card {
	seen := map[ids.CardID]bool{}
	out := make([]Card, 0, d.Size())
	for _, c := range d.cards {
		if seen[c.ID()] {
			continue
		}
		seen[c.ID()] = true
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// CardGroup is a set of card IDs that share a printed name (i.e. all pitch variants of one
// card). Pair mutations enumerate cross-products across two groups so every variant
// combination becomes its own mutation candidate.
type CardGroup []ids.CardID

// CardPair is two card groups whose members an anneal mutation should add together because
// either alone has weak realised value but the combination unlocks a hidden rider.
type CardPair struct {
	First  CardGroup
	Second CardGroup
}

// Card-group definitions for registered pairs. Splitting them out as named vars keeps the
// CardPairs registry compact and lets future pairs reuse a group on both sides if needed
// (e.g. a self-pair like "two copies of X").
var (
	moonWishGroup     = CardGroup{ids.MoonWishRed, ids.MoonWishYellow, ids.MoonWishBlue}
	sunKissGroup      = CardGroup{ids.SunKissRed, ids.SunKissYellow, ids.SunKissBlue}
	belittleGroup     = CardGroup{ids.BelittleRed, ids.BelittleYellow, ids.BelittleBlue}
	minnowismGroup    = CardGroup{ids.MinnowismRed, ids.MinnowismYellow, ids.MinnowismBlue}
	nimblismGroup     = CardGroup{ids.NimblismRed, ids.NimblismYellow, ids.NimblismBlue}
	sloggismGroup     = CardGroup{ids.SloggismRed, ids.SloggismYellow, ids.SloggismBlue}
	nimbleStrikeGroup = CardGroup{
		ids.NimbleStrikeRed, ids.NimbleStrikeYellow, ids.NimbleStrikeBlue,
	}
	regurgitatingSlogGroup = CardGroup{
		ids.RegurgitatingSlogRed, ids.RegurgitatingSlogYellow, ids.RegurgitatingSlogBlue,
	}
	amuletOfHavencallGroup = CardGroup{ids.AmuletOfHavencallBlue}
	rallyTheRearguardGroup = CardGroup{
		ids.RallyTheRearguardRed, ids.RallyTheRearguardYellow, ids.RallyTheRearguardBlue,
	}
)

// CardPairs is the registry of synergy pairs the anneal mutation generator considers as
// units — every pair where one card's printed text mentions the other by name. Order
// matters for deterministic mutation output: pairs are emitted in registry order, with
// cross-product (firstVariant, secondVariant) enumeration in group-slice order and
// (i, j) deck-index iteration in ascending order.
//
// Pair-variant combos whose halves carry NotImplemented are skipped at enumeration time
// (the registry's legal pool excludes them), so a pair stays valid to register even when
// one half is still a stub — the combo activates once both halves are fully implemented.
var CardPairs = []CardPair{
	{First: moonWishGroup, Second: sunKissGroup},
	{First: belittleGroup, Second: minnowismGroup},
	{First: nimbleStrikeGroup, Second: nimblismGroup},
	{First: regurgitatingSlogGroup, Second: sloggismGroup},
	{First: amuletOfHavencallGroup, Second: rallyTheRearguardGroup},
}

// pairDedupeKey is the canonical (sorted removed IDs, sorted add IDs) tuple identifying a
// pair-mutation candidate by content; see pairSwapMutations for the iteration shape.
type pairDedupeKey struct {
	rmA, rmB   ids.CardID
	addA, addB ids.CardID
}

// pairSwapMutations emits paired add mutations for every entry in CardPairs by taking the
// cross-product of every (i, j) deck-index pair (i < j) with every (firstVariant,
// secondVariant) combo from the pair's two groups. Each emitted mutation removes the cards
// at positions i and j from a fresh copy of d's cards and appends the chosen pair variants.
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
// Pair-variant adds are limited to pool.byID — NotImplemented printings are absent from the
// registry's legal pool, and the caller's legal predicate is honoured via the pool build.
// Removal targets aren't filtered — same convention as singleSwapMutations: a deck that
// arrived holding a banned card can still have it removed.
//
// maxCopies enforcement is NOT applied here; AllMutations runs filterMaxCopiesViolations on
// the combined output so single-slot and pair candidates share one cap-checking pass.
//
// Returned decks share no backing slices with d or each other.
func pairSwapMutations(d *Deck, pairs []CardPair, pool legalCardPool) []Mutation {
	if len(pairs) == 0 || d.Size() < 2 {
		return nil
	}
	seen := map[pairDedupeKey]bool{}
	var out []Mutation
	for _, pair := range pairs {
		for _, firstID := range pair.First {
			first, ok := pool.byID[firstID]
			if !ok {
				continue
			}
			for _, secondID := range pair.Second {
				second, ok := pool.byID[secondID]
				if !ok {
					continue
				}
				addA, addB := sortedIDPair(firstID, secondID)
				for i := 0; i < len(d.cards); i++ {
					for j := i + 1; j < len(d.cards); j++ {
						idI, idJ := d.cards[i].ID(), d.cards[j].ID()
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
						newCards := pairSwapByIndex(d.cards, i, j, first, second)
						nd := New(d.Hero, d.Weapons, newCards)
						nd.Sideboard = d.Sideboard
						nd.Equipment = d.Equipment
						out = append(out, Mutation{
							Deck: nd,
							Description: fmt.Sprintf("-1 %s, -1 %s, +1 %s, +1 %s",
								d.cards[i].DisplayName(), d.cards[j].DisplayName(),
								first.DisplayName(), second.DisplayName()),
						})
					}
				}
			}
		}
	}
	return out
}

// sortedIDPair returns (a, b) sorted ascending so callers can build canonical
// order-independent keys.
func sortedIDPair(a, b ids.CardID) (ids.CardID, ids.CardID) {
	if b < a {
		return b, a
	}
	return a, b
}

// pairSwapByIndex returns a fresh slice equal to src with positions i and j removed and
// first and second appended. i and j must be distinct and in range; callers guarantee this
// via i < j enumeration over a sized loop.
func pairSwapByIndex(src []Card, i, j int, first, second Card) []Card {
	out := make([]Card, 0, len(src))
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
func cardMultisetKey(cs []Card) string {
	counts := map[ids.CardID]int{}
	for _, c := range cs {
		counts[c.ID()]++
	}
	keys := make([]int, 0, len(counts))
	for id := range counts {
		keys = append(keys, int(id))
	}
	sort.Ints(keys)
	var b []byte
	for _, id := range keys {
		b = append(b, fmt.Sprintf("%d:%d,", id, counts[ids.CardID(id)])...)
	}
	return string(b)
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
//
// Production callers go through AllMutations; filterMaxCopiesViolations is exported for
// tests that exercise the cap filter directly.
func filterMaxCopiesViolations(muts []Mutation, maxCopies int) []Mutation {
	out := make([]Mutation, 0, len(muts))
	for _, m := range muts {
		if respectsMaxCopies(m.Deck.cards, maxCopies) {
			out = append(out, m)
		}
	}
	return out
}

// respectsMaxCopies reports whether every distinct ID in cs appears at most maxCopies
// times. Returns false at the first overshoot so a single hot card short-circuits the
// count. Exported for tests pinning the short-circuit behaviour.
func respectsMaxCopies(cs []Card, maxCopies int) bool {
	counts := map[ids.CardID]int{}
	for _, c := range cs {
		counts[c.ID()]++
		if counts[c.ID()] > maxCopies {
			return false
		}
	}
	return true
}
