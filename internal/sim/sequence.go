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
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/item"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/turnlogger"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// perItemAbilityCap caps how many instances of one item's activated ability the chain
// runner enumerates per turn, bounding the wmask 2^k explosion when an item's Count
// gets large. Realistic counts in play tend to 1-3; 4 leaves headroom without letting a
// pathological hand blow up the per-leaf mask loop.
const perItemAbilityCap = 4

// FormatLogEntry renders a LogEntry into its display string. Chain entries with N=0
// drop the "(+0)" suffix; trigger entries carry a "(from <source>)" tail.
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

// bestAttackWithWeapons enumerates phase / weapon masks for one partition leaf and
// returns the best (damage, futureValue, budget, swungWeapons, winnerState, legal,
// cacheable) tuple. masterState holds the start-of-turn carryover (hero, arsenal, auras,
// items, banished, graveyard, opponentMarked) the chain runner reads from; each per-leaf
// state branches off via masterState.Copy().
func bestAttackWithWeapons(
	masterState *gameengine.GameState,
	weapons []weapon.Weapon,
	attackers, defenders, pitched, held []card.Card,
	d *deck.Deck,
	bufs *attackBufs,
	mp Matchup,
	blockTotal, arsenalInIdx, arsenalDefenderIdx int,
	arsenalAtChainStart card.Card,
	skipLog bool,
) (int, int, chainBudget, []string, *gameengine.GameState, bool, bool) {
	runechantCarryover := auraCountByNameInState(masterState, "Runechant")
	ctx := &sequenceContext{
		hero:                masterState.Hero().(hero.Hero),
		pitched:             pitched,
		deck:                d,
		handStart:           held,
		arsenalAtChainStart: arsenalAtChainStart,
		bufs:                bufs,
		runechantCarryover:  runechantCarryover,
		matchup:             mp,
		blockTotal:          blockTotal,
		arsenalInIdx:        arsenalInIdx,
		priorOpponentMarked: masterState.OpponentMarked(),
		priorBanish:         masterState.Banished(),
		priorGraveyard:      masterState.Graveyard(),
		defenders:           defenders,
		skipLog:             skipLog,
		cacheable:           true,
	}
	// Extend bufs.activatedAbilities with item ability instances for this Best call.
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

	// Build a per-leaf state from the master. Defense (if non-modal blockers) mutates
	// this state so post-defense auras / graveyard propagate into the chain start state
	// via leafState.Copy() per permutation. After defense the leaf state's deck is
	// dropped so per-permutation Copy doesn't waste cycles copying it — preparePermState
	// installs a fresh ctx.deck.ShallowCopy() instead (matching old behaviour where defense's
	// deck mutations don't bleed into the chain).
	hasDRs := containsDefenseReaction(defenders)
	hasModalBlocker := containsModalBlocker(defenders)
	leafState := masterState.Copy()
	ctx.leafState = leafState

	var defenseDealtConst int
	defenseCacheableConst := true
	if !hasModalBlocker && len(defenders) > 0 {
		defenseDealtConst, defenseCacheableConst = ctx.runDefense(defenders, pitched, d, mp.IncomingDamage, noBlockBudgetCap, arsenalDefenderIdx)
	}
	leafState.SetDeck(nil)
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
	var bestWinner *gameengine.GameState
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
			dealt, futureValue, winner, legal := ctx.bestSequence(allAttackers)
			if !legal {
				continue
			}
			// Cost the DRs against the prior-turn runechant carryover. Build a one-shot
			// probe engine seeded with the runechant aura so variable-cost DRs read
			// RunechantCount() off this aura.
			drCost := 0
			if len(defenders) > 0 {
				probe := newDRCostProbe(ctx.runechantCarryover)
				for _, def := range defenders {
					if !attackerMetaPtrFor(def).actsAsDR {
						continue
					}
					drCost += def.Cost(probe)
				}
			}
			if drCost > phase.defendBudget {
				continue
			}
			if hasModalBlocker {
				defenseDealt, defenseCacheable = ctx.runDefense(defenders, pitched, d, mp.IncomingDamage, phase.defendBudget-drCost, arsenalDefenderIdx)
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

	_ = attackersMaxCost
	if !foundFeasible {
		return 0, 0, chainBudget{}, nil, nil, false, defenseCacheable
	}
	return bestDealt, defenseDealt, bestBudget, bestSwung, bestWinner, true, ctx.cacheable && defenseCacheable
}

// newDRCostProbe returns a fresh empty *GameEngine seeded only with the runechant aura
// (when runechants > 0) for variable-cost DR cost probing. Defense-reactions read
// RunechantCount() off this engine to decide their Cost; no other state matters.
func newDRCostProbe(runechants int) *gameengine.GameEngine {
	ge := gameengine.New()
	if runechants > 0 {
		ge.CreateAura(token.NewRunechant(runechants))
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
	matchup               Matchup
	blockTotal            int
	arsenalInIdx          int
	priorOpponentMarked   bool
	priorBanish           []card.Card
	priorGraveyard        []card.Card
	activatedAbilities    []card.Card
	activatedAbilityCosts []int
	defenders             []card.Card
	leafState             *gameengine.GameState
	skipLog               bool
	cacheable             bool
	// permState is the active per-permutation state. Set by preparePermState before
	// each permutation's chain run; nil otherwise.
	permState *gameengine.GameState
}

// newPermLogger returns a fresh logger when ctx is recording, or nil for the find-best
// pass. Each permutation gets its own logger so the winning permutation's log doesn't
// get overwritten by subsequent permutations.
func (ctx *sequenceContext) newPermLogger() *turnlogger.TurnLogger {
	if ctx.skipLog {
		return nil
	}
	return turnlogger.New()
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
	state.SetLogger(ctx.newPermLogger())
	state.SetDeck(deckPile)
	state.SetIncomingDamage(matchupIncomingDamage)
	ge := state.Engine()
	cs := &ctx.bufs.drCardStateScratch

	total := 0
	cacheable := true

	// Per-DR view: graveyard = defenders so DRs that scan graveyard see the defender
	// set.
	drGraveyard := append([]card.Card(nil), defenders...)
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
	chainGraveyard := append([]card.Card(nil), ctx.priorGraveyard...)
	chainGraveyard = append(chainGraveyard, defenders...)
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
func (ctx *sequenceContext) preparePermState(playedAttackers []*card.CardState, n int) *gameengine.GameState {
	s := ctx.leafState.CopyPersistentState()
	s.ResetEphemeralState()
	s.SetHero(ctx.hero)
	s.SetArsenal(ctx.arsenalAtChainStart)
	s.SetOpponentMarked(ctx.priorOpponentMarked)
	s.SetBlockTotal(ctx.blockTotal)
	s.SetDeck(ctx.deck.ShallowCopy())
	hand := append([]card.Card(nil), ctx.handStart...)
	for k := 0; k < n; k++ {
		hand = append(hand, playedAttackers[k].Card)
	}
	for _, c := range ctx.attackPitchPerm {
		hand = append(hand, c)
	}
	s.SetPitched(ctx.pitched)
	s.SetHand(hand)
	s.SetLogger(ctx.newPermLogger())
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
		ge := permState.Engine()
		if ge.HasEndOfTurnFire() {
			ge.FireEndOfTurn()
		}
		return 0, pendingFutureValueFromState(permState), permState, true
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
			if winner != nil && ctx.cacheable && !winner.IsCacheable() {
				ctx.cacheable = false
			}
			if !legal {
				return
			}
			drawn := winner.CardsDrawn()
			cmp := chainScoreCmp(dmg, drawn, futureValue, best, bestCardsDrawn, bestFutureValue)
			if !foundLegal || cmp > 0 {
				best = dmg
				bestCardsDrawn = drawn
				bestFutureValue = futureValue
				foundLegal = true
				bestWinner = winner
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
func (ctx *sequenceContext) playSequence(order []card.Card) (damage int, futureValue int, residualBudget int, legal bool) {
	n := len(order)
	pcBuf := ctx.bufs.pcBuf
	meta := ctx.bufs.permMeta[:n]
	for i, c := range order {
		meta[i] = attackerMetaPtrFor(c)
		ctx.seedChainEntry(&pcBuf[i], c, i)
	}
	d, fv, rb, _, lg := ctx.playSequenceWithMeta(n)
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
	ctx.permState = state
	ge := state.Engine()
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
			state.SetCardsPlayed(append(state.CardsPlayed(), pc.Card))
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
		state.SetCardsPlayed(append(state.CardsPlayed(), pc.Card))
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
