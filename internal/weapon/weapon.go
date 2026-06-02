// Package weapon owns two pieces of the Flesh and Blood weapon model, mirroring the
// internal/aura and internal/item split:
//
//   - Card: the platonic weapon card (a card.Card that also reports Hands and an activated
//     Ability). Each concrete weapon in internal/weapon/weapons satisfies it; the card lives
//     in the card universe and goes to the graveyard when its equipped object is destroyed.
//   - Weapon: the concrete mutable engine-side object built when the card is equipped. It
//     embeds trigger.Trigger so an end-phase handler (Talishar's rust-counter self-destruct)
//     can subscribe, and carries the per-turn counter total. A deck equips 0–2 weapons under
//     the "2 × 1H or 1 × 2H" rule.
//
// This package must NOT import gameengine — same layering rule as internal/aura /
// internal/item. The concrete Weapon's Destroy routes back through card.GameEngine, which
// gameengine.GameEngine satisfies structurally.
package weapon

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/trigger"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Card is the platonic weapon card. It is a full card.Card (so it has a CardID, a
// DisplayName, a Types line, and lives in the card universe / graveyard) plus the two
// weapon-specific accessors the attack-turn runner reads: Hands for the equipment-slot rule
// and Ability for the swing card it enqueues each turn. A weapon registers its equipped object
// (and any trigger) from its Play via ge.CreateWeapon — the same way an item card's Play calls
// ge.CreateItem — so there is no weapon-trigger accessor on the card.
type Card interface {
	card.Card
	Hands() int
	Ability() card.Card
}

// Handler is the typed weapon handler signature — the func stored on a Weapon and called at
// every Fire. Mirrors aura.Handler / the item handler.
type Handler func(card.GameEngine, card.Logger, card.Weapon, card.FireContext)

// Weapon is the concrete mutable entry the engine stores in its equipped-weapon list. The
// embedded trigger.Trigger holds the source card and, for a triggered weapon, the firing
// event and handler; Weapon adds the counter total and caches the platonic card's Hands /
// Ability so the attack-turn runner reads them without re-asserting the source.
//
// activeEngine is set by Fire so a handler-side Destroy routes back through
// engine.DestroyWeapon without allocating a per-fire wrapper. Cleared after the handler
// returns.
type Weapon struct {
	trigger.Trigger[card.Weapon]
	activeEngine card.GameEngine
	source       Card
	count        int
}

// NewFromCard builds the engine-side Weapon for the supplied platonic card. The handler fires
// on every event in tt's bit set; pass tt=0 / fire=nil for an untriggered weapon. oncePerTurn
// caps it to one fire per turn; typeFilter narrows the firing site (pass nil for none).
// count starts at 0. Mirrors item.NewFromCard.
func NewFromCard(source Card, tt triggertype.Type, fire Handler, oncePerTurn bool, typeFilter trigger.TypeFilter) *Weapon {
	return &Weapon{
		Trigger: trigger.FromCard[card.Weapon](source, tt, fire, oncePerTurn, typeFilter),
		source:  source,
	}
}

func (w *Weapon) Count() int     { return w.count }
func (w *Weapon) SetCount(n int) { w.count = n }

// Name / Hands / Ability surface the platonic card's weapon attributes. The attack-turn
// runner enqueues Ability() as a 1-AP playable each turn; Hands gates the equipment-slot
// rule; Name attributes the swing log.
func (w *Weapon) Name() string       { return w.source.Name() }
func (w *Weapon) Hands() int         { return w.source.Hands() }
func (w *Weapon) Ability() card.Card { return w.source.Ability() }

// SourceCard returns the platonic weapon card boxed as any — the card sent to the graveyard
// on destroy. Mirrors Aura / Item; consumers assert to card.Card.
func (w *Weapon) SourceCard() any { return w.source }

// Fire invokes the stored handler with this weapon as the typed receiver, passing the
// FireContext so multi-trigger handlers can dispatch. activeEngine is set for the handler's
// duration and cleared afterward.
func (w *Weapon) Fire(engine card.GameEngine, logger card.Logger, ctx card.FireContext) {
	w.activeEngine = engine
	w.Invoke(engine, logger, w, ctx)
	w.activeEngine = nil
}

// Destroy removes this weapon from the arena, by object identity, through its bound engine;
// panics if no engine is bound (Fire binds it for a trigger handler).
func (w *Weapon) Destroy(addToGraveyard bool) {
	w.activeEngine.DestroyWeaponObject(w, addToGraveyard)
}

// Copy returns a deep copy boxed as any so the gameengine.Weapon interface can avoid
// referencing the concrete type. activeEngine is cleared on the copy.
func (w *Weapon) Copy() any {
	out := *w
	out.activeEngine = nil
	return &out
}

// CopyInto rewrites *dst with w's fields and returns dst as any. When dst is a *Weapon from
// a prior permutation's slot the per-perm reset reuses the existing allocation instead of
// paying for a fresh one — the equip loadout is identical across perms, so only the mutable
// counter / fire flag change. dst types other than *Weapon fall back to Copy.
func (w *Weapon) CopyInto(dst any) any {
	d, ok := dst.(*Weapon)
	if !ok {
		return w.Copy()
	}
	*d = *w
	d.activeEngine = nil
	return d
}

// Compile-time check that *Weapon satisfies card.Weapon.
var _ card.Weapon = (*Weapon)(nil)
