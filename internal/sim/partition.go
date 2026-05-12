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

func (e *Evaluator) findBest(hero Hero, weapons []Weapon, hand []card.Card, mp Matchup, d *deck.Deck, prior gameengine.Spec, skipLog bool) TurnSummary {
	// Install the hero on the scratch engine so cards' g.CurrentHeroClass /
	// g.HeroWantsLowerHealth reads return the right value across every permutation
	// the upcoming search drives through the engine.
	bufs := e.getAttackBufs(len(hand), weapons)
	bufs.state.SetHero(hero)
	bufs.drScratch.SetHero(hero)
	var cacheKey evalCacheKey
	cacheUsable := e.cache != nil
	if cacheUsable {
		var keyOK bool
		cacheKey, keyOK = makeCacheKey(hero, weapons, hand, prior)
		if !keyOK {
			cacheUsable = false
		}
	}
	if cacheUsable {
		if entry, ok := e.cache.lookup(cacheKey); ok {
			e.cache.hits.Add(1)
			return e.replayBest(entry, hero, weapons, hand, mp, d, prior, skipLog)
		}
		e.cache.misses.Add(1)
	}

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
		State: CarryState{
			Hand:    append([]card.Card(nil), hand...),
			Deck:    d.Copy(),
			Arsenal: arsenalCardIn,
			Auras:   auraSliceFromEngine(prior.Auras),
			Items:   itemSliceFromEngine(prior.Items),
		},
	}
	cacheable := true
	var bestSwung []string
	for i := 0; i < n; i++ {
		best.BestLine[i] = CardAssignment{Card: hand[i], Role: Held}
	}
	if arsenalCardIn != nil {
		best.BestLine[n] = CardAssignment{Card: arsenalCardIn, Role: Arsenal, FromArsenal: true}
	}

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
			attackDealt, defenseDealt, swung, carry, ok, leafCacheable, arsenalAtChainStart := e.evaluatePartition(
				hero, weapons, hand, d,
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
			arsenalCard := arsenalAtChainStart
			futureValuePlayed := pendingFutureValue(carry.Auras, carry.Items)
			if runningSeen {
				cmp := chainScoreCmp(v, carry.CardsDrawn, futureValuePlayed,
					best.Value, best.State.CardsDrawn, runningFutureValue)
				if cmp < 0 {
					return
				}
				if cmp == 0 {
					willOccupy := arsenalCard != nil || len(carry.Hand) > 0
					runningWillOccupy := best.State.Arsenal != nil || len(best.State.Hand) > 0
					if !willOccupy || runningWillOccupy {
						return
					}
				}
			}
			best.State.CopyFrom(&carry)
			best.State.Arsenal = arsenalCard
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
	if best.State.Arsenal == nil {
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

// promoteRandomHandCardToArsenal picks one card from best.State.Hand and moves it into
// best.State.Arsenal, removing it from State.Hand. Deterministic per-hand pick.
func promoteRandomHandCardToArsenal(best *TurnSummary, startingHand []card.Card, arsenalCardIn card.Card) {
	if len(best.State.Hand) == 0 {
		return
	}
	eligible := make([]int, 0, len(best.State.Hand))
	for i, c := range best.State.Hand {
		t := c.Types(nil)
		if t.Has(card.TypeBlock) || t.IsResource() {
			continue
		}
		eligible = append(eligible, i)
	}
	if len(eligible) == 0 {
		return
	}
	pick := eligible[int(arsenalPromotionHash(startingHand, best.State.Hand, arsenalCardIn)%uint64(len(eligible)))]
	chosen := best.State.Hand[pick]
	best.State.Arsenal = chosen
	best.State.Hand = append(best.State.Hand[:pick:pick], best.State.Hand[pick+1:]...)
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

// defendersDamage tallies the total Value contribution of the partition's defense phase.
// DRs resolve first via Play (their ApplyAndLogEffectiveDefense decrements
// state.incomingDamage and credits the block); plain blocks then consume whatever incoming
// damage is left, capped per card. The engine is caller-provided (from attackBufs) and
// reset per DR via a PermutationSeed scoped to this DR's view (pitched, deckPile copy,
// defenders-as-graveyard, running incoming damage).
//
// blockBudget is the remaining defense-phase pitch supply after the caller has subtracted
// DR costs. Modal blockers enumerate their modes within blockBudget and pick the one
// yielding the highest BonusDefense; non-modal Blockers run their hook unchanged.
//
// Returns the per-DR cacheable status as a sticky bit — once a DR reads deck or graveyard,
// the partition's defense-phase output isn't safe to cache.
func defendersDamage(defenders, pitched []card.Card, deckPile *deck.Deck, state *gameengine.GameEngine, gravBuf []card.Card, cs *card.CardState, incomingDamage, blockBudget, arsenalDefenderIdx int) (int, []card.Card, bool) {
	total := 0
	remaining := incomingDamage
	cacheable := true
	for i, def := range defenders {
		if !attackerMetaPtrFor(def).actsAsDR {
			continue
		}
		gravBuf = append(gravBuf[:0], defenders...)
		// Preserve the running aura set and logger across the per-DR reset so chain-created
		// auras persist for the caller and the recording / find-best logger choice carries
		// through.
		preservedAuras := state.Auras()
		preservedLogger := state.Logger()
		state.Reset(gameengine.PermutationSeed{
			Deck:           deckPile,
			Graveyard:      gravBuf,
			Auras:          preservedAuras,
			Pitched:        pitched,
			Defenders:      defenders,
			IncomingDamage: remaining,
			Logger:         preservedLogger,
		})
		*cs = card.CardState{Card: def, FromArsenal: i == arsenalDefenderIdx}
		state.ResolveChainStep(state.Logger(), cs)
		total += state.Value()
		remaining = state.IncomingDamage()
		if !state.IsCacheable() {
			cacheable = false
		}
	}
	state.SetDefenders(defenders)
	for _, def := range defenders {
		if attackerMetaPtrFor(def).actsAsDR {
			continue
		}
		bestMode, bestCost := pickBlockerMode(def, state, cs, blockBudget)
		blockBudget -= bestCost
		*cs = card.CardState{Card: def, Mode: bestMode}
		if b, ok := def.(card.Blocker); ok {
			b.Block(state, state.Logger(), cs)
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
func pickBlockerMode(d card.Card, state *gameengine.GameEngine, cs *card.CardState, blockBudget int) (int8, int) {
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
		b.Block(state, state.Logger(), cs)
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
