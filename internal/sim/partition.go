package sim

// Top-level hand enumeration: findBest walks every partition (Pitch / Attack / Defend /
// Held / Arsenal assignment) and delegates each leaf's chain-feasibility check to
// bestAttackWithWeapons. Post-enumeration helpers decide how an empty arsenal slot gets
// filled, plus the roleAllowed policy function that shapes the partition tree.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

func (e *Evaluator) findBest(weapons []Weapon, hand []card.Card, mp Matchup, d *deck.Deck, prior Prior, skipLog bool) TurnSummary {
	masterState := newTurnMasterState(prior, mp, d)

	var cacheKey evalCacheKey
	cacheUsable := e.cache != nil
	if cacheUsable {
		var keyOK bool
		cacheKey, keyOK = makeCacheKey(weapons, hand, prior)
		if !keyOK {
			cacheUsable = false
		}
	}
	if cacheUsable {
		if entry, ok := e.cache.lookup(cacheKey); ok {
			e.cache.hits.Add(1)
			return e.replayBest(entry, weapons, hand, mp, d, prior, masterState, skipLog)
		}
		e.cache.misses.Add(1)
	}

	bufs := e.getAttackBufs(len(hand), weapons)
	arsenalCardIn := prior.Arsenal
	n := len(hand)
	totalN := n
	if arsenalCardIn != nil {
		totalN = n + 1
	}

	best := TurnSummary{
		BestLine:       make([]CardAssignment, totalN),
		IncomingDamage: mp.IncomingDamage,
		Cacheable:      true,
	}
	for i := 0; i < n; i++ {
		best.BestLine[i] = CardAssignment{Card: hand[i], Role: Held}
	}
	if arsenalCardIn != nil {
		best.BestLine[n] = CardAssignment{Card: arsenalCardIn, Role: Arsenal, FromArsenal: true}
	}
	cacheable := true
	var bestSwung []string
	var runningSeen bool
	var runningFutureValue int

	rolesBuf := bufs.rolesBuf[:totalN]
	pvals := bufs.pitchVals[:totalN]
	dvals := bufs.defenseVals[:totalN]
	isDR := bufs.isDRBuf[:totalN]
	canAttack := bufs.canAttackBuf[:totalN]

	fillPartitionPerCardBufs(hand, n, totalN, arsenalCardIn, pvals, dvals, isDR, canAttack)

	var recurse func(i, pitchSum, defenseSum int)
	recurse = func(i, pitchSum, defenseSum int) {
		if i == totalN {
			attackDealt, defenseDealt, swung, winner, ok, leafCacheable, arsenalAtChainStart := e.evaluatePartition(
				masterState, weapons, hand, d,
				rolesBuf, n, bufs,
				mp, defenseSum,
				prior, skipLog,
			)
			if !leafCacheable {
				cacheable = false
			}
			if !ok {
				return
			}

			v := attackDealt + defenseDealt
			futureValuePlayed := pendingFutureValueFromState(winner)
			if runningSeen {
				cmp := chainScoreCmp(v, winner.CardsDrawn(), futureValuePlayed,
					best.Value, best.State.CardsDrawn(), runningFutureValue)
				if cmp < 0 {
					return
				}
				if cmp == 0 {
					willOccupy := arsenalAtChainStart != nil || len(winner.Hand()) > 0
					runningWillOccupy := best.State.Arsenal() != nil || len(best.State.Hand()) > 0
					if !willOccupy || runningWillOccupy {
						return
					}
				}
			}
			winner.SetArsenal(arsenalAtChainStart)
			best.State = winner
			best.Value = v
			runningFutureValue = futureValuePlayed
			runningSeen = true
			bestSwung = swung
			for j := 0; j < totalN; j++ {
				best.BestLine[j].Role = rolesBuf[j]
			}
			return
		}
		isArsenalSlot := i == n && arsenalCardIn != nil
		maxRole := Held
		if isArsenalSlot {
			maxRole = Arsenal
		}
		for r := Role(0); r <= maxRole; r++ {
			if !roleAllowed(r, isArsenalSlot, isDR[i], canAttack[i]) {
				continue
			}
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
	best.Cacheable = cacheable
	if best.State == nil {
		// No feasible partition was found — synthesise an "untouched" trailing state
		// that holds the starting hand and prior values. Callers (deck_eval) still need
		// a non-nil State to read end-of-turn fields.
		fallback := masterState.Copy()
		fallback.SetHand(append([]card.Card(nil), hand...))
		fallback.SetArsenal(arsenalCardIn)
		best.State = fallback
	}
	if best.State.Arsenal() == nil {
		promoteRandomHandCardToArsenal(&best, hand, arsenalCardIn)
	}
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
		if best.BestLine[i].Role == Held && best.BestLine[i].Card.ID() == chosen.ID() {
			best.BestLine[i].Role = Arsenal
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

// groupByRoleInto appends hand cards into caller-provided pitched/attackers/defenders
// slices.
func groupByRoleInto(hand []card.Card, roles []Role, pitched, attackers, defenders []card.Card) ([]card.Card, []card.Card, []card.Card) {
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

// gatherHeldCards appends every hand card with role Held into the caller-provided held slice.
func gatherHeldCards(hand []card.Card, roles []Role, held []card.Card) []card.Card {
	for i, c := range hand {
		if roles[i] == Held {
			held = append(held, c)
		}
	}
	return held
}

// findArsenalCard returns the arsenal-in card when it stays in the arsenal slot, nil
// otherwise.
func findArsenalCard(rolesBuf []Role, arsenalCardIn card.Card, n int) card.Card {
	if arsenalCardIn != nil && rolesBuf[n] == Arsenal {
		return arsenalCardIn
	}
	return nil
}

// roleAllowed decides whether the partition enumerator may assign role r to the current
// card.
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
		return true
	}
	return r != Attack || canAttack
}

// defendersDamage tallies the total Value contribution of the partition's defense phase
// against the caller-supplied state engine. DRs resolve first; plain blocks then consume
// whatever incoming damage is left, capped per card. The engine is mutated in place: each
// DR's resolution updates the running aura set and the engine's IncomingDamage; the chain
// phase will read the post-defense graveyard via the engine's left-behind state.
//
// blockBudget is the remaining defense-phase pitch supply after the caller has subtracted
// DR costs. Modal blockers enumerate their modes within blockBudget and pick the one
// yielding the highest BonusDefense; non-modal Blockers run their hook unchanged.
//
// Returns the per-DR cacheable status as a sticky bit — once a DR reads deck or graveyard,
// the partition's defense-phase output isn't safe to cache.
func defendersDamage(defenders, pitched []card.Card, deckPile *deck.Deck, ge *gameengine.GameEngine, gravBuf []card.Card, cs *card.CardState, incomingDamage, blockBudget, arsenalDefenderIdx int) (int, []card.Card, bool) {
	total := 0
	remaining := incomingDamage
	cacheable := true
	ge.SetDeck(deckPile)
	for i, def := range defenders {
		if !attackerMetaPtrFor(def).actsAsDR {
			continue
		}
		gravBuf = append(gravBuf[:0], defenders...)
		ge.SetGraveyard(gravBuf)
		ge.SetPitched(pitched)
		ge.SetDefenders(defenders)
		ge.SetValue(0)
		ge.SetIncomingDamage(remaining)
		ge.SetCacheable(true)
		*cs = card.CardState{Card: def, FromArsenal: i == arsenalDefenderIdx}
		ge.ResolveChainStep(ge.Logger(), cs)
		total += ge.Value()
		remaining = ge.IncomingDamage()
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

// chainScoreCmp lexicographically compares two leaf scores by (value, cardsDrawn,
// futureValue).
func chainScoreCmp(aValue, aCardsDrawn, aFutureValue, bValue, bCardsDrawn, bFutureValue int) int {
	if aValue != bValue {
		if aValue > bValue {
			return 1
		}
		return -1
	}
	if aCardsDrawn != bCardsDrawn {
		if aCardsDrawn > bCardsDrawn {
			return 1
		}
		return -1
	}
	if aFutureValue != bFutureValue {
		if aFutureValue > bFutureValue {
			return 1
		}
		return -1
	}
	return 0
}
