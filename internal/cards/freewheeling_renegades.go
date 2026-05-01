// Freewheeling Renegades — Generic Action - Attack. Cost 1. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If this is defended by an action card, this has -2{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var freewheelingRenegadesTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type FreewheelingRenegadesRed struct{}

func (FreewheelingRenegadesRed) ID() ids.CardID          { return ids.FreewheelingRenegadesRed }
func (FreewheelingRenegadesRed) Name() string            { return "Freewheeling Renegades" }
func (FreewheelingRenegadesRed) Cost(*sim.TurnState) int { return 1 }
func (FreewheelingRenegadesRed) Pitch() int              { return 1 }
func (FreewheelingRenegadesRed) Attack() int             { return 6 }
func (FreewheelingRenegadesRed) Defense() int            { return 2 }
func (FreewheelingRenegadesRed) Types() card.TypeSet     { return freewheelingRenegadesTypes }
func (FreewheelingRenegadesRed) GoAgain() bool           { return false }

// not implemented: defended-by-action-card -2{p} rider
func (FreewheelingRenegadesRed) NotImplemented() {}
func (c FreewheelingRenegadesRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}

type FreewheelingRenegadesYellow struct{}

func (FreewheelingRenegadesYellow) ID() ids.CardID          { return ids.FreewheelingRenegadesYellow }
func (FreewheelingRenegadesYellow) Name() string            { return "Freewheeling Renegades" }
func (FreewheelingRenegadesYellow) Cost(*sim.TurnState) int { return 1 }
func (FreewheelingRenegadesYellow) Pitch() int              { return 2 }
func (FreewheelingRenegadesYellow) Attack() int             { return 5 }
func (FreewheelingRenegadesYellow) Defense() int            { return 2 }
func (FreewheelingRenegadesYellow) Types() card.TypeSet     { return freewheelingRenegadesTypes }
func (FreewheelingRenegadesYellow) GoAgain() bool           { return false }

// not implemented: defended-by-action-card -2{p} rider
func (FreewheelingRenegadesYellow) NotImplemented() {}
func (c FreewheelingRenegadesYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}

type FreewheelingRenegadesBlue struct{}

func (FreewheelingRenegadesBlue) ID() ids.CardID          { return ids.FreewheelingRenegadesBlue }
func (FreewheelingRenegadesBlue) Name() string            { return "Freewheeling Renegades" }
func (FreewheelingRenegadesBlue) Cost(*sim.TurnState) int { return 1 }
func (FreewheelingRenegadesBlue) Pitch() int              { return 3 }
func (FreewheelingRenegadesBlue) Attack() int             { return 4 }
func (FreewheelingRenegadesBlue) Defense() int            { return 2 }
func (FreewheelingRenegadesBlue) Types() card.TypeSet     { return freewheelingRenegadesTypes }
func (FreewheelingRenegadesBlue) GoAgain() bool           { return false }

// not implemented: defended-by-action-card -2{p} rider
func (FreewheelingRenegadesBlue) NotImplemented() {}
func (c FreewheelingRenegadesBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}
