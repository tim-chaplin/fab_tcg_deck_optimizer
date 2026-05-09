// Warmonger's Recital — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next attack action card you play this turn gains +N{p} and "When this hits, put it on
// the bottom of its owner's deck." **Go again**" (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// warmongersRecitalRecycleOnHit pulls the buffed attack out of graveyard onto the deck bottom.
func warmongersRecitalRecycleOnHit(s *sim.TurnState, self *sim.CardState, _ *sim.OnHitHandler) {
	target := self.Card
	if _, ok := s.RecycleFromGraveyardToBottom(func(c sim.Card) bool { return c == target }); !ok {
		return
	}
	s.LogPostTrigger(self.Card.DisplayName(), "Recycled to bottom of deck on hit", 0)
}

// warmongersRecitalPlay grants the next attack action +n{p} and the on-hit recycle rider.
// Fizzles silently if no attack action follows in CardsRemaining.
func warmongersRecitalPlay(s *sim.TurnState, self *sim.CardState, source sim.Card, n int) {
	for _, pc := range s.CardsRemaining {
		if pc.Card.Types().IsAttackAction() {
			pc.BonusAttack += n
			pc.OnHit = append(pc.OnHit, sim.OnHitHandler{
				Fire:   warmongersRecitalRecycleOnHit,
				Source: source,
			})
			break
		}
	}
	n2 := self.DealEffectiveAttack(s)
	s.Log(self, n2)
}

func (c WarmongersRecitalRed) Play(s *sim.TurnState, self *sim.CardState) {
	warmongersRecitalPlay(s, self, c, 3)
}

func (c WarmongersRecitalYellow) Play(s *sim.TurnState, self *sim.CardState) {
	warmongersRecitalPlay(s, self, c, 2)
}

func (c WarmongersRecitalBlue) Play(s *sim.TurnState, self *sim.CardState) {
	warmongersRecitalPlay(s, self, c, 1)
}
