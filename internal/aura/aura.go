// Package aura owns the concrete Aura type — a persistent hook entry that fires at a
// scheduled trigger type (start of turn, attack, attack action, end of turn). Card-backed
// and token-backed auras share the same struct; SourceCard distinguishes them.
package aura

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/trigger"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Handler is the typed aura handler signature: the func stored on each Aura and called at
// every Fire. The card.FireContext carries the firing event (ctx.FiringType — multi-trigger
// handlers dispatch on it), the triggering card (ctx.TriggeringCard, nil for turn-boundary
// events), and any trigger-specific payload.
type Handler func(card.GameEngine, card.Logger, card.Aura, card.FireContext)

// Aura is the concrete entry the engine stores in its persistent hook list. The embedded
// trigger.Trigger holds the shared core; Aura adds the multi-fire count.
//
// activeEngine is set by Fire to the engine driving the current firing event so Destroy can
// route back without allocating a per-fire wrapper struct. Single-threaded per attack-turn runner.
type Aura struct {
	trigger.Trigger[card.Aura]
	activeEngine card.GameEngine
	count        int
}

// NewFromCard builds a card-backed aura. source is the originating card — typically the
// Card field of a CardState. SourceCard surfaces it back so engines can route it into the
// graveyard on destroy. typeFilter narrows the firing site; pass nil for no filter.
func NewFromCard(source card.Card, tt triggertype.Type, fire Handler, count int, oncePerTurn bool, typeFilter trigger.TypeFilter) *Aura {
	return &Aura{
		Trigger: trigger.FromCard[card.Aura](source, tt, fire, oncePerTurn, typeFilter),
		count:   count,
	}
}

// NewFromToken builds a token aura — no originating card. CardName returns the supplied
// name; CardID returns tokenID so cache keys distinguish each token kind. typeFilter
// narrows the firing site; pass nil for no filter.
func NewFromToken(name string, tokenID ids.CardID, tt triggertype.Type, fire Handler, count int, typeFilter trigger.TypeFilter) *Aura {
	return &Aura{
		Trigger: trigger.FromToken[card.Aura](name, tokenID, tt, fire, false, typeFilter),
		count:   count,
	}
}

func (a *Aura) Count() int     { return a.count }
func (a *Aura) SetCount(n int) { a.count = n }
func (a *Aura) DecrementCount() int {
	a.count--
	return a.count
}

// Fire invokes the stored handler with this aura as the typed receiver, passing the
// FireContext so multi-trigger handlers can dispatch. activeEngine is set for the
// handler's duration and cleared afterward so a Destroy outside the firing window can't
// call into a stale engine.
func (a *Aura) Fire(engine card.GameEngine, logger card.Logger, ctx card.FireContext) {
	a.activeEngine = engine
	a.Invoke(engine, logger, a, ctx)
	a.activeEngine = nil
}

// Copy returns a deep copy boxed as any so the consumer interface can avoid referencing
// its own concrete type. activeEngine is cleared on the copy.
func (a *Aura) Copy() any {
	out := *a
	out.activeEngine = nil
	return &out
}

// CopyInto rewrites *dst with a's fields and returns dst as any. When dst is a *Aura
// from a prior permutation's slot the per-perm reset reuses the existing allocation
// instead of paying for a fresh one. dst types other than *Aura fall back to Copy.
func (a *Aura) CopyInto(dst any) any {
	d, ok := dst.(*Aura)
	if !ok {
		return a.Copy()
	}
	*d = *a
	d.activeEngine = nil
	return d
}

// Compile-time check that *Aura satisfies card.Aura.
var _ card.Aura = (*Aura)(nil)

// Destroy ends the aura currently being fired, routed through activeEngine. Calling it
// outside a Fire window panics deterministically rather than silently no-oping.
func (a *Aura) Destroy(addToGraveyard bool) {
	a.activeEngine.DestroyAura(addToGraveyard)
}

// SetFromCard rewrites a in place for reuse as a card-sourced aura — the entry the
// gameengine aura pool hands back to a CreateAura caller after locating a free slot.
// Clears activeEngine and firedThisTurn so the recycled slot fires cleanly.
func (a *Aura) SetFromCard(source card.Card, tt triggertype.Type, fire Handler, count int, oncePerTurn bool, typeFilter trigger.TypeFilter) {
	a.Trigger.SetFromCard(source, tt, fire, oncePerTurn, typeFilter)
	a.count = count
	a.activeEngine = nil
}

// Clear zeros every field so the slot reads as "free" for the gameengine aura pool's
// next-free-slot scan. SourceCard() == nil and Count() == 0 are the canonical dead-slot
// tests pool readers gate on.
func (a *Aura) Clear() {
	a.Trigger.Clear()
	a.count = 0
	a.activeEngine = nil
}
