// Package aura owns the concrete Aura type — a persistent hook entry that fires at a
// scheduled trigger type (start of turn, attack, attack action, end of turn). Card-backed
// and token-backed auras share the same struct; SourceCard distinguishes them.
//
// Handler is the aura-domain function shape; the package's own ctx struct implements
// card.Aura so handler bodies receive the typed surface directly without an outer
// wrapping layer.
package aura

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Handler is the typed aura handler signature. Lives in v2/aura (not v2/card) because
// the shape is meaningful only to aura: it's the signature stored on each Aura and called
// at every Fire. card.GameEngine.CreateStartOfTurnAura et al. inline the function type in
// their parameter declarations so v2/card doesn't need to import v2/aura.
type Handler func(card.GameEngine, card.Logger, card.Aura)

// Aura is the concrete entry the engine stores in its persistent hook list. source is
// non-nil for card-backed auras; tokenName / tokenID are populated for token auras.
type Aura struct {
	triggerType   triggertype.Type
	fire          Handler
	source        card.Card
	tokenName     string
	tokenID       ids.CardID
	count         int
	oncePerTurn   bool
	firedThisTurn bool
}

// NewFromCard builds a card-backed aura. source is the originating card — typically the
// Card field of a CardState. SourceCard surfaces it back so engines can route it into the
// graveyard on destroy.
func NewFromCard(source card.Card, tt triggertype.Type, fire Handler, count int, oncePerTurn bool) *Aura {
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
// stored card.Card otherwise. Boxed to `any` so the gameengine.Aura interface declaration
// can avoid referencing card.Card directly; consumers assert back.
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

// Fire invokes the stored handler. The handler signature is typed against v2/card so the
// engine and logger arguments pass straight through without a per-fire assertion.
func (a *Aura) Fire(engine card.GameEngine, logger card.Logger) {
	a.fire(engine, logger, &ctx{a: a, engine: engine})
}

// Copy returns a deep copy boxed as any so the gameengine.Aura interface declaration can
// avoid referencing its own type — letting concrete impls satisfy without importing.
func (a *Aura) Copy() any {
	out := *a
	return &out
}

// ctx is the per-fire value Fire passes to the handler. Implements card.Aura so handler
// bodies receive the typed surface directly; Destroy routes back to the engine's
// DestroyAura through the stored card.GameEngine reference.
type ctx struct {
	a      *Aura
	engine card.GameEngine
}

// Compile-time check that ctx satisfies card.Aura.
var _ card.Aura = (*ctx)(nil)

func (c *ctx) Count() int          { return c.a.count }
func (c *ctx) DecrementCount() int { c.a.count--; return c.a.count }
func (c *ctx) CardName() string    { return c.a.CardName() }
func (c *ctx) CardID() ids.CardID  { return c.a.CardID() }
func (c *ctx) Destroy(addToGraveyard bool) {
	c.engine.DestroyAura(addToGraveyard)
}
