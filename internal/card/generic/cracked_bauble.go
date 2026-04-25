// Cracked Bauble — Generic Resource. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "*(A player may add any number of Cracked Baubles to their card-pool in sealed deck or
// booster draft formats.)*"
//
// Stub only — marked NotImplemented so the optimizer skips it. The printed effect isn't modelled;
// Play returns 0.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var crackedBaubleTypes = card.NewTypeSet(card.TypeGeneric)

type CrackedBaubleYellow struct{}

func (CrackedBaubleYellow) ID() card.ID                               { return card.CrackedBaubleYellow }
func (CrackedBaubleYellow) Name() string                              { return "Cracked Bauble (Yellow)" }
func (CrackedBaubleYellow) Cost(*card.TurnState) int                  { return 0 }
func (CrackedBaubleYellow) Pitch() int                                { return 2 }
func (CrackedBaubleYellow) Attack() int                               { return 0 }
func (CrackedBaubleYellow) Defense() int                              { return 0 }
func (CrackedBaubleYellow) Types() card.TypeSet                       { return crackedBaubleTypes }
func (CrackedBaubleYellow) GoAgain() bool                             { return false }
func (CrackedBaubleYellow) NotImplemented()                           {}
func (CrackedBaubleYellow) Play(*card.TurnState, *card.CardState) int { return 0 }
