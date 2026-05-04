// Talishar, the Lost Prince — Generic Weapon - Sword (2H). Power 4.
//
// Text: "**Once per Turn Action** - {r}{r}, put a rust counter on Talishar, the Lost Prince:
// **Attack** At the beginning of your end phase, if Talishar, the Lost Prince has 3 or more rust
// counters on it, destroy it."

package weapons

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var talisharTypes = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeSword, card.TypeTwoHand)

type Talishar struct{}

func (Talishar) ID() ids.WeaponID    { return ids.TalisharID }
func (Talishar) Name() string        { return "Talishar, the Lost Prince" }
func (Talishar) Types() card.TypeSet { return talisharTypes }
func (Talishar) Hands() int          { return 2 }
func (Talishar) Ability() sim.Card   { return talisharAbility }

// Cached at package init — see nebula_blade.go for the alloc-free rationale.
var talisharAbility sim.Card = TalisharAbility{}

// not implemented: rust-counter activation cost and end-phase self-destruct at 3+ counters
func (Talishar) NotImplemented() {}

var talisharAbilityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeSword, card.TypeTwoHand, card.TypeAttack)

type TalisharAbility struct{}

func (TalisharAbility) ID() ids.CardID          { return ids.TalisharAbilityID }
func (TalisharAbility) Name() string            { return "Talishar, the Lost Prince" }
func (TalisharAbility) Cost(*sim.TurnState) int { return 0 }
func (TalisharAbility) Pitch() int              { return 0 }
func (TalisharAbility) Attack() int             { return 4 }
func (TalisharAbility) Defense() int            { return 0 }
func (TalisharAbility) Types() card.TypeSet     { return talisharAbilityTypes }
func (TalisharAbility) GoAgain() bool           { return false }
func (TalisharAbility) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
