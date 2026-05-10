package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// Pre-allocated scratch buffers threaded through the attack-evaluation pipeline (findBest
// partition loop, bestAttackWithWeapons phase/weapon masks, bestSequence permutation search).
// Pooled on the Evaluator so one sizing amortises across every hand a long-running iterate pass
// evaluates.
//
// Fields are grouped by lifetime via embedded sub-structs:
//
//   - shapeBufs         constructed once per (handSize, weapons) pair, reused across all
//                       calls against the same shape. Sized at construction; never grows.
//   - permBufs          backing arrays for the per-permutation TurnState slice fields.
//                       resetStateForPermutation re-slices these to [:0] on every
//                       permutation; mid-chain growth past the pre-sized cap allocates a
//                       fresh backing array (rare).
//   - carryWinnerBufs   sliding-window CarryState scratches, one per nesting level
//                       (sequence / mask-combo / leaf). Each level's scratch is reused
//                       across iterations at that level via CarryState.CopyFrom; ownership
//                       rules are documented per field.
//
// Embedded so call sites address fields directly (bufs.pcBuf) without a sub-struct prefix.

// shapeBufs holds the buffers sized once at construction from (handSize, weapons). They
// stay shape-stable across every call against this attackBufs and never need re-sizing.
type shapeBufs struct {
	pcBuf  []CardState
	ptrBuf []*CardState
	state  *TurnState
	// permMeta parallels pcBuf: each entry points into the global cardMetaCache so playSequence's
	// inner loop skips interface dispatch on Types / GoAgain and reads cached cost bounds.
	// Pointer-valued so bestSequence's permutation swaps move 8 bytes instead of a full struct.
	permMeta []*attackerMeta
	// attackerBuf is the per-mask-combo working slice that bestAttackWithWeapons fills with
	// the partition's attackers + the weapon-mask's selected weapons before handing off to
	// bestSequence. Sized at construction; the slice header re-slices to [:n] per call.
	attackerBuf []Card
	// weaponNames[mask] is the pre-built []string of weapon names indexed by the
	// weapon-prefix bits of the wmask, used for SwungWeapons display. Items don't
	// contribute to swung weapons so this stays weapon-only.
	weaponNames [][]string
	// activatedAbilities is the unified activated-ability list — weapons (positions
	// 0..len(weapons)-1) materialised at construction; items appended per Best call
	// from priorItems. The wmask iterates over the whole slice; index j's bit selects
	// activatedAbilities[j]. An "activated ability" is the same chain step whether it
	// came from a weapon or an item, so they share one list and one path. Per-Best
	// assembly re-slices back to the weapon prefix length before appending items,
	// leaving the cached weapon entries reusable across calls.
	activatedAbilities    []Card
	activatedAbilityCosts []int
	// weaponAbilityCount is len(weapons) at construction — the size of the cached
	// weapon prefix in activatedAbilities. Per-Best assembly re-slices both
	// activatedAbilities and activatedAbilityCosts back to this length before
	// appending items.
	weaponAbilityCount int
	// Partition-loop buffers, consumed by findBest. Sized handSize+1 to cover the optional
	// arsenal-in slot the enumerator treats as index n. isDRBuf caches card.TypeDefenseReaction
	// membership; canAttackBuf caches "this card can take Attack role" (Action or Weapon
	// subtype — Block-typed cards lack both).
	rolesBuf     []Role
	pitchVals    []int
	defenseVals  []int
	isDRBuf      []bool
	canAttackBuf []bool
	// pitchedValsScratch backs the per-leaf "pitched values" slice consumed by phase-mask
	// enumeration. Re-sliced to [:0] at the start of every leaf to eliminate a per-leaf alloc.
	pitchedValsScratch []int
	pitchedBuf         []Card
	// pitchPermBuf is the per-leaf pitch-ordering scratch — the chain runner permutes
	// pitched cards in place via Heap's algorithm, just like the attacker permutation. Sized
	// to handSize+1 so any pitch-role count the partition produces fits.
	pitchPermBuf []Card
	// pitchPermValsBuf parallels pitchPermBuf: pitchPermValsBuf[i] is the cached Pitch() of
	// pitchPermBuf[i]. The pitch Heap swaps both slices together so playSequenceWithMeta
	// can fetch pitch values without re-entering the Card.Pitch() interface call on every
	// pop — millions of permutations make per-pop dispatch one of the costliest paths.
	pitchPermValsBuf []int
	// pitchAttrBuf is the per-permutation flat backing array for per-CardState PitchedToPlay
	// slices. Layout: card 0's slice = pitchAttrBuf[s0:s1], card 1's = pitchAttrBuf[s1:s2],
	// adjacent windows. Total length across cards = pitches consumed (≤ len(pitched)). Sized
	// at construction to handSize+1 so append never reallocates the backing array — slice
	// headers stored on CardState stay valid across the active permutation.
	pitchAttrBuf []Card
	attackersBuf []Card
	defendersBuf []Card
	heldBuf      []Card
	// defenseGravScratch backs state.Graveyard during DR Plays. Reset via [:0]+append per
	// iteration so card effects can freely mutate their view without leaking into the next one.
	defenseGravScratch []Card
}

// permBufs holds the per-permutation slice backings. resetStateForPermutation seeds the
// TurnState's slice fields from these (via append([:0], ...)) so an unmodified permutation
// never reallocates: only mid-chain growth past the pre-sized cap forces a new backing
// array.
type permBufs struct {
	handBacking        []Card
	graveBacking       []Card
	banishBacking      []Card
	cardsPlayedBacking []Card
	logBacking         []LogEntry
	// logger is the recording-mode *TurnLogger reused across permutations. Allocated
	// once in newAttackBufs and rebound each permutation via SetBuffer(logBacking).
	// resetStateForPermutation points TurnState.logger at this when ctx is recording
	// or leaves TurnState.logger nil when ctx is skipping.
	logger              *turnlogger.TurnLogger
	auraTriggersBacking []Aura
	// itemsBacking backs TurnState.Items per permutation — seeded from priorItems and
	// freely mutated by item-ability Plays.
	itemsBacking []Item
	// defenderAurasBacking is the per-partition scratch for the post-defense aura set,
	// aliased by sequenceContext.defenderAuras.
	defenderAurasBacking []Aura
	// triggersBacking backs TurnState.Triggers. Sized small — Triggers are one-shot, so
	// the typical queue depth is one or two riders.
	triggersBacking []Trigger
}

// carryWinnerBufs holds the running-winner CarryState scratches — one per nesting level
// in the partition / mask-combo / permutation hierarchy. Each level's scratch is updated
// allocation-free via CarryState.CopyFrom; the ownership / aliasing rules per scratch
// are documented inline.
type carryWinnerBufs struct {
	// carryWinnerScratch is the per-Best-call sliding window into which bestSequence
	// snapshots the current winning permutation's end-of-chain state via
	// CarryState.SnapshotFromTurn. Slice backing arrays grow once across the lifetime of
	// the Evaluator's cached attackBufs and stay reused on every snapshot — per-Best
	// snapshots avoid reallocation. Reset (lengths to 0, scalars zeroed) at the top of
	// each bestSequence call so a stale value can't leak through when no permutation
	// lands a new best.
	carryWinnerScratch CarryState
	// bestCarryScratch is the mask-combo-level sliding window inside
	// bestAttackWithWeapons. Each new-best (pmask, wmask) update copies carryWinnerScratch
	// into this scratch via CarryState.CopyFrom (allocation-free after the first sizing).
	// bestAttackWithWeapons returns it as an alias — invalidated by the next call into
	// bestAttackWithWeapons against the same bufs, so callers that need the data to
	// outlive the next call must copy it out.
	bestCarryScratch CarryState
}

// attackBufs is the pooled scratch the attack-evaluation pipeline threads through every
// call. The Evaluator caches it so a single sizing amortises across every hand a deck
// eval visits.
type attackBufs struct {
	shapeBufs
	permBufs
	carryWinnerBufs
	// drScratch is a pooled TurnState for defense-reaction cost probing inside the
	// (pmask × wmask) loop; reusing its heap slot avoids a per-iteration alloc caused by
	// interface-call escape. Doesn't fit a sub-struct cleanly — it's a one-off TurnState
	// rather than a slice backing or a winner scratch.
	drScratch TurnState
	// drScratchAuras backs drScratch.Auras with the runechant aura entry when leftover
	// runechants > 0, so DR Cost reads s.Runechants() off this aura. Reused across
	// (pmask × wmask) iterations to avoid per-iter allocs.
	drScratchAuras []Aura
	// drCardStateScratch is a pooled *CardState handed to DR Card.Play calls. Each Play
	// takes a *CardState through an interface boundary so a literal &CardState{} would
	// escape and heap-alloc once per DR per partition — reusing this slot keeps the
	// whole defense-phase replay allocation-free. Reset per call by the caller.
	drCardStateScratch CardState
}

func newAttackBufs(handSize, weaponCount int, weapons []Weapon) *attackBufs {
	// +1 reserves a slot for the arsenal-in card, which joins attackers or defenders when the
	// enumerator plays it from arsenal. +maxDrawnExtensions leaves headroom for mid-turn-drawn
	// cards that play as chain extensions — cheap cycling cards (cost 0, Go again, draws a
	// card) can extend a chain well past the starting hand size.
	const maxDrawnExtensions = 32
	maxAttackers := handSize + weaponCount + 1 + maxDrawnExtensions
	// Materialise the weapon prefix of activatedAbilities once: w.Ability() returns the
	// stable package-cached interface value (per dev-standards), and its cost flows
	// through attackerMetaPtrFor's global cache. Per-Best assembly re-slices back to
	// this prefix length before appending items.
	activatedAbilities := make([]Card, len(weapons))
	activatedAbilityCosts := make([]int, len(weapons))
	for i, w := range weapons {
		ab := w.Ability()
		activatedAbilities[i] = ab
		activatedAbilityCosts[i] = attackerMetaPtrFor(ab).maxCost
	}
	// weaponNames stays shape-stable since it's purely a function of the weapons loadout
	// — the SwungWeapons display reads it per-mask without rebuilding.
	numMasks := 1 << weaponCount
	weaponNames := make([][]string, numMasks)
	for mask := 0; mask < numMasks; mask++ {
		var names []string
		for i, w := range weapons {
			if mask&(1<<i) != 0 {
				names = append(names, w.Name())
			}
		}
		weaponNames[mask] = names
	}
	pcBuf := make([]CardState, maxAttackers)
	ptrBuf := make([]*CardState, maxAttackers)
	// Wire the ptrBuf entries to their pcBuf slots once — the mapping is stable across every
	// permutation so playSequenceWithMeta doesn't need to rewrite it per call.
	for i := range pcBuf {
		ptrBuf[i] = &pcBuf[i]
	}
	// Per-permutation backing capacity for the Log slice. Log accumulates one entry per
	// chain step + every rider/trigger (typically 2-6 per step) so 64 covers the long tail.
	const logBackingCap = 64
	return &attackBufs{
		shapeBufs: shapeBufs{
			pcBuf:                 pcBuf,
			ptrBuf:                ptrBuf,
			state:                 &TurnState{},
			permMeta:              make([]*attackerMeta, maxAttackers),
			attackerBuf:           make([]Card, maxAttackers),
			weaponNames:           weaponNames,
			activatedAbilities:    activatedAbilities,
			activatedAbilityCosts: activatedAbilityCosts,
			weaponAbilityCount:    len(weapons),
			rolesBuf:              make([]Role, handSize+1),
			pitchVals:             make([]int, handSize+1),
			defenseVals:           make([]int, handSize+1),
			isDRBuf:               make([]bool, handSize+1),
			canAttackBuf:          make([]bool, handSize+1),
			pitchedValsScratch:    make([]int, 0, handSize+1),
			pitchedBuf:            make([]Card, 0, handSize+1),
			pitchPermBuf:          make([]Card, 0, handSize+1),
			pitchPermValsBuf:      make([]int, 0, handSize+1),
			pitchAttrBuf:          make([]Card, 0, handSize+1),
			attackersBuf:          make([]Card, 0, handSize+1),
			defendersBuf:          make([]Card, 0, handSize+1),
			heldBuf:               make([]Card, 0, handSize+1),
			defenseGravScratch:    make([]Card, 0, handSize+1),
		},
		permBufs: permBufs{
			handBacking:          make([]Card, 0, maxAttackers),
			graveBacking:         make([]Card, 0, maxAttackers),
			banishBacking:        make([]Card, 0, handSize+1),
			cardsPlayedBacking:   make([]Card, 0, maxAttackers),
			logBacking:           make([]LogEntry, 0, logBackingCap),
			logger:               turnlogger.New(),
			auraTriggersBacking:  make([]Aura, 0, handSize+1),
			itemsBacking:         make([]Item, 0, 4),
			defenderAurasBacking: make([]Aura, 0, handSize+1),
			triggersBacking:      make([]Trigger, 0, 4),
		},
		// carryWinnerBufs starts zero-valued — the slice backings grow on first use.
	}
}

// getAttackBufs returns the Evaluator's cached attackBufs when (handSize, weapons) match the
// last call; otherwise allocates a fresh one and caches it. weapons are zero-size structs so
// interface equality is a stable identity check. For a single deck eval (10k shuffles, same
// hand size + same weapons) this allocates once and reuses on every subsequent call —
// attackBufs is the second-biggest allocator after the eval-time slice copies.
func (e *Evaluator) getAttackBufs(handSize int, weapons []Weapon) *attackBufs {
	if e.cachedBufs != nil && e.cachedHandSize == handSize && sameWeapons(e.cachedWeapons, weapons) {
		return e.cachedBufs
	}
	e.cachedBufs = newAttackBufs(handSize, len(weapons), weapons)
	e.cachedHandSize = handSize
	// Snapshot the weapons slice header — caller may reuse the slice across calls. The
	// underlying Weapon values are zero-size structs; interface equality compares
	// (type, nil-data) tuples that are stable across calls.
	e.cachedWeapons = append(e.cachedWeapons[:0], weapons...)
	return e.cachedBufs
}

// sameWeapons reports whether two weapon slices contain the same weapons in the same order.
// Element-wise interface equality works because every weapon implementation is a zero-size
// struct, making interface values comparable and stable across calls.
func sameWeapons(a, b []Weapon) bool {
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
// Pitch / Defense magnitudes, Defense-Reaction membership, and Attack-role eligibility.
// Computing them up front keeps the recurse's inner body free of card-method / type-assert
// calls, which would otherwise repeat on every leaf. totalN covers the optional arsenal-in
// slot at index n; when present, its Defense picks up ArsenalDefenseBonus so the partition /
// capping pipeline sees the effective value.
func fillPartitionPerCardBufs(hand []Card, n, totalN int, arsenalCardIn Card, pvals, dvals []int, isDR, canAttack []bool) {
	for i := 0; i < totalN; i++ {
		var c Card
		if i < n {
			c = hand[i]
		} else {
			c = arsenalCardIn
		}
		pvals[i] = c.Pitch()
		dvals[i] = c.Defense()
		// Arsenal slot (i == n) lives at the end. Defense Reactions whose +N{d} rider only fires
		// when played from arsenal opt in via ArsenalDefenseBonus; bump the static Defense() up
		// here so the partition / capping pipeline sees the effective value.
		if i == n {
			dvals[i] += arsenalDefenseBonusOf(c)
		}
		m := attackerMetaPtrFor(c)
		ts := m.types
		isDR[i] = m.actsAsDR
		// Attack role covers every hand card that can resolve during the action chain —
		// Action subtypes (attack actions, non-attack actions, Auras, Instants) and Attack
		// Reactions (free chain steps). Block-typed cards and Defense Reactions have neither
		// subtype and can only pitch / block. Weapons never appear in hand — equipped
		// weapons stream their Ability() in via the wmask path, not the partition.
		canAttack[i] = ts.Has(card.TypeAction) || ts.Has(card.TypeAttackReaction)
	}
}
