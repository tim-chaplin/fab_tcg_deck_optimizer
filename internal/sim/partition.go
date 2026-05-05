package sim

// Top-level hand enumeration: findBest walks every partition (Pitch / Attack / Defend /
// Held / Arsenal assignment) and delegates each leaf's chain-feasibility check to
// bestAttackWithWeapons. Post-enumeration helpers decide how an empty arsenal slot gets
// filled, plus the roleAllowed policy function that shapes the partition tree. Tiebreak
// order across partition leaves lives on runningCarry (see running_carry.go): Value →
// leftover runechants (future arcane) → cards drawn → more pending future value at end
// of chain (non-token aura Counts plus unspent item Counts — hidden later-turn payoff
// the current-turn Value misses) → arsenal slot ending occupied (saves a hand slot next
// refill; covers both arsenal-in-stayed and Held-for-promotion).

import ()

func (e *Evaluator) findBest(hero Hero, weapons []Weapon, hand []Card, mp Matchup, deck []Card, arsenalCardIn Card, priorAuras []Aura, priorItems []Item, skipLog bool) TurnSummary {
	// Cache fast-path. Bypassed when disabled (e.cache nil) or when any input overflows
	// a fixed-size cache-key slot (hand size, weapons, auras, items).
	var cacheKey evalCacheKey
	cacheUsable := e.cache != nil
	if cacheUsable {
		var keyOK bool
		cacheKey, keyOK = makeCacheKey(hero, weapons, hand, arsenalCardIn, priorAuras, priorItems)
		if !keyOK {
			cacheUsable = false
		}
	}
	if cacheUsable {
		if entry, ok := e.cache.lookup(cacheKey); ok {
			e.cache.hits.Add(1)
			return e.replayBest(entry, hero, weapons, hand, mp, deck, arsenalCardIn, priorAuras, priorItems, skipLog)
		}
		e.cache.misses.Add(1)
	}

	n := len(hand)
	// The partition recurse treats the arsenal-in card as an extra entry at index n with a
	// restricted role menu (Arsenal / Attack / Defend), so everything about it is decided inside
	// the enumeration. totalN is the effective size of BestLine.
	totalN := n
	if arsenalCardIn != nil {
		totalN = n + 1
	}

	// Seed best.State.Runechants() with the carryover: partitions with no attacks don't reduce
	// it, so carryover is the baseline to beat. BestLine starts with every hand card Held and
	// the arsenal-in card (if any) staying in the slot, so a hand with no Value-adding
	// partition still reports sensible "nothing played, nothing pitched" assignments.
	// Cacheable starts true: the no-feasible-line fallback ran no chain and read no hidden
	// state, so the seed is trivially cacheable. Each leaf's cacheable bit ANDs into a
	// running sticky and the final value stamps best.Cacheable.
	best := TurnSummary{
		BestLine:       make([]CardAssignment, totalN),
		IncomingDamage: mp.IncomingDamage,
		Cacheable:      true,
		State: CarryState{
			Hand:    append([]Card(nil), hand...),
			Deck:    append([]Card(nil), deck...),
			Arsenal: arsenalCardIn,
			Auras:   append([]Aura(nil), priorAuras...),
			Items:   append([]Item(nil), priorItems...),
		},
	}
	cacheable := true
	// bestSwung holds the winning partition's swung weapon names — surfaced on the summary so
	// the printout can list weapons that swung this turn (weapons have no BestLine entry).
	var bestSwung []string
	for i := 0; i < n; i++ {
		best.BestLine[i] = CardAssignment{Card: hand[i], Role: Held}
	}
	if arsenalCardIn != nil {
		best.BestLine[n] = CardAssignment{Card: arsenalCardIn, Role: Arsenal, FromArsenal: true}
	}

	// bufs is the pooled scratch space for this deck evaluation. Partition scratch is sized
	// handSize+1, big enough for totalN when an arsenal-in card inflates the effective hand.
	// Each field is re-sliced and rewritten below, so carry-over from prior calls can't leak.
	bufs := e.getAttackBufs(n, weapons)
	// running tracks the new-best leaf as the recurse explores partitions. It seeds with
	// the no-feasible-leaf fallback's state (arsenal-in stays, runechant carryover, every
	// hand card Held so a post-hoc promotion would fill arsenal) so the tiebreaker treats
	// the seed as a valid baseline. Finalize clones the scratch into best.State once at
	// the end so the returned TurnSummary owns independent backing.
	// Seed the running carry's scratch with the no-feasible-line fallback's hand —
	// the partition recurse's "all Held" baseline — so willOccupy reads the seed
	// state correctly when no chain has been promoted yet.
	bufs.findBestCarryScratch.Reset()
	bufs.findBestCarryScratch.Hand = append(bufs.findBestCarryScratch.Hand[:0], hand...)
	running := runningCarry{
		scratch:            &bufs.findBestCarryScratch,
		leftoverRunechants: tokenCountIn(priorAuras, TokenTypeRunechant),
		arsenal:            arsenalCardIn,
	}
	rolesBuf := bufs.rolesBuf[:totalN]
	pvals := bufs.pitchVals[:totalN]
	dvals := bufs.defenseVals[:totalN]
	isDR := bufs.isDRBuf[:totalN]
	canAttack := bufs.canAttackBuf[:totalN]

	fillPartitionPerCardBufs(hand, n, totalN, arsenalCardIn, pvals, dvals, isDR, canAttack)

	var recurse func(i, pitchSum, defenseSum int)
	recurse = func(i, pitchSum, defenseSum int) {
		if i == totalN {
			attackDealt, defenseDealt, leftoverRunechants, swung, carry, ok, leafCacheable, arsenalAtChainStart := e.evaluatePartition(
				hero, weapons, hand, deck, arsenalCardIn,
				rolesBuf, n, bufs,
				mp, defenseSum,
				priorAuras, priorItems, skipLog,
			)
			// Aggregate per leaf — an infeasible attack chain still surfaces its DR-side
			// reads (defendersDamage runs before the feasibility gate inside
			// bestAttackWithWeapons) so a DR scanning the graveyard pins cacheable=false
			// even when the partition's chain rejects.
			if !leafCacheable {
				cacheable = false
			}
			if !ok {
				return
			}

			v := attackDealt + defenseDealt
			arsenalCard := arsenalAtChainStart
			futureValuePlayed := pendingFutureValue(carry.Auras, carry.Items)
			// willOccupy reads end-of-chain carry rather than rolesBuf so chains that drew
			// a card mid-chain (Pitch + Gold-spend) register the drawn card as filling
			// next turn's arsenal — same outcome as a Held card promoting via the post-hoc
			// step.
			willOccupy := arsenalCard != nil || len(carry.Hand) > 0
			if !running.Beats(v, leftoverRunechants, futureValuePlayed, carry.CardsDrawn, willOccupy) {
				return
			}
			running.Promote(v, leftoverRunechants, futureValuePlayed, carry.CardsDrawn, arsenalCard, &carry)
			best.Value = v
			bestSwung = swung
			// Cards and FromArsenal flags were populated at construction; Role is the only
			// field that varies per-permutation.
			for j := 0; j < totalN; j++ {
				best.BestLine[j].Role = rolesBuf[j]
			}
			return
		}
		isArsenalSlot := i == n && arsenalCardIn != nil
		// Hand cards can't take Arsenal role (post-hoc promotion handles that). Cap the range at
		// Held for hand slots to skip the roleAllowed-rejection work for Arsenal.
		maxRole := Held
		if isArsenalSlot {
			maxRole = Arsenal
		}
		for r := Role(0); r <= maxRole; r++ {
			if !roleAllowed(r, isArsenalSlot, isDR[i], canAttack[i]) {
				continue
			}
			// FaB rule: defense reactions and plain blocks only happen during the defend step
			// of an attack chain. With 0 incoming there is no defend step, so the Defend role
			// is illegal for every card type — DR or otherwise. Skipping it here also avoids
			// no-op DR Plays leaving an empty hand to a downstream len(s.Hand)==0 gate.
			if r == Defend && mp.IncomingDamage == 0 {
				continue
			}
			rolesBuf[i] = r
			switch r {
			case Pitch:
				recurse(i+1, pitchSum+pvals[i], defenseSum)
			case Defend:
				recurse(i+1, pitchSum, defenseSum+dvals[i])
			case Attack, Held, Arsenal:
				recurse(i+1, pitchSum, defenseSum)
			}
		}
	}
	recurse(0, 0, 0)
	best.SwungWeapons = bestSwung
	// Promote the running winner's CarryState into best.State. Finalize clones the scratch
	// once so best.State owns independent backing, or leaves the seed in place when no
	// leaf produced a feasible chain.
	running.Finalize(&best.State)
	// Stamp Cacheable last from the AND-aggregated sticky bit so every leaf the search
	// touched (feasible or rejected) contributes. The post-hoc arsenal promotion below
	// doesn't run a chain, so it doesn't move the bit.
	best.Cacheable = cacheable
	// If the arsenal slot is empty after the chain runs, promote one card from State.Hand
	// into it (deterministic per-hand pick). State.Hand at this point holds the partition's
	// Held cards plus anything tutored mid-chain; both are equivalent future-turn value, so
	// the promotion picks across the combined pool.
	if best.State.Arsenal == nil {
		promoteRandomHandCardToArsenal(&best, hand, arsenalCardIn)
	}
	// Cache store happens after post-hoc arsenal promotion so the cached BestLine reflects
	// final roles (one Held entry may have flipped to Arsenal). Only stores when the chain
	// reported Cacheable=true at end of search — uncacheable results would be unsafe to
	// reuse for a future call with the same key but different deck contents.
	if cacheUsable {
		if best.Cacheable {
			e.cache.store(cacheKey, evalCacheEntry{
				line:         append([]CardAssignment(nil), best.BestLine...),
				swungWeapons: append([]string(nil), best.SwungWeapons...),
			})
		} else {
			e.cache.uncacheable.Add(1)
		}
	}
	return best
}

// promoteRandomHandCardToArsenal picks one card from best.State.Hand (the chain's end-of-turn
// hand — partition Held cards plus anything tutored mid-chain) and moves it into
// best.State.Arsenal, removing it from State.Hand. Deterministic per-hand pick (hashed from
// sorted starting-hand IDs + Hand IDs + arsenal-in ID) so equivalent inputs always promote
// the same card. No-op when State.Hand is empty.
//
// When the promoted card matches a Held entry in BestLine, that entry's Role flips to
// Arsenal so the per-card display still attributes the slot. Tutored cards (not in BestLine)
// just live in State.Arsenal without a Role flip — there's no BestLine entry to update.
func promoteRandomHandCardToArsenal(best *TurnSummary, startingHand []Card, arsenalCardIn Card) {
	if len(best.State.Hand) == 0 {
		return
	}
	pick := int(arsenalPromotionHash(startingHand, best.State.Hand, arsenalCardIn) % uint64(len(best.State.Hand)))
	chosen := best.State.Hand[pick]
	best.State.Arsenal = chosen
	best.State.Hand = append(best.State.Hand[:pick:pick], best.State.Hand[pick+1:]...)
	// Flip the matching BestLine entry from Held to Arsenal so per-card displays show the
	// correct role. Match the first Held entry whose card ID equals chosen — covers tutored
	// cards too if they happen to share an ID with a Held hand card, but harmlessly no-ops
	// when the chosen card is purely a tutored printing.
	for i := range best.BestLine {
		if best.BestLine[i].Role == Held && best.BestLine[i].Card.ID() == chosen.ID() {
			best.BestLine[i].Role = Arsenal
			break
		}
	}
}

// arsenalPromotionHash computes the deterministic bucket seed that picks which hand card
// fills an empty arsenal slot. FNV-1a over the starting-hand IDs + state-Hand IDs + arsenal-
// in ID — the only requirement is a uniform spread across bucket counts so the same hand
// always picks the same slot.
func arsenalPromotionHash(startingHand, stateHand []Card, arsenalCardIn Card) uint64 {
	const (
		fnvOffsetBasis uint64 = 1469598103934665603
		fnvPrime       uint64 = 1099511628211
	)
	h := fnvOffsetBasis
	for _, c := range startingHand {
		h ^= uint64(c.ID())
		h *= fnvPrime
	}
	for _, c := range stateHand {
		h ^= uint64(c.ID())
		h *= fnvPrime
	}
	if arsenalCardIn != nil {
		h ^= uint64(arsenalCardIn.ID())
		h *= fnvPrime
	}
	return h
}

// groupByRoleInto appends hand cards into caller-provided pitched/attackers/defenders slices
// (passed pre-reset to length 0) to avoid per-partition heap allocation.
func groupByRoleInto(hand []Card, roles []Role, pitched, attackers, defenders []Card) ([]Card, []Card, []Card) {
	for i, c := range hand {
		switch roles[i] {
		case Pitch:
			pitched = append(pitched, c)
		case Attack:
			attackers = append(attackers, c)
		case Defend:
			defenders = append(defenders, c)
		}
	}
	return pitched, attackers, defenders
}

// gatherHeldCards appends every hand card with role Held into the caller-provided held slice
// (passed pre-reset to length 0) and returns it. Threads the partition's Held set into
// bestAttackWithWeapons so alt-cost effects can consult it via TurnState.Hand.
func gatherHeldCards(hand []Card, roles []Role, held []Card) []Card {
	for i, c := range hand {
		if roles[i] == Held {
			held = append(held, c)
		}
	}
	return held
}

// findArsenalCard returns the arsenal-in card when it stays in the arsenal slot, nil otherwise.
// Hand cards never take Arsenal role during enumeration (post-hoc promotion handles that), so
// the only slot that can be Arsenal is the arsenal-in slot at index n.
func findArsenalCard(rolesBuf []Role, arsenalCardIn Card, n int) Card {
	if arsenalCardIn != nil && rolesBuf[n] == Arsenal {
		return arsenalCardIn
	}
	return nil
}

// roleAllowed decides whether the partition enumerator may assign role r to the current card.
// The arsenal-in slot may only take Arsenal (stay), Attack (an Action / Weapon — non-attack
// actions play fine from arsenal on your turn), or Defend (Defense Reactions only — plain-
// blocking from arsenal isn't legal). Hand cards take any role except Attack for cards that
// can't take Attack role at all (DRs and Block-typed cards lack the Action / Weapon subtype);
// their role loop caps at Held, so the "which Held card gets arsenaled" choice happens
// post-hoc and doesn't bias toward low-ID slots.
func roleAllowed(r Role, isArsenalSlot, isDefenseReaction, canAttack bool) bool {
	if isArsenalSlot {
		switch r {
		case Pitch, Held:
			return false
		case Attack:
			return canAttack
		case Defend:
			return isDefenseReaction
		}
		return true // Arsenal is always allowed on the arsenal-in slot.
	}
	return r != Attack || canAttack
}

// defendersDamage tallies the total Value contribution of the partition's defense phase. DRs
// resolve first via Play (their ApplyAndLogEffectiveDefense decrements state.IncomingDamage
// and credits the block, with arcane / runechant riders adding their own Value); plain blocks
// then consume whatever incoming damage is left, capped per card. Played in isolation — no
// attack ordering; per-DR TurnState carries Pitched / deck plus a fresh copy of the defenders
// list in graveyard so DRs that scan for banish targets see the same shape across iterations.
//
// blockBudget is the remaining defense-phase pitch supply after the caller has subtracted DR
// costs. Modal blockers (Blocker + ModalCard + BlockCost) enumerate their modes within
// blockBudget and pick the one yielding the highest BonusDefense; non-modal Blockers run
// their hook unchanged. Pass MaxInt to disable the budget check entirely (no modal blockers
// in the partition).
//
// arsenalDefenderIdx is the position of the arsenal-in card in defenders when it took the
// Defend role (-1 otherwise) — used to flag the matching CardState.FromArsenal so
// EffectiveDefense picks up the ArsenalDefenseBonus rider.
//
// state is caller-provided (from attackBufs) and reset per DR. gravBuf is the caller-owned
// scratch backing state's graveyard; the returned slice is the (possibly grown) buffer for
// reuse. The threaded incoming-damage counter persists across the DR loop and into the plain-
// block loop so DRs see the full incoming pool first (maximising any +1{d} riders) and plain
// blocks pick up the residual.
//
// Returns the per-DR cacheable status as a sticky bit — once a DR reads deck or graveyard
// via the accessors, the partition's defense-phase output isn't safe to cache; aggregated up
// through bestAttackWithWeapons.
func defendersDamage(defenders, pitched, deck []Card, state *TurnState, gravBuf []Card, cs *CardState, incomingDamage, blockBudget, arsenalDefenderIdx int) (int, []Card, bool) {
	total := 0
	remaining := incomingDamage
	cacheable := true
	// state.Auras is caller-seeded with priorAuras so DR Plays consolidate created
	// auras against the carryover entries. The per-DR reset below preserves the
	// running list across iterations.
	for i, d := range defenders {
		if !attackerMetaPtrFor(d).actsAsDR {
			continue
		}
		gravBuf = append(gravBuf[:0], defenders...)
		// Per-DR seed starts cacheable; the DR's Play flips it via accessors if it reads
		// deck or graveyard. Set explicitly because TurnState's zero-value is uncacheable.
		// Auras carry across the per-DR reset so created auras persist for the caller.
		preservedAuras := state.Auras
		*state = TurnState{Pitched: pitched, deck: deck, graveyard: gravBuf, IncomingDamage: remaining, cacheable: true, Defenders: defenders, Auras: preservedAuras}
		*cs = CardState{Card: d, FromArsenal: i == arsenalDefenderIdx}
		d.Play(state, cs)
		total += state.Value
		remaining = state.IncomingDamage
		if !state.IsCacheable() {
			cacheable = false
		}
	}
	// Plain blocks contribute their printed Defense plus any BonusDefense the optional
	// Blocker hook flipped after scanning state.Defenders. Modal blockers with BlockCost
	// pick the highest-bonus mode that fits the remaining budget; non-modal Blockers run
	// their hook unchanged. Reuse the caller-provided cs scratch so the interface call
	// doesn't escape a fresh CardState per plain blocker.
	state.Defenders = defenders
	for _, d := range defenders {
		if attackerMetaPtrFor(d).actsAsDR {
			continue
		}
		bestMode, bestCost := pickBlockerMode(d, state, cs, blockBudget)
		blockBudget -= bestCost
		*cs = CardState{Card: d, Mode: bestMode}
		if b, ok := d.(Blocker); ok {
			b.Block(state, cs)
		}
		block := cs.EffectiveDefense()
		if block > remaining {
			block = remaining
		}
		if block > 0 {
			total += block
			remaining -= block
		}
	}
	return total, gravBuf, cacheable
}

// pickBlockerMode returns the mode index and resource cost yielding the highest BonusDefense
// for d within blockBudget. Non-modal blockers and blockers without BlockCost return (0, 0)
// — the chain runner runs Block once at mode 0. Modal+BlockCost cards probe each mode by
// running Block on the cs scratch with the candidate mode, observing the resulting
// BonusDefense; the caller re-runs Block with the chosen mode for the actual contribution.
func pickBlockerMode(d Card, state *TurnState, cs *CardState, blockBudget int) (int8, int) {
	mc, ok := d.(ModalCard)
	if !ok {
		return 0, 0
	}
	bc, ok := d.(BlockCost)
	if !ok {
		return 0, 0
	}
	b, ok := d.(Blocker)
	if !ok {
		return 0, 0
	}
	bestBonus := -1
	bestMode := int8(0)
	bestCost := 0
	for mode := int8(0); mode < int8(mc.Modes()); mode++ {
		cost := bc.BlockCost(mode)
		if cost > blockBudget {
			continue
		}
		*cs = CardState{Card: d, Mode: mode}
		b.Block(state, cs)
		if cs.BonusDefense > bestBonus {
			bestBonus = cs.BonusDefense
			bestMode = mode
			bestCost = cost
		}
	}
	return bestMode, bestCost
}

// chainBudget captures the winning phase-split's attack-chain resource state.
type chainBudget struct {
	resource         int
	maxPitch         int
	hasAttackPitches bool
}

// phaseBudgets is one (pmask) configuration's split of pitched-resource totals across the
// attack and defense phases. Each side tracks both its running total and the largest single
// pitch assigned to it — the "largest pitch" feeds the pitch-timing waste check (if the
// residual budget after paying all costs is at least that value, one pitch could have been
// Held, and the partition is illegal).
type phaseBudgets struct {
	attackBudget, defendBudget         int
	maxAttackPitch, maxDefendPitch     int
	hasAttackPitches, hasDefendPitches bool
}

// splitPitchesAcrossPhases assigns each pitch to the attack or defense phase based on the
// bitmask and computes the per-phase resource summary. Bit i set → pitchedVals[i] funds
// defense; bit i clear → it funds attack. phaseCount==1 forces every pitch to the attack
// phase (no DRs present or no pitches to split) regardless of pmask.
func splitPitchesAcrossPhases(pitchedVals []int, pmask, phaseCount int) phaseBudgets {
	var p phaseBudgets
	for i, v := range pitchedVals {
		if phaseCount > 1 && pmask&(1<<i) != 0 {
			p.defendBudget += v
			if v > p.maxDefendPitch {
				p.maxDefendPitch = v
			}
			p.hasDefendPitches = true
		} else {
			p.attackBudget += v
			if v > p.maxAttackPitch {
				p.maxAttackPitch = v
			}
			p.hasAttackPitches = true
		}
	}
	return p
}

// containsDefenseReaction reports whether any card in cards participates in the defense
// phase via the Play hook (printed Defense Reactions or DefensiveInstant-marked Instants).
// The partition-leaf precompute uses this to decide whether the defense-phase pitch
// enumeration needs to split budgets at all — no such cards means every pitch funds the
// attack phase.
func containsDefenseReaction(cards []Card) bool {
	for _, c := range cards {
		if attackerMetaPtrFor(c).actsAsDR {
			return true
		}
	}
	return false
}

// containsModalBlocker reports whether any card in cards is a modal blocker — Blocker
// implementations whose mode count varies via ModalCard and whose mode cost varies via
// BlockCost. Partitions with at least one modal blocker recompute defendersDamage per
// (pmask, wmask) so each candidate's spare defense budget gates the mode pick; partitions
// without one keep the once-per-leaf defendersDamage shortcut. Defenders are short
// (≤ 4 cards in practice), so the inline type-assertion chain stays cheap and fast-fails
// on the BlockCost gate for the common no-modal-blocker case.
func containsModalBlocker(cards []Card) bool {
	for _, c := range cards {
		if _, ok := c.(BlockCost); !ok {
			continue
		}
		if mc, ok := c.(ModalCard); ok && mc.Modes() > 1 {
			if _, ok := c.(Blocker); ok {
				return true
			}
		}
	}
	return false
}

// noBlockBudgetCap is the sentinel passed to defendersDamage when the partition has no
// modal blockers — the per-mode cost path is unused, so any large value works. Picked to
// be obviously beyond any real defense budget without depending on math.MaxInt overflow
// arithmetic the budget arithmetic might do.
const noBlockBudgetCap = 1 << 30

// pendingFutureValue sums the Count of every non-token Aura plus every Item at end of
// chain — the partition tiebreaker's "hidden later-turn payoff" signal. Token auras
// (Runechant) credit Value at creation, so they're dropped. Card auras carry Count ==
// fires-remaining for counters and == 1 for one-shots. Items only credit Value when the
// activated ability is spent (drawing a card, consuming a Gold), so unspent items at end
// of chain represent payoff the partition is saving for later — counted here so a
// partition that leaves a Gold token in play beats one that arsenaled the creator and
// made nothing.
func pendingFutureValue(auras []Aura, items []Item) int {
	total := 0
	for _, a := range auras {
		if a.Self.IsToken() {
			continue
		}
		total += a.Count
	}
	for _, it := range items {
		total += it.Count
	}
	return total
}
