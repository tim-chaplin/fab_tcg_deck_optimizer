package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// auraEntry is the canonical engine-owned Aura impl. Card-driven Create*Aura methods on
// *GameEngine build one of these and append it to g.auras; sim's token-aura factories
// (NewRunechantAura, NewPonderAura, …) build via NewTokenAura which also returns an
// auraEntry. Both flavours satisfy the public Aura interface.
type auraEntry struct {
	triggerType   TriggerType
	handler       card.AuraHandler
	source        card.Card // nil for token auras
	tokenName     string    // "Runechant" / "Ponder" / "" for card auras
	tokenID       ids.CardID
	count         int
	oncePerTurn   bool
	firedThisTurn bool
}

// NewCardAura builds a card-backed aura — source is self.Card. The aura's CardName and
// CardID surface self.Card's DisplayName / ID. On destroy with addToGraveyard=true the
// source card lands in the graveyard.
func NewCardAura(self *card.CardState, tt TriggerType, handler card.AuraHandler, count int, oncePerTurn bool) Aura {
	return &auraEntry{
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
func NewTokenAura(name string, tokenID ids.CardID, tt TriggerType, handler card.AuraHandler, count int) Aura {
	return &auraEntry{
		triggerType: tt,
		handler:     handler,
		tokenName:   name,
		tokenID:     tokenID,
		count:       count,
	}
}

func (a *auraEntry) TriggerType() TriggerType { return a.triggerType }
func (a *auraEntry) OncePerTurn() bool        { return a.oncePerTurn }
func (a *auraEntry) FiredThisTurn() bool      { return a.firedThisTurn }
func (a *auraEntry) SetFiredThisTurn(v bool)  { a.firedThisTurn = v }

func (a *auraEntry) CardName() string {
	if a.source != nil {
		return a.source.DisplayName()
	}
	return a.tokenName
}

func (a *auraEntry) CardID() ids.CardID {
	if a.source != nil {
		return a.source.ID()
	}
	return a.tokenID
}

func (a *auraEntry) SourceCard() card.Card { return a.source }
func (a *auraEntry) Count() int            { return a.count }
func (a *auraEntry) SetCount(n int)        { a.count = n }
func (a *auraEntry) DecrementCount() int   { a.count--; return a.count }

func (a *auraEntry) Fire(g *GameEngine, l card.Logger) {
	ctx := &auraCtx{a: a, g: g}
	a.handler(g, l, ctx)
}

func (a *auraEntry) OnDestroy(g *GameEngine) {
	if a.source != nil {
		g.AppendGraveyard(a.source)
	}
}

func (a *auraEntry) Clone() Aura {
	out := *a
	return &out
}

// auraCtx is the engine-built adapter handed to each aura handler — it satisfies the
// card.Aura interface so cards interact with the firing aura through a stable surface
// regardless of the underlying entry's storage.
type auraCtx struct {
	a *auraEntry
	g *GameEngine
}

func (c *auraCtx) Count() int          { return c.a.count }
func (c *auraCtx) DecrementCount() int { c.a.count--; return c.a.count }
func (c *auraCtx) CardName() string    { return c.a.CardName() }
func (c *auraCtx) CardID() ids.CardID  { return c.a.CardID() }
func (c *auraCtx) Destroy(addToGraveyard bool) {
	c.g.DestroyAura(addToGraveyard)
}

// === Card-facing aura creation methods on GameEngine ===

// CreateStartOfTurnAura registers a TriggerStartOfTurn aura: the handler fires at the
// start of each subsequent turn.
func (g *GameEngine) CreateStartOfTurnAura(self *card.CardState, handler card.AuraHandler, count int) {
	g.AppendAura(NewCardAura(self, TriggerStartOfTurn, handler, count, false))
}

// CreateOncePerTurnAttackActionAura registers a TriggerAttackAction aura with the
// OncePerTurn gate set — fires at most once per turn regardless of how many attack
// actions resolve.
func (g *GameEngine) CreateOncePerTurnAttackActionAura(self *card.CardState, handler card.AuraHandler, count int) {
	g.AppendAura(NewCardAura(self, TriggerAttackAction, handler, count, true))
}
