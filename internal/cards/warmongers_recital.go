// Warmonger's Recital — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next attack action card you play this turn gains +N{p} and "When this hits, put it on
// the bottom of its owner's deck." **Go again**" (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// warmongersRecitalRecycleOnHit pulls the buffed attack out of graveyard onto the deck bottom.
func warmongersRecitalRecycleOnHit(s card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	target := self.Card
	if _, ok := s.RecycleFromGraveyardToBottom(func(c card.Card) bool { return c == target }); !ok {
		return
	}
	l.AppendPostTrigger(self.Card.DisplayName(), "Recycled to bottom of deck on hit", 0)
}

// warmongersRecitalPlay grants the next attack action +n{p} and the on-hit recycle rider.
// Fizzles silently if no attack action follows in CardsRemaining.
func warmongersRecitalPlay(s card.GameEngine, l card.Logger, self *card.CardState, source card.Card, n int) {
	for _, pc := range s.CardsRemaining() {
		if pc.Card.Types().IsAttackAction() {
			pc.BonusAttack += n
			pc.OnHit = append(pc.OnHit, card.OnHitHandler{
				Fire:   warmongersRecitalRecycleOnHit,
				Source: source,
			})
			break
		}
	}
}

func (c WarmongersRecitalRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	warmongersRecitalPlay(s, l, self, c, 3)
}

func (c WarmongersRecitalYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	warmongersRecitalPlay(s, l, self, c, 2)
}

func (c WarmongersRecitalBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	warmongersRecitalPlay(s, l, self, c, 1)
}
