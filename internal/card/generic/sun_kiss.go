// Sun Kiss — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
// Printed health-gain: Red 3{h}, Yellow 2{h}, Blue 1{h}.
//
// Text: "Gain N{h}. If you have played a card named Moon Wish this turn, draw a card and Sun Kiss
// gains **go again**." (N is the printed variant value above.)
//
// The synergy is pitch-agnostic: it triggers off any Moon Wish printing in the same turn's
// CardsPlayed.

package generic

import (
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

var sunKissTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// sunKissPlay emits the chain step ("Sun Kiss [R]: PLAY"), writes the heal as a sub-line
// "Gained N health" under it, and — when Moon Wish has already played this turn — fires
// an extra mid-turn draw and a go-again grant on self.
func sunKissPlay(heal int, s *card.TurnState, self *card.CardState) {
	if playedMoonWishThisTurn(s) {
		s.DrawOne()
		self.GrantedGoAgain = true
	}
	s.LogPlay(self)
	s.ApplyAndLogRiderOnPlay(self, fmt.Sprintf("Gained %d health", heal), heal)
}

// playedMoonWishThisTurn reports whether any prior card resolved this turn was a Moon Wish
// printing. Exact-match on Name() works because all three Moon Wish printings share the
// base name "Moon Wish" — the pitch suffix lives in DisplayName, not Name.
func playedMoonWishThisTurn(s *card.TurnState) bool {
	for _, c := range s.CardsPlayed {
		if c.Name() == "Moon Wish" {
			return true
		}
	}
	return false
}

type SunKissRed struct{}

func (SunKissRed) ID() card.ID              { return card.SunKissRed }
func (SunKissRed) Name() string             { return "Sun Kiss" }
func (SunKissRed) Cost(*card.TurnState) int { return 0 }
func (SunKissRed) Pitch() int               { return 1 }
func (SunKissRed) Attack() int              { return 0 }
func (SunKissRed) Defense() int             { return 2 }
func (SunKissRed) Types() card.TypeSet      { return sunKissTypes }
func (SunKissRed) GoAgain() bool            { return false }
func (SunKissRed) Play(s *card.TurnState, self *card.CardState) {
	sunKissPlay(3, s, self)
}

type SunKissYellow struct{}

func (SunKissYellow) ID() card.ID              { return card.SunKissYellow }
func (SunKissYellow) Name() string             { return "Sun Kiss" }
func (SunKissYellow) Cost(*card.TurnState) int { return 0 }
func (SunKissYellow) Pitch() int               { return 2 }
func (SunKissYellow) Attack() int              { return 0 }
func (SunKissYellow) Defense() int             { return 2 }
func (SunKissYellow) Types() card.TypeSet      { return sunKissTypes }
func (SunKissYellow) GoAgain() bool            { return false }
func (SunKissYellow) Play(s *card.TurnState, self *card.CardState) {
	sunKissPlay(2, s, self)
}

type SunKissBlue struct{}

func (SunKissBlue) ID() card.ID              { return card.SunKissBlue }
func (SunKissBlue) Name() string             { return "Sun Kiss" }
func (SunKissBlue) Cost(*card.TurnState) int { return 0 }
func (SunKissBlue) Pitch() int               { return 3 }
func (SunKissBlue) Attack() int              { return 0 }
func (SunKissBlue) Defense() int             { return 2 }
func (SunKissBlue) Types() card.TypeSet      { return sunKissTypes }
func (SunKissBlue) GoAgain() bool            { return false }
func (SunKissBlue) Play(s *card.TurnState, self *card.CardState) {
	sunKissPlay(1, s, self)
}
