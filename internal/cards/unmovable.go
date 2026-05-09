// Unmovable — Generic Defense Reaction. Cost 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 7, Yellow 6, Blue 5.
// Text: "If Unmovable is played from arsenal, it gains +1{d}."
//
// +1{d} when played from arsenal via sim.ArsenalDefenseBonus (docs/dev-standards.md).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (UnmovableRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}
func (UnmovableRed) ArsenalDefenseBonus() int { return 1 }

func (UnmovableYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}
func (UnmovableYellow) ArsenalDefenseBonus() int { return 1 }

func (UnmovableBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}
func (UnmovableBlue) ArsenalDefenseBonus() int { return 1 }
