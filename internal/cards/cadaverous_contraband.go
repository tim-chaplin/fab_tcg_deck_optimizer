// Cadaverous Contraband — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Cadaverous Contraband hits, you may put a 'non-attack' action card from your graveyard
// on top of your deck."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var cadaverousContrabandTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type CadaverousContrabandRed struct{}

func (CadaverousContrabandRed) ID() ids.CardID           { return ids.CadaverousContrabandRed }
func (CadaverousContrabandRed) Name() string             { return "Cadaverous Contraband" }
func (CadaverousContrabandRed) Cost(*card.TurnState) int { return 2 }
func (CadaverousContrabandRed) Pitch() int               { return 1 }
func (CadaverousContrabandRed) Attack() int              { return 6 }
func (CadaverousContrabandRed) Defense() int             { return 2 }
func (CadaverousContrabandRed) Types() card.TypeSet      { return cadaverousContrabandTypes }
func (CadaverousContrabandRed) GoAgain() bool            { return false }

// not implemented: on-hit graveyard → top-of-deck rider
func (CadaverousContrabandRed) NotImplemented() {}
func (c CadaverousContrabandRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type CadaverousContrabandYellow struct{}

func (CadaverousContrabandYellow) ID() ids.CardID           { return ids.CadaverousContrabandYellow }
func (CadaverousContrabandYellow) Name() string             { return "Cadaverous Contraband" }
func (CadaverousContrabandYellow) Cost(*card.TurnState) int { return 2 }
func (CadaverousContrabandYellow) Pitch() int               { return 2 }
func (CadaverousContrabandYellow) Attack() int              { return 5 }
func (CadaverousContrabandYellow) Defense() int             { return 2 }
func (CadaverousContrabandYellow) Types() card.TypeSet      { return cadaverousContrabandTypes }
func (CadaverousContrabandYellow) GoAgain() bool            { return false }

// not implemented: on-hit graveyard → top-of-deck rider
func (CadaverousContrabandYellow) NotImplemented() {}
func (c CadaverousContrabandYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type CadaverousContrabandBlue struct{}

func (CadaverousContrabandBlue) ID() ids.CardID           { return ids.CadaverousContrabandBlue }
func (CadaverousContrabandBlue) Name() string             { return "Cadaverous Contraband" }
func (CadaverousContrabandBlue) Cost(*card.TurnState) int { return 2 }
func (CadaverousContrabandBlue) Pitch() int               { return 3 }
func (CadaverousContrabandBlue) Attack() int              { return 4 }
func (CadaverousContrabandBlue) Defense() int             { return 2 }
func (CadaverousContrabandBlue) Types() card.TypeSet      { return cadaverousContrabandTypes }
func (CadaverousContrabandBlue) GoAgain() bool            { return false }

// not implemented: on-hit graveyard → top-of-deck rider
func (CadaverousContrabandBlue) NotImplemented() {}
func (c CadaverousContrabandBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
