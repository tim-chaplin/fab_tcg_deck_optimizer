package sim

// Top-level hand enumeration: findBest walks every partition (Pitch / Attack / Defend /
// Held / Arsenal assignment) and delegates each leaf's chain-feasibility check to
// bestAttackWithWeapons. Post-enumeration helpers decide how an empty arsenal slot gets
// filled, plus the roleAllowed policy function that shapes the partition tree.

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
	role        deck.Role
	pitchVal    int
	defenseVal  int
	isDR        bool
	canAttack   bool
	fromArsenal bool
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
			e.cache.hits.Add(1)
			return e.replayBest(entry, weapons, hand, d, masterState)
		}
		e.cache.misses.Add(1)
	}

	bufs := e.getAttackBufs(len(hand), weapons)
	arsenalCardIn := masterState.Arsenal()
	incoming := masterState.IncomingDamage()
	n := len(hand)
	totalN := n
	if arsenalCardIn != nil {
		totalN = n + 1
	}

	best := TurnSummary{
		BestLine:       make([]deck.CardAssignment, totalN),
		IncomingDamage: incoming,
		Cacheable:      true,
	}
	for i := 0; i < n; i++ {
		best.BestLine[i] = deck.CardAssignment{Card: hand[i], Role: deck.Held}
	}
	if arsenalCardIn != nil {
		best.BestLine[n] = deck.CardAssignment{Card: arsenalCardIn, Role: deck.Arsenal, FromArsenal: true}
	}
	cacheable := true
	var bestSwung []string
	var runningSeen bool
	var runningScore chainScore

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
			score := chainScoreOf(winner, v)
			if runningSeen && score.cmp(runningScore) <= 0 {
				return
			}
			winner.SetArsenal(arsenalAtChainStart)
			best.State = winner
			best.Value = v
			runningScore = score
			runningSeen = true
			bestSwung = swung
			for j := 0; j < totalN; j++ {
				best.BestLine[j].Role = pcards[j].role
			}
			return
		}
		pc := &pcards[i]
		maxRole := deck.Held
		if pc.fromArsenal {
			maxRole = deck.Arsenal
		}
		for r := deck.Role(0); r <= maxRole; r++ {
			if !roleAllowed(r, pc.fromArsenal, pc.isDR, pc.canAttack) {
				continue
			}
			if r == deck.Defend && incoming == 0 {
				continue
			}
			pc.role = r
			switch r {
			case deck.Pitch:
				recurse(i+1, pitchSum+pc.pitchVal, defenseSum)
			case deck.Defend:
				recurse(i+1, pitchSum, defenseSum+pc.defenseVal)
			case deck.Attack, deck.Held, deck.Arsenal:
				recurse(i+1, pitchSum, defenseSum)
			}
		}
	}
	recurse(0, 0, 0)
	best.SwungWeapons = bestSwung
	best.Cacheable = cacheable
	if best.State == nil {
		// No feasible partition was found — synthesise an "untouched" trailing state
		// that holds the starting hand and prior values, so callers can read end-of-turn
		// fields off a non-nil State. Carry masterState.Value() onto best.Value so the
		// no-chain path still credits any start-of-action-phase aura tick.
		fallback := masterState.Copy()
		fallback.SetHand(append([]card.Card(nil), hand...))
		fallback.SetArsenal(arsenalCardIn)
		best.State = fallback
		best.Value = masterState.Value()
	}
	if best.State.Arsenal() == nil {
		promoteRandomHandCardToArsenal(&best, hand, arsenalCardIn)
	}
	if cacheUsable {
		if best.Cacheable {
			e.cache.store(cacheKey, evalCacheEntry{
				line:         append([]deck.CardAssignment(nil), best.BestLine...),
				swungWeapons: append([]string(nil), best.SwungWeapons...),
			})
		} else {
			e.cache.uncacheable.Add(1)
		}
	}
	return best
}

// promoteRandomHandCardToArsenal picks one card from best.State.Hand() and moves it
// into best.State's arsenal slot, removing it from the state's hand. Deterministic
// per-hand pick.
func promoteRandomHandCardToArsenal(best *TurnSummary, startingHand []card.Card, arsenalCardIn card.Card) {
	handState := best.State.Hand()
	if len(handState) == 0 {
		return
	}
	eligible := make([]int, 0, len(handState))
	for i, c := range handState {
		t := c.Types(nil)
		if t.Has(card.TypeBlock) || t.IsResource() {
			continue
		}
		eligible = append(eligible, i)
	}
	if len(eligible) == 0 {
		return
	}
	pick := eligible[int(arsenalPromotionHash(startingHand, handState, arsenalCardIn)%uint64(len(eligible)))]
	chosen := handState[pick]
	best.State.SetArsenal(chosen)
	newHand := append(handState[:pick:pick], handState[pick+1:]...)
	best.State.SetHand(newHand)
	for i := range best.BestLine {
		if best.BestLine[i].Role == deck.Held && best.BestLine[i].Card.ID() == chosen.ID() {
			best.BestLine[i].Role = deck.Arsenal
			break
		}
	}
}

// arsenalPromotionHash computes the deterministic bucket seed via FNV-1a.
func arsenalPromotionHash(startingHand, stateHand []card.Card, arsenalCardIn card.Card) uint64 {
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

// groupByRole appends each card to the caller-provided pitched / attackers / defenders
// slice matching its enumerated role; Held and Arsenal cards are skipped.
func groupByRole(pcards []partitionCard, pitched, attackers, defenders []card.Card) ([]card.Card, []card.Card, []card.Card) {
	for _, pc := range pcards {
		switch pc.role {
		case deck.Pitch:
			pitched = append(pitched, pc.card)
		case deck.Attack:
			attackers = append(attackers, pc.card)
		case deck.Defend:
			defenders = append(defenders, pc.card)
		}
	}
	return pitched, attackers, defenders
}

// gatherHeldCards appends every card with role Held into the caller-provided held slice.
func gatherHeldCards(pcards []partitionCard, held []card.Card) []card.Card {
	for _, pc := range pcards {
		if pc.role == deck.Held {
			held = append(held, pc.card)
		}
	}
	return held
}

// findArsenalCard returns the arsenal-in card when its slot kept the Arsenal role, nil
// otherwise (it was reassigned to Attack / Defend, or no card started in the arsenal).
func findArsenalCard(pcards []partitionCard, n int) card.Card {
	if len(pcards) > n && pcards[n].role == deck.Arsenal {
		return pcards[n].card
	}
	return nil
}

// roleAllowed decides whether the partition enumerator may assign role r to the current
// card.
func roleAllowed(r deck.Role, isArsenalSlot, isDefenseReaction, canAttack bool) bool {
	if isArsenalSlot {
		switch r {
		case deck.Pitch, deck.Held:
			return false
		case deck.Attack:
			return canAttack
		case deck.Defend:
			return isDefenseReaction
		}
		return true
	}
	return r != deck.Attack || canAttack
}

// defendersDamage tallies the total Value contribution of the partition's defense phase
// against the caller-supplied state engine. DRs resolve first; plain blocks then consume
// whatever incoming damage is left, capped per card. The engine is mutated in place: the
// matchup figure rides in via SetIncomingDamage (which zeroes the damage-blocked
// accumulator), each DR's resolution and each plain block bank into that accumulator, and
// the chain phase reads the post-defense graveyard via the engine's left-behind state.
//
// blockBudget is the remaining defense-phase pitch supply after the caller has subtracted
// DR costs. Modal blockers enumerate their modes within blockBudget and pick the one
// yielding the highest BonusDefense; non-modal Blockers run their hook unchanged.
//
// Returns the per-DR cacheable status as a sticky bit — once a DR reads deck or graveyard,
// the partition's defense-phase output isn't safe to cache.
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

// chainScore is a leaf's comparable score. cmp ranks it lexicographically:
//   - value: the chain-step damage / block credited this turn.
//   - cardsPlayed: cards played this turn. A played card does something useful even when
//     the payoff lands next turn (an aura ticks later, a token mints currency), so
//     playing > holding when value ties.
//   - totalCards: cards available next turn — the post-refill hand (held cards topped up
//     to intellect) plus an occupied arsenal slot.
//   - totalCounters: summed Count of every Aura plus every Item — pending aura fires +
//     token stockpile, the weakest signal.
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
