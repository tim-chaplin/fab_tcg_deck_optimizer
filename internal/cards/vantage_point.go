// Vantage Point — Runeblade Action - Attack.
//
// Text: "If you've played or created an aura this turn, this gets **overpower**."
//
// Credits sim.OverpowerValue (0) for the granted Overpower; flag still flips on s.Overpower
// for any future consumer.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var vantagePointTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// vantagePointPlay flips s.Overpower when an aura has been played or created this turn, then
// emits the chain step.
func vantagePointPlay(s *sim.TurnState, self *sim.CardState) {
	if s.HasPlayedOrCreatedAura() {
		s.Overpower = true
		s.AddValue(sim.OverpowerValue)
	}
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type VantagePointRed struct{}

func (VantagePointRed) ID() ids.CardID          { return ids.VantagePointRed }
func (VantagePointRed) Name() string            { return "Vantage Point" }
func (VantagePointRed) Cost(*sim.TurnState) int { return 3 }
func (VantagePointRed) Pitch() int              { return 1 }
func (VantagePointRed) Attack() int             { return 7 }
func (VantagePointRed) Defense() int            { return 3 }
func (VantagePointRed) Types() card.TypeSet     { return vantagePointTypes }
func (VantagePointRed) GoAgain() bool           { return false }
func (VantagePointRed) Play(s *sim.TurnState, self *sim.CardState) {
	vantagePointPlay(s, self)
}

type VantagePointYellow struct{}

func (VantagePointYellow) ID() ids.CardID          { return ids.VantagePointYellow }
func (VantagePointYellow) Name() string            { return "Vantage Point" }
func (VantagePointYellow) Cost(*sim.TurnState) int { return 3 }
func (VantagePointYellow) Pitch() int              { return 2 }
func (VantagePointYellow) Attack() int             { return 6 }
func (VantagePointYellow) Defense() int            { return 3 }
func (VantagePointYellow) Types() card.TypeSet     { return vantagePointTypes }
func (VantagePointYellow) GoAgain() bool           { return false }
func (VantagePointYellow) Play(s *sim.TurnState, self *sim.CardState) {
	vantagePointPlay(s, self)
}

type VantagePointBlue struct{}

func (VantagePointBlue) ID() ids.CardID          { return ids.VantagePointBlue }
func (VantagePointBlue) Name() string            { return "Vantage Point" }
func (VantagePointBlue) Cost(*sim.TurnState) int { return 3 }
func (VantagePointBlue) Pitch() int              { return 3 }
func (VantagePointBlue) Attack() int             { return 5 }
func (VantagePointBlue) Defense() int            { return 3 }
func (VantagePointBlue) Types() card.TypeSet     { return vantagePointTypes }
func (VantagePointBlue) GoAgain() bool           { return false }
func (VantagePointBlue) Play(s *sim.TurnState, self *sim.CardState) {
	vantagePointPlay(s, self)
}
