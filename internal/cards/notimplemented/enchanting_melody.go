// Enchanting Melody — Generic Action - Aura. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue
// 3. Defense 2.
//
// Text: "**Go again** If your hero would be dealt damage, instead destroy Enchanting Melody and
// prevent 4 damage that source would deal. At the beginning of your end phase, destroy Enchanting
// Melody unless you have played a 'non-attack' action card this turn."
//
// Sets s.AuraCreated so same-turn aura-readers see the entry.

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
)

var enchantingMelodyTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAura)

type EnchantingMelodyRed struct{}

func (EnchantingMelodyRed) ID() ids.CardID          { return ids.EnchantingMelodyRed }
func (EnchantingMelodyRed) Name() string            { return "Enchanting Melody" }
func (EnchantingMelodyRed) Cost(*sim.TurnState) int { return 2 }
func (EnchantingMelodyRed) Pitch() int              { return 1 }
func (EnchantingMelodyRed) Attack() int             { return 0 }
func (EnchantingMelodyRed) Defense() int            { return 2 }
func (EnchantingMelodyRed) Types() card.TypeSet     { return enchantingMelodyTypes }
func (EnchantingMelodyRed) GoAgain() bool           { return true }

// not implemented: damage-prevention trigger, end-phase destruction clause
func (EnchantingMelodyRed) NotImplemented() {}
func (EnchantingMelodyRed) Play(s *sim.TurnState, self *sim.CardState) {
	cards.SetAuraCreated(s)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type EnchantingMelodyYellow struct{}

func (EnchantingMelodyYellow) ID() ids.CardID          { return ids.EnchantingMelodyYellow }
func (EnchantingMelodyYellow) Name() string            { return "Enchanting Melody" }
func (EnchantingMelodyYellow) Cost(*sim.TurnState) int { return 2 }
func (EnchantingMelodyYellow) Pitch() int              { return 2 }
func (EnchantingMelodyYellow) Attack() int             { return 0 }
func (EnchantingMelodyYellow) Defense() int            { return 2 }
func (EnchantingMelodyYellow) Types() card.TypeSet     { return enchantingMelodyTypes }
func (EnchantingMelodyYellow) GoAgain() bool           { return true }

// not implemented: damage-prevention trigger, end-phase destruction clause
func (EnchantingMelodyYellow) NotImplemented() {}
func (EnchantingMelodyYellow) Play(s *sim.TurnState, self *sim.CardState) {
	cards.SetAuraCreated(s)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type EnchantingMelodyBlue struct{}

func (EnchantingMelodyBlue) ID() ids.CardID          { return ids.EnchantingMelodyBlue }
func (EnchantingMelodyBlue) Name() string            { return "Enchanting Melody" }
func (EnchantingMelodyBlue) Cost(*sim.TurnState) int { return 2 }
func (EnchantingMelodyBlue) Pitch() int              { return 3 }
func (EnchantingMelodyBlue) Attack() int             { return 0 }
func (EnchantingMelodyBlue) Defense() int            { return 2 }
func (EnchantingMelodyBlue) Types() card.TypeSet     { return enchantingMelodyTypes }
func (EnchantingMelodyBlue) GoAgain() bool           { return true }

// not implemented: damage-prevention trigger, end-phase destruction clause
func (EnchantingMelodyBlue) NotImplemented() {}
func (EnchantingMelodyBlue) Play(s *sim.TurnState, self *sim.CardState) {
	cards.SetAuraCreated(s)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
