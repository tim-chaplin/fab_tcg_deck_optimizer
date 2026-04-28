// Brandish — Generic Action - Attack. Cost 1. Printed power: Red 3, Yellow 2, Blue 1. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Brandish hits, your next weapon attack this turn gains +1{p}. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var brandishTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BrandishRed struct{}

func (BrandishRed) ID() ids.CardID          { return ids.BrandishRed }
func (BrandishRed) Name() string            { return "Brandish" }
func (BrandishRed) Cost(*sim.TurnState) int { return 1 }
func (BrandishRed) Pitch() int              { return 1 }
func (BrandishRed) Attack() int             { return 3 }
func (BrandishRed) Defense() int            { return 2 }
func (BrandishRed) Types() card.TypeSet     { return brandishTypes }
func (BrandishRed) GoAgain() bool           { return true }

// not implemented: next-weapon-attack +1{p} grant (weapon chain not scanned)
func (BrandishRed) NotImplemented() {}
func (c BrandishRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type BrandishYellow struct{}

func (BrandishYellow) ID() ids.CardID          { return ids.BrandishYellow }
func (BrandishYellow) Name() string            { return "Brandish" }
func (BrandishYellow) Cost(*sim.TurnState) int { return 1 }
func (BrandishYellow) Pitch() int              { return 2 }
func (BrandishYellow) Attack() int             { return 2 }
func (BrandishYellow) Defense() int            { return 2 }
func (BrandishYellow) Types() card.TypeSet     { return brandishTypes }
func (BrandishYellow) GoAgain() bool           { return true }

// not implemented: next-weapon-attack +1{p} grant (weapon chain not scanned)
func (BrandishYellow) NotImplemented() {}
func (c BrandishYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type BrandishBlue struct{}

func (BrandishBlue) ID() ids.CardID          { return ids.BrandishBlue }
func (BrandishBlue) Name() string            { return "Brandish" }
func (BrandishBlue) Cost(*sim.TurnState) int { return 1 }
func (BrandishBlue) Pitch() int              { return 3 }
func (BrandishBlue) Attack() int             { return 1 }
func (BrandishBlue) Defense() int            { return 2 }
func (BrandishBlue) Types() card.TypeSet     { return brandishTypes }
func (BrandishBlue) GoAgain() bool           { return true }

// not implemented: next-weapon-attack +1{p} grant (weapon chain not scanned)
func (BrandishBlue) NotImplemented() {}
func (c BrandishBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
