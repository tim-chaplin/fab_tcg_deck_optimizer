package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Pre-allocated scratch threaded through the attack-evaluation pipeline (findBest, the
// pmask/wmask loop, bestSequence). Pooled on the Evaluator so one sizing amortises across
// every hand. Engine state isn't pooled — each permutation runs against a fresh leaf copy
// so per-perm mutations stay isolated.

// attackBufs holds the partition-level + chain-permutation scratch slices the search reuses
// across calls. Sized at construction; the slice headers re-slice to [:n] per call.
type attackBufs struct {
	// Chain-permutation scratch — pcBuf / ptrBuf back each chain step's CardState.
	pcBuf  []card.CardState
	ptrBuf []*card.CardState
	// permMeta parallels pcBuf, each entry pointing into cardMetaCache so playSequence's
	// inner loop skips interface dispatch on Types / GoAgain. Pointer-valued so perm swaps
	// move 8 bytes instead of a full struct.
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
	// 0..len(weapons)-1) materialised at construction; items appended per Best call. Per-Best
	// assembly re-slices to the weapon prefix before appending items.
	activatedAbilities    []card.Card
	activatedAbilityCosts []int
	// weaponAbilityCount is len(weapons) at construction — the size of the cached weapon
	// prefix in activatedAbilities.
	weaponAbilityCount int
	// maxResourceBonus bounds the resource points the turn's cards can add beyond printed
	// pitch (see resourceBonusUpperBound); the attack-budget prune relaxes by it.
	maxResourceBonus int

	// Partition-loop buffer, consumed by findBest. Sized handSize+1 to cover the optional
	// arsenal-in slot the enumerator treats as index n.
	partitionCards     []partitionCard
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
	// statePool hands out per-permutation *GameState scratch. preparePermState pulls a state
	// out (and CopyPersistentStateFrom seeds it from the leaf); losing perms / wmasks /
	// partitions and displaced best-winners are returned via Put so the same handful of
	// states cycle through the chain runner without re-allocating. runOneShuffle calls
	// FreeAll at end of shuffle to release the per-turn winner chain that escaped the pool.
	// Pointer-valued so a zero-value attackBufs literal isn't a usable Pool — callers
	// must build via newAttackBufs which seeds the cap.
	statePool *gameengine.Pool
	// pooledDeck is the recycled *Deck wrapper; preparePermState rebinds its slice headers
	// to alias ctx.deck. The winning perm's wrapper is cloned out via promoteWinnerDeck so
	// the next perm can rebind safely.
	pooledDeck *deck.Deck
	// pooledLeafDeck is the per-leaf scratch *Deck wrapper, ShallowCopied from the master
	// deck before runDefense so DR Plays (Rise Above's PrependToDeck, an Opt-ing DR, …)
	// mutate it instead of the shared master. ctx.deck rebinds to this wrapper so
	// preparePermState's ShallowCopyFrom propagates the post-DR state into the per-perm chain.
	pooledLeafDeck *deck.Deck
	// pooledEngine wraps the active permState; rebound per perm. The wrapper itself holds
	// no state, and no caller stashes ge across perms.
	pooledEngine *gameengine.GameEngine
	// pooledGravBuf backs the per-perm graveyard. preparePermState seeds it from
	// leafState.Graveyard() once per Best call and installs it on the pool state;
	// mid-chain AppendGraveyard grows past the prefix. promoteWinnerState clones the
	// winner's graveyard so the next reset can't scribble it. Hand and cardsPlayed
	// backings live on each pooled GameState now and don't need their own buf field.
	pooledGravBuf []card.Card
	// runDefenseDRGravBuf and runDefenseChainGravBuf back the two graveyard views the
	// defense pass installs on leafState: the DR-loop view (just defenders) and the
	// post-defense view (priorGraveyard + defenders). Recycled across runDefense calls.
	// preparePermState's gravBuf seed copies leafState's graveyard into its own backing
	// before the perm runs, so aliasing the leafState's slice to these buffers is safe.
	runDefenseDRGravBuf    []card.Card
	runDefenseChainGravBuf []card.Card
	// runDefenseHandBuf backs the role-tagged defense hand (held + attackers + pitched)
	// installed on leafState for the DR + plain-block phases. Recycled across runDefense
	// calls.
	runDefenseHandBuf []card.CardState
	// runDefensePostDRHeldBuf backs the post-DR Held-only view of state.HandStates().
	// The plain-block survivingHeld computation and the chain phase's handStart both
	// consume this slice, so a Held card a DR Play removed is automatically absent
	// from both. Recycled across runDefense calls.
	runDefensePostDRHeldBuf []card.Card
	// pooledLeafState is the per-Best leafState recycled across every call. bestAttackWithWeapons
	// resets it from masterState via CopyFrom so the per-call masterState.Copy()
	// allocation goes away. Defense mutations write through to this pool slot; preparePermState
	// then ResetForPermutationFrom this slot per perm.
	pooledLeafState *gameengine.GameState
	// pooledDRCostProbe is the recycled empty *GameEngine the variable-cost DR cost path
	// reads RunechantCount off. Lazy-init on first DR-cost call; per call we rewrite the
	// aura count rather than allocating a fresh engine + aura per probe.
	pooledDRCostProbe *gameengine.GameEngine
	// pooledDRProbeAura is the recycled runechant *aura.Aura the probe holds in its
	// auras slot, lazy-init alongside pooledDRCostProbe. Per probe we rewrite its Count
	// instead of building a fresh aura via token.NewRunechant.
	pooledDRProbeAura *aura.Aura
	// pooledDRProbeAuras is the 1-element slice header pooledDRProbeAura is installed
	// through via GameState.SetAuras — pooled so the probe stays allocation-free.
	pooledDRProbeAuras []gameengine.Aura
	// pooledSequenceCtx is the recycled per-partition-leaf sequenceContext.
	// newSequenceContext zeroes it and fills the active fields in place so the per-leaf
	// alloc goes away. Sole users are bestAttackWithWeapons and the one-shot replay path;
	// neither retains the ctx past its call, so pool reuse is safe.
	pooledSequenceCtx *sequenceContext
	// Cache-solution scratch: seqAttack/seqPitch hold the winning attacker order+modes and
	// pitch order, defModes the per-defender blocker modes. partSolution folds them in for
	// the winning partition; bestSolution holds the overall winner for the cache store.
	seqAttack    []playedCard
	seqPitch     []card.Card
	defModes     []playedCard
	partSolution cacheSolution
	bestSolution cacheSolution
}

// statePoolCap sizes the per-Evaluator *GameState pool. Measured peak in-flight on a
// 200-shuffle Viserai run with all Put-back sites + FreeAll-per-shuffle is 16; the cap
// here adds ~50% headroom for higher-fanout decks. Out-of-budget runs panic via the
// Pool's Get check — that's a hard signal to raise this constant rather than silently
// degrade into per-call allocation.
const statePoolCap = 24

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
		statePool:             gameengine.NewPool(statePoolCap),
		partitionCards:        make([]partitionCard, handSize+1),
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

// fillPartitionCards populates pcards[:totalN] from the hand and the optional arsenal-in
// card: each card's static Pitch / Defense magnitudes, Defense-Reaction membership, and
// Attack-role eligibility. The role field is left for the recurse to enumerate.
func fillPartitionCards(hand []card.Card, n, totalN int, arsenalCardIn card.Card, pcards []partitionCard) {
	for i := 0; i < totalN; i++ {
		fromArsenal := i == n
		c := arsenalCardIn
		if !fromArsenal {
			c = hand[i]
		}
		defenseVal := c.Defense()
		if fromArsenal {
			defenseVal += card.ArsenalDefenseBonusOf(c)
		}
		m := attackerMetaPtrFor(c)
		ts := m.types
		pcards[i] = partitionCard{
			card:        c,
			pitchVal:    c.Pitch(),
			defenseVal:  defenseVal,
			isDR:        m.actsAsDR,
			canAttack:   ts.Has(card.TypeAction) || ts.Has(card.TypeAttackReaction),
			fromArsenal: fromArsenal,
		}
	}
}
