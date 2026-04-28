// Wounding Blow — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 3.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var woundingBlowTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type WoundingBlowRed struct{}

func (WoundingBlowRed) ID() ids.CardID           { return ids.WoundingBlowRed }
func (WoundingBlowRed) Name() string             { return "Wounding Blow" }
func (WoundingBlowRed) Cost(*card.TurnState) int { return 0 }
func (WoundingBlowRed) Pitch() int               { return 1 }
func (WoundingBlowRed) Attack() int              { return 4 }
func (WoundingBlowRed) Defense() int             { return 3 }
func (WoundingBlowRed) Types() card.TypeSet      { return woundingBlowTypes }
func (WoundingBlowRed) GoAgain() bool            { return false }
func (c WoundingBlowRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type WoundingBlowYellow struct{}

func (WoundingBlowYellow) ID() ids.CardID           { return ids.WoundingBlowYellow }
func (WoundingBlowYellow) Name() string             { return "Wounding Blow" }
func (WoundingBlowYellow) Cost(*card.TurnState) int { return 0 }
func (WoundingBlowYellow) Pitch() int               { return 2 }
func (WoundingBlowYellow) Attack() int              { return 3 }
func (WoundingBlowYellow) Defense() int             { return 3 }
func (WoundingBlowYellow) Types() card.TypeSet      { return woundingBlowTypes }
func (WoundingBlowYellow) GoAgain() bool            { return false }
func (c WoundingBlowYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type WoundingBlowBlue struct{}

func (WoundingBlowBlue) ID() ids.CardID           { return ids.WoundingBlowBlue }
func (WoundingBlowBlue) Name() string             { return "Wounding Blow" }
func (WoundingBlowBlue) Cost(*card.TurnState) int { return 0 }
func (WoundingBlowBlue) Pitch() int               { return 3 }
func (WoundingBlowBlue) Attack() int              { return 2 }
func (WoundingBlowBlue) Defense() int             { return 3 }
func (WoundingBlowBlue) Types() card.TypeSet      { return woundingBlowTypes }
func (WoundingBlowBlue) GoAgain() bool            { return false }
func (c WoundingBlowBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
