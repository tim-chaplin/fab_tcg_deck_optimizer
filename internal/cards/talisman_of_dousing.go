// Talisman of Dousing — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** **Spellvoid 1**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var talismanOfDousingTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type TalismanOfDousingYellow struct{}

func (TalismanOfDousingYellow) ID() ids.CardID           { return ids.TalismanOfDousingYellow }
func (TalismanOfDousingYellow) Name() string             { return "Talisman of Dousing" }
func (TalismanOfDousingYellow) Cost(*card.TurnState) int { return 0 }
func (TalismanOfDousingYellow) Pitch() int               { return 2 }
func (TalismanOfDousingYellow) Attack() int              { return 0 }
func (TalismanOfDousingYellow) Defense() int             { return 0 }
func (TalismanOfDousingYellow) Types() card.TypeSet      { return talismanOfDousingTypes }
func (TalismanOfDousingYellow) GoAgain() bool            { return true }

// not implemented: passive Spellvoid 1
func (TalismanOfDousingYellow) NotImplemented()                              {}
func (TalismanOfDousingYellow) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
