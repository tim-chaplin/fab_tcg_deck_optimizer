// Fiddler's Green — Generic Block. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 1. Printed on-graveyard gain: Red 3{h}, Yellow 2{h}, Blue 1{h}.
//
// Text: "When this is put into your graveyard from anywhere, gain N{h}." (N is the printed
// variant value above.)
//
// Modelling: using this card to defend sends it to the graveyard, so the N{h} gain fires on
// the DR Play path — credited as +N damage equivalent. Pitched copies go to the bottom of the
// deck instead, so they don't trigger the rider.

package generic

import (
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// fiddlersGreenPlay emits the chain step then writes the printed N{h} as a "Gained N
// health (graveyard trigger)" sub-line under self. Health is valued 1-to-1 with damage.
func fiddlersGreenPlay(s *card.TurnState, self *card.CardState, heal int) {
	s.ApplyAndLogEffectiveDefense(self)
	s.LogRiderOnPlay(self, fmt.Sprintf("Gained %d health (graveyard trigger)", heal), heal)
}

type FiddlersGreenRed struct{}

func (FiddlersGreenRed) ID() card.ID              { return card.FiddlersGreenRed }
func (FiddlersGreenRed) Name() string             { return "Fiddler's Green" }
func (FiddlersGreenRed) Cost(*card.TurnState) int { return 0 }
func (FiddlersGreenRed) Pitch() int               { return 1 }
func (FiddlersGreenRed) Attack() int              { return 0 }
func (FiddlersGreenRed) Defense() int             { return 1 }
func (FiddlersGreenRed) Types() card.TypeSet      { return defenseReactionTypes }
func (FiddlersGreenRed) GoAgain() bool            { return false }
func (FiddlersGreenRed) NotSilverAgeLegal()       {}
func (FiddlersGreenRed) Play(s *card.TurnState, self *card.CardState) {
	fiddlersGreenPlay(s, self, 3)
}

type FiddlersGreenYellow struct{}

func (FiddlersGreenYellow) ID() card.ID              { return card.FiddlersGreenYellow }
func (FiddlersGreenYellow) Name() string             { return "Fiddler's Green" }
func (FiddlersGreenYellow) Cost(*card.TurnState) int { return 0 }
func (FiddlersGreenYellow) Pitch() int               { return 2 }
func (FiddlersGreenYellow) Attack() int              { return 0 }
func (FiddlersGreenYellow) Defense() int             { return 1 }
func (FiddlersGreenYellow) Types() card.TypeSet      { return defenseReactionTypes }
func (FiddlersGreenYellow) GoAgain() bool            { return false }
func (FiddlersGreenYellow) NotSilverAgeLegal()       {}
func (FiddlersGreenYellow) Play(s *card.TurnState, self *card.CardState) {
	fiddlersGreenPlay(s, self, 2)
}

type FiddlersGreenBlue struct{}

func (FiddlersGreenBlue) ID() card.ID              { return card.FiddlersGreenBlue }
func (FiddlersGreenBlue) Name() string             { return "Fiddler's Green" }
func (FiddlersGreenBlue) Cost(*card.TurnState) int { return 0 }
func (FiddlersGreenBlue) Pitch() int               { return 3 }
func (FiddlersGreenBlue) Attack() int              { return 0 }
func (FiddlersGreenBlue) Defense() int             { return 1 }
func (FiddlersGreenBlue) Types() card.TypeSet      { return defenseReactionTypes }
func (FiddlersGreenBlue) GoAgain() bool            { return false }
func (FiddlersGreenBlue) NotSilverAgeLegal()       {}
func (FiddlersGreenBlue) Play(s *card.TurnState, self *card.CardState) {
	fiddlersGreenPlay(s, self, 1)
}
