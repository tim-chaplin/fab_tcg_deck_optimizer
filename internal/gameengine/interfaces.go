// Package gameengine owns GameEngine — the per-turn shared state threaded through every
// Card.Play — and the local interfaces (Aura, Trigger, Item, Hero) it consumes. Concrete
// impls live in dedicated leaf packages (internal/aura, internal/trigger, internal/hero, internal/token) that do
// not import gameengine; sim wires builders in via init so the engine never imports the
// concrete types directly. internal/card.GameEngine is a narrow subset of GameEngine's method
// set — the card-facing surface.
package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/trigger"
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
	trigger.Hook
	CardName() string
	CardID() ids.CardID
	SourceCard() any
	Count() int
	SetCount(int)
	DecrementCount() int
	Copy() any
	CopyInto(dst any) any
	// Clear zeros source / count / fire-closure so the underlying pool slot reads as
	// free for the GameState aura pool's next-free-slot scan. Called by DestroyAura
	// after the slot leaves play.
	Clear()
}

// EphemeralTrigger is the engine's view of a one-shot deferred handler. The trigger.Hook
// trio (OncePerTurn / FiredThisTurn / SetFiredThisTurn) is part of the shared dispatch
// surface fireHooks walks — ephemeral triggers never gate on it (they fire once and are
// removed), so OncePerTurn is always false.
type EphemeralTrigger interface {
	trigger.Hook
	CardName() string
}

// Item is the engine's view of an in-play permanent. Token items carry an activated
// ability (returned as `any`; callers assert back to card.Card); card-sourced items carry
// a trigger FireTriggers dispatches through, mirroring Aura.
type Item interface {
	trigger.Hook
	CardName() string
	CardID() ids.CardID
	Count() int
	SetCount(int)
	Ability() any
	SourceCard() any
	Copy() any
}

// Weapon is the engine's view of an equipped weapon — the mutable object built when the
// platonic weapon card is equipped at game start. Mirrors Aura / Item: it carries a trigger
// (no weapon registers a handler yet), a counter total, and the source card sent to the
// graveyard on destroy. Name / Hands / Ability surface the platonic card's weapon attributes
// the attack-turn runner reads — Ability is the swing card enqueued each turn. Copy uses any
// so State.Copy can deep-copy without referencing the concrete *weapon.Weapon. CopyInto is
// the per-perm reset fast path: it rewrites a pooled slot's existing concrete value in place
// (the loadout is identical across perms) so the hot path allocates nothing.
type Weapon interface {
	trigger.Hook
	CardName() string
	CardID() ids.CardID
	Count() int
	SetCount(int)
	SourceCard() any
	Name() string
	Hands() int
	Ability() card.Card
	Copy() any
	CopyInto(dst any) any
}

// Hero is the engine's view of the active hero. The trigger-dispatch surface mirrors
// the methods Aura and Item expose through their embedded trigger — FireTriggers walks
// the hero through the same path.
type Hero interface {
	trigger.Hook
	ID() ids.HeroID
	Name() string
	Class() card.CardType
	Types() card.TypeSet
	Intelligence() int
	Opt(cards []card.Card) (top, bottom []card.Card)
}

// LowerHealthWanter is a Hero marker. Heroes whose strategy revolves around staying at
// lower {h} than the opponent opt in; cards with "less {h} than an opposing hero" riders
// assume the clause fires for these heroes and never for anyone else.
type LowerHealthWanter interface {
	WantsLowerHealth()
}
