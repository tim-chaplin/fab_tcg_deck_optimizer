package hand

// Top-level hand enumeration: bestUncached walks every partition (Pitch / Attack / Defend /
// Held / Arsenal assignment) and delegates each leaf's chain-feasibility check to
// bestAttackWithWeapons. Post-enumeration helpers decide how an empty arsenal slot gets
// filled, plus the beatsBest / roleAllowed policy functions that shape the partition tree.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

func (e *Evaluator) bestUncached(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card, priorAuraTriggers []card.AuraTrigger) TurnSummary {
	n := len(hand)
	// The partition recurse treats the arsenal-in card as an extra entry at index n with a
	// restricted role menu (Arsenal / Attack / Defend), so everything about it is decided inside
	// the enumeration. totalN is the effective size of BestLine.
	totalN := n
	if arsenalCardIn != nil {
		totalN = n + 1
	}

	// Seed best.LeftoverRunechants with the carryover: partitions with no attacks don't reduce
	// it, so carryover is the baseline to beat. BestLine starts with every hand card Held and
	// the arsenal-in card (if any) staying in the slot, so a hand with no Value-adding partition
	// still reports sensible "nothing played, nothing pitched" assignments.
	best := TurnSummary{BestLine: make([]CardAssignment, totalN), LeftoverRunechants: runechantCarryover}
	// bestSwung holds the winning partition's swung weapon names so fillContributions can rebuild
	// the chain it runs bestSequence over. Lives outside TurnSummary since weapons are
	// recoverable from AttackChain once fillContributions finishes. bestBudget captures the
	// winning phase-split's chain-resource state; the replay re-seeds ctx with it so
	// bestSequence finds the exact permutation that won during enumeration.
	var bestSwung []string
	var bestBudget chainBudget
	// bestHasHeld tracks whether the current best has at least one Held hand card — lets
	// beatsBest distinguish "arsenal will be occupied post-hoc" from "arsenal will be empty."
	// Seeded true when the hand is non-empty: the initial best puts every hand card into Held,
	// so a post-hoc promotion would fill arsenal. Candidates need both a Value/leftover tie and
	// some way to end with arsenal occupied to displace it.
	bestHasHeld := n > 0
	// bestFutureValuePlayed tracks how many card.AddsFutureValue cards the current best is
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
		best.ArsenalCard = arsenalCardIn
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
	pitched := bufs.pitchedBuf
	attackers := bufs.attackersBuf
	defenders := bufs.defendersBuf
	held := bufs.heldBuf

	var recurse func(i, pitchSum, defenseSum int)
	recurse = func(i, pitchSum, defenseSum int) {
		if i == totalN {
			prevented := defenseSum
			if prevented > incomingDamage {
				prevented = incomingDamage
			}
			// Group roles into played / pitched / defending buckets. Iterates the hand (size n),
			// then layers in the arsenal slot (index n) based on its assigned role. Arsenal-role
			// cards contribute nothing this turn whether they came from hand or the slot.
			var p, a, d []card.Card
			hasAnyDefender := hasReactions
			if !hasAnyDefender && arsenalCardIn != nil && rolesBuf[n] == Defend {
				hasAnyDefender = true
			}
			if hasAnyDefender {
				p, a, d = groupByRoleInto(hand, rolesBuf[:n], pitched[:0], attackers[:0], defenders[:0])
				if arsenalCardIn != nil {
					switch rolesBuf[n] {
					case Attack:
						a = append(a, arsenalCardIn)
					case Defend:
						d = append(d, arsenalCardIn)
					}
				}
			} else {
				p, a = groupPitchAttack(hand, rolesBuf[:n], pitched[:0], attackers[:0])
				if arsenalCardIn != nil && rolesBuf[n] == Attack {
					a = append(a, arsenalCardIn)
				}
			}
			// Arsenal-in is appended last to a / d above, so its index in the attackers slice is
			// len(a)-1 when present in the chain. -1 means no arsenal-in card in the attackers
			// (either no arsenal-in card at all, or it took a different role).
			arsenalInIdx := -1
			if arsenalCardIn != nil && rolesBuf[n] == Attack {
				arsenalInIdx = len(a) - 1
			}
			// Held cards (hand cards left without a Pitch / Attack / Defend role) thread
			// through to TurnState.Held so alt-cost effects (e.g. Moon Wish's "use a Held
			// card") can read len > 0. Arsenal-in can never be Held (roleAllowed bars it).
			h := gatherHeldCards(hand, rolesBuf[:n], held[:0])
			attackDealt, defenseDealt, leftoverRunechants, budget, swung, ok := bestAttackWithWeapons(hero, weapons, a, d, p, h, deck, bufs, runechantCarryover, incomingDamage, defenseSum, arsenalInIdx, priorAuraTriggers)
			if !ok {
				return
			}

			v := attackDealt + defenseDealt + prevented
			arsenalCard := findArsenalCard(rolesBuf, arsenalCardIn, n)
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
			bestWillOccupy := best.ArsenalCard != nil || bestHasHeld
			if !beatsBest(v, leftoverRunechants, futureValuePlayed, willOccupy, best, bestFutureValuePlayed, bestWillOccupy) {
				return
			}
			best.Value = v
			bestSwung = swung
			bestBudget = budget
			best.LeftoverRunechants = leftoverRunechants
			best.ArsenalCard = arsenalCard
			bestHasHeld = hasHeld
			bestFutureValuePlayed = futureValuePlayed
			// Write the winning roles into BestLine. Cards and FromArsenal flags were populated
			// at construction; only Role varies. Contribution is cleared here and filled by
			// fillContributions below for the winning line.
			for j := 0; j < totalN; j++ {
				best.BestLine[j].Role = rolesBuf[j]
				best.BestLine[j].Contribution = 0
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
	// Once per Best call, on the winning line only, attribute per-card contribution.
	if len(best.BestLine) > 0 {
		fillContributions(&best, hero, weapons, bestSwung, bestBudget, deck, bufs, incomingDamage, runechantCarryover, priorAuraTriggers)
	}
	// If the arsenal slot is empty after enumeration, promote one Held card into it. Held hand
	// cards and Held mid-turn-drawn cards are treated as one pool — neither source is preferred,
	// because both end the turn as a single card of equivalent future-turn value. The pick is
	// deterministic per-hand (hashed from sorted card IDs + drawn card IDs + arsenal-in ID) so
	// the memo stays consistent, but spreads across candidates to avoid a lowest-ID bias.
	if best.ArsenalCard == nil {
		promoteRandomHeldToArsenal(&best, hand, n, arsenalCardIn)
	}
	return best
}

// promoteRandomHeldToArsenal picks one Held card — a hand card in best.BestLine or a mid-turn-
// drawn card in best.Drawn — and flips its role to Arsenal. Both sources share a single
// candidate pool so the draw isn't preferred over hand Helds (nor the other way around). Cards
// in best.ReturnedToTopOfDeck are skipped on the BestLine side because alt-cost effects already
// re-routed them; their copies typically reappear in best.Drawn (e.g. Moon Wish's alt cost
// puts a Held card on top of deck and Sun Kiss's tutor draws it). No-op when nothing is Held.
func promoteRandomHeldToArsenal(best *TurnSummary, hand []card.Card, n int, arsenalCardIn card.Card) {
	handHeldCount := countHeldInBestLine(best.BestLine, n, best.ReturnedToTopOfDeck)
	drawnHeldCount := countHeldInDrawn(best.Drawn)
	total := handHeldCount + drawnHeldCount
	if total == 0 {
		return
	}
	pick := int(arsenalPromotionHash(hand, best.Drawn, arsenalCardIn) % uint64(total))
	// Walk hand Helds first (in BestLine order), then drawn Helds (in draw order), mapping pick
	// to the matching slot.
	if pick < handHeldCount {
		promoteNthHeldInBestLine(best, n, pick)
		return
	}
	promoteNthHeldInDrawn(best, pick-handHeldCount)
}

// countHeldInBestLine returns how many of the first n BestLine entries are still Role=Held
// after enumeration AND haven't been moved to deck top by an alt-cost effect (per
// returnedToTopOfDeck). The first-n restriction excludes any arsenal-in entry (which lives
// at index n and is never Held); the returnedToTopOfDeck skip prevents the moved copies
// from competing for arsenal against the drawn copies they spawned.
func countHeldInBestLine(line []CardAssignment, n int, returnedToTopOfDeck []card.Card) int {
	c := 0
	consumed := append([]card.Card(nil), returnedToTopOfDeck...)
	for i := 0; i < n; i++ {
		if line[i].Role != Held {
			continue
		}
		if idx := indexOfCard(consumed, line[i].Card); idx >= 0 {
			consumed = append(consumed[:idx], consumed[idx+1:]...)
			continue
		}
		c++
	}
	return c
}

// indexOfCard returns the position of the first card in cs whose ID matches c, or -1 when
// none. Used by the returnedToTopOfDeck match in countHeldInBestLine and promoteNthHeldInBestLine.
func indexOfCard(cs []card.Card, c card.Card) int {
	for i, x := range cs {
		if x.ID() == c.ID() {
			return i
		}
	}
	return -1
}

// countHeldInDrawn returns how many mid-turn-drawn cards are still Role=Held after the winning
// chain resolved.
func countHeldInDrawn(drawn []CardAssignment) int {
	c := 0
	for i := range drawn {
		if drawn[i].Role == Held {
			c++
		}
	}
	return c
}

// arsenalPromotionHash computes the deterministic bucket seed that picks which Held card fills
// an empty arsenal slot. Uses FNV-1a over the sorted hand IDs + drawn card IDs + arsenal-in ID —
// the only requirement is a uniform spread across bucket counts 1..total so no lowest-ID bias
// creeps in while the memo stays consistent per hand.
func arsenalPromotionHash(hand []card.Card, drawn []CardAssignment, arsenalCardIn card.Card) uint64 {
	const (
		fnvOffsetBasis uint64 = 1469598103934665603
		fnvPrime       uint64 = 1099511628211
	)
	h := fnvOffsetBasis
	for _, c := range hand {
		h ^= uint64(c.ID())
		h *= fnvPrime
	}
	for _, d := range drawn {
		h ^= uint64(d.Card.ID())
		h *= fnvPrime
	}
	if arsenalCardIn != nil {
		h ^= uint64(arsenalCardIn.ID())
		h *= fnvPrime
	}
	return h
}

// promoteNthHeldInBestLine flips the pick-th non-consumed Held hand card (in BestLine order)
// to Arsenal, and records it on best.ArsenalCard. Caller guarantees pick < count of eligible
// Held entries (per countHeldInBestLine, which applies the same returnedToTopOfDeck skip).
func promoteNthHeldInBestLine(best *TurnSummary, n, pick int) {
	idx := 0
	consumed := append([]card.Card(nil), best.ReturnedToTopOfDeck...)
	for i := 0; i < n; i++ {
		if best.BestLine[i].Role != Held {
			continue
		}
		if k := indexOfCard(consumed, best.BestLine[i].Card); k >= 0 {
			consumed = append(consumed[:k], consumed[k+1:]...)
			continue
		}
		if idx == pick {
			best.BestLine[i].Role = Arsenal
			best.ArsenalCard = best.BestLine[i].Card
			return
		}
		idx++
	}
}

// promoteNthHeldInDrawn flips the pick-th Held mid-turn-drawn card (in draw order) to Arsenal,
// and records it on best.ArsenalCard. Caller guarantees pick < count of Held drawn entries.
func promoteNthHeldInDrawn(best *TurnSummary, pick int) {
	idx := 0
	for i := range best.Drawn {
		if best.Drawn[i].Role != Held {
			continue
		}
		if idx == pick {
			best.Drawn[i].Role = Arsenal
			best.ArsenalCard = best.Drawn[i].Card
			return
		}
		idx++
	}
}

// groupByRoleInto appends hand cards into caller-provided pitched/attackers/defenders slices
// (passed pre-reset to length 0) to avoid per-partition heap allocation.
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

// gatherHeldCards appends every hand card with role Held into the caller-provided held slice
// (passed pre-reset to length 0) and returns it. Threads the partition's Held set into
// bestAttackWithWeapons so alt-cost effects can consult it via TurnState.Held.
func gatherHeldCards(hand []card.Card, roles []Role, held []card.Card) []card.Card {
	for i, c := range hand {
		if roles[i] == Held {
			held = append(held, c)
		}
	}
	return held
}

// groupPitchAttack is the reaction-free leaf's grouping step: skips the defenders bucket (only
// needed for Defense-Reaction-Play dispatch, which this path doesn't run).
func groupPitchAttack(hand []card.Card, roles []Role, pitched, attackers []card.Card) ([]card.Card, []card.Card) {
	for i, c := range hand {
		switch roles[i] {
		case Pitch:
			pitched = append(pitched, c)
		case Attack:
			attackers = append(attackers, c)
		}
	}
	return pitched, attackers
}

// findArsenalCard returns the arsenal-in card when it stays in the arsenal slot, nil otherwise.
// Hand cards never take Arsenal role during enumeration (post-hoc promotion handles that), so
// the only slot that can be Arsenal is the arsenal-in slot at index n.
func findArsenalCard(rolesBuf []Role, arsenalCardIn card.Card, n int) card.Card {
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
	if leftoverRunechants > best.LeftoverRunechants {
		return true
	}
	if leftoverRunechants < best.LeftoverRunechants {
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

// defenseReactionDamage runs Play() for every Defense Reaction in defenders and sums the damage
// they deal back to the attacker (e.g. a banish-an-aura-for-arcane rider). Played in isolation
// — no attack ordering; TurnState carries Pitched / Deck plus a per-DR fresh copy of the
// defenders list in Graveyard so effects that scan the graveyard see plain blocks and other
// defenders. Uncapped: this damage is dealt, not prevented.
//
// state is caller-provided (from attackBufs) and reset per call. gravBuf is the caller-owned
// scratch backing state.Graveyard; the returned slice is the (possibly grown) buffer for reuse.
func defenseReactionDamage(defenders, pitched, deck []card.Card, state *card.TurnState, gravBuf []card.Card, cs *card.CardState) (int, []card.Card) {
	total := 0
	for _, d := range defenders {
		if !d.Types().IsDefenseReaction() {
			continue
		}
		gravBuf = append(gravBuf[:0], defenders...)
		*state = card.TurnState{Pitched: pitched, Deck: deck, Graveyard: gravBuf}
		*cs = card.CardState{Card: d}
		total += d.Play(state, cs)
	}
	return total, gravBuf
}

// chainBudget captures the winning phase-split's attack-chain resource state. Reusing it to seed
// the replay ctx in fillContributions ensures playSequenceWithMeta finds the exact permutation
// that won during partition enumeration — critical for per-card attribution since different
// permutations can deal different per-card damage.
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
func containsDefenseReaction(cards []card.Card) bool {
	for _, c := range cards {
		if c.Types().IsDefenseReaction() {
			return true
		}
	}
	return false
}
