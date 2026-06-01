package deck

// Mutation enumeration for the anneal hill climb. Emits every alternative weapon loadout,
// every (remove one, add one) single-card swap, and every (-1/-1, +1/+1) synergy pair swap
// a deck admits. Order is ascending ids.CardID for stability; the anneal driver shuffles
// the returned slice each round to debias exploration.
//
// CARD PAIRS: a pair is two CardGroups (variant lists). The pair generator enumerates the
// cross-product of (deck-index pair) × (firstVariant, secondVariant), so duplicates at
// distinct positions can both be removed in one mutation. Pair mutations express the
// "atomic 2-for-2 swap" the single-slot generator can't, the escape hatch for synergies
// whose halves are individually weaker than competitors. See CardPairs.

import (
	"fmt"
	"sort"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// mutationKind tags Mutation's tagged-union body so Deck can dispatch to the right
// materialiser. Values are package-private; production code never inspects them.
type mutationKind uint8

const (
	kindSingleSwap mutationKind = iota
	kindPairSwap
	kindWeaponLoadout
)

// Mutation is one candidate single-slot, pair, or weapon-loadout change: a tagged-union spec
// the consumer materialises and labels on demand. The anneal driver shuffles thousands of
// mutations per round and typically accepts within the first handful, so eagerly building
// either the deck or the human-readable summary would burn work per spec that almost always
// goes unused. Deck and Description both run only when actually needed.
type Mutation struct {
	base     *Deck
	kind     mutationKind
	addA     Card
	addB     Card
	removeA  Card
	removeID ids.CardID
	iA       int
	iB       int
	weapons  []Weapon
}

// Description renders the human-readable summary of this mutation (e.g.
// "-1 Aether Slash [R], +1 Arcanic Spike [R]"). Built on demand — the anneal driver only
// reads it for the one mutation per round that gets accepted, so 99% of the per-round
// Sprintf work the eager field used to do is avoided.
func (m Mutation) Description() string {
	switch m.kind {
	case kindSingleSwap:
		return fmt.Sprintf("-1 %s, +1 %s", m.removeA.DisplayName(), m.addA.DisplayName())
	case kindPairSwap:
		return fmt.Sprintf("-1 %s, -1 %s, +1 %s, +1 %s",
			m.base.cards[m.iA].DisplayName(), m.base.cards[m.iB].DisplayName(),
			m.addA.DisplayName(), m.addB.DisplayName())
	case kindWeaponLoadout:
		return fmt.Sprintf("swapped weapons from %s to %s",
			loadoutLabel(m.base.Weapons), loadoutLabel(m.weapons))
	}
	panic(fmt.Sprintf("deck: unknown mutationKind %d", m.kind))
}

// Swap returns the (added, removed) card IDs of a single-card-swap mutation — the pool card being
// tried and the deck card it replaces. ok is false for weapon-loadout and pair mutations, which
// have no single added/removed pair, so the ranking search skips them.
func (m Mutation) Swap() (added, removed ids.CardID, ok bool) {
	if m.kind != kindSingleSwap {
		return 0, 0, false
	}
	return m.addA.ID(), m.removeID, true
}

// Deck builds and returns the mutated deck this spec describes. Each call constructs a fresh
// *Deck, so the caller can shuffle / mutate the result without aliasing the spec or any
// sibling materialisation.
func (m Mutation) Deck() *Deck {
	switch m.kind {
	case kindSingleSwap:
		return m.materialiseSingleSwap()
	case kindPairSwap:
		return m.materialisePairSwap()
	case kindWeaponLoadout:
		return m.materialiseWeaponLoadout()
	}
	panic(fmt.Sprintf("deck: unknown mutationKind %d", m.kind))
}

// materialiseSingleSwap puts the added card in the removed card's slot, keeping every other
// position fixed. Positional alignment matters: the anneal evaluator shuffles a mutant and its
// incumbent from the same seed (common random numbers), so an array differing at exactly one
// index yields shuffles differing at exactly one drawn position. Shifting later slots would
// break that coupling.
func (m Mutation) materialiseSingleSwap() *Deck {
	newCards := make([]Card, len(m.base.cards))
	replaced := false
	for i, c := range m.base.cards {
		if !replaced && c.ID() == m.removeID {
			newCards[i] = m.addA
			replaced = true
			continue
		}
		newCards[i] = c
	}
	nd := New(m.base.Hero, m.base.Weapons, newCards)
	nd.Format = m.base.Format
	nd.Sideboard = m.base.Sideboard
	nd.Equipment = m.base.Equipment
	return nd
}

func (m Mutation) materialisePairSwap() *Deck {
	newCards := pairSwapByIndex(m.base.cards, m.iA, m.iB, m.addA, m.addB)
	nd := New(m.base.Hero, m.base.Weapons, newCards)
	nd.Format = m.base.Format
	nd.Sideboard = m.base.Sideboard
	nd.Equipment = m.base.Equipment
	return nd
}

func (m Mutation) materialiseWeaponLoadout() *Deck {
	newCards := append([]Card(nil), m.base.cards...)
	nd := New(m.base.Hero, m.weapons, newCards)
	nd.Format = m.base.Format
	nd.Sideboard = m.base.Sideboard
	nd.Equipment = m.base.Equipment
	return nd
}

// AllMutations returns every single-card, weapon-loadout, and pair mutation of d in a
// deterministic order: alternative weapon loadouts (sorted), then (removeID, addID) single
// swaps, then synergy pair swaps. removeID must be in the deck; removeID==addID is skipped.
// Ascending-CardID ordering for stability — anneal shuffles the result.
//
// Single-card swaps let the climber reach odd per-card counts (1×X + 3×Y at maxCopies=3).
// The pair layer is the orthogonal escape hatch for synergies whose halves are weaker than
// competitors individually.
//
// reg supplies the legal pool; banned cards in the deck can still be swapped out. maxCopies
// is enforced inline by the single-swap and pair generators (a weapon-loadout change leaves
// card counts untouched, so it's always cap-safe). includePairs gates the pair layer.
func AllMutations(d *Deck, maxCopies int, includePairs bool, reg Registry) []Mutation {
	pool := buildLegalByID(d, reg)
	baseCounts := cardCounts(d.cards)
	out := weaponLoadoutMutations(d, reg)
	out = append(out, singleSwapMutations(d, pool, baseCounts, maxCopies)...)
	if includePairs {
		out = append(out, pairSwapMutations(d, CardPairs, pool, baseCounts, maxCopies)...)
	}
	return out
}

// cardCounts tallies cs by CardID. AllMutations threads the result into the single-swap and
// pair generators so each candidate's inline maxCopies check is O(1).
func cardCounts(cs []Card) map[ids.CardID]int {
	out := make(map[ids.CardID]int, len(cs))
	for _, c := range cs {
		out[c.ID()]++
	}
	return out
}

// buildLegalByID materialises the registry's legal pool (for d's format + hero) as an ID→Card
// map plus IDs in ascending order. The map gives generators O(1) "is this swap-in eligible"
// checks.
func buildLegalByID(d *Deck, reg Registry) legalCardPool {
	cards := reg.LegalCardsFor(d.Format, d.Hero)
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

// legalCardPool is the materialised view of the registry's legal pool. ids is the
// ascending-ID list for stable enumeration; byID resolves an ID for the swap-in side.
type legalCardPool struct {
	byID map[ids.CardID]Card
	ids  []ids.CardID
}

// weaponLoadoutMutations emits one Mutation per distinct weapon loadout that isn't the
// current one. Loadouts are canonicalised by weaponKey (names sorted) and processed in key
// order so the output is deterministic regardless of map-iteration randomness.
func weaponLoadoutMutations(d *Deck, reg Registry) []Mutation {
	loadouts := weaponLoadouts(reg.LegalWeaponsFor(d.Format, d.Hero))
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
		out = append(out, Mutation{
			base:    d,
			kind:    kindWeaponLoadout,
			weapons: l.weapons,
		})
	}
	return out
}

// singleSwapMutations emits every single-card remove+add mutation the deck admits. Remove
// targets iterate in ascending ids.CardID for stability (no value-based bias; the anneal
// driver shuffles afterward). Add candidates skip no-ops (same ID); the maxCopies cap is
// enforced by filterMaxCopiesViolations downstream so this generator stays cap-blind.
func singleSwapMutations(d *Deck, pool legalCardPool, baseCounts map[ids.CardID]int, maxCopies int) []Mutation {
	uniqueRemovals := uniqueDeckCards(d)
	var out []Mutation
	for _, removed := range uniqueRemovals {
		removeID := removed.ID()
		for _, addID := range pool.ids {
			if addID == removeID {
				continue // no-op: remove one and add one of the same card.
			}
			replacement := pool.byID[addID]
			// Inline maxCopies check: only addID's count grows (+1). Other counts unchanged or
			// reduced, so they can't cross the cap.
			if baseCounts[addID]+1 > effectiveMaxCopies(replacement, maxCopies) {
				continue
			}
			out = append(out, Mutation{
				base:     d,
				kind:     kindSingleSwap,
				addA:     replacement,
				removeA:  removed,
				removeID: removeID,
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
func pairSwapMutations(d *Deck, pairs []CardPair, pool legalCardPool, baseCounts map[ids.CardID]int, maxCopies int) []Mutation {
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
				// Inline maxCopies check: firstID and secondID each grow by 1; if the IDs
				// collide (a self-pair), the same card grows by 2.
				if firstID == secondID {
					if baseCounts[firstID]+2 > effectiveMaxCopies(first, maxCopies) {
						continue
					}
				} else {
					if baseCounts[firstID]+1 > effectiveMaxCopies(first, maxCopies) {
						continue
					}
					if baseCounts[secondID]+1 > effectiveMaxCopies(second, maxCopies) {
						continue
					}
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
						out = append(out, Mutation{
							base: d,
							kind: kindPairSwap,
							iA:   i,
							iB:   j,
							addA: first,
							addB: second,
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

// respectsMaxCopies reports whether every distinct ID in cs appears at most maxCopies
// times. Returns false at the first overshoot so a single hot card short-circuits the
// count. Exported for tests pinning the short-circuit behaviour.
func respectsMaxCopies(cs []Card, maxCopies int) bool {
	return respectsMaxCopiesInto(cs, maxCopies, map[ids.CardID]int{})
}

// respectsMaxCopiesInto is respectsMaxCopies with a caller-supplied counts map so a tight
// filtering loop can reuse one map instead of allocating per check. counts is cleared on
// entry; its contents on return are unspecified.
func respectsMaxCopiesInto(cs []Card, maxCopies int, counts map[ids.CardID]int) bool {
	clear(counts)
	for _, c := range cs {
		counts[c.ID()]++
		if counts[c.ID()] > effectiveMaxCopies(c, maxCopies) {
			return false
		}
	}
	return true
}
