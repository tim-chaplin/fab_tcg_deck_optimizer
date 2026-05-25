package sim

// Top-level hand enumeration: findBest walks every Pitch / Attack / Defend / Held / Arsenal
// assignment and delegates each leaf to bestAttackWithWeapons. Plus roleAllowed (the
// partition-tree shape policy) and the empty-arsenal-promotion helpers.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// partitionCard is one slot in the partition enumerator's working set: a hand (or
// arsenal-in) card, the role the recurse is currently assigning it, and the static
// per-card facts that roleAllowed and the budget sums read.
type partitionCard struct {
	card        card.Card
	role        card.Role
	pitchVal    int
	defenseVal  int
	isDR        bool
	canAttack   bool
	fromArsenal bool
}

// resourceBonusUpperBound sums MaxResourcePoints over every card.ResourceSource among the
// turn's cards — hand, the arsenal-in card, and in-play items' source cards.
func resourceBonusUpperBound(hand []card.Card, arsenalIn card.Card, items []gameengine.Item) int {
	sum := 0
	add := func(c card.Card) {
		if rs, ok := c.(card.ResourceSource); ok {
			sum += rs.MaxResourcePoints()
		}
	}
	for _, c := range hand {
		add(c)
	}
	if arsenalIn != nil {
		add(arsenalIn)
	}
	for _, it := range items {
		if sc, ok := it.SourceCard().(card.Card); ok {
			add(sc)
		}
	}
	return sum
}

func (e *Evaluator) findBest(weapons []weapon.Weapon, hand []card.Card, d *deck.Deck, masterState *gameengine.GameState) TurnSummary {
	var cacheKey evalCacheKey
	cacheUsable := e.cache != nil
	if cacheUsable {
		var keyOK bool
		cacheKey, keyOK = makeCacheKey(weapons, hand, masterState)
		if !keyOK {
			cacheUsable = false
		}
	}
	if cacheUsable {
		if entry, ok := e.cache.lookup(cacheKey); ok {
			// Replay only when the caller's deck is at least as deep as the chain consumed
			// during caching; a shallower deck would hit a deck-empty branch and diverge.
			if d == nil || d.Size() >= entry.cardsRemovedFromDeck {
				e.cache.hits.Add(1)
				return e.replayBest(entry, weapons, hand, d, masterState)
			}
		}
		e.cache.misses.Add(1)
	}

	bufs := e.getAttackBufs(len(hand), weapons)
	// Cleared so a Best with no feasible partition stores an empty solution, not a stale one.
	bufs.bestSolution.reset()
	arsenalCardIn := masterState.Arsenal()
	bufs.maxResourceBonus = resourceBonusUpperBound(hand, arsenalCardIn, masterState.Items())
	incoming := masterState.IncomingDamage()
	n := len(hand)
	totalN := n
	if arsenalCardIn != nil {
		totalN = n + 1
	}

	best := TurnSummary{
		BestLine:       make([]card.CardAssignment, totalN),
		IncomingDamage: incoming,
		Cacheable:      true,
	}
	for i := 0; i < n; i++ {
		best.BestLine[i] = card.CardAssignment{Card: hand[i], Role: card.Held}
	}
	if arsenalCardIn != nil {
		best.BestLine[n] = card.CardAssignment{Card: arsenalCardIn, Role: card.Arsenal, FromArsenal: true}
	}
	cacheable := true
	var bestSwung []string
	var runningSeen bool
	var runningScore chainScore

	// Defend role is valid whenever any incoming damage — physical or arcane — could be
	// prevented.
	anyIncoming := incoming > 0 || masterState.ArcaneIncomingDamage() > 0

	pcards := bufs.partitionCards[:totalN]
	fillPartitionCards(hand, n, totalN, arsenalCardIn, pcards)

	var recurse func(i, pitchSum, defenseSum int)
	recurse = func(i, pitchSum, defenseSum int) {
		if i == totalN {
			attackDealt, defenseDealt, swung, winner, ok, leafCacheable, arsenalAtChainStart := e.evaluatePartition(
				masterState, weapons, d,
				pcards, n, bufs,
				defenseSum,
			)
			if !leafCacheable {
				cacheable = false
			}
			if !ok {
				return
			}

			v := attackDealt + defenseDealt
			winner.SetArsenal(arsenalAtChainStart)
			var promoted card.Card
			if winner.Arsenal() == nil {
				promoted = promoteHeldToArsenal(winner, hand, arsenalCardIn)
			}
			score := chainScoreOf(winner, v)
			if runningSeen && score.cmp(runningScore) <= 0 {
				return
			}
			best.State = winner
			best.Value = v
			runningScore = score
			runningSeen = true
			bestSwung = swung
			bufs.bestSolution.copyFrom(&bufs.partSolution)
			for j := 0; j < totalN; j++ {
				best.BestLine[j].Role = pcards[j].role
			}
			markPromotedInBestLine(best.BestLine, promoted)
			return
		}
		pc := &pcards[i]
		maxRole := card.Held
		if pc.fromArsenal {
			maxRole = card.Arsenal
		}
		for r := card.Role(0); r <= maxRole; r++ {
			if !roleAllowed(r, pc.fromArsenal, pc.isDR, pc.canAttack) {
				continue
			}
			if r == card.Defend && !anyIncoming {
				continue
			}
			pc.role = r
			switch r {
			case card.Pitch:
				recurse(i+1, pitchSum+pc.pitchVal, defenseSum)
			case card.Defend:
				recurse(i+1, pitchSum, defenseSum+pc.defenseVal)
			case card.Attack, card.Held, card.Arsenal:
				recurse(i+1, pitchSum, defenseSum)
			}
		}
	}
	recurse(0, 0, 0)
	best.SwungWeapons = bestSwung
	best.Cacheable = cacheable
	if best.State == nil {
		// No feasible partition — synthesise an untouched trailing state holding the starting
		// hand and prior values so callers can read end-of-turn fields off a non-nil State.
		// best.Value carries masterState.Value() to credit any start-of-action-phase aura tick.
		fallback := masterState.Copy()
		fallback.SetHand(append([]card.Card(nil), hand...))
		fallback.SetArsenal(arsenalCardIn)
		best.State = fallback
		best.Value = masterState.Value()
		if best.State.Arsenal() == nil {
			markPromotedInBestLine(best.BestLine, promoteHeldToArsenal(best.State, hand, arsenalCardIn))
		}
	}
	if cacheUsable {
		if best.Cacheable {
			cardsRemoved := 0
			if best.State != nil {
				cardsRemoved = best.State.CardsRemovedFromDeck()
			}
			e.cache.store(cacheKey, evalCacheEntry{
				line:                 append([]card.CardAssignment(nil), best.BestLine...),
				swungWeapons:         append([]string(nil), best.SwungWeapons...),
				attackOrder:          append([]playedCard(nil), bufs.bestSolution.attack...),
				pitchOrder:           append([]card.Card(nil), bufs.bestSolution.pitch...),
				defenders:            append([]playedCard(nil), bufs.bestSolution.defenders...),
				cardsRemovedFromDeck: cardsRemoved,
			})
		} else {
			e.cache.uncacheable.Add(1)
		}
	}
	return best
}

// promoteHeldToArsenal moves an arsenal-eligible Held card from state's hand into its empty
// arsenal slot and returns it, or nil when none are eligible. The pick is deterministic per
// hand. Callers tracking a best line mark the returned card via markPromotedInBestLine.
func promoteHeldToArsenal(state *gameengine.GameState, startingHand []card.Card, arsenalCardIn card.Card) card.Card {
	hand := state.HandStates()
	eligible := 0
	for i := range hand {
		if isArsenalEligible(hand[i].Card) {
			eligible++
		}
	}
	if eligible == 0 {
		return nil
	}
	target := int(arsenalPromotionHash(startingHand, hand, arsenalCardIn) % uint64(eligible))
	for i := range hand {
		if !isArsenalEligible(hand[i].Card) {
			continue
		}
		if target == 0 {
			chosen := hand[i].Card
			state.SetArsenal(chosen)
			// Swap-remove: the hand is read by membership / length, never by index.
			last := len(hand) - 1
			hand[i] = hand[last]
			state.SetHandStates(hand[:last])
			return chosen
		}
		target--
	}
	return nil
}

// isArsenalEligible reports whether c may fill an empty arsenal slot: any card that is
// neither a block nor a resource.
func isArsenalEligible(c card.Card) bool {
	t := c.Types(nil)
	return !t.Has(card.TypeBlock) && !t.IsResource()
}

// markPromotedInBestLine flips promoted's Held assignment to Arsenal in the best line.
func markPromotedInBestLine(line []card.CardAssignment, promoted card.Card) {
	if promoted == nil {
		return
	}
	for i := range line {
		if line[i].Role == card.Held && line[i].Card == promoted {
			line[i].Role = card.Arsenal
			return
		}
	}
}

// arsenalPromotionHash computes the deterministic bucket seed via FNV-1a.
func arsenalPromotionHash(startingHand []card.Card, stateHand []card.CardState, arsenalCardIn card.Card) uint64 {
	const (
		fnvOffsetBasis uint64 = 1469598103934665603
		fnvPrime       uint64 = 1099511628211
	)
	h := fnvOffsetBasis
	for _, c := range startingHand {
		h ^= uint64(c.ID())
		h *= fnvPrime
	}
	for _, hc := range stateHand {
		h ^= uint64(hc.Card.ID())
		h *= fnvPrime
	}
	if arsenalCardIn != nil {
		h ^= uint64(arsenalCardIn.ID())
		h *= fnvPrime
	}
	return h
}

// groupByRole appends each card to the caller-provided pitched / attackers / defenders
// slice matching its enumerated role; Held and Arsenal cards are skipped.
func groupByRole(pcards []partitionCard, pitched, attackers, defenders []card.Card) ([]card.Card, []card.Card, []card.Card) {
	for _, pc := range pcards {
		switch pc.role {
		case card.Pitch:
			pitched = append(pitched, pc.card)
		case card.Attack:
			attackers = append(attackers, pc.card)
		case card.Defend:
			defenders = append(defenders, pc.card)
		}
	}
	return pitched, attackers, defenders
}

// gatherHeldCards appends every card with role Held into the caller-provided held slice.
func gatherHeldCards(pcards []partitionCard, held []card.Card) []card.Card {
	for _, pc := range pcards {
		if pc.role == card.Held {
			held = append(held, pc.card)
		}
	}
	return held
}

// findArsenalCard returns the arsenal-in card when its slot kept the Arsenal role, nil
// otherwise (it was reassigned to Attack / Defend, or no card started in the arsenal).
func findArsenalCard(pcards []partitionCard, n int) card.Card {
	if len(pcards) > n && pcards[n].role == card.Arsenal {
		return pcards[n].card
	}
	return nil
}

// roleAllowed decides whether the partition enumerator may assign role r to the current
// card.
func roleAllowed(r card.Role, isArsenalSlot, isDefenseReaction, canAttack bool) bool {
	if isArsenalSlot {
		switch r {
		case card.Pitch, card.Held:
			return false
		case card.Attack:
			return canAttack
		case card.Defend:
			return isDefenseReaction
		}
		return true
	}
	return r != card.Attack || canAttack
}

// defendersDamage tallies the defense phase's total Value contribution against the supplied
// engine. DRs resolve first; plain blocks then consume whatever incoming damage is left,
// capped per card. blockBudget is the resource pool the modal-blocker mode pick spends from.
// Returns the sticky cacheable bit — once a DR reads deck or graveyard, the partition's
// defense output is uncacheable.
func defendersDamage(defenders, pitched []card.Card, deckPile *deck.Deck, ge *gameengine.GameEngine, gravBuf []card.Card, cs *card.CardState, incomingDamage, blockBudget, arsenalDefenderIdx int) (int, []card.Card, bool) {
	total := 0
	cacheable := true
	ge.SetDeck(deckPile)
	ge.SetIncomingDamage(incomingDamage)
	for i, def := range defenders {
		if !attackerMetaPtrFor(def).actsAsDR {
			continue
		}
		gravBuf = append(gravBuf[:0], defenders...)
		ge.SetGraveyard(gravBuf)
		ge.SetPitched(pitched)
		ge.SetDefenders(defenders)
		ge.SetValue(0)
		ge.SetCacheable(true)
		*cs = card.CardState{Card: def, FromArsenal: i == arsenalDefenderIdx}
		ge.ResolveChainStep(ge.Logger(), cs)
		total += ge.Value()
		if !ge.IsCacheable() {
			cacheable = false
		}
	}
	ge.SetDefenders(defenders)
	for _, def := range defenders {
		if attackerMetaPtrFor(def).actsAsDR {
			continue
		}
		bestMode, bestCost := pickBlockerMode(def, ge, cs, blockBudget)
		blockBudget -= bestCost
		*cs = card.CardState{Card: def, Mode: bestMode}
		if b, ok := def.(card.Blocker); ok {
			b.Block(ge, ge.Logger(), cs)
		}
		block := cs.EffectiveDefense()
		if rem := ge.RemainingUnblockedDamage(); block > rem {
			block = rem
		}
		if block > 0 {
			total += block
			ge.AddDamageBlocked(block)
		}
	}
	return total, gravBuf, cacheable
}

// pickBlockerMode returns the mode index and resource cost yielding the highest
// BonusDefense for d within blockBudget.
func pickBlockerMode(d card.Card, ge *gameengine.GameEngine, cs *card.CardState, blockBudget int) (int8, int) {
	mc, ok := d.(card.Modal)
	if !ok {
		return 0, 0
	}
	bc, ok := d.(card.BlockCost)
	if !ok {
		return 0, 0
	}
	b, ok := d.(card.Blocker)
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
		*cs = card.CardState{Card: d, Mode: mode}
		b.Block(ge, ge.Logger(), cs)
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
// attack and defense phases.
type phaseBudgets struct {
	attackBudget, defendBudget         int
	maxAttackPitch, maxDefendPitch     int
	hasAttackPitches, hasDefendPitches bool
}

// splitPitchesAcrossPhases assigns each pitch to the attack or defense phase based on the
// bitmask.
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
// phase via the Play hook.
func containsDefenseReaction(cards []card.Card) bool {
	for _, c := range cards {
		if attackerMetaPtrFor(c).actsAsDR {
			return true
		}
	}
	return false
}

// containsModalBlocker reports whether any card in cards is a modal blocker.
func containsModalBlocker(cards []card.Card) bool {
	for _, c := range cards {
		if _, ok := c.(card.BlockCost); !ok {
			continue
		}
		if mc, ok := c.(card.Modal); ok && mc.Modes() > 1 {
			if _, ok := c.(card.Blocker); ok {
				return true
			}
		}
	}
	return false
}

// noBlockBudgetCap is the sentinel passed to defendersDamage when the partition has no
// modal blockers.
const noBlockBudgetCap = 1 << 30

// chainScore ranks leaves lexicographically:
//   - value: chain-step damage / block credited this turn.
//   - cardsPlayed: tiebreaker favouring "play > hold" when payoff lands next turn (aura ticks,
//     minted tokens).
//   - totalCards: post-refill hand + arsenal slot available next turn.
//   - totalCounters: summed Aura + Item Count — the weakest signal.
type chainScore struct {
	value         int
	cardsPlayed   int
	totalCards    int
	totalCounters int
}

// cmp returns 1 when a outranks b, -1 when b outranks a, 0 when they tie.
func (a chainScore) cmp(b chainScore) int {
	for _, d := range [...]int{
		a.value - b.value,
		a.cardsPlayed - b.cardsPlayed,
		a.totalCards - b.totalCards,
		a.totalCounters - b.totalCounters,
	} {
		if d > 0 {
			return 1
		}
		if d < 0 {
			return -1
		}
	}
	return 0
}
