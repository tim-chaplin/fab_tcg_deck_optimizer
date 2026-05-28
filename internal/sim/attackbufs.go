package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// Pre-allocated scratch threaded through the attack-evaluation pipeline (findBest, the
// pmask/wmask loop, bestSequence). Pooled on the Evaluator so one sizing amortises across
// every hand. All *GameState instances come from statePool (gameengine.NewPrewarmedPool).

// attackBufs holds the partition-level + attack-turn-permutation scratch slices the search reuses
// across calls. Sized at construction; the slice headers re-slice to [:n] per call.
type attackBufs struct {
	// Attack-turn-permutation scratch — pcBuf / ptrBuf back each attack step's CardState.
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
	// attackerWeaponIdxBuf parallels attackerBuf: the equipped-weapon index for each appended
	// weapon swing, or -1 for partition attackers / item abilities. bestSequence seeds each
	// CardState.WeaponIdx from it so the per-perm weapon object resolves correctly.
	attackerWeaponIdxBuf []int
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
	// pitchPcBuf is the per-leaf backing for the *CardState wrappers around pitched cards.
	// groupByRole writes pcBuf[i] = CardState{Card: pc.card} for each Pitch-role card and
	// stores the &pcBuf[i] pointer in pitchedBuf. Sized handSize+1 so it can carry every
	// hand card simultaneously even if every card pitches.
	pitchPcBuf       []card.CardState
	pitchedBuf       []*card.CardState
	pitchPermBuf     []*card.CardState
	pitchPermValsBuf []int
	pitchAttrBuf     []*card.CardState
	attackersBuf     []card.Card
	defendersBuf     []card.Card
	heldBuf          []card.Card
	// defenseGravScratch backs the per-defender graveyard view inside defendersDamage.
	defenseGravScratch []card.Card
	// drCardStateScratch is a pooled *CardState handed to DR Card.Play calls.
	drCardStateScratch card.CardState
	// statePool aliases the Evaluator's single *GameState pool. Every attackBufs the
	// Evaluator caches shares the same pointer, so all *GameState borrows (leafState,
	// per-perm attack-turn scratch, per-shuffle carry) draw from one fixed budget.
	statePool *gameengine.Pool
	// pooledEngine wraps the active permState; rebound per perm. The wrapper itself holds
	// no state, and no caller stashes ge across perms.
	pooledEngine *gameengine.GameEngine
	// runDefensePostDRHeldBuf backs the post-DR Held-only view of state.HandStates().
	// The plain-block survivingHeld computation and the attack turn phase's handStart both
	// consume this slice, so a Held card a DR Play removed is automatically absent
	// from both. Recycled across runDefense calls.
	runDefensePostDRHeldBuf []card.Card
	// pooledDRCostProbe is the recycled empty *GameEngine the variable-cost DR cost path
	// reads RunechantCount off. Lazy-init on first DR-cost call; per call we rewrite the
	// Runechant token slot's count via SetRunechantCount, no per-probe allocation.
	pooledDRCostProbe *gameengine.GameEngine
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

func newAttackBufs(handSize, weaponCount int, weapons []gameengine.Weapon, statePool *gameengine.Pool) *attackBufs {
	// +1 reserves a slot for the arsenal-in card; +maxDrawnExtensions leaves headroom for
	// mid-turn-drawn cards that play as attack-turn extensions.
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
		attackerWeaponIdxBuf:  make([]int, maxAttackers),
		weaponNames:           weaponNames,
		activatedAbilities:    activatedAbilities,
		activatedAbilityCosts: activatedAbilityCosts,
		weaponAbilityCount:    len(weapons),
		statePool:             statePool,
		partitionCards:        make([]partitionCard, handSize+1),
		pitchedValsScratch:    make([]int, 0, handSize+1),
		pitchPcBuf:            make([]card.CardState, handSize+1),
		pitchedBuf:            make([]*card.CardState, 0, handSize+1),
		pitchPermBuf:          make([]*card.CardState, 0, handSize+1),
		pitchPermValsBuf:      make([]int, 0, handSize+1),
		pitchAttrBuf:          make([]*card.CardState, 0, handSize+1),
		attackersBuf:          make([]card.Card, 0, handSize+1),
		defendersBuf:          make([]card.Card, 0, handSize+1),
		heldBuf:               make([]card.Card, 0, handSize+1),
		defenseGravScratch:    make([]card.Card, 0, handSize+1),
	}
}

// getAttackBufs returns the Evaluator's cached attackBufs when (handSize, weapons) match
// the last call; otherwise allocates a fresh one and caches it.
func (e *Evaluator) getAttackBufs(handSize int, weapons []gameengine.Weapon) *attackBufs {
	if e.cachedBufs != nil && e.cachedHandSize == handSize && sameWeapons(e.cachedWeapons, weapons) {
		return e.cachedBufs
	}
	e.cachedBufs = newAttackBufs(handSize, len(weapons), weapons, e.statePool)
	e.cachedHandSize = handSize
	e.cachedWeapons = append(e.cachedWeapons[:0], weapons...)
	return e.cachedBufs
}

// sameWeapons reports whether two weapon slices hold the same weapon kinds in the same
// order. Compares by CardID, not pointer: equip rebuilds fresh weapon objects every turn
// (CopyFrom deep-copies), so pointer identity would miss the cache every turn even though
// the loadout is unchanged.
func sameWeapons(a, b []gameengine.Weapon) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].CardID() != b[i].CardID() {
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
