// Surging Militia — Generic Action - Attack. Cost 2. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "Surging Militia has +1{p} for each non-equipment card defending it."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var surgingMilitiaTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type SurgingMilitiaRed struct{}

func (SurgingMilitiaRed) ID() ids.CardID          { return ids.SurgingMilitiaRed }
func (SurgingMilitiaRed) Name() string            { return "Surging Militia" }
func (SurgingMilitiaRed) Cost(*sim.TurnState) int { return 2 }
func (SurgingMilitiaRed) Pitch() int              { return 1 }
func (SurgingMilitiaRed) Attack() int             { return 5 }
func (SurgingMilitiaRed) Defense() int            { return 2 }
func (SurgingMilitiaRed) Types() card.TypeSet     { return surgingMilitiaTypes }
func (SurgingMilitiaRed) GoAgain() bool           { return false }

// not implemented: defended-by +N{p} rider (defender's hand contents not exposed)
func (SurgingMilitiaRed) NotImplemented() {}
func (c SurgingMilitiaRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type SurgingMilitiaYellow struct{}

func (SurgingMilitiaYellow) ID() ids.CardID          { return ids.SurgingMilitiaYellow }
func (SurgingMilitiaYellow) Name() string            { return "Surging Militia" }
func (SurgingMilitiaYellow) Cost(*sim.TurnState) int { return 2 }
func (SurgingMilitiaYellow) Pitch() int              { return 2 }
func (SurgingMilitiaYellow) Attack() int             { return 4 }
func (SurgingMilitiaYellow) Defense() int            { return 2 }
func (SurgingMilitiaYellow) Types() card.TypeSet     { return surgingMilitiaTypes }
func (SurgingMilitiaYellow) GoAgain() bool           { return false }

// not implemented: defended-by +N{p} rider (defender's hand contents not exposed)
func (SurgingMilitiaYellow) NotImplemented() {}
func (c SurgingMilitiaYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type SurgingMilitiaBlue struct{}

func (SurgingMilitiaBlue) ID() ids.CardID          { return ids.SurgingMilitiaBlue }
func (SurgingMilitiaBlue) Name() string            { return "Surging Militia" }
func (SurgingMilitiaBlue) Cost(*sim.TurnState) int { return 2 }
func (SurgingMilitiaBlue) Pitch() int              { return 3 }
func (SurgingMilitiaBlue) Attack() int             { return 3 }
func (SurgingMilitiaBlue) Defense() int            { return 2 }
func (SurgingMilitiaBlue) Types() card.TypeSet     { return surgingMilitiaTypes }
func (SurgingMilitiaBlue) GoAgain() bool           { return false }

// not implemented: defended-by +N{p} rider (defender's hand contents not exposed)
func (SurgingMilitiaBlue) NotImplemented() {}
func (c SurgingMilitiaBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
