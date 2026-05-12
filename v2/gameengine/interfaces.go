// Package gameengine owns GameEngine — the per-turn shared state threaded through every
// Card.Play — and the local interfaces (Aura, Trigger, Item, Hero) it consumes. Concrete
// impls of those interfaces live in the caller (today internal/sim; eventually their own
// v2/ subpackages). v2/card.GameEngine is a narrow subset of GameEngine's method set — the
// card-facing surface; sim imports this package directly and gets the full API including
// Reset / Snapshot / Fire* / SetHero.
package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Aura is the engine's view of a persistent hook entry. Concrete aura impls expose
// what the engine needs through these methods directly.
type Aura interface {
	TriggerType() triggertype.Type
	OncePerTurn() bool
	FiredThisTurn() bool
	SetFiredThisTurn(bool)
	// CardName is the originating card or token's display name, used for log attribution.
	CardName() string
	// CardID is the originating card's registry ID; tokens carry their reserved engine
	// token-ID (RunechantTokenID, …) so cache keys distinguish each kind without an
	// extra discriminator.
	CardID() ids.CardID
	// SourceCard returns the originating card.Card for card-backed auras, or nil for
	// token auras. The eval-loop start-of-turn aura listing uses it to attribute aura
	// fires back to their source card in the printout.
	SourceCard() card.Card
	Count() int
	SetCount(int)
	DecrementCount() int
	// Fire invokes the aura's handler with a card.Aura context built by the engine.
	Fire(g *GameEngine, l card.Logger)
	// OnDestroy is called by the engine when the aura is destroyed with addToGraveyard=true.
	// Card-backed auras push their source card into the graveyard via g.AppendGraveyard;
	// token auras no-op.
	OnDestroy(g *GameEngine)
	// Copy returns a deep copy. GameEngine.Copy walks the aura list and stores the
	// result so per-permutation Count / FiredThisTurn mutations don't leak across copies.
	Copy() Aura
}

// Trigger is the engine's view of a one-shot deferred handler.
type Trigger interface {
	TriggerType() triggertype.Type
	// CardName is the source card's display name, used for log attribution.
	CardName() string
	// Matches narrows triggertype.Hit fires to a card-type predicate. Returns true when the
	// trigger has no type filter.
	Matches(types card.TypeSet) bool
	// Fire invokes the trigger's handler with a card.Trigger context built by the engine.
	Fire(g *GameEngine, l card.Logger)
}

// Item is the engine's view of an in-play permanent with an activated ability.
type Item interface {
	CardName() string
	CardID() ids.CardID
	Count() int
	SetCount(int)
	// Ability is the activated-ability card the chain runner enqueues each turn.
	Ability() card.Card
	// Copy returns a deep copy. GameEngine.Copy walks the item list and stores the
	// result so per-permutation Count mutations don't leak across copies.
	Copy() Item
}

// Hero is the engine's view of the active hero. Concrete heroes (e.g. internal/heroes/*)
// satisfy this. The engine reads Class / Types on demand (Universal cards fold the active
// hero's class into their own Types) and calls OnCardPlayed before each chain step resolves
// so hero abilities fire ahead of the card itself.
type Hero interface {
	ID() ids.HeroID
	Name() string
	Class() card.CardType
	Types() card.TypeSet
	Intelligence() int
	OnCardPlayed(played card.Card, g *GameEngine, l card.Logger) int
	Opt(cards []card.Card) (top, bottom []card.Card)
}

// LowerHealthWanter is a Hero marker. Heroes whose strategy revolves around staying at
// lower {h} than the opponent opt in; cards with "less {h} than an opposing hero" riders
// assume the clause fires for these heroes and never for anyone else.
type LowerHealthWanter interface {
	WantsLowerHealth()
}
