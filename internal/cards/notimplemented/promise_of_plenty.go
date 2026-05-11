// Promise of Plenty — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Promise of Plenty hits, each hero who doesn't have a card in their arsenal puts the top
// card of their deck face down into their arsenal. If Promise of Plenty is played from arsenal, it
// gains **go again**."
//
// Modelling: the arsenal-placement rider isn't modelled (arsenal/deck content tracking would
// be required). Standard played-from-arsenal go-again (docs/dev-standards.md).

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func promiseOfPlentyPlay(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.GrantGoAgainIfFromArsenal()
}

// not implemented: on-hit arsenal-placement rider (arsenal/deck content tracking would be required)

func (PromiseOfPlentyRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	promiseOfPlentyPlay(s, l, self)
}

// not implemented: on-hit arsenal-placement rider (arsenal/deck content tracking would be required)

func (PromiseOfPlentyYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	promiseOfPlentyPlay(s, l, self)
}

// not implemented: on-hit arsenal-placement rider (arsenal/deck content tracking would be required)

func (PromiseOfPlentyBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	promiseOfPlentyPlay(s, l, self)
}
