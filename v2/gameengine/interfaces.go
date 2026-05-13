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
	"github.com/tim-chaplin/fab-deck-optimizer/v2/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Aura is the engine's view of a persistent hook entry. Engine, logger, and source-card
// values flow through as `any` so concrete impls (in v2/aura) can satisfy this without
// importing gameengine OR v2/card — the engine asserts the boxed values back to the typed
// interfaces it uses. Copy uses any for the same reason; State.Copy asserts back to Aura.
type Aura interface {
	TriggerType() triggertype.Type
	OncePerTurn() bool
	FiredThisTurn() bool
	SetFiredThisTurn(bool)
	CardName() string
	CardID() ids.CardID
	// SourceCard returns the originating card boxed as any for card-backed auras, or nil
	// for token auras. DestroyAura asserts it back to card.Card to route into graveyard;
	// sim's start-of-turn listing reads it to attribute aura fires to their source card.
	SourceCard() any
	Count() int
	SetCount(int)
	DecrementCount() int
	Fire(engine, logger any)
	Copy() any
}

// Trigger is the engine's view of a one-shot deferred handler. As with Aura, engine /
// logger / type-set values flow through as `any`.
type Trigger interface {
	TriggerType() triggertype.Type
	CardName() string
	Matches(types any) bool
	Fire(engine, logger any)
}

// Item is the engine's view of an in-play permanent with an activated ability. Ability
// is returned as `any` so concrete Item impls (in v2/token) can satisfy this without
// importing v2/card — engine and chain-runner callers assert back to card.Card.
type Item interface {
	CardName() string
	CardID() ids.CardID
	Count() int
	SetCount(int)
	Ability() any
	Copy() any
}

// Hero is the engine's view of the active hero. OnCardPlayed takes hero.GameEngine and
// hero.Logger — narrow surfaces concrete heroes consume directly so v2/hero doesn't
// reference v2/card.GameEngine / v2/card.Logger. *GameEngine satisfies hero.GameEngine
// structurally.
type Hero interface {
	ID() ids.HeroID
	Name() string
	Class() card.CardType
	Types() card.TypeSet
	Intelligence() int
	OnCardPlayed(played card.Card, ge hero.GameEngine, l hero.Logger) int
	Opt(cards []card.Card) (top, bottom []card.Card)
}

// LowerHealthWanter is a Hero marker. Heroes whose strategy revolves around staying at
// lower {h} than the opponent opt in; cards with "less {h} than an opposing hero" riders
// assume the clause fires for these heroes and never for anyone else.
type LowerHealthWanter interface {
	WantsLowerHealth()
}
