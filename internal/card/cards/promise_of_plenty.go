// Promise of Plenty — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Promise of Plenty hits, each hero who doesn't have a card in their arsenal puts
// the top card of their deck face down into their arsenal. If Promise of Plenty is played from
// arsenal, it gains **go again**."
//
// Modelling:
//   - Opponent arsenal isn't tracked; the rider only fires our own side. Our half is a pure
//     deck→arsenal move; the engine accessor (MoveFromTopOfDeckToArsenal) doesn't flip
//     IsCacheable because we never observe the popped card's identity, so the result stays
//     reproducible from the cache key.
//   - "From arsenal: go again" — call self.GrantGoAgainIfFromArsenal() at the top of Play.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func promiseOfPlentyOnHitFillArsenal(ge card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	if !ge.MoveFromTopOfDeckToArsenal() {
		return
	}
	l.AppendPostTrigger(self.Card.DisplayName(), "Filled empty arsenal from top of deck", 0)
}

func promiseOfPlentyPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.GrantGoAgainIfFromArsenal()
	self.RegisterOnHit(promiseOfPlentyOnHitFillArsenal)
}

func (PromiseOfPlentyRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	promiseOfPlentyPlay(ge, l, self)
}

func (PromiseOfPlentyYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	promiseOfPlentyPlay(ge, l, self)
}

func (PromiseOfPlentyBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	promiseOfPlentyPlay(ge, l, self)
}
