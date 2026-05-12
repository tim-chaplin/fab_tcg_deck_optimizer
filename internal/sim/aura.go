package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Aura is the sim concrete impl of gameengine.Aura. Card-driven Create*Aura methods on
// *gameengine.GameEngine build one of these via sim's registered builders; sim's token-aura
// factories (NewRunechantAura, NewPonderAura, …) build via NewTokenAura. Both flavours
// satisfy gameengine.Aura.
type Aura struct {
	triggerType   gameengine.TriggerType
	handler       card.AuraHandler
	source        card.Card // nil for token auras
	tokenName     string    // "Runechant" / "Ponder" / "" for card auras
	tokenID       ids.CardID
	count         int
	oncePerTurn   bool
	firedThisTurn bool
}

// NewCardAura builds a card-backed aura — source is self.Card. CardName / CardID surface
// self.Card's DisplayName / ID. On destroy with addToGraveyard=true the source card lands
// in the graveyard.
func NewCardAura(self *card.CardState, tt gameengine.TriggerType, handler card.AuraHandler, count int, oncePerTurn bool) *Aura {
	return &Aura{
		triggerType: tt,
		handler:     handler,
		source:      self.Card,
		count:       count,
		oncePerTurn: oncePerTurn,
	}
}

// NewTokenAura builds a token aura — there is no underlying card. CardName returns the
// supplied name (e.g. "Runechant"); CardID returns the supplied tokenID so cache keys
// distinguish each token kind without an extra discriminator. On destroy the aura no-ops
// in the graveyard pass (tokens don't head to graveyard).
func NewTokenAura(name string, tokenID ids.CardID, tt gameengine.TriggerType, handler card.AuraHandler, count int) *Aura {
	return &Aura{
		triggerType: tt,
		handler:     handler,
		tokenName:   name,
		tokenID:     tokenID,
		count:       count,
	}
}

// auraSliceAsEngine converts a []*Aura to []gameengine.Aura for engine-API call sites
// (PermutationSeed.Auras, Spec.Auras). The underlying entries are unchanged — each *Aura
// satisfies gameengine.Aura — so this is just a per-element box.
func auraSliceAsEngine(src []*Aura) []gameengine.Aura {
	if len(src) == 0 {
		return nil
	}
	out := make([]gameengine.Aura, len(src))
	for i, a := range src {
		out[i] = a
	}
	return out
}

// auraSliceFromEngine type-asserts every entry of an engine-returned []gameengine.Aura
// back to *Aura. Sim's runtime is the only registered Aura impl, so each assertion is
// guaranteed to succeed in production.
func auraSliceFromEngine(src []gameengine.Aura) []*Aura {
	if len(src) == 0 {
		return nil
	}
	out := make([]*Aura, len(src))
	for i, a := range src {
		out[i] = a.(*Aura)
	}
	return out
}

func (a *Aura) TriggerType() gameengine.TriggerType { return a.triggerType }
func (a *Aura) OncePerTurn() bool                   { return a.oncePerTurn }
func (a *Aura) FiredThisTurn() bool                 { return a.firedThisTurn }
func (a *Aura) SetFiredThisTurn(v bool)             { a.firedThisTurn = v }

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

func (a *Aura) Fire(g *gameengine.GameEngine, l card.Logger) {
	ctx := &auraCtx{a: a, g: g}
	a.handler(g, l, ctx)
}

func (a *Aura) OnDestroy(g *gameengine.GameEngine) {
	if a.source != nil {
		g.AppendGraveyard(a.source)
	}
}

func (a *Aura) Clone() gameengine.Aura {
	out := *a
	return &out
}

// auraCtx is the adapter handed to each aura handler — satisfies card.Aura so cards
// interact with the firing aura through a stable surface regardless of the underlying
// entry's storage.
type auraCtx struct {
	a *Aura
	g *gameengine.GameEngine
}

func (c *auraCtx) Count() int          { return c.a.count }
func (c *auraCtx) DecrementCount() int { c.a.count--; return c.a.count }
func (c *auraCtx) CardName() string    { return c.a.CardName() }
func (c *auraCtx) CardID() ids.CardID  { return c.a.CardID() }
func (c *auraCtx) Destroy(addToGraveyard bool) {
	c.g.DestroyAura(addToGraveyard)
}
