// Package hand evaluates the value of a hand of Flesh and Blood cards
// played in isolation.
package hand

import (
	"slices"
	"strings"
	"sync"

	"github.com/timch/fab-deck-optimizer/internal/card"
)

// Role is what a card does on a given turn cycle.
type Role uint8

const (
	Pitch Role = iota
	Attack
	Defend
)

// HandSize is the fixed number of cards in a hand. Best() assumes this
// exact size; raise and adjust memo key types together if it ever changes.
const HandSize = 4

// Play is the chosen partition for a hand: one role per card, plus the
// resulting damage dealt and damage prevented. Roles are aligned to the
// hand after Best has sorted it into canonical order.
type Play struct {
	Roles     []Role
	Dealt     int
	Prevented int
}

// Value returns the total value of the play (damage dealt + damage prevented).
func (p Play) Value() int { return p.Dealt + p.Prevented }

// Best returns the optimal Play for the given hand against an opponent that
// will attack for incomingDamage on their next turn. The hand is sorted
// in place into canonical order; the returned Roles align with that order.
//
// Cards are partitioned into three roles:
//   - Pitch: contributes its Pitch value as resources
//   - Attack: consumes Cost resources, contributes Attack to damage dealt
//   - Defend: held in hand, contributes Defend to damage prevented (capped
//     at incomingDamage; excess block is wasted)
//
// Results are memoized on the sorted hand + incoming damage.
func Best(hand []card.Card, incomingDamage int) Play {
	slices.SortFunc(hand, cardCompare)

	var key memoKey
	copy(key.cards[:], hand)
	key.incoming = int32(incomingDamage)

	if v, ok := memoLoad(key); ok {
		roles := make([]Role, HandSize)
		copy(roles, v.roles[:])
		return Play{Roles: roles, Dealt: int(v.dealt), Prevented: int(v.prevented)}
	}

	p := solve(hand, incomingDamage)

	var v memoVal
	v.dealt = int32(p.Dealt)
	v.prevented = int32(p.Prevented)
	copy(v.roles[:], p.Roles)
	memoStore(key, v)

	return p
}

// solve is the uncached brute-force optimizer (3^N partitions).
func solve(hand []card.Card, incomingDamage int) Play {
	n := len(hand)
	best := Play{Roles: make([]Role, n)}
	roles := make([]Role, n)

	var recurse func(i int)
	recurse = func(i int) {
		if i == n {
			score(hand, roles, incomingDamage, &best)
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

func score(hand []card.Card, roles []Role, incoming int, best *Play) {
	var resources, costs, dealt, defense int
	for i, c := range hand {
		switch roles[i] {
		case Pitch:
			resources += c.Pitch
		case Attack:
			costs += c.Cost
			dealt += c.Attack
		case Defend:
			defense += c.Defend
		}
	}
	if resources < costs {
		return
	}
	prevented := defense
	if prevented > incoming {
		prevented = incoming
	}
	v := dealt + prevented
	if v > best.Dealt+best.Prevented {
		best.Dealt = dealt
		best.Prevented = prevented
		copy(best.Roles, roles)
	}
}

// cardCompare defines a total order over Card for canonicalization.
// Returns negative / zero / positive per slices.SortFunc convention.
func cardCompare(a, b card.Card) int {
	if a.Cost != b.Cost {
		return a.Cost - b.Cost
	}
	if a.Pitch != b.Pitch {
		return a.Pitch - b.Pitch
	}
	if a.Attack != b.Attack {
		return a.Attack - b.Attack
	}
	if a.Defend != b.Defend {
		return a.Defend - b.Defend
	}
	return strings.Compare(a.Name, b.Name)
}

// --- memo cache ---

type memoKey struct {
	cards    [HandSize]card.Card
	incoming int32
}

type memoVal struct {
	dealt     int32
	prevented int32
	roles     [HandSize]Role
}

var (
	memoMu sync.RWMutex
	memo   = map[memoKey]memoVal{}
)

func memoLoad(k memoKey) (memoVal, bool) {
	memoMu.RLock()
	v, ok := memo[k]
	memoMu.RUnlock()
	return v, ok
}

func memoStore(k memoKey, v memoVal) {
	memoMu.Lock()
	memo[k] = v
	memoMu.Unlock()
}
