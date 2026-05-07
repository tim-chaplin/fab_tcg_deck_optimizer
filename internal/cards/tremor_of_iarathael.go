// Tremor of íArathael — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If a card has been put into your banished zone this turn, Tremor of íArathael gains
// +2{p}."
//
// Snapshot at Play: only banishes earlier in the chain trigger the +2{p}; later banishes don't
// retroactively buff this attack (matches the past-tense "has been put").

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var tremorOfIArathaelTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func tremorOfIArathaelPlay(s *sim.TurnState, self *sim.CardState) {
	if len(s.Banish) > 0 {
		self.BonusAttack += 2
	}
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type TremorOfIArathaelRed struct{}

func (TremorOfIArathaelRed) ID() ids.CardID          { return ids.TremorOfIArathaelRed }
func (TremorOfIArathaelRed) Name() string            { return "Tremor of íArathael" }
func (TremorOfIArathaelRed) Cost(*sim.TurnState) int { return 1 }
func (TremorOfIArathaelRed) Pitch() int              { return 1 }
func (TremorOfIArathaelRed) Attack() int             { return 4 }
func (TremorOfIArathaelRed) Defense() int            { return 2 }
func (TremorOfIArathaelRed) Types() card.TypeSet     { return tremorOfIArathaelTypes }
func (TremorOfIArathaelRed) GoAgain() bool           { return false }
func (TremorOfIArathaelRed) Play(s *sim.TurnState, self *sim.CardState) {
	tremorOfIArathaelPlay(s, self)
}

type TremorOfIArathaelYellow struct{}

func (TremorOfIArathaelYellow) ID() ids.CardID          { return ids.TremorOfIArathaelYellow }
func (TremorOfIArathaelYellow) Name() string            { return "Tremor of íArathael" }
func (TremorOfIArathaelYellow) Cost(*sim.TurnState) int { return 1 }
func (TremorOfIArathaelYellow) Pitch() int              { return 2 }
func (TremorOfIArathaelYellow) Attack() int             { return 3 }
func (TremorOfIArathaelYellow) Defense() int            { return 2 }
func (TremorOfIArathaelYellow) Types() card.TypeSet     { return tremorOfIArathaelTypes }
func (TremorOfIArathaelYellow) GoAgain() bool           { return false }
func (TremorOfIArathaelYellow) Play(s *sim.TurnState, self *sim.CardState) {
	tremorOfIArathaelPlay(s, self)
}

type TremorOfIArathaelBlue struct{}

func (TremorOfIArathaelBlue) ID() ids.CardID          { return ids.TremorOfIArathaelBlue }
func (TremorOfIArathaelBlue) Name() string            { return "Tremor of íArathael" }
func (TremorOfIArathaelBlue) Cost(*sim.TurnState) int { return 1 }
func (TremorOfIArathaelBlue) Pitch() int              { return 3 }
func (TremorOfIArathaelBlue) Attack() int             { return 2 }
func (TremorOfIArathaelBlue) Defense() int            { return 2 }
func (TremorOfIArathaelBlue) Types() card.TypeSet     { return tremorOfIArathaelTypes }
func (TremorOfIArathaelBlue) GoAgain() bool           { return false }
func (TremorOfIArathaelBlue) Play(s *sim.TurnState, self *sim.CardState) {
	tremorOfIArathaelPlay(s, self)
}
