// Hocus Pocus — Runeblade Action - Attack. Cost 0, Defense 3.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "When this attacks, create a Runechant token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var hocusPocusTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type HocusPocusRed struct{}

func (HocusPocusRed) ID() ids.CardID          { return ids.HocusPocusRed }
func (HocusPocusRed) Name() string            { return "Hocus Pocus" }
func (HocusPocusRed) Cost(*sim.TurnState) int { return 0 }
func (HocusPocusRed) Pitch() int              { return 1 }
func (HocusPocusRed) Attack() int             { return 3 }
func (HocusPocusRed) Defense() int            { return 3 }
func (HocusPocusRed) Types() card.TypeSet     { return hocusPocusTypes }
func (HocusPocusRed) GoAgain() bool           { return false }
func (HocusPocusRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	s.AddValue(s.CreateRunechants(1))
	s.LogRider(self, 1, "Created a runechant")
}

type HocusPocusYellow struct{}

func (HocusPocusYellow) ID() ids.CardID          { return ids.HocusPocusYellow }
func (HocusPocusYellow) Name() string            { return "Hocus Pocus" }
func (HocusPocusYellow) Cost(*sim.TurnState) int { return 0 }
func (HocusPocusYellow) Pitch() int              { return 2 }
func (HocusPocusYellow) Attack() int             { return 2 }
func (HocusPocusYellow) Defense() int            { return 3 }
func (HocusPocusYellow) Types() card.TypeSet     { return hocusPocusTypes }
func (HocusPocusYellow) GoAgain() bool           { return false }
func (HocusPocusYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	s.AddValue(s.CreateRunechants(1))
	s.LogRider(self, 1, "Created a runechant")
}

type HocusPocusBlue struct{}

func (HocusPocusBlue) ID() ids.CardID          { return ids.HocusPocusBlue }
func (HocusPocusBlue) Name() string            { return "Hocus Pocus" }
func (HocusPocusBlue) Cost(*sim.TurnState) int { return 0 }
func (HocusPocusBlue) Pitch() int              { return 3 }
func (HocusPocusBlue) Attack() int             { return 1 }
func (HocusPocusBlue) Defense() int            { return 3 }
func (HocusPocusBlue) Types() card.TypeSet     { return hocusPocusTypes }
func (HocusPocusBlue) GoAgain() bool           { return false }
func (HocusPocusBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	s.AddValue(s.CreateRunechants(1))
	s.LogRider(self, 1, "Created a runechant")
}
