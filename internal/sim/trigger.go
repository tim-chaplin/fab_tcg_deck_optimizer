package sim

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// Trigger holds the firing-data shared between Auras and standalone one-shot riders.
// Source identifies the card that registered the trigger; TriggerType matches the firing
// site; TypeFilter narrows the qualifying events to a card-type predicate (used by
// TriggerHit hit-riders today, nil = no filter); Handler runs when the trigger fires.
//
// Aura embeds Trigger by value (see aura.go) and adds persistent-instance state
// (Self, OncePerTurn, FiredThisTurn, Count). Standalone Triggers — the sim removes them
// after firing — are stored on TurnState.Triggers and added via AddTrigger.
//
// Handlers receive (s, t, a): for aura fires a is the firing aura, for standalone
// trigger fires a is nil. Handlers that mutate s.Auras (e.g. via CreateRunechants, which
// appends a Runechant aura) MUST read state through a BEFORE the mutating call — a slice
// realloc invalidates the parameter.
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
	// Handler runs when TriggerType (and TypeFilter, if set) match.
	Handler TriggerHandler
}

// TriggerHandler is the business-logic callback for a Trigger. Aura fires pass the
// firing aura as a; standalone trigger fires pass nil. l is the cards-facing log sink
// the chain runner threaded through Play; handlers emit their pre/post-trigger lines
// through it so a nil l silently elides logging during the find-best pass.
type TriggerHandler func(s *TurnState, l Logger, t *Trigger, a *Aura)
