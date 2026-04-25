// Talisman of Balance — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** At the beginning of your end phase, if you have less cards in arsenal than an
// opposing hero, destroy Talisman of Balance and put the top card of your deck into an empty
// arsenal zone you control."
//
// Stub only — marked NotImplemented so the optimizer skips it. The printed effect isn't modelled;
// Play returns 0.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var talismanOfBalanceTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type TalismanOfBalanceBlue struct{}

func (TalismanOfBalanceBlue) ID() card.ID                               { return card.TalismanOfBalanceBlue }
func (TalismanOfBalanceBlue) Name() string                              { return "Talisman of Balance (Blue)" }
func (TalismanOfBalanceBlue) Cost(*card.TurnState) int                  { return 0 }
func (TalismanOfBalanceBlue) Pitch() int                                { return 3 }
func (TalismanOfBalanceBlue) Attack() int                               { return 0 }
func (TalismanOfBalanceBlue) Defense() int                              { return 0 }
func (TalismanOfBalanceBlue) Types() card.TypeSet                       { return talismanOfBalanceTypes }
func (TalismanOfBalanceBlue) GoAgain() bool                             { return true }
func (TalismanOfBalanceBlue) NotImplemented()                           {}
func (TalismanOfBalanceBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
