// Package aura owns the concrete Aura type — a persistent hook entry that fires at a
// scheduled trigger type (start of turn, attack, attack action, end of turn). Card-backed
// and token-backed auras share the same struct; SourceCard distinguishes them.
//
// Handler is the aura-domain function shape; the package's own ctx struct implements
// card.Aura so handler bodies receive the typed surface directly without an outer
// wrapping layer.
package aura

import (
	"sync"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Handler is the typed aura handler signature. Lives in internal/aura (not internal/card) because
// the shape is meaningful only to aura: it's the signature stored on each Aura and called
// at every Fire. card.GameEngine.CreateStartOfTurnAura et al. inline the function type in
// their parameter declarations so internal/card doesn't need to import internal/aura.
type Handler func(card.GameEngine, card.Logger, card.Aura)

// Aura is the concrete entry the engine stores in its persistent hook list. source is
// non-nil for card-backed auras; tokenName / tokenID are populated for token auras.
//
// activeEngine is set by Fire to the engine driving the current firing event so the
// aura's own card.Aura-interface methods (Destroy) can route back without allocating a
// per-fire wrapper struct. Single-threaded per chain-runner; a copied aura clears
// activeEngine via Copy.
type Aura struct {
	triggerType   triggertype.Type
	fire          Handler
	source        card.Card
	tokenName     string
	tokenID       ids.CardID
	activeEngine  card.GameEngine
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
// name; CardID returns tokenID so cache keys distinguish each token kind. Pulls from
// auraPool to recycle *Aura allocations across mid-chain create/destroy cycles.
func NewFromToken(name string, tokenID ids.CardID, tt triggertype.Type, fire Handler, count int) *Aura {
	a := auraPool.Get().(*Aura)
	*a = Aura{
		triggerType: tt,
		fire:        fire,
		tokenName:   name,
		tokenID:     tokenID,
		count:       count,
	}
	return a
}

// auraPool recycles *Aura values destroyed mid-chain via the deferred-release path the
// engine drives (Release called only after Fire returns, so the post-fire activeEngine
// clear can't race a concurrent goroutine pulling the recycled pointer via Get).
var auraPool = sync.Pool{New: func() any { return &Aura{} }}

// Release returns a to the recycling pool. Only safe to call after Fire has fully
// returned; the engine's fireMatching loop enforces this by deferring Release until
// after Fire's activeEngine = nil completes. Card-backed auras (source != nil) skip
// the pool because they may live across turns and be referenced by snapshots.
func (a *Aura) Release() {
	if a == nil || a.source != nil {
		return
	}
	auraPool.Put(a)
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

// Fire invokes the stored handler. *Aura itself satisfies card.Aura — the activeEngine
// field is set so Destroy() can route back through engine.DestroyAura without allocating
// a per-fire wrapper. Cleared after the handler returns so a stray reference to the
// Aura outside its firing window can't accidentally call into a stale engine.
func (a *Aura) Fire(engine card.GameEngine, logger card.Logger) {
	a.activeEngine = engine
	a.fire(engine, logger, a)
	a.activeEngine = nil
}

// Copy returns a deep copy boxed as any so the gameengine.Aura interface declaration can
// avoid referencing its own type — letting concrete impls satisfy without importing.
// activeEngine is cleared on the copy so a per-permutation aura clone starts with no
// stale firing-engine pointer.
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

// Destroy ends the aura currently being fired. Routed through the engine reference Fire
// installed in activeEngine; calling Destroy outside a Fire window panics deterministically
// rather than silently no-oping.
func (a *Aura) Destroy(addToGraveyard bool) {
	a.activeEngine.DestroyAura(addToGraveyard)
}
