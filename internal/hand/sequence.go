package hand

// Attack-chain search: bestAttackWithWeapons evaluates one partition leaf across all phase /
// weapon masks, bestSequence picks the best ordering of attackers via Heap's algorithm, and
// playSequence* replay a single permutation through TurnState while firing hero triggers and
// AuraTrigger / EphemeralAttackTrigger handlers.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Phase masks: when no Defense Reactions are present (or no pitches exist), all pitches go to
// the attack phase, so we visit one configuration. Otherwise we enumerate 2^|pitched| splits.
func bestAttackWithWeapons(hero hero.Hero, weapons []weapon.Weapon, attackers, defenders, pitched, deck []card.Card, bufs *attackBufs, runechantCarryover, incomingDamage, blockTotal, arsenalInIdx int, priorAuraTriggers []card.AuraTrigger) (int, int, int, chainBudget, []string, bool) {
	ctx := &sequenceContext{
		hero:               hero,
		pitched:            pitched,
		deck:               deck,
		bufs:               bufs,
		runechantCarryover: runechantCarryover,
		incomingDamage:     incomingDamage,
		blockTotal:         blockTotal,
		arsenalInIdx:       arsenalInIdx,
		priorAuraTriggers:  priorAuraTriggers,
		// Borrow bufs' pre-sized winner scratch so the eval closure's append-winner step reuses
		// one backing array per Best call instead of allocating per sequenceContext.
		drawnWinner:        bufs.drawnWinnerScratch[:0],
		auraTriggersWinner: bufs.auraTriggersWinnerScratch[:0],
	}
	// Hoist leaf-constant TurnState fields out of the per-permutation reset in
	// playSequenceWithMeta.
	ctx.seedState()

	// Defense Reactions fire independently of ordering and attack chain (each sees a fresh
	// TurnState with only Pitched + Deck), so their Play-return damage is constant across phase /
	// weapon masks. Compute it once; reseed ctx state for the attack chain afterwards.
	hasDRs := containsDefenseReaction(defenders)
	var defenseDealt int
	if hasDRs {
		defenseDealt, bufs.defenseGravScratch = defenseReactionDamage(defenders, pitched, deck, bufs.state, bufs.defenseGravScratch, &bufs.drCardStateScratch)
		ctx.seedState()
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
			dealt, leftoverRunechants, legal := ctx.bestSequence(allAttackers, nil, nil, nil, nil)
			if !legal {
				continue
			}
			// Cost the DRs against the chain's final runechant count. DRs with variable cost
			// read state.Runechants inside their Cost; static DRs return a constant. Reuse
			// bufs.drScratch instead of allocating a fresh TurnState per mask iteration — the
			// interface call boxes the pointer, so a stack allocation would escape and heap-alloc
			// every loop.
			bufs.drScratch = card.TurnState{Runechants: leftoverRunechants}
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
				foundFeasible = true
			}
		}
	}

	if !foundFeasible {
		return 0, 0, 0, chainBudget{}, nil, false
	}
	return bestDealt, defenseDealt, bestLeftoverRunechants, bestBudget, bestSwung, true
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
	hero               hero.Hero
	pitched, deck      []card.Card
	bufs               *attackBufs
	resourceBudget     int
	runechantCarryover int
	incomingDamage     int
	blockTotal         int
	hasAttackPitches   bool
	maxAttackPitch     int
	// arsenalInIdx is the index in the attackers slice (the slice passed to bestSequence) of
	// the card that came from the arsenal slot at start of turn, or -1 when no arsenal-in card
	// is in the chain. Lets bestSequence flag the matching pcBuf entry's FromArsenal as the
	// permutation moves it around.
	arsenalInIdx int
	// priorAuraTriggers are the AuraTriggers carried in from the previous turn (e.g. an
	// AttackAction trigger from a Malefic Incantation played a turn ago). Each permutation
	// seeds state.AuraTriggers with a fresh copy of this slice so mid-chain firing can
	// decrement Count / set FiredThisTurn without leaking those mutations across permutations.
	priorAuraTriggers []card.AuraTrigger
	// drawnWinner snapshots the winning permutation's drawn cards so fillContributions can
	// surface them on summary.Drawn. Populated from state.Drawn during the winner path because
	// Heap's algorithm keeps iterating after the winner is chosen and state.Drawn reflects the
	// last permutation's draws. Every entry in drawnWinner is assigned Role=Held by the
	// caller (post-Best, promoteRandomHeldToArsenal may flip one to Arsenal).
	drawnWinner []card.Card
	// auraTriggersWinner snapshots the winning permutation's final state.AuraTriggers so the
	// deck loop can carry them into next turn. Includes both inherited triggers from
	// priorAuraTriggers (with mutated Count / FiredThisTurn) and ones added by Play.
	auraTriggersWinner []card.AuraTrigger
}

// fireAttackActionTriggers walks state.AuraTriggers after an attack action card resolves and
// invokes every TriggerAttackAction entry whose OncePerTurn gate is open. Each fire
// decrements the trigger's Count; when Count hits zero the aura drops out of the list and
// Self lands in the graveyard so downstream same-turn effects see the destroy. Returns the
// summed Damage from all fires, folded into chain damage by the caller.
//
// Slice mutation: a survivors prefix is built in place over the existing slice; entries
// kept after firing are written back at increasing indices, exhausted ones are skipped.
func fireAttackActionTriggers(state *card.TurnState) int {
	total := 0
	triggers := state.AuraTriggers
	dst := triggers[:0]
	for i := range triggers {
		t := triggers[i]
		if t.Type != card.TriggerAttackAction || (t.OncePerTurn && t.FiredThisTurn) {
			dst = append(dst, t)
			continue
		}
		total += t.Handler(state)
		t.FiredThisTurn = true
		t.Count--
		if t.Count <= 0 {
			state.AddToGraveyard(t.Self)
			continue
		}
		dst = append(dst, t)
	}
	state.AuraTriggers = dst
	return total
}

// fireEphemeralAttackTriggers walks state.EphemeralAttackTriggers after an attack action
// card resolves and invokes every entry whose Matches predicate accepts the attacker. Each
// fire consumes the trigger (fire-once semantics) and routes its damage to the source's
// perCardOut slot via SourceIndex — Mauvrion Skies's "if hits" Runechants, for instance,
// surface on Mauvrion's BestLine entry rather than the attacker's. Non-matching entries stay
// in the slice for a later attack action; anything still in the list at end of chain
// fizzles silently (no graveyard bookkeeping — the source was already graveyarded when its
// own Play resolved).
//
// Slice mutation parallels fireAttackActionTriggers: a survivors prefix is built in place
// over the existing slice, with fired entries skipped.
func fireEphemeralAttackTriggers(state *card.TurnState, target *card.CardState, perCardOut []float64) int {
	total := 0
	triggers := state.EphemeralAttackTriggers
	dst := triggers[:0]
	for i := range triggers {
		t := triggers[i]
		if t.Matches != nil && !t.Matches(target) {
			dst = append(dst, t)
			continue
		}
		dmg := t.Handler(state, target)
		total += dmg
		if perCardOut != nil && t.SourceIndex >= 0 && t.SourceIndex < len(perCardOut) {
			perCardOut[t.SourceIndex] += float64(dmg)
		}
	}
	state.EphemeralAttackTriggers = dst
	return total
}

// seedState writes the leaf-constant TurnState fields (pitched / deck refs, incoming damage,
// block total) so the per-permutation reset in playSequenceWithMeta can skip them.
func (ctx *sequenceContext) seedState() {
	s := ctx.bufs.state
	s.Pitched = ctx.pitched
	s.Deck = ctx.deck
	s.IncomingDamage = ctx.incomingDamage
	s.BlockTotal = ctx.blockTotal
}

// bestSequence tries every ordering of attackers and returns the max total damage plus the
// runechant count at the end of the winning permutation. Between each card's Play() and its
// append to CardsPlayed, the hero's OnCardPlayed hook fires so triggered abilities contribute.
// legal=true when at least one ordering is playable; false when every permutation is rejected
// by playSequenceWithMeta's resource / go-again / pitch-waste checks.
//
// Uses Heap's algorithm (iterative) — no closure/callback alloc, no recursive call per perm.
//
// When winnerOrderOut is non-nil (len >= len(attackers)) the winning permutation is copied into
// it. perCardOut / perCardTriggerOut / perCardAuraTriggerOut (same size rule) receive the
// winning line's per-card Play damage, hero-trigger damage, and mid-chain aura-trigger damage.
// fillContributions uses these; the partition-loop caller passes nil for all four so the
// permutation search stays allocation-free.
func (ctx *sequenceContext) bestSequence(attackers, winnerOrderOut []card.Card, perCardOut, perCardTriggerOut, perCardAuraTriggerOut []float64) (int, int, bool) {
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
		pcBuf[idx] = card.CardState{Card: c, FromArsenal: idx == ctx.arsenalInIdx}
	}

	// Scratch buffers are playSequence's per-card outputs, overwritten every permutation. On a
	// new winner we copy them into the caller's perCardOut / perCardTriggerOut /
	// perCardAuraTriggerOut. Only populated when the caller asked to track.
	var scratch, triggerScratch, auraTriggerScratch []float64
	if perCardOut != nil {
		scratch = ctx.bufs.perCardScratch[:n]
	}
	if perCardTriggerOut != nil {
		if cap(ctx.bufs.perCardTriggerScratch) < n {
			ctx.bufs.perCardTriggerScratch = make([]float64, n)
		}
		triggerScratch = ctx.bufs.perCardTriggerScratch[:n]
	}
	if perCardAuraTriggerOut != nil {
		if cap(ctx.bufs.perCardAuraTriggerScratch) < n {
			ctx.bufs.perCardAuraTriggerScratch = make([]float64, n)
		}
		auraTriggerScratch = ctx.bufs.perCardAuraTriggerScratch[:n]
	}

	best := 0
	bestLeftoverRunechants := ctx.runechantCarryover
	foundLegal := false
	ctx.drawnWinner = ctx.drawnWinner[:0]
	ctx.auraTriggersWinner = ctx.auraTriggersWinner[:0]
	eval := func() {
		dmg, leftoverRunechants, _, legal := ctx.playSequenceWithMeta(n, scratch, triggerScratch, auraTriggerScratch)
		if !legal {
			return
		}
		if !foundLegal || dmg > best ||
			(dmg == best && leftoverRunechants > bestLeftoverRunechants) {
			best = dmg
			bestLeftoverRunechants = leftoverRunechants
			foundLegal = true
			ctx.drawnWinner = append(ctx.drawnWinner[:0], ctx.bufs.state.Drawn...)
			ctx.auraTriggersWinner = append(ctx.auraTriggersWinner[:0], ctx.bufs.state.AuraTriggers...)
			if winnerOrderOut != nil {
				for i := 0; i < n; i++ {
					winnerOrderOut[i] = pcBuf[i].Card
				}
			}
			if perCardOut != nil {
				copy(perCardOut[:n], scratch)
			}
			if perCardTriggerOut != nil {
				copy(perCardTriggerOut[:n], triggerScratch)
			}
			if perCardAuraTriggerOut != nil {
				copy(perCardAuraTriggerOut[:n], auraTriggerScratch)
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
// When perCardOut is non-nil (len >= n) each entry is the card's Play return for that
// position; perCardTriggerOut (same size rule) receives the hero's OnCardPlayed return;
// perCardAuraTriggerOut (same size rule) receives the mid-chain AuraTrigger return (e.g.
// Malefic Incantation's TriggerAttackAction). The hot partition-loop callers pass nil.
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
func (ctx *sequenceContext) playSequence(order []card.Card, perCardOut, perCardTriggerOut, perCardAuraTriggerOut []float64) (damage int, leftoverRunechants int, residualBudget int, legal bool) {
	ctx.seedState()
	n := len(order)
	pcBuf := ctx.bufs.pcBuf
	meta := ctx.bufs.permMeta[:n]
	for i, c := range order {
		meta[i] = attackerMetaPtrFor(c)
		pcBuf[i] = card.CardState{Card: c, FromArsenal: i == ctx.arsenalInIdx}
	}
	return ctx.playSequenceWithMeta(n, perCardOut, perCardTriggerOut, perCardAuraTriggerOut)
}

// playSequenceWithMeta runs the permutation currently held in ctx.bufs.pcBuf[:n] with
// aligned permMeta[:n]. CardState (Card + FromArsenal) persists across permutations, so any
// field a prior card's Play flips on a future card needs a per-permutation reset:
// GrantedGoAgain (next-attack go-again grants) and BonusAttack (next-attack +N{p} grants).
//
// Per-card output attribution:
//   - perCardOut[i] = card's own Play return (plus any EphemeralAttackTrigger damage routed
//     back to this slot via SourceIndex — ephemeral triggers credit their source card, not
//     the attacker that happened to consume them).
//   - perCardTriggerOut[i] = hero.OnCardPlayed return for this card.
//   - perCardAuraTriggerOut[i] = fireAttackActionTriggers return fired by this card. Credited
//     to the attack action card that resolved because prior-turn auras have no BestLine entry
//     this turn — attributing to the aura would drop the damage off the chain display.
func (ctx *sequenceContext) playSequenceWithMeta(n int, perCardOut, perCardTriggerOut, perCardAuraTriggerOut []float64) (damage int, leftoverRunechants int, residualBudget int, legal bool) {
	pcBuf := ctx.bufs.pcBuf
	ptrBuf := ctx.bufs.ptrBuf
	meta := ctx.bufs.permMeta[:n]
	for i := 0; i < n; i++ {
		pcBuf[i].GrantedGoAgain = false
		pcBuf[i].BonusAttack = 0
		if perCardOut != nil {
			perCardOut[i] = 0
		}
		if perCardTriggerOut != nil {
			perCardTriggerOut[i] = 0
		}
		if perCardAuraTriggerOut != nil {
			perCardAuraTriggerOut[i] = 0
		}
	}
	played := ptrBuf[:n]
	state := ctx.bufs.state
	// Per-permutation reset. Only touch fields the cards mutate; leaf-stable fields (Pitched,
	// Deck, IncomingDamage, BlockTotal) come from seedState. A full-struct replace here
	// memcpies big slice headers on every permutation and profiles dramatically slower.
	state.CardsPlayed = ctx.bufs.cardsPlayedBuf[:0]
	state.Runechants = ctx.runechantCarryover
	state.ArcaneDamageDealt = false
	state.AuraCreated = false
	state.Overpower = false
	state.NonAttackActionPlayed = false
	// Deck and Drawn reset per permutation: DrawOne mutates them, so a prior permutation's
	// consumption would poison the next.
	state.Deck = ctx.deck
	state.Drawn = nil
	// Graveyard and Banish reset per permutation: cards append themselves to Graveyard as
	// they resolve, and graveyard-banish effects shift cards into Banish. Reusing the scratch
	// backing array keeps the reset allocation-free.
	state.Graveyard = ctx.bufs.attackGravScratch[:0]
	state.Banish = nil
	// AuraTriggers reset per permutation: seeded with a copy of priorAuraTriggers so
	// mid-chain attack-action triggers can fire without their Count / FiredThisTurn
	// mutations leaking across permutations. Cards adding triggers via AddAuraTrigger extend
	// the same scratch slice.
	state.AuraTriggers = append(ctx.bufs.auraTriggersScratch[:0], ctx.priorAuraTriggers...)
	// EphemeralAttackTriggers reset per permutation as empty — fire-once triggers never
	// carry across turns, so there's nothing to seed from prior state.
	state.EphemeralAttackTriggers = ctx.bufs.ephemeralTriggersScratch[:0]
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
		// triggers. Cards that deal arcane damage via their Play text flip the flag themselves.
		isAttackOrWeapon := m.isAttackOrWeapon
		if isAttackOrWeapon && state.Runechants > 0 {
			state.ArcaneDamageDealt = true
		}

		// Hero ability fires BEFORE the card's own Play so "aura created this turn" checks
		// inside the card's Play see the runechant (or other aura) the hero just made.
		// Viserai's "another non-attack action" gate still excludes the current card because
		// NonAttackActionPlayed isn't flipped until the end of the iteration.
		triggerDmg := ctx.hero.OnCardPlayed(pc.Card, state)
		ephemeralsBefore := len(state.EphemeralAttackTriggers)
		playDmg := pc.Card.Play(state, pc)
		// Stamp SourceIndex on any EphemeralAttackTriggers the card registered during Play
		// so fireEphemeralAttackTriggers can route their damage back to this card's
		// perCardOut slot.
		for k := ephemeralsBefore; k < len(state.EphemeralAttackTriggers); k++ {
			state.EphemeralAttackTriggers[k].SourceIndex = i
		}
		auraTriggerDmg := 0
		ephemeralDmg := 0
		if m.isAttackAction {
			auraTriggerDmg = fireAttackActionTriggers(state)
			// Fire ephemeral triggers AFTER hero and aura triggers so the handler sees the
			// fully-resolved attacker state (Dominate grants, hero-created auras, fresh
			// Runechants from aura triggers). Damage is routed back to each trigger's
			// source via SourceIndex, so perCardOut is updated in place inside the helper.
			ephemeralDmg = fireEphemeralAttackTriggers(state, pc, perCardOut)
		}
		// BonusAttack is granted by a prior card's "next attack +N{p}" rider. The grant is
		// folded in here (not by the target's Play) so the +N is attributed to the attack
		// receiving the buff rather than the granter, and so any "if this hits" rider
		// inside the target's Play can read self.EffectiveAttack() consistently. Applied
		// unconditionally — picking who can legally receive a +N{p} grant is the grantor's
		// responsibility (it scans CardsRemaining and matches the appropriate type, e.g.
		// attack actions for Come to Fight, weapon swings for Brandish), and a future card
		// that grants damage to a non-attack source shouldn't have to fight a solver-side
		// type gate.
		//
		// Reuses pc.EffectiveAttack() for the printed-power-plus-bonus calculation
		// (which already enforces the FaB attack-power floor: a -3 grant on a 1-power
		// attack resolves as 0, not -2) and folds in any rider damage Play returned
		// beyond Card.Attack — that "extra" component sits outside the attack-power
		// floor so it's added on top of the clamped value.
		cardContrib := pc.EffectiveAttack() + (playDmg - pc.Card.Attack())
		if cardContrib < 0 {
			cardContrib = 0
		}
		damage += cardContrib + triggerDmg + auraTriggerDmg + ephemeralDmg
		if perCardOut != nil {
			perCardOut[i] = float64(cardContrib)
		}
		if perCardTriggerOut != nil {
			perCardTriggerOut[i] = float64(triggerDmg)
		}
		if perCardAuraTriggerOut != nil {
			perCardAuraTriggerOut[i] = float64(auraTriggerDmg)
		}
		state.CardsPlayed = append(state.CardsPlayed, pc.Card)
		if m.types.IsNonAttackAction() {
			state.NonAttackActionPlayed = true
		}
		// Weapons and persistent card types (Auras, Items) stay in their zone when they
		// resolve; any destroy event that should send them to the graveyard is a separate
		// trigger. Everything else — Actions, Attack Reactions, Defense Reactions, Blocks,
		// Instants — heads to the graveyard immediately.
		if !m.types.PersistsInPlay() {
			state.Graveyard = append(state.Graveyard, pc.Card)
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

	// Mid-turn-drawn cards always carry to the next hand as Held or compete for the empty
	// arsenal slot; they never pitch or extend the chain. If they could, the solver's best line
	// would depend on Deck[0], which the player commits before the draw reveals.

	// Pitch-timing rule: every Pitch-role card must have paid for something on the stack. If the
	// chain's leftover budget is at least the max attack-phase pitch, one pitch could have been
	// Held instead — this permutation violates FaB's rules.
	if ctx.hasAttackPitches && resources >= ctx.maxAttackPitch {
		return 0, 0, 0, false
	}
	return damage, state.Runechants, resources, true
}
