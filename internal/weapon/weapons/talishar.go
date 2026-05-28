// Talishar, the Lost Prince — Generic Weapon - Sword (2H). Power 4.
//
// Text: "**Once per Turn Action** - {r}{r}, put a rust counter on Talishar, the Lost Prince:
// **Attack** At the beginning of your end phase, if Talishar, the Lost Prince has 3 or more rust
// counters on it, destroy it."

package weapons

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// talisharRustLimit is the rust-counter total at which Talishar self-destructs in the end
// phase. Hardcoded here — it's specific to this weapon, not a generic mechanic.
const talisharRustLimit = 3

var talisharTypes = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeSword, card.TypeTwoHand)

// Talishar is the platonic weapon card — the equipped permanent. Its activated ability is
// what the attack-turn runner enqueues each turn.
type Talishar struct{}

func (Talishar) ID() ids.CardID                                     { return ids.TalisharID }
func (Talishar) Name() string                                       { return "Talishar, the Lost Prince" }
func (Talishar) DisplayName() string                                { return "Talishar, the Lost Prince" }
func (Talishar) Cost() int                                          { return 0 }
func (Talishar) Pitch() int                                         { return 0 }
func (Talishar) Attack() int                                        { return 0 }
func (Talishar) Defense() int                                       { return 0 }
func (Talishar) Types(card.GameEngine) card.TypeSet                 { return talisharTypes }
func (Talishar) GoAgain(card.GameEngine) bool                       { return false }
func (Talishar) Play(card.GameEngine, card.Logger, *card.CardState) {}
func (Talishar) Hands() int                                         { return 2 }
func (Talishar) Ability() card.Card                                 { return talisharAbility }

// WeaponTrigger registers Talishar's end-phase self-destruct: at end of turn, if it has
// accumulated talisharRustLimit rust counters, destroy it (the platonic card goes to the
// graveyard and the weapon can't be swung next turn).
func (Talishar) WeaponTrigger() (triggertype.Type, weapon.Handler) {
	return triggertype.EndOfTurn, talisharEndPhaseSelfDestruct
}

func talisharEndPhaseSelfDestruct(_ card.GameEngine, l card.Logger, self card.Weapon, _ card.FireContext) {
	if self.Count() < talisharRustLimit {
		return
	}
	l.AppendPostTrigger(self.CardName(), "Talishar destroyed at 3 rust counters", 0)
	self.Destroy(true)
}

var talisharAbility card.Card = TalisharAbility{}

var talisharAbilityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeSword, card.TypeTwoHand, card.TypeAttack)

type TalisharAbility struct{}

func (TalisharAbility) ID() ids.CardID                     { return ids.TalisharAbilityID }
func (TalisharAbility) Name() string                       { return "Talishar, the Lost Prince" }
func (TalisharAbility) DisplayName() string                { return "Talishar, the Lost Prince" }
func (TalisharAbility) Cost() int                          { return 2 }
func (TalisharAbility) Pitch() int                         { return 0 }
func (TalisharAbility) Attack() int                        { return 4 }
func (TalisharAbility) Defense() int                       { return 0 }
func (TalisharAbility) Types(card.GameEngine) card.TypeSet { return talisharAbilityTypes }
func (TalisharAbility) GoAgain(card.GameEngine) bool       { return false }

// Play puts a rust counter on the equipped Talishar — the activation's "put a rust counter"
// clause. The 4-power attack itself is resolved by the attack-turn runner from the ability's
// printed Attack; this handler only bumps the weapon's counter.
func (TalisharAbility) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if w := ge.WeaponByID(ids.TalisharID); w != nil {
		w.SetCount(w.Count() + 1)
	}
}
