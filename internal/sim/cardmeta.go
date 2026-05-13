package sim

// Per-card metadata cache: scalar attributes (types, cost bounds, GoAgain, attack-action
// membership) that playSequence reads in its hot inner loop, hoisted out of interface-dispatch
// via a lazily-populated table sized for the full card-ID space.

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// attackerMeta caches the scalar card attributes playSequence reads on every permutation. With
// this hoisted to a per-attacker lookup, the hot inner loop skips Types / GoAgain interface
// dispatch; the one meta build amortises across all N! permutations.
//
// minCost / maxCost are static bounds on Card.Cost(s). For cards implementing VariableCost
// the solver uses them for O(1) partition pre-screens and falls through to Cost(state) in
// the chain inner loop. For non-VariableCost cards, minCost == maxCost == Cost(&TurnState{})
// and the cached value is used directly (no interface call per play).
type attackerMeta struct {
	types      card.TypeSet
	card       card.Card // held for variable-cost / modal-cost chain-time Cost calls
	minCost    int
	maxCost    int
	isVariable bool
	// isAttack is the "this chain step is an attack" test driving fireAttackAuras — true on
	// any card carrying TypeAttack (attack action cards and weapon abilities both).
	isAttack bool
	// isAttackAction is the "attack action card" test (Action+Attack, no Weapon) gating
	// TriggerAttackAction Auras like Malefic Incantation. Weapon abilities carry TypeAttack
	// but aren't attack action CARDS, so they miss this gate.
	isAttackAction bool
	// isFreeChainStep is set on cards that resolve in the chain without paying an Action
	// Point — Instants and Attack Reactions (both 0 AP per FaB rules). Action cards and
	// weapon swings cost 1 AP and don't set this.
	isFreeChainStep bool
	// isModalCost is set when the card implements ModalCost — costAt dispatches on self.Mode
	// instead of the static maxCost / Cost(s) paths.
	isModalCost bool
	// isModalBlocker is set when the card has a per-mode block-time cost (Blocker +
	// Modal + BlockCost). Cached so containsModalBlocker / defendersDamage's mode-pick
	// fast paths avoid per-leaf type assertions.
	isModalBlocker bool
	// actsAsDR is true for printed Defense Reactions and DefensiveInstant-marked Instants.
	// Cached so per-leaf hot-path defender checks fold the type-assertion into the
	// once-per-ID slow path.
	actsAsDR bool
	// modes is the mode count for a Modal, 1 for non-modal cards. Sized int8 so it
	// packs into the bool block's padding without growing attackerMeta — every chain step
	// reads permMeta[i] in the inner loop, and every extra cache line through that table
	// shows up in the anneal bench.
	modes int8
}

// costAt returns the card's effective cost given the current TurnState and chosen mode.
// ModalCost cards dispatch on the mode index; static cards return the cached value
// directly; VariableCost cards defer to Cost(s) so every game-state-dependent costing rule
// lives inside the card, not the solver.
func (m *attackerMeta) costAt(s *gameengine.GameEngine, mode int8) int {
	if m.isModalCost {
		return m.card.(card.ModalCost).ModalCost(mode)
	}
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
//
// Test-only cards sharing ids.InvalidCard collide on cache slot 0; the InvalidCard branch
// builds a fresh meta inline so each gets its own. Production never sees InvalidCard.
func attackerMetaPtrFor(c card.Card) *attackerMeta {
	id := c.ID()
	if id == ids.InvalidCard {
		m := buildAttackerMeta(c)
		return &m
	}
	if atomic.LoadUint32(&cardMetaReady[id]) == 1 {
		return &cardMetaCache[id]
	}
	cardMetaSlowPath(c, id)
	return &cardMetaCache[id]
}

// cardMetaSlowPath populates the cache entry under cardMetaMu and returns the computed meta.
func cardMetaSlowPath(c card.Card, id ids.CardID) attackerMeta {
	cardMetaMu.Lock()
	defer cardMetaMu.Unlock()
	// Re-check under lock: another goroutine may have populated between the atomic load and here.
	if atomic.LoadUint32(&cardMetaReady[id]) == 1 {
		return cardMetaCache[id]
	}
	m := buildAttackerMeta(c)
	cardMetaCache[id] = m
	atomic.StoreUint32(&cardMetaReady[id], 1)
	return m
}

// buildAttackerMeta computes a fresh attackerMeta from c. Shared by the cache slow path
// and the InvalidCard bypass in attackerMetaPtrFor.
func buildAttackerMeta(c card.Card) attackerMeta {
	t := c.Types(nil)
	_, isDefensiveInstant := c.(card.DefensiveInstant)
	m := attackerMeta{
		types:           t,
		card:            c,
		isAttack:        t.Has(card.TypeAttack),
		isAttackAction:  t.IsAttackAction(),
		isFreeChainStep: t.Has(card.TypeInstant) || t.IsAttackReaction(),
		actsAsDR:        t.IsDefenseReaction() || isDefensiveInstant,
		modes:           1,
	}
	if mc, ok := c.(card.Modal); ok {
		// A Modal must expose at least two modes — the marker exists to enumerate
		// across them. Returning 0 would silently zero the chain (outer loop runs zero
		// iterations); returning 1 makes the marker pointless. Panic so the bug surfaces
		// at first encounter rather than corrupting solver output.
		n := mc.Modes()
		if n < 2 {
			panic(fmt.Sprintf("Modal %s: Modes() = %d, want >= 2", c.Name(), n))
		}
		m.modes = int8(n)
		// Modal blockers (also Blocker + BlockCost) get a flag for the defendersDamage
		// fast path — it scans defenders per leaf and the type assertions add up.
		if _, ok := c.(card.Blocker); ok {
			if _, ok := c.(card.BlockCost); ok {
				m.isModalBlocker = true
			}
		}
	}
	if mc, ok := c.(card.ModalCost); ok {
		// ModalCost overrides VariableCost / static Cost in costAt — folds the per-mode
		// costs into the partition pre-screen so MinCost / MaxCost still bound the search.
		m.isModalCost = true
		minC, maxC := mc.ModalCost(0), mc.ModalCost(0)
		for i := int8(1); i < m.modes; i++ {
			c := mc.ModalCost(i)
			if c < minC {
				minC = c
			}
			if c > maxC {
				maxC = c
			}
		}
		m.minCost = minC
		m.maxCost = maxC
		m.isVariable = minC != maxC
	} else if vc, ok := c.(card.VariableCost); ok {
		m.minCost = vc.MinCost()
		m.maxCost = vc.MaxCost()
		m.isVariable = m.minCost != m.maxCost
	} else {
		// Static cost: any TurnState probe returns the same value.
		fixed := c.Cost(&gameengine.GameEngine{GameState: gameengine.NewState()})
		m.minCost = fixed
		m.maxCost = fixed
	}
	return m
}
