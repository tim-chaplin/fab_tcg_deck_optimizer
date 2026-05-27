package sim

// Attack-turn search: bestAttackWithWeapons evaluates one partition leaf across all
// phase / weapon masks; bestSequence picks the best attacker ordering; playSequence*
// replays one permutation against a fresh per-permutation GameState copy while firing
// hero triggers, Aura handlers, and per-attack OnHit closures.
//
// State lifecycle:
//   - findBest builds a master *GameState once per Best call.
//   - evaluatePartition copies the master into a per-leaf state, runs defense, then
//     enumerates attack-turn permutations against fresh leafState copies so per-permutation
//     mutations stay isolated. The winning copy's *GameState is the partition's result.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
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
	attackers, defenders []card.Card,
	pitched []*card.CardState,
	held []card.Card,
	d *deck.Deck,
	bufs *attackBufs,
	blockTotal, arsenalInIdx int,
	arsenalAtAttackTurnStart card.Card,
) *sequenceContext {
	ctx := bufs.pooledSequenceCtx
	if ctx == nil {
		ctx = &sequenceContext{}
		bufs.pooledSequenceCtx = ctx
	}
	*ctx = sequenceContext{
		pitched:                  pitched,
		attackers:                attackers,
		deck:                     d,
		handStart:                held,
		arsenalAtAttackTurnStart: arsenalAtAttackTurnStart,
		bufs:                     bufs,
		runechantCarryover:       masterState.RunechantCount(),
		blockTotal:               blockTotal,
		arsenalInIdx:             arsenalInIdx,
		priorOpponentMarked:      masterState.OpponentMarked(),
		priorBanish:              masterState.Banished(),
		priorGraveyard:           masterState.Graveyard(),
		defenders:                defenders,
		startOfTurnValue:         masterState.Value(),
		cacheable:                true,
	}
	abilities := bufs.activatedAbilities[:bufs.weaponAbilityCount]
	abilityCosts := bufs.activatedAbilityCosts[:bufs.weaponAbilityCount]
	appendAbilities := func(it gameengine.Item) {
		ability := it.Ability()
		if ability == nil {
			// Triggered items (Talisman of Recompense) have no activated ability — they fire
			// through FireTriggers, with nothing to enqueue as a playable.
			return
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
	for _, it := range masterState.Items() {
		appendAbilities(it)
	}
	masterState.ForEachTokenItem(func(it gameengine.Item) {
		appendAbilities(it)
	})
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

// installLeafDeck rebinds leafState's owned *deck.Deck wrapper to alias master deck d.
// Run before runDefense so DR Plays (Rise Above's PrependToDeck, an Opt-ing DR, ...)
// mutate the leaf-scoped wrapper rather than the master shared across leaves.
func installLeafDeck(ctx *sequenceContext, _ *attackBufs, d *deck.Deck) {
	ctx.leafState.Deck().ShallowCopyFrom(d)
	ctx.deck = ctx.leafState.Deck()
}

// bestAttackWithWeapons enumerates phase / weapon masks for one partition leaf and
// returns the best (damage, defenseDealt, budget, swungWeapons, winnerState, legal,
// cacheable) tuple. Each per-leaf state branches off via masterState.Copy().
func bestAttackWithWeapons(
	masterState *gameengine.GameState,
	weapons []weapon.Weapon,
	attackers, defenders []card.Card,
	pitched []*card.CardState,
	held []card.Card,
	d *deck.Deck,
	bufs *attackBufs,
	blockTotal, arsenalInIdx, arsenalDefenderIdx int,
	arsenalAtAttackTurnStart card.Card,
) (int, int, attackTurnBudget, []string, *gameengine.GameState, bool, bool) {
	ctx := newSequenceContext(masterState, weapons, attackers, defenders, pitched, held, d, bufs, blockTotal, arsenalInIdx, arsenalAtAttackTurnStart)
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
		// goes through the DamageAboutToBeTaken → DamageTaken sequence.
		defenseDealtConst = ctx.fireUndefendedDamageTriggers()
	}
	defenseDealt := defenseDealtConst
	defenseCacheable := defenseCacheableConst
	// Stay paired with bestWinner / sol.defenders; the loop-scoped defenseDealt is
	// clobbered every modal-blocker pmask iteration.
	bestDefenseDealt := defenseDealtConst
	bestDefenseCacheable := defenseCacheableConst

	pitchedVals := bufs.pitchedValsScratch[:0]
	for _, c := range pitched {
		pitchedVals = append(pitchedVals, c.Card.Pitch())
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

	var bestScore, bestTotalScore attackTurnScore
	var bestSwung []string
	var bestBudget attackTurnBudget
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
			// Defense resolves before the attack turn. A modal blocker's block depends on
			// phase.defendBudget, so the defense pass runs once per phase here.
			installLeafDeck(ctx, bufs, d)
			defenseDealt, defenseCacheable, ctx.handStart = ctx.runDefense(defenders, pitched, held, ctx.deck, incoming, phase.defendBudget-drCost, arsenalDefenderIdx, nil)
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
				bestBudget = attackTurnBudget{resource: phase.attackBudget, maxPitch: phase.maxAttackPitch, hasAttackPitches: phase.hasAttackPitches}
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
		return 0, 0, attackTurnBudget{}, nil, nil, false, defenseCacheable
	}
	return bestScore.value, bestDefenseDealt, bestBudget, bestSwung, bestWinner, true, ctx.cacheable && bestDefenseCacheable
}

// drCostProbe returns the pooled *GameEngine with its Runechant token slot set to the
// given count for variable-cost DR cost probing. DRs read only RunechantCount() to decide
// Cost. The engine is lazily built once and reused across probes; per call just rewrites
// the slot's count via SetRunechantCount.
func (ctx *sequenceContext) drCostProbe(runechants int) *gameengine.GameEngine {
	bufs := ctx.bufs
	ge := bufs.pooledDRCostProbe
	if ge == nil {
		ge = gameengine.New()
		bufs.pooledDRCostProbe = ge
	}
	ge.GameState.SetRunechantCount(runechants)
	return ge
}

// sequenceContext carries the stable per-partition-leaf environment.
type sequenceContext struct {
	pitched                  []*card.CardState
	attackers                []card.Card
	deck                     *deck.Deck
	handStart                []card.Card
	arsenalAtAttackTurnStart card.Card
	bufs                     *attackBufs
	attackPitchPerm          []*card.CardState
	attackPitchVals          []int
	resourceBudget           int
	runechantCarryover       int
	blockTotal               int
	arsenalInIdx             int
	priorOpponentMarked      bool
	priorBanish              []card.Card
	priorGraveyard           []card.Card
	activatedAbilities       []card.Card
	activatedAbilityCosts    []int
	defenders                []card.Card
	leafState                *gameengine.GameState
	// startOfTurnValue is masterState.Value() captured at construction and re-seeded into each
	// per-perm state after ResetEphemeralState — attack-turn accumulators ride on top of the
	// start-of-action-phase aura tick, so summary.Value includes that baseline.
	startOfTurnValue int
	cacheable        bool
	// replayLogger, when non-nil, is installed on each per-perm state so log emissions
	// stream inline. PrintBestTurn sets it to a stdout StreamLogger; the eval hot path
	// leaves it nil so the state's NoopLogger keeps emissions free.
	replayLogger card.Logger
	// permState records the last *GameState playSequence ran the attack turn against, so the
	// test-only PermEngine accessor can read it. The hot bestSequence path threads the
	// winner through return values and leaves this nil.
	permState *gameengine.GameState
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

// fireUndefendedDamageTriggers runs the DamageAboutToBeTaken → DamageTaken sequence
// against ctx.leafState on the no-defender path: DamageAboutToBeTaken subscribers may
// call PreventIncomingDamage to absorb the swing before DamageTaken fires. Zeroes
// state.Value before the fires so handler AddValues are captured cleanly into the
// returned defense-dealt total, then restores ctx.leafState.Value to its prior level.
func (ctx *sequenceContext) fireUndefendedDamageTriggers() int {
	ctx.leafState.SetIsMyTurn(false)
	leafEngine := ctx.permEngine(ctx.leafState)
	ctx.leafState.SetValue(0)
	if ctx.leafState.RemainingUnblockedDamage() > 0 {
		leafEngine.FireTriggers(triggertype.DamageAboutToBeTaken, nil)
	}
	if ctx.leafState.RemainingUnblockedDamage() > 0 {
		leafEngine.FireTriggers(triggertype.DamageTaken, nil)
	}
	return ctx.leafState.Value()
}

// runDefense mutates ctx.leafState through the defender list, accumulating per-DR Value.
// Auras grow with any DR-added entries; graveyard is left as priorGraveyard + defenders for
// the attack turn phase. Per-permutation per-turn-locals reset via ResetEphemeralState, so runDefense
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
func (ctx *sequenceContext) runDefense(defenders []card.Card, pitched []*card.CardState, held []card.Card, deckPile *deck.Deck, matchupIncomingDamage, blockBudget, arsenalDefenderIdx int, cachedModes []playedCard) (int, bool, []card.Card) {
	state := ctx.leafState
	state.SetIsMyTurn(false)
	if ctx.replayLogger != nil {
		state.SetLogger(ctx.replayLogger)
	}
	state.SetDeck(deckPile)
	state.SetIncomingDamage(matchupIncomingDamage)
	// Baseline leafState's graveyard to priorGraveyard. Defending cards move to
	// graveyard only when the attack turn closes (the post-block append below).
	state.SetGraveyard(append(state.Graveyard()[:0], ctx.priorGraveyard...))
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
		defenseHand = append(defenseHand, card.CardState{Card: c.Card, Role: card.Pitch})
	}
	state.SetHandStates(defenseHand)

	// DR loop: graveyard reads see only priorGraveyard.
	for i, def := range defenders {
		if !attackerMetaPtrFor(def).actsAsDR {
			continue
		}
		ctx.bufs.defModes[i] = playedCard{card: def}
		state.SetPitched(pitched)
		state.SetDefenders(defenders)
		state.SetValue(0)
		state.SetCacheable(true)
		*cs = card.CardState{Card: def, FromArsenal: i == arsenalDefenderIdx}
		ge.ResolveAttackStep(state.Logger(), cs)
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
	// Discards already entered graveyard via ge.Discard mid-block-loop.
	discarded := len(origHeld) - (len(state.HandStates()) - len(ctx.attackers) - len(pitched))
	survivingHeld := origHeld[discarded:]

	// Combat chain closes — all defending cards move to graveyard now, simultaneously.
	for _, def := range defenders {
		state.AppendGraveyard(def)
	}

	// Defense phase over: fire DamageAboutToBeTaken while the unblocked figure is still
	// mutable — subscribers may call PreventIncomingDamage to absorb. Then DamageTaken
	// fires only if damage still gets through. Zero state.Value before the fires so any
	// AddValue from the handlers is captured cleanly into total.
	state.SetValue(0)
	if state.RemainingUnblockedDamage() > 0 {
		ge.FireTriggers(triggertype.DamageAboutToBeTaken, nil)
	}
	if state.RemainingUnblockedDamage() > 0 {
		ge.FireTriggers(triggertype.DamageTaken, nil)
	}
	total += state.Value()

	return total, cacheable, survivingHeld
}

// preparePermState returns a fresh per-permutation *GameState for the attack turn run. The state
// inherits leafState's post-defense auras / items / graveyard / banished / hero / arsenal;
// ResetEphemeralState wipes the previous perm's play state, then this perm's inputs install.
// Hand = attack-turn attackers + attack-phase pitched bag, so Hand() reads see the upcoming bag.
//
// IncomingDamage is not re-installed: the matchup figure rode in constant on leafState and
// ResetEphemeralState zeroed the damage-blocked accumulator.
func (ctx *sequenceContext) preparePermState(playedAttackers []*card.CardState, n int) *gameengine.GameState {
	bufs := ctx.bufs
	s := bufs.statePool.Get()
	s.CopyPersistentStateFrom(ctx.leafState)
	s.ResetEphemeralState()
	s.SetValue(ctx.startOfTurnValue)
	// Hero / opponentMarked already mirror leafState (which mirrors masterState). Only
	// arsenal and blockTotal need explicit setting: arsenal may have been promoted out of
	// the leaf via findArsenalCard, and blockTotal is zeroed by ResetEphemeralState.
	s.SetArsenal(ctx.arsenalAtAttackTurnStart)
	s.SetBlockTotal(ctx.blockTotal)
	// graveyard / banished arrived via CopyPersistentStateFrom, copied into s's own
	// prewarmed backing — attack-turn runner appends mutate this slot's storage only. The deck
	// wrapper s owns is rebound in place to alias ctx.deck's backing.
	s.Deck().ShallowCopyFrom(ctx.deck)
	// Build the per-perm hand into s's prewarmed backing. A cap shortfall means the
	// workload outgrew gameengine's defaultHandCap — bump it rather than fall back to
	// per-call alloc.
	needed := len(ctx.handStart) + n + len(ctx.attackPitchPerm)
	hand := s.HandStates()
	if cap(hand) < needed {
		panic("sim: pooled state hand backing too small; raise gameengine.defaultHandCap")
	}
	hand = hand[:0]
	for _, c := range ctx.handStart {
		hand = append(hand, card.CardState{Card: c, Role: card.Held})
	}
	for k := 0; k < n; k++ {
		hand = append(hand, card.CardState{Card: playedAttackers[k].Card, Role: card.Attack})
	}
	for _, c := range ctx.attackPitchPerm {
		hand = append(hand, card.CardState{Card: c.Card, Role: card.Pitch})
	}
	s.SetHandStates(hand)
	// CardsPlayed is similarly prewarmed; the cap check guards mid-attack-turn growth.
	cpNeeded := n + len(ctx.attackPitchPerm)
	cp := s.CardsPlayed()
	if cap(cp) < cpNeeded {
		panic("sim: pooled state cardsPlayed backing too small; raise gameengine.defaultCardsPlayedCap")
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
func (ctx *sequenceContext) captureWinningSeq(pcBuf []card.CardState, pitchPerm []*card.CardState) {
	b := ctx.bufs
	b.seqAttack = b.seqAttack[:0]
	for i := range pcBuf {
		b.seqAttack = append(b.seqAttack, playedCard{
			card: pcBuf[i].Card, mode: pcBuf[i].Mode, fromArsenal: pcBuf[i].FromArsenal,
		})
	}
	b.seqPitch = b.seqPitch[:0]
	for _, pc := range pitchPerm {
		b.seqPitch = append(b.seqPitch, pc.Card)
	}
}

// bestSequence tries every ordering of attackers and returns the winning permutation's
// attackTurnScore. legal=true when at least one ordering is playable. Returns the winning
// *GameState via the second return value.
func (ctx *sequenceContext) bestSequence(attackers []card.Card) (attackTurnScore, *gameengine.GameState, bool) {
	n := len(attackers)
	if n == 0 {
		if len(ctx.attackPitchPerm) > 0 {
			return attackTurnScore{}, nil, false
		}
		emptyAttackers := ctx.bufs.ptrBuf[:0]
		permState := ctx.preparePermState(emptyAttackers, 0)
		ge := ctx.permEngine(permState)
		ge.FireTriggers(triggertype.EndOfTurn, nil)
		ctx.captureWinningSeq(nil, nil)
		// permState.Value() carries the seeded baseline plus any EndOfTurn fire delta.
		return attackTurnScoreOf(permState, permState.Value()), permState, true
	}
	pcBuf := ctx.bufs.pcBuf[:n]
	permMeta := ctx.bufs.permMeta[:n]
	for idx, c := range attackers {
		permMeta[idx] = attackerMetaPtrFor(c)
		ctx.seedAttackStepEntry(&pcBuf[idx], c, idx)
	}

	var bestScore attackTurnScore
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
			score := attackTurnScoreOf(winner, dmg)
			if !foundLegal || score.cmp(bestScore) > 0 {
				bestScore = score
				foundLegal = true
				// Hand the superseded prior best back to the pool so the next perm's
				// preparePermState reuses it; each pool slot owns its own graveyard /
				// banished / hand / deck backings, so no per-promotion clone is needed.
				ctx.bufs.statePool.Put(bestWinner)
				bestWinner = winner
				ctx.captureWinningSeq(pcBuf, pitchPerm)
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
// Pitch entries are *CardState; comparison reads through .Card.ID().
func sortPitchByID(perm []*card.CardState, vals []int) {
	for i := 1; i < len(perm); i++ {
		for j := i; j > 0 && perm[j].Card.ID() < perm[j-1].Card.ID(); j-- {
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
// Pitch entries are *CardState; comparison reads through .Card.ID().
func nextPermPitches(perm []*card.CardState, vals []int) bool {
	n := len(perm)
	if n < 2 {
		return false
	}
	i := n - 2
	for i >= 0 && perm[i].Card.ID() >= perm[i+1].Card.ID() {
		i--
	}
	if i < 0 {
		return false
	}
	j := n - 1
	for perm[j].Card.ID() <= perm[i].Card.ID() {
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
		ctx.seedAttackStepEntry(&pcBuf[i], c, i)
	}
	d, tc, rb, winner, lg := ctx.playSequenceWithMeta(n)
	ctx.permState = winner
	return d, tc, rb, lg
}

// playSequenceModal is playSequence for a cache replay: it seeds each attack step from a
// playedCard, applying the cached modal Mode rather than the default 0. The winning state
// lands on ctx.permState.
func (ctx *sequenceContext) playSequenceModal(order []playedCard) (damage int, totalCounters int, residualBudget int, legal bool) {
	n := len(order)
	pcBuf := ctx.bufs.pcBuf
	meta := ctx.bufs.permMeta[:n]
	for i, pc := range order {
		meta[i] = attackerMetaPtrFor(pc.card)
		ctx.seedAttackStepEntry(&pcBuf[i], pc.card, i)
		pcBuf[i].Mode = pc.mode
		pcBuf[i].FromArsenal = pc.fromArsenal
	}
	d, tc, rb, winner, lg := ctx.playSequenceWithMeta(n)
	ctx.permState = winner
	return d, tc, rb, lg
}

// seedAttackStepEntry initialises one pcBuf slot: bind (Card, FromArsenal, Mode) and zero every
// ephemeral field. Mode is reseeded per modal tuple by the attack-turn runner; initial 0 covers
// non-modal attackers.
func (ctx *sequenceContext) seedAttackStepEntry(pc *card.CardState, c card.Card, idx int) {
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
		// Self-gating — no-op on a miss. Runs before the OnHit / FireTriggers block so
		// those handlers see the updated HitThisTurn flag + DamageDealt count.
		state.RegisterPhysicalDamage(activeAttack.EffectiveAttack(), activeAttack.EffectiveDominate())
		if hit {
			for i := range activeAttack.OnHit {
				h := &activeAttack.OnHit[i]
				h.Fire(ge, state.Logger(), activeAttack, h)
			}
			ge.FireTriggers(triggertype.Hit, activeAttack)
		}
		activeAttack = nil
	}
	for i, pc := range played {
		m := meta[i]
		// RemoveFromHand returns false when an earlier attack step moved this card out of
		// hand (e.g. DiscardToTopOfDeck alt cost). The partition planned against the
		// pre-attack-turn hand; if the card is gone, reject so the optimiser doesn't credit a
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
			if !state.RemoveFromHand(pool.perm[k].Card) {
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
			ge.FireTriggers(triggertype.CardOrAbility, pc)
			state.SetAttackReactionTarget(activeAttack)
			ge.ResolveAttackStep(state.Logger(), pc)
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
		// to set GrantedInstant. The free-attack-step check dispatches per-mode for
		// ModalTypes cards (Tip-Off mode 1 reads as Instant → 0 AP).
		if !m.isFreeAttackStepWithMode(pc.Mode) && !pc.GrantedInstant {
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
		ge.FireTriggers(triggertype.CardOrAbility, pc)
		ge.ResolveAttackStep(state.Logger(), pc)
		// ModalTypes cards (Tip-Off) read different is-attack / type-line values per Mode.
		// Resolve once and route the attack / non-attack-action / persistence checks off
		// the same TypeSet.
		modeTypes := m.typesWithMode(pc.Mode)
		if modeTypes.Has(card.TypeAttack) {
			// Mark is consumed only when the marked hero takes damage. A 0-power swing can't
			// strip the mark, so downstream attack steps can still read it.
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

// pendingTotalCountersFromState sums Count over Auras + Items (card-backed + token
// slots) at end of attack turn — the partition's secondary tiebreaker. Counts pending
// aura fires and token stockpile at 1:1 weight, weaker than a real card in hand (see
// pendingTotalCardsFromState).
func pendingTotalCountersFromState(gs *gameengine.GameState) int {
	if gs == nil {
		return 0
	}
	return gs.TotalAuraCount() + gs.TotalItemCount()
}

// attackTurnScoreOf builds a leaf's attackTurnScore from its end-of-attack-turn winner state and the
// damage / block value credited this turn.
func attackTurnScoreOf(winner *gameengine.GameState, value int) attackTurnScore {
	return attackTurnScore{
		value:         value,
		cardsPlayed:   len(winner.CardsPlayed()),
		totalCards:    pendingTotalCardsFromState(winner),
		totalCounters: pendingTotalCountersFromState(winner),
	}
}

// pendingTotalCardsFromState projects cards available next turn: post-refill hand (held
// topped up to intellect) plus an occupied arsenal. Scoring post-refill rather than the bare
// end-of-attack-turn hand lets the tiebreaker credit an attack turn that empties hand into attacks.
// Refill is uncapped, so a near-decked-out attack turn is scored slightly optimistically.
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
