package hand

// Entry points and memoization: package-level Best / BestWithTriggers, the Evaluator that owns
// per-goroutine scratch buffers, and the shared memo keyed on (hero, weapons, sorted hand,
// incoming damage, runechant carryover, arsenal-in ID).

import (
	"sync"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Best returns the optimal TurnSummary for the given hand against an opponent that will attack
// for incomingDamage on their next turn. Equipped weapons may be swung for their Cost if
// resources allow.
//
// Cards partition into five roles:
//   - Pitch: contributes its Pitch value as resource paying for a played card this turn or a
//     Defense Reaction on the opponent's turn.
//   - Attack: consumes Cost on our turn; resolved by calling Card.Play in an order the optimizer
//     chooses. TurnState effects carry forward to later attacks in the same sequence.
//   - Defend: contributes Defense to damage prevented (capped at incomingDamage; excess wasted).
//     Plain blocking is free; Defense Reactions must pay Cost and contribute any Play() damage.
//   - Held: stays in hand for next turn. Contributes nothing this turn.
//   - Arsenal: moves into the arsenal slot at end of turn, or for an arsenal-in card, stays.
//     Contributes nothing this turn.
//
// Pitch resources split across two phases because resources don't carry between turns: attack
// pitches pay for this turn's played cards, defense pitches pay for Defense Reactions on the
// opponent's turn. A card can only be pitched while some unpaid card is on the stack in the
// matching phase, so a hand with no plays in a phase has no Pitch-role cards in that phase; any
// card that can't be legally pitched becomes Held.
//
// Results are memoized on (heroID, sorted weapon IDs, sorted card IDs, incomingDamage,
// runechantCarryover, arsenal-in ID) so repeat evaluations short-circuit. The hand is sorted in
// place into canonical order first; BestLine's hand entries align with that post-sort order.
// Every card in the hand must be registered in package cards or Best panics.
//
// runechantCarryover is the Runechant token count carrying in from the previous turn.
// TurnSummary.LeftoverRunechants is the count at end of the chosen chain; feed it back as the
// next turn's carryover.
//
// arsenalCardIn is the card sitting in the arsenal slot at start of turn (nil if empty). The
// enumerator pulls it in as an extra CardAssignment with restricted role options — Arsenal
// (stay), Attack (any non-DR card), or Defend (Defense Reactions only). Never Pitch or Held. A
// hand card may also take the Arsenal role so long as at most one BestLine entry ends up there.
func Best(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card) TurnSummary {
	return sharedEvaluator.Best(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn)
}

// BestWithTriggers is the package-level counterpart of Evaluator.BestWithTriggers, using the
// shared evaluator. Pass priorAuraTriggers to feed the cross-turn trigger carry into the
// search — those triggers may fire mid-chain (e.g. Malefic Incantation's TriggerAttackAction
// rune) and contribute damage to this turn's Value.
func BestWithTriggers(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card, priorAuraTriggers []card.AuraTrigger) TurnSummary {
	return sharedEvaluator.BestWithTriggers(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn, priorAuraTriggers)
}

// Best is the method form of the package-level Best: same semantics, uses this Evaluator's
// scratch buffers so concurrent goroutines can each hold their own.
func (e *Evaluator) Best(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card) TurnSummary {
	return e.BestWithTriggers(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn, nil)
}

// BestWithTriggers is Best plus an explicit priorAuraTriggers input — the AuraTriggers
// carrying in from the previous turn. Non-empty priorAuraTriggers disables memoization: the
// triggers contain Handler closures that aren't comparable, and the sim mutates trigger
// Count / FiredThisTurn mid-chain. With nil priorAuraTriggers, BestWithTriggers matches
// Best exactly (fully memoable).
func (e *Evaluator) BestWithTriggers(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card, priorAuraTriggers []card.AuraTrigger) TurnSummary {
	// IDs go into a fixed-size stack array to avoid a per-call slice alloc. Hand size is capped
	// at 8 (matches memoKey.cardIDs); larger hands panic out of the inner loops.
	n := len(hand)
	var ids [8]card.ID
	memoable := len(priorAuraTriggers) == 0
	for i, c := range hand {
		ids[i] = c.ID()
		if _, ok := c.(card.NoMemo); ok {
			memoable = false
		}
	}
	if arsenalCardIn != nil {
		if _, ok := arsenalCardIn.(card.NoMemo); ok {
			memoable = false
		}
	}

	sortHandByID(hand, ids[:], n)

	key := makeMemoKey(hero, weapons, &ids, n, incomingDamage, runechantCarryover, arsenalCardIn)
	if memoable {
		memoMu.RLock()
		cached, hit := memo[key]
		memoMu.RUnlock()
		if hit {
			return cached
		}
	}
	result := e.bestUncached(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn, priorAuraTriggers)
	if memoable {
		memoMu.Lock()
		memo[key] = result
		memoMu.Unlock()
	}
	return result
}

// sortHandByID sorts the first n entries of hand and ids in parallel by ascending id, in place.
// Insertion sort — for n ≤ 8 this beats sort.Sort and avoids boxing slices through sort.Interface.
// Canonicalizing the hand order is what lets the memo key collapse permutations onto one entry.
func sortHandByID(hand []card.Card, ids []card.ID, n int) {
	for i := 1; i < n; i++ {
		for j := i; j > 0 && ids[j-1] > ids[j]; j-- {
			ids[j-1], ids[j] = ids[j], ids[j-1]
			hand[j-1], hand[j] = hand[j], hand[j-1]
		}
	}
}

// memoKey is the comparable map key for the shared memo. Hand size is capped at 8 cards. Hero
// ID + weapon IDs live in the key so Evaluators with different (hero, weapons) tuples coexist in
// the memo without a scope wipe — distinct scopes just produce distinct keys. The uint16 hero.ID
// keeps the whole key a fixed-size integer struct; no string hashing per lookup.
type memoKey struct {
	heroID             hero.ID
	weaponIDs          [2]card.ID
	cardIDs            [8]card.ID
	cardCount          uint8
	incoming           int
	runechantCarryover int
	// arsenalInID is card.Invalid when the slot is empty, otherwise the ID of the starting
	// arsenal card — different arsenal-ins give distinct cache entries.
	arsenalInID card.ID
}

// Evaluator owns the per-goroutine mutable state hand.Best threads through an evaluation: a
// scratch-buffer cache keyed by (handSize, weapons). The memo cache is shared across all
// Evaluators so every worker benefits from cached hands — only the mutated scratch buffers
// must be per-goroutine. A long-lived Evaluator avoids reallocating ~20 scratch slices per call.
type Evaluator struct {
	// bufs holds the pre-allocated scratch slices for bestUncached / bestAttackWithWeapons /
	// bestSequence. Keyed by (handSize, weaponCount, weaponIDs); recreated when any differ so the
	// scratch sizing stays correct.
	bufs            *attackBufs
	bufsHandSize    int
	bufsWeaponIDs   [2]card.ID
	bufsWeaponCount int
	bufsValid       bool
}

// NewEvaluator returns a fresh Evaluator with empty scratch buffers. One Evaluator per goroutine
// is safe for concurrent use and still shares the global memo.
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// sharedEvaluator backs the package-level Best — single-threaded callers don't need to construct
// their own. Parallel callers create one per goroutine for bufs but all share the global memo.
var sharedEvaluator = NewEvaluator()

// memo caches canonical-order TurnSummary results keyed by memoKey. Shared across all Evaluators
// and goroutines; protected by memoMu. Hero + weapon IDs live in memoKey so distinct (hero,
// weapons) scopes coexist without a wipe step.
var (
	memo   = map[memoKey]TurnSummary{}
	memoMu sync.RWMutex
)

// ClearMemo drops every cached TurnSummary. Callers (iterate-mode, benchmarks) use this to cap
// memo growth across unrelated runs; cross-run hit rate is near zero so nothing of value is lost.
func ClearMemo() {
	memoMu.Lock()
	clear(memo)
	memoMu.Unlock()
}

// MemoLen returns the current number of cached entries, for diagnostic logging.
func MemoLen() int {
	memoMu.RLock()
	n := len(memo)
	memoMu.RUnlock()
	return n
}

// makeMemoKey builds a comparable memo key. The hand must already be sorted by card ID.
// sortedIDs is a pointer to the caller's stack [8]card.ID to avoid a slice-header escape. Weapon
// IDs are sorted into the two fixed slots so loadouts in any order hash to the same key.
func makeMemoKey(hero hero.Hero, weapons []weapon.Weapon, sortedIDs *[8]card.ID, n int, incoming int, runechantCarryover int, arsenalCardIn card.Card) memoKey {
	k := memoKey{
		heroID:             hero.ID(),
		incoming:           incoming,
		runechantCarryover: runechantCarryover,
		cardCount:          uint8(n),
		cardIDs:            *sortedIDs,
	}
	if arsenalCardIn != nil {
		k.arsenalInID = arsenalCardIn.ID()
	}
	switch len(weapons) {
	case 1:
		k.weaponIDs[0] = weapons[0].ID()
	case 2:
		a, b := weapons[0].ID(), weapons[1].ID()
		if a > b {
			a, b = b, a
		}
		k.weaponIDs[0], k.weaponIDs[1] = a, b
	}
	return k
}
