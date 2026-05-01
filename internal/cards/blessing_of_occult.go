// Blessing of Occult — Runeblade Action - Aura. Cost 1, Defense 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "At the start of your turn, destroy Blessing of Occult then create N Runechant tokens."
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Handler creates N Runechants next turn.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var blessingOfOccultTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

// blessingOfOccultTriggerText pre-formats the trigger log text for each Runechant count
// (1 = Blue, 2 = Yellow, 3 = Red). The text is captured into a per-play closure on every
// cast, so a constant lookup avoids a Sprintf alloc per cast.
var blessingOfOccultTriggerText = [...]string{
	1: "Created a runechant",
	2: "Created 2 runechants",
	3: "Created 3 runechants",
}

type BlessingOfOccultRed struct{}

func (BlessingOfOccultRed) ID() ids.CardID          { return ids.BlessingOfOccultRed }
func (BlessingOfOccultRed) Name() string            { return "Blessing of Occult" }
func (BlessingOfOccultRed) Cost(*sim.TurnState) int { return 1 }
func (BlessingOfOccultRed) Pitch() int              { return 1 }
func (BlessingOfOccultRed) Attack() int             { return 0 }
func (BlessingOfOccultRed) Defense() int            { return 2 }
func (BlessingOfOccultRed) Types() card.TypeSet     { return blessingOfOccultTypes }
func (BlessingOfOccultRed) GoAgain() bool           { return false }
func (BlessingOfOccultRed) AddsFutureValue()        {}
func (c BlessingOfOccultRed) Play(s *sim.TurnState, self *sim.CardState) {
	blessingOfOccultPlay(s, self, c, 3)
}

type BlessingOfOccultYellow struct{}

func (BlessingOfOccultYellow) ID() ids.CardID          { return ids.BlessingOfOccultYellow }
func (BlessingOfOccultYellow) Name() string            { return "Blessing of Occult" }
func (BlessingOfOccultYellow) Cost(*sim.TurnState) int { return 1 }
func (BlessingOfOccultYellow) Pitch() int              { return 2 }
func (BlessingOfOccultYellow) Attack() int             { return 0 }
func (BlessingOfOccultYellow) Defense() int            { return 2 }
func (BlessingOfOccultYellow) Types() card.TypeSet     { return blessingOfOccultTypes }
func (BlessingOfOccultYellow) GoAgain() bool           { return false }
func (BlessingOfOccultYellow) AddsFutureValue()        {}
func (c BlessingOfOccultYellow) Play(s *sim.TurnState, self *sim.CardState) {
	blessingOfOccultPlay(s, self, c, 2)
}

type BlessingOfOccultBlue struct{}

func (BlessingOfOccultBlue) ID() ids.CardID          { return ids.BlessingOfOccultBlue }
func (BlessingOfOccultBlue) Name() string            { return "Blessing of Occult" }
func (BlessingOfOccultBlue) Cost(*sim.TurnState) int { return 1 }
func (BlessingOfOccultBlue) Pitch() int              { return 3 }
func (BlessingOfOccultBlue) Attack() int             { return 0 }
func (BlessingOfOccultBlue) Defense() int            { return 2 }
func (BlessingOfOccultBlue) Types() card.TypeSet     { return blessingOfOccultTypes }
func (BlessingOfOccultBlue) GoAgain() bool           { return false }
func (BlessingOfOccultBlue) AddsFutureValue()        {}
func (c BlessingOfOccultBlue) Play(s *sim.TurnState, self *sim.CardState) {
	blessingOfOccultPlay(s, self, c, 1)
}

// blessingOfOccultPlay registers the shared next-turn trigger that creates n Runechants
// and emits the same-turn chain step (no value contribution; all credit is deferred to
// the trigger).
func blessingOfOccultPlay(s *sim.TurnState, selfState *sim.CardState, selfCard sim.Card, n int) {
	s.RegisterStartOfTurn(selfCard, 1, blessingOfOccultTriggerText[n], func(s *sim.TurnState) int { return s.CreateRunechants(n) })
	s.LogPlay(selfState)
}
