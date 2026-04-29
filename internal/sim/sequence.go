package sim

// Attack-chain search: bestAttackWithWeapons evaluates one partition leaf across all phase /
// weapon masks, bestSequence picks the best ordering of attackers via Heap's algorithm, and
// playSequence* replay a single permutation through TurnState while firing hero triggers and
// AuraTrigger / EphemeralAttackTrigger handlers.

import (
	"fmt"
)

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
func bestAttackWithWeapons(hero Hero, weapons []Weapon, attackers, defenders, pitched, held, deck []Card, bufs *attackBufs, runechantCarryover, incomingDamage, blockTotal, arsenalInIdx, arsenalDefenderIdx int, arsenalAtChainStart Card, priorAuraTriggers []AuraTrigger, skipLog bool) (int, int, int, chainBudget, []string, CarryState, bool, bool) {
	ctx := &sequenceContext{
		hero:                hero,
		pitched:             pitched,
		deck:                deck,
		handStart:           held,
		arsenalAtChainStart: arsenalAtChainStart,
		bufs:                bufs,
		runechantCarryover:  runechantCarryover,
		incomingDamage:      incomingDamage,
		blockTotal:          blockTotal,
		arsenalInIdx:        arsenalInIdx,
		priorAuraTriggers:   priorAuraTriggers,
		skipLog:             skipLog,
		cacheable:           true,
		// Point carryWinner at the bufs-persistent scratch so SnapshotFromTurn reuses
		// backing arrays across leaves and Best calls (bufs is Evaluator-cached).
		carryWinner: &bufs.carryWinnerScratch,
	}
	// Defenders fire independently of ordering and attack chain — DRs through Play, plain
	// blocks as raw block credit — so their total Value contribution is constant across phase
	// / weapon masks. Compute it once. Includes DR blocks + arcane / runechant riders + plain-
	// block residual against the partition's incoming damage; over-blocked excess is discarded
	// by the per-card cap.
	hasDRs := containsDefenseReaction(defenders)
	var defenseDealt int
	// defenseCacheable defaults to true — a partition with no defenders runs no DR Plays,
	// so nothing in the defense phase reads hidden state.
	defenseCacheable := true
	if len(defenders) > 0 {
		defenseDealt, bufs.defenseGravScratch, defenseCacheable = defendersDamage(defenders, pitched, deck, bufs.state, bufs.defenseGravScratch, &bufs.drCardStateScratch, incomingDamage, arsenalDefenderIdx)
	}

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
	bestLeftoverRunechants := runechantCarryover
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

		for wmask := 0; wmask < 1<<len(weapons); wmask++ {
			weaponCost := bufs.weaponCosts[wmask] // weapons are static-cost
			// Lower bound on total chain cost (sum of MinCost across attackers + weapons). If the
			// attack budget can't cover even this floor, no permutation is feasible. Mid-turn
			// draws can pitch on top of the committed hand pitch ("hopeful" partitions) but
			// can't reduce the base cost, so this MinCost prune is safe. No matching pitch-timing
			// pre-screen here: drawn cards play as chain extensions and consume the residual, so
			// playSequenceWithMeta enforces pitch-timing post-extension instead.
			if attackersMinCost+weaponCost > phase.attackBudget {
				continue
			}
			allAttackers := bufs.attackerBuf[:len(attackers)]
			for i, w := range weapons {
				if wmask&(1<<i) != 0 {
					allAttackers = append(allAttackers, w)
				}
			}
			dealt, leftoverRunechants, legal := ctx.bestSequence(allAttackers)
			if !legal {
				continue
			}
			// Cost the DRs against the chain's final runechant count. DRs with variable cost
			// read state.Runechants inside their Cost; static DRs return a constant. Reuse
			// bufs.drScratch instead of allocating a fresh TurnState per mask iteration — the
			// interface call boxes the pointer, so a stack allocation would escape and heap-alloc
			// every loop.
			bufs.drScratch = TurnState{Runechants: leftoverRunechants}
			drCost := 0
			for _, d := range defenders {
				if !d.Types().IsDefenseReaction() {
					continue
				}
				drCost += d.Cost(&bufs.drScratch)
			}
			if drCost > phase.defendBudget {
				continue
			}
			if phase.hasDefendPitches && phase.defendBudget-drCost >= phase.maxDefendPitch {
				continue
			}
			if !foundFeasible || dealt > bestDealt ||
				(dealt == bestDealt && leftoverRunechants > bestLeftoverRunechants) {
				bestDealt = dealt
				bestLeftoverRunechants = leftoverRunechants
				bestSwung = bufs.weaponNames[wmask]
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
		return 0, 0, 0, chainBudget{}, nil, CarryState{}, false, defenseCacheable
	}
	// Return bestCarryScratch as an alias — the caller (findBest's recurse, replayBest)
	// must copy or clone before the next bestAttackWithWeapons call against the same bufs.
	// findBest's recurse calls bufs.findBestCarryScratch.CopyFrom(carry) on a new-best
	// leaf and clones once at end of findBest; the replayBest path consumes the alias
	// before any second call to bestAttackWithWeapons.
	return bestDealt, defenseDealt, bestLeftoverRunechants, bestBudget, bestSwung, bufs.bestCarryScratch, true, ctx.cacheable && defenseCacheable
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
	hero          Hero
	pitched, deck []Card
	// handStart is the partition's Held-role hand cards — what state.Hand starts as before
	// the chain runs. Cards mutating state.Hand mid-chain (DrawOne, Moon Wish tutor) work
	// against a per-permutation copy so the next permutation gets handStart back.
	handStart []Card
	// arsenalAtChainStart is the card sitting in the arsenal slot at the start of the chain
	// — set when the partition assigned arsenalCardIn the Arsenal role, nil otherwise.
	// state.Arsenal starts as this value; cards that destroy or replace arsenal contents
	// during Play would mutate state.Arsenal, but the simulator doesn't model that today.
	arsenalAtChainStart Card
	bufs                *attackBufs
	// attackPitchPerm is the active pitch ordering for the attack phase — the pmask-selected
	// subset of ctx.pitched, populated by bestAttackWithWeapons in original order and
	// permuted in place by bestSequence's pitch Heap loop. Backing array is bufs.pitchPermBuf
	// so per-leaf reuse never allocates.
	attackPitchPerm []Card
	// attackPitchVals parallels attackPitchPerm: attackPitchVals[i] is the cached Pitch()
	// of attackPitchPerm[i]. Permuted in lockstep with attackPitchPerm so the per-pop
	// resource math reads ints instead of going through the Card.Pitch() interface call.
	attackPitchVals    []int
	resourceBudget     int
	runechantCarryover int
	incomingDamage     int
	blockTotal         int
	// arsenalInIdx is the index in the attackers slice (the slice passed to bestSequence) of
	// the card that came from the arsenal slot at start of turn, or -1 when no arsenal-in card
	// is in the chain. Lets bestSequence flag the matching pcBuf entry's FromArsenal as the
	// permutation moves it around.
	arsenalInIdx int
	// priorAuraTriggers are the AuraTriggers carried in from the previous turn (e.g. an
	// AttackAction trigger from a Malefic Incantation played a turn ago). Each permutation
	// seeds state.AuraTriggers with a fresh copy of this slice so mid-chain firing can
	// decrement Count / set FiredThisTurn without leaking those mutations across permutations.
	priorAuraTriggers []AuraTrigger
	// carryWinner is a slice header POINTING into bufs.carryWinnerScratch — the persistent
	// snapshot buffer that survives across Best calls via the Evaluator's cached attackBufs.
	// Heap's algorithm keeps iterating past the winner and the shared state.* fields reflect
	// whatever ordering ran last, so the snapshot has to happen the moment a new winner is
	// found; reusing the bufs-owned backing arrays makes that snapshot allocation-free
	// after the first sizing.
	carryWinner *CarryState
	// skipLog propagates into TurnState.SkipLog on every permutation reset. When true,
	// chains run with Log appends elided (Value still credited); the caller is replaying
	// later with skipLog=false to materialise the printout.
	skipLog bool
	// cacheable is a sticky bit ANDed in after every permutation in bestSequence. Starts
	// true on context construction; flips to false the first time a permutation's chain
	// reports !state.IsCacheable() at end of chain — once a card in any sibling chain reads
	// hidden state, the partition's output isn't safe to cache. Carries across phase /
	// weapon masks within the same leaf because the solver explores all configurations and
	// the cache key would have to disambiguate which the winner came from.
	cacheable bool
}

// fireAttackActionTriggers walks state.AuraTriggers after an attack action card resolves
// and invokes every TriggerAttackAction entry whose OncePerTurn gate is open. Each fire
// decrements the trigger's Count; when Count hits zero the aura drops out of the list and
// Self lands in the graveyard so downstream same-turn effects see the destroy. The sim
// publishes the triggering card via state.TriggeringCard before each handler runs and
// clears it after; handlers read it through s.AddPreTriggerLogEntry to attribute their
// log line back to the triggering card.
//
// Slice mutation: a survivors prefix is built in place over the existing slice; entries
// kept after firing are written back at increasing indices, exhausted ones are skipped.
func fireAttackActionTriggers(state *TurnState, triggeringCard Card) {
	triggers := state.AuraTriggers
	dst := triggers[:0]
	for i := range triggers {
		t := triggers[i]
		if t.Type != TriggerAttackAction || (t.OncePerTurn && t.FiredThisTurn) {
			dst = append(dst, t)
			continue
		}
		state.TriggeringCard = triggeringCard
		t.Handler(state)
		state.TriggeringCard = nil
		t.FiredThisTurn = true
		t.Count--
		if t.Count <= 0 {
			// Direct field write — the framework destroying an exhausted aura is
			// deterministic from cards played, not a card-driven content read, so no
			// cacheable flip.
			state.graveyard = append(state.graveyard, t.Self)
			continue
		}
		dst = append(dst, t)
	}
	state.AuraTriggers = dst
}

// fireEphemeralAttackTriggers walks state.EphemeralAttackTriggers after an attack action
// card resolves and invokes every entry whose Matches predicate accepts the attacker. Each
// fire consumes the trigger (fire-once semantics). Handlers receive target as a direct
// arg and call s.AddPostTriggerLogEntry themselves to log their damage-equivalent. Non-matching
// entries stay in the slice for a later attack action; anything still in the list at end
// of chain fizzles silently (no graveyard bookkeeping — the source was already graveyarded
// when its own Play resolved).
//
// Slice mutation parallels fireAttackActionTriggers: a survivors prefix is built in place
// over the existing slice, with fired entries skipped.
func fireEphemeralAttackTriggers(state *TurnState, target *CardState) {
	triggers := state.EphemeralAttackTriggers
	dst := triggers[:0]
	for i := range triggers {
		t := triggers[i]
		if t.Matches != nil && !t.Matches(target) {
			dst = append(dst, t)
			continue
		}
		t.Handler(state, target)
	}
	state.EphemeralAttackTriggers = dst
}

// resetStateForPermutation rewrites every TurnState field to its per-permutation starting
// value. Hand and Deck are deep-copied so card-driven mutations (DrawOne, tutors, alt-cost
// prepends) don't leak to the next permutation. The leaf-stable read-only fields (Pitched,
// IncomingDamage, BlockTotal) come from ctx; AuraTriggers gets a fresh copy of
// priorAuraTriggers so mid-chain firing's Count / FiredThisTurn mutations stay scoped.
// Value resets to 0 so the dispatcher can use it as the permutation's running damage total.
//
// The transient slices (Hand, Deck, Graveyard, Banish, CardsPlayed, Log, AuraTriggers) all
// borrow pre-allocated backing arrays from attackBufs via append([:0], src...) so unchanged
// permutations don't allocate fresh slices. snapshotCarry clones the winning permutation's
// slices before the next permutation overwrites these buffers; mid-chain growth past the
// pre-sized cap is the only path that allocates a new backing array.
func (ctx *sequenceContext) resetStateForPermutation() {
	s := ctx.bufs.state
	bufs := ctx.bufs
	*s = TurnState{
		Hand:                    append(bufs.handBacking[:0], ctx.handStart...),
		deck:                    append(bufs.deckBacking[:0], ctx.deck...),
		Arsenal:                 ctx.arsenalAtChainStart,
		graveyard:               bufs.graveBacking[:0],
		Banish:                  bufs.banishBacking[:0],
		CardsPlayed:             bufs.cardsPlayedBacking[:0],
		Log:                     bufs.logBacking[:0],
		Pitched:                 ctx.pitched,
		IncomingDamage:          ctx.incomingDamage,
		BlockTotal:              ctx.blockTotal,
		Runechants:              ctx.runechantCarryover,
		AuraTriggers:            append(bufs.auraTriggersBacking[:0], ctx.priorAuraTriggers...),
		EphemeralAttackTriggers: bufs.ephemeralBacking[:0],
		SkipLog:                 ctx.skipLog,
		// Permutation seed starts cacheable; the first card-driven deck / graveyard read
		// in this permutation flips it to false. Set explicitly because zero-value is false.
		cacheable: true,
	}
}

// bestSequence tries every ordering of attackers and returns the max total damage plus the
// runechant count at the end of the winning permutation. Between each card's Play() and its
// append to CardsPlayed, the hero's OnCardPlayed hook fires so triggered abilities contribute.
// legal=true when at least one ordering is playable; false when every permutation is rejected
// by playSequenceWithMeta's resource / go-again / pitch-waste checks.
//
// Uses Heap's algorithm (iterative) — no closure/callback alloc, no recursive call per perm.
// The winning permutation's end-of-chain CarryState lands in ctx.carryWinner so callers can
// adopt the snapshot for next-turn state.
func (ctx *sequenceContext) bestSequence(attackers []Card) (int, int, bool) {
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
		ctx.carryWinner.SnapshotFromTurn(ctx.bufs.state)
		return 0, ctx.runechantCarryover, true
	}
	pcBuf := ctx.bufs.pcBuf[:n]
	permMeta := ctx.bufs.permMeta[:n]
	for idx, c := range attackers {
		permMeta[idx] = attackerMetaPtrFor(c)
		pcBuf[idx] = CardState{Card: c, FromArsenal: idx == ctx.arsenalInIdx}
	}

	best := 0
	bestLeftoverRunechants := ctx.runechantCarryover
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
	// eval runs the active attack permutation against every pitch ordering it can legally
	// pair with — initial ordering plus Heap's enumeration over attackPitchPerm. The pitch
	// Heap walks indices [0, pn) so 0 / 1 pitch counts naturally collapse to the single
	// initial call without entering the inner loop body. The pitch-perm body is inlined
	// twice (initial + post-swap) to keep the closure capture set small — single closure
	// for the outer Heap to drive.
	eval := func() {
		dmg, leftoverRunechants, _, legal := ctx.playSequenceWithMeta(n)
		if ctx.cacheable && !state.IsCacheable() {
			ctx.cacheable = false
		}
		if legal && (!foundLegal || dmg > best ||
			(dmg == best && leftoverRunechants > bestLeftoverRunechants)) {
			best = dmg
			bestLeftoverRunechants = leftoverRunechants
			foundLegal = true
			ctx.carryWinner.SnapshotFromTurn(state)
		}
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
				dmg, leftoverRunechants, _, legal := ctx.playSequenceWithMeta(n)
				if ctx.cacheable && !state.IsCacheable() {
					ctx.cacheable = false
				}
				if legal && (!foundLegal || dmg > best ||
					(dmg == best && leftoverRunechants > bestLeftoverRunechants)) {
					best = dmg
					bestLeftoverRunechants = leftoverRunechants
					foundLegal = true
					ctx.carryWinner.SnapshotFromTurn(state)
				}
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
	return best, bestLeftoverRunechants, foundLegal
}

// playSequence plays `order` as a sequence of cards, reusing ctx.bufs' pooled buffers.
// Buffers are mutated in place; the caller must not read them concurrently.
//
// Runechant flow:
//   - state.Runechants starts at ctx.runechantCarryover.
//   - Play / OnCardPlayed calling CreateRunechants increments the count AND returns n damage
//     — tokens are credited exactly once, at creation.
//   - After each Attack / Weapon card resolves, all current tokens fire and are destroyed;
//     state.Runechants is zeroed but damage is NOT re-added (tokens were credited at
//     creation).
//   - At end of the sequence, state.Runechants is the leftover count carrying into next turn.
//
// Resource flow lives on playSequenceWithMeta; this wrapper just forwards.
//
// Populates permMeta from order and then calls playSequenceWithMeta. The hot path
// (bestSequence) builds meta once and calls playSequenceWithMeta directly to amortise
// interface dispatch across the N! permutations.
func (ctx *sequenceContext) playSequence(order []Card) (damage int, leftoverRunechants int, residualBudget int, legal bool) {
	n := len(order)
	pcBuf := ctx.bufs.pcBuf
	meta := ctx.bufs.permMeta[:n]
	for i, c := range order {
		meta[i] = attackerMetaPtrFor(c)
		pcBuf[i] = CardState{Card: c, FromArsenal: i == ctx.arsenalInIdx}
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
// Resource flow: the chain's pitch pool is ctx.attackPitchPerm in its current Heap-permuted
// order. Each chain step pays its cost by drawing from the front of the pool — popping new
// pitches off attackPitchPerm one at a time when the front exhausts. Every pitched card
// whose resources contribute (even partially) to a step's payment lands in that step's
// PitchedToPlay slice. So pitching one 3-resource non-attack to fund three 1-cost plays
// attributes the non-attack to all three, not just the one whose payment popped it.
// Leftover front-of-pool resources between steps roll forward — they aren't wasted, just
// reused. At end of chain every pitched CARD must have popped (FaB's pitch-timing rule);
// a leftover pitched card rejects the permutation, but residual carry from the last
// popped pitch is fine (it's surplus, not a held-back pitch).
//
// ctx.resourceBudget seeds the front-of-pool resources without any backing card. In
// production it's always 0 (real pitches fund every step). Tests set it for synthetic-
// budget plays that bypass the pitch pool entirely.
//
// Damage flows through state.Value: the dispatcher records the chain step's
// Play+BonusAttack contribution via state.AddLogEntry; pre-trigger handlers (hero, aura)
// credit themselves through AddPreTriggerLogEntry, post-trigger handlers (ephemeral)
// through AddPostTriggerLogEntry. The returned damage is just state.Value at end of chain.
func (ctx *sequenceContext) playSequenceWithMeta(n int) (damage int, leftoverRunechants int, residualBudget int, legal bool) {
	pcBuf := ctx.bufs.pcBuf
	ptrBuf := ctx.bufs.ptrBuf
	meta := ctx.bufs.permMeta[:n]
	for i := 0; i < n; i++ {
		pcBuf[i].GrantedGoAgain = false
		pcBuf[i].BonusAttack = 0
		pcBuf[i].PitchedToPlay = nil
	}
	played := ptrBuf[:n]
	// Per-permutation reset: full-state rewrite. Hand and Deck are deep-copied so cards can
	// mutate them freely without leaking to the next permutation. state.Value resets to 0.
	ctx.resetStateForPermutation()
	state := ctx.bufs.state
	pitchPerm := ctx.attackPitchPerm
	pitchVals := ctx.attackPitchVals
	pitchIdx := 0
	pn := len(pitchPerm)
	// frontCard / frontRemaining track the unfunded balance of the partially-consumed
	// pitched card carried over from a previous chain step. Between chain steps either
	// (a) front is empty (frontCard==nil && frontRemaining==0) or (b) one pitched card sits
	// at the front with leftover resources. Within a single chain step's payment, the
	// front may pop and be replaced by the next pitch as carry exhausts. The synthetic
	// budget path (resourceBudget != 0, tests with no real pitches) seeds frontRemaining
	// with no backing card so attribution stays empty for those test plays.
	var frontCard Card
	frontRemaining := ctx.resourceBudget
	// attrBuf is the per-permutation flat backing for pc.PitchedToPlay slices. Pre-sized
	// once at construction to handSize+1 so append never reallocates — the slice headers
	// pcBuf entries hold stay valid for the duration of this permutation.
	attrBuf := ctx.bufs.pitchAttrBuf[:0]
	for i, pc := range played {
		m := meta[i]
		cost := m.costAt(state)
		attrStart := len(attrBuf)
		// Pay this step's cost by drawing from the front of the pitch pool. Any pitched
		// card whose resources contribute even partially to this step lands in the
		// attribution slice — including the front carrying over from a prior step. So
		// pitching one Malefic (3) to fund three 1-cost plays attributes Malefic to all
		// three, not just the one whose payment popped it off the deck.
		remaining := cost
		for remaining > 0 {
			if frontCard == nil && frontRemaining == 0 {
				if pitchIdx >= pn {
					return 0, 0, 0, false
				}
				frontCard = pitchPerm[pitchIdx]
				frontRemaining = pitchVals[pitchIdx]
				pitchIdx++
			}
			if frontCard != nil {
				attrBuf = append(attrBuf, frontCard)
			}
			if frontRemaining > remaining {
				frontRemaining -= remaining
				remaining = 0
			} else {
				remaining -= frontRemaining
				frontRemaining = 0
				frontCard = nil
			}
		}
		pc.PitchedToPlay = attrBuf[attrStart:len(attrBuf)]

		state.CardsRemaining = played[i+1:]

		// If this card is an attack or weapon and any Runechant is live, those tokens fire on
		// its damage step. Set ArcaneDamageDealt now — before Play and OnCardPlayed — so Play
		// effects that read "if you've dealt arcane damage this turn" see the flag for same-hand
		// triggers. Cards that deal arcane damage via their Play text flip the flag themselves
		// through DealArcaneDamage. The flip is gated on LikelyDamageHits — same hit-likelihood
		// model the rest of the arcane plumbing uses — so a single live runechant on an attack
		// satisfies the gate but a live cluster the model treats as blockable doesn't.
		isAttackOrWeapon := m.isAttackOrWeapon
		if isAttackOrWeapon && state.Runechants > 0 && LikelyDamageHits(state.Runechants, false) {
			state.ArcaneDamageDealt = true
		}

		// Hero ability fires BEFORE the card's own Play so "aura created this turn" checks
		// inside the card's Play see the runechant (or other aura) the hero just made.
		// Viserai's "another non-attack action" gate still excludes the current card because
		// NonAttackActionPlayed isn't flipped until the end of the iteration. The hero
		// handler logs its own contribution via state.AddPreTriggerLogEntry; its int return
		// is unused.
		ctx.hero.OnCardPlayed(pc.Card, state)
		ephemeralsBefore := len(state.EphemeralAttackTriggers)
		// Card.Play owns its chain-step log line and value contribution: it calls
		// state.ApplyAndLogEffectiveAttack / LogPlay before returning. The dispatcher
		// just sequences the surrounding triggers.
		pc.Card.Play(state, pc)
		// Stamp SourceIndex on any EphemeralAttackTriggers the card registered during Play
		// so fireEphemeralAttackTriggers can attribute the fire back to this card.
		for k := ephemeralsBefore; k < len(state.EphemeralAttackTriggers); k++ {
			state.EphemeralAttackTriggers[k].SourceIndex = i
		}
		// Log order matches FaB's stack-resolution order, not the dispatcher's call order:
		// hero / aura triggers fire when the card is played (LL1), go on top of the stack,
		// and resolve before the card itself. Ephemeral "if hits" triggers fire after the
		// attack lands and log below. Each trigger handler authors its own log line via
		// AddPreTriggerLogEntry / AddPostTriggerLogEntry, with LogEntry.Source naming the
		// triggering card so appendGroupedChainEntries clusters the trigger underneath the
		// chain entry that names that card.
		if m.isAttackAction {
			fireAttackActionTriggers(state, pc.Card)
			// Fire ephemeral triggers AFTER hero and aura triggers so the handler sees the
			// fully-resolved attacker state (Dominate grants, hero-created auras, fresh
			// Runechants from aura triggers).
			fireEphemeralAttackTriggers(state, pc)
		}
		state.CardsPlayed = append(state.CardsPlayed, pc.Card)
		if m.types.IsNonAttackAction() {
			state.NonAttackActionPlayed = true
		}
		// Weapons and persistent card types (Auras, Items) stay in their zone when they
		// resolve; any destroy event that should send them to the graveyard is a separate
		// trigger. Everything else — Actions, Attack Reactions, Defense Reactions, Blocks,
		// Instants — heads to the graveyard immediately. Direct field write — the framework
		// driving the chain isn't a card-driven content read, so no cacheable flip.
		if !m.types.PersistsInPlay() {
			state.graveyard = append(state.graveyard, pc.Card)
		}

		// Attacks and weapon swings consume all runechants in play. Damage isn't re-added: each
		// token was credited +1 at creation time, so this is pure state cleanup.
		if isAttackOrWeapon {
			state.Runechants = 0
		}

		if i < n-1 && !(m.baseGoAgain || pc.GrantedGoAgain) {
			return 0, 0, 0, false
		}
	}

	// Pitch-timing rule: every Pitch-role card must have paid for something on the stack. If
	// the chain finished with pitches still queued, one of them was held back without funding
	// any cost — illegal in FaB. Leftover carry on the front (frontRemaining > 0) is fine —
	// that's the over-pitch surplus on the last popped card, not a held-back pitch.
	if pitchIdx < pn {
		return 0, 0, 0, false
	}
	return state.Value, state.Runechants, frontRemaining, true
}
