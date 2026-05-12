package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// Reset rewrites every GameEngine field to its per-permutation starting value derived from
// seed. The engine reuses its own internal slice backings via append([:0], src...) so
// repeated Reset calls don't allocate after the first sizing — only mid-chain growth past
// the current cap forces a fresh backing.
//
// Deck handling: the engine owns a scratch *deck.Deck; the first Reset copies seed.Deck via
// Copy(), subsequent Resets refill via CopyFrom so the cards / weapons backing arrays stay
// reused. seed.Deck nil leaves the scratch in place (zero-sized).
//
// Logger: nil seed.Logger leaves g.logger nil — the find-best pass's short-circuit signal.
// A non-nil logger gets installed verbatim; callers are expected to point it at the right
// per-permutation backing buffer themselves before passing.
//
// Cacheable: starts true. The first card-driven deck / graveyard read in this permutation
// flips it false via the appropriate accessor's contract.
func (g *GameEngine) Reset(seed PermutationSeed) {
	g.hand = append(g.hand[:0], seed.Hand...)
	if g.deck == nil {
		// Always allocate a scratch deck so card-driven mutations (PrependToDeck, Opt, …)
		// don't panic on nil. seed.Deck nil is treated as an empty source: deck.Copy on a
		// nil pointer returns a fresh empty deck.
		g.deck = seed.Deck.Copy()
	} else {
		g.deck.CopyFrom(seed.Deck)
	}
	g.arsenal = seed.Arsenal
	g.graveyard = append(g.graveyard[:0], seed.Graveyard...)
	g.banished = append(g.banished[:0], seed.Banished...)
	g.auras = g.auras[:0]
	for _, a := range seed.Auras {
		g.auras = append(g.auras, a.Clone())
	}
	g.triggers = g.triggers[:0]
	g.items = g.items[:0]
	for _, i := range seed.Items {
		g.items = append(g.items, i.Clone())
	}
	g.cardsPlayed = g.cardsPlayed[:0]
	g.cardsRemaining = nil
	g.pitched = seed.Pitched
	g.defenders = seed.Defenders

	g.logger = seed.Logger
	g.triggeringCard = nil
	g.attackReactionTarget = nil

	g.actionPoints = 1
	g.value = 0
	g.cardsDrawn = 0
	g.incomingDamage = seed.IncomingDamage
	g.arcaneIncomingDamage = seed.ArcaneIncomingDamage
	g.blockTotal = seed.BlockTotal
	g.currentAuraIdx = -1

	g.cardBanished = false
	g.arcaneDamageDealt = false
	g.opponentMarked = seed.OpponentMarked
	g.auraCreated = false
	g.nonAttackActionPlayed = false
	g.currentAuraDestroyed = false
	g.currentStepRerouted = false
	g.cacheable = true
}

// Snapshot returns the engine's persistent state — the values that carry into next turn.
// Sim's CarryState.SnapshotFromTurn copies fields out of the returned PersistentSnapshot.
// Deck gets a fresh Copy() so subsequent permutations' deck mutations don't reach back
// into the snapshot.
func (g *GameEngine) Snapshot() PersistentSnapshot {
	var d *deck.Deck
	if g.deck != nil {
		d = g.deck.Copy()
	}
	return PersistentSnapshot{
		Hand:           appendCopy(nil, g.hand),
		Deck:           d,
		Arsenal:        g.arsenal,
		Graveyard:      appendCopy(nil, g.graveyard),
		Banished:       appendCopy(nil, g.banished),
		Auras:          appendAuraCopy(nil, g.auras),
		Items:          appendItemCopy(nil, g.items),
		CardsDrawn:     g.cardsDrawn,
		OpponentMarked: g.opponentMarked,
		LogEntries:     appendLogCopy(nil, g.logger.Entries()),
	}
}

// FireAttack walks the trigger entries (today, aura entries; in the future, other
// engine-level subscribers) with TriggerType()==TriggerAttack and invokes every one whose
// OncePerTurn gate is open. The triggering card is published on g.triggeringCard so
// handlers can attribute log lines back to the source. Cursor-based iteration so a
// handler-side splice (Destroy) advances only when the slice length didn't change.
func (g *GameEngine) FireAttack(triggeringCard card.Card) {
	g.fireMatching(triggeringCard, TriggerAttack)
}

// FireAttackAction is the TriggerAttackAction counterpart to FireAttack: walks the
// trigger entries matching that type and fires those whose OncePerTurn gate is open.
func (g *GameEngine) FireAttackAction(triggeringCard card.Card) {
	g.fireMatching(triggeringCard, TriggerAttackAction)
}

// fireMatching is the shared aura-fire walk for FireAttack / FireAttackAction /
// FireEndOfTurn. Iterates g.auras with a cursor so handler-side splicing (Destroy mutates
// g.auras in place, shifting the next entry down to the cursor's index) advances only
// when the slice length didn't change.
func (g *GameEngine) fireMatching(triggeringCard card.Card, trigger TriggerType) {
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

// HasEndOfTurnFire reports whether either Auras or Triggers carries a TriggerEndOfTurn
// entry. Lets the chain runner skip the end-of-turn walk when nothing would fire.
func (g *GameEngine) HasEndOfTurnFire() bool {
	for _, a := range g.auras {
		if a.TriggerType() == TriggerEndOfTurn {
			return true
		}
	}
	for _, t := range g.triggers {
		if t.TriggerType() == TriggerEndOfTurn {
			return true
		}
	}
	return false
}

// FireEndOfTurn runs after the chain has finished resolving (and the legality gates have
// passed) but before Snapshot captures the next-turn state. Walks Auras and Triggers in
// one pass each:
//
//   - Aura entries respect OncePerTurn / FiredThisTurn semantics; the handler owns
//     destruction via the engine's destroyAura path.
//   - Trigger entries are one-shot; fired entries are removed afterward. Snapshotting
//     len(g.triggers) before iterating keeps a handler that calls AddXxxTrigger from firing
//     its newcomer on the same pass — newcomers stay queued for the next matching event.
func (g *GameEngine) FireEndOfTurn() {
	g.fireMatching(nil, TriggerEndOfTurn)
	n := len(g.triggers)
	for i := 0; i < n; i++ {
		tr := g.triggers[i]
		if tr.TriggerType() != TriggerEndOfTurn {
			continue
		}
		tr.Fire(g, g.logger)
	}
	kept := g.triggers[:0]
	for i, tr := range g.triggers {
		if i < n && tr.TriggerType() == TriggerEndOfTurn {
			continue
		}
		kept = append(kept, tr)
	}
	g.triggers = kept
}

// FireHit walks the one-shot trigger queue and invokes every TriggerHit entry whose type
// filter matches the attacking card's types. Surviving entries (filter mismatch) are kept;
// fired entries are removed.
func (g *GameEngine) FireHit(attackerTypes card.TypeSet) {
	kept := g.triggers[:0]
	for i := range g.triggers {
		t := g.triggers[i]
		if t.TriggerType() != TriggerHit || !t.Matches(attackerTypes) {
			kept = append(kept, t)
			continue
		}
		t.Fire(g, g.logger)
	}
	g.triggers = kept
}

// FireStartOfTurn walks g.auras and invokes every TriggerStartOfTurn entry, calling
// onFire with each entry's pre-state snapshot so sim can attribute damage / draws /
// log lines back to the firing aura. Auras that destroy themselves splice out;
// FiredThisTurn flips reset on each fresh turn boundary (sim calls AdvanceTurnBoundary).
//
// The onFire callback receives:
//   - pre is the index in g.auras of the firing entry at the time of the call.
//   - damage is g.value's delta during this handler — the partition tiebreaker uses it.
//   - drawnCard is the first card the handler appended to hand, or nil. Used by
//     processAurasAtStartOfTurn to surface "revealed" entries.
//   - newLogEntries is the slice of LogEntries the handler appended (caller may copy out).
func (g *GameEngine) FireStartOfTurn(onFire func(idx int, damage int, drawnCard card.Card, newLogEntries []turnlogger.LogEntry)) {
	for i := 0; i < len(g.auras); {
		a := g.auras[i]
		// Re-arm OncePerTurn at each turn boundary.
		g.auras[i].SetFiredThisTurn(false)
		if a.TriggerType() != TriggerStartOfTurn {
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

// AdvanceTurnBoundary clears the per-turn FiredThisTurn flag on every persisted aura. The
// chain runner calls this when advancing across the turn boundary so the OncePerTurn gate
// rearms.
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
// Called by the card.Aura context the engine threads into each handler; cards do not call
// this directly.
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

// AppendAura appends a to the engine's aura list. Used by sim's card-aura constructors
// (CreateStartOfTurnAura, CreateAttackActionAura, …) which build the typed entry in sim
// and hand it to the engine. Flips AuraCreated so same-turn "if you've played or created
// an aura" riders see the entry.
func (g *GameEngine) AppendAura(a Aura) {
	g.auras = append(g.auras, a)
	g.auraCreated = true
}

// AppendTrigger appends t to the engine's one-shot trigger queue.
func (g *GameEngine) AppendTrigger(t Trigger) {
	g.triggers = append(g.triggers, t)
}

// AppendItem appends i to the engine's item list. Used by sim's token-item helpers; the
// chain runner enqueues each item's Ability as a playable activated ability each turn.
func (g *GameEngine) AppendItem(i Item) {
	g.items = append(g.items, i)
}

// SetArcaneDamageDealt flips the sticky flag. Used by sim's token aura handlers (the
// runechant aura sets it when its consume reaches the LikelyDamageHits window) and by the
// chain runner's arcane-damage helper.
func (g *GameEngine) SetArcaneDamageDealt(v bool) { g.arcaneDamageDealt = v }

// === Helpers internal to the engine ===

func appendCopy(dst []card.Card, src []card.Card) []card.Card {
	if len(src) == 0 {
		return dst
	}
	return append(dst, src...)
}

func appendAuraCopy(dst []Aura, src []Aura) []Aura {
	if len(src) == 0 {
		return dst
	}
	return append(dst, src...)
}

func appendItemCopy(dst []Item, src []Item) []Item {
	if len(src) == 0 {
		return dst
	}
	return append(dst, src...)
}

func appendLogCopy(dst []turnlogger.LogEntry, src []turnlogger.LogEntry) []turnlogger.LogEntry {
	if len(src) == 0 {
		return dst
	}
	return append(dst, src...)
}
