// Even Bigger Than That! — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2,
// Blue 3.
//
// Text: "Play Even Bigger Than That! only if you've dealt {p} this turn. **Opt 3**, then reveal
// the top card of your deck. If it has {p} greater than the amount of damage you've dealt this
// turn, create a Quicken token and draw a card."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (EvenBiggerThanThatRed) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return ge.HitThisTurn()
}

func (EvenBiggerThanThatYellow) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return ge.HitThisTurn()
}

func (EvenBiggerThanThatBlue) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return ge.HitThisTurn()
}

func evenBiggerThanThatPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	dealt := ge.DamageDealt()
	ge.OptWith(l, 3, func(cs []card.Card) (top, bottom []card.Card) {
		for _, c := range cs {
			if c.Attack() > dealt {
				top = append(top, c)
			} else {
				bottom = append(bottom, c)
			}
		}
		return
	})
	revealed, ok := ge.PeekDeck()
	if !ok {
		return
	}
	if revealed.Attack() <= dealt {
		return
	}
	ge.CreateQuicken(1)
	ge.DrawOne()
	l.AppendPostTriggerf(self.Card.DisplayName(), 0,
		"Top card %s power %d > dealt %d — created Quicken and drew",
		revealed.DisplayName(), revealed.Attack(), dealt)
}

func (EvenBiggerThanThatRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	evenBiggerThanThatPlay(ge, l, self)
}

func (EvenBiggerThanThatYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	evenBiggerThanThatPlay(ge, l, self)
}

func (EvenBiggerThanThatBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	evenBiggerThanThatPlay(ge, l, self)
}
