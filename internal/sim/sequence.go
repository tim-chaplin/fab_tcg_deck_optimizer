package sim

// Attack-chain search: bestAttackWithWeapons evaluates one partition leaf across all
// phase / weapon masks; bestSequence picks the best attacker ordering; playSequence*
// replays one permutation against a fresh per-permutation GameState copy while firing
// hero triggers, Aura handlers, and per-attack OnHit closures.
//
// State lifecycle:
//   - findBest builds a master *GameState once per Best call.
//   - evaluatePartition copies the master into a per-leaf state, runs defense, then
//     enumerates chain permutations against fresh leafState copies so per-permutation
//     mutations stay isolated. The winning copy's *GameState is the partition's result.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// perItemAbilityCap caps per-turn ability instances per item to bound the wmask 2^k blowup
// when Count is large. Realistic in-play counts are 1-3; 4 leaves headroom.
const perItemAbilityCap = 4

// newSequenceContext builds the sequenceContext shared by the search and print-time replay
// paths. Folds in item-ability instances and refreshes the pooled leafState from master.
// Does NOT run defense or seed the graveyard buf — callers do that after attaching per-perm
// overrides (pmask, wmask, replayLogger).
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
			// Triggered items (Talisman of Recompense) have no activated ability — they fire
			// through FireTriggers, with nothing to enqueue as a playable.
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

	// leafState borrows a slot from bufs.statePool. ResetEphemeralState rearms per-turn
	// state before CopyPersistentStateFrom rewrites cross-turn carryover. The defense
	// pass installs hand / deck / pitched / defenders from bufs scratch before reading,
	// so CopyPersistentState's nil-out of those fields is fine. Caller must pair with
	// releaseLeafState (typically `defer ctx.releaseLeafState()`).
	leaf := bufs.statePool.Get()
	leaf.ResetEphemeralState()
	leaf.CopyPersistentStateFrom(masterState)
	ctx.leafState = leaf
	return ctx
}

// releaseLeafState returns the borrowed leafState slot to the pool. Safe on a nil slot.
func (ctx *sequenceContext) releaseLeafState() {
	if ctx.leafState == nil {
		return
	}
	ctx.bufs.statePool.Put(ctx.leafState)
	ctx.leafState = nil
}

// installLeafDeck refreshes bufs.pooledLeafDeck from the master deck d and rebinds ctx.deck.
// Run before runDefense so DR Plays (Rise Above's PrependToDeck, an Opt-ing DR, ...) mutate
// the leaf-scoped wrapper rather than the master shared across leaves. The wrapper is
// recycled across leaves; only its slice headers reset.
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
	defer ctx.releaseLeafState()
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
		ctx.permEngine(ctx.leafState).FireTriggers(triggertype.DamageTaken, nil)
	}
	ctx.leafState.SetDeck(nil)
	defenseDealt := defenseDealtConst
	defenseCacheable := defenseCacheableConst
	// Stay paired with bestWinner / sol.defenders; the loop-scoped defenseDealt is
	// clobbered every modal-blocker pmask iteration.
	bestDefenseDealt := defenseDealtConst
	bestDefenseCacheable := defenseCacheableConst

	// Seed pooledGravBuf's prefix with leafState's graveyard once per Best call. The chain
	// runner only appends past the prefix so it stays intact across perms; preparePermState
	// just re-slices to the prefix length. Modal-blocker leaves re-seed inside the wmask
	// loop after each runDefense rewrites leafGrav. Sized for worst-case mid-chain growth.
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

	var bestScore, bestTotalScore chainScore
	var bestSwung []string
	var bestBudget chainBudget
	var bestWinner *gameengine.GameState
	foundFeasible := false

	// Loop invariants hoisted out of pmask × wmask: weaponBitsMask / totalAbilityMasks
	// (constant per leaf), and drCost (constant given defenders + carryover runechants).
	// The probe engine + runechant aura is only built / re-seeded when a defender acts as
	// a DR — plain-block-only leaves skip it.
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
			// Resource producers can lift the budget above printed pitch; relax the prune
			// by maxResourceBonus and let pay do the exact funding check.
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
			// Rank pmasks by attack + defense, not attack alone, so a pmask funding a stronger
			// modal blocker can win on equal attack.
			totalScore := score
			totalScore.value += defenseDealt
			if !foundFeasible || totalScore.cmp(bestTotalScore) > 0 {
				bestScore = score
				bestTotalScore = totalScore
				bestDefenseDealt = defenseDealt
				bestDefenseCacheable = defenseCacheable
				bestSwung = bufs.weaponNames[wmask&weaponBitsMask]
				bestBudget = chainBudget{resource: phase.attackBudget, maxPitch: phase.maxAttackPitch, hasAttackPitches: phase.hasAttackPitches}
				// Hand the superseded leaf-best back to the pool so the next preparePermState
				// reuses its struct + slice backings.
				bufs.statePool.Put(bestWinner)
				bestWinner = winner
				foundFeasible = true
				sol := &bufs.partSolution
				sol.attack = append(sol.attack[:0], bufs.seqAttack...)
				sol.pitch = append(sol.pitch[:0], bufs.seqPitch...)
				sol.defenders = append(sol.defenders[:0], bufs.defModes...)
			} else {
				// Losing wmask: return its winner to the pool so the next pmask × wmask iter
				// reuses it.
				bufs.statePool.Put(winner)
			}
		}
	}

	if !foundFeasible {
		return 0, 0, chainBudget{}, nil, nil, false, defenseCacheable
	}
	return bestScore.value, bestDefenseDealt, bestBudget, bestSwung, bestWinner, true, ctx.cacheable && bestDefenseCacheable
}

// drCostProbe returns the pooled *GameEngine seeded with a runechant aura at the given
// count for variable-cost DR cost probing. DRs read only RunechantCount() to decide Cost.
// The engine + aura are lazily built once and reused across probes; per call rewrites Count.
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
	// startOfTurnValue is masterState.Value() captured at construction and re-seeded into each
	// per-perm state after ResetEphemeralState — chain accumulators ride on top of the
	// start-of-action-phase aura tick, so summary.Value includes that baseline.
	startOfTurnValue int
	cacheable        bool
	// replayLogger, when non-nil, is installed on each per-perm state so log emissions
	// stream inline. PrintBestTurn sets it to a stdout StreamLogger; the eval hot path
	// leaves it nil so the state's NoopLogger keeps emissions free.
	replayLogger card.Logger
	// permState records the last *GameState playSequence ran the chain against, so the
	// test-only PermEngine accessor can read it. The hot bestSequence path threads the
	// winner through return values and leaves this nil.
	permState *gameengine.GameState
}

// seedPoolGravBuf grows bufs.pooledGravBuf if needed, copies leafState.Graveyard() into its
// prefix, and slices to that length. preparePermState reuses the prefix verbatim across all
// permutations. maxChainAttackers + maxPitchPerm size the trailing headroom required to keep
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

// promoteWinnerDeck swaps winner's pooled-deck pointer for a fresh wrapper so the next perm
// can reset bufs.pooledDeck without trampling the winner. The clone reuses the same shared
// slice backings (cap=len), so append-only mutations on either side still allocate fresh.
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

// promoteWinnerState clones the winner's graveyard / banished into independent backings
// so the next preparePermState — which would rewrite the shared graveyard scratch the
// winner aliases (and leafState's banished [:n:n] alias) — can't trample them.
//
// Hand and cardsPlayed already live on the winner's own pooled GameState backing (see
// preparePermState) and the statePool keeps the winner out of circulation until the
// consumer drops it, so neither needs a per-promotion clone.
func (ctx *sequenceContext) promoteWinnerState(winner *gameengine.GameState) {
	if winner == nil {
		return
	}
	if g := winner.Graveyard(); len(g) > 0 {
		winner.SetGraveyard(cloneCardSlice(g))
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

// appendExcludingMultiset appends src to dst, skipping the first occurrence of each card in
// exclude (multiset semantics). Used by runDefense to keep DR-banished defenders out of the
// post-defense graveyard.
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

// runDefense mutates ctx.leafState through the defender list, accumulating per-DR Value.
// Auras grow with any DR-added entries; graveyard is left as priorGraveyard + defenders for
// the chain phase. Per-permutation chain-locals reset via ResetEphemeralState, so runDefense
// doesn't restore them.
//
// SetIncomingDamage installs the matchup figure once and zeroes the damage-blocked
// accumulator; each DR + plain block banks into that accumulator, so
// RemainingUnblockedDamage() reads the unblocked remainder while the matchup figure stays
// constant.
//
// Before the plain-block loop, the role-tagged defense hand (held + attackers + pitched) is
// installed so Discard consumes only Held. Returns the Held cards left after any Blocker
// discards. cachedModes, when non-nil, supplies each plain blocker's mode for a cache replay;
// nil drives the normal pickBlockerMode search.
func (ctx *sequenceContext) runDefense(defenders, pitched, held []card.Card, deckPile *deck.Deck, matchupIncomingDamage, blockBudget, arsenalDefenderIdx int, cachedModes []playedCard) (int, bool, []card.Card) {
	state := ctx.leafState
	state.SetIsMyTurn(false)
	if ctx.replayLogger != nil {
		state.SetLogger(ctx.replayLogger)
	}
	state.SetDeck(deckPile)
	state.SetIncomingDamage(matchupIncomingDamage)
	ge := ctx.permEngine(state)
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
	// DiscardToTopOfDeck (variable-cost DR Plays use these to remove a Held card)
	// see only the partition's Held subset, not masterState's full hand defaulted to Held.
	// Built into the slot's own prewarmed hand backing so the cap survives Put → Get.
	state.SetDefenders(defenders)
	defenseHand := state.HandStates()[:0]
	for _, c := range held {
		defenseHand = append(defenseHand, card.CardState{Card: c, Role: card.Held})
	}
	for _, c := range ctx.attackers {
		defenseHand = append(defenseHand, card.CardState{Card: c, Role: card.Attack})
	}
	for _, c := range pitched {
		defenseHand = append(defenseHand, card.CardState{Card: c, Role: card.Pitch})
	}
	state.SetHandStates(defenseHand)

	// Per-DR view: graveyard = defenders so DRs that scan graveyard see the defender
	// set. runDefenseDRGravBuf is recycled across runDefense calls. drGraveyard carries
	// across the DR loop so a banish or destroy by an earlier DR is reflected in the
	// view the next DR scans.
	drGraveyard := append(ctx.bufs.runDefenseDRGravBuf[:0], defenders...)
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
		drGraveyard = state.Graveyard()
	}
	ctx.bufs.runDefenseDRGravBuf = drGraveyard

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

	// Leave state with graveyard = priorGraveyard + defenders + discarded cards for the chain
	// phase, EXCLUDING any defender a DR banished — those cards already moved to banished, and
	// duplicating them in graveyard would let BanishFromGraveyard / RecycleFromGraveyard see
	// a phantom copy. runDefenseChainGravBuf is recycled across runDefense calls; per-perm
	// the slice gets copied into bufs.pooledGravBuf so leafState's aliasing is safe.
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

// preparePermState returns a fresh per-permutation *GameState for the chain run. The state
// inherits leafState's post-defense auras / items / graveyard / banished / hero / arsenal;
// ResetEphemeralState wipes the previous perm's play state, then this perm's inputs install.
// Hand = chain attackers + attack-phase pitched bag, so Hand() reads see the upcoming bag.
//
// IncomingDamage is not re-installed: the matchup figure rode in constant on leafState and
// ResetEphemeralState zeroed the damage-blocked accumulator.
//
// Deck wrapper recycling: ctx.bufs.pooledDeck is the scratch wrapper reused across perms;
// promoteWinnerDeck clones the wrapper out before the next perm runs, so this slot is free.
func (ctx *sequenceContext) preparePermState(playedAttackers []*card.CardState, n int) *gameengine.GameState {
	bufs := ctx.bufs
	s := bufs.statePool.Get()
	s.CopyPersistentStateFrom(ctx.leafState)
	s.ResetEphemeralState()
	s.SetValue(ctx.startOfTurnValue)
	// Hero / opponentMarked already mirror leafState (which mirrors masterState). Only
	// arsenal and blockTotal need explicit setting: arsenal may have been promoted out of
	// the leaf via findArsenalCard, and blockTotal is zeroed by ResetEphemeralState.
	s.SetArsenal(ctx.arsenalAtChainStart)
	s.SetBlockTotal(ctx.blockTotal)
	// pooledGravBuf's prefix was seeded with leafState.Graveyard() by bestAttackWithWeapons
	// / runDefense (modal path); the chain runner only appends past the prefix, so per perm
	// we just re-slice back to the prefix length.
	grav := bufs.pooledGravBuf
	s.SetGraveyard(grav)
	if bufs.pooledDeck == nil {
		bufs.pooledDeck = ctx.deck.ShallowCopy()
	} else {
		bufs.pooledDeck.ShallowCopyFrom(ctx.deck)
	}
	s.SetDeck(bufs.pooledDeck)
	// Build the per-perm hand into s's prewarmed backing. A cap shortfall means the
	// workload outgrew pooledStateHandCap — bump it rather than fall back to per-call
	// alloc.
	needed := len(ctx.handStart) + n + len(ctx.attackPitchPerm)
	hand := s.HandStates()
	if cap(hand) < needed {
		panic("sim: pooled state hand backing too small; raise pooledStateHandCap")
	}
	hand = hand[:0]
	for _, c := range ctx.handStart {
		hand = append(hand, card.CardState{Card: c, Role: card.Held})
	}
	for k := 0; k < n; k++ {
		hand = append(hand, card.CardState{Card: playedAttackers[k].Card, Role: card.Attack})
	}
	for _, c := range ctx.attackPitchPerm {
		hand = append(hand, card.CardState{Card: c, Role: card.Pitch})
	}
	s.SetHandStates(hand)
	// CardsPlayed is similarly prewarmed; the cap check guards mid-chain growth.
	cpNeeded := n + len(ctx.attackPitchPerm)
	cp := s.CardsPlayed()
	if cap(cp) < cpNeeded {
		panic("sim: pooled state cardsPlayed backing too small; raise pooledStateCardsPlayedCap")
	}
	s.SetCardsPlayed(cp[:0])
	s.SetPitched(ctx.pitched)
	// ResetEphemeralState set s.logger to NoopLogger; PrintBestTurn runs install a
	// StreamLogger so emissions stream to the writer inline.
	if ctx.replayLogger != nil {
		s.SetLogger(ctx.replayLogger)
	}
	return s
}

// captureWinningSeq records the winning permutation's attacker order + modes and pitch
// ordering into attackBufs scratch — the raw material for a verbatim cache replay.
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
	// Canonicalise ascending by card ID so lex-next-permutation visits each distinct
	// ordering exactly once — duplicates skip the redundant swap.
	sortAttackersByID(pcBuf, permMeta)
	sortPitchByID(pitchPerm, pitchVals)
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
				// Hand the superseded prior best back to the pool so the next perm's
				// preparePermState reuses it. promoteWinnerState already cloned its slices,
				// so the recycled state has independent backings.
				ctx.bufs.statePool.Put(bestWinner)
				bestWinner = winner
				ctx.captureWinningSeq(pcBuf, pitchPerm)
				ctx.promoteWinnerDeck(winner)
				ctx.promoteWinnerState(winner)
			} else {
				// Loser: state is no longer referenced — return to the pool so the next
				// perm reuses it instead of allocating fresh.
				ctx.bufs.statePool.Put(winner)
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
		// nextPermPitches leaves pitchPerm in lex-max state; reset to ascending so the
		// enumeration re-enters from the start for each attacker permutation.
		sortPitchByID(pitchPerm, pitchVals)
		tryPitchOrdering()
		for nextPermPitches(pitchPerm, pitchVals) {
			tryPitchOrdering()
		}
	}
	eval()
	for nextPermAttackers(pcBuf, permMeta) {
		eval()
	}
	return bestScore, bestWinner, foundLegal
}

// attackerKey packs (Card.ID, FromArsenal) into a uint32 — same-ID entries can still differ
// when played from arsenal (cost / rider changes), so the symmetry break must distinguish.
func attackerKey(pc *card.CardState) uint32 {
	k := uint32(pc.Card.ID()) << 1
	if pc.FromArsenal {
		k |= 1
	}
	return k
}

// sortAttackersByID insertion-sorts pcBuf and permMeta in lockstep, ascending by
// (Card.ID, FromArsenal). n is small (≤ 8 attackers), so insertion sort beats
// sort.Slice's closure-allocation overhead.
func sortAttackersByID(pcBuf []card.CardState, permMeta []*attackerMeta) {
	for i := 1; i < len(pcBuf); i++ {
		for j := i; j > 0 && attackerKey(&pcBuf[j]) < attackerKey(&pcBuf[j-1]); j-- {
			pcBuf[j-1], pcBuf[j] = pcBuf[j], pcBuf[j-1]
			permMeta[j-1], permMeta[j] = permMeta[j], permMeta[j-1]
		}
	}
}

// sortPitchByID is sortAttackersByID for the (pitchPerm, pitchVals) parallel slices.
func sortPitchByID(perm []card.Card, vals []int) {
	for i := 1; i < len(perm); i++ {
		for j := i; j > 0 && perm[j].ID() < perm[j-1].ID(); j-- {
			perm[j-1], perm[j] = perm[j], perm[j-1]
			vals[j-1], vals[j] = vals[j], vals[j-1]
		}
	}
}

// nextPermAttackers advances (pcBuf, permMeta) in lockstep to the lex-next permutation by
// attackerKey, returning false once the slice is in descending order. Equal-key entries
// skip the redundant swap, so duplicates yield each distinct ordering exactly once.
func nextPermAttackers(pcBuf []card.CardState, permMeta []*attackerMeta) bool {
	n := len(pcBuf)
	if n < 2 {
		return false
	}
	i := n - 2
	for i >= 0 && attackerKey(&pcBuf[i]) >= attackerKey(&pcBuf[i+1]) {
		i--
	}
	if i < 0 {
		return false
	}
	pivot := attackerKey(&pcBuf[i])
	j := n - 1
	for attackerKey(&pcBuf[j]) <= pivot {
		j--
	}
	pcBuf[i], pcBuf[j] = pcBuf[j], pcBuf[i]
	permMeta[i], permMeta[j] = permMeta[j], permMeta[i]
	for l, r := i+1, n-1; l < r; l, r = l+1, r-1 {
		pcBuf[l], pcBuf[r] = pcBuf[r], pcBuf[l]
		permMeta[l], permMeta[r] = permMeta[r], permMeta[l]
	}
	return true
}

// nextPermPitches is nextPermAttackers for the (pitchPerm, pitchVals) parallel slices.
func nextPermPitches(perm []card.Card, vals []int) bool {
	n := len(perm)
	if n < 2 {
		return false
	}
	i := n - 2
	for i >= 0 && perm[i].ID() >= perm[i+1].ID() {
		i--
	}
	if i < 0 {
		return false
	}
	j := n - 1
	for perm[j].ID() <= perm[i].ID() {
		j--
	}
	perm[i], perm[j] = perm[j], perm[i]
	vals[i], vals[j] = vals[j], vals[i]
	for l, r := i+1, n-1; l < r; l, r = l+1, r-1 {
		perm[l], perm[r] = perm[r], perm[l]
		vals[l], vals[r] = vals[r], vals[l]
	}
	return true
}

// playSequence builds permMeta and calls playSequenceWithMeta, recording the result on
// ctx.permState for the test-only PermEngine accessor. bestSequence's hot path threads the
// winner through return values instead and skips this write.
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

// seedChainEntry initialises one pcBuf slot: bind (Card, FromArsenal, Mode) and zero every
// ephemeral field. Mode is reseeded per modal tuple by the chain runner; initial 0 covers
// non-modal attackers.
func (ctx *sequenceContext) seedChainEntry(pc *card.CardState, c card.Card, idx int) {
	pc.Card = c
	pc.FromArsenal = idx == ctx.arsenalInIdx
	pc.Mode = 0
	pc.Ephemeral.Reset()
}

// playSequenceWithMeta runs the permutation currently held in ctx.bufs.pcBuf[:n] with
// aligned permMeta[:n]. Returns (damage, totalCounters, residualBudget, winner, legal).
// A nil winner indicates an infeasible permutation.
func (ctx *sequenceContext) playSequenceWithMeta(n int) (damage int, totalCounters int, residualBudget int, winner *gameengine.GameState, legal bool) {
	pcBuf := ctx.bufs.pcBuf
	ptrBuf := ctx.bufs.ptrBuf
	meta := ctx.bufs.permMeta[:n]
	for i := 0; i < n; i++ {
		pcBuf[i].Ephemeral.Reset()
	}
	played := ptrBuf[:n]
	state := ctx.preparePermState(played, n)
	state.SetIsMyTurn(true)
	ge := ctx.permEngine(state)
	// infeasible returns state to the pool and surfaces the legal=false signal. Used at
	// every early-return below — the caller (bestSequence) only handles Put for legal
	// outcomes (winner / loser), so the infeasible Put is owned here.
	infeasible := func() (int, int, int, *gameengine.GameState, bool) {
		ctx.bufs.statePool.Put(state)
		return 0, 0, 0, nil, false
	}
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
		// RemoveFromHand returns false when an earlier chain step moved this card out of
		// hand (e.g. DiscardToTopOfDeck alt cost). The partition planned against the
		// pre-chain hand; if the card is gone, reject so the optimiser doesn't credit a
		// phantom play.
		if !state.RemoveFromHand(pc.Card) {
			return infeasible()
		}
		prevPitchIdx := pool.idx
		contrib, ok := pool.pay(ge, m.costAt(ge, pc.Mode))
		if !ok {
			return infeasible()
		}
		pc.PitchedToPlay = contrib
		for k := prevPitchIdx; k < pool.idx; k++ {
			if !state.RemoveFromHand(pool.perm[k]) {
				return infeasible()
			}
		}
		if m.typesWithMode(pc.Mode).IsAttackReaction() {
			if m.hasPlayPrecondition {
				if !pc.Card.(card.PlayPrecondition).PlayPrecondition(ge, pc) {
					return infeasible()
				}
			}
			ar, ok := pc.Card.(card.AttackReaction)
			if !ok || activeAttack == nil || !ar.ARTargetAllowed(ge, activeAttack.Card, pc.Mode) {
				return infeasible()
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
		// Runs after finalizeActiveAttack so an earlier attack's on-hit rider has its chance
		// to set GrantedInstant. The free-chain-step check dispatches per-mode for
		// ModalTypes cards (Tip-Off mode 1 reads as Instant → 0 AP).
		if !m.isFreeChainStepWithMode(pc.Mode) && !pc.GrantedInstant {
			if state.ActionPoints() <= 0 {
				return infeasible()
			}
			state.AddActionPoints(-1)
		}
		if m.hasPlayPrecondition {
			if !pc.Card.(card.PlayPrecondition).PlayPrecondition(ge, pc) {
				return infeasible()
			}
		}

		state.SetCardsRemaining(played[i+1:])

		state.SetCurrentStepRerouted(false)
		// CardOrAbility fires before the card resolves so play-triggered effects land ahead
		// of the played card's own effect.
		ge.FireTriggers(triggertype.CardOrAbility, pc.Card)
		ge.ResolveChainStep(state.Logger(), pc)
		// ModalTypes cards (Tip-Off) read different is-attack / type-line values per Mode.
		// Resolve once and route the attack / non-attack-action / persistence checks off
		// the same TypeSet.
		modeTypes := m.typesWithMode(pc.Mode)
		if modeTypes.Has(card.TypeAttack) {
			// Mark is consumed only when the marked hero takes damage. A 0-power swing can't
			// strip the mark, so downstream chain steps can still read it.
			if pc.EffectiveAttack() > 0 {
				state.ClearOpponentMarked()
			}
			activeAttack = pc
		}
		state.AppendCardsPlayed(pc.Card)
		if modeTypes.IsNonAttackAction() {
			state.SetNonAttackActionPlayed(true)
		}
		if !modeTypes.PersistsInPlay() && !state.CurrentStepRerouted() {
			state.AppendGraveyard(pc.Card)
		}
		if pc.EffectiveGoAgain(ge) {
			state.AddActionPoints(1)
		}
	}
	finalizeActiveAttack()

	if pool.idx < pool.n {
		return infeasible()
	}
	ge.FireTriggers(triggertype.EndOfTurn, nil)
	return state.Value(), pendingTotalCountersFromState(state), pool.remaining, state, true
}

// pendingTotalCountersFromState sums Count over Auras + Items at end of chain — the
// partition's secondary tiebreaker. Counts pending aura fires and token stockpile at
// 1:1 weight, weaker than a real card in hand (see pendingTotalCardsFromState).
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

// pendingTotalCardsFromState projects cards available next turn: post-refill hand (held
// topped up to intellect) plus an occupied arsenal. Scoring post-refill rather than the bare
// end-of-chain hand lets the tiebreaker credit a chain that empties hand into attacks.
// Refill is uncapped, so a near-decked-out chain is scored slightly optimistically.
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
