// Clap 'Em in Irons — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** When this enters the arena, {t} target Pirate hero or ally. It can't {u}
// while this is in the arena. At the start of your turn, destroy this."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var clapEmInIronsTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type ClapEmInIronsBlue struct{}

func (ClapEmInIronsBlue) ID() card.ID                               { return card.ClapEmInIronsBlue }
func (ClapEmInIronsBlue) Name() string                              { return "Clap 'Em in Irons (Blue)" }
func (ClapEmInIronsBlue) Cost(*card.TurnState) int                  { return 0 }
func (ClapEmInIronsBlue) Pitch() int                                { return 3 }
func (ClapEmInIronsBlue) Attack() int                               { return 0 }
func (ClapEmInIronsBlue) Defense() int                              { return 0 }
func (ClapEmInIronsBlue) Types() card.TypeSet                       { return clapEmInIronsTypes }
func (ClapEmInIronsBlue) GoAgain() bool                             { return true }
// not implemented: passive Pirate-target tap rider; self-destroys at upkeep
func (ClapEmInIronsBlue) NotImplemented()                           {}
func (ClapEmInIronsBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
