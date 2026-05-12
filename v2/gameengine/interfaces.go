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
)

// TriggerType categorises when an Aura or Trigger fires.
type TriggerType int

const (
	// TriggerStartOfTurn fires at the start of the owning player's action phase, before the
	// best-line search.
	TriggerStartOfTurn TriggerType = iota
	// TriggerAttackAction fires each time an attack action card resolves during the chain.
	TriggerAttackAction
	// TriggerAttack fires when ANY attack resolves — attack action card or weapon swing.
	TriggerAttack
	// TriggerEndOfTurn fires after the chain finishes resolving, before the carry snapshot.
	TriggerEndOfTurn
	// TriggerHit fires when an attack hits (post-AR-buff EffectiveAttack survives blocks).
	TriggerHit
)

// Aura is the engine's view of a persistent hook entry. Concrete aura impls expose
// what the engine needs through these methods directly.
type Aura interface {
	TriggerType() TriggerType
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
	// Clone returns a deep copy of this aura. The engine clones every aura on Reset so
	// per-permutation Count / FiredThisTurn mutations don't leak across the shared
	// priorAuras slice the seed carries.
	Clone() Aura
}

// Trigger is the engine's view of a one-shot deferred handler.
type Trigger interface {
	TriggerType() TriggerType
	// CardName is the source card's display name, used for log attribution.
	CardName() string
	// Matches narrows TriggerHit fires to a card-type predicate. Returns true when the
	// trigger has no type filter.
	Matches(types card.TypeSet) bool
	// Fire invokes the trigger's handler with a card.Trigger context built by the engine.
	Fire(g *GameEngine, l card.Logger)
	// Clone returns a deep copy of this trigger. Same per-permutation isolation rationale
	// as Aura.Clone.
	Clone() Trigger
}

// Item is the engine's view of an in-play permanent with an activated ability.
type Item interface {
	CardName() string
	CardID() ids.CardID
	Count() int
	SetCount(int)
	// Ability is the activated-ability card the chain runner enqueues each turn.
	Ability() card.Card
	// Clone returns a deep copy of this item. Same per-permutation isolation rationale
	// as Aura.Clone.
	Clone() Item
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
