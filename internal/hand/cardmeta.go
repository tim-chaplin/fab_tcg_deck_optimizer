package hand

// Per-card metadata cache: scalar attributes (types, cost bounds, GoAgain, attack-action
// membership) that playSequence reads in its hot inner loop, hoisted out of interface-dispatch
// via a lazily-populated table sized for the full card-ID space.

import (
	"sync"
	"sync/atomic"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// attackerMeta caches the scalar card attributes playSequence reads on every permutation. With
// this hoisted to a per-attacker lookup, the hot inner loop skips Types / GoAgain interface
// dispatch; the one meta build amortises across all N! permutations.
//
// minCost / maxCost are static bounds on Card.Cost(s). For cards implementing card.VariableCost
// the solver uses them for O(1) partition pre-screens and falls through to card.Cost(state) in
// the chain inner loop. For non-VariableCost cards, minCost == maxCost == Cost(&TurnState{})
// and the cached value is used directly (no interface call per play).
type attackerMeta struct {
	types            card.TypeSet
	card             card.Card // held for variable-cost chain-time Cost(state) calls
	minCost          int
	maxCost          int
	isVariable       bool
	baseGoAgain      bool
	isAttackOrWeapon bool
	// isAttackAction is the "attack action card" test (Action+Attack, no Weapon) the sim uses
	// to pick which Play resolutions fire TriggerAttackAction AuraTriggers. Weapons carry
	// TypeAttack but aren't attack action CARDS; only the Action+Attack bitmask matches the
	// printed trigger text on cards like Malefic Incantation.
	isAttackAction bool
}

// costAt returns the card's effective cost given the current TurnState. Static cards return the
// cached value directly; variable-cost cards defer to card.Cost(s) so every game-state-dependent
// costing rule lives inside the card, not the solver.
func (m *attackerMeta) costAt(s *card.TurnState) int {
	if m.isVariable {
		return m.card.Cost(s)
	}
	return m.maxCost
}

// cardMetaCache / cardMetaReady are shared, read-only-after-init card metadata tables. Populated
// lazily via cardMetaSlowPath on first encounter, then read from all goroutines without sync.
// Sized for the full uint16 ID space so lookups are plain bounds-checked reads (~2 MB total).
const cardMetaCacheSize = 1 << 16

var (
	cardMetaCache [cardMetaCacheSize]attackerMeta
	cardMetaReady [cardMetaCacheSize]uint32 // written once (atomically) per ID; 0 = unready, 1 = ready
	cardMetaMu    sync.Mutex
)

// attackerMetaPtrFor returns a pointer to cached metadata for c, populating on first encounter.
// Hands back a direct pointer into the global cache so permutation swaps move 8 bytes instead of
// a full attackerMeta struct. The target is read-only after initialisation. Safe from multiple
// goroutines: the first writer per ID holds the mutex, later readers see the ready flag set with
// a release barrier and read the immutable meta entry directly.
func attackerMetaPtrFor(c card.Card) *attackerMeta {
	id := c.ID()
	if atomic.LoadUint32(&cardMetaReady[id]) == 1 {
		return &cardMetaCache[id]
	}
	cardMetaSlowPath(c, id)
	return &cardMetaCache[id]
}

// cardMetaSlowPath populates the cache entry under cardMetaMu and returns the computed meta.
func cardMetaSlowPath(c card.Card, id card.ID) attackerMeta {
	cardMetaMu.Lock()
	defer cardMetaMu.Unlock()
	// Re-check under lock: another goroutine may have populated between the atomic load and here.
	if atomic.LoadUint32(&cardMetaReady[id]) == 1 {
		return cardMetaCache[id]
	}
	t := c.Types()
	m := attackerMeta{
		types:            t,
		card:             c,
		baseGoAgain:      c.GoAgain(),
		isAttackOrWeapon: t.Has(card.TypeAttack) || t.Has(card.TypeWeapon),
		isAttackAction:   t.IsAttackAction(),
	}
	if vc, ok := c.(card.VariableCost); ok {
		m.minCost = vc.MinCost()
		m.maxCost = vc.MaxCost()
		m.isVariable = m.minCost != m.maxCost
	} else {
		// Static cost: any TurnState probe returns the same value. Cache once.
		fixed := c.Cost(&card.TurnState{})
		m.minCost = fixed
		m.maxCost = fixed
	}
	cardMetaCache[id] = m
	atomic.StoreUint32(&cardMetaReady[id], 1)
	return m
}
