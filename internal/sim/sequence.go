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
	}
	// Defenders fire independently of ordering and attack chain — DRs through Play, plain
	// blocks as raw block credit — so their total Value contribution is constant across phase
	// / weapon masks. Compute it once. Includes DR blocks + arcane / runechant riders + plain-
	// block residual against the partition's incoming damage; over-blocked excess is discarded
	// by the per-card cap.
	hasDRs := containsDefenseReaction(defenders)
	var defenseDealt int
	var defenseUncacheable bool
	if len(defenders) > 0 {
		defenseDealt, bufs.defenseGravScratch, defenseUncacheable = defendersDamage(defenders, pitched, deck, bufs.state, bufs.defenseGravScratch, &bufs.drCardStateScratch, incomingDamage, arsenalDefenderIdx)
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
	var bestCarry CarryState
	foundFeasible := false

	for pmask := 0; pmask < phaseCount; pmask++ {
		phase := splitPitchesAcrossPhases(pitchedVals, pmask, phaseCount)

		ctx.resourceBudget = phase.attackBudget
		ctx.hasAttackPitches = phase.hasAttackPitches
		ctx.maxAttackPitch = phase.maxAttackPitch

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
				bestCarry = ctx.carryWinner
				foundFeasible = true
			}
		}
	}

	if !foundFeasible {
		// No-feasible-line leaves still surface the defense-phase cacheable bit — DR Plays
		// ran independently of the (rejected) attack chain so a DR that read graveyard
		// poisons the result regardless of attack-feasibility.
		return 0, 0, 0, chainBudget{}, nil, CarryState{}, false, defenseUncacheable
	}
	return bestDealt, defenseDealt, bestLeftoverRunechants, bestBudget, bestSwung, bestCarry, true, ctx.uncacheable || defenseUncacheable
}

// sequenceContext carries the stable per-partition-leaf environment: hero (for OnCardPlayed
// triggers), pitched / deck refs for Card.Play, shared scratch buffers, and the numeric budgets
// that persist across permutation and mask iterations. Built once per leaf so the hot inner
// calls (playSequence, bestSequence) shrink to their varying inputs and tracking outputs.
//
// resourceBudget / hasAttackPitches / maxAttackPitch are rewritten by bestAttackWithWeapons on
// each phase-mask iteration: they fund the attack chain and let playSequenceWithMeta reject
// permutations whose final residual breaks FaB's pitch-timing rule (excess >= max pitch means
// one pitch could have been Held instead).
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
	resourceBudget      int
	runechantCarryover  int
	incomingDamage      int
	blockTotal          int
	hasAttackPitches    bool
	maxAttackPitch      int
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
	// carryWinner snapshots the winning permutation's end-of-chain TurnState slice — every
	// field that survives the turn boundary (Hand, Deck, Arsenal, Graveyard, Banish,
	// Runechants, AuraTriggers). Heap's algorithm keeps iterating past the winner and the
	// shared state.* fields reflect whatever ordering ran last, so the snapshot has to
	// happen the moment a new winner is found.
	carryWinner CarryState
	// skipLog propagates into TurnState.SkipLog on every permutation reset. When true,
	// chains run with Log appends elided (Value still credited); the caller is replaying
	// later with skipLog=false to materialise the printout.
	skipLog bool
	// uncacheable is a sticky bit ORed in after every permutation in bestSequence. Any
	// permutation reporting !state.IsCacheable() at chain end pins the leaf as uncacheable
	// — once a card in any sibling chain reads hidden state, the partition's output isn't
	// safe to cache. Carries across phase / weapon masks within the same leaf because the
	// solver explores all configurations and the cache key would have to disambiguate
	// which the winner came from.
	uncacheable bool
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
			state.AddToGraveyard(t.Self)
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
		// No attackers means no chain costs are deducted — the attack phase spends zero from the
		// budget. If any attack-phase pitches exist they over-pay (residual == budget >= maxPitch
		// since the budget is the sum of those pitches); pitch-timing fails.
		if ctx.hasAttackPitches && ctx.resourceBudget >= ctx.maxAttackPitch {
			return 0, 0, false
		}
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
	ctx.carryWinner = CarryState{}
	state := ctx.bufs.state
	eval := func() {
		dmg, leftoverRunechants, _, legal := ctx.playSequenceWithMeta(n)
		// OR in this permutation's cacheable status before the legality gate — even an
		// illegal permutation that read hidden state during Play poisons the leaf, since
		// the solver's pre-screens didn't catch it and a real run would have reached that
		// read. Short-circuited once already poisoned to avoid the field read on
		// already-uncacheable hands.
		if !ctx.uncacheable && !state.IsCacheable() {
			ctx.uncacheable = true
		}
		if !legal {
			return
		}
		if !foundLegal || dmg > best ||
			(dmg == best && leftoverRunechants > bestLeftoverRunechants) {
			best = dmg
			bestLeftoverRunechants = leftoverRunechants
			foundLegal = true
			ctx.carryWinner = snapshotCarry(state)
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
// Resource flow: ctx.resourceBudget is the starting pool; each card deducts
// attackerMeta.costAt(state). Negative remaining budget returns legal=false.
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
// GrantedGoAgain (next-attack go-again grants) and BonusAttack (next-attack +N{p} grants).
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
	}
	played := ptrBuf[:n]
	// Per-permutation reset: full-state rewrite. Hand and Deck are deep-copied so cards can
	// mutate them freely without leaking to the next permutation. state.Value resets to 0.
	ctx.resetStateForPermutation()
	state := ctx.bufs.state
	resources := ctx.resourceBudget
	for i, pc := range played {
		m := meta[i]
		cost := m.costAt(state)
		resources -= cost
		if resources < 0 {
			return 0, 0, 0, false
		}

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

	// Pitch-timing rule: every Pitch-role card must have paid for something on the stack. If the
	// chain's leftover budget is at least the max attack-phase pitch, one pitch could have been
	// Held instead — this permutation violates FaB's rules.
	if ctx.hasAttackPitches && resources >= ctx.maxAttackPitch {
		return 0, 0, 0, false
	}
	return state.Value, state.Runechants, resources, true
}

// snapshotCarry copies every persistent TurnState field that survives the turn boundary into
// a CarryState. The slice copies are intentional: mid-chain state.* slices alias attackBufs
// scratch storage and the next permutation will overwrite them. The deck loop adopts these
// slices wholesale into the next-turn state. Reads s.deck / s.graveyard directly so the
// snapshot itself doesn't poison cacheable — the per-permutation cacheable check ran before
// snapshotCarry, and resetStateForPermutation clears the bit before the next permutation.
func snapshotCarry(s *TurnState) CarryState {
	return CarryState{
		Hand:         append([]Card(nil), s.Hand...),
		Deck:         append([]Card(nil), s.deck...),
		Arsenal:      s.Arsenal,
		Graveyard:    append([]Card(nil), s.graveyard...),
		Banish:       append([]Card(nil), s.Banish...),
		Runechants:   s.Runechants,
		AuraTriggers: append([]AuraTrigger(nil), s.AuraTriggers...),
		Log:          append([]LogEntry(nil), s.Log...),
	}
}
