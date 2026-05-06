package sim

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// Trigger holds the firing-data shared between Auras and standalone one-shot riders.
// Source identifies the card that registered the trigger; TriggerType matches the firing
// site; TypeFilter narrows the qualifying events to a card-type predicate (used by
// TriggerHit hit-riders today, nil = no filter); Handler runs when the trigger fires.
//
// Aura embeds Trigger by value (see aura.go) and adds persistent-instance state
// (Self, Count, OncePerTurn, FiredThisTurn). Standalone Triggers — the sim removes
// them after firing — are stored on TurnState.Triggers and added via AddTrigger.
//
// Handlers receive a *Trigger; aura-bound handlers recover the parent Aura via
// s.AuraFor(t) when they need Self / Count.
//
// Handlers MUST NOT call AddTrigger from within a fire walk: the firing helper
// snapshots len(s.Triggers) before iterating, so a trigger added during fire stays
// queued for the next matching event but isn't seen on the current pass.
type Trigger struct {
	// Source is the originating card, surfaced in the trace via Source.DisplayName().
	// Auras whose Self is a token leave Source nil; Source is the production source for
	// card-based triggers / auras.
	Source Card
	// TriggerType matches the firing site (TriggerEndOfTurn, TriggerAttack, …).
	TriggerType TriggerType
	// TypeFilter optionally narrows TriggerHit to a card-type predicate. nil = no filter.
	TypeFilter func(card.TypeSet) bool
	// Count is a handler-defined counter — Malefic / Runeblood read it as fires remaining
	// and decrement themselves; Runechant reads it as a copy count; standalone triggers
	// like High Striker store the per-variant payload (e.g. token count). Sim treats it
	// as opaque storage.
	Count int
	// Handler runs when TriggerType (and TypeFilter, if set) match.
	Handler TriggerHandler
}

// TriggerHandler is the business-logic callback for a Trigger. The same signature is
// used for both standalone Triggers and Aura-embedded Triggers — aura handlers call
// s.AuraFor(t) to recover the parent Aura when they need Self / Count.
type TriggerHandler func(s *TurnState, t *Trigger)
