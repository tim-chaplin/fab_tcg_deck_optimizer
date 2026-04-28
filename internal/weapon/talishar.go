// Talishar, the Lost Prince — Generic Weapon - Sword (2H). Power 4.
//
// Text: "**Once per Turn Action** - {r}{r}, put a rust counter on Talishar, the Lost Prince:
// **Attack** At the beginning of your end phase, if Talishar, the Lost Prince has 3 or more rust
// counters on it, destroy it."

package weapon

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var talisharTypes = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeSword, card.TypeTwoHand)

type Talishar struct{}

func (Talishar) ID() ids.WeaponID         { return ids.TalisharID }
func (Talishar) Name() string             { return "Talishar, the Lost Prince" }
func (Talishar) Cost(*card.TurnState) int { return 0 }
func (Talishar) Pitch() int               { return 0 }
func (Talishar) Attack() int              { return 4 }
func (Talishar) Defense() int             { return 0 }
func (Talishar) Types() card.TypeSet      { return talisharTypes }
func (Talishar) GoAgain() bool            { return false }
func (Talishar) Hands() int               { return 2 }

// not implemented: rust-counter activation cost and end-phase self-destruct at 3+ counters
func (Talishar) NotImplemented() {}
func (Talishar) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
