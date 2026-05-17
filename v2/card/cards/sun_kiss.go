// Sun Kiss — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
// Printed health-gain: Red 3{h}, Yellow 2{h}, Blue 1{h}.
//
// Text: "Gain N{h}. If you have played a card named Moon Wish this turn, draw a card and Sun Kiss
// gains **go again**." (N is the printed variant value above.)
//
// The synergy is pitch-agnostic: it triggers off any Moon Wish printing in the same turn's
// CardsPlayed.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// sunKissPlay credits the heal as a sub-line under self. When Moon Wish has already played
// this turn, additionally fires an extra mid-turn draw and grants go-again on self.
func sunKissPlay(heal int, ge card.GameEngine, l card.Logger, self *card.CardState) {
	if playedMoonWishThisTurn(ge) {
		ge.DrawOne()
		self.GrantedGoAgain = true
	}
	ge.AddValue(heal)
	l.AppendPostTriggerf(self.Card.DisplayName(), heal, "Gained %d health", heal)
}

// playedMoonWishThisTurn reports whether any prior card resolved this turn is a Moon Wish
// printing. Match on Name() so all three pitch printings count.
func playedMoonWishThisTurn(ge card.GameEngine) bool {
	for _, c := range ge.CardsPlayed() {
		if c.Name() == "Moon Wish" {
			return true
		}
	}
	return false
}

func (SunKissRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	sunKissPlay(3, ge, l, self)
}

func (SunKissYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	sunKissPlay(2, ge, l, self)
}

func (SunKissBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	sunKissPlay(1, ge, l, self)
}
