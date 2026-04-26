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
	// perCardScratch is sized maxAttackers (handSize + weaponCount). Written by playSequence only
	// when the caller passes a non-nil perCardOut; bestSequence snapshots the winning
	// permutation's per-card damage from here into the caller's output buffer. The partition-loop
	// hot path passes nil and never touches this slice.
	perCardScratch []float64
	// perCardTriggerScratch parallels perCardScratch for hero-trigger damage (OnCardPlayed
	// return). Only written when the caller tracks.
	perCardTriggerScratch []float64
	// perCardAuraTriggerScratch parallels perCardScratch for mid-chain AuraTrigger damage
	// (fireAttackActionTriggers return). Split out from perCardTriggerScratch so the display
	// can separately attribute hero OnCardPlayed and mid-chain aura triggers.
	perCardAuraTriggerScratch []float64
	// fillContribWinnerOrder / fillContribPerCard are output buffers for bestSequence during
	// fillContributions's tracked replay. Kept on attackBufs so each Best call reuses the slab.
	fillContribWinnerOrder    []card.Card
	fillContribPerCard        []float64
	fillContribTriggerDmg     []float64
	fillContribAuraTriggerDmg []float64
	// fillContribUsed marks hand indices already assigned during chain→hand mapping. Sized
	// handSize; reset with clear before each fillContributions pass.
	fillContribUsed []bool
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
				names = append(names, w.Name())
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
	return &attackBufs{
		permMeta:                  make([]*attackerMeta, maxAttackers),
		pcBuf:                     pcBuf,
		ptrBuf:                    ptrBuf,
		state:                     &card.TurnState{},
		attackerBuf:               make([]card.Card, maxAttackers),
		weaponCosts:               weaponCosts,
		weaponNames:               weaponNames,
		rolesBuf:                  make([]Role, handSize+1),
		pitchVals:                 make([]int, handSize+1),
		defenseVals:               make([]int, handSize+1),
		isDRBuf:                   make([]bool, handSize+1),
		addsFutureValueBuf:        make([]bool, handSize+1),
		pitchedValsScratch:        make([]int, 0, handSize+1),
		pitchedBuf:                make([]card.Card, 0, handSize+1),
		attackersBuf:              make([]card.Card, 0, handSize+1),
		defendersBuf:              make([]card.Card, 0, handSize+1),
		heldBuf:                   make([]card.Card, 0, handSize+1),
		defenseGravScratch:        make([]card.Card, 0, handSize+1),
		perCardScratch:            make([]float64, maxAttackers),
		perCardTriggerScratch:     make([]float64, maxAttackers),
		perCardAuraTriggerScratch: make([]float64, maxAttackers),
		fillContribWinnerOrder:    make([]card.Card, maxAttackers),
		fillContribPerCard:        make([]float64, maxAttackers),
		fillContribTriggerDmg:     make([]float64, maxAttackers),
		fillContribAuraTriggerDmg: make([]float64, maxAttackers),
		fillContribUsed:           make([]bool, handSize),
	}
}

// getAttackBufs returns a fresh attackBufs sized for this hand. Callers allocate fresh per
// Best call.
func (e *Evaluator) getAttackBufs(handSize int, weapons []weapon.Weapon) *attackBufs {
	return newAttackBufs(handSize, len(weapons), weapons)
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
