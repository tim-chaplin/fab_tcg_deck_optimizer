// Package aura owns the concrete Aura type — a persistent hook entry that fires at a
// scheduled trigger type (start of turn, attack, attack action, end of turn). Card-backed
// and token-backed auras share the same struct; SourceCard distinguishes them. Satisfies
// gameengine.Aura structurally without importing gameengine — the Fire signature uses
// card.GameEngine, Copy returns any, and the ctx-side Destroy hop type-asserts the engine
// to a local interface.
//
// FaB's two aura-flavored tokens (Runechant, Ponder) are also built here in tokens.go —
// they share the concrete Aura type with card-backed auras.
package aura

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Aura is the concrete entry the engine stores in its persistent hook list. source is
// non-nil for card-backed auras; tokenName / tokenID are populated for token auras.
type Aura struct {
	triggerType   triggertype.Type
	handler       card.AuraHandler
	source        card.Card
	tokenName     string
	tokenID       ids.CardID
	count         int
	oncePerTurn   bool
	firedThisTurn bool
}

// NewCard builds a card-backed aura — SourceCard surfaces self.Card. When DestroyAura
// fires with addToGraveyard=true the engine pushes SourceCard into the graveyard.
func NewCard(self *card.CardState, tt triggertype.Type, handler card.AuraHandler, count int, oncePerTurn bool) *Aura {
	return &Aura{
		triggerType: tt,
		handler:     handler,
		source:      self.Card,
		count:       count,
		oncePerTurn: oncePerTurn,
	}
}

// NewToken builds a token aura — there is no underlying card. CardName returns the supplied
// name; CardID returns tokenID so cache keys distinguish each token kind.
func NewToken(name string, tokenID ids.CardID, tt triggertype.Type, handler card.AuraHandler, count int) *Aura {
	return &Aura{
		triggerType: tt,
		handler:     handler,
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

func (a *Aura) SourceCard() card.Card { return a.source }
func (a *Aura) Count() int            { return a.count }
func (a *Aura) SetCount(n int)        { a.count = n }
func (a *Aura) DecrementCount() int   { a.count--; return a.count }

// Fire invokes the aura's registered handler with a card.Aura context. The engine satisfies
// card.GameEngine so the engine value flows through unchanged.
func (a *Aura) Fire(ge card.GameEngine, l card.Logger) {
	ctx := &auraCtx{a: a, ge: ge}
	a.handler(ge, l, ctx)
}

// Copy returns a deep copy boxed as any so the gameengine.Aura interface declaration can
// avoid referencing its own type — letting concrete impls satisfy without importing.
func (a *Aura) Copy() any {
	out := *a
	return &out
}

// destroyableEngine is the slice of the engine surface the aura context needs to wire
// Destroy back through. *gameengine.GameEngine satisfies it structurally.
type destroyableEngine interface {
	DestroyAura(addToGraveyard bool)
}

// auraCtx adapts (Aura, engine) into the card.Aura surface handlers see. Cards reach for
// Count / DecrementCount / CardName / CardID through this and request destruction via
// Destroy — the latter routes back to the engine's DestroyAura.
type auraCtx struct {
	a  *Aura
	ge card.GameEngine
}

func (c *auraCtx) Count() int          { return c.a.count }
func (c *auraCtx) DecrementCount() int { c.a.count--; return c.a.count }
func (c *auraCtx) CardName() string    { return c.a.CardName() }
func (c *auraCtx) CardID() ids.CardID  { return c.a.CardID() }
func (c *auraCtx) Destroy(addToGraveyard bool) {
	c.ge.(destroyableEngine).DestroyAura(addToGraveyard)
}
