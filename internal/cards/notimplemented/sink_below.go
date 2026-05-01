// Sink Below — Generic Defense Reaction. Cost 0.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "You may put a card from your hand on the bottom of your deck. If you do, draw a card."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
)

type SinkBelowRed struct{}

func (SinkBelowRed) ID() ids.CardID          { return ids.SinkBelowRed }
func (SinkBelowRed) Name() string            { return "Sink Below" }
func (SinkBelowRed) Cost(*sim.TurnState) int { return 0 }
func (SinkBelowRed) Pitch() int              { return 1 }
func (SinkBelowRed) Attack() int             { return 0 }
func (SinkBelowRed) Defense() int            { return 4 }
func (SinkBelowRed) Types() card.TypeSet     { return cards.DefenseReactionTypes }
func (SinkBelowRed) GoAgain() bool           { return false }
func (SinkBelowRed) NotSilverAgeLegal()      {}

// not implemented: discard-to-cycle rider (hand cycling not modelled)
func (SinkBelowRed) NotImplemented() {}
func (SinkBelowRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}

type SinkBelowYellow struct{}

func (SinkBelowYellow) ID() ids.CardID          { return ids.SinkBelowYellow }
func (SinkBelowYellow) Name() string            { return "Sink Below" }
func (SinkBelowYellow) Cost(*sim.TurnState) int { return 0 }
func (SinkBelowYellow) Pitch() int              { return 2 }
func (SinkBelowYellow) Attack() int             { return 0 }
func (SinkBelowYellow) Defense() int            { return 3 }
func (SinkBelowYellow) Types() card.TypeSet     { return cards.DefenseReactionTypes }
func (SinkBelowYellow) GoAgain() bool           { return false }
func (SinkBelowYellow) NotSilverAgeLegal()      {}

// not implemented: discard-to-cycle rider (hand cycling not modelled)
func (SinkBelowYellow) NotImplemented() {}
func (SinkBelowYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}

type SinkBelowBlue struct{}

func (SinkBelowBlue) ID() ids.CardID          { return ids.SinkBelowBlue }
func (SinkBelowBlue) Name() string            { return "Sink Below" }
func (SinkBelowBlue) Cost(*sim.TurnState) int { return 0 }
func (SinkBelowBlue) Pitch() int              { return 3 }
func (SinkBelowBlue) Attack() int             { return 0 }
func (SinkBelowBlue) Defense() int            { return 2 }
func (SinkBelowBlue) Types() card.TypeSet     { return cards.DefenseReactionTypes }
func (SinkBelowBlue) GoAgain() bool           { return false }
func (SinkBelowBlue) NotSilverAgeLegal()      {}

// not implemented: discard-to-cycle rider (hand cycling not modelled)
func (SinkBelowBlue) NotImplemented() {}
func (SinkBelowBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}
