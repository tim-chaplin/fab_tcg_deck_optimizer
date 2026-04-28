package sim

// Top-level hand enumeration: findBest walks every partition (Pitch / Attack / Defend /
// Held / Arsenal assignment) and delegates each leaf's chain-feasibility check to
// bestAttackWithWeapons. Post-enumeration helpers decide how an empty arsenal slot gets
// filled, plus the beatsBest / roleAllowed policy functions that shape the partition tree.

import ()

func (e *Evaluator) findBest(hero Hero, weapons []Weapon, hand []Card, incomingDamage int, deck []Card, runechantCarryover int, arsenalCardIn Card, priorAuraTriggers []AuraTrigger, skipLog bool) TurnSummary {
	// Cache fast-path. priorAuraTriggers != 0 disables the cache entirely — carryover
	// triggers add a hidden state input the key doesn't model, and serializing
	// AuraTrigger.Handler closures isn't worth the complexity. Cache disabled (e.cache nil)
	// or hand too big to fingerprint also bypass.
	var cacheKey evalCacheKey
	cacheUsable := e.cache != nil && len(priorAuraTriggers) == 0
	if cacheUsable {
		var keyOK bool
		cacheKey, keyOK = makeCacheKey(hero, weapons, hand, runechantCarryover, arsenalCardIn)
		if !keyOK {
			cacheUsable = false
		}
	}
	if cacheUsable {
		if entry, ok := e.cache.lookup(cacheKey); ok {
			e.cache.hits++
			return e.replayBest(entry, hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn, skipLog)
		}
		e.cache.misses++
	} else if e.cache != nil && len(priorAuraTriggers) > 0 {
		e.cache.skipsTriggers++
	}

	n := len(hand)
	// The partition recurse treats the arsenal-in card as an extra entry at index n with a
	// restricted role menu (Arsenal / Attack / Defend), so everything about it is decided inside
	// the enumeration. totalN is the effective size of BestLine.
	totalN := n
	if arsenalCardIn != nil {
		totalN = n + 1
	}

	// Seed best.State.Runechants with the carryover: partitions with no attacks don't reduce
	// it, so carryover is the baseline to beat. BestLine starts with every hand card Held and
	// the arsenal-in card (if any) staying in the slot, so a hand with no Value-adding
	// partition still reports sensible "nothing played, nothing pitched" assignments.
	// Cacheable starts true: the no-feasible-line fallback ran no chain and read no hidden
	// state, so the seed is trivially cacheable. Each leaf's cacheable bit ANDs into a
	// running sticky and the final value stamps best.Cacheable.
	best := TurnSummary{
		BestLine:       make([]CardAssignment, totalN),
		IncomingDamage: incomingDamage,
		Cacheable:      true,
		State: CarryState{
			Hand:         append([]Card(nil), hand...),
			Deck:         append([]Card(nil), deck...),
			Arsenal:      arsenalCardIn,
			Runechants:   runechantCarryover,
			AuraTriggers: append([]AuraTrigger(nil), priorAuraTriggers...),
		},
	}
	cacheable := true
	// bestSwung holds the winning partition's swung weapon names — surfaced on the summary so
	// the printout can list weapons that swung this turn (weapons have no BestLine entry).
	var bestSwung []string
	// bestHasHeld tracks whether the current best has at least one Held hand card — lets
	// beatsBest distinguish "arsenal will be occupied post-hoc" from "arsenal will be empty."
	// Seeded true when the hand is non-empty: the initial best puts every hand card into Held,
	// so a post-hoc promotion would fill arsenal. Candidates need both a Value/leftover tie and
	// some way to end with arsenal occupied to displace it.
	bestHasHeld := n > 0
	// bestFutureValuePlayed tracks how many AddsFutureValue cards the current best is
	// playing (Role=Attack). Seeded 0 because the initial best assigns every card Held. The
	// beatsBest tiebreaker prefers partitions that play MORE future-value cards at equal
	// Value/leftover — their hidden-later-turn payoff is invisible to the current-turn
	// score, so without this bias a lone sigil loses to Held → arsenal promotion on the
	// arsenal-occupancy tiebreak.
	bestFutureValuePlayed := 0
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
	rolesBuf := bufs.rolesBuf[:totalN]
	pvals := bufs.pitchVals[:totalN]
	dvals := bufs.defenseVals[:totalN]
	isDR := bufs.isDRBuf[:totalN]
	addsFutureValue := bufs.addsFutureValueBuf[:totalN]

	hasReactions := fillPartitionPerCardBufs(hand, n, totalN, arsenalCardIn, pvals, dvals, isDR, addsFutureValue)

	var recurse func(i, pitchSum, defenseSum int)
	recurse = func(i, pitchSum, defenseSum int) {
		if i == totalN {
			attackDealt, defenseDealt, leftoverRunechants, swung, carry, ok, leafCacheable, arsenalAtChainStart := e.evaluatePartition(
				hero, weapons, hand, deck, arsenalCardIn,
				rolesBuf, n, bufs,
				runechantCarryover, incomingDamage, defenseSum,
				priorAuraTriggers, skipLog,
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
			// Hand cards never take Arsenal role during enumeration, so arsenalCard is only set
			// when arsenal-in stayed; post-hoc promotion potential is tracked via hasHeld.
			hasHeld := false
			futureValuePlayed := 0
			for j := 0; j < n; j++ {
				if rolesBuf[j] == Held {
					hasHeld = true
				}
				if rolesBuf[j] == Attack && addsFutureValue[j] {
					futureValuePlayed++
				}
			}
			if arsenalCardIn != nil && rolesBuf[n] == Attack && addsFutureValue[n] {
				futureValuePlayed++
			}
			willOccupy := arsenalCard != nil || hasHeld
			bestWillOccupy := best.State.Arsenal != nil || bestHasHeld
			if !beatsBest(v, leftoverRunechants, futureValuePlayed, willOccupy, best, bestFutureValuePlayed, bestWillOccupy) {
				return
			}
			best.Value = v
			bestSwung = swung
			// Adopt the winner's CarryState wholesale; arsenal-in occupancy overrides the
			// snapshot's Arsenal so an arsenal-in card that stayed is preserved.
			best.State = carry
			best.State.Arsenal = arsenalCard
			bestHasHeld = hasHeld
			bestFutureValuePlayed = futureValuePlayed
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
			if !roleAllowed(r, isArsenalSlot, isDR[i]) {
				continue
			}
			// With no damage coming in and no Defense Reactions in the hand, a non-DR card's
			// Defend contribution is 0 — same as Held — and nothing scans the defender set,
			// so the two partitions produce the same Value / leftover / futureValuePlayed and
			// Held wins the arsenal-occupancy tiebreaker. Skip the dominated Defend branch.
			// DR-present hands keep Defend because DR Play effects scan defenders as a
			// graveyard seed (e.g. Weeping Battleground banishing an aura a non-DR blocker
			// put there).
			if r == Defend && incomingDamage == 0 && !isDR[i] && !hasReactions {
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
			e.cache.uncacheable++
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

// beatsBest decides whether a candidate partition displaces the current best. Tiebreak
// order: Value → leftover runechants (future arcane) → more AddsFutureValue cards played
// (hidden later-turn payoff the current-turn Value misses) → arsenal slot ending occupied
// (saves a hand slot next refill; covers both arsenal-in-stayed and Held-for-promotion).
func beatsBest(v, leftoverRunechants, futureValuePlayed int, willOccupyArsenal bool, best TurnSummary, bestFutureValuePlayed int, bestWillOccupyArsenal bool) bool {
	if v > best.Value {
		return true
	}
	if v < best.Value {
		return false
	}
	if leftoverRunechants > best.State.Runechants {
		return true
	}
	if leftoverRunechants < best.State.Runechants {
		return false
	}
	if futureValuePlayed > bestFutureValuePlayed {
		return true
	}
	if futureValuePlayed < bestFutureValuePlayed {
		return false
	}
	return willOccupyArsenal && !bestWillOccupyArsenal
}

// roleAllowed decides whether the partition enumerator may assign role r to the current card.
// The arsenal-in slot may only take Arsenal (stay), Attack (any non-DR card — non-attack actions
// play fine from arsenal on your turn), or Defend (Defense Reactions only — plain-blocking from
// arsenal isn't legal). Hand cards take any role except Attack for Defense Reactions (DRs only
// fire on the opponent's turn); their role loop caps at Held, so the "which Held card gets
// arsenaled" choice happens post-hoc and doesn't bias toward low-ID slots.
func roleAllowed(r Role, isArsenalSlot, isDefenseReaction bool) bool {
	if isArsenalSlot {
		switch r {
		case Pitch, Held:
			return false
		case Attack:
			return !isDefenseReaction
		case Defend:
			return isDefenseReaction
		}
		return true // Arsenal is always allowed on the arsenal-in slot.
	}
	return !(r == Attack && isDefenseReaction)
}

// defendersDamage tallies the total Value contribution of the partition's defense phase. DRs
// resolve first via Play (their ApplyAndLogEffectiveDefense decrements state.IncomingDamage
// and credits the block, with arcane / runechant riders adding their own Value); plain blocks
// then consume whatever incoming damage is left, capped per card. Played in isolation — no
// attack ordering; per-DR TurnState carries Pitched / deck plus a fresh copy of the defenders
// list in graveyard so DRs that scan for banish targets see the same shape across iterations.
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
func defendersDamage(defenders, pitched, deck []Card, state *TurnState, gravBuf []Card, cs *CardState, incomingDamage, arsenalDefenderIdx int) (int, []Card, bool) {
	total := 0
	remaining := incomingDamage
	cacheable := true
	for i, d := range defenders {
		if !d.Types().IsDefenseReaction() {
			continue
		}
		gravBuf = append(gravBuf[:0], defenders...)
		// Per-DR seed starts cacheable; the DR's Play flips it via accessors if it reads
		// deck or graveyard. Set explicitly because TurnState's zero-value is uncacheable.
		*state = TurnState{Pitched: pitched, deck: deck, graveyard: gravBuf, IncomingDamage: remaining, cacheable: true}
		*cs = CardState{Card: d, FromArsenal: i == arsenalDefenderIdx}
		d.Play(state, cs)
		total += state.Value
		remaining = state.IncomingDamage
		if !state.IsCacheable() {
			cacheable = false
		}
	}
	for _, d := range defenders {
		if d.Types().IsDefenseReaction() {
			continue
		}
		block := d.Defense()
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

// containsDefenseReaction reports whether any card in cards is a Defense Reaction. The
// partition-leaf precompute uses this to decide whether the defense-phase pitch enumeration
// needs to split budgets at all (no DRs means every pitch funds the attack phase).
func containsDefenseReaction(cards []Card) bool {
	for _, c := range cards {
		if c.Types().IsDefenseReaction() {
			return true
		}
	}
	return false
}

