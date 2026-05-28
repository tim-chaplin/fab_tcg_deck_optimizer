// Package trigger owns the shared trigger machinery. Trigger[T] is the embeddable core
// every triggered entry carries — one-shot ephemeral triggers, auras, and items — pairing
// a firing event with a typed handler, the source identity, and an optional type filter.
// EphemeralTrigger is the one-shot kind: the engine fires it on the next matching event
// and then drops it from the queue.
package trigger

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// TypeFilter narrows the firing site to a card-type predicate. nil means any matching
// event qualifies.
type TypeFilter func(card.TypeSet) bool

// Hook is the trigger-dispatch surface every triggered entry exposes — auras, items,
// one-shot triggers, and the hero. Concrete entries get it for free by embedding
// Trigger[T] (whose pointer-receiver methods cover everything except Fire); each entry
// type supplies its own Fire that calls Invoke with itself as the typed receiver. The
// engine's dispatch loop walks entries through this interface and threads both the
// firing trigger type (a single bit, even when the registration's TriggerType
// subscribes to several via bitmask OR) and the triggering card (the CardState whose
// resolution raised the event, or nil for turn-boundary events) so handlers can
// dispatch on which event invoked them and read in-attack-turn state on the source.
type Hook interface {
	TriggerType() triggertype.Type
	OncePerTurn() bool
	FiredThisTurn() bool
	SetFiredThisTurn(bool)
	Matches(types card.TypeSet) bool
	Fire(engine card.GameEngine, logger card.Logger, triggeringCard *card.CardState, firingType triggertype.Type)
}

// Trigger is the embeddable core a triggered entry carries. T is the concrete surface the
// handler receives — card.EphemeralTrigger, card.Aura, card.Item — so handlers stay typed
// with no assertion. The embedding type supplies a Fire method that calls Invoke with
// itself as the typed receiver and the firing trigger type so multi-trigger subscribers
// can dispatch.
//
// triggerType is the subscription bitmask — handlers may subscribe to several events by
// OR-ing types (e.g. DamageAboutToBeTaken | EndOfTurn). oncePerTurn caps the entry to
// one fire per turn across all subscribed events; firedThisTurn is the gate, set after
// a fire and re-armed at the turn boundary.
type Trigger[T any] struct {
	triggerType   triggertype.Type
	fire          func(card.GameEngine, card.Logger, T, *card.CardState, triggertype.Type)
	source        card.Card
	tokenName     string
	tokenID       ids.CardID
	typeFilter    TypeFilter
	oncePerTurn   bool
	firedThisTurn bool
}

// FromCard builds a card-sourced Trigger core. oncePerTurn caps it to one fire per turn;
// typeFilter narrows the firing site, pass nil for no filter. The fire handler receives
// the triggering card (the CardState whose resolution raised the event, or nil for
// turn-boundary events) and the firing trigger type so multi-trigger subscribers can
// dispatch.
func FromCard[T any](source card.Card, tt triggertype.Type, fire func(card.GameEngine, card.Logger, T, *card.CardState, triggertype.Type), oncePerTurn bool, typeFilter TypeFilter) Trigger[T] {
	return Trigger[T]{triggerType: tt, fire: fire, source: source, oncePerTurn: oncePerTurn, typeFilter: typeFilter}
}

// FromToken builds a token-sourced Trigger core — no originating card. CardName returns
// the supplied name; CardID returns tokenID so cache keys distinguish each token kind.
func FromToken[T any](name string, tokenID ids.CardID, tt triggertype.Type, fire func(card.GameEngine, card.Logger, T, *card.CardState, triggertype.Type), oncePerTurn bool, typeFilter TypeFilter) Trigger[T] {
	return Trigger[T]{triggerType: tt, fire: fire, tokenName: name, tokenID: tokenID, oncePerTurn: oncePerTurn, typeFilter: typeFilter}
}

// FromHero builds a Trigger core for a hero ability — no card or token identity. The
// embedding Hero type owns its own ID / Name, so the source / token slots stay zero;
// CardName returns "" and CardID returns 0 on a hero-sourced trigger.
func FromHero[T any](tt triggertype.Type, fire func(card.GameEngine, card.Logger, T, *card.CardState, triggertype.Type), oncePerTurn bool, typeFilter TypeFilter) Trigger[T] {
	return Trigger[T]{triggerType: tt, fire: fire, oncePerTurn: oncePerTurn, typeFilter: typeFilter}
}

func (t *Trigger[T]) TriggerType() triggertype.Type { return t.triggerType }
func (t *Trigger[T]) OncePerTurn() bool             { return t.oncePerTurn }
func (t *Trigger[T]) FiredThisTurn() bool           { return t.firedThisTurn }
func (t *Trigger[T]) SetFiredThisTurn(v bool)       { t.firedThisTurn = v }

// Matches reports whether the type filter accepts the firing event's type set.
func (t *Trigger[T]) Matches(types card.TypeSet) bool {
	if t.typeFilter == nil {
		return true
	}
	return t.typeFilter(types)
}

// CardName is the originating card or token's display name — used for log attribution.
func (t *Trigger[T]) CardName() string {
	if t.source != nil {
		return t.source.DisplayName()
	}
	return t.tokenName
}

// CardID is the originating card's registry ID, or the token ID for token-sourced entries.
func (t *Trigger[T]) CardID() ids.CardID {
	if t.source != nil {
		return t.source.ID()
	}
	return t.tokenID
}

// SourceCard returns the originating card boxed as any, or nil for token-sourced entries.
func (t *Trigger[T]) SourceCard() any {
	if t.source == nil {
		return nil
	}
	return t.source
}

// Invoke runs the handler, passing self as the typed receiver plus the triggering card
// (or nil for turn-boundary events) and the firing trigger type so multi-trigger
// handlers can dispatch on which event invoked them.
func (t *Trigger[T]) Invoke(engine card.GameEngine, logger card.Logger, self T, triggeringCard *card.CardState, firingType triggertype.Type) {
	t.fire(engine, logger, self, triggeringCard, firingType)
}

// SetFromCard rewrites t in place as a fresh card-sourced trigger. firedThisTurn
// re-arms so a reused pool slot fires cleanly. Token-identity fields zero out — a slot
// returned to the card-aura pool stops claiming any token name / ID.
func (t *Trigger[T]) SetFromCard(source card.Card, tt triggertype.Type, fire func(card.GameEngine, card.Logger, T, *card.CardState, triggertype.Type), oncePerTurn bool, typeFilter TypeFilter) {
	t.triggerType = tt
	t.fire = fire
	t.source = source
	t.tokenName = ""
	t.tokenID = 0
	t.typeFilter = typeFilter
	t.oncePerTurn = oncePerTurn
	t.firedThisTurn = false
}

// Clear zeros every field so the slot reads as "free" for pool reuse. SourceCard()
// returning nil is the canonical "dead slot" test the gameengine aura pool keys on.
func (t *Trigger[T]) Clear() {
	var zero Trigger[T]
	*t = zero
}
