// Looking for a Scrap — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Looking for a Scrap, you may banish a card with 1{p} from
// your graveyard. When you do, this gains +1{p} and **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var lookingForAScrapTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type LookingForAScrapRed struct{}

func (LookingForAScrapRed) ID() ids.CardID          { return ids.LookingForAScrapRed }
func (LookingForAScrapRed) Name() string            { return "Looking for a Scrap" }
func (LookingForAScrapRed) Cost(*sim.TurnState) int { return 1 }
func (LookingForAScrapRed) Pitch() int              { return 1 }
func (LookingForAScrapRed) Attack() int             { return 4 }
func (LookingForAScrapRed) Defense() int            { return 2 }
func (LookingForAScrapRed) Types() card.TypeSet     { return lookingForAScrapTypes }
func (LookingForAScrapRed) GoAgain() bool           { return false }

// not implemented: graveyard-banish additional cost and the +1{p}/go-again bonus rider
func (LookingForAScrapRed) NotImplemented() {}
func (c LookingForAScrapRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type LookingForAScrapYellow struct{}

func (LookingForAScrapYellow) ID() ids.CardID          { return ids.LookingForAScrapYellow }
func (LookingForAScrapYellow) Name() string            { return "Looking for a Scrap" }
func (LookingForAScrapYellow) Cost(*sim.TurnState) int { return 1 }
func (LookingForAScrapYellow) Pitch() int              { return 2 }
func (LookingForAScrapYellow) Attack() int             { return 3 }
func (LookingForAScrapYellow) Defense() int            { return 2 }
func (LookingForAScrapYellow) Types() card.TypeSet     { return lookingForAScrapTypes }
func (LookingForAScrapYellow) GoAgain() bool           { return false }

// not implemented: graveyard-banish additional cost and the +1{p}/go-again bonus rider
func (LookingForAScrapYellow) NotImplemented() {}
func (c LookingForAScrapYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type LookingForAScrapBlue struct{}

func (LookingForAScrapBlue) ID() ids.CardID          { return ids.LookingForAScrapBlue }
func (LookingForAScrapBlue) Name() string            { return "Looking for a Scrap" }
func (LookingForAScrapBlue) Cost(*sim.TurnState) int { return 1 }
func (LookingForAScrapBlue) Pitch() int              { return 3 }
func (LookingForAScrapBlue) Attack() int             { return 2 }
func (LookingForAScrapBlue) Defense() int            { return 2 }
func (LookingForAScrapBlue) Types() card.TypeSet     { return lookingForAScrapTypes }
func (LookingForAScrapBlue) GoAgain() bool           { return false }

// not implemented: graveyard-banish additional cost and the +1{p}/go-again bonus rider
func (LookingForAScrapBlue) NotImplemented() {}
func (c LookingForAScrapBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
