package sim

// Attack-chain search: bestAttackWithWeapons evaluates one partition leaf across all phase /
// weapon masks, bestSequence picks the best ordering of attackers via Heap's algorithm, and
// playSequence* replay a single permutation through GameEngine while firing hero triggers,
// Aura handlers, and per-attack OnHit closures.

import (
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// perItemAbilityCap caps how many instances of one item's activated ability the chain
// runner enumerates per turn, bounding the wmask 2^k explosion when an item's Count gets
// large. Realistic counts in play tend to 1-3; 4 leaves headroom without letting a
// pathological hand blow up the per-leaf mask loop.
const perItemAbilityCap = 4

// FormatLogEntry renders a LogEntry into its display string. Chain entries with N=0 drop
// the "(+0)" suffix; trigger entries carry a "(from <source>)" tail.
func FormatLogEntry(e turnlogger.LogEntry) string {
	if e.Kind == turnlogger.LogEntryChainStep {
		if e.N == 0 {
			return e.Text
		}
		return fmt.Sprintf("%s (+%d)", e.Text, e.N)
	}
	if e.N == 0 {
		return fmt.Sprintf("%s (from %s)", e.Text, e.Source)
	}
	return fmt.Sprintf("%s (+%d) (from %s)", e.Text, e.N, e.Source)
}

// bestAttackWithWeapons enumerates phase / weapon masks for one partition leaf and returns
// the best (damage, futureValue, budget, swungWeapons, carryState, legal, cacheable) tuple.
func bestAttackWithWeapons(hero Hero, weapons []Weapon, attackers, defenders, pitched, held []card.Card, d *deck.Deck, bufs *attackBufs, mp Matchup, blockTotal, arsenalInIdx, arsenalDefenderIdx int, arsenalAtChainStart card.Card, prior gameengine.Spec, skipLog bool) (int, int, chainBudget, []string, CarryState, bool, bool) {
	priorAuras := auraSliceFromEngine(prior.Auras)
	priorItems := itemSliceFromEngine(prior.Items)
	runechantCarryover := auraCountByName(priorAuras, "Runechant")
	ctx := &sequenceContext{
		hero:                hero,
		pitched:             pitched,
		deck:                d,
		handStart:           held,
		arsenalAtChainStart: arsenalAtChainStart,
		bufs:                bufs,
		runechantCarryover:  runechantCarryover,
		matchup:             mp,
		blockTotal:          blockTotal,
		arsenalInIdx:        arsenalInIdx,
		priorAuras:          priorAuras,
		priorItems:          priorItems,
		priorOpponentMarked: prior.OpponentMarked,
		priorBanish:         prior.Banished,
		priorGraveyard:      prior.Graveyard,
		defenderAuras:       bufs.defenderAurasScratch[:0],
		defenders:           defenders,
		skipLog:             skipLog,
		cacheable:           true,
		carryWinner:         &bufs.carryWinnerScratch,
	}
	// Extend bufs.activatedAbilities with item ability instances for this Best call.
	abilities := bufs.activatedAbilities[:bufs.weaponAbilityCount]
	abilityCosts := bufs.activatedAbilityCosts[:bufs.weaponAbilityCount]
	for _, it := range priorItems {
		copies := it.Count()
		if copies > perItemAbilityCap {
			copies = perItemAbilityCap
		}
		ab := it.Ability()
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
	// Non-modal defender contribution is constant across phase / weapon masks.
	hasDRs := containsDefenseReaction(defenders)
	hasModalBlocker := containsModalBlocker(defenders)
	var defenseDealtConst int
	defenseCacheableConst := true
	if !hasModalBlocker && len(defenders) > 0 {
		defenseDealtConst, defenseCacheableConst = ctx.runDefense(defenders, pitched, d, priorAuras, mp.IncomingDamage, noBlockBudgetCap, arsenalDefenderIdx)
	}
	defenseDealt := defenseDealtConst
	defenseCacheable := defenseCacheableConst

	pitchedVals := bufs.pitchedValsScratch[:0]
	for _, c := range pitched {
		pitchedVals = append(pitchedVals, c.Pitch())
	}

	phaseCount := 1
	if hasDRs && len(pitched) > 0 {
		phaseCount = 1 << len(pitched)
	}

	attackersMinCost := 0
	attackersMaxCost := 0
	for _, a := range attackers {
		m := attackerMetaPtrFor(a)
		attackersMinCost += m.minCost
		attackersMaxCost += m.maxCost
	}

	copy(bufs.attackerBuf, attackers)

	bestDealt := 0
	bestFutureValue := 0
	var bestSwung []string
	var bestBudget chainBudget
	foundFeasible := false

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

		weaponBitsMask := (1 << len(weapons)) - 1
		totalAbilityMasks := 1 << len(ctx.activatedAbilities)
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
			dealt, futureValue, legal := ctx.bestSequence(allAttackers)
			if !legal {
				continue
			}
			// Cost the DRs against the prior-turn runechant carryover. Seed the DR scratch
			// engine with a runechant aura when the carryover is positive so variable-cost
			// DRs read RunechantCount() off this aura.
			drScratchAuras := bufs.drScratchAuras[:0]
			if ctx.runechantCarryover > 0 {
				drScratchAuras = append(drScratchAuras, NewRunechantAura(ctx.runechantCarryover))
			}
			bufs.drScratchAuras = drScratchAuras
			bufs.drScratch.Reset(gameengine.PermutationSeed{Auras: auraSliceAsEngine(drScratchAuras)})
			drCost := 0
			for _, d := range defenders {
				if !attackerMetaPtrFor(d).actsAsDR {
					continue
				}
				drCost += d.Cost(bufs.drScratch)
			}
			if drCost > phase.defendBudget {
				continue
			}
			if hasModalBlocker {
				defenseDealt, defenseCacheable = ctx.runDefense(defenders, pitched, d, priorAuras, mp.IncomingDamage, phase.defendBudget-drCost, arsenalDefenderIdx)
			}
			if phase.hasDefendPitches && phase.defendBudget-drCost >= phase.maxDefendPitch {
				continue
			}
			candidateDrawn := ctx.carryWinner.CardsDrawn
			bestDrawn := bufs.bestCarryScratch.CardsDrawn
			candidateHand := len(ctx.carryWinner.Hand)
			bestHand := len(bufs.bestCarryScratch.Hand)
			cmp := chainScoreCmp(dealt, candidateDrawn, futureValue, bestDealt, bestDrawn, bestFutureValue)
			if !foundFeasible || cmp > 0 || (cmp == 0 && candidateHand > bestHand) {
				bestDealt = dealt
				bestFutureValue = futureValue
				bestSwung = bufs.weaponNames[wmask&weaponBitsMask]
				bestBudget = chainBudget{resource: phase.attackBudget, maxPitch: phase.maxAttackPitch, hasAttackPitches: phase.hasAttackPitches}
				bufs.bestCarryScratch.CopyFrom(ctx.carryWinner)
				foundFeasible = true
			}
		}
	}

	_ = attackersMaxCost
	if !foundFeasible {
		return 0, 0, chainBudget{}, nil, CarryState{}, false, defenseCacheable
	}
	return bestDealt, defenseDealt, bestBudget, bestSwung, bufs.bestCarryScratch, true, ctx.cacheable && defenseCacheable
}

// sequenceContext carries the stable per-partition-leaf environment.
type sequenceContext struct {
	hero                  Hero
	pitched               []card.Card
	deck                  *deck.Deck
	handStart             []card.Card
	arsenalAtChainStart   card.Card
	bufs                  *attackBufs
	attackPitchPerm       []card.Card
	attackPitchVals       []int
	resourceBudget        int
	runechantCarryover    int
	matchup               Matchup
	blockTotal            int
	arsenalInIdx          int
	priorAuras            []*Aura
	priorItems            []*Item
	priorOpponentMarked   bool
	priorBanish           []card.Card
	priorGraveyard        []card.Card
	activatedAbilities    []card.Card
	activatedAbilityCosts []int
	defenderAuras         []*Aura
	defenders             []card.Card
	carryWinner           *CarryState
	skipLog               bool
	cacheable             bool
}

// activeLogger returns the *turnlogger.TurnLogger threaded into each permutation — nil
// during the find-best pass so cards' log helpers short-circuit, the bufs-owned recorder
// during replay.
func (ctx *sequenceContext) activeLogger() *turnlogger.TurnLogger {
	if ctx.skipLog {
		return nil
	}
	return ctx.bufs.logger
}

// runDefense seeds the engine's aura set with priorAuras, runs defendersDamage, and
// captures the post-defense aura set into ctx.defenderAuras for later chain seeding.
func (ctx *sequenceContext) runDefense(defenders, pitched []card.Card, deckPile *deck.Deck, priorAuras []*Aura, incomingDamage, blockBudget, arsenalDefenderIdx int) (int, bool) {
	bufs := ctx.bufs
	bufs.state.SetAuras(auraSliceAsEngine(append(ctx.defenderAuras[:0], priorAuras...)))
	bufs.state.SetLogger(ctx.activeLogger())
	dealt, gravScratch, cacheable := defendersDamage(defenders, pitched, deckPile, bufs.state, bufs.defenseGravScratch, &bufs.drCardStateScratch, incomingDamage, blockBudget, arsenalDefenderIdx)
	bufs.defenseGravScratch = gravScratch
	ctx.defenderAuras = auraSliceFromEngine(bufs.state.Auras())
	bufs.defenderAurasScratch = ctx.defenderAuras
	return dealt, cacheable
}

// resetStateForPermutation rewrites every engine field to its per-permutation starting
// value by building a gameengine.PermutationSeed and calling engine.Reset. The engine owns
// its own slice backings, so an unmodified permutation never reallocates.
func (ctx *sequenceContext) resetStateForPermutation() {
	bufs := ctx.bufs
	auraSeed := ctx.priorAuras
	if len(ctx.defenderAuras) > 0 {
		auraSeed = ctx.defenderAuras
	}
	// Seed graveyard with prior entries first, defenders second (the partition has already
	// dropped defenders into the graveyard for this leaf).
	graveSeed := append(ctx.bufs.defenseGravScratch[:0], ctx.priorGraveyard...)
	_ = graveSeed // capture is incidental
	var logger *turnlogger.TurnLogger
	if l := ctx.activeLogger(); l != nil {
		l.SetBuffer(bufs.logBacking)
		logger = l
	}
	bufs.state.Reset(gameengine.PermutationSeed{
		Hand:                 ctx.handStart,
		Deck:                 ctx.deck,
		Arsenal:              ctx.arsenalAtChainStart,
		Graveyard:            ctx.priorGraveyard,
		Banished:             ctx.priorBanish,
		OpponentMarked:       ctx.priorOpponentMarked,
		Auras:                auraSliceAsEngine(auraSeed),
		Items:                itemSliceAsEngine(ctx.priorItems),
		Pitched:              ctx.pitched,
		Defenders:            nil,
		IncomingDamage:       ctx.matchup.IncomingDamage,
		ArcaneIncomingDamage: ctx.matchup.ArcaneIncomingDamage,
		BlockTotal:           ctx.blockTotal,
		Logger:               logger,
	})
	// Re-append defenders to the graveyard so chain cards that scan it see this turn's
	// defenders alongside prior entries.
	for _, c := range ctx.defenders {
		bufs.state.AppendGraveyard(c)
	}
}

// bestSequence tries every ordering of attackers and returns the max total damage plus
// the pendingFutureValue at the end of the winning permutation. legal=true when at least
// one ordering is playable.
func (ctx *sequenceContext) bestSequence(attackers []card.Card) (int, int, bool) {
	n := len(attackers)
	if n == 0 {
		if len(ctx.attackPitchPerm) > 0 {
			return 0, 0, false
		}
		ctx.resetStateForPermutation()
		st := ctx.bufs.state
		if st.HasEndOfTurnFire() {
			st.FireEndOfTurn()
		}
		ctx.carryWinner.SnapshotFromTurn(st)
		return 0, pendingFutureValue(auraSliceFromEngine(st.Auras()), itemSliceFromEngine(st.Items())), true
	}
	pcBuf := ctx.bufs.pcBuf[:n]
	permMeta := ctx.bufs.permMeta[:n]
	for idx, c := range attackers {
		permMeta[idx] = attackerMetaPtrFor(c)
		ctx.seedChainEntry(&pcBuf[idx], c, idx)
	}

	best := 0
	bestFutureValue := 0
	foundLegal := false
	ctx.carryWinner.Reset()
	state := ctx.bufs.state
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
		if !hasModal {
			dmg, futureValue, _, legal := ctx.playSequenceWithMeta(n)
			if ctx.cacheable && !state.IsCacheable() {
				ctx.cacheable = false
			}
			if !legal {
				return
			}
			cmp := chainScoreCmp(dmg, state.CardsDrawn(), futureValue, best, bestCardsDrawn, bestFutureValue)
			if !foundLegal || cmp > 0 {
				best = dmg
				bestCardsDrawn = state.CardsDrawn()
				bestFutureValue = futureValue
				foundLegal = true
				ctx.carryWinner.SnapshotFromTurn(state)
			}
			return
		}
		for tuple := 0; tuple < tupleCount; tuple++ {
			rem := tuple
			for i := 0; i < n; i++ {
				modes := int(permMeta[i].modes)
				pcBuf[i].Mode = int8(rem % modes)
				rem /= modes
			}
			dmg, futureValue, _, legal := ctx.playSequenceWithMeta(n)
			if ctx.cacheable && !state.IsCacheable() {
				ctx.cacheable = false
			}
			if !legal {
				continue
			}
			cmp := chainScoreCmp(dmg, state.CardsDrawn(), futureValue, best, bestCardsDrawn, bestFutureValue)
			if !foundLegal || cmp > 0 {
				best = dmg
				bestCardsDrawn = state.CardsDrawn()
				bestFutureValue = futureValue
				foundLegal = true
				ctx.carryWinner.SnapshotFromTurn(state)
			}
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
	return best, bestFutureValue, foundLegal
}

// playSequence is a thin wrapper that builds permMeta and calls playSequenceWithMeta.
// The hot path (bestSequence) builds meta once and calls playSequenceWithMeta directly.
func (ctx *sequenceContext) playSequence(order []card.Card) (damage int, futureValue int, residualBudget int, legal bool) {
	n := len(order)
	pcBuf := ctx.bufs.pcBuf
	meta := ctx.bufs.permMeta[:n]
	for i, c := range order {
		meta[i] = attackerMetaPtrFor(c)
		ctx.seedChainEntry(&pcBuf[i], c, i)
	}
	return ctx.playSequenceWithMeta(n)
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
// aligned permMeta[:n]. Returns (damage, futureValue, residualBudget, legal).
func (ctx *sequenceContext) playSequenceWithMeta(n int) (damage int, futureValue int, residualBudget int, legal bool) {
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
	ctx.resetStateForPermutation()
	state := ctx.bufs.state
	state.SetHero(ctx.hero)
	// Seed engine.hand with the upcoming chain attackers + the pitched bag so each chain
	// step's Play sees an accurate "in hand right now" snapshot.
	for k := 0; k < n; k++ {
		state.AppendHandRaw(played[k].Card)
	}
	for _, c := range ctx.attackPitchPerm {
		state.AppendHandRaw(c)
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
		if gameengine.LikelyToHit(activeAttack) {
			for i := range activeAttack.OnHit {
				h := &activeAttack.OnHit[i]
				h.Fire(state, state.Logger(), activeAttack, h)
			}
			// Drain matching TriggerHit Triggers.
			types := activeAttack.Card.Types(nil)
			prevTriggering := state.TriggeringCard()
			state.SetTriggeringCard(activeAttack.Card)
			state.FireHit(types)
			state.SetTriggeringCard(prevTriggering)
		}
		activeAttack = nil
	}
	for i, pc := range played {
		m := meta[i]
		if !m.isFreeChainStep {
			if state.ActionPoints() <= 0 {
				return 0, 0, 0, false
			}
			state.AddActionPoints(-1)
		}
		// Remove the playing card from engine.hand before resolving.
		state.RemoveFromHand(pc.Card)
		prevPitchIdx := pool.idx
		contrib, ok := pool.pay(m.costAt(state, pc.Mode))
		if !ok {
			return 0, 0, 0, false
		}
		pc.PitchedToPlay = contrib
		for k := prevPitchIdx; k < pool.idx; k++ {
			state.RemoveFromHand(pool.perm[k])
		}
		if m.types.IsAttackReaction() {
			if pre, ok := pc.Card.(card.PlayPrecondition); ok {
				if !pre.PlayPrecondition(state, pc) {
					return 0, 0, 0, false
				}
			}
			ar, ok := pc.Card.(AttackReaction)
			if !ok || activeAttack == nil || !ar.ARTargetAllowed(state, activeAttack.Card, pc.Mode) {
				return 0, 0, 0, false
			}
			ctx.hero.OnCardPlayed(pc.Card, state, state.Logger())
			state.SetAttackReactionTarget(activeAttack)
			state.ResolveChainStep(state.Logger(), pc)
			state.SetAttackReactionTarget(nil)
			state.SetCardsPlayed(append(state.CardsPlayed(), pc.Card))
			state.AppendGraveyard(pc.Card)
			if pc.EffectiveGoAgain(state) {
				state.AddActionPoints(1)
			}
			continue
		}

		finalizeActiveAttack()
		if pre, ok := pc.Card.(card.PlayPrecondition); ok {
			if !pre.PlayPrecondition(state, pc) {
				return 0, 0, 0, false
			}
		}

		state.SetCardsRemaining(played[i+1:])

		ctx.hero.OnCardPlayed(pc.Card, state, state.Logger())
		state.SetCurrentStepRerouted(false)
		state.ResolveChainStep(state.Logger(), pc)
		if m.isAttack {
			state.FireAttack(pc.Card)
			state.ClearOpponentMarked()
		}
		if m.isAttackAction {
			state.FireAttackAction(pc.Card)
		}
		if m.isAttack {
			activeAttack = pc
		}
		state.SetCardsPlayed(append(state.CardsPlayed(), pc.Card))
		if m.types.IsNonAttackAction() {
			state.SetNonAttackActionPlayed(true)
		}
		if !m.types.PersistsInPlay() && !state.CurrentStepRerouted() {
			state.AppendGraveyard(pc.Card)
		}
		if pc.EffectiveGoAgain(state) {
			state.AddActionPoints(1)
		}
	}
	finalizeActiveAttack()

	if pool.idx < pool.n {
		return 0, 0, 0, false
	}
	if state.HasEndOfTurnFire() {
		state.FireEndOfTurn()
	}
	return state.Value(), pendingFutureValue(auraSliceFromEngine(state.Auras()), itemSliceFromEngine(state.Items())), pool.remaining, true
}

// pendingFutureValue sums the Count of every Aura plus every Item at end of chain — the
// partition tiebreaker's "hidden later-turn payoff" signal.
func pendingFutureValue(auras []*Aura, items []*Item) int {
	total := 0
	for _, a := range auras {
		total += a.Count()
	}
	for _, it := range items {
		total += it.Count()
	}
	return total
}
