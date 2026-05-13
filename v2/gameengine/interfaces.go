// Package gameengine owns GameEngine — the per-turn shared state threaded through every
// Card.Play — and the local interfaces (Aura, Trigger, Item, Hero) it consumes. Concrete
// impls live in dedicated leaf packages (v2/aura, v2/trigger, v2/hero, v2/token) that do
// not import gameengine; sim wires builders in via init so the engine never imports the
// concrete types directly. v2/card.GameEngine is a narrow subset of GameEngine's method
// set — the card-facing surface.
package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Aura is the engine's view of a persistent hook entry. The Fire signature uses
// card.GameEngine so concrete aura impls satisfy this without importing gameengine —
// *GameEngine satisfies card.GameEngine. Copy returns any for the same reason; the engine
// type-asserts back to Aura when consuming.
type Aura interface {
	TriggerType() triggertype.Type
	OncePerTurn() bool
	FiredThisTurn() bool
	SetFiredThisTurn(bool)
	CardName() string
	CardID() ids.CardID
	// SourceCard returns the originating card.Card for card-backed auras, or nil for token
	// auras. DestroyAura uses this to route the source card into the graveyard; sim's
	// start-of-turn listing reads it to attribute aura fires to their source card.
	SourceCard() card.Card
	Count() int
	SetCount(int)
	DecrementCount() int
	Fire(ge card.GameEngine, l card.Logger)
	// Copy returns a deep copy boxed as any so concrete impls satisfy without importing
	// gameengine. State.Copy asserts the result back to Aura.
	Copy() any
}

// Trigger is the engine's view of a one-shot deferred handler.
type Trigger interface {
	TriggerType() triggertype.Type
	CardName() string
	Matches(types card.TypeSet) bool
	Fire(ge card.GameEngine, l card.Logger)
}

// Item is the engine's view of an in-play permanent with an activated ability.
type Item interface {
	CardName() string
	CardID() ids.CardID
	Count() int
	SetCount(int)
	Ability() card.Card
	// Copy returns a deep copy boxed as any. State.Copy asserts the result back to Item.
	Copy() any
}

// Hero is the engine's view of the active hero. OnCardPlayed takes card.GameEngine so
// concrete heroes don't have to import gameengine — *GameEngine satisfies card.GameEngine.
type Hero interface {
	ID() ids.HeroID
	Name() string
	Class() card.CardType
	Types() card.TypeSet
	Intelligence() int
	OnCardPlayed(played card.Card, ge card.GameEngine, l card.Logger) int
	Opt(cards []card.Card) (top, bottom []card.Card)
}

// LowerHealthWanter is a Hero marker. Heroes whose strategy revolves around staying at
// lower {h} than the opponent opt in; cards with "less {h} than an opposing hero" riders
// assume the clause fires for these heroes and never for anyone else.
type LowerHealthWanter interface {
	WantsLowerHealth()
}
