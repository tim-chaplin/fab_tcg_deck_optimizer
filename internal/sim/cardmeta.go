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
// minCost / maxCost are static bounds on the actual paid cost (the cheapest of Card.Cost()
// and any AlternativeCost / VariableCost / ModalCost branch). The solver uses them for
// O(1) partition pre-screens, then falls through to costAt at play time for the live
// figure. Non-variable, non-modal, non-alt-cost cards have minCost == maxCost == Cost().
type attackerMeta struct {
	types      card.TypeSet
	card       card.Card // held for variable-cost / modal-cost / alt-cost play-time Cost calls
	minCost    int
	maxCost    int
	isVariable bool
	// hasAlternativeCost is set when the card implements card.AlternativeCost. costAt
	// picks the cheaper of Card.Cost() and AlternativeCost(g) at play time, flipping
	// PaidAlternativeCost on the CardState when the alt branch wins. Static bound: alt
	// branches reduce the cost (and never raise it), so minCost folds the alt's lower
	// bound in conservatively.
	hasAlternativeCost bool
	// isAttack is the "this attack step is an attack" test driving fireAttackAuras — true on
	// any card carrying TypeAttack (attack action cards and weapon abilities both). For
	// ModalTypes cards, this is the mode-0 value; the attack-turn runner uses isAttackAt(mode)
	// instead so per-mode type-line changes propagate.
	isAttack bool
	// isFreeAttackStep is set on cards that resolve in the attack turn without paying an Action
	// Point — Instants and Attack Reactions (both 0 AP per FaB rules). Action cards and
	// weapon swings cost 1 AP and don't set this. For ModalTypes cards, this is the
	// mode-0 value; the attack-turn runner uses isFreeAttackStepAt(mode) instead.
	isFreeAttackStep bool
	// isModalCost is set when the card implements ModalCost — costAt dispatches on self.Mode
	// instead of the static maxCost / Cost(s) paths.
	isModalCost bool
	// isModalTypes is set when the card implements ModalTypes — typesAt / isAttackAt /
	// isFreeAttackStepAt dispatch on self.Mode instead of returning the cached static fields.
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
	// to gate the per-attack-step type assertion — only ~6 cards in the registry implement
	// it, so the 99%+ case folds into a single bool read.
	hasPlayPrecondition bool
	// modes is the mode count for a Modal, 1 for non-modal cards. Sized int8 so it
	// packs into the bool block's padding alongside the bools above — every attack step
	// reads permMeta[i] in the inner loop, and every extra cache line through that table
	// shows up in the anneal bench.
	modes int8
	// typesByMode is the per-mode TypeSet table for ModalTypes cards (nil for everyone
	// else). Length equals modes; entry i is the type-line when self.Mode == i. The
	// attack-turn runner reads it via typesWithMode(mode) when isModalTypes is set; non-modal
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

// isFreeAttackStepWithMode reports whether the card-at-mode resolves without paying an
// AP. Reads the per-mode TypeSet for ModalTypes cards; falls through to the cached
// isFreeAttackStep otherwise.
func (m *attackerMeta) isFreeAttackStepWithMode(mode int8) bool {
	if m.isModalTypes {
		t := m.typesByMode[mode]
		return t.Has(card.TypeInstant) || t.IsAttackReaction()
	}
	return m.isFreeAttackStep
}

// costAt returns the card's effective cost given the current engine state and chosen
// mode, plus whether the AlternativeCost branch was chosen. ModalCost dispatches on mode;
// VariableCost defers to EffectiveCost(g); AlternativeCost picks min(Cost, alt) when alt is
// available. paidAlt is true only when the AlternativeCost branch wins — the attack-turn
// runner mirrors it onto pc.PaidAlternativeCost so the card's Play body branches correctly.
func (m *attackerMeta) costAt(ge *gameengine.GameEngine, mode int8) (cost int, paidAlt bool) {
	if m.isModalCost {
		return m.card.(card.ModalCost).ModalCost(mode), false
	}
	if m.isVariable {
		return m.card.(card.VariableCost).EffectiveCost(ge), false
	}
	base := m.maxCost
	if m.hasAlternativeCost {
		if alt, ok := m.card.(card.AlternativeCost).AlternativeCost(ge); ok && alt < base {
			return alt, true
		}
	}
	return base, false
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
		isFreeAttackStep:    t.Has(card.TypeInstant) || t.IsAttackReaction(),
		actsAsDR:            t.IsDefenseReaction() || isDefensiveInstant,
		hasPlayPrecondition: hasPlayPre,
		modes:               1,
	}
	if mc, ok := c.(card.Modal); ok {
		// A Modal must expose at least two modes — the marker exists to enumerate
		// across them. Returning 0 would silently zero the attack turn (outer loop runs zero
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
		// attack-turn runner reads typesByMode[mode] directly. Must be paired with Modal to
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
		// VariableCost: MinCost is the static lower bound; the printed Cost() is the upper.
		m.minCost = vc.MinCost()
		m.maxCost = c.Cost()
		m.isVariable = m.minCost != m.maxCost
	} else {
		// Static or alternative cost: Cost() is the printed value, identical at every call.
		// The static lower bound is Cost() unless an AlternativeCost branch is wired in,
		// which can probe to 0 — the conservative bound for alt-cost cards.
		fixed := c.Cost()
		m.maxCost = fixed
		m.minCost = fixed
	}
	if _, ok := c.(card.AlternativeCost); ok {
		m.hasAlternativeCost = true
		// Alt branches never raise the cost; pre-screen with 0 as the worst-case lower
		// bound. Real branches usually return 0 (Moon Wish, Rise Above); a card returning
		// a higher alt cost would still be over-bounded by 0 — sound, just looser.
		m.minCost = 0
	}
	return m
}
