// Sirens of Safe Harbor — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this is put into your graveyard from anywhere, gain 1{h}."
//
// Modelling: the card hits the graveyard after resolving as an attack, so the 1{h} gain fires
// on every Play — credited as +1 damage equivalent. Pitched copies go to the bottom of the
// deck instead of the graveyard, so they don't trigger the rider.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (SirensOfSafeHarborRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
	s.AddValue(1)
	l.LogRider(self, 1, "Gained 1 health (graveyard trigger)")
}

func (SirensOfSafeHarborYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
	s.AddValue(1)
	l.LogRider(self, 1, "Gained 1 health (graveyard trigger)")
}

func (SirensOfSafeHarborBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
	s.AddValue(1)
	l.LogRider(self, 1, "Gained 1 health (graveyard trigger)")
}
