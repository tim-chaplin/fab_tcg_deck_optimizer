// Clap 'Em in Irons — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** When this enters the arena, {t} target Pirate hero or ally. It can't {u}
// while this is in the arena. At the start of your turn, destroy this."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var clapEmInIronsTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type ClapEmInIronsBlue struct{}

func (ClapEmInIronsBlue) ID() ids.CardID           { return ids.ClapEmInIronsBlue }
func (ClapEmInIronsBlue) Name() string             { return "Clap 'Em in Irons" }
func (ClapEmInIronsBlue) Cost(*card.TurnState) int { return 0 }
func (ClapEmInIronsBlue) Pitch() int               { return 3 }
func (ClapEmInIronsBlue) Attack() int              { return 0 }
func (ClapEmInIronsBlue) Defense() int             { return 0 }
func (ClapEmInIronsBlue) Types() card.TypeSet      { return clapEmInIronsTypes }
func (ClapEmInIronsBlue) GoAgain() bool            { return true }

// not implemented: passive tap-target Pirate; can't unfreeze; self-destroys at start of turn
func (ClapEmInIronsBlue) NotImplemented()                              {}
func (ClapEmInIronsBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
