// Titanium Bauble — Generic Resource. Cost 0. Printed pitch variants: Blue 3. Defense 3.

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var titaniumBaubleTypes = card.NewTypeSet(card.TypeGeneric)

type TitaniumBaubleBlue struct{}

func (TitaniumBaubleBlue) ID() card.ID              { return card.TitaniumBaubleBlue }
func (TitaniumBaubleBlue) Name() string             { return "Titanium Bauble" }
func (TitaniumBaubleBlue) Cost(*card.TurnState) int { return 0 }
func (TitaniumBaubleBlue) Pitch() int               { return 3 }
func (TitaniumBaubleBlue) Attack() int              { return 0 }
func (TitaniumBaubleBlue) Defense() int             { return 3 }
func (TitaniumBaubleBlue) Types() card.TypeSet      { return titaniumBaubleTypes }
func (TitaniumBaubleBlue) GoAgain() bool            { return false }

// not implemented: pitch-3 resource with 3{d}; no other effect
func (TitaniumBaubleBlue) NotImplemented()                              {}
func (TitaniumBaubleBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
