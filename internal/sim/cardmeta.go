package sim

// Per-card metadata cache: scalar attributes (types, cost bounds, GoAgain, attack-action)
// playSequence reads in its hot loop, hoisted out of interface dispatch via a lazily-
// populated table sized for the full card-ID space.

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// attackerMeta caches the scalar card attributes playSequence reads per permutation. The
// hot loop skips Types / GoAgain interface dispatch; one meta build amortises across N!.
//
// minCost / maxCost are static bounds on Card.Cost. VariableCost cards: solver uses them
// for O(1) partition pre-screens, then falls through to Cost(state) in the chain loop.
// Non-VariableCost: minCost == maxCost == Cost(&TurnState{}), used directly.
type attackerMeta struct {
	types      card.TypeSet
	card       card.Card // held for variable-cost / modal-cost chain-time Cost calls
	minCost    int
	maxCost    int
	isVariable bool
	// isAttack is the "this chain step is an attack" test driving fireAttackAuras — true on
	// any card carrying TypeAttack (attack action cards and weapon abilities both). For
	// ModalTypes cards, this is the mode-0 value; the chain runner uses isAttackAt(mode)
	// instead so per-mode type-line changes propagate.
	isAttack bool
	// isFreeChainStep is set on cards that resolve in the chain without paying an Action
	// Point — Instants and Attack Reactions (both 0 AP per FaB rules). Action cards and
	// weapon swings cost 1 AP and don't set this. For ModalTypes cards, this is the
	// mode-0 value; the chain runner uses isFreeChainStepAt(mode) instead.
	isFreeChainStep bool
	// isModalCost is set when the card implements ModalCost — costAt dispatches on self.Mode
	// instead of the static maxCost / Cost(s) paths.
	isModalCost bool
	// isModalTypes is set when the card implements ModalTypes — typesAt / isAttackAt /
	// isFreeChainStepAt dispatch on self.Mode instead of returning the cached static fields.
	isModalTypes bool
	// isModalBlocker is set when the card has a per-mode block-time cost (Blocker +
	// Modal + BlockCost). Cached so containsModalBlocker / defendersDamage's mode-pick
	// fast paths avoid per-leaf type assertions.
	isModalBlocker bool
	// actsAsDR is true for printed Defense Reactions and DefensiveInstant-marked Instants.
	// Cached so per-leaf hot-path defender checks fold the type-assertion into the
	// once-per-ID slow path.
	actsAsDR bool
	// hasPlayPrecondition is true when the card implements card.PlayPrecondition. Cached
	// to gate the per-chain-step type assertion — only ~6 cards in the registry implement
	// it, so the 99%+ case folds into a single bool read.
	hasPlayPrecondition bool
	// modes is the mode count for a Modal, 1 for non-modal cards. Sized int8 so it
	// packs into the bool block's padding alongside the bools above — every chain step
	// reads permMeta[i] in the inner loop, and every extra cache line through that table
	// shows up in the anneal bench.
	modes int8
	// typesByMode is the per-mode TypeSet table for ModalTypes cards (nil for everyone
	// else). Length equals modes; entry i is the type-line when self.Mode == i. The
	// chain runner reads it via typesWithMode(mode) when isModalTypes is set; non-modal
	// cards stay on the .types fast path.
	typesByMode []card.TypeSet
}

// typesWithMode returns the TypeSet this card resolves with for the given mode. Falls
// back to the cached static .types for non-ModalTypes cards.
func (m *attackerMeta) typesWithMode(mode int8) card.TypeSet {
	if m.isModalTypes {
		return m.typesByMode[mode]
	}
	return m.types
}

// isFreeChainStepWithMode reports whether the card-at-mode resolves without paying an
// AP. Reads the per-mode TypeSet for ModalTypes cards; falls through to the cached
// isFreeChainStep otherwise.
func (m *attackerMeta) isFreeChainStepWithMode(mode int8) bool {
	if m.isModalTypes {
		t := m.typesByMode[mode]
		return t.Has(card.TypeInstant) || t.IsAttackReaction()
	}
	return m.isFreeChainStep
}

// verifyStaticCost arms the static-cost assertion in costAt. A divergence between the
// re-probed Cost(ge) and the cached bound means the card varies cost with game state but
// declares neither VariableCost nor ModalCost, so the attack-budget prune bounds it from a
// wrong fixed cost and can skip a fundable line. Off in production to avoid the probe.
var verifyStaticCost bool

// costAt returns the card's effective cost given the current TurnState and chosen mode.
// ModalCost cards dispatch on the mode index; static cards return the cached value
// directly; VariableCost cards defer to Cost(s) so every game-state-dependent costing rule
// lives inside the card, not the solver.
func (m *attackerMeta) costAt(ge *gameengine.GameEngine, mode int8) int {
	if m.isModalCost {
		return m.card.(card.ModalCost).ModalCost(mode)
	}
	if m.isVariable {
		return m.card.Cost(ge)
	}
	if verifyStaticCost {
		if live := m.card.Cost(ge); live != m.maxCost {
			panic(fmt.Sprintf("card %q varies Cost(g) with game state (cached %d, live %d) "+
				"but implements neither VariableCost nor ModalCost — declare one so the "+
				"attack-budget prune bounds it correctly", m.card.Name(), m.maxCost, live))
		}
	}
	return m.maxCost
}

// cardMetaCache / cardMetaReady are read-only-after-init metadata tables, populated lazily
// by cardMetaSlowPath. Sized for the full uint16 ID space so lookups are bounds-checked
// reads (~2 MB total).
const cardMetaCacheSize = 1 << 16

var (
	cardMetaCache [cardMetaCacheSize]attackerMeta
	cardMetaReady [cardMetaCacheSize]uint32 // written once (atomically) per ID; 0 = unready, 1 = ready
	cardMetaMu    sync.Mutex
)

// attackerMetaPtrFor returns a cache pointer for c, populating on first encounter. Direct
// pointer return lets perm swaps move 8 bytes instead of a full struct. Read-only post-init,
// safe from multiple goroutines via the cardMetaReady atomic flag.
//
// Test-only cards sharing ids.InvalidCard would collide on slot 0; that branch builds a
// fresh meta inline. Production never sees InvalidCard.
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
	_, hasPlayPre := c.(card.PlayPrecondition)
	m := attackerMeta{
		types:               t,
		card:                c,
		isAttack:            t.Has(card.TypeAttack),
		isFreeChainStep:     t.Has(card.TypeInstant) || t.IsAttackReaction(),
		actsAsDR:            t.IsDefenseReaction() || isDefensiveInstant,
		hasPlayPrecondition: hasPlayPre,
		modes:               1,
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
	if mt, ok := c.(card.ModalTypes); ok {
		// ModalTypes cards have a per-mode TypeSet table; precompute it once so the
		// chain runner reads typesByMode[mode] directly. Must be paired with Modal to
		// give the table its length.
		if m.modes < 2 {
			panic(fmt.Sprintf("ModalTypes %s: also needs Modal with Modes() >= 2 (got modes=%d)", c.Name(), m.modes))
		}
		m.isModalTypes = true
		m.typesByMode = make([]card.TypeSet, m.modes)
		for i := int8(0); i < m.modes; i++ {
			m.typesByMode[i] = mt.TypesForMode(nil, i)
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
		fixed := c.Cost(gameengine.New())
		m.minCost = fixed
		m.maxCost = fixed
	}
	return m
}
