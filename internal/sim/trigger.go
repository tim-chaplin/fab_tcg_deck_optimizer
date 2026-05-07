package sim

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// Trigger holds the firing-data shared between Auras and standalone one-shot riders.
// Source identifies the card that registered the trigger; TriggerType matches the firing
// site; TypeFilter narrows the qualifying events to a card-type predicate (used by
// TriggerHit hit-riders today, nil = no filter); Handler runs when the trigger fires.
//
// Aura embeds Trigger by value (see aura.go) and adds persistent-instance state
// (Self, OncePerTurn, FiredThisTurn). Standalone Triggers — the sim removes them after
// firing — are stored on TurnState.Triggers and added via AddTrigger.
//
// Handlers receive a *Trigger; aura-bound handlers reach the parent Aura through
// t.Aura. The fire loop sets t.Aura just before invoking the handler, so the pointer
// is valid for the duration of the handler — handlers that mutate s.Auras (e.g. via
// CreateRunechants, which appends a Runechant aura) MUST read aura state via t.Aura
// BEFORE the mutating call, since a slice realloc would invalidate the back-pointer.
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
	// Aura points at the parent Aura when this Trigger is an aura's embedded firing-data;
	// nil for standalone Triggers. The fire loop sets it before invoking Handler.
	Aura *Aura
}

// TriggerHandler is the business-logic callback for a Trigger. The same signature is
// used for both standalone Triggers and Aura-embedded Triggers — aura handlers reach
// the parent Aura via t.Aura.
type TriggerHandler func(s *TurnState, t *Trigger)
