// Blessing of Occult — Runeblade Action - Aura. Cost 1, Defense 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "At the start of your turn, destroy Blessing of Occult then create N Runechant tokens."
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Handler creates N Runechants next turn.

package runeblade

import (
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

var blessingOfOccultTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type BlessingOfOccultRed struct{}

func (BlessingOfOccultRed) ID() card.ID              { return card.BlessingOfOccultRed }
func (BlessingOfOccultRed) Name() string             { return "Blessing of Occult" }
func (BlessingOfOccultRed) Cost(*card.TurnState) int { return 1 }
func (BlessingOfOccultRed) Pitch() int               { return 1 }
func (BlessingOfOccultRed) Attack() int              { return 0 }
func (BlessingOfOccultRed) Defense() int             { return 2 }
func (BlessingOfOccultRed) Types() card.TypeSet      { return blessingOfOccultTypes }
func (BlessingOfOccultRed) GoAgain() bool            { return false }
func (BlessingOfOccultRed) AddsFutureValue()         {}
func (c BlessingOfOccultRed) Play(s *card.TurnState, self *card.CardState) {
	blessingOfOccultPlay(s, self, c, 3)
}

type BlessingOfOccultYellow struct{}

func (BlessingOfOccultYellow) ID() card.ID              { return card.BlessingOfOccultYellow }
func (BlessingOfOccultYellow) Name() string             { return "Blessing of Occult" }
func (BlessingOfOccultYellow) Cost(*card.TurnState) int { return 1 }
func (BlessingOfOccultYellow) Pitch() int               { return 2 }
func (BlessingOfOccultYellow) Attack() int              { return 0 }
func (BlessingOfOccultYellow) Defense() int             { return 2 }
func (BlessingOfOccultYellow) Types() card.TypeSet      { return blessingOfOccultTypes }
func (BlessingOfOccultYellow) GoAgain() bool            { return false }
func (BlessingOfOccultYellow) AddsFutureValue()         {}
func (c BlessingOfOccultYellow) Play(s *card.TurnState, self *card.CardState) {
	blessingOfOccultPlay(s, self, c, 2)
}

type BlessingOfOccultBlue struct{}

func (BlessingOfOccultBlue) ID() card.ID              { return card.BlessingOfOccultBlue }
func (BlessingOfOccultBlue) Name() string             { return "Blessing of Occult" }
func (BlessingOfOccultBlue) Cost(*card.TurnState) int { return 1 }
func (BlessingOfOccultBlue) Pitch() int               { return 3 }
func (BlessingOfOccultBlue) Attack() int              { return 0 }
func (BlessingOfOccultBlue) Defense() int             { return 2 }
func (BlessingOfOccultBlue) Types() card.TypeSet      { return blessingOfOccultTypes }
func (BlessingOfOccultBlue) GoAgain() bool            { return false }
func (BlessingOfOccultBlue) AddsFutureValue()         {}
func (c BlessingOfOccultBlue) Play(s *card.TurnState, self *card.CardState) {
	blessingOfOccultPlay(s, self, c, 1)
}

// blessingOfOccultPlay registers the shared next-turn trigger that creates n Runechants
// and emits the same-turn chain step (no value contribution; all credit is deferred to
// the trigger).
func blessingOfOccultPlay(s *card.TurnState, selfState *card.CardState, selfCard card.Card, n int) {
	text := "Created a runechant"
	if n > 1 {
		text = fmt.Sprintf("Created %d runechants", n)
	}
	s.RegisterStartOfTurn(selfCard, 1, text, func(s *card.TurnState) int { return s.CreateRunechants(n) })
	s.LogPlay(selfState)
}
