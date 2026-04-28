// Brutal Assault — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 3.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var brutalAssaultTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BrutalAssaultRed struct{}

func (BrutalAssaultRed) ID() ids.CardID          { return ids.BrutalAssaultRed }
func (BrutalAssaultRed) Name() string            { return "Brutal Assault" }
func (BrutalAssaultRed) Cost(*sim.TurnState) int { return 2 }
func (BrutalAssaultRed) Pitch() int              { return 1 }
func (BrutalAssaultRed) Attack() int             { return 6 }
func (BrutalAssaultRed) Defense() int            { return 3 }
func (BrutalAssaultRed) Types() card.TypeSet     { return brutalAssaultTypes }
func (BrutalAssaultRed) GoAgain() bool           { return false }
func (c BrutalAssaultRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type BrutalAssaultYellow struct{}

func (BrutalAssaultYellow) ID() ids.CardID          { return ids.BrutalAssaultYellow }
func (BrutalAssaultYellow) Name() string            { return "Brutal Assault" }
func (BrutalAssaultYellow) Cost(*sim.TurnState) int { return 2 }
func (BrutalAssaultYellow) Pitch() int              { return 2 }
func (BrutalAssaultYellow) Attack() int             { return 5 }
func (BrutalAssaultYellow) Defense() int            { return 3 }
func (BrutalAssaultYellow) Types() card.TypeSet     { return brutalAssaultTypes }
func (BrutalAssaultYellow) GoAgain() bool           { return false }
func (c BrutalAssaultYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type BrutalAssaultBlue struct{}

func (BrutalAssaultBlue) ID() ids.CardID          { return ids.BrutalAssaultBlue }
func (BrutalAssaultBlue) Name() string            { return "Brutal Assault" }
func (BrutalAssaultBlue) Cost(*sim.TurnState) int { return 2 }
func (BrutalAssaultBlue) Pitch() int              { return 3 }
func (BrutalAssaultBlue) Attack() int             { return 4 }
func (BrutalAssaultBlue) Defense() int            { return 3 }
func (BrutalAssaultBlue) Types() card.TypeSet     { return brutalAssaultTypes }
func (BrutalAssaultBlue) GoAgain() bool           { return false }
func (c BrutalAssaultBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
