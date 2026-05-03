// Stony Woottonhog — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "While Stony Woottonhog is defended by less than 2 non-equipment cards, it has +1{p}."
//
// Conservative model: the +1{p} bonus is dropped — assume the defender always blocks with
// 2+ non-equipment cards.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var stonyWoottonhogTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func stonyWoottonhogPlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type StonyWoottonhogRed struct{}

func (StonyWoottonhogRed) ID() ids.CardID          { return ids.StonyWoottonhogRed }
func (StonyWoottonhogRed) Name() string            { return "Stony Woottonhog" }
func (StonyWoottonhogRed) Cost(*sim.TurnState) int { return 2 }
func (StonyWoottonhogRed) Pitch() int              { return 1 }
func (StonyWoottonhogRed) Attack() int             { return 6 }
func (StonyWoottonhogRed) Defense() int            { return 2 }
func (StonyWoottonhogRed) Types() card.TypeSet     { return stonyWoottonhogTypes }
func (StonyWoottonhogRed) GoAgain() bool           { return false }
func (StonyWoottonhogRed) Play(s *sim.TurnState, self *sim.CardState) {
	stonyWoottonhogPlay(s, self)
}

type StonyWoottonhogYellow struct{}

func (StonyWoottonhogYellow) ID() ids.CardID          { return ids.StonyWoottonhogYellow }
func (StonyWoottonhogYellow) Name() string            { return "Stony Woottonhog" }
func (StonyWoottonhogYellow) Cost(*sim.TurnState) int { return 2 }
func (StonyWoottonhogYellow) Pitch() int              { return 2 }
func (StonyWoottonhogYellow) Attack() int             { return 5 }
func (StonyWoottonhogYellow) Defense() int            { return 2 }
func (StonyWoottonhogYellow) Types() card.TypeSet     { return stonyWoottonhogTypes }
func (StonyWoottonhogYellow) GoAgain() bool           { return false }
func (StonyWoottonhogYellow) Play(s *sim.TurnState, self *sim.CardState) {
	stonyWoottonhogPlay(s, self)
}

type StonyWoottonhogBlue struct{}

func (StonyWoottonhogBlue) ID() ids.CardID          { return ids.StonyWoottonhogBlue }
func (StonyWoottonhogBlue) Name() string            { return "Stony Woottonhog" }
func (StonyWoottonhogBlue) Cost(*sim.TurnState) int { return 2 }
func (StonyWoottonhogBlue) Pitch() int              { return 3 }
func (StonyWoottonhogBlue) Attack() int             { return 4 }
func (StonyWoottonhogBlue) Defense() int            { return 2 }
func (StonyWoottonhogBlue) Types() card.TypeSet     { return stonyWoottonhogTypes }
func (StonyWoottonhogBlue) GoAgain() bool           { return false }
func (StonyWoottonhogBlue) Play(s *sim.TurnState, self *sim.CardState) {
	stonyWoottonhogPlay(s, self)
}
