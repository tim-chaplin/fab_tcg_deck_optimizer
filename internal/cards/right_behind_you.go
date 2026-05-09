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
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

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

func (RightBehindYouRed) Block(s *sim.TurnState, self *sim.CardState) {
	rightBehindYouBlock(s, self)
}
func (RightBehindYouRed) Play(s *sim.TurnState, self *sim.CardState) {
	rightBehindYouPlay(s, self)
}

func (RightBehindYouYellow) Block(s *sim.TurnState, self *sim.CardState) {
	rightBehindYouBlock(s, self)
}
func (RightBehindYouYellow) Play(s *sim.TurnState, self *sim.CardState) {
	rightBehindYouPlay(s, self)
}

func (RightBehindYouBlue) Block(s *sim.TurnState, self *sim.CardState) {
	rightBehindYouBlock(s, self)
}
func (RightBehindYouBlue) Play(s *sim.TurnState, self *sim.CardState) {
	rightBehindYouPlay(s, self)
}
