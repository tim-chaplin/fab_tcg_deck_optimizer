// Package hand evaluates the value of a hand of Flesh and Blood cards played in isolation.
package hand

import (
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Role is what a card does on a given turn cycle.
type Role uint8

const (
	Pitch Role = iota
	Attack
	Defend
)

// Play is the chosen partition for a hand: one role per card, plus the resulting damage dealt and
// damage prevented. Best sorts the caller's hand into canonical order in place, and Roles are
// aligned to that post-sort order. (Weapon swing decisions are not reported in Roles — they're
// consumed only for their damage contribution.)
type Play struct {
	Roles     []Role
	Dealt     int
	Prevented int
}

// Value returns the total value of the play (damage dealt + damage prevented).
func (p Play) Value() int { return p.Dealt + p.Prevented }

// String returns a human-readable role name ("PITCH", "ATTACK", "DEFEND").
func (r Role) String() string {
	switch r {
	case Pitch:
		return "PITCH"
	case Attack:
		return "ATTACK"
	case Defend:
		return "DEFEND"
	}
	return "UNKNOWN"
}

// FormatRoles pairs each card in hand with its assigned role for debug output, e.g.
// "Hocus Pocus (Blue): PITCH, Runic Reaping (Red): ATTACK".
func FormatRoles(hand []card.Card, roles []Role) string {
	parts := make([]string, len(hand))
	for i, c := range hand {
		parts[i] = c.Name() + ": " + roles[i].String()
	}
	return strings.Join(parts, ", ")
}

// Best returns the optimal Play for the given hand against an opponent that will attack for
// incomingDamage on their next turn. Any equipped weapons may also be swung for their Cost if
// resources allow.
//
// Cards are partitioned into three roles:
//   - Pitch: contributes its Pitch value as resources.
//   - Attack: consumes Cost resources; the attack is resolved by calling Card.Play in some order
//     the optimizer chooses. Effects on TurnState carry forward to later attacks in the same
//     sequence.
//   - Defend: contributes Defense to damage prevented (capped at incomingDamage; excess block is
//     wasted).
//
// The optimizer brute-forces all 3^N partitions, then for each legal partition enumerates every
// subset of weapons to swing and every ordering of the combined attacker list. For N=4 with 0–2
// weapons that remains well under 10k evaluations.
//
// Results are memoized on (hero name, sorted weapon names, sorted card IDs, incomingDamage) so
// that repeated evaluations of the same hand across shuffles short-circuit. The hand is sorted
// in place into canonical order first — Roles in the returned Play align with that post-sort
// order. Every card in the hand must be registered in package cards; Best panics otherwise.
func Best(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card) Play {
	ids := handIDs(hand)
	sort.Sort(&handByID{hand: hand, ids: ids})

	// Unmemoable hands skip the cache read but still write — the stale entry is harmless since
	// future unmemoable lookups will skip the read too.
	key := formatMemoKey(hero, weapons, ids, incomingDamage)
	if isMemoable(hand) {
		memoMu.RLock()
		cached, hit := memo[key]
		memoMu.RUnlock()
		if hit {
			// Clone Roles so callers can't mutate the cached slice.
			roles := make([]Role, len(cached.Roles))
			copy(roles, cached.Roles)
			return Play{Roles: roles, Dealt: cached.Dealt, Prevented: cached.Prevented}
		}
	}
	result := bestUncached(hero, weapons, hand, incomingDamage, deck)
	memoMu.Lock()
	memo[key] = result
	memoMu.Unlock()
	return result
}

// isMemoable reports whether a hand's Best result can be cached. Any card implementing
// card.NoMemo (e.g. one whose Play depends on deck composition not captured by the memo key)
// disqualifies the whole hand.
func isMemoable(hand []card.Card) bool {
	for _, c := range hand {
		if _, ok := c.(card.NoMemo); ok {
			return false
		}
	}
	return true
}

// handIDs returns the registry IDs for each card in hand, preserving order. Panics if any card
// isn't registered in package cards.
func handIDs(hand []card.Card) []cards.ID {
	ids := make([]cards.ID, len(hand))
	for i, c := range hand {
		id, found := cards.ByName(c.Name())
		if !found {
			panic("hand: card not in index: " + c.Name())
		}
		ids[i] = id
	}
	return ids
}

// handByID sorts a parallel (hand, ids) pair by ascending ID.
type handByID struct {
	hand []card.Card
	ids  []cards.ID
}

func (h *handByID) Len() int           { return len(h.ids) }
func (h *handByID) Less(i, j int) bool { return h.ids[i] < h.ids[j] }
func (h *handByID) Swap(i, j int) {
	h.ids[i], h.ids[j] = h.ids[j], h.ids[i]
	h.hand[i], h.hand[j] = h.hand[j], h.hand[i]
}

// memo caches canonical-order results keyed by formatMemoKey.
var (
	memoMu sync.RWMutex
	memo   = map[string]Play{}
)

// formatMemoKey serializes the canonical key fields into a string. The hand must already be
// sorted by card ID; weapon names are sorted here.
func formatMemoKey(hero hero.Hero, weapons []weapon.Weapon, sortedIDs []cards.ID, incoming int) string {
	wnames := make([]string, len(weapons))
	for i, w := range weapons {
		wnames[i] = w.Name()
	}
	sort.Strings(wnames)

	var b strings.Builder
	b.WriteString(hero.Name())
	b.WriteByte('|')
	b.WriteString(strings.Join(wnames, ","))
	b.WriteByte('|')
	for i, id := range sortedIDs {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.FormatUint(uint64(id), 10))
	}
	b.WriteByte('|')
	b.WriteString(strconv.Itoa(incoming))
	return b.String()
}

func bestUncached(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card) Play {
	n := len(hand)
	best := Play{Roles: make([]Role, n)}
	roles := make([]Role, n)

	var recurse func(i int)
	recurse = func(i int) {
		if i == n {
			evalPartition(hero, weapons, hand, roles, incomingDamage, deck, &best)
			return
		}
		for r := Role(0); r <= Defend; r++ {
			roles[i] = r
			recurse(i + 1)
		}
	}
	recurse(0)
	return best
}

// canAfford reports whether the pitch resources produced by pitched cover the combined Cost of
// every card in attackers. A partition is legal only if its attackers (hand cards plus any weapon
// swings joined into the same list) can all be paid for.
func canAfford(pitched, attackers []card.Card) bool {
	resources := 0
	for _, c := range pitched {
		resources += c.Pitch()
	}
	cost := 0
	for _, c := range attackers {
		cost += c.Cost()
	}
	return resources >= cost
}

func evalPartition(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, roles []Role, incoming int, deck []card.Card, best *Play) {
	pitched, attackers, defenders := groupByRole(hand, roles)
	if !canAfford(pitched, attackers) {
		return
	}
	prevented := preventedDamage(defenders, incoming)
	defenseDealt := defenseReactionDamage(defenders, pitched, deck)
	attackDealt := bestAttackWithWeapons(hero, weapons, attackers, pitched, deck)

	totalDealt := attackDealt + defenseDealt
	v := totalDealt + prevented
	if v > best.Dealt+best.Prevented {
		best.Dealt = totalDealt
		best.Prevented = prevented
		copy(best.Roles, roles)
	}
}

// groupByRole buckets `hand` by `roles` (parallel slices) into pitchers, attackers, and
// defenders. Order within each bucket matches the cards' positions in `hand`.
func groupByRole(hand []card.Card, roles []Role) (pitched, attackers, defenders []card.Card) {
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
	return
}

// preventedDamage is the damage a wall of `defenders` blocks against `incoming`: the sum of
// printed Defense, capped at incoming (excess block is wasted).
func preventedDamage(defenders []card.Card, incoming int) int {
	total := 0
	for _, d := range defenders {
		total += d.Defense()
	}
	if total > incoming {
		return incoming
	}
	return total
}

// defenseReactionDamage runs the Play() hook of every Defense Reaction in `defenders` and sums
// the damage they deal back to the attacker (e.g. Weeping Battleground's 1 arcane on banish).
// Played in isolation — no attack ordering; TurnState only carries Pitched/Deck so effects that
// check "what was pitched" work. Uncapped: this damage is dealt, not prevented.
func defenseReactionDamage(defenders, pitched, deck []card.Card) int {
	total := 0
	for _, d := range defenders {
		if !d.Types()["Defense Reaction"] {
			continue
		}
		state := card.TurnState{Pitched: pitched, Deck: deck}
		total += d.Play(&state)
	}
	return total
}

// bestAttackWithWeapons enumerates every subset of `weapons` to swing alongside `attackers` and
// returns the max damage over all (affordable) weapon masks. Each selected weapon adds its Cost
// and joins the attacker permutation inside bestAttackDamage.
func bestAttackWithWeapons(hero hero.Hero, weapons []weapon.Weapon, attackers, pitched, deck []card.Card) int {
	best := 0
	for mask := 0; mask < 1<<len(weapons); mask++ {
		allAttackers := append([]card.Card(nil), attackers...)
		for i, w := range weapons {
			if mask&(1<<i) != 0 {
				allAttackers = append(allAttackers, w)
			}
		}
		if !canAfford(pitched, allAttackers) {
			continue
		}
		if dealt := bestAttackDamage(hero, allAttackers, pitched, deck); dealt > best {
			best = dealt
		}
	}
	return best
}

// bestAttackDamage tries every ordering of attackers and returns the max total damage after Play
// is called on each in sequence. Between each attacker's Play() and its append to CardsPlayed,
// the hero's OnCardPlayed hook fires so triggered abilities (e.g. Viserai's Runechants)
// contribute.
//
// Each permutation gets its own freshly-allocated []*PlayedCard so that any grants applied via
// CardsRemaining mutation don't leak across permutations.
func bestAttackDamage(hero hero.Hero, attackers, pitched, deck []card.Card) int {
	if len(attackers) == 0 {
		return 0
	}
	perm := make([]card.Card, len(attackers))
	copy(perm, attackers)
	best := 0
	permute(perm, 0, func(order []card.Card) {
		if !isLegalOrder(hero, pitched, deck, order) {
			return
		}
		if d := evaluateAttackDamage(hero, pitched, deck, order); d > best {
			best = d
		}
	})
	return best
}

// evaluateAttackDamage plays `order` as an attack sequence and returns total damage dealt —
// every card's Play() plus the hero's OnCardPlayed trigger. Assumes the ordering is legal (see
// isLegalOrder); chain-legality is not rechecked here. Each call allocates fresh *PlayedCard
// wrappers so nothing leaks back to the caller.
func evaluateAttackDamage(hero hero.Hero, pitched, deck []card.Card, order []card.Card) int {
	played := make([]*card.PlayedCard, len(order))
	for i, c := range order {
		played[i] = &card.PlayedCard{Card: c}
	}
	state := card.TurnState{Pitched: pitched, Deck: deck}
	total := 0
	for i, pc := range played {
		state.CardsRemaining = played[i+1:]
		state.Self = pc
		total += pc.Card.Play(&state)
		total += hero.OnCardPlayed(pc.Card, &state)
		state.CardsPlayed = append(state.CardsPlayed, pc.Card)
	}
	return total
}

// isLegalOrder plays `order` as an attack sequence (same loop as bestAttackDamage's permutation
// callback) and reports whether every non-final card had EffectiveGoAgain when it finished
// playing. Damage is not tallied — this is a pure legality check, useful for asserting that
// specific orderings the solver would have to consider are rejected. Each call allocates fresh
// *PlayedCard wrappers.
func isLegalOrder(hero hero.Hero, pitched, deck []card.Card, order []card.Card) bool {
	played := make([]*card.PlayedCard, len(order))
	for i, c := range order {
		played[i] = &card.PlayedCard{Card: c}
	}
	state := card.TurnState{Pitched: pitched, Deck: deck}
	for i, pc := range played {
		state.CardsRemaining = played[i+1:]
		state.Self = pc
		pc.Card.Play(&state)
		hero.OnCardPlayed(pc.Card, &state)
		state.CardsPlayed = append(state.CardsPlayed, pc.Card)
		if i < len(played)-1 && !pc.EffectiveGoAgain() {
			return false
		}
	}
	return true
}

func permute(a []card.Card, k int, emit func([]card.Card)) {
	if k == len(a)-1 {
		emit(a)
		return
	}
	for i := k; i < len(a); i++ {
		a[k], a[i] = a[i], a[k]
		permute(a, k+1, emit)
		a[k], a[i] = a[i], a[k]
	}
}
