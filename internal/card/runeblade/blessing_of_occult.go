// Blessing of Occult — Runeblade Action - Aura. Cost 1, Defense 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "At the start of your turn, destroy Blessing of Occult then create N Runechant tokens."
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: assume the aura ticks down and produces its full Runechant payout. Play returns
// N (Red=3, Yellow=2, Blue=1).
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var blessingOfOccultTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type BlessingOfOccultRed struct{}

func (BlessingOfOccultRed) Name() string              { return "Blessing of Occult (Red)" }
func (BlessingOfOccultRed) Cost() int                 { return 1 }
func (BlessingOfOccultRed) Pitch() int                { return 1 }
func (BlessingOfOccultRed) Attack() int               { return 0 }
func (BlessingOfOccultRed) Defense() int              { return 2 }
func (BlessingOfOccultRed) Types() card.TypeSet    { return blessingOfOccultTypes }
func (BlessingOfOccultRed) GoAgain() bool             { return false }
func (BlessingOfOccultRed) Play(*card.TurnState) int  { return 3 }

type BlessingOfOccultYellow struct{}

func (BlessingOfOccultYellow) Name() string             { return "Blessing of Occult (Yellow)" }
func (BlessingOfOccultYellow) Cost() int                { return 1 }
func (BlessingOfOccultYellow) Pitch() int               { return 2 }
func (BlessingOfOccultYellow) Attack() int              { return 0 }
func (BlessingOfOccultYellow) Defense() int             { return 2 }
func (BlessingOfOccultYellow) Types() card.TypeSet   { return blessingOfOccultTypes }
func (BlessingOfOccultYellow) GoAgain() bool            { return false }
func (BlessingOfOccultYellow) Play(*card.TurnState) int { return 2 }

type BlessingOfOccultBlue struct{}

func (BlessingOfOccultBlue) Name() string             { return "Blessing of Occult (Blue)" }
func (BlessingOfOccultBlue) Cost() int                { return 1 }
func (BlessingOfOccultBlue) Pitch() int               { return 3 }
func (BlessingOfOccultBlue) Attack() int              { return 0 }
func (BlessingOfOccultBlue) Defense() int             { return 2 }
func (BlessingOfOccultBlue) Types() card.TypeSet   { return blessingOfOccultTypes }
func (BlessingOfOccultBlue) GoAgain() bool            { return false }
func (BlessingOfOccultBlue) Play(*card.TurnState) int { return 1 }
