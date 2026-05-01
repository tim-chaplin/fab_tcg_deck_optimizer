// Right Behind You — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this defends together with another card from hand, this gets +1{d} and you may look
// at the top card of your deck. You may put it on the bottom."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var rightBehindYouTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type RightBehindYouRed struct{}

func (RightBehindYouRed) ID() ids.CardID          { return ids.RightBehindYouRed }
func (RightBehindYouRed) Name() string            { return "Right Behind You" }
func (RightBehindYouRed) Cost(*sim.TurnState) int { return 3 }
func (RightBehindYouRed) Pitch() int              { return 1 }
func (RightBehindYouRed) Attack() int             { return 7 }
func (RightBehindYouRed) Defense() int            { return 2 }
func (RightBehindYouRed) Types() card.TypeSet     { return rightBehindYouTypes }
func (RightBehindYouRed) GoAgain() bool           { return false }

// not implemented: defend-together +1{d} buff and deck-bottom peek rider
func (RightBehindYouRed) NotImplemented() {}
func (c RightBehindYouRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type RightBehindYouYellow struct{}

func (RightBehindYouYellow) ID() ids.CardID          { return ids.RightBehindYouYellow }
func (RightBehindYouYellow) Name() string            { return "Right Behind You" }
func (RightBehindYouYellow) Cost(*sim.TurnState) int { return 3 }
func (RightBehindYouYellow) Pitch() int              { return 2 }
func (RightBehindYouYellow) Attack() int             { return 6 }
func (RightBehindYouYellow) Defense() int            { return 2 }
func (RightBehindYouYellow) Types() card.TypeSet     { return rightBehindYouTypes }
func (RightBehindYouYellow) GoAgain() bool           { return false }

// not implemented: defend-together +1{d} buff and deck-bottom peek rider
func (RightBehindYouYellow) NotImplemented() {}
func (c RightBehindYouYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type RightBehindYouBlue struct{}

func (RightBehindYouBlue) ID() ids.CardID          { return ids.RightBehindYouBlue }
func (RightBehindYouBlue) Name() string            { return "Right Behind You" }
func (RightBehindYouBlue) Cost(*sim.TurnState) int { return 3 }
func (RightBehindYouBlue) Pitch() int              { return 3 }
func (RightBehindYouBlue) Attack() int             { return 5 }
func (RightBehindYouBlue) Defense() int            { return 2 }
func (RightBehindYouBlue) Types() card.TypeSet     { return rightBehindYouTypes }
func (RightBehindYouBlue) GoAgain() bool           { return false }

// not implemented: defend-together +1{d} buff and deck-bottom peek rider
func (RightBehindYouBlue) NotImplemented() {}
func (c RightBehindYouBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
