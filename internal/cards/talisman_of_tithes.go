// Talisman of Tithes — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** If an opponent would draw 1 or more cards during your action phase, instead
// destroy Talisman of Tithes and they draw that many cards minus 1."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var talismanOfTithesTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type TalismanOfTithesBlue struct{}

func (TalismanOfTithesBlue) ID() ids.CardID          { return ids.TalismanOfTithesBlue }
func (TalismanOfTithesBlue) Name() string            { return "Talisman of Tithes" }
func (TalismanOfTithesBlue) Cost(*sim.TurnState) int { return 0 }
func (TalismanOfTithesBlue) Pitch() int              { return 3 }
func (TalismanOfTithesBlue) Attack() int             { return 0 }
func (TalismanOfTithesBlue) Defense() int            { return 0 }
func (TalismanOfTithesBlue) Types() card.TypeSet     { return talismanOfTithesTypes }
func (TalismanOfTithesBlue) GoAgain() bool           { return true }

// not implemented: self-destroys on an opposing draw during your action phase → opponent draws
// minus 1
func (TalismanOfTithesBlue) NotImplemented()                            {}
func (TalismanOfTithesBlue) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }
