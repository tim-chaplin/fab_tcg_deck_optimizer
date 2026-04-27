package hand

// Pre-allocated scratch buffers threaded through the attack-evaluation pipeline (bestUncached
// partition loop, bestAttackWithWeapons phase/weapon masks, bestSequence permutation search).
// Pooled on the Evaluator so one sizing amortises across every hand a long-running iterate pass
// evaluates.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// attackBufs holds pre-allocated buffers for the attack-evaluation pipeline (bestSequence →
// playSequence) and the partition loop in bestUncached. Allocated once and cached on the
// Evaluator so a deck eval reuses them across every partition, mask, and permutation.
type attackBufs struct {
	pcBuf  []card.CardState
	ptrBuf []*card.CardState
	state  *card.TurnState
	// drScratch is a pooled TurnState for defense-reaction cost probing inside the
	// (pmask × wmask) loop; reusing its heap slot avoids a per-iteration alloc caused by
	// interface-call escape.
	drScratch card.TurnState
	// drCardStateScratch is a pooled *CardState handed to DR Card.Play calls. Each Play takes
	// a *CardState through an interface boundary so a literal &card.CardState{} would escape
	// and heap-alloc once per DR per partition — reusing this slot keeps the whole defense-phase
	// replay allocation-free. Reset per call by the caller.
	drCardStateScratch card.CardState
	attackerBuf        []card.Card // for bestAttackWithWeapons mask iteration
	// Pre-computed per-mask weapon data. Indexed by bitmask (0 to 2^len(weapons)-1):
	// weaponCosts[mask] is total Cost; weaponNames[mask] is the pre-built []string of names.
	weaponCosts []int
	weaponNames [][]string
	// permMeta parallels pcBuf: each entry points into the global cardMetaCache so playSequence's
	// inner loop skips interface dispatch on Types / GoAgain and reads cached cost bounds.
	// Pointer-valued so bestSequence's permutation swaps move 8 bytes instead of a full struct.
	permMeta []*attackerMeta
	// Partition-loop buffers, consumed by bestUncached. Sized handSize+1 to cover the optional
	// arsenal-in slot the enumerator treats as index n. isDRBuf caches TypeDefenseReaction
	// membership to skip Types().Has calls; addsFutureValueBuf caches
	// card.AddsFutureValue implementation so the beatsBest tiebreaker can count how many
	// hidden-future-value cards a partition queues.
	rolesBuf           []Role
	pitchVals          []int
	defenseVals        []int
	isDRBuf            []bool
	addsFutureValueBuf []bool
	// pitchedValsScratch backs the per-leaf "pitched values" slice consumed by phase-mask
	// enumeration. Re-sliced to [:0] at the start of every leaf to eliminate a per-leaf alloc.
	pitchedValsScratch []int
	pitchedBuf         []card.Card
	attackersBuf       []card.Card
	defendersBuf       []card.Card
	heldBuf            []card.Card
	// defenseGravScratch backs state.Graveyard during DR Plays. Reset via [:0]+append per
	// iteration so card effects can freely mutate their view without leaking into the next one.
	defenseGravScratch []card.Card
	// Per-permutation backing slices reused across every Heap's-algorithm permutation in a
	// leaf. resetStateForPermutation seeds the TurnState's slice fields from these (via
	// append([:0], ...)) so an unmodified permutation never reallocates: only mid-chain
	// growth past the pre-sized cap forces a new backing array. snapshotCarry clones the
	// winning permutation's slices before the next permutation overwrites them.
	deckBacking         []card.Card
	handBacking         []card.Card
	graveBacking        []card.Card
	banishBacking       []card.Card
	cardsPlayedBacking  []card.Card
	logBacking          []card.LogEntry
	auraTriggersBacking []card.AuraTrigger
	ephemeralBacking    []card.EphemeralAttackTrigger
}

func newAttackBufs(handSize, weaponCount int, weapons []weapon.Weapon) *attackBufs {
	// +1 reserves a slot for the arsenal-in card, which joins attackers or defenders when the
	// enumerator plays it from arsenal. +maxDrawnExtensions leaves headroom for mid-turn-drawn
	// cards that play as chain extensions — cheap cycling cards (cost 0, Go again, draws a
	// card) can extend a chain well past the starting hand size.
	const maxDrawnExtensions = 32
	maxAttackers := handSize + weaponCount + 1 + maxDrawnExtensions
	numMasks := 1 << weaponCount
	weaponCosts := make([]int, numMasks)
	weaponNames := make([][]string, numMasks)
	for mask := 0; mask < numMasks; mask++ {
		cost := 0
		var names []string
		for i, w := range weapons {
			if mask&(1<<i) != 0 {
				cost += w.Cost(&card.TurnState{})
				names = append(names, card.DisplayName(w))
			}
		}
		weaponCosts[mask] = cost
		weaponNames[mask] = names
	}
	pcBuf := make([]card.CardState, maxAttackers)
	ptrBuf := make([]*card.CardState, maxAttackers)
	// Wire the ptrBuf entries to their pcBuf slots once — the mapping is stable across every
	// permutation so playSequenceWithMeta doesn't need to rewrite it per call.
	for i := range pcBuf {
		ptrBuf[i] = &pcBuf[i]
	}
	// Per-permutation backing capacities. Decks run ~40 cards but mid-turn draws and tutors can
	// grow them, so size to a safe headroom. Log accumulates one entry per chain step + every
	// rider/trigger (typically 2-6 per step) so 64 covers the long tail.
	const (
		deckBackingCap = 64
		logBackingCap  = 64
	)
	return &attackBufs{
		permMeta:            make([]*attackerMeta, maxAttackers),
		pcBuf:               pcBuf,
		ptrBuf:              ptrBuf,
		state:               &card.TurnState{},
		attackerBuf:         make([]card.Card, maxAttackers),
		weaponCosts:         weaponCosts,
		weaponNames:         weaponNames,
		rolesBuf:            make([]Role, handSize+1),
		pitchVals:           make([]int, handSize+1),
		defenseVals:         make([]int, handSize+1),
		isDRBuf:             make([]bool, handSize+1),
		addsFutureValueBuf:  make([]bool, handSize+1),
		pitchedValsScratch:  make([]int, 0, handSize+1),
		pitchedBuf:          make([]card.Card, 0, handSize+1),
		attackersBuf:        make([]card.Card, 0, handSize+1),
		defendersBuf:        make([]card.Card, 0, handSize+1),
		heldBuf:             make([]card.Card, 0, handSize+1),
		defenseGravScratch:  make([]card.Card, 0, handSize+1),
		deckBacking:         make([]card.Card, 0, deckBackingCap),
		handBacking:         make([]card.Card, 0, maxAttackers),
		graveBacking:        make([]card.Card, 0, maxAttackers),
		banishBacking:       make([]card.Card, 0, handSize+1),
		cardsPlayedBacking:  make([]card.Card, 0, maxAttackers),
		logBacking:          make([]card.LogEntry, 0, logBackingCap),
		auraTriggersBacking: make([]card.AuraTrigger, 0, handSize+1),
		// Ephemeral attack triggers (Mauvrion Skies, Runic Reaping) typically register one per
		// applicable card — pre-sized cap avoids the per-Play slice grow.
		ephemeralBacking: make([]card.EphemeralAttackTrigger, 0, handSize+1),
	}
}

// getAttackBufs returns the Evaluator's cached attackBufs when (handSize, weapons) match the
// last call; otherwise allocates a fresh one and caches it. weapons are zero-size structs so
// interface equality is a stable identity check. For a single deck eval (10k shuffles, same
// hand size + same weapons) this allocates once and reuses on every subsequent call —
// attackBufs is the second-biggest allocator after the eval-time slice copies.
func (e *Evaluator) getAttackBufs(handSize int, weapons []weapon.Weapon) *attackBufs {
	if e.cachedBufs != nil && e.cachedHandSize == handSize && sameWeapons(e.cachedWeapons, weapons) {
		return e.cachedBufs
	}
	e.cachedBufs = newAttackBufs(handSize, len(weapons), weapons)
	e.cachedHandSize = handSize
	// Snapshot the weapons slice header — caller may reuse the slice across calls. The
	// underlying weapon.Weapon values are zero-size structs; interface equality compares
	// (type, nil-data) tuples that are stable across calls.
	e.cachedWeapons = append(e.cachedWeapons[:0], weapons...)
	return e.cachedBufs
}

// sameWeapons reports whether two weapon slices contain the same weapons in the same order.
// Element-wise interface equality works because every weapon implementation is a zero-size
// struct, making interface values comparable and stable across calls.
func sameWeapons(a, b []weapon.Weapon) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// fillPartitionPerCardBufs writes the per-card values the partition recurse reads at each leaf:
// Pitch / Defense magnitudes, Defense-Reaction membership, and AddsFutureValue interface
// satisfaction. Computing them up front keeps the recurse's inner body free of card-method /
// type-assert calls, which would otherwise repeat on every leaf. totalN covers the optional
// arsenal-in slot at index n; when present, its Defense picks up ArsenalDefenseBonus so the
// partition / capping pipeline sees the effective value. Returns whether any card is a
// Defense Reaction so the leaf branch can pick between the full three-bucket grouper and the
// faster reaction-free grouper.
func fillPartitionPerCardBufs(hand []card.Card, n, totalN int, arsenalCardIn card.Card, pvals, dvals []int, isDR, addsFutureValue []bool) bool {
	hasReactions := false
	for i := 0; i < totalN; i++ {
		var c card.Card
		if i < n {
			c = hand[i]
		} else {
			c = arsenalCardIn
		}
		pvals[i] = c.Pitch()
		dvals[i] = c.Defense()
		// Arsenal slot (i == n) lives at the end. Defense Reactions whose +N{d} rider only fires
		// when played from arsenal (Unmovable, Springboard Somersault) opt in via
		// card.ArsenalDefenseBonus; bump the static Defense() up here so the partition / capping
		// pipeline sees the effective value.
		if i == n {
			if ab, ok := c.(card.ArsenalDefenseBonus); ok {
				dvals[i] += ab.ArsenalDefenseBonus()
			}
		}
		isDR[i] = c.Types().IsDefenseReaction()
		if isDR[i] {
			hasReactions = true
		}
		_, addsFutureValue[i] = c.(card.AddsFutureValue)
	}
	return hasReactions
}
