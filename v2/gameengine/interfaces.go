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

// Aura is the engine's view of a persistent hook entry. Fire takes typed
// card.GameEngine / card.Logger arguments — the concrete *aura.Aura imports v2/card
// directly so the handler signature is end-to-end typed (no wrap closure). SourceCard
// stays `any` so token auras can return nil cleanly; consumers assert to card.Card.
// Copy uses any so State.Copy can box without referencing the concrete type.
type Aura interface {
	TriggerType() triggertype.Type
	OncePerTurn() bool
	FiredThisTurn() bool
	SetFiredThisTurn(bool)
	CardName() string
	CardID() ids.CardID
	SourceCard() any
	Count() int
	SetCount(int)
	DecrementCount() int
	Fire(engine card.GameEngine, logger card.Logger)
	Copy() any
}

// Trigger is the engine's view of a one-shot deferred handler. Like Aura, Fire takes
// typed card.GameEngine / card.Logger arguments; Matches takes the firing card's
// TypeSet directly.
type Trigger interface {
	TriggerType() triggertype.Type
	CardName() string
	Matches(types card.TypeSet) bool
	Fire(engine card.GameEngine, logger card.Logger)
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
