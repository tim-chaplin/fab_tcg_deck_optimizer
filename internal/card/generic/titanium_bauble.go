// Titanium Bauble — Generic Resource. Cost 0. Printed pitch variants: Blue 3. Defense 3.
//
// Stub only — marked NotImplemented so the optimizer skips it. The printed effect isn't modelled;
// Play returns 0.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var titaniumBaubleTypes = card.NewTypeSet(card.TypeGeneric)

type TitaniumBaubleBlue struct{}

func (TitaniumBaubleBlue) ID() card.ID                               { return card.TitaniumBaubleBlue }
func (TitaniumBaubleBlue) Name() string                              { return "Titanium Bauble (Blue)" }
func (TitaniumBaubleBlue) Cost(*card.TurnState) int                  { return 0 }
func (TitaniumBaubleBlue) Pitch() int                                { return 3 }
func (TitaniumBaubleBlue) Attack() int                               { return 0 }
func (TitaniumBaubleBlue) Defense() int                              { return 3 }
func (TitaniumBaubleBlue) Types() card.TypeSet                       { return titaniumBaubleTypes }
func (TitaniumBaubleBlue) GoAgain() bool                             { return false }
func (TitaniumBaubleBlue) NotImplemented()                           {}
func (TitaniumBaubleBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
