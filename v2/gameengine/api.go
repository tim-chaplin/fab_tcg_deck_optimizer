package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// Reset rewrites every GameEngine field to its per-permutation starting value derived
// from seed. Aura / Trigger / Item lists are cleared; the caller follows up with
// CreateAura / CreateTrigger / CreateItem per entry that should carry into this
// permutation.
//
// Deck: seed.Deck nil collapses to a fresh empty deck (deck.Copy on a nil pointer
// returns &Deck{}); otherwise a fresh copy is made.
//
// Logger: nil seed.Logger leaves g.logger nil — the find-best pass's short-circuit
// signal. A non-nil logger gets installed verbatim.
//
// Cacheable: starts true. The first card-driven deck / graveyard read in this
// permutation flips it false via the appropriate accessor's contract.
func (g *GameEngine) Reset(seed PermutationSeed) {
	g.hand = append([]card.Card(nil), seed.Hand...)
	g.deck = seed.Deck.Copy()
	g.arsenal = seed.Arsenal
	g.graveyard = append([]card.Card(nil), seed.Graveyard...)
	g.banished = append([]card.Card(nil), seed.Banished...)
	g.auras = nil
	g.triggers = nil
	g.items = nil
	g.cardsPlayed = nil
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

// Copy returns a deep copy of the engine. Slice and *deck.Deck fields get fresh
// backing storage; Aura / Item entries are deep-copied via their Copy() methods so
// per-permutation Count / FiredThisTurn mutations stay isolated. Triggers are
// effectively immutable after construction, so only the slice header is duplicated.
// Logger is reset to nil — the caller installs a fresh per-clone logger when
// recording, leaving find-best clones allocation-free.
//
// Cards / sim drive every persistent mutation through this method when they want a
// fresh sandbox: the chain runner copies the start-of-turn engine per permutation,
// plays out, compares Value, and keeps the winning copy as next turn's master.
func (g *GameEngine) Copy() *GameEngine {
	out := *g
	out.hand = appendCopy(nil, g.hand)
	if g.deck != nil {
		out.deck = g.deck.Copy()
	}
	out.graveyard = appendCopy(nil, g.graveyard)
	out.banished = appendCopy(nil, g.banished)
	out.pitched = appendCopy(nil, g.pitched)
	out.defenders = appendCopy(nil, g.defenders)
	out.cardsPlayed = appendCopy(nil, g.cardsPlayed)
	if len(g.cardsRemaining) > 0 {
		out.cardsRemaining = append([]*card.CardState(nil), g.cardsRemaining...)
	} else {
		out.cardsRemaining = nil
	}
	if len(g.auras) > 0 {
		out.auras = make([]Aura, len(g.auras))
		for i, a := range g.auras {
			out.auras[i] = a.Copy()
		}
	} else {
		out.auras = nil
	}
	if len(g.triggers) > 0 {
		out.triggers = append([]Trigger(nil), g.triggers...)
	} else {
		out.triggers = nil
	}
	if len(g.items) > 0 {
		out.items = make([]Item, len(g.items))
		for i, it := range g.items {
			out.items[i] = it.Copy()
		}
	} else {
		out.items = nil
	}
	out.logger = nil
	return &out
}

// BeginPermutation resets per-chain locals on the engine in preparation for a fresh
// permutation's chain run. Auras, items, banished, graveyard, deck, arsenal, pitched,
// hero, and OpponentMarked carry over from the caller-side Copy() — they represent the
// leaf's pre-chain state. logger is installed verbatim (pass nil for the find-best path,
// a fresh logger for the recording path).
func (g *GameEngine) BeginPermutation(hand []card.Card, incomingDamage int, logger *turnlogger.TurnLogger) {
	g.hand = hand
	g.cardsPlayed = nil
	g.cardsRemaining = nil
	g.triggers = nil
	g.triggeringCard = nil
	g.attackReactionTarget = nil
	g.actionPoints = 1
	g.value = 0
	g.cardsDrawn = 0
	g.incomingDamage = incomingDamage
	g.cardBanished = false
	g.arcaneDamageDealt = false
	g.nonAttackActionPlayed = false
	g.currentAuraDestroyed = false
	g.currentStepRerouted = false
	g.currentAuraIdx = -1
	g.cacheable = true
	g.logger = logger
	g.auraCreated = len(g.auras) > 0
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

// CreateAura appends a to the engine's aura list. Used by sim's card-aura constructors
// (CreateStartOfTurnAura, CreateAttackActionAura, …) which build the typed entry in sim
// and hand it to the engine. Flips AuraCreated so same-turn "if you've played or created
// an aura" riders see the entry.
func (g *GameEngine) CreateAura(a Aura) {
	g.auras = append(g.auras, a)
	g.auraCreated = true
}

// CreateTrigger appends t to the engine's one-shot trigger queue.
func (g *GameEngine) CreateTrigger(t Trigger) {
	g.triggers = append(g.triggers, t)
}

// CreateItem appends i to the engine's item list. Used by sim's token-item helpers; the
// chain runner enqueues each item's Ability as a playable activated ability each turn.
func (g *GameEngine) CreateItem(i Item) {
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
