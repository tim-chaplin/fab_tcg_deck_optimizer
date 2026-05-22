// Package gameengine owns GameEngine — the per-turn shared state threaded through every
// Card.Play — and the local interfaces (Aura, Trigger, Item, Hero) it consumes. Concrete
// impls live in dedicated leaf packages (internal/aura, internal/trigger, internal/hero, internal/token) that do
// not import gameengine; sim wires builders in via init so the engine never imports the
// concrete types directly. internal/card.GameEngine is a narrow subset of GameEngine's method
// set — the card-facing surface.
package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Aura is the engine's view of a persistent hook entry. Fire takes typed
// card.GameEngine / card.Logger arguments — the concrete *aura.Aura imports internal/card
// directly so the handler signature is end-to-end typed (no wrap closure). SourceCard
// stays `any` so token auras can return nil cleanly; consumers assert to card.Card.
// Copy uses any so State.Copy can box without referencing the concrete type. CopyInto
// is the per-perm reset fast path: when ResetForPermutationFrom has a pooled slot from
// the prior perm it hands the existing concrete value back as `any` so the implementation
// can rewrite that value's fields in place and return it boxed without allocating.
type Aura interface {
	TriggerType() triggertype.Type
	OncePerTurn() bool
	FiredThisTurn() bool
	SetFiredThisTurn(bool)
	Matches(types card.TypeSet) bool
	CardName() string
	CardID() ids.CardID
	SourceCard() any
	Count() int
	SetCount(int)
	DecrementCount() int
	Fire(engine card.GameEngine, logger card.Logger)
	Copy() any
	CopyInto(dst any) any
}

// EphemeralTrigger is the engine's view of a one-shot deferred handler. Like Aura, Fire
// takes typed card.GameEngine / card.Logger arguments; Matches takes the firing card's
// TypeSet directly.
type EphemeralTrigger interface {
	TriggerType() triggertype.Type
	CardName() string
	Matches(types card.TypeSet) bool
	Fire(engine card.GameEngine, logger card.Logger)
}

// Item is the engine's view of an in-play permanent. Token items carry an activated
// ability (returned as `any`; callers assert back to card.Card); card-sourced items carry
// a trigger FireTriggers dispatches through, mirroring Aura. A token item leaves the
// trigger fields zero so its TriggerType never matches a firing event.
type Item interface {
	CardName() string
	CardID() ids.CardID
	Count() int
	SetCount(int)
	Ability() any
	SourceCard() any
	TriggerType() triggertype.Type
	Matches(types card.TypeSet) bool
	OncePerTurn() bool
	FiredThisTurn() bool
	SetFiredThisTurn(bool)
	Fire(engine card.GameEngine, logger card.Logger)
	Copy() any
}

// Hero is the engine's view of the active hero. OnCardPlayed takes hero.GameEngine and
// hero.Logger — narrow surfaces concrete heroes consume directly so internal/hero doesn't
// reference internal/card.GameEngine / internal/card.Logger. *GameEngine satisfies hero.GameEngine
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
