// Cut Down to Size — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, if they have 4 or more cards in hand, they discard a card."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var cutDownToSizeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type CutDownToSizeRed struct{}

func (CutDownToSizeRed) ID() ids.CardID          { return ids.CutDownToSizeRed }
func (CutDownToSizeRed) Name() string            { return "Cut Down to Size" }
func (CutDownToSizeRed) Cost(*sim.TurnState) int { return 2 }
func (CutDownToSizeRed) Pitch() int              { return 1 }
func (CutDownToSizeRed) Attack() int             { return 6 }
func (CutDownToSizeRed) Defense() int            { return 2 }
func (CutDownToSizeRed) Types() card.TypeSet     { return cutDownToSizeTypes }
func (CutDownToSizeRed) GoAgain() bool           { return false }

// not implemented: on-hit opponent discard (conditional on hand size)
func (CutDownToSizeRed) NotImplemented() {}
func (CutDownToSizeRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type CutDownToSizeYellow struct{}

func (CutDownToSizeYellow) ID() ids.CardID          { return ids.CutDownToSizeYellow }
func (CutDownToSizeYellow) Name() string            { return "Cut Down to Size" }
func (CutDownToSizeYellow) Cost(*sim.TurnState) int { return 2 }
func (CutDownToSizeYellow) Pitch() int              { return 2 }
func (CutDownToSizeYellow) Attack() int             { return 5 }
func (CutDownToSizeYellow) Defense() int            { return 2 }
func (CutDownToSizeYellow) Types() card.TypeSet     { return cutDownToSizeTypes }
func (CutDownToSizeYellow) GoAgain() bool           { return false }

// not implemented: on-hit opponent discard (conditional on hand size)
func (CutDownToSizeYellow) NotImplemented() {}
func (CutDownToSizeYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type CutDownToSizeBlue struct{}

func (CutDownToSizeBlue) ID() ids.CardID          { return ids.CutDownToSizeBlue }
func (CutDownToSizeBlue) Name() string            { return "Cut Down to Size" }
func (CutDownToSizeBlue) Cost(*sim.TurnState) int { return 2 }
func (CutDownToSizeBlue) Pitch() int              { return 3 }
func (CutDownToSizeBlue) Attack() int             { return 4 }
func (CutDownToSizeBlue) Defense() int            { return 2 }
func (CutDownToSizeBlue) Types() card.TypeSet     { return cutDownToSizeTypes }
func (CutDownToSizeBlue) GoAgain() bool           { return false }

// not implemented: on-hit opponent discard (conditional on hand size)
func (CutDownToSizeBlue) NotImplemented() {}
func (CutDownToSizeBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
