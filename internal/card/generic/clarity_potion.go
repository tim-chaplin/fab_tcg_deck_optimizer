// Clarity Potion — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy Clarity Potion: **Opt 2**"
//
// Stub only — marked NotImplemented so the optimizer skips it. The printed effect isn't modelled;
// Play returns 0.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var clarityPotionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type ClarityPotionBlue struct{}

func (ClarityPotionBlue) ID() card.ID                               { return card.ClarityPotionBlue }
func (ClarityPotionBlue) Name() string                              { return "Clarity Potion (Blue)" }
func (ClarityPotionBlue) Cost(*card.TurnState) int                  { return 0 }
func (ClarityPotionBlue) Pitch() int                                { return 3 }
func (ClarityPotionBlue) Attack() int                               { return 0 }
func (ClarityPotionBlue) Defense() int                              { return 0 }
func (ClarityPotionBlue) Types() card.TypeSet                       { return clarityPotionTypes }
func (ClarityPotionBlue) GoAgain() bool                             { return false }
func (ClarityPotionBlue) NotImplemented()                           {}
func (ClarityPotionBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
