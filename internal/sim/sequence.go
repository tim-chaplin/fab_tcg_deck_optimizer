package sim

// Attack-chain search: bestAttackWithWeapons evaluates one partition leaf across all phase /
// weapon masks, bestSequence picks the best ordering of attackers via Heap's algorithm, and
// playSequence* replay a single permutation through TurnState while firing hero triggers,
// Aura handlers, and per-attack OnHit closures.

import (
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// perItemAbilityCap caps how many instances of one item's activated ability the chain
// runner enumerates per turn, bounding the wmask 2^k explosion when an item's Count gets
// large. Realistic Gold counts in play tend to 1-3; 4 leaves headroom without letting a
// pathological hand blow up the per-leaf mask loop.
const perItemAbilityCap = 4

// FormatLogEntry renders a LogEntry into its display string. Chain entries with N=0 drop
// the "(+0)" suffix; trigger entries carry a "(from <source>)" tail. The grouped MyTurn
// renderer prefers formatTextWithDelta for trigger entries that get clustered under their
// parent chain line; FormatLogEntry is the fallback for orphan triggers and external
// callers that just need the verbose string.
func FormatLogEntry(e LogEntry) string {
	if e.Kind == LogEntryChainStep {
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

// Phase masks: when no Defense Reactions are present (or no pitches exist), all pitches go to
// the attack phase, so we visit one configuration. Otherwise we enumerate 2^|pitched| splits.
//
// arsenalAtChainStart is the card sitting in the arsenal slot at the start of the chain — set
// when the partition assigned arsenalCardIn the Arsenal role (it's staying), nil otherwise
// (no arsenal-in, or arsenal-in is playing as Attack/Defend).
func bestAttackWithWeapons(hero Hero, weapons []Weapon, attackers, defenders, pitched, held []card.Card, d *deck.Deck, bufs *attackBufs, mp Matchup, blockTotal, arsenalInIdx, arsenalDefenderIdx int, arsenalAtChainStart card.Card, prior TurnState, skipLog bool) (int, int, chainBudget, []string, CarryState, bool, bool) {
	runechantCarryover := tokenCountIn(prior.auras, TokenTypeRunechant)
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
		priorAuras:          prior.auras,
		priorItems:          prior.Items,
		priorOpponentMarked: prior.opponentMarked,
		priorBanish:         prior.banished,
		priorGraveyard:      prior.graveyard,
		// defenderAuras shares backing with bufs.defenderAurasBacking so the per-partition
		// capture is alloc-free across Best calls.
		defenderAuras: bufs.defenderAurasBacking[:0],
		defenders:     defenders,
		skipLog:       skipLog,
		cacheable:     true,
		// Point carryWinner at the bufs-persistent scratch so SnapshotFromTurn reuses
		// backing arrays across leaves and Best calls (bufs is Evaluator-cached).
		carryWinner: &bufs.carryWinnerScratch,
	}
	// Extend bufs.activatedAbilities with item ability instances for this Best call —
	// the weapon prefix (positions 0..weaponAbilityCount-1) is materialised once at
	// attackBufs construction; items append on top. The unified slice drives both the
	// wmask cost summing and the chain assembly downstream — an activated ability is
	// the same chain step whether it came from a weapon or an item, so they share one
	// list and one path. The only weapon-specific bit is the SwungWeapons name lookup,
	// which keys off bufs.weaponNames[wmask & weaponBitsMask].
	//
	// Per-instance ability costs come from attackerMetaPtrFor's global cardMetaCache
	// so the per-Best assembly stays alloc-free.
	abilities := bufs.activatedAbilities[:bufs.weaponAbilityCount]
	abilityCosts := bufs.activatedAbilityCosts[:bufs.weaponAbilityCount]
	// In-play items contribute min(Count, perItemAbilityCap) ability instances each so
	// the wmask can pick "play it 0..N times this turn".
	for _, it := range prior.Items {
		copies := it.Count
		if copies > perItemAbilityCap {
			copies = perItemAbilityCap
		}
		cost := attackerMetaPtrFor(it.Ability).maxCost
		for i := 0; i < copies; i++ {
			abilities = append(abilities, it.Ability)
			abilityCosts = append(abilityCosts, cost)
		}
	}
	bufs.activatedAbilities = abilities
	bufs.activatedAbilityCosts = abilityCosts
	ctx.activatedAbilities = abilities
	ctx.activatedAbilityCosts = abilityCosts
	// Non-modal defender contribution is constant across phase / weapon masks — DRs through
	// Play, plain blocks as raw block credit — so we compute it once at the top. Modal
	// blockers (Blocker + Modal + BlockCost) need a per-pmask budget to pick their
	// mode, so partitions with any modal blocker recompute defendersDamage per (pmask,
	// wmask) below; the once-per-leaf shortcut applies to the common no-modal-blocker case.
	hasDRs := containsDefenseReaction(defenders)
	hasModalBlocker := containsModalBlocker(defenders)
	var defenseDealtConst int
	// defenseCacheable defaults to true — a partition with no defenders runs no DR Plays,
	// so nothing in the defense phase reads hidden state.
	defenseCacheableConst := true
	if !hasModalBlocker && len(defenders) > 0 {
		defenseDealtConst, defenseCacheableConst = ctx.runDefense(defenders, pitched, d, prior.auras, mp.IncomingDamage, noBlockBudgetCap, arsenalDefenderIdx)
	}
	defenseDealt := defenseDealtConst
	defenseCacheable := defenseCacheableConst

	pitchedVals := bufs.pitchedValsScratch[:0]
	for _, c := range pitched {
		pitchedVals = append(pitchedVals, c.Pitch())
	}

	// Phase splits only matter when there is actually a defense phase to fund (a DR exists) AND
	// there are pitches to split. Otherwise every pitch goes to the attack phase and we visit a
	// single configuration.
	phaseCount := 1
	if hasDRs && len(pitched) > 0 {
		phaseCount = 1 << len(pitched)
	}

	// Pre-screen precomputation: printed-cost sums let us reject doomed (pmask, wmask) pairs in
	// O(1) before spinning up bestSequence's N! permutation loop. attackersMinCost sums the
	// floor-cost of each attacker (non-discount: printed Cost; discount: 0), a safe under-estimate
	// of chain cost. attackersPrinted is the no-discount upper bound, used for the pitch-waste
	// upper bound check.
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

		// Production path uses real pitched cards via attackPitchPerm; resourceBudget is the
		// synthetic-carry escape hatch reserved for tests, so leave it 0 here.
		ctx.resourceBudget = 0
		// Populate attackPitchPerm with the pmask-selected attack-phase pitches in original
		// order, plus a parallel int slice with their Pitch() values cached. bestSequence's
		// nested Heap permutes both slices in lockstep.
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

		// wmask enumerates the subset of activated abilities to play this turn — the
		// low len(weapons) bits select weapon abilities, higher bits select item
		// abilities. Both kinds live in ctx.activatedAbilities, so a single bit-by-bit
		// loop drives both cost summing and chain assembly. weaponBitsMask isolates
		// the weapon prefix for the SwungWeapons name lookup downstream.
		weaponBitsMask := (1 << len(weapons)) - 1
		totalAbilityMasks := 1 << len(ctx.activatedAbilities)
		for wmask := 0; wmask < totalAbilityMasks; wmask++ {
			abilityCost := 0
			for j := range ctx.activatedAbilities {
				if wmask&(1<<j) != 0 {
					abilityCost += ctx.activatedAbilityCosts[j]
				}
			}
			// Lower bound on total chain cost (sum of MinCost across attackers + selected
			// abilities). If the attack budget can't cover even this floor, no permutation is
			// feasible. Mid-turn draws can pitch on top of the committed hand pitch ("hopeful"
			// partitions) but can't reduce the base cost, so this MinCost prune is safe. No
			// matching pitch-timing pre-screen here: drawn cards play as chain extensions and
			// consume the residual, so playSequenceWithMeta enforces pitch-timing post-extension
			// instead.
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
			// Cost the DRs against the prior-turn runechant carryover (cost is paid
			// before the card resolves, so DRs can't fund themselves with auras they
			// create). Variable-cost DRs read s.RunechantCount() inside their Cost; static
			// DRs return a constant. bufs.drScratch is reused per mask iteration — the
			// interface call boxes the pointer, so a stack allocation would escape.
			bufs.drScratch = TurnState{}
			if ctx.runechantCarryover > 0 {
				bufs.drScratchAuras = append(bufs.drScratchAuras[:0], NewRunechantAura(ctx.runechantCarryover))
				bufs.drScratch.auras = bufs.drScratchAuras
			}
			drCost := 0
			for _, d := range defenders {
				if !attackerMetaPtrFor(d).actsAsDR {
					continue
				}
				drCost += d.Cost(&bufs.drScratch)
			}
			if drCost > phase.defendBudget {
				continue
			}
			// Modal blockers (Brothers in Arms, …) draw from the same defense pitch supply
			// the DRs do; recompute defendersDamage per (pmask, wmask) here so each partition
			// candidate sees the right spare budget. Non-modal-blocker partitions stick with
			// the once-per-leaf defenseDealtConst computed above.
			if hasModalBlocker {
				defenseDealt, defenseCacheable = ctx.runDefense(defenders, pitched, d, prior.auras, mp.IncomingDamage, phase.defendBudget-drCost, arsenalDefenderIdx)
			}
			if phase.hasDefendPitches && phase.defendBudget-drCost >= phase.maxDefendPitch {
				continue
			}
			// Tiebreaker order across (pmask, wmask) candidates: chainScoreCmp on
			// (dealt, cardsDrawn, futureValue) → end-of-chain hand size. The hand-size
			// fallback breaks "spend or hoard" ties at zero damage by preferring the chain
			// that drew an extra card into Hand (which post-hoc-promotes into arsenal).
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
				// Reuse bufs.bestCarryScratch's backing arrays so the per-mask-combo
				// update is allocation-free. The mask-combo loop runs to completion
				// against this scratch; findBest's recurse clones it into its own
				// scratch on a new-leaf-best so the alias here is fine.
				bufs.bestCarryScratch.CopyFrom(ctx.carryWinner)
				foundFeasible = true
			}
		}
	}

	if !foundFeasible {
		// No-feasible-line leaves still surface the defense-phase cacheable bit — DR Plays
		// ran independently of the (rejected) attack chain so a DR that read graveyard
		// poisons the result regardless of attack-feasibility.
		return 0, 0, chainBudget{}, nil, CarryState{}, false, defenseCacheable
	}
	// Return bestCarryScratch as an alias — the caller must copy or clone before the
	// next bestAttackWithWeapons call against the same bufs.
	return bestDealt, defenseDealt, bestBudget, bestSwung, bufs.bestCarryScratch, true, ctx.cacheable && defenseCacheable
}

// sequenceContext carries the stable per-partition-leaf environment: hero (for OnCardPlayed
// triggers), pitched / deck refs for Card.Play, shared scratch buffers, and the active pitch
// ordering that funds the attack chain. Built once per leaf so the hot inner calls
// (playSequence, bestSequence) shrink to their varying inputs and tracking outputs.
//
// attackPitchPerm is rewritten by bestAttackWithWeapons on each pmask iteration with the
// attack-phase pitched cards in their original order, then permuted in place by bestSequence's
// pitch Heap loop. playSequenceWithMeta walks it left-to-right, popping cards as costs come up
// and carrying any over-pitch forward; per-card attribution lands in CardState.PitchedToPlay.
// A permutation is rejected if a chain step needs more resources than the remaining pitch
// pool can supply or if any pitch card stays unpopped at end of chain (FaB's pitch-timing
// rule).
//
// resourceBudget is the synthetic starting carry — 0 in the production path (real pitches
// fund every chain step) but set by tests that drive playSequence with a budget number
// instead of a real pitched bag.
type sequenceContext struct {
	hero    Hero
	pitched []card.Card
	// deck is the per-leaf starting deck. resetStateForPermutation gives each
	// permutation a fresh *deck.Deck wrapper via ShallowCopy — the slice headers are
	// per-permutation so Draw advances independently, while the read-only backing
	// array stays shared (no permutation Shuffles, and the mutating PutTop /
	// PutBottom / Tutor methods always allocate fresh backing).
	deck *deck.Deck
	// handStart is the partition's Held-role hand cards — what state.Hand starts as before
	// the chain runs. Cards mutating state.Hand mid-chain (DrawOne, Moon Wish tutor) work
	// against a per-permutation copy so the next permutation gets handStart back.
	handStart []card.Card
	// arsenalAtChainStart is the card sitting in the arsenal slot at the start of the chain
	// — set when the partition assigned arsenalCardIn the Arsenal role, nil otherwise.
	// state.Arsenal starts as this value; cards that destroy or replace arsenal contents
	// during Play would mutate state.Arsenal, but the simulator doesn't model that today.
	arsenalAtChainStart card.Card
	bufs                *attackBufs
	// attackPitchPerm is the active pitch ordering for the attack phase — the pmask-selected
	// subset of ctx.pitched, populated by bestAttackWithWeapons in original order and
	// permuted in place by bestSequence's pitch Heap loop. Backing array is bufs.pitchPermBuf
	// so per-leaf reuse never allocates.
	attackPitchPerm []card.Card
	// attackPitchVals parallels attackPitchPerm: attackPitchVals[i] is the cached Pitch()
	// of attackPitchPerm[i]. Permuted in lockstep with attackPitchPerm so the per-pop
	// resource math reads ints instead of going through the Card.Pitch() interface call.
	attackPitchVals    []int
	resourceBudget     int
	runechantCarryover int
	matchup            Matchup
	blockTotal         int
	// arsenalInIdx is the index in the attackers slice (the slice passed to bestSequence) of
	// the card that came from the arsenal slot at start of turn, or -1 when no arsenal-in card
	// is in the chain. Lets bestSequence flag the matching pcBuf entry's FromArsenal as the
	// permutation moves it around.
	arsenalInIdx int
	// priorAuras are the Auras carried in from the previous turn (e.g. an
	// AttackAction trigger from a Malefic Incantation played a turn ago). Each permutation
	// seeds state.auras with a fresh copy of this slice so mid-chain firing can
	// decrement Count / set FiredThisTurn without leaking those mutations across permutations.
	priorAuras []Aura
	// priorItems are the Items carried in from the previous turn. Each permutation seeds
	// state.Items with a fresh copy so mid-chain ability plays (decrementing Count,
	// destroying items) don't leak across permutations.
	priorItems []Item
	// priorOpponentMarked seeds TurnState.OpponentMarked at each permutation's start —
	// non-zero when the previous turn's chain left a Mark on the opposing hero.
	priorOpponentMarked bool
	// priorBanish is the cards-already-in-banished-zone snapshot from the previous turn.
	// Each permutation seeds state.banished with a fresh copy so "count cards in your
	// banished zone" readers see the carried entries; CardBanished stays reset so
	// "did anything banish THIS turn" stays false until BanishFromGraveyard fires.
	priorBanish []card.Card
	// priorGraveyard is the cards-already-in-graveyard snapshot from the previous turn.
	// Each permutation seeds state.graveyard with a fresh copy AHEAD of this turn's
	// defenders (which the partition's defendersDamage loop has already deposited) so
	// graveyard-scanning riders (BanishFromGraveyard, Sigil of Silphidae's enter
	// trigger) see prior-turn entries alongside same-turn blocks.
	priorGraveyard []card.Card
	// activatedAbilities is the unified weapon + item ability list materialised at the
	// top of bestAttackWithWeapons — weapons' Ability() Cards followed by per-priorItem
	// ability instances (one per Count, capped at perItemAbilityCap). The wmask iterates
	// over this slice; index j's bit selects activatedAbilities[j].
	activatedAbilities []card.Card
	// activatedAbilityCosts parallels activatedAbilities — cached cost for each entry
	// so the per-wmask budget check reads ints rather than re-entering Card.Cost.
	activatedAbilityCosts []int
	// defenderAuras is the post-defense aura set: priorAuras consolidated with the
	// auras DR / plain-block Plays create during defendersDamage. resetStateForPermutation
	// seeds state.auras from this so chain cards see auras created during defense.
	defenderAuras []Aura
	// defenders is the partition's defender slice (DRs + plain blocks together). After
	// the defense phase resolves, every defender lands in the graveyard;
	// resetStateForPermutation seeds the chain's state.graveyard from this so chain cards
	// that scan the graveyard (e.g. recycle riders) see defender cards already there.
	defenders []card.Card
	// carryWinner is a slice header POINTING into bufs.carryWinnerScratch — the persistent
	// snapshot buffer that survives across Best calls via the Evaluator's cached attackBufs.
	// Heap's algorithm keeps iterating past the winner and the shared state.* fields reflect
	// whatever ordering ran last, so the snapshot has to happen the moment a new winner is
	// found; reusing the bufs-owned backing arrays makes that snapshot allocation-free
	// after the first sizing.
	carryWinner *CarryState
	// skipLog drives the per-permutation logger choice. When true, every permutation
	// reset leaves TurnState.logger nil so chains run with Log appends elided (Value
	// still credited); when false, TurnState.logger points at bufs.logger so the
	// per-permutation chain materialises into the bufs-owned scratch. The caller is
	// replaying later with skipLog=false to materialise the printout.
	skipLog bool
	// cacheable is a sticky bit ANDed in after every permutation in bestSequence. Starts
	// true on context construction; flips to false the first time a permutation's chain
	// reports !state.IsCacheable() at end of chain — once a card in any sibling chain reads
	// hidden state, the partition's output isn't safe to cache. Carries across phase /
	// weapon masks within the same leaf because the solver explores all configurations and
	// the cache key would have to disambiguate which the winner came from.
	cacheable bool
}

// activeLogger returns the *turnlogger.TurnLogger the chain runner threads through
// Card.Play / Block / OnHitHandler.Fire / TriggerHandler / Hero.OnCardPlayed for this
// pass — bufs.logger when recording, a typed-nil during the find-best pass so every
// Append* helper short-circuits via TurnLogger's nil-receiver guards (and cards skip
// expensive arg construction via l.Recording()). Centralised here so
// resetStateForPermutation and runDefense pick the same value off ctx.skipLog without
// each spelling the conditional inline.
func (ctx *sequenceContext) activeLogger() *turnlogger.TurnLogger {
	if ctx.skipLog {
		return nil
	}
	return ctx.bufs.logger
}

// runDefense seeds bufs.state.auras with priorAuras (so DR Plays' aura-create helpers
// consolidate against existing tokens), runs defendersDamage, and captures the post-
// defense aura set into ctx.defenderAuras (and bufs.defenderAurasBacking) for later
// chain seeding. Re-bind both headers after the call to track any growth-driven realloc.
func (ctx *sequenceContext) runDefense(defenders, pitched []card.Card, deckPile *deck.Deck, priorAuras []Aura, incomingDamage, blockBudget, arsenalDefenderIdx int) (int, bool) {
	bufs := ctx.bufs
	bufs.state.auras = append(ctx.defenderAuras[:0], priorAuras...)
	// Seed the logger so the first DR's per-defender reset (and the plain-block phase
	// that reuses bufs.state afterwards) sees the chain-level skip / record decision;
	// defendersDamage preserves it through its struct-literal reset. Nil during the
	// find-best pass; bufs.logger during replay.
	bufs.state.logger = ctx.activeLogger()
	dealt, gravScratch, cacheable := defendersDamage(defenders, pitched, deckPile, bufs.state, bufs.defenseGravScratch, &bufs.drCardStateScratch, incomingDamage, blockBudget, arsenalDefenderIdx)
	bufs.defenseGravScratch = gravScratch
	ctx.defenderAuras = bufs.state.auras
	bufs.defenderAurasBacking = bufs.state.auras
	return dealt, cacheable
}

// fireAttackActionAuras walks state.auras after an attack action card resolves
// and invokes every TriggerAttackAction entry whose OncePerTurn gate is open. Each fire
// decrements the trigger's Count; when Count hits zero the aura drops out of the list and
// Self lands in the graveyard so downstream same-turn effects see the destroy. The sim
// publishes the triggering card via state.triggeringCard before each handler runs and
// clears it after; handlers read it through s.AddPreTriggerLogEntry to attribute their
// log line back to the triggering card.
//
// Iterates state.auras with a cursor that handles handler-side splicing: a handler
// calling s.DestroyAura mutates state.auras in place (shifting the next entry down to
// the cursor's index), so the loop only advances when the slice length didn't change.
func fireAttackActionAuras(state *TurnState, triggeringCard card.Card) {
	for i := 0; i < len(state.auras); {
		a := &state.auras[i]
		if a.TriggerType != TriggerAttackAction || (a.OncePerTurn && a.FiredThisTurn) {
			i++
			continue
		}
		state.triggeringCard = triggeringCard
		state.currentAuraIdx = i
		state.currentAuraDestroyed = false
		ctx := auraCtx{a: a, s: state}
		a.Handler(state, state.logger, &ctx)
		state.currentAuraIdx = -1
		state.triggeringCard = nil
		if !state.currentAuraDestroyed {
			state.auras[i].FiredThisTurn = true
			i++
		}
	}
}

// hasEndOfTurnFire reports whether either Auras or Triggers carries a TriggerEndOfTurn
// entry. Lets the chain runner skip the end-of-turn walk when nothing would fire.
func hasEndOfTurnFire(state *TurnState) bool {
	for _, a := range state.auras {
		if a.TriggerType == TriggerEndOfTurn {
			return true
		}
	}
	for _, t := range state.triggers {
		if t.TriggerType == TriggerEndOfTurn {
			return true
		}
	}
	return false
}

// fireEndOfTurn runs after the chain has finished resolving (and the legality gates
// have passed) but before snapshotCarry captures the next-turn state. Walks Auras and
// Triggers in one pass each:
//
//   - Aura entries respect OncePerTurn / FiredThisTurn semantics; the handler owns
//     destruction via s.DestroyAura.
//   - Trigger entries are one-shot; the sim removes each fired entry afterward.
//     Snapshotting len(state.triggers) before iterating keeps a handler that calls
//     AddTrigger from firing its newcomer on the same pass — newcomers stay queued for
//     the next matching event.
func fireEndOfTurn(state *TurnState) {
	for i := 0; i < len(state.auras); {
		a := &state.auras[i]
		if a.TriggerType != TriggerEndOfTurn || (a.OncePerTurn && a.FiredThisTurn) {
			i++
			continue
		}
		state.currentAuraIdx = i
		state.currentAuraDestroyed = false
		ctx := auraCtx{a: a, s: state}
		a.Handler(state, state.logger, &ctx)
		state.currentAuraIdx = -1
		if !state.currentAuraDestroyed {
			state.auras[i].FiredThisTurn = true
			i++
		}
	}
	n := len(state.triggers)
	for i := 0; i < n; i++ {
		tr := state.triggers[i]
		if tr.TriggerType != TriggerEndOfTurn {
			continue
		}
		ctx := triggerCtx{t: &tr}
		tr.Handler(state, state.logger, &ctx)
	}
	kept := state.triggers[:0]
	for i, tr := range state.triggers {
		if i < n && tr.TriggerType == TriggerEndOfTurn {
			continue
		}
		kept = append(kept, tr)
	}
	state.triggers = kept
}

// fireAttackAuras is the TriggerAttack counterpart to fireAttackActionAuras: walks
// state.auras when ANY attack resolves (attack action OR weapon swing) and invokes every
// TriggerAttack entry. The runechant token aura uses this trigger. Same cursor / splice
// semantics as fireAttackActionAuras.
func fireAttackAuras(state *TurnState, triggeringCard card.Card) {
	for i := 0; i < len(state.auras); {
		a := &state.auras[i]
		if a.TriggerType != TriggerAttack || (a.OncePerTurn && a.FiredThisTurn) {
			i++
			continue
		}
		state.triggeringCard = triggeringCard
		state.currentAuraIdx = i
		state.currentAuraDestroyed = false
		ctx := auraCtx{a: a, s: state}
		a.Handler(state, state.logger, &ctx)
		state.currentAuraIdx = -1
		state.triggeringCard = nil
		if !state.currentAuraDestroyed {
			state.auras[i].FiredThisTurn = true
			i++
		}
	}
}

// resetStateForPermutation rewrites every TurnState field to its per-permutation starting
// value. Hand is deep-copied so card-driven mutations (DrawOne, alt-cost prepends) don't
// leak to the next permutation. The leaf-stable read-only fields (Pitched, IncomingDamage,
// BlockTotal) come from ctx; Auras gets a fresh copy of priorAuras so
// mid-chain firing's Count / FiredThisTurn mutations stay scoped. Value resets to 0 so the
// dispatcher can use it as the permutation's running damage total.
//
// The transient slices (Hand, Graveyard, Banish, CardsPlayed, Log, Auras) all
// borrow pre-allocated backing arrays from attackBufs via append([:0], src...) so unchanged
// permutations don't allocate fresh slices. snapshotCarry clones the winning permutation's
// slices before the next permutation overwrites these buffers; mid-chain growth past the
// pre-sized cap is the only path that allocates a new backing array.
//
// deck is a fresh Copy of ctx.deck per permutation so card-driven mutations
// (PrependToDeck, TutorFromDeck, Opt, PopDeckTop) stay scoped to this permutation.
func (ctx *sequenceContext) resetStateForPermutation() {
	s := ctx.bufs.state
	bufs := ctx.bufs
	// Field-by-field assignment skips the implicit zero-fill the struct-literal form does
	// over every TurnState slot before applying the listed fields. The few fields whose
	// permutation-start value is zero (Value, AuraCreated, CardsRemaining, Overpower,
	// NonAttackActionPlayed, ArcaneDamageDealt, Revealed, TriggeringCard) are assigned
	// explicitly so any leftover state from a previous permutation gets cleared.
	s.hand = append(bufs.handBacking[:0], ctx.handStart...)
	// ShallowCopy shares the backing array with ctx.deck; safe because no permutation
	// path Shuffles (only Draw / PeekTop, which re-slice the header on this Deck only)
	// and PutTop / PutBottom / Tutor always allocate fresh backing. Cuts the per-
	// permutation deck.Copy slice allocation that dominated the CPU profile.
	s.deck = ctx.deck.ShallowCopy()
	s.Arsenal = ctx.arsenalAtChainStart
	s.graveyard = append(bufs.graveBacking[:0], ctx.priorGraveyard...)
	s.graveyard = append(s.graveyard, ctx.defenders...)
	s.banished = append(bufs.banishBacking[:0], ctx.priorBanish...)
	s.cardBanished = false
	s.actionPoints = 1
	s.arcaneDamageDealt = false
	s.opponentMarked = ctx.priorOpponentMarked
	// Seed from defenderAuras when defendersDamage ran for this leaf (it already
	// includes priorAuras consolidated with DR-created auras); fall back to priorAuras
	// otherwise.
	auraSeed := ctx.priorAuras
	if len(ctx.defenderAuras) > 0 {
		auraSeed = ctx.defenderAuras
	}
	s.auras = append(bufs.auraTriggersBacking[:0], auraSeed...)
	// Triggers are one-shot and don't carry across permutations — reset to empty so a
	// prior permutation's Triggers don't leak.
	s.triggers = bufs.triggersBacking[:0]
	s.Items = append(bufs.itemsBacking[:0], ctx.priorItems...)
	s.CardsDrawn = 0
	s.currentAuraIdx = -1
	s.value = 0
	// Bind the logger to the per-permutation backing buffer when recording, leaving
	// s.logger nil during the find-best pass so every Log helper short-circuits at
	// the call site. ctx.activeLogger encapsulates the skip / record decision.
	if l := ctx.activeLogger(); l != nil {
		l.SetBuffer(bufs.logBacking)
		s.logger = l
	} else {
		s.logger = nil
	}
	s.cardsPlayed = bufs.cardsPlayedBacking[:0]
	s.auraCreated = false
	s.cardsRemaining = nil
	s.pitched = ctx.pitched
	s.nonAttackActionPlayed = false
	s.incomingDamage = ctx.matchup.IncomingDamage
	s.arcaneIncomingDamage = ctx.matchup.ArcaneIncomingDamage
	s.blockTotal = ctx.blockTotal
	s.attackReactionTarget = nil
	s.triggeringCard = nil
	// Permutation seed starts cacheable; the first card-driven deck / graveyard read
	// in this permutation flips it to false. Set explicitly because zero-value is false.
	s.cacheable = true
}

// bestSequence tries every ordering of attackers and returns the max total damage plus the
// pendingFutureValue at the end of the winning permutation. Between each card's Play() and
// its append to CardsPlayed, the hero's OnCardPlayed hook fires so triggered abilities
// contribute. legal=true when at least one ordering is playable; false when every permutation
// is rejected by playSequenceWithMeta's resource / go-again / pitch-waste checks.
//
// Uses Heap's algorithm (iterative) — no closure/callback alloc, no recursive call per perm.
// The winning permutation's end-of-chain CarryState lands in ctx.carryWinner so callers can
// adopt the snapshot for next-turn state.
func (ctx *sequenceContext) bestSequence(attackers []card.Card) (int, int, bool) {
	n := len(attackers)
	if n == 0 {
		// No chain steps means no costs to pay. Any unspent pitch card in the attack phase
		// breaks FaB's pitch-timing rule — pitching is only legal to fund a cost on the stack
		// — so a non-empty attackPitchPerm rejects the empty chain.
		if len(ctx.attackPitchPerm) > 0 {
			return 0, 0, false
		}
		// Empty-chain leaves still need a populated CarryState — the cache-replay path
		// adopts ctx.carryWinner directly and the snapshot must reflect the held cards in
		// state.Hand so post-hoc arsenal promotion has something to pick from. Reset+snapshot
		// mirrors the per-permutation work eval() does for n>0 chains.
		ctx.resetStateForPermutation()
		st := ctx.bufs.state
		if hasEndOfTurnFire(st) {
			fireEndOfTurn(st)
		}
		ctx.carryWinner.SnapshotFromTurn(st)
		return 0, pendingFutureValue(st.auras, st.Items), true
	}
	pcBuf := ctx.bufs.pcBuf[:n]
	permMeta := ctx.bufs.permMeta[:n]
	for idx, c := range attackers {
		permMeta[idx] = attackerMetaPtrFor(c)
		// Field-by-field assignment preserves pcBuf[idx].OnHit's backing array across
		// Best calls — a struct-literal assignment would drop it and force every Play
		// that registers an OnHit to allocate a fresh slice on the hot anneal path.
		// Other sliced fields (PitchedToPlay) get reset by playSequenceWithMeta.
		pcBuf[idx].Card = c
		pcBuf[idx].FromArsenal = idx == ctx.arsenalInIdx
		pcBuf[idx].GrantedGoAgain = false
		pcBuf[idx].GrantedDominate = false
		pcBuf[idx].GrantedOverpower = false
		pcBuf[idx].BonusAttack = 0
		pcBuf[idx].BonusDefense = 0
		pcBuf[idx].PitchedToPlay = nil
		pcBuf[idx].OnHit = pcBuf[idx].OnHit[:0]
		pcBuf[idx].Mode = 0
	}

	best := 0
	bestFutureValue := 0
	foundLegal := false
	// Zero ctx.carryWinner's contents (preserving slice backing arrays) so a stale value
	// from a previous Best call's leaf can't leak through when no permutation lands a new
	// best in this leaf. The slice lengths drop to 0 but backing arrays survive — the
	// next SnapshotFromTurn refills via append([:0], src...) without allocating.
	ctx.carryWinner.Reset()
	state := ctx.bufs.state
	pitchPerm := ctx.attackPitchPerm
	pitchVals := ctx.attackPitchVals
	pn := len(pitchPerm)
	// tupleCount is the cartesian product of permMeta[i].modes across the chain. The product
	// is permutation-invariant (multiplication commutes), so it's computed once here and
	// reused across every (attack, pitch) ordering. hasModal short-circuits the per-tuple
	// mode decode for chains with no ModalCards — pcBuf[i].Mode is already 0 from the per-
	// permutation reset, so the fast path skips the mixed-radix decode entirely.
	tupleCount := 1
	for i := 0; i < n; i++ {
		tupleCount *= int(permMeta[i].modes)
	}
	hasModal := tupleCount > 1
	bestCardsDrawn := 0
	// tryPitchOrdering plays the chain against the current attack-permutation × pitch-
	// permutation pair, threads cacheable, and folds a legal result into the running
	// best via chainScoreCmp on (dmg, cardsDrawn, futureValue). cardsDrawn varies across
	// permutations because PlayPrecondition can reject in some orderings (Gold ability
	// fires only when the order resolves a token-creator first); ranking by it makes
	// the chain runner prefer a perm that drew over one that didn't at equal damage.
	// Modal chains additionally enumerate the cartesian product of Modal mode
	// indices via a mixed-radix decode of `tuple` over permMeta[i].modes. The per-perm
	// resolve+fold body is inlined into both branches: hoisting into a closure spills
	// the capture set to the heap on every call inside Heap's permutation loop.
	tryPitchOrdering := func() {
		if !hasModal {
			dmg, futureValue, _, legal := ctx.playSequenceWithMeta(n)
			if ctx.cacheable && !state.IsCacheable() {
				ctx.cacheable = false
			}
			if !legal {
				return
			}
			cmp := chainScoreCmp(dmg, state.CardsDrawn, futureValue, best, bestCardsDrawn, bestFutureValue)
			if !foundLegal || cmp > 0 {
				best = dmg
				bestCardsDrawn = state.CardsDrawn
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
			cmp := chainScoreCmp(dmg, state.CardsDrawn, futureValue, best, bestCardsDrawn, bestFutureValue)
			if !foundLegal || cmp > 0 {
				best = dmg
				bestCardsDrawn = state.CardsDrawn
				bestFutureValue = futureValue
				foundLegal = true
				ctx.carryWinner.SnapshotFromTurn(state)
			}
		}
	}
	// eval runs the active attack permutation against every pitch ordering — initial
	// ordering plus Heap's enumeration over attackPitchPerm. pn ∈ {0,1} naturally collapse
	// to the single initial call (the inner loop's bound rejects). The pitch Heap swaps
	// pitchPerm and pitchVals in lockstep so the cached Pitch() values stay aligned.
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
	// Heap's algorithm, iterative: c[] counts how many times each stack frame has iterated.
	// pcBuf and permMeta swap together so playSequenceWithMeta sees meta aligned with the
	// current permutation. FromArsenal rides inside pcBuf (one byte), so it permutes for free;
	// no separate permFromArsenal slice to maintain.
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

// playSequence plays `order` as a sequence of cards, reusing ctx.bufs' pooled buffers.
// Buffers are mutated in place; the caller must not read them concurrently.
//
// Runechant flow:
//   - state.RunechantCount() starts at ctx.runechantCarryover.
//   - CreateRunechants bumps the count and credits +n damage at creation time.
//   - Each Attack / Weapon resolution fires all current tokens and destroys them; no
//     re-credit (tokens were credited at creation).
//   - Surviving runechants at end of chain feed into pendingFutureValue along with
//     other auras and items.
//
// Resource flow lives on playSequenceWithMeta; this wrapper just forwards.
//
// Populates permMeta from order and then calls playSequenceWithMeta. The hot path
// (bestSequence) builds meta once and calls playSequenceWithMeta directly to amortise
// interface dispatch across the N! permutations.
func (ctx *sequenceContext) playSequence(order []card.Card) (damage int, futureValue int, residualBudget int, legal bool) {
	n := len(order)
	pcBuf := ctx.bufs.pcBuf
	meta := ctx.bufs.permMeta[:n]
	for i, c := range order {
		meta[i] = attackerMetaPtrFor(c)
		// Field-by-field — preserve pcBuf[i].OnHit backing across calls (see bestSequence).
		pcBuf[i].Card = c
		pcBuf[i].FromArsenal = i == ctx.arsenalInIdx
		pcBuf[i].GrantedGoAgain = false
		pcBuf[i].GrantedDominate = false
		pcBuf[i].GrantedOverpower = false
		pcBuf[i].BonusAttack = 0
		pcBuf[i].BonusDefense = 0
		pcBuf[i].PitchedToPlay = nil
		pcBuf[i].OnHit = pcBuf[i].OnHit[:0]
		pcBuf[i].Mode = 0
	}
	return ctx.playSequenceWithMeta(n)
}

// playSequenceWithMeta runs the permutation currently held in ctx.bufs.pcBuf[:n] with
// aligned permMeta[:n]. CardState (Card + FromArsenal) persists across permutations, so any
// field a prior card's Play flips on a future card needs a per-permutation reset:
// GrantedGoAgain (next-attack go-again grants), BonusAttack (next-attack +N{p} grants),
// and PitchedToPlay (per-card pitch attribution recomputed against the active pitch
// ordering).
//
// Resource flow: each chain step's cost is paid by pitchPool.pay against the pool seeded
// from ctx.attackPitchPerm × ctx.attackPitchVals (and ctx.resourceBudget for the test-
// only synthetic-budget path). The pool returns the pitched cards that funded this step,
// which land directly on pc.PitchedToPlay. End-of-chain validity: pool.idx must equal
// pool.n — a pitched card held back without funding any cost rejects the permutation.
//
// Damage flows through state.value: the dispatcher records the chain step's
// Play+BonusAttack contribution via state.AddLogEntry; pre-trigger handlers (hero, aura)
// credit themselves through AddPreTriggerLogEntry, post-trigger handlers (OnHit, AR
// buffs) through AddPostTriggerLogEntry. The returned damage is just state.value at end
// of chain.
func (ctx *sequenceContext) playSequenceWithMeta(n int) (damage int, futureValue int, residualBudget int, legal bool) {
	pcBuf := ctx.bufs.pcBuf
	ptrBuf := ctx.bufs.ptrBuf
	meta := ctx.bufs.permMeta[:n]
	for i := 0; i < n; i++ {
		pcBuf[i].GrantedGoAgain = false
		pcBuf[i].BonusAttack = 0
		pcBuf[i].PitchedToPlay = nil
		pcBuf[i].OnHit = pcBuf[i].OnHit[:0]
	}
	played := ptrBuf[:n]
	// Per-permutation reset: full-state rewrite. Hand and Deck are deep-copied so cards can
	// mutate them freely without leaking to the next permutation. state.value resets to 0.
	ctx.resetStateForPermutation()
	state := ctx.bufs.state
	// Seed state.hand with the upcoming chain attackers so each chain step's Play sees an
	// accurate "in hand right now" snapshot — committed cards (pitched, defending, already
	// played, the playing card) are out, but cards going to be played later in this chain
	// stay in. The current card gets removed at the top of each iteration. handStart (Held
	// cards) is already in state.hand from resetStateForPermutation; mid-chain DrawOne
	// continues to append; chain attackers join here.
	for k := 0; k < n; k++ {
		state.hand = append(state.hand, played[k].Card)
	}
	// Pitched cards stay in hand until the pool actually pops them to fund a cost — a card
	// in the partition's Pitch role isn't "pitched" yet, just queued. Mid-chain cards
	// reading state.hand should see the as-yet-unconsumed pitches alongside upcoming chain
	// steps and the Held cards.
	for _, c := range ctx.attackPitchPerm {
		state.hand = append(state.hand, c)
	}
	pool := pitchPool{
		perm:      ctx.attackPitchPerm,
		vals:      ctx.attackPitchVals,
		n:         len(ctx.attackPitchPerm),
		remaining: ctx.resourceBudget,
		attr:      ctx.bufs.pitchAttrBuf[:0],
	}
	// activeAttack is the most recent attack/weapon CardState awaiting OnHit firing. ARs
	// played later buff it; finalizeActiveAttack flushes it (fires OnHit when LikelyToHit
	// is true on the post-buff EffectiveAttack) on the next non-AR card or at end of chain.
	var activeAttack *card.CardState
	finalizeActiveAttack := func() {
		if activeAttack == nil {
			return
		}
		if LikelyToHit(activeAttack) {
			for i := range activeAttack.OnHit {
				h := &activeAttack.OnHit[i]
				h.Fire(state, state.logger, activeAttack, h)
			}
			// Drain matching TriggerHit Triggers — each rider's TypeFilter narrows the
			// qualifying hits (e.g. "attack action card" vs broader "attack" wording).
			// Triggers that don't match stay queued for a later qualifying hit. The
			// firing handler reads state.triggeringCard for attack identity.
			if len(state.triggers) > 0 {
				types := activeAttack.Card.Types(nil)
				prevTriggering := state.triggeringCard
				state.triggeringCard = activeAttack.Card
				kept := state.triggers[:0]
				for i := range state.triggers {
					t := state.triggers[i]
					if t.TriggerType != TriggerHit ||
						(t.TypeFilter != nil && !t.TypeFilter(types)) {
						kept = append(kept, t)
						continue
					}
					ctx := triggerCtx{t: &t}
					t.Handler(state, state.logger, &ctx)
				}
				state.triggers = kept
				state.triggeringCard = prevTriggering
			}
		}
		activeAttack = nil
	}
	for i, pc := range played {
		m := meta[i]
		// Action Point gate: paying chain steps cost 1 AP; free steps (Instants, Attack
		// Reactions) cost 0. A paying card resolving with no AP available rejects the
		// permutation. Go again and AP-grant effects restock the pool for later steps.
		if !m.isFreeChainStep {
			if state.actionPoints <= 0 {
				return 0, 0, 0, false
			}
			state.actionPoints--
		}
		// Remove the playing card from state.hand before resolving — it's leaving the hand
		// to enter the chain. Linear search by interface equality works because every card
		// implementation is a zero-sized struct, so two copies compare equal and any one of
		// them is fine to drop.
		for j := range state.hand {
			if state.hand[j] == pc.Card {
				state.hand = append(state.hand[:j], state.hand[j+1:]...)
				break
			}
		}
		prevPitchIdx := pool.idx
		contrib, ok := pool.pay(m.costAt(state, pc.Mode))
		if !ok {
			return 0, 0, 0, false
		}
		pc.PitchedToPlay = contrib
		// Drop pitches the pool freshly popped on this card's behalf. Same interface-equality
		// removal as the playing-card drop above; Carry-from-prior-step pitches were already
		// removed when their own pay call popped them.
		for k := prevPitchIdx; k < pool.idx; k++ {
			popped := pool.perm[k]
			for j := range state.hand {
				if state.hand[j] == popped {
					state.hand = append(state.hand[:j], state.hand[j+1:]...)
					break
				}
			}
		}
		// PlayPrecondition for ARs: ARs don't trigger finalizeActiveAttack (the AR needs
		// the active attack alive), so any AR precondition runs against the pre-OnHit
		// state — the AR's own predicate (ARTargetAllowed) handles target-shape gates.
		// Non-AR preconditions wait until after finalizeActiveAttack fires below so the
		// check reads OnHit-mutated state (e.g. an Item created by the previous attack's
		// hit).
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
			ctx.hero.OnCardPlayed(pc.Card, state, state.logger)
			state.attackReactionTarget = activeAttack
			ResolveChainStep(state, state.logger, pc)
			state.attackReactionTarget = nil
			state.cardsPlayed = append(state.cardsPlayed, pc.Card)
			state.graveyard = append(state.graveyard, pc.Card)
			// Go again is not printed on ARs but honour the flag if granted.
			if pc.EffectiveGoAgain(state) {
				state.actionPoints++
			}
			continue
		}

		// Non-AR card: flush any pending OnHit (uses the previous attack's post-buff
		// EffectiveAttack) before the new card's precondition runs — so a precondition
		// gating on Items / Auras created by the previous attack's OnHit (Gold ability
		// after Strike Gold) reads the post-OnHit state.
		finalizeActiveAttack()
		if pre, ok := pc.Card.(card.PlayPrecondition); ok {
			if !pre.PlayPrecondition(state, pc) {
				return 0, 0, 0, false
			}
		}

		state.cardsRemaining = played[i+1:]

		// Hero ability fires BEFORE the card's own Play so "aura created this turn" checks
		// inside the card's Play see the runechant (or other aura) the hero just made.
		// Viserai's "another non-attack action" gate still excludes the current card because
		// NonAttackActionPlayed isn't flipped until the end of the iteration. The hero
		// handler logs its own contribution via state.AddPreTriggerLogEntry; its int return
		// is unused.
		ctx.hero.OnCardPlayed(pc.Card, state, state.logger)
		// Card.Play owns its chain-step log line and pre-buff damage contribution. ARs top
		// up the parent chain step's delta; OnHit funcs fire later via finalizeActiveAttack.
		// fireAttackAuras runs after Play so continuous "while you control an aura"
		// modifiers (Yinti Yanti +1{p}) see live token auras before TriggerAttack
		// handlers consume them.
		state.currentStepRerouted = false
		ResolveChainStep(state, state.logger, pc)
		if m.isAttack {
			fireAttackAuras(state, pc.Card)
			// Clear after Play so any "if defending hero is marked" rider on this card
			// already read the pre-clear flag. OnHit handlers fire later via
			// finalizeActiveAttack, after this clear, so on-hit Mark riders survive.
			state.opponentMarked = false
		}
		if m.isAttackAction {
			fireAttackActionAuras(state, pc.Card)
		}
		if m.isAttack {
			activeAttack = pc
		}
		state.cardsPlayed = append(state.cardsPlayed, pc.Card)
		if m.types.IsNonAttackAction() {
			state.nonAttackActionPlayed = true
		}
		// Weapons and persistent card types (Auras, Items) stay in their zone when they
		// resolve; any destroy event that should send them to the graveyard is a separate
		// trigger. Everything else — Actions, Attack Reactions, Defense Reactions, Blocks,
		// Instants — heads to the graveyard immediately, unless the card's Play routed it
		// elsewhere via RecycleToDeckBottom (which sets state.currentStepRerouted).
		// Direct field write — the framework driving the chain isn't a card-driven content
		// read, so no cacheable flip.
		if !m.types.PersistsInPlay() && !state.currentStepRerouted {
			state.graveyard = append(state.graveyard, pc.Card)
		}

		// Go again grants 1 AP after the card resolves. EffectiveGoAgain folds in both
		// printed Go again and mid-chain conditional grants — by the time we get here the
		// card's Play has had a chance to flip GrantedGoAgain on itself.
		if pc.EffectiveGoAgain(state) {
			state.actionPoints++
		}
	}
	// Flush any attack still pending OnHit at end of chain.
	finalizeActiveAttack()

	// Pitch-timing rule: every Pitch-role card must have paid for something on the stack. If
	// the chain finished with pitches still queued, one of them was held back without funding
	// any cost — illegal in FaB. Leftover carry on the front (frontRemaining > 0) is fine —
	// that's the over-pitch surplus on the last popped card, not a held-back pitch.
	if pool.idx < pool.n {
		return 0, 0, 0, false
	}
	// End-of-turn fire after the chain settles, before snapshot: aura entries
	// (Ponders, chain-created end-of-turn auras) and one-shot Triggers.
	if hasEndOfTurnFire(state) {
		fireEndOfTurn(state)
	}
	return state.value, pendingFutureValue(state.auras, state.Items), pool.remaining, true
}
