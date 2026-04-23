// Enchanting Melody — Generic Action - Aura. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue
// 3. Defense 2.
//
// Text: "**Go again** If your hero would be dealt damage, instead destroy Enchanting Melody and
// prevent 4 damage that source would deal. At the beginning of your end phase, destroy Enchanting
// Melody unless you have played a 'non-attack' action card this turn."
//
// The aura-created flag is set so same-turn aura-readers (Yinti Yanti, Runerager Swarm, etc.)
// see the entry.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var enchantingMelodyTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAura)

type EnchantingMelodyRed struct{}

func (EnchantingMelodyRed) ID() card.ID                 { return card.EnchantingMelodyRed }
func (EnchantingMelodyRed) Name() string                { return "Enchanting Melody (Red)" }
func (EnchantingMelodyRed) Cost(*card.TurnState) int                   { return 2 }
func (EnchantingMelodyRed) Pitch() int                  { return 1 }
func (EnchantingMelodyRed) Attack() int                 { return 0 }
func (EnchantingMelodyRed) Defense() int                { return 2 }
func (EnchantingMelodyRed) Types() card.TypeSet         { return enchantingMelodyTypes }
func (EnchantingMelodyRed) GoAgain() bool               { return true }
// not implemented: damage-prevention trigger, end-phase destruction clause
func (EnchantingMelodyRed) NotImplemented()             {}
func (EnchantingMelodyRed) Play(s *card.TurnState, _ *card.CardState) int { return setAuraCreated(s) }

type EnchantingMelodyYellow struct{}

func (EnchantingMelodyYellow) ID() card.ID                 { return card.EnchantingMelodyYellow }
func (EnchantingMelodyYellow) Name() string                { return "Enchanting Melody (Yellow)" }
func (EnchantingMelodyYellow) Cost(*card.TurnState) int                   { return 2 }
func (EnchantingMelodyYellow) Pitch() int                  { return 2 }
func (EnchantingMelodyYellow) Attack() int                 { return 0 }
func (EnchantingMelodyYellow) Defense() int                { return 2 }
func (EnchantingMelodyYellow) Types() card.TypeSet         { return enchantingMelodyTypes }
func (EnchantingMelodyYellow) GoAgain() bool               { return true }
// not implemented: damage-prevention trigger, end-phase destruction clause
func (EnchantingMelodyYellow) NotImplemented()             {}
func (EnchantingMelodyYellow) Play(s *card.TurnState, _ *card.CardState) int { return setAuraCreated(s) }

type EnchantingMelodyBlue struct{}

func (EnchantingMelodyBlue) ID() card.ID                 { return card.EnchantingMelodyBlue }
func (EnchantingMelodyBlue) Name() string                { return "Enchanting Melody (Blue)" }
func (EnchantingMelodyBlue) Cost(*card.TurnState) int                   { return 2 }
func (EnchantingMelodyBlue) Pitch() int                  { return 3 }
func (EnchantingMelodyBlue) Attack() int                 { return 0 }
func (EnchantingMelodyBlue) Defense() int                { return 2 }
func (EnchantingMelodyBlue) Types() card.TypeSet         { return enchantingMelodyTypes }
func (EnchantingMelodyBlue) GoAgain() bool               { return true }
// not implemented: damage-prevention trigger, end-phase destruction clause
func (EnchantingMelodyBlue) NotImplemented()             {}
func (EnchantingMelodyBlue) Play(s *card.TurnState, _ *card.CardState) int { return setAuraCreated(s) }
