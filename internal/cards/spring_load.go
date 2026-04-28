// Spring Load — Generic Action - Attack. Cost 1. Printed power: Red 2, Yellow 2, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, if you have no cards in hand, it gains +3{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var springLoadTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type SpringLoadRed struct{}

func (SpringLoadRed) ID() ids.CardID          { return ids.SpringLoadRed }
func (SpringLoadRed) Name() string            { return "Spring Load" }
func (SpringLoadRed) Cost(*sim.TurnState) int { return 1 }
func (SpringLoadRed) Pitch() int              { return 1 }
func (SpringLoadRed) Attack() int             { return 2 }
func (SpringLoadRed) Defense() int            { return 2 }
func (SpringLoadRed) Types() card.TypeSet     { return springLoadTypes }
func (SpringLoadRed) GoAgain() bool           { return false }

// not implemented: +3{p} 'no cards in hand' rider never fires
func (SpringLoadRed) NotImplemented() {}
func (c SpringLoadRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type SpringLoadYellow struct{}

func (SpringLoadYellow) ID() ids.CardID          { return ids.SpringLoadYellow }
func (SpringLoadYellow) Name() string            { return "Spring Load" }
func (SpringLoadYellow) Cost(*sim.TurnState) int { return 1 }
func (SpringLoadYellow) Pitch() int              { return 2 }
func (SpringLoadYellow) Attack() int             { return 2 }
func (SpringLoadYellow) Defense() int            { return 2 }
func (SpringLoadYellow) Types() card.TypeSet     { return springLoadTypes }
func (SpringLoadYellow) GoAgain() bool           { return false }

// not implemented: +3{p} 'no cards in hand' rider never fires
func (SpringLoadYellow) NotImplemented() {}
func (c SpringLoadYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type SpringLoadBlue struct{}

func (SpringLoadBlue) ID() ids.CardID          { return ids.SpringLoadBlue }
func (SpringLoadBlue) Name() string            { return "Spring Load" }
func (SpringLoadBlue) Cost(*sim.TurnState) int { return 1 }
func (SpringLoadBlue) Pitch() int              { return 3 }
func (SpringLoadBlue) Attack() int             { return 2 }
func (SpringLoadBlue) Defense() int            { return 2 }
func (SpringLoadBlue) Types() card.TypeSet     { return springLoadTypes }
func (SpringLoadBlue) GoAgain() bool           { return false }

// not implemented: +3{p} 'no cards in hand' rider never fires
func (SpringLoadBlue) NotImplemented() {}
func (c SpringLoadBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
