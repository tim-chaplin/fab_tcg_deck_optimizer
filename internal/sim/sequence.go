package sim

// Attack-chain search: bestAttackWithWeapons evaluates one partition leaf across all
// phase / weapon masks, bestSequence picks the best ordering of attackers via Heap's
// algorithm, and playSequence* replay a single permutation through a fresh per-
// permutation GameState copy while firing hero triggers, Aura handlers, and per-attack
// OnHit closures.
//
// State lifecycle:
//   - findBest builds a master *GameState once per Best call (prior auras / items /
//     state + matchup config).
//   - evaluatePartition copies the master into a per-leaf state, runs defense on it
//     (which mutates auras / graveyard / value / incomingDamage), then enumerates
//     chain permutations.
//   - Each chain permutation runs against a fresh leafState.Copy() so per-permutation
//     mutations stay isolated. The winning copy's *GameState pointer is the partition's
//     result. Card hooks see the engine via state.Engine() — the engine wrapper carries
//     the rules-engine API the card.GameEngine interface demands.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/item"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// perItemAbilityCap caps how many instances of one item's activated ability the chain
// runner enumerates per turn, bounding the wmask 2^k explosion when an item's Count
// gets large. Realistic counts in play tend to 1-3; 4 leaves headroom without letting a
// pathological hand blow up the per-leaf mask loop.
const perItemAbilityCap = 4

// newSequenceContext builds the sequenceContext shared between the search path and the
// print-time replay path. It folds in item-ability instances and refreshes the pooled
// leafState from master, but does NOT run defense or seed the graveyard buf — callers do
// those steps after attaching any per-perm overrides (pmask, wmask, replayLogger).
func newSequenceContext(
	masterState *gameengine.GameState,
	weapons []weapon.Weapon,
	attackers, defenders, pitched, held []card.Card,
	d *deck.Deck,
	bufs *attackBufs,
	blockTotal, arsenalInIdx int,
	arsenalAtChainStart card.Card,
) *sequenceContext {
	ctx := &sequenceContext{
		hero:                masterState.Hero().(hero.Hero),
		pitched:             pitched,
		deck:                d,
		handStart:           held,
		arsenalAtChainStart: arsenalAtChainStart,
		bufs:                bufs,
		runechantCarryover:  auraCountByNameInState(masterState, "Runechant"),
		blockTotal:          blockTotal,
		arsenalInIdx:        arsenalInIdx,
		priorOpponentMarked: masterState.OpponentMarked(),
		priorBanish:         masterState.Banished(),
		priorGraveyard:      masterState.Graveyard(),
		defenders:           defenders,
		startOfTurnValue:    masterState.Value(),
		cacheable:           true,
	}
	abilities := bufs.activatedAbilities[:bufs.weaponAbilityCount]
	abilityCosts := bufs.activatedAbilityCosts[:bufs.weaponAbilityCount]
	for _, it := range masterState.Items() {
		copies := it.Count()
		if copies > perItemAbilityCap {
			copies = perItemAbilityCap
		}
		ab := it.Ability().(card.Card)
		cost := attackerMetaPtrFor(ab).maxCost
		for i := 0; i < copies; i++ {
			abilities = append(abilities, ab)
			abilityCosts = append(abilityCosts, cost)
		}
	}
	bufs.activatedAbilities = abilities
	bufs.activatedAbilityCosts = abilityCosts
	ctx.activatedAbilities = abilities
	ctx.activatedAbilityCosts = abilityCosts

	if bufs.pooledLeafState == nil {
		bufs.pooledLeafState = masterState.Copy()
	} else {
		bufs.pooledLeafState.CopyFrom(masterState)
	}
	ctx.leafState = bufs.pooledLeafState
	return ctx
}

// bestAttackWithWeapons enumerates phase / weapon masks for one partition leaf and
// returns the best (damage, futureValue, budget, swungWeapons, winnerState, legal,
// cacheable) tuple. Each per-leaf state branches off via masterState.Copy().
func bestAttackWithWeapons(
	masterState *gameengine.GameState,
	weapons []weapon.Weapon,
	attackers, defenders, pitched, held []card.Card,
	d *deck.Deck,
	bufs *attackBufs,
	blockTotal, arsenalInIdx, arsenalDefenderIdx int,
	arsenalAtChainStart card.Card,
) (int, int, chainBudget, []string, *gameengine.GameState, bool, bool) {
	ctx := newSequenceContext(masterState, weapons, attackers, defenders, pitched, held, d, bufs, blockTotal, arsenalInIdx, arsenalAtChainStart)
	hasDRs := containsDefenseReaction(defenders)
	hasModalBlocker := containsModalBlocker(defenders)
	incoming := masterState.IncomingDamage()

	var defenseDealtConst int
	defenseCacheableConst := true
	if !hasModalBlocker && len(defenders) > 0 {
		defenseDealtConst, defenseCacheableConst = ctx.runDefense(defenders, pitched, d, incoming, noBlockBudgetCap, arsenalDefenderIdx)
	}
	ctx.leafState.SetDeck(nil)
	defenseDealt := defenseDealtConst
	defenseCacheable := defenseCacheableConst

	// Seed bufs.pooledGravBuf's prefix with leafState's current graveyard once per Best
	// call. preparePermState re-slices the pool to len(leafGrav) per perm; the chain
	// runner only appends past the prefix so the prefix stays intact across perms and
	// the copy needn't be repeated. modal-blocker leaves re-seed inside the wmask loop
	// after each runDefense call (which rewrites leafGrav). Sized for the worst-case
	// mid-chain growth so preparePermState never needs to realloc.
	ctx.seedPoolGravBuf(len(attackers)+len(ctx.activatedAbilities), len(pitched))

	pitchedVals := bufs.pitchedValsScratch[:0]
	for _, c := range pitched {
		pitchedVals = append(pitchedVals, c.Pitch())
	}

	phaseCount := 1
	if hasDRs && len(pitched) > 0 {
		phaseCount = 1 << len(pitched)
	}

	attackersMinCost := 0
	for _, a := range attackers {
		attackersMinCost += attackerMetaPtrFor(a).minCost
	}

	copy(bufs.attackerBuf, attackers)

	bestDealt := 0
	bestFutureValue := 0
	var bestSwung []string
	var bestBudget chainBudget
	var bestWinner *gameengine.GameState
	foundFeasible := false

	// Loop invariants hoisted out of the pmask×wmask loops. weaponBitsMask /
	// totalAbilityMasks depend only on the constant weapons / abilities lists; drCost
	// depends only on defenders and the carryover runechants so it can be costed once
	// up front. The probe engine + runechant aura on bufs is only built / re-seeded
	// when at least one defender actually acts as a DR — leaves with only plain
	// blockers skip it entirely.
	weaponBitsMask := (1 << len(weapons)) - 1
	totalAbilityMasks := 1 << len(ctx.activatedAbilities)
	drCost := 0
	hasDRDefender := false
	for _, def := range defenders {
		if attackerMetaPtrFor(def).actsAsDR {
			hasDRDefender = true
			break
		}
	}
	if hasDRDefender {
		probe := ctx.drCostProbe(ctx.runechantCarryover)
		for _, def := range defenders {
			if !attackerMetaPtrFor(def).actsAsDR {
				continue
			}
			drCost += def.Cost(probe)
		}
	}

	for pmask := 0; pmask < phaseCount; pmask++ {
		phase := splitPitchesAcrossPhases(pitchedVals, pmask, phaseCount)

		ctx.resourceBudget = 0
		attackPitchPerm := bufs.pitchPermBuf[:0]
		attackPitchVals := bufs.pitchPermValsBuf[:0]
		for i, c := range pitched {
			if phaseCount > 1 && pmask&(1<<i) != 0 {
				continue
			}
			attackPitchPerm = append(attackPitchPerm, c)
			attackPitchVals = append(attackPitchVals, pitchedVals[i])
		}
		ctx.attackPitchPerm = attackPitchPerm
		ctx.attackPitchVals = attackPitchVals

		for wmask := 0; wmask < totalAbilityMasks; wmask++ {
			abilityCost := 0
			for j := range ctx.activatedAbilities {
				if wmask&(1<<j) != 0 {
					abilityCost += ctx.activatedAbilityCosts[j]
				}
			}
			if attackersMinCost+abilityCost > phase.attackBudget {
				continue
			}
			allAttackers := bufs.attackerBuf[:len(attackers)]
			for j, ab := range ctx.activatedAbilities {
				if wmask&(1<<j) != 0 {
					allAttackers = append(allAttackers, ab)
				}
			}
			dealt, futureValue, winner, legal := ctx.bestSequence(allAttackers)
			if !legal {
				continue
			}
			if drCost > phase.defendBudget {
				continue
			}
			if hasModalBlocker {
				defenseDealt, defenseCacheable = ctx.runDefense(defenders, pitched, d, incoming, phase.defendBudget-drCost, arsenalDefenderIdx)
				ctx.seedPoolGravBuf(len(allAttackers), len(attackPitchPerm))
			}
			if phase.hasDefendPitches && phase.defendBudget-drCost >= phase.maxDefendPitch {
				continue
			}
			candidateDrawn := winner.CardsDrawn()
			var bestDrawn int
			if bestWinner != nil {
				bestDrawn = bestWinner.CardsDrawn()
			}
			candidateHand := len(winner.Hand())
			var bestHand int
			if bestWinner != nil {
				bestHand = len(bestWinner.Hand())
			}
			cmp := chainScoreCmp(dealt, candidateDrawn, futureValue, bestDealt, bestDrawn, bestFutureValue)
			if !foundFeasible || cmp > 0 || (cmp == 0 && candidateHand > bestHand) {
				bestDealt = dealt
				bestFutureValue = futureValue
				bestSwung = bufs.weaponNames[wmask&weaponBitsMask]
				bestBudget = chainBudget{resource: phase.attackBudget, maxPitch: phase.maxAttackPitch, hasAttackPitches: phase.hasAttackPitches}
				bestWinner = winner
				foundFeasible = true
			}
		}
	}

	if !foundFeasible {
		return 0, 0, chainBudget{}, nil, nil, false, defenseCacheable
	}
	return bestDealt, defenseDealt, bestBudget, bestSwung, bestWinner, true, ctx.cacheable && defenseCacheable
}

// drCostProbe returns the pooled *GameEngine seeded with a runechant aura at count
// runechants (when > 0) for variable-cost DR cost probing. Defense-reactions read
// RunechantCount() off this engine to decide their Cost; no other state matters. The
// engine and its single runechant aura are lazily built on first call and reused across
// every Best-call probe — per call we rewrite the aura's Count instead of allocating.
func (ctx *sequenceContext) drCostProbe(runechants int) *gameengine.GameEngine {
	bufs := ctx.bufs
	ge := bufs.pooledDRCostProbe
	if ge == nil {
		ge = gameengine.New()
		bufs.pooledDRCostProbe = ge
		bufs.pooledDRProbeAura = token.NewRunechant(0)
	}
	gs := ge.GameState
	gs.ClearAuras()
	if runechants > 0 {
		bufs.pooledDRProbeAura.SetCount(runechants)
		gs.CreateAura(bufs.pooledDRProbeAura)
	}
	return ge
}

// sequenceContext carries the stable per-partition-leaf environment.
type sequenceContext struct {
	hero                  hero.Hero
	pitched               []card.Card
	deck                  *deck.Deck
	handStart             []card.Card
	arsenalAtChainStart   card.Card
	bufs                  *attackBufs
	attackPitchPerm       []card.Card
	attackPitchVals       []int
	resourceBudget        int
	runechantCarryover    int
	blockTotal            int
	arsenalInIdx          int
	priorOpponentMarked   bool
	priorBanish           []card.Card
	priorGraveyard        []card.Card
	activatedAbilities    []card.Card
	activatedAbilityCosts []int
	defenders             []card.Card
	leafState             *gameengine.GameState
	// startOfTurnValue is masterState.Value() captured at construction and re-seeded into
	// each per-perm state after ResetEphemeralState. Chain accumulators ride on top of the
	// start-of-action-phase aura tick, so summary.Value from Best includes that baseline.
	startOfTurnValue int
	cacheable        bool
	// replayLogger, when non-nil, is installed on each per-perm state so cards' log
	// emissions stream to it inline. PrintBestTurn sets it to a *gameengine.StreamLogger
	// pointed at stdout for single-chain replay; the eval hot path leaves it nil so the
	// state's default NoopLogger keeps every emission free.
	replayLogger card.Logger
	// permState is the last *GameState a playSequence call ran the chain against. Set
	// by playSequence so the test-only PermEngine accessor can read the post-chain
	// state; the hot bestSequence path threads the winner through return values and
	// leaves this nil.
	permState *gameengine.GameState
}

// seedPoolGravBuf grows bufs.pooledGravBuf if needed, copies leafState.Graveyard() into
// its prefix, and sets the slice length to len(leafGrav). preparePermState then re-uses
// the prefix verbatim across all permutations (chain appends only past the prefix, so
// the prefix never gets overwritten). maxChainAttackers and maxPitchPerm are upper
// bounds on the per-perm chain that determine the trailing headroom required to keep
// AppendGraveyard alloc-free.
func (ctx *sequenceContext) seedPoolGravBuf(maxChainAttackers, maxPitchPerm int) {
	bufs := ctx.bufs
	leafGrav := ctx.leafState.Graveyard()
	gravNeeded := len(leafGrav) + maxChainAttackers + maxPitchPerm
	if cap(bufs.pooledGravBuf) < gravNeeded {
		bufs.pooledGravBuf = make([]card.Card, len(leafGrav), gravNeeded)
		copy(bufs.pooledGravBuf, leafGrav)
		return
	}
	bufs.pooledGravBuf = bufs.pooledGravBuf[:len(leafGrav)]
	copy(bufs.pooledGravBuf, leafGrav)
}

// permEngine returns ctx.bufs.pooledEngine with its embedded *GameState rebound to state.
// Lazily allocates the wrapper on first call.
func (ctx *sequenceContext) permEngine(state *gameengine.GameState) *gameengine.GameEngine {
	if ctx.bufs.pooledEngine == nil {
		ctx.bufs.pooledEngine = state.Engine()
		return ctx.bufs.pooledEngine
	}
	ctx.bufs.pooledEngine.GameState = state
	return ctx.bufs.pooledEngine
}

// promoteWinnerDeck swaps winner's pooled-deck pointer for a freshly-allocated copy so
// ctx.bufs.pooledDeck stays free for the next permutation to reset. Without this hand-off,
// the next preparePermState's ShallowCopyFrom call would mutate the wrapper bestWinner
// still holds. The clone reuses the same shared slice backings (cap=len), so the chain
// runner's append-only mutations on either side still allocate fresh.
func (ctx *sequenceContext) promoteWinnerDeck(winner *gameengine.GameState) {
	if winner == nil {
		return
	}
	winnerDeck := winner.Deck()
	if winnerDeck != ctx.bufs.pooledDeck {
		return
	}
	winner.SetDeck(winnerDeck.ShallowCopy())
}

// promoteWinnerState hands off ctx.bufs.pooledState to bestWinner so the next permutation's
// CopyPersistentStateFrom doesn't trample it. The pool pointer is cleared; next perm
// reallocates a fresh state via CopyPersistentState. Wins are rare relative to losses
// (best-of-N! permutations), so the recycled-on-loss path is the common case and the
// allocation only happens once per new best.
//
// Also clones the winner's hand and graveyard so the next preparePermState's reseed
// can't trample the winning permutation's recorded state. The clones are unconditional
// because the slices may or may not still alias ctx.bufs.pooledHandBuf / pooledGravBuf
// after mid-chain growth — if growth happened they're already independent and the extra
// clone is a no-op cost; if not, the clone is the only thing keeping winner's state
// stable across the next perm's reset.
func (ctx *sequenceContext) promoteWinnerState(winner *gameengine.GameState) {
	if winner == nil {
		return
	}
	if winner == ctx.bufs.pooledState {
		ctx.bufs.pooledState = nil
	}
	winnerHand := winner.Hand()
	if len(winnerHand) > 0 {
		clone := make([]card.Card, len(winnerHand))
		copy(clone, winnerHand)
		winner.SetHand(clone)
	}
	winnerGrav := winner.Graveyard()
	if len(winnerGrav) > 0 {
		clone := make([]card.Card, len(winnerGrav))
		copy(clone, winnerGrav)
		winner.SetGraveyard(clone)
	}
	winnerCardsPlayed := winner.CardsPlayed()
	if len(winnerCardsPlayed) > 0 {
		clone := make([]card.Card, len(winnerCardsPlayed))
		copy(clone, winnerCardsPlayed)
		winner.SetCardsPlayed(clone)
	}
	// Banished aliases the pooled leafState's backing via CopyForPermutation's [:n:n]
	// slice and isn't reset per perm. Clone it so the next Best call's CopyFrom on
	// the leafState pool can't overwrite winner.banished's data.
	winnerBanished := winner.Banished()
	if len(winnerBanished) > 0 {
		clone := make([]card.Card, len(winnerBanished))
		copy(clone, winnerBanished)
		winner.SetBanished(clone)
	}
}

// runDefense mutates ctx.leafState through the defender list, accumulating per-DR Value
// into total. Auras grow with any DR-added entries; graveyard is left as priorGraveyard
// + defenders for the chain phase. Chain-locals (value, action points, …) get reset
// per permutation via ResetEphemeralState, so runDefense doesn't bother restoring them.
//
// SetIncomingDamage installs the matchup figure once and zeroes the damage-blocked
// accumulator; each DR's resolution and each plain block then bank into that accumulator,
// so leafState.RemainingUnblockedDamage() reads the unblocked remainder as defense
// proceeds while the matchup figure itself stays constant.
func (ctx *sequenceContext) runDefense(defenders, pitched []card.Card, deckPile *deck.Deck, matchupIncomingDamage, blockBudget, arsenalDefenderIdx int) (int, bool) {
	state := ctx.leafState
	if ctx.replayLogger != nil {
		state.SetLogger(ctx.replayLogger)
	}
	state.SetDeck(deckPile)
	state.SetIncomingDamage(matchupIncomingDamage)
	ge := state.Engine()
	cs := &ctx.bufs.drCardStateScratch

	total := 0
	cacheable := true

	// Per-DR view: graveyard = defenders so DRs that scan graveyard see the defender
	// set. runDefenseDRGravBuf is recycled across runDefense calls.
	drGraveyard := append(ctx.bufs.runDefenseDRGravBuf[:0], defenders...)
	ctx.bufs.runDefenseDRGravBuf = drGraveyard
	for i, def := range defenders {
		if !attackerMetaPtrFor(def).actsAsDR {
			continue
		}
		state.SetGraveyard(drGraveyard)
		state.SetPitched(pitched)
		state.SetDefenders(defenders)
		state.SetValue(0)
		state.SetCacheable(true)
		*cs = card.CardState{Card: def, FromArsenal: i == arsenalDefenderIdx}
		ge.ResolveChainStep(state.Logger(), cs)
		total += state.Value()
		if !state.IsCacheable() {
			cacheable = false
		}
	}

	// Plain blocks: walk surviving defenders, picking the best mode within blockBudget.
	state.SetDefenders(defenders)
	for _, def := range defenders {
		if attackerMetaPtrFor(def).actsAsDR {
			continue
		}
		bestMode, bestCost := pickBlockerMode(def, ge, cs, blockBudget)
		blockBudget -= bestCost
		*cs = card.CardState{Card: def, Mode: bestMode}
		if b, ok := def.(card.Blocker); ok {
			b.Block(ge, state.Logger(), cs)
		}
		block := cs.EffectiveDefense()
		if rem := state.RemainingUnblockedDamage(); block > rem {
			block = rem
		}
		if block > 0 {
			total += block
			state.AddDamageBlocked(block)
		}
	}

	// Leave state with graveyard = priorGraveyard + defenders for the chain phase.
	// runDefenseChainGravBuf is recycled across runDefense calls; preparePermState
	// copies the installed slice into bufs.pooledGravBuf before each perm so the
	// aliasing on leafState's graveyard is safe.
	chainGraveyard := append(ctx.bufs.runDefenseChainGravBuf[:0], ctx.priorGraveyard...)
	chainGraveyard = append(chainGraveyard, defenders...)
	ctx.bufs.runDefenseChainGravBuf = chainGraveyard
	state.SetGraveyard(chainGraveyard)

	return total, cacheable
}

// preparePermState returns a fresh per-permutation *GameState for the chain run. The
// state inherits the leafState's post-defense auras / items / graveyard / banished /
// hero / arsenal; ResetEphemeralState wipes the previous permutation's play state, then
// this permutation's inputs are installed. Hand is set to the chain attackers + the
// attack-phase pitched bag so each chain step's Hand() read sees the upcoming bag.
//
// IncomingDamage needs no re-install: the matchup figure rode in constant on leafState,
// and ResetEphemeralState zeroed the damage-blocked accumulator, so the attack chain
// already sees the full matchup figure.
//
// Deck wrapper recycling: ctx.bufs.pooledDeck is the scratch wrapper reused across
// permutations and across every partition leaf of this Best call. The bestWinner tracker
// clones the wrapper out before next perm runs (see bestSequence's tryOnce), so the pool
// slot is always free here.
func (ctx *sequenceContext) preparePermState(playedAttackers []*card.CardState, n int) *gameengine.GameState {
	bufs := ctx.bufs
	if bufs.pooledState == nil {
		bufs.pooledState = ctx.leafState.CopyPersistentState()
	} else {
		bufs.pooledState.CopyPersistentStateFrom(ctx.leafState)
	}
	s := bufs.pooledState
	s.ResetEphemeralState()
	s.SetValue(ctx.startOfTurnValue)
	// Hero and opponentMarked already mirror leafState (via the CopyPersistentStateFrom
	// above or the freshly-allocated copy) and leafState's values match ctx.hero /
	// ctx.priorOpponentMarked respectively — they all root back to masterState. Only
	// arsenal and blockTotal need explicit setting: arsenal may have been promoted out
	// of the leaf via findArsenalCard (nil-ed when the slot was reassigned to Attack/
	// Defend), and blockTotal is zeroed by ResetEphemeralState.
	s.SetArsenal(ctx.arsenalAtChainStart)
	s.SetBlockTotal(ctx.blockTotal)
	// pooledGravBuf's prefix was seeded with leafState.Graveyard() by
	// bestAttackWithWeapons / runDefense (modal path), and the chain runner only ever
	// appends past that prefix — so per perm we just re-slice back to the prefix's
	// length without re-copying.
	grav := bufs.pooledGravBuf
	s.SetGraveyard(grav)
	if bufs.pooledDeck == nil {
		bufs.pooledDeck = ctx.deck.ShallowCopy()
	} else {
		bufs.pooledDeck.ShallowCopyFrom(ctx.deck)
	}
	s.SetDeck(bufs.pooledDeck)
	needed := len(ctx.handStart) + n + len(ctx.attackPitchPerm)
	if cap(bufs.pooledHandBuf) < needed {
		bufs.pooledHandBuf = make([]card.Card, 0, needed)
	}
	hand := bufs.pooledHandBuf[:0]
	hand = append(hand, ctx.handStart...)
	for k := 0; k < n; k++ {
		hand = append(hand, playedAttackers[k].Card)
	}
	for _, c := range ctx.attackPitchPerm {
		hand = append(hand, c)
	}
	bufs.pooledHandBuf = hand
	// Seed cardsPlayed with the pooled backing. Headroom = n + len(attackPitchPerm)
	// covers the chain runner's per-step appends; mid-chain free chain steps that don't
	// consume hand cards (rare) grow this buffer normally.
	cpNeeded := n + len(ctx.attackPitchPerm)
	if cap(bufs.pooledCardsPlayedBuf) < cpNeeded {
		bufs.pooledCardsPlayedBuf = make([]card.Card, 0, cpNeeded)
	}
	bufs.pooledCardsPlayedBuf = bufs.pooledCardsPlayedBuf[:0]
	s.SetCardsPlayed(bufs.pooledCardsPlayedBuf)
	s.SetPitched(ctx.pitched)
	s.SetHand(hand)
	// ResetEphemeralState set s.logger to NoopLogger; PrintBestTurn-driven runs install
	// a StreamLogger here so every emission streams to the writer inline.
	if ctx.replayLogger != nil {
		s.SetLogger(ctx.replayLogger)
	}
	return s
}

// bestSequence tries every ordering of attackers and returns the max total damage plus
// the pendingFutureValue at the end of the winning permutation. legal=true when at
// least one ordering is playable. Returns the winning *GameState via the third return
// value.
func (ctx *sequenceContext) bestSequence(attackers []card.Card) (int, int, *gameengine.GameState, bool) {
	n := len(attackers)
	if n == 0 {
		if len(ctx.attackPitchPerm) > 0 {
			return 0, 0, nil, false
		}
		emptyAttackers := ctx.bufs.ptrBuf[:0]
		permState := ctx.preparePermState(emptyAttackers, 0)
		ge := ctx.permEngine(permState)
		if ge.HasEndOfTurnFire() {
			ge.FireEndOfTurn()
		}
		ctx.promoteWinnerDeck(permState)
		ctx.promoteWinnerState(permState)
		// permState.Value() carries the seeded baseline plus any EndOfTurn fire delta.
		return permState.Value(), pendingFutureValueFromState(permState), permState, true
	}
	pcBuf := ctx.bufs.pcBuf[:n]
	permMeta := ctx.bufs.permMeta[:n]
	for idx, c := range attackers {
		permMeta[idx] = attackerMetaPtrFor(c)
		ctx.seedChainEntry(&pcBuf[idx], c, idx)
	}

	best := 0
	bestFutureValue := 0
	var bestWinner *gameengine.GameState
	foundLegal := false
	pitchPerm := ctx.attackPitchPerm
	pitchVals := ctx.attackPitchVals
	pn := len(pitchPerm)
	tupleCount := 1
	for i := 0; i < n; i++ {
		tupleCount *= int(permMeta[i].modes)
	}
	hasModal := tupleCount > 1
	bestCardsDrawn := 0
	tryPitchOrdering := func() {
		tryOnce := func() {
			dmg, futureValue, _, winner, legal := ctx.playSequenceWithMeta(n)
			if !legal {
				return
			}
			if ctx.cacheable && !winner.IsCacheable() {
				ctx.cacheable = false
			}
			drawn := winner.CardsDrawn()
			cmp := chainScoreCmp(dmg, drawn, futureValue, best, bestCardsDrawn, bestFutureValue)
			if !foundLegal || cmp > 0 {
				best = dmg
				bestCardsDrawn = drawn
				bestFutureValue = futureValue
				foundLegal = true
				bestWinner = winner
				ctx.promoteWinnerDeck(winner)
				ctx.promoteWinnerState(winner)
			}
		}
		if !hasModal {
			tryOnce()
			return
		}
		for tuple := 0; tuple < tupleCount; tuple++ {
			rem := tuple
			for i := 0; i < n; i++ {
				modes := int(permMeta[i].modes)
				pcBuf[i].Mode = int8(rem % modes)
				rem /= modes
			}
			tryOnce()
		}
	}
	eval := func() {
		tryPitchOrdering()
		var pc [8]int
		pi := 0
		for pi < pn {
			if pc[pi] < pi {
				if pi&1 == 0 {
					pitchPerm[0], pitchPerm[pi] = pitchPerm[pi], pitchPerm[0]
					pitchVals[0], pitchVals[pi] = pitchVals[pi], pitchVals[0]
				} else {
					pitchPerm[pc[pi]], pitchPerm[pi] = pitchPerm[pi], pitchPerm[pc[pi]]
					pitchVals[pc[pi]], pitchVals[pi] = pitchVals[pi], pitchVals[pc[pi]]
				}
				tryPitchOrdering()
				pc[pi]++
				pi = 0
			} else {
				pc[pi] = 0
				pi++
			}
		}
	}
	eval()
	var c [8]int
	i := 0
	for i < n {
		if c[i] < i {
			if i&1 == 0 {
				pcBuf[0], pcBuf[i] = pcBuf[i], pcBuf[0]
				permMeta[0], permMeta[i] = permMeta[i], permMeta[0]
			} else {
				pcBuf[c[i]], pcBuf[i] = pcBuf[i], pcBuf[c[i]]
				permMeta[c[i]], permMeta[i] = permMeta[i], permMeta[c[i]]
			}
			eval()
			c[i]++
			i = 0
		} else {
			c[i] = 0
			i++
		}
	}
	return best, bestFutureValue, bestWinner, foundLegal
}

// playSequence is a thin wrapper that builds permMeta and calls playSequenceWithMeta.
// Records the last-played state on ctx.permState so the test-only PermEngine accessor can
// inspect the chain run's final state; bestSequence's hot path skips this write since the
// winner state is plumbed through return values instead.
func (ctx *sequenceContext) playSequence(order []card.Card) (damage int, futureValue int, residualBudget int, legal bool) {
	n := len(order)
	pcBuf := ctx.bufs.pcBuf
	meta := ctx.bufs.permMeta[:n]
	for i, c := range order {
		meta[i] = attackerMetaPtrFor(c)
		ctx.seedChainEntry(&pcBuf[i], c, i)
	}
	d, fv, rb, winner, lg := ctx.playSequenceWithMeta(n)
	ctx.permState = winner
	return d, fv, rb, lg
}

// seedChainEntry initialises one pcBuf slot for a fresh chain pass.
func (ctx *sequenceContext) seedChainEntry(pc *card.CardState, c card.Card, idx int) {
	pc.Card = c
	pc.FromArsenal = idx == ctx.arsenalInIdx
	pc.GrantedGoAgain = false
	pc.GrantedDominate = false
	pc.GrantedOverpower = false
	pc.BonusAttack = 0
	pc.BonusDefense = 0
	pc.PitchedToPlay = nil
	pc.OnHit = pc.OnHit[:0]
	pc.Mode = 0
}

// playSequenceWithMeta runs the permutation currently held in ctx.bufs.pcBuf[:n] with
// aligned permMeta[:n]. Returns (damage, futureValue, residualBudget, winner, legal).
// A nil winner indicates an infeasible permutation.
func (ctx *sequenceContext) playSequenceWithMeta(n int) (damage int, futureValue int, residualBudget int, winner *gameengine.GameState, legal bool) {
	pcBuf := ctx.bufs.pcBuf
	ptrBuf := ctx.bufs.ptrBuf
	meta := ctx.bufs.permMeta[:n]
	for i := 0; i < n; i++ {
		pcBuf[i].GrantedGoAgain = false
		pcBuf[i].GrantedDominate = false
		pcBuf[i].GrantedOverpower = false
		pcBuf[i].BonusAttack = 0
		pcBuf[i].BonusDefense = 0
		pcBuf[i].PitchedToPlay = nil
		pcBuf[i].OnHit = pcBuf[i].OnHit[:0]
	}
	played := ptrBuf[:n]
	state := ctx.preparePermState(played, n)
	ge := ctx.permEngine(state)
	pool := pitchPool{
		perm:      ctx.attackPitchPerm,
		vals:      ctx.attackPitchVals,
		n:         len(ctx.attackPitchPerm),
		remaining: ctx.resourceBudget,
		attr:      ctx.bufs.pitchAttrBuf[:0],
	}
	var activeAttack *card.CardState
	finalizeActiveAttack := func() {
		if activeAttack == nil {
			return
		}
		if gameengine.LikelyToHit(activeAttack) {
			for i := range activeAttack.OnHit {
				h := &activeAttack.OnHit[i]
				h.Fire(ge, state.Logger(), activeAttack, h)
			}
			types := activeAttack.Card.Types(nil)
			prevTriggering := state.TriggeringCard()
			state.SetTriggeringCard(activeAttack.Card)
			ge.FireHit(types)
			state.SetTriggeringCard(prevTriggering)
		}
		activeAttack = nil
	}
	for i, pc := range played {
		m := meta[i]
		if !m.isFreeChainStep {
			if state.ActionPoints() <= 0 {
				return 0, 0, 0, nil, false
			}
			state.AddActionPoints(-1)
		}
		state.RemoveFromHand(pc.Card)
		prevPitchIdx := pool.idx
		contrib, ok := pool.pay(m.costAt(ge, pc.Mode))
		if !ok {
			return 0, 0, 0, nil, false
		}
		pc.PitchedToPlay = contrib
		for k := prevPitchIdx; k < pool.idx; k++ {
			state.RemoveFromHand(pool.perm[k])
		}
		if m.types.IsAttackReaction() {
			if pre, ok := pc.Card.(card.PlayPrecondition); ok {
				if !pre.PlayPrecondition(ge, pc) {
					return 0, 0, 0, nil, false
				}
			}
			ar, ok := pc.Card.(card.AttackReaction)
			if !ok || activeAttack == nil || !ar.ARTargetAllowed(ge, activeAttack.Card, pc.Mode) {
				return 0, 0, 0, nil, false
			}
			ctx.hero.OnCardPlayed(pc.Card, ge, state.Logger())
			state.SetAttackReactionTarget(activeAttack)
			ge.ResolveChainStep(state.Logger(), pc)
			state.SetAttackReactionTarget(nil)
			state.AppendCardsPlayed(pc.Card)
			state.AppendGraveyard(pc.Card)
			if pc.EffectiveGoAgain(ge) {
				state.AddActionPoints(1)
			}
			continue
		}

		finalizeActiveAttack()
		if pre, ok := pc.Card.(card.PlayPrecondition); ok {
			if !pre.PlayPrecondition(ge, pc) {
				return 0, 0, 0, nil, false
			}
		}

		state.SetCardsRemaining(played[i+1:])

		ctx.hero.OnCardPlayed(pc.Card, ge, state.Logger())
		state.SetCurrentStepRerouted(false)
		ge.ResolveChainStep(state.Logger(), pc)
		if m.isAttack {
			ge.FireAttack(pc.Card)
			state.ClearOpponentMarked()
		}
		if m.isAttackAction {
			ge.FireAttackAction(pc.Card)
		}
		if m.isAttack {
			activeAttack = pc
		}
		state.AppendCardsPlayed(pc.Card)
		if m.types.IsNonAttackAction() {
			state.SetNonAttackActionPlayed(true)
		}
		if !m.types.PersistsInPlay() && !state.CurrentStepRerouted() {
			state.AppendGraveyard(pc.Card)
		}
		if pc.EffectiveGoAgain(ge) {
			state.AddActionPoints(1)
		}
	}
	finalizeActiveAttack()

	if pool.idx < pool.n {
		return 0, 0, 0, nil, false
	}
	if ge.HasEndOfTurnFire() {
		ge.FireEndOfTurn()
	}
	return state.Value(), pendingFutureValueFromState(state), pool.remaining, state, true
}

// pendingFutureValueFromState sums the Count of every Aura plus every Item on the state
// at end of chain — the partition tiebreaker's "hidden later-turn payoff" signal.
func pendingFutureValueFromState(gs *gameengine.GameState) int {
	if gs == nil {
		return 0
	}
	total := 0
	for _, a := range gs.Auras() {
		total += a.Count()
	}
	for _, it := range gs.Items() {
		total += it.Count()
	}
	return total
}

// pendingFutureValue sums the Count of every Aura plus every Item — used by
// partition.go when comparing partitions whose winning states are already in hand.
func pendingFutureValue(auras []*aura.Aura, items []*item.Item) int {
	total := 0
	for _, a := range auras {
		total += a.Count()
	}
	for _, it := range items {
		total += it.Count()
	}
	return total
}
