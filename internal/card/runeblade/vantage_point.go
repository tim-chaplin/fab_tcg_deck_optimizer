// Vantage Point — Runeblade Action - Attack. Cost 3, Defense 3.
// Printed power: Red 7, Yellow 6, Blue 5.
// Text: "If you've played or created an aura this turn, this gets **overpower**."
//
// Sets TurnState.Overpower when the aura condition is met.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var vantagePointTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type VantagePointRed struct{}

func (VantagePointRed) ID() card.ID              { return card.VantagePointRed }
func (VantagePointRed) Name() string             { return "Vantage Point" }
func (VantagePointRed) Cost(*card.TurnState) int { return 3 }
func (VantagePointRed) Pitch() int               { return 1 }
func (VantagePointRed) Attack() int              { return 7 }
func (VantagePointRed) Defense() int             { return 3 }
func (VantagePointRed) Types() card.TypeSet      { return vantagePointTypes }
func (VantagePointRed) GoAgain() bool            { return false }

// not implemented: Overpower flag is set on the aura condition but never consumed by the solver
func (VantagePointRed) NotImplemented() {}
func (c VantagePointRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, vantagePointPlay(c.Attack(), s)-self.Card.Attack())
}

type VantagePointYellow struct{}

func (VantagePointYellow) ID() card.ID              { return card.VantagePointYellow }
func (VantagePointYellow) Name() string             { return "Vantage Point" }
func (VantagePointYellow) Cost(*card.TurnState) int { return 3 }
func (VantagePointYellow) Pitch() int               { return 2 }
func (VantagePointYellow) Attack() int              { return 6 }
func (VantagePointYellow) Defense() int             { return 3 }
func (VantagePointYellow) Types() card.TypeSet      { return vantagePointTypes }
func (VantagePointYellow) GoAgain() bool            { return false }

// not implemented: Overpower flag is set on the aura condition but never consumed by the solver
func (VantagePointYellow) NotImplemented() {}
func (c VantagePointYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, vantagePointPlay(c.Attack(), s)-self.Card.Attack())
}

type VantagePointBlue struct{}

func (VantagePointBlue) ID() card.ID              { return card.VantagePointBlue }
func (VantagePointBlue) Name() string             { return "Vantage Point" }
func (VantagePointBlue) Cost(*card.TurnState) int { return 3 }
func (VantagePointBlue) Pitch() int               { return 3 }
func (VantagePointBlue) Attack() int              { return 5 }
func (VantagePointBlue) Defense() int             { return 3 }
func (VantagePointBlue) Types() card.TypeSet      { return vantagePointTypes }
func (VantagePointBlue) GoAgain() bool            { return false }

// not implemented: Overpower flag is set on the aura condition but never consumed by the solver
func (VantagePointBlue) NotImplemented() {}
func (c VantagePointBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, vantagePointPlay(c.Attack(), s)-self.Card.Attack())
}
func vantagePointPlay(base int, s *card.TurnState) int {
	if s.HasAuraInPlay() {
		s.Overpower = true
	}
	return base
}
