package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// Rules-orchestration methods on *GameEngine — Fire* walks, DestroyAura, the turn-
// boundary advance. Each operates on the embedded *GameState's slices but applies
// game-rule semantics (cursor iteration for handler-side splices, OncePerTurn gating,
// FiredThisTurn accounting, post-fire trigger drainage).

// FireAttack walks the aura entries with TriggerType()==triggertype.Attack and invokes
// every one whose OncePerTurn gate is open. The triggering card is published on
// triggeringCard so handlers can attribute log lines back to the source. Cursor-based
// iteration so a handler-side splice (Destroy) advances only when the slice length
// didn't change.
func (g *GameEngine) FireAttack(triggeringCard card.Card) {
	g.fireMatching(triggeringCard, triggertype.Attack)
}

// FireAttackAction is the triggertype.AttackAction counterpart to FireAttack: walks the
// aura entries matching that type and fires those whose OncePerTurn gate is open.
func (g *GameEngine) FireAttackAction(triggeringCard card.Card) {
	g.fireMatching(triggeringCard, triggertype.AttackAction)
}

// fireMatching is the shared aura-fire walk for FireAttack / FireAttackAction /
// FireEndOfTurn. Iterates auras with a cursor so handler-side splicing (Destroy mutates
// the auras slice in place, shifting the next entry down to the cursor's index)
// advances only when the slice length didn't change.
func (g *GameEngine) fireMatching(triggeringCard card.Card, trigger triggertype.Type) {
	for i := 0; i < len(g.auras); {
		a := g.auras[i]
		if a.TriggerType() != trigger || (a.OncePerTurn() && a.FiredThisTurn()) {
			i++
			continue
		}
		g.triggeringCard = triggeringCard
		g.currentAuraIdx = i
		g.currentAuraDestroyed = false
		a.Fire(g, g.logger)
		g.currentAuraIdx = -1
		g.triggeringCard = nil
		if !g.currentAuraDestroyed {
			g.auras[i].SetFiredThisTurn(true)
			i++
		}
	}
}

// HasEndOfTurnFire reports whether either Auras or Triggers carries a
// triggertype.EndOfTurn entry. Lets the chain runner skip the end-of-turn walk when
// nothing would fire.
func (g *GameEngine) HasEndOfTurnFire() bool {
	for _, a := range g.auras {
		if a.TriggerType() == triggertype.EndOfTurn {
			return true
		}
	}
	for _, t := range g.triggers {
		if t.TriggerType() == triggertype.EndOfTurn {
			return true
		}
	}
	return false
}

// FireEndOfTurn runs after the chain has finished resolving (and the legality gates
// have passed) but before the carry state is captured. Walks Auras and Triggers in one
// pass each:
//
//   - Aura entries respect OncePerTurn / FiredThisTurn semantics; the handler owns
//     destruction via the engine's destroyAura path.
//   - Trigger entries are one-shot; fired entries are removed afterward. Snapshotting
//     len(g.triggers) before iterating keeps a handler that calls AddXxxTrigger from
//     firing its newcomer on the same pass — newcomers stay queued for the next matching
//     event.
func (g *GameEngine) FireEndOfTurn() {
	g.fireMatching(nil, triggertype.EndOfTurn)
	n := len(g.triggers)
	for i := 0; i < n; i++ {
		tr := g.triggers[i]
		if tr.TriggerType() != triggertype.EndOfTurn {
			continue
		}
		tr.Fire(g, g.logger)
	}
	kept := g.triggers[:0]
	for i, tr := range g.triggers {
		if i < n && tr.TriggerType() == triggertype.EndOfTurn {
			continue
		}
		kept = append(kept, tr)
	}
	g.triggers = kept
}

// FireHit walks the one-shot trigger queue and invokes every triggertype.Hit entry
// whose type filter matches the attacking card's types. Surviving entries (filter
// mismatch) are kept; fired entries are removed.
func (g *GameEngine) FireHit(attackerTypes card.TypeSet) {
	kept := g.triggers[:0]
	for i := range g.triggers {
		t := g.triggers[i]
		if t.TriggerType() != triggertype.Hit || !t.Matches(attackerTypes) {
			kept = append(kept, t)
			continue
		}
		t.Fire(g, g.logger)
	}
	g.triggers = kept
}

// FireStartOfTurn walks g.auras and invokes every triggertype.StartOfTurn entry,
// calling onFire with each entry's pre-state snapshot so sim can attribute damage /
// draws / log lines back to the firing aura. Auras that destroy themselves splice out;
// FiredThisTurn flips reset on each fresh turn boundary.
//
// The onFire callback receives:
//   - pre is the index in g.auras of the firing entry at the time of the call.
//   - damage is g.value's delta during this handler — the partition tiebreaker uses it.
//   - drawnCard is the first card the handler appended to hand, or nil. Used by
//     processAurasAtStartOfTurn to surface "revealed" entries.
//   - newLogEntries is the slice of LogEntries the handler appended (caller may copy
//     out).
func (g *GameEngine) FireStartOfTurn(onFire func(idx int, damage int, drawnCard card.Card, newLogEntries []turnlogger.LogEntry)) {
	for i := 0; i < len(g.auras); {
		a := g.auras[i]
		g.auras[i].SetFiredThisTurn(false)
		if a.TriggerType() != triggertype.StartOfTurn {
			i++
			continue
		}
		preHand := len(g.hand)
		preLog := 0
		if g.logger != nil {
			preLog = len(g.logger.Entries())
		}
		preValue := g.value
		g.currentAuraIdx = i
		g.currentAuraDestroyed = false
		a.Fire(g, g.logger)
		g.currentAuraIdx = -1

		damage := g.value - preValue
		var drawn card.Card
		if len(g.hand) > preHand {
			drawn = g.hand[preHand]
		}
		var newEntries []turnlogger.LogEntry
		if g.logger != nil {
			if entries := g.logger.Entries(); len(entries) > preLog {
				newEntries = entries[preLog:]
			}
		}
		if onFire != nil {
			onFire(i, damage, drawn, newEntries)
		}
		if !g.currentAuraDestroyed {
			i++
		}
	}
}

// AdvanceTurnBoundary clears the per-turn FiredThisTurn flag on every persisted aura.
// The chain runner calls this when advancing across the turn boundary so the
// OncePerTurn gate rearms.
func (g *GameEngine) AdvanceTurnBoundary() {
	for i := range g.auras {
		g.auras[i].SetFiredThisTurn(false)
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
	i := g.currentAuraIdx
	if i < 0 || i >= len(g.auras) {
		return
	}
	if addToGraveyard {
		g.auras[i].OnDestroy(g)
	}
	g.auras = append(g.auras[:i], g.auras[i+1:]...)
	g.currentAuraDestroyed = true
}
