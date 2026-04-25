// Potion of Ironhide — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy Potion of Ironhide: Attack action cards you own gain +1{d} this
// turn."
//
// Stub only — marked NotImplemented so the optimizer skips it. The printed effect isn't modelled;
// Play returns 0.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var potionOfIronhideTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type PotionOfIronhideBlue struct{}

func (PotionOfIronhideBlue) ID() card.ID                               { return card.PotionOfIronhideBlue }
func (PotionOfIronhideBlue) Name() string                              { return "Potion of Ironhide (Blue)" }
func (PotionOfIronhideBlue) Cost(*card.TurnState) int                  { return 0 }
func (PotionOfIronhideBlue) Pitch() int                                { return 3 }
func (PotionOfIronhideBlue) Attack() int                               { return 0 }
func (PotionOfIronhideBlue) Defense() int                              { return 0 }
func (PotionOfIronhideBlue) Types() card.TypeSet                       { return potionOfIronhideTypes }
func (PotionOfIronhideBlue) GoAgain() bool                             { return false }
func (PotionOfIronhideBlue) NotImplemented()                           {}
func (PotionOfIronhideBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
