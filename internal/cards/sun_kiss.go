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
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// sunKissPlay emits the chain step ("Sun Kiss [R]: PLAY"), writes the heal as a sub-line
// "Gained N health" under it, and — when Moon Wish has already played this turn — fires
// an extra mid-turn draw and a go-again grant on self.
func sunKissPlay(heal int, s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	if playedMoonWishThisTurn(s) {
		s.DrawOne()
		self.GrantedGoAgain = true
	}
	s.AddValue(heal)
	l.AppendPostTriggerf(self.Card.DisplayName(), heal, "Gained %d health", heal)
}

// playedMoonWishThisTurn reports whether any prior card resolved this turn is a Moon Wish
// printing. Exact-match on Name() works because all three Moon Wish printings share the
// base name "Moon Wish" — the pitch suffix lives in DisplayName, not Name.
func playedMoonWishThisTurn(s sim.GameEngine) bool {
	for _, c := range s.CardsPlayed() {
		if c.Name() == "Moon Wish" {
			return true
		}
	}
	return false
}

func (SunKissRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	sunKissPlay(3, s, l, self)
}

func (SunKissYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	sunKissPlay(2, s, l, self)
}

func (SunKissBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	sunKissPlay(1, s, l, self)
}
