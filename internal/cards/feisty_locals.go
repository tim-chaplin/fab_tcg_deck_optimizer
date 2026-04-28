// Feisty Locals — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If this is defended by an action card, this gets +2{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var feistyLocalsTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type FeistyLocalsRed struct{}

func (FeistyLocalsRed) ID() ids.CardID          { return ids.FeistyLocalsRed }
func (FeistyLocalsRed) Name() string            { return "Feisty Locals" }
func (FeistyLocalsRed) Cost(*sim.TurnState) int { return 0 }
func (FeistyLocalsRed) Pitch() int              { return 1 }
func (FeistyLocalsRed) Attack() int             { return 3 }
func (FeistyLocalsRed) Defense() int            { return 2 }
func (FeistyLocalsRed) Types() card.TypeSet     { return feistyLocalsTypes }
func (FeistyLocalsRed) GoAgain() bool           { return false }

// not implemented: defended-by-action-card +2{p} rider
func (FeistyLocalsRed) NotImplemented() {}
func (c FeistyLocalsRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type FeistyLocalsYellow struct{}

func (FeistyLocalsYellow) ID() ids.CardID          { return ids.FeistyLocalsYellow }
func (FeistyLocalsYellow) Name() string            { return "Feisty Locals" }
func (FeistyLocalsYellow) Cost(*sim.TurnState) int { return 0 }
func (FeistyLocalsYellow) Pitch() int              { return 2 }
func (FeistyLocalsYellow) Attack() int             { return 2 }
func (FeistyLocalsYellow) Defense() int            { return 2 }
func (FeistyLocalsYellow) Types() card.TypeSet     { return feistyLocalsTypes }
func (FeistyLocalsYellow) GoAgain() bool           { return false }

// not implemented: defended-by-action-card +2{p} rider
func (FeistyLocalsYellow) NotImplemented() {}
func (c FeistyLocalsYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type FeistyLocalsBlue struct{}

func (FeistyLocalsBlue) ID() ids.CardID          { return ids.FeistyLocalsBlue }
func (FeistyLocalsBlue) Name() string            { return "Feisty Locals" }
func (FeistyLocalsBlue) Cost(*sim.TurnState) int { return 0 }
func (FeistyLocalsBlue) Pitch() int              { return 3 }
func (FeistyLocalsBlue) Attack() int             { return 1 }
func (FeistyLocalsBlue) Defense() int            { return 2 }
func (FeistyLocalsBlue) Types() card.TypeSet     { return feistyLocalsTypes }
func (FeistyLocalsBlue) GoAgain() bool           { return false }

// not implemented: defended-by-action-card +2{p} rider
func (FeistyLocalsBlue) NotImplemented() {}
func (c FeistyLocalsBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
