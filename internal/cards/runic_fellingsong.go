// Runic Fellingsong — Runeblade Action - Attack. Cost 3, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 7, Yellow 6, Blue 5.
// Text: "When this attacks, you may banish an aura from your graveyard. If you do, deal 1 arcane
// damage to target hero."
//
// Play credits Attack() plus 1 arcane when banishAuraFromGraveyard lands an aura in
// the banished zone.
// No aura in the graveyard → the banish rider fizzles and Play returns just Attack().

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// runicFellingsongPlay emits the chain step at printed power, then writes the banish-for-
// arcane rider as a sub-line under self when an aura was successfully banished from the
// graveyard. banishAuraFromGraveyard flips ArcaneDamageDealt internally as part of its
// arcane-damage payload.
func runicFellingsongPlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	if n := banishAuraFromGraveyard(s); n > 0 {
		s.AddValue(n)
		s.LogRider(self, n, "Banished an aura, dealt 1 arcane damage")
	}
}

func (RunicFellingsongRed) Play(s *sim.TurnState, self *sim.CardState) {
	runicFellingsongPlay(s, self)
}

func (RunicFellingsongYellow) Play(s *sim.TurnState, self *sim.CardState) {
	runicFellingsongPlay(s, self)
}

func (RunicFellingsongBlue) Play(s *sim.TurnState, self *sim.CardState) {
	runicFellingsongPlay(s, self)
}
