// Package aura owns the concrete Aura type — a persistent hook entry that fires at a
// scheduled trigger type (start of turn, attack, attack action, end of turn). Card-backed
// and token-backed auras share the same struct; SourceCard distinguishes them.
//
// The package is a generic mechanism — it defines its own narrow interfaces (Card,
// GameEngine in interfaces.go) and does not import v2/card or v2/gameengine. Stored
// handlers are typed as Handler — engine, logger, and the per-fire ctx all flow through
// as `any` so aura never depends on the consumer's broader surfaces. Each consumer defines
// its own interface (e.g. v2/card.Aura for the per-fire ctx) and asserts at the handler
// boundary; the factory layer (sim/init.go) wraps typed user handlers into Handler closures.
package aura

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Handler is the per-aura stored handler. Every argument is typed as `any` — engine,
// logger, and the per-fire ctx — so aura doesn't depend on the consumer's interfaces. The
// factory layer wraps typed user handlers into closures of this shape; handler bodies
// type-assert to whatever richer interfaces they need.
type Handler func(engine, logger, ctx any)

// Aura is the concrete entry the engine stores in its persistent hook list. source is
// non-nil for card-backed auras; tokenName / tokenID are populated for token auras.
type Aura struct {
	triggerType   triggertype.Type
	fire          Handler
	source        Card
	tokenName     string
	tokenID       ids.CardID
	count         int
	oncePerTurn   bool
	firedThisTurn bool
}

// NewFromCard builds a card-backed aura. source is the originating card — typically the
// Card field of a CardState. SourceCard surfaces it back so engines can route it into the
// graveyard on destroy.
func NewFromCard(source Card, tt triggertype.Type, fire Handler, count int, oncePerTurn bool) *Aura {
	return &Aura{
		triggerType: tt,
		fire:        fire,
		source:      source,
		count:       count,
		oncePerTurn: oncePerTurn,
	}
}

// NewFromToken builds a token aura — no originating card. CardName returns the supplied
// name; CardID returns tokenID so cache keys distinguish each token kind.
func NewFromToken(name string, tokenID ids.CardID, tt triggertype.Type, fire Handler, count int) *Aura {
	return &Aura{
		triggerType: tt,
		fire:        fire,
		tokenName:   name,
		tokenID:     tokenID,
		count:       count,
	}
}

func (a *Aura) TriggerType() triggertype.Type { return a.triggerType }
func (a *Aura) OncePerTurn() bool             { return a.oncePerTurn }
func (a *Aura) FiredThisTurn() bool           { return a.firedThisTurn }
func (a *Aura) SetFiredThisTurn(v bool)       { a.firedThisTurn = v }

func (a *Aura) CardName() string {
	if a.source != nil {
		return a.source.DisplayName()
	}
	return a.tokenName
}

func (a *Aura) CardID() ids.CardID {
	if a.source != nil {
		return a.source.ID()
	}
	return a.tokenID
}

// SourceCard returns the originating source as `any` — nil for token-backed auras, the
// stored Card boxed otherwise. The engine asserts the value back to its richer card type
// (typically v2/card.Card) when routing into graveyard.
func (a *Aura) SourceCard() any {
	if a.source == nil {
		return nil
	}
	return a.source
}
func (a *Aura) Count() int     { return a.count }
func (a *Aura) SetCount(n int) { a.count = n }
func (a *Aura) DecrementCount() int {
	a.count--
	return a.count
}

// Fire invokes the stored handler. engine, logger, and the per-fire ctx pass through as
// `any` — the handler (typically wrapped at the factory layer) asserts them to its
// preferred typed interfaces.
func (a *Aura) Fire(engine, logger any) {
	a.fire(engine, logger, &ctx{a: a, engine: engine})
}

// Copy returns a deep copy boxed as any so the gameengine.Aura interface declaration can
// avoid referencing its own type — letting concrete impls satisfy without importing.
func (a *Aura) Copy() any {
	out := *a
	return &out
}

// ctx is the concrete per-fire value Fire passes to the handler as `any`. Consumers
// define their own interface (e.g. v2/card.Aura) and assert this value to it; the method
// set here is what every such interface can rely on. Destroy routes back to the engine's
// DestroyAura through the package's local GameEngine interface.
type ctx struct {
	a      *Aura
	engine any
}

func (c *ctx) Count() int          { return c.a.count }
func (c *ctx) DecrementCount() int { c.a.count--; return c.a.count }
func (c *ctx) CardName() string    { return c.a.CardName() }
func (c *ctx) CardID() ids.CardID  { return c.a.CardID() }
func (c *ctx) Destroy(addToGraveyard bool) {
	c.engine.(GameEngine).DestroyAura(addToGraveyard)
}
