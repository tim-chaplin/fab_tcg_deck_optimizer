// Right Behind You — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this defends together with another card from hand, this gets +1{d} and you may look
// at the top card of your deck. You may put it on the bottom."
//
// Block scans s.Defenders for a second plain blocker; if at least one is present (DRs
// alongside don't count), the +1{d} fires by bumping self.BonusDefense. The deck-top
// peek/bottom rider is dropped — the optimizer would never bottom a card it'd rather
// draw next.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var rightBehindYouTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func rightBehindYouPlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

// rightBehindYouBlock fires the +1{d} together-bonus when at least two plain blockers
// share the defenders slot. Short-circuits on the second non-DR sighting.
func rightBehindYouBlock(s *sim.TurnState, self *sim.CardState) {
	plainCount := 0
	for _, d := range s.Defenders {
		if d.Types().IsDefenseReaction() {
			continue
		}
		plainCount++
		if plainCount >= 2 {
			self.BonusDefense += 1
			return
		}
	}
}

type RightBehindYouRed struct{}

func (RightBehindYouRed) ID() ids.CardID          { return ids.RightBehindYouRed }
func (RightBehindYouRed) Name() string            { return "Right Behind You" }
func (RightBehindYouRed) Cost(*sim.TurnState) int { return 3 }
func (RightBehindYouRed) Pitch() int              { return 1 }
func (RightBehindYouRed) Attack() int             { return 7 }
func (RightBehindYouRed) Defense() int            { return 2 }
func (RightBehindYouRed) Types() card.TypeSet     { return rightBehindYouTypes }
func (RightBehindYouRed) GoAgain() bool           { return false }
func (RightBehindYouRed) Block(s *sim.TurnState, self *sim.CardState) {
	rightBehindYouBlock(s, self)
}
func (RightBehindYouRed) Play(s *sim.TurnState, self *sim.CardState) {
	rightBehindYouPlay(s, self)
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
func (RightBehindYouYellow) Block(s *sim.TurnState, self *sim.CardState) {
	rightBehindYouBlock(s, self)
}
func (RightBehindYouYellow) Play(s *sim.TurnState, self *sim.CardState) {
	rightBehindYouPlay(s, self)
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
func (RightBehindYouBlue) Block(s *sim.TurnState, self *sim.CardState) {
	rightBehindYouBlock(s, self)
}
func (RightBehindYouBlue) Play(s *sim.TurnState, self *sim.CardState) {
	rightBehindYouPlay(s, self)
}
