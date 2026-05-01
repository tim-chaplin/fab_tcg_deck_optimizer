// Barraging Brawnhide — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "While Barraging Brawnhide is defended by less than 2 non-equipment cards, it has +1{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var barragingBrawnhideTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BarragingBrawnhideRed struct{}

func (BarragingBrawnhideRed) ID() ids.CardID          { return ids.BarragingBrawnhideRed }
func (BarragingBrawnhideRed) Name() string            { return "Barraging Brawnhide" }
func (BarragingBrawnhideRed) Cost(*sim.TurnState) int { return 3 }
func (BarragingBrawnhideRed) Pitch() int              { return 1 }
func (BarragingBrawnhideRed) Attack() int             { return 7 }
func (BarragingBrawnhideRed) Defense() int            { return 2 }
func (BarragingBrawnhideRed) Types() card.TypeSet     { return barragingBrawnhideTypes }
func (BarragingBrawnhideRed) GoAgain() bool           { return false }

// not implemented: defended-by-<2-non-equipment +1{p} rider
func (BarragingBrawnhideRed) NotImplemented() {}
func (c BarragingBrawnhideRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}

type BarragingBrawnhideYellow struct{}

func (BarragingBrawnhideYellow) ID() ids.CardID          { return ids.BarragingBrawnhideYellow }
func (BarragingBrawnhideYellow) Name() string            { return "Barraging Brawnhide" }
func (BarragingBrawnhideYellow) Cost(*sim.TurnState) int { return 3 }
func (BarragingBrawnhideYellow) Pitch() int              { return 2 }
func (BarragingBrawnhideYellow) Attack() int             { return 6 }
func (BarragingBrawnhideYellow) Defense() int            { return 2 }
func (BarragingBrawnhideYellow) Types() card.TypeSet     { return barragingBrawnhideTypes }
func (BarragingBrawnhideYellow) GoAgain() bool           { return false }

// not implemented: defended-by-<2-non-equipment +1{p} rider
func (BarragingBrawnhideYellow) NotImplemented() {}
func (c BarragingBrawnhideYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}

type BarragingBrawnhideBlue struct{}

func (BarragingBrawnhideBlue) ID() ids.CardID          { return ids.BarragingBrawnhideBlue }
func (BarragingBrawnhideBlue) Name() string            { return "Barraging Brawnhide" }
func (BarragingBrawnhideBlue) Cost(*sim.TurnState) int { return 3 }
func (BarragingBrawnhideBlue) Pitch() int              { return 3 }
func (BarragingBrawnhideBlue) Attack() int             { return 5 }
func (BarragingBrawnhideBlue) Defense() int            { return 2 }
func (BarragingBrawnhideBlue) Types() card.TypeSet     { return barragingBrawnhideTypes }
func (BarragingBrawnhideBlue) GoAgain() bool           { return false }

// not implemented: defended-by-<2-non-equipment +1{p} rider
func (BarragingBrawnhideBlue) NotImplemented() {}
func (c BarragingBrawnhideBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}
