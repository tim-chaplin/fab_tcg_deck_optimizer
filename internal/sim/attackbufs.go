package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Pre-allocated scratch buffers threaded through the attack-evaluation pipeline (findBest
// partition loop, bestAttackWithWeapons phase/weapon masks, bestSequence permutation
// search). Pooled on the Evaluator so one sizing amortises across every hand a long-running
// iterate pass evaluates. The engine state itself isn't pooled — each permutation runs
// against a fresh leafEngine.Copy() so per-permutation mutations stay isolated.

// attackBufs holds the partition-level + chain-permutation scratch slices the search reuses
// across calls. Sized at construction; the slice headers re-slice to [:n] per call.
type attackBufs struct {
	// Chain-permutation scratch — pcBuf / ptrBuf back each chain step's CardState.
	pcBuf  []card.CardState
	ptrBuf []*card.CardState
	// permMeta parallels pcBuf: each entry points into the global cardMetaCache so
	// playSequence's inner loop skips interface dispatch on Types / GoAgain and reads
	// cached cost bounds. Pointer-valued so bestSequence's permutation swaps move 8 bytes
	// instead of a full struct.
	permMeta []*attackerMeta
	// attackerBuf is the per-mask-combo working slice that bestAttackWithWeapons fills
	// with the partition's attackers + the weapon-mask's selected weapons before handing
	// off to bestSequence. Sized at construction; the slice header re-slices to [:n] per
	// call.
	attackerBuf []card.Card
	// weaponNames[mask] is the pre-built []string of weapon names indexed by the
	// weapon-prefix bits of the wmask, used for SwungWeapons display.
	weaponNames [][]string
	// activatedAbilities is the unified activated-ability list — weapons (positions
	// 0..len(weapons)-1) materialised at construction; items appended per Best call from
	// priorItems. Per-Best assembly re-slices back to the weapon prefix length before
	// appending items, leaving the cached weapon entries reusable across calls.
	activatedAbilities    []card.Card
	activatedAbilityCosts []int
	// weaponAbilityCount is len(weapons) at construction — the size of the cached weapon
	// prefix in activatedAbilities.
	weaponAbilityCount int

	// Partition-loop buffers, consumed by findBest. Sized handSize+1 to cover the optional
	// arsenal-in slot the enumerator treats as index n.
	rolesBuf           []deck.Role
	pitchVals          []int
	defenseVals        []int
	isDRBuf            []bool
	canAttackBuf       []bool
	pitchedValsScratch []int
	pitchedBuf         []card.Card
	pitchPermBuf       []card.Card
	pitchPermValsBuf   []int
	pitchAttrBuf       []card.Card
	attackersBuf       []card.Card
	defendersBuf       []card.Card
	heldBuf            []card.Card
	// defenseGravScratch backs the per-defender graveyard view inside defendersDamage.
	defenseGravScratch []card.Card
	// drCardStateScratch is a pooled *CardState handed to DR Card.Play calls.
	drCardStateScratch card.CardState
	// Per-permutation scratch recycled across every partition leaf the Best call evaluates.
	// Each preparePermState rebinds these to the active leaf's state shape, so a single set
	// of backings amortises across the whole findBest pass instead of per-leaf reallocation.
	// pooledState is the active recycled *GameState; preparePermState calls
	// CopyPersistentStateFrom(leafState) to rewrite it for the current leaf's perm. The
	// promoteWinnerState path clears this so the next perm allocates a fresh state and the
	// winner stays untouched.
	pooledState *gameengine.GameState
	// pooledDeck is the recycled *Deck wrapper; preparePermState rebinds its slice headers
	// to alias ctx.deck. The winning perm's wrapper is cloned out via promoteWinnerDeck so
	// the next perm can rebind safely.
	pooledDeck *deck.Deck
	// pooledEngine wraps the active permState; rebound per perm so chain runs use a single
	// engine wrapper. The wrapper holds no state of its own; no caller stashes ge across
	// perms.
	pooledEngine *gameengine.GameEngine
	// pooledHandBuf / pooledGravBuf back the per-perm hand and graveyard slices.
	// preparePermState rebuilds them via append into [:0] and SetHand / SetGraveyard them
	// on the pool state; mid-chain growth allocates fresh backing without trampling these.
	// promoteWinnerState clones the winner's hand and graveyard so the next reset can't
	// scribble the winning state.
	pooledHandBuf []card.Card
	pooledGravBuf []card.Card
	// pooledCardsPlayedBuf backs the per-perm cardsPlayed accumulator. The chain runner
	// appends to it via AppendCardsPlayed every chain step; pooling skips one allocation
	// per perm (first append spawns a fresh backing otherwise) plus growth events.
	// promoteWinnerState clones the winner's slice when the winner aliases this backing.
	pooledCardsPlayedBuf []card.Card
	// runDefenseDRGravBuf and runDefenseChainGravBuf back the two graveyard views the
	// defense pass installs on leafState: the DR-loop view (just defenders) and the
	// post-defense view (priorGraveyard + defenders). Recycled across runDefense calls.
	// preparePermState's gravBuf seed copies leafState's graveyard into its own backing
	// before the perm runs, so aliasing the leafState's slice to these buffers is safe.
	runDefenseDRGravBuf    []card.Card
	runDefenseChainGravBuf []card.Card
	// pooledLeafState is the per-Best leafState recycled across every call. bestAttackWithWeapons
	// resets it from masterState via CopyFrom so the per-call masterState.Copy()
	// allocation goes away. Defense mutations write through to this pool slot; preparePermState
	// then ResetForPermutationFrom this slot per perm.
	pooledLeafState *gameengine.GameState
}

func newAttackBufs(handSize, weaponCount int, weapons []weapon.Weapon) *attackBufs {
	// +1 reserves a slot for the arsenal-in card; +maxDrawnExtensions leaves headroom for
	// mid-turn-drawn cards that play as chain extensions.
	const maxDrawnExtensions = 32
	maxAttackers := handSize + weaponCount + 1 + maxDrawnExtensions
	activatedAbilities := make([]card.Card, len(weapons))
	activatedAbilityCosts := make([]int, len(weapons))
	for i, w := range weapons {
		ab := w.Ability()
		activatedAbilities[i] = ab
		activatedAbilityCosts[i] = attackerMetaPtrFor(ab).maxCost
	}
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
	pcBuf := make([]card.CardState, maxAttackers)
	ptrBuf := make([]*card.CardState, maxAttackers)
	for i := range pcBuf {
		ptrBuf[i] = &pcBuf[i]
	}
	return &attackBufs{
		pcBuf:                 pcBuf,
		ptrBuf:                ptrBuf,
		permMeta:              make([]*attackerMeta, maxAttackers),
		attackerBuf:           make([]card.Card, maxAttackers),
		weaponNames:           weaponNames,
		activatedAbilities:    activatedAbilities,
		activatedAbilityCosts: activatedAbilityCosts,
		weaponAbilityCount:    len(weapons),
		rolesBuf:              make([]deck.Role, handSize+1),
		pitchVals:             make([]int, handSize+1),
		defenseVals:           make([]int, handSize+1),
		isDRBuf:               make([]bool, handSize+1),
		canAttackBuf:          make([]bool, handSize+1),
		pitchedValsScratch:    make([]int, 0, handSize+1),
		pitchedBuf:            make([]card.Card, 0, handSize+1),
		pitchPermBuf:          make([]card.Card, 0, handSize+1),
		pitchPermValsBuf:      make([]int, 0, handSize+1),
		pitchAttrBuf:          make([]card.Card, 0, handSize+1),
		attackersBuf:          make([]card.Card, 0, handSize+1),
		defendersBuf:          make([]card.Card, 0, handSize+1),
		heldBuf:               make([]card.Card, 0, handSize+1),
		defenseGravScratch:    make([]card.Card, 0, handSize+1),
	}
}

// getAttackBufs returns the Evaluator's cached attackBufs when (handSize, weapons) match
// the last call; otherwise allocates a fresh one and caches it.
func (e *Evaluator) getAttackBufs(handSize int, weapons []weapon.Weapon) *attackBufs {
	if e.cachedBufs != nil && e.cachedHandSize == handSize && sameWeapons(e.cachedWeapons, weapons) {
		return e.cachedBufs
	}
	e.cachedBufs = newAttackBufs(handSize, len(weapons), weapons)
	e.cachedHandSize = handSize
	e.cachedWeapons = append(e.cachedWeapons[:0], weapons...)
	return e.cachedBufs
}

// sameWeapons reports whether two weapon slices contain the same weapons in the same
// order.
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

// fillPartitionPerCardBufs writes the per-card values the partition recurse reads at each
// leaf: Pitch / Defense magnitudes, Defense-Reaction membership, and Attack-role
// eligibility.
func fillPartitionPerCardBufs(hand []card.Card, n, totalN int, arsenalCardIn card.Card, pvals, dvals []int, isDR, canAttack []bool) {
	for i := 0; i < totalN; i++ {
		var c card.Card
		if i < n {
			c = hand[i]
		} else {
			c = arsenalCardIn
		}
		pvals[i] = c.Pitch()
		dvals[i] = c.Defense()
		if i == n {
			dvals[i] += card.ArsenalDefenseBonusOf(c)
		}
		m := attackerMetaPtrFor(c)
		ts := m.types
		isDR[i] = m.actsAsDR
		canAttack[i] = ts.Has(card.TypeAction) || ts.Has(card.TypeAttackReaction)
	}
}
