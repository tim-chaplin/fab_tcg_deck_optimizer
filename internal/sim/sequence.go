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
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
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
	ctx := bufs.pooledSequenceCtx
	if ctx == nil {
		ctx = &sequenceContext{}
		bufs.pooledSequenceCtx = ctx
	}
	*ctx = sequenceContext{
		pitched:             pitched,
		attackers:           attackers,
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
		ability := it.Ability()
		if ability == nil {
			// A triggered item (e.g. Talisman of Recompense) has no activated ability —
			// it fires through FireTriggers, so there is nothing to enqueue as a playable.
			continue
		}
		copies := it.Count()
		if copies > perItemAbilityCap {
			copies = perItemAbilityCap
		}
		ab := ability.(card.Card)
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
	// pooledState / recycledState are mutable per-perm scratch. A prior Best call's winner
	// threads in as the next turn's masterState with no intervening copy, so a pooled slot
	// can still alias this call's masterState. Drop the alias here so preparePermState
	// allocates a fresh state rather than mutating the caller's live state in place.
	if bufs.pooledState == masterState {
		bufs.pooledState = nil
	}
	if bufs.recycledState == masterState {
		bufs.recycledState = nil
	}
	return ctx
}

// installLeafDeck refreshes bufs.pooledLeafDeck from the master deck d and rebinds
// ctx.deck to it. Run before runDefense so DR Plays (Rise Above's PrependToDeck, an Opt-ing
// DR, ...) mutate the leaf-scoped wrapper rather than the master shared across leaves; the
// chain phase then reads the post-DR leaf deck via preparePermState's ShallowCopyFrom on
// ctx.deck. The wrapper itself is recycled across leaves; only its slice headers reset.
func installLeafDeck(ctx *sequenceContext, bufs *attackBufs, d *deck.Deck) {
	if bufs.pooledLeafDeck == nil {
		bufs.pooledLeafDeck = d.ShallowCopy()
	} else {
		bufs.pooledLeafDeck.ShallowCopyFrom(d)
	}
	ctx.deck = bufs.pooledLeafDeck
}

// bestAttackWithWeapons enumerates phase / weapon masks for one partition leaf and
// returns the best (damage, defenseDealt, budget, swungWeapons, winnerState, legal,
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
	// Cleared up front so a partition with no defenders leaves an empty defender capture.
	bufs.defModes = bufs.defModes[:0]
	hasDRs := containsDefenseReaction(defenders)
	hasModalBlocker := containsModalBlocker(defenders)
	incoming := masterState.IncomingDamage()

	var defenseDealtConst int
	defenseCacheableConst := true
	if !hasModalBlocker && len(defenders) > 0 {
		installLeafDeck(ctx, bufs, d)
		defenseDealtConst, defenseCacheableConst, ctx.handStart = ctx.runDefense(defenders, pitched, held, ctx.deck, incoming, noBlockBudgetCap, arsenalDefenderIdx, nil)
	} else if !hasModalBlocker && incoming > 0 {
		// No defenders, so runDefense doesn't run — but unblocked incoming damage still
		// fires DamageTaken so auras destroyed by taking damage leave the arena.
		ctx.leafState.SetIsMyTurn(false)
		ctx.leafState.Engine().FireTriggers(triggertype.DamageTaken, nil)
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

	var bestScore chainScore
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

		if drCost > phase.defendBudget {
			continue
		}
		if phase.hasDefendPitches && phase.defendBudget-drCost >= phase.maxDefendPitch {
			continue
		}
		if hasModalBlocker {
			// Defense resolves before the attack chain. A modal blocker's block depends on
			// phase.defendBudget, so the defense pass runs once per phase here.
			installLeafDeck(ctx, bufs, d)
			defenseDealt, defenseCacheable, ctx.handStart = ctx.runDefense(defenders, pitched, held, ctx.deck, incoming, phase.defendBudget-drCost, arsenalDefenderIdx, nil)
			ctx.seedPoolGravBuf(len(attackers)+len(ctx.activatedAbilities), len(attackPitchPerm))
		}

		for wmask := 0; wmask < totalAbilityMasks; wmask++ {
			abilityCost := 0
			for j := range ctx.activatedAbilities {
				if wmask&(1<<j) != 0 {
					abilityCost += ctx.activatedAbilityCosts[j]
				}
			}
			// A resource-producing card lifts the budget above printed pitch, so the prune
			// relaxes by maxResourceBonus (the declared upper bound); pay does the exact
			// funding check.
			if attackersMinCost+abilityCost > phase.attackBudget+bufs.maxResourceBonus {
				continue
			}
			allAttackers := bufs.attackerBuf[:len(attackers)]
			for j, ab := range ctx.activatedAbilities {
				if wmask&(1<<j) != 0 {
					allAttackers = append(allAttackers, ab)
				}
			}
			score, winner, legal := ctx.bestSequence(allAttackers)
			if !legal {
				continue
			}
			if !foundFeasible || score.cmp(bestScore) > 0 {
				bestScore = score
				bestSwung = bufs.weaponNames[wmask&weaponBitsMask]
				bestBudget = chainBudget{resource: phase.attackBudget, maxPitch: phase.maxAttackPitch, hasAttackPitches: phase.hasAttackPitches}
				// Hand the superseded leaf-best to bufs.recycledState so the next
				// preparePermState can reuse its struct + aura backing instead of
				// allocating a fresh CopyPersistentState.
				if bestWinner != nil && bufs.recycledState == nil {
					bufs.recycledState = bestWinner
				}
				bestWinner = winner
				foundFeasible = true
				sol := &bufs.partSolution
				sol.attack = append(sol.attack[:0], bufs.seqAttack...)
				sol.pitch = append(sol.pitch[:0], bufs.seqPitch...)
				sol.defenders = append(sol.defenders[:0], bufs.defModes...)
			}
		}
	}

	if !foundFeasible {
		return 0, 0, chainBudget{}, nil, nil, false, defenseCacheable
	}
	return bestScore.value, defenseDealt, bestBudget, bestSwung, bestWinner, true, ctx.cacheable && defenseCacheable
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
		bufs.pooledDRProbeAuras = []gameengine.Aura{bufs.pooledDRProbeAura}
	}
	gs := ge.GameState
	if runechants > 0 {
		bufs.pooledDRProbeAura.SetCount(runechants)
		gs.SetAuras(bufs.pooledDRProbeAuras)
	} else {
		gs.ClearAuras()
	}
	return ge
}

// sequenceContext carries the stable per-partition-leaf environment.
type sequenceContext struct {
	pitched               []card.Card
	attackers             []card.Card
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
	// Hand / Graveyard / CardsPlayed alias the pooled per-perm scratch buffers; Banished
	// aliases the pooled leafState's [:n:n] backing. The next Best call's CopyFrom on the
	// leafState pool or the next perm's preparePermState would trample them in place,
	// so clone each into an independent backing here.
	if h := winner.HandStates(); len(h) > 0 {
		winner.SetHandStates(cloneCardSlice(h))
	}
	if g := winner.Graveyard(); len(g) > 0 {
		winner.SetGraveyard(cloneCardSlice(g))
	}
	if cp := winner.CardsPlayed(); len(cp) > 0 {
		winner.SetCardsPlayed(cloneCardSlice(cp))
	}
	if b := winner.Banished(); len(b) > 0 {
		winner.SetBanished(cloneCardSlice(b))
	}
}

// cloneCardSlice returns an independent backing copy of src. Callers should gate on
// len(src) > 0 when they want to preserve the field's existing slice header on empty.
func cloneCardSlice[T any](src []T) []T {
	out := make([]T, len(src))
	copy(out, src)
	return out
}

// appendExcludingMultiset appends every entry in src to dst, except for the first
// occurrence of each card that appears in exclude (treated as a multiset). When exclude
// is empty the inner loop short-circuits to a single append. Used by runDefense to keep
// defenders that DRs banished out of the post-defense chain graveyard.
func appendExcludingMultiset(dst, src, exclude []card.Card) []card.Card {
	if len(exclude) == 0 {
		return append(dst, src...)
	}
	// Defender lists are tiny (<= handSize+1) so an alloc-free linear scan beats a map.
	skip := make([]bool, len(exclude))
	for _, c := range src {
		excluded := false
		for j, e := range exclude {
			if !skip[j] && e == c {
				skip[j] = true
				excluded = true
				break
			}
		}
		if !excluded {
			dst = append(dst, c)
		}
	}
	return dst
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
//
// Before the plain-block loop the genuine role-tagged defense hand — held + attackers +
// pitched — is installed; Discard consumes only a Held card. Returns the Held cards left
// after any Blocker discards.
// cachedModes, when non-nil, supplies each plain blocker's mode (parallel to defenders) so
// a cache replay skips the pickBlockerMode search; nil drives the normal mode pick.
func (ctx *sequenceContext) runDefense(defenders, pitched, held []card.Card, deckPile *deck.Deck, matchupIncomingDamage, blockBudget, arsenalDefenderIdx int, cachedModes []playedCard) (int, bool, []card.Card) {
	state := ctx.leafState
	state.SetIsMyTurn(false)
	if ctx.replayLogger != nil {
		state.SetLogger(ctx.replayLogger)
	}
	state.SetDeck(deckPile)
	state.SetIncomingDamage(matchupIncomingDamage)
	ge := state.Engine()
	cs := &ctx.bufs.drCardStateScratch

	// defModes captures each defender's resolved blocker mode, parallel to defenders. DRs
	// resolve at mode 0.
	if cap(ctx.bufs.defModes) < len(defenders) {
		ctx.bufs.defModes = make([]playedCard, len(defenders))
	}
	ctx.bufs.defModes = ctx.bufs.defModes[:len(defenders)]

	total := 0
	cacheable := true

	// Install the role-tagged defense hand before the DR loop so HeldHand() /
	// PopHandAt() (variable-cost DR Plays use these to remove a Held card) see only
	// the partition's Held subset, not masterState's full hand defaulted to Held.
	state.SetDefenders(defenders)
	defenseHand := ctx.bufs.runDefenseHandBuf[:0]
	for _, c := range held {
		defenseHand = append(defenseHand, card.CardState{Card: c, Role: card.Held})
	}
	for _, c := range ctx.attackers {
		defenseHand = append(defenseHand, card.CardState{Card: c, Role: card.Attack})
	}
	for _, c := range pitched {
		defenseHand = append(defenseHand, card.CardState{Card: c, Role: card.Pitch})
	}
	ctx.bufs.runDefenseHandBuf = defenseHand
	state.SetHandStates(defenseHand)

	// Per-DR view: graveyard = defenders so DRs that scan graveyard see the defender
	// set. runDefenseDRGravBuf is recycled across runDefense calls.
	drGraveyard := append(ctx.bufs.runDefenseDRGravBuf[:0], defenders...)
	ctx.bufs.runDefenseDRGravBuf = drGraveyard
	for i, def := range defenders {
		if !attackerMetaPtrFor(def).actsAsDR {
			continue
		}
		ctx.bufs.defModes[i] = playedCard{card: def}
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
	// origHeld is the Held-role slice of the post-DR HandStates; any card a DR Play
	// removed from hand drops out automatically.
	origHeld := ctx.bufs.runDefensePostDRHeldBuf[:0]
	for _, hs := range state.HandStates() {
		if hs.Role == card.Held {
			origHeld = append(origHeld, hs.Card)
		}
	}
	ctx.bufs.runDefensePostDRHeldBuf = origHeld
	for i, def := range defenders {
		if attackerMetaPtrFor(def).actsAsDR {
			continue
		}
		var bestMode int8
		if cachedModes != nil {
			bestMode = cachedModes[i].mode
		} else {
			var bestCost int
			bestMode, bestCost = pickBlockerMode(def, ge, cs, blockBudget)
			blockBudget -= bestCost
		}
		ctx.bufs.defModes[i] = playedCard{card: def, mode: bestMode}
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

	// A Blocker may have discarded leading Held cards; the survivors are origHeld's tail.
	discarded := len(origHeld) - (len(state.HandStates()) - len(ctx.attackers) - len(pitched))
	survivingHeld := origHeld[discarded:]

	// Leave state with graveyard = priorGraveyard + defenders + discarded cards for the
	// chain phase, EXCLUDING any defender a DR banished during the defense pass — those
	// cards moved to the banished zone and must not also appear in the graveyard, or
	// subsequent BanishFromGraveyard / RecycleFromGraveyard scans would see a phantom
	// copy (and the total card count across all zones would drift). runDefenseChainGravBuf
	// is recycled across runDefense calls; preparePermState copies the installed slice into
	// bufs.pooledGravBuf before each perm so the aliasing on leafState's graveyard is safe.
	chainGraveyard := append(ctx.bufs.runDefenseChainGravBuf[:0], ctx.priorGraveyard...)
	drBanished := state.Banished()[len(ctx.priorBanish):]
	chainGraveyard = appendExcludingMultiset(chainGraveyard, defenders, drBanished)
	chainGraveyard = append(chainGraveyard, origHeld[:discarded]...)
	ctx.bufs.runDefenseChainGravBuf = chainGraveyard
	state.SetGraveyard(chainGraveyard)

	// Defense phase over: if any incoming damage got through, fire DamageTaken so auras
	// destroyed by taking damage (Arcane Cussing, Bloodspill Invocation) leave the arena.
	if state.RemainingUnblockedDamage() > 0 {
		ge.FireTriggers(triggertype.DamageTaken, nil)
	}

	return total, cacheable, survivingHeld
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
		// Drain the recycled slot before falling back to a fresh struct alloc: superseded
		// best-winners (handed off by bestSequence's tryOnce when a better perm displaces
		// them) already carry independent slice backings, so CopyPersistentStateFrom can
		// rewrite the struct in place without leaking any prior-winner reference.
		if bufs.recycledState != nil {
			bufs.pooledState = bufs.recycledState
			bufs.recycledState = nil
			bufs.pooledState.CopyPersistentStateFrom(ctx.leafState)
		} else {
			bufs.pooledState = ctx.leafState.CopyPersistentState()
		}
	} else {
		bufs.pooledState.CopyPersistentStateFrom(ctx.leafState)
	}
	s := bufs.pooledState
	s.ResetEphemeralState()
	s.SetValue(ctx.startOfTurnValue)
	// Hero and opponentMarked already mirror leafState (via the CopyPersistentStateFrom
	// above or the freshly-allocated copy) and leafState's values match masterState's
	// (which leafState was copied from) and ctx.priorOpponentMarked respectively. Only
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
		bufs.pooledHandBuf = make([]card.CardState, 0, needed)
	}
	hand := bufs.pooledHandBuf[:0]
	for _, c := range ctx.handStart {
		hand = append(hand, card.CardState{Card: c, Role: card.Held})
	}
	for k := 0; k < n; k++ {
		hand = append(hand, card.CardState{Card: playedAttackers[k].Card, Role: card.Attack})
	}
	for _, c := range ctx.attackPitchPerm {
		hand = append(hand, card.CardState{Card: c, Role: card.Pitch})
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
	s.SetHandStates(hand)
	// ResetEphemeralState set s.logger to NoopLogger; PrintBestTurn-driven runs install
	// a StreamLogger here so every emission streams to the writer inline.
	if ctx.replayLogger != nil {
		s.SetLogger(ctx.replayLogger)
	}
	return s
}

// captureWinningSeq records the winning permutation's attacker order+modes and attack-pitch
// ordering into attackBufs scratch. Called at each new best, so when the search finishes
// the scratch describes the winning sequence — the raw material for a verbatim cache replay.
func (ctx *sequenceContext) captureWinningSeq(pcBuf []card.CardState, pitchPerm []card.Card) {
	b := ctx.bufs
	b.seqAttack = b.seqAttack[:0]
	for i := range pcBuf {
		b.seqAttack = append(b.seqAttack, playedCard{
			card: pcBuf[i].Card, mode: pcBuf[i].Mode, fromArsenal: pcBuf[i].FromArsenal,
		})
	}
	b.seqPitch = append(b.seqPitch[:0], pitchPerm...)
}

// bestSequence tries every ordering of attackers and returns the winning permutation's
// chainScore. legal=true when at least one ordering is playable. Returns the winning
// *GameState via the second return value.
func (ctx *sequenceContext) bestSequence(attackers []card.Card) (chainScore, *gameengine.GameState, bool) {
	n := len(attackers)
	if n == 0 {
		if len(ctx.attackPitchPerm) > 0 {
			return chainScore{}, nil, false
		}
		emptyAttackers := ctx.bufs.ptrBuf[:0]
		permState := ctx.preparePermState(emptyAttackers, 0)
		ge := ctx.permEngine(permState)
		ge.FireTriggers(triggertype.EndOfTurn, nil)
		ctx.captureWinningSeq(nil, nil)
		ctx.promoteWinnerDeck(permState)
		ctx.promoteWinnerState(permState)
		// permState.Value() carries the seeded baseline plus any EndOfTurn fire delta.
		return chainScoreOf(permState, permState.Value()), permState, true
	}
	pcBuf := ctx.bufs.pcBuf[:n]
	permMeta := ctx.bufs.permMeta[:n]
	for idx, c := range attackers {
		permMeta[idx] = attackerMetaPtrFor(c)
		ctx.seedChainEntry(&pcBuf[idx], c, idx)
	}

	var bestScore chainScore
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
	tryPitchOrdering := func() {
		tryOnce := func() {
			dmg, _, _, winner, legal := ctx.playSequenceWithMeta(n)
			if !legal {
				return
			}
			if ctx.cacheable && !winner.IsCacheable() {
				ctx.cacheable = false
			}
			score := chainScoreOf(winner, dmg)
			if !foundLegal || score.cmp(bestScore) > 0 {
				bestScore = score
				foundLegal = true
				// A prior best is superseded — hand its *GameState to bufs.recycledState so
				// the next preparePermState reuses it instead of allocating a fresh struct.
				// promoteWinnerState already cloned that prior winner's slices, so the
				// recycle slot's backings are independent.
				if bestWinner != nil && ctx.bufs.recycledState == nil {
					ctx.bufs.recycledState = bestWinner
				}
				bestWinner = winner
				ctx.captureWinningSeq(pcBuf, pitchPerm)
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
	return bestScore, bestWinner, foundLegal
}

// playSequence is a thin wrapper that builds permMeta and calls playSequenceWithMeta.
// Records the last-played state on ctx.permState so the test-only PermEngine accessor can
// inspect the chain run's final state; bestSequence's hot path skips this write since the
// winner state is plumbed through return values instead.
func (ctx *sequenceContext) playSequence(order []card.Card) (damage int, totalCounters int, residualBudget int, legal bool) {
	n := len(order)
	pcBuf := ctx.bufs.pcBuf
	meta := ctx.bufs.permMeta[:n]
	for i, c := range order {
		meta[i] = attackerMetaPtrFor(c)
		ctx.seedChainEntry(&pcBuf[i], c, i)
	}
	d, tc, rb, winner, lg := ctx.playSequenceWithMeta(n)
	ctx.permState = winner
	return d, tc, rb, lg
}

// playSequenceModal is playSequence for a cache replay: it seeds each chain step from a
// playedCard, applying the cached modal Mode rather than the default 0. The winning state
// lands on ctx.permState.
func (ctx *sequenceContext) playSequenceModal(order []playedCard) (damage int, totalCounters int, residualBudget int, legal bool) {
	n := len(order)
	pcBuf := ctx.bufs.pcBuf
	meta := ctx.bufs.permMeta[:n]
	for i, pc := range order {
		meta[i] = attackerMetaPtrFor(pc.card)
		ctx.seedChainEntry(&pcBuf[i], pc.card, i)
		pcBuf[i].Mode = pc.mode
		pcBuf[i].FromArsenal = pc.fromArsenal
	}
	d, tc, rb, winner, lg := ctx.playSequenceWithMeta(n)
	ctx.permState = winner
	return d, tc, rb, lg
}

// seedChainEntry initialises one pcBuf slot for a fresh chain pass.
func (ctx *sequenceContext) seedChainEntry(pc *card.CardState, c card.Card, idx int) {
	pc.Card = c
	pc.FromArsenal = idx == ctx.arsenalInIdx
	pc.Mode = 0
	resetPerPermChainStepFields(pc)
}

// resetPerPermChainStepFields zeroes the resolution-scratch fields a per-perm chain run
// dirties on a *CardState — granted keywords, bonus accumulators, OnHit registrations, and
// the pitch-to-play list. Card / FromArsenal / Mode are left alone since they ride the
// permutation's identity and are reseeded by the outer permutation loop.
func resetPerPermChainStepFields(pc *card.CardState) {
	pc.GrantedGoAgain = false
	pc.GrantedDominate = false
	pc.GrantedOverpower = false
	pc.GrantedInstant = false
	pc.BonusAttack = 0
	pc.BonusDefense = 0
	pc.PitchedToPlay = nil
	pc.OnHit = pc.OnHit[:0]
}

// playSequenceWithMeta runs the permutation currently held in ctx.bufs.pcBuf[:n] with
// aligned permMeta[:n]. Returns (damage, totalCounters, residualBudget, winner, legal).
// A nil winner indicates an infeasible permutation.
func (ctx *sequenceContext) playSequenceWithMeta(n int) (damage int, totalCounters int, residualBudget int, winner *gameengine.GameState, legal bool) {
	pcBuf := ctx.bufs.pcBuf
	ptrBuf := ctx.bufs.ptrBuf
	meta := ctx.bufs.permMeta[:n]
	for i := 0; i < n; i++ {
		resetPerPermChainStepFields(&pcBuf[i])
	}
	played := ptrBuf[:n]
	state := ctx.preparePermState(played, n)
	state.SetIsMyTurn(true)
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
		hit := gameengine.LikelyToHit(activeAttack)
		state.SetLastAttackHit(hit)
		if hit {
			for i := range activeAttack.OnHit {
				h := &activeAttack.OnHit[i]
				h.Fire(ge, state.Logger(), activeAttack, h)
			}
			ge.FireTriggers(triggertype.Hit, activeAttack.Card)
		}
		activeAttack = nil
	}
	for i, pc := range played {
		m := meta[i]
		// RemoveFromHand returns false when an earlier chain step's Play moved this card
		// out of hand (e.g. a hand-on-top alt cost via PopHandAt + PrependToDeck). The
		// partition planned this card as a play / pitch using the pre-chain hand; if the
		// card is no longer there the partition is no longer realisable — reject it so
		// the optimiser doesn't credit a phantom play.
		if !state.RemoveFromHand(pc.Card) {
			return 0, 0, 0, nil, false
		}
		prevPitchIdx := pool.idx
		contrib, ok := pool.pay(ge, m.costAt(ge, pc.Mode))
		if !ok {
			return 0, 0, 0, nil, false
		}
		pc.PitchedToPlay = contrib
		for k := prevPitchIdx; k < pool.idx; k++ {
			if !state.RemoveFromHand(pool.perm[k]) {
				return 0, 0, 0, nil, false
			}
		}
		if m.types.IsAttackReaction() {
			if m.hasPlayPrecondition {
				if !pc.Card.(card.PlayPrecondition).PlayPrecondition(ge, pc) {
					return 0, 0, 0, nil, false
				}
			}
			ar, ok := pc.Card.(card.AttackReaction)
			if !ok || activeAttack == nil || !ar.ARTargetAllowed(ge, activeAttack.Card, pc.Mode) {
				return 0, 0, 0, nil, false
			}
			ge.FireTriggers(triggertype.CardOrAbility, pc.Card)
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
		// Runs after finalizeActiveAttack so an earlier attack's on-hit rider has had its
		// chance to set this card's GrantedInstant.
		if !m.isFreeChainStep && !pc.GrantedInstant {
			if state.ActionPoints() <= 0 {
				return 0, 0, 0, nil, false
			}
			state.AddActionPoints(-1)
		}
		if m.hasPlayPrecondition {
			if !pc.Card.(card.PlayPrecondition).PlayPrecondition(ge, pc) {
				return 0, 0, 0, nil, false
			}
		}

		state.SetCardsRemaining(played[i+1:])

		state.SetCurrentStepRerouted(false)
		// CardOrAbility fires once before the card resolves so play-triggered effects
		// land ahead of the played card's own effect.
		ge.FireTriggers(triggertype.CardOrAbility, pc.Card)
		ge.ResolveChainStep(state.Logger(), pc)
		if m.isAttack {
			state.ClearOpponentMarked()
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
	ge.FireTriggers(triggertype.EndOfTurn, nil)
	return state.Value(), pendingTotalCountersFromState(state), pool.remaining, state, true
}

// pendingTotalCountersFromState sums the Count of every Aura plus every Item on the
// state at end of chain — the partition's secondary tiebreaker. Counts pending aura
// fires alongside token stockpile at 1:1 weight; structurally less valuable than a real
// card in hand or arsenal (see pendingTotalCardsFromState).
func pendingTotalCountersFromState(gs *gameengine.GameState) int {
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

// chainScoreOf builds a leaf's chainScore from its end-of-chain winner state and the
// damage / block value credited this turn.
func chainScoreOf(winner *gameengine.GameState, value int) chainScore {
	return chainScore{
		value:         value,
		cardsPlayed:   len(winner.CardsPlayed()),
		totalCards:    pendingTotalCardsFromState(winner),
		totalCounters: pendingTotalCountersFromState(winner),
	}
}

// pendingTotalCardsFromState projects the cards available next turn: the post-refill hand
// (held cards topped up to intellect) plus an occupied arsenal. Scoring the post-refill
// hand — not the bare end-of-chain hand — is what lets the tiebreaker credit a chain that
// empties its hand into attacks and arsenals a card, since the emptied hand refills back
// up regardless. The refill is projected uncapped: a near-decked-out chain is scored
// slightly optimistically.
func pendingTotalCardsFromState(gs *gameengine.GameState) int {
	if gs == nil {
		return 0
	}
	held := len(gs.HandStates())
	intellect := gs.Hero().(hero.Hero).Intelligence()
	n := held + endOfTurnDraws(held, intellect)
	if gs.Arsenal() != nil {
		n++
	}
	return n
}
