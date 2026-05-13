package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// Rules-orchestration methods on *GameEngine — Fire* walks, DestroyAura, the turn-
// boundary advance. Each operates on g.state's slices but applies game-rule semantics
// (cursor iteration for handler-side splices, OncePerTurn gating, FiredThisTurn
// accounting, post-fire trigger drainage).

// FireAttack walks the trigger entries (today, aura entries; in the future, other
// engine-level subscribers) with TriggerType()==triggertype.Attack and invokes every one
// whose OncePerTurn gate is open. The triggering card is published on
// g.state.triggeringCard so handlers can attribute log lines back to the source.
// Cursor-based iteration so a handler-side splice (Destroy) advances only when the slice
// length didn't change.
func (g *GameEngine) FireAttack(triggeringCard card.Card) {
	g.fireMatching(triggeringCard, triggertype.Attack)
}

// FireAttackAction is the triggertype.AttackAction counterpart to FireAttack: walks the
// trigger entries matching that type and fires those whose OncePerTurn gate is open.
func (g *GameEngine) FireAttackAction(triggeringCard card.Card) {
	g.fireMatching(triggeringCard, triggertype.AttackAction)
}

// fireMatching is the shared aura-fire walk for FireAttack / FireAttackAction /
// FireEndOfTurn. Iterates g.state.auras with a cursor so handler-side splicing
// (Destroy mutates the auras slice in place, shifting the next entry down to the
// cursor's index) advances only when the slice length didn't change.
func (g *GameEngine) fireMatching(triggeringCard card.Card, trigger triggertype.Type) {
	s := g.state
	for i := 0; i < len(s.auras); {
		a := s.auras[i]
		if a.TriggerType() != trigger || (a.OncePerTurn() && a.FiredThisTurn()) {
			i++
			continue
		}
		s.triggeringCard = triggeringCard
		s.currentAuraIdx = i
		s.currentAuraDestroyed = false
		a.Fire(g, s.logger)
		s.currentAuraIdx = -1
		s.triggeringCard = nil
		if !s.currentAuraDestroyed {
			s.auras[i].SetFiredThisTurn(true)
			i++
		}
	}
}

// HasEndOfTurnFire reports whether either Auras or Triggers carries a
// triggertype.EndOfTurn entry. Lets the chain runner skip the end-of-turn walk when
// nothing would fire.
func (g *GameEngine) HasEndOfTurnFire() bool {
	for _, a := range g.state.auras {
		if a.TriggerType() == triggertype.EndOfTurn {
			return true
		}
	}
	for _, t := range g.state.triggers {
		if t.TriggerType() == triggertype.EndOfTurn {
			return true
		}
	}
	return false
}

// FireEndOfTurn runs after the chain has finished resolving (and the legality gates have
// passed) but before the carry state is captured. Walks Auras and Triggers in one pass
// each:
//
//   - Aura entries respect OncePerTurn / FiredThisTurn semantics; the handler owns
//     destruction via the engine's destroyAura path.
//   - Trigger entries are one-shot; fired entries are removed afterward. Snapshotting
//     len(g.state.triggers) before iterating keeps a handler that calls AddXxxTrigger from
//     firing its newcomer on the same pass — newcomers stay queued for the next matching
//     event.
func (g *GameEngine) FireEndOfTurn() {
	g.fireMatching(nil, triggertype.EndOfTurn)
	s := g.state
	n := len(s.triggers)
	for i := 0; i < n; i++ {
		tr := s.triggers[i]
		if tr.TriggerType() != triggertype.EndOfTurn {
			continue
		}
		tr.Fire(g, s.logger)
	}
	kept := s.triggers[:0]
	for i, tr := range s.triggers {
		if i < n && tr.TriggerType() == triggertype.EndOfTurn {
			continue
		}
		kept = append(kept, tr)
	}
	s.triggers = kept
}

// FireHit walks the one-shot trigger queue and invokes every triggertype.Hit entry
// whose type filter matches the attacking card's types. Surviving entries (filter
// mismatch) are kept; fired entries are removed.
func (g *GameEngine) FireHit(attackerTypes card.TypeSet) {
	s := g.state
	kept := s.triggers[:0]
	for i := range s.triggers {
		t := s.triggers[i]
		if t.TriggerType() != triggertype.Hit || !t.Matches(attackerTypes) {
			kept = append(kept, t)
			continue
		}
		t.Fire(g, s.logger)
	}
	s.triggers = kept
}

// FireStartOfTurn walks g.state.auras and invokes every triggertype.StartOfTurn entry,
// calling onFire with each entry's pre-state snapshot so sim can attribute damage /
// draws / log lines back to the firing aura. Auras that destroy themselves splice out;
// FiredThisTurn flips reset on each fresh turn boundary.
//
// The onFire callback receives:
//   - pre is the index in g.state.auras of the firing entry at the time of the call.
//   - damage is g.state.value's delta during this handler — the partition tiebreaker
//     uses it.
//   - drawnCard is the first card the handler appended to hand, or nil. Used by
//     processAurasAtStartOfTurn to surface "revealed" entries.
//   - newLogEntries is the slice of LogEntries the handler appended (caller may copy
//     out).
func (g *GameEngine) FireStartOfTurn(onFire func(idx int, damage int, drawnCard card.Card, newLogEntries []turnlogger.LogEntry)) {
	s := g.state
	for i := 0; i < len(s.auras); {
		a := s.auras[i]
		s.auras[i].SetFiredThisTurn(false)
		if a.TriggerType() != triggertype.StartOfTurn {
			i++
			continue
		}
		preHand := len(s.hand)
		preLog := 0
		if s.logger != nil {
			preLog = len(s.logger.Entries())
		}
		preValue := s.value
		s.currentAuraIdx = i
		s.currentAuraDestroyed = false
		a.Fire(g, s.logger)
		s.currentAuraIdx = -1

		damage := s.value - preValue
		var drawn card.Card
		if len(s.hand) > preHand {
			drawn = s.hand[preHand]
		}
		var newEntries []turnlogger.LogEntry
		if s.logger != nil {
			if entries := s.logger.Entries(); len(entries) > preLog {
				newEntries = entries[preLog:]
			}
		}
		if onFire != nil {
			onFire(i, damage, drawn, newEntries)
		}
		if !s.currentAuraDestroyed {
			i++
		}
	}
}

// AdvanceTurnBoundary clears the per-turn FiredThisTurn flag on every persisted aura.
// The chain runner calls this when advancing across the turn boundary so the
// OncePerTurn gate rearms.
func (g *GameEngine) AdvanceTurnBoundary() {
	for i := range g.state.auras {
		g.state.auras[i].SetFiredThisTurn(false)
	}
}

// DestroyAura removes the aura currently being fired and, when addToGraveyard==true,
// invokes the aura's OnDestroy hook to push the aura's source card into the graveyard
// (token auras no-op). Direct splice (no cacheable flip) — destruction is deterministic
// from the triggering event, not hidden state.
//
// Called by the card.Aura context the engine threads into each handler; cards do not
// call this directly.
func (g *GameEngine) DestroyAura(addToGraveyard bool) {
	s := g.state
	i := s.currentAuraIdx
	if i < 0 || i >= len(s.auras) {
		return
	}
	if addToGraveyard {
		s.auras[i].OnDestroy(g)
	}
	s.auras = append(s.auras[:i], s.auras[i+1:]...)
	s.currentAuraDestroyed = true
}

// CreateAura forwards to GameState.CreateAura. Cards reach this through the named
// CreateStartOfTurnAura / CreateOncePerTurnAttackActionAura helpers; sim's per-
// permutation seed calls it directly to re-add prior auras.
func (g *GameEngine) CreateAura(a Aura) { g.state.CreateAura(a) }

// CreateTrigger forwards to GameState.CreateTrigger.
func (g *GameEngine) CreateTrigger(t Trigger) { g.state.CreateTrigger(t) }

// CreateItem forwards to GameState.CreateItem.
func (g *GameEngine) CreateItem(i Item) { g.state.CreateItem(i) }
