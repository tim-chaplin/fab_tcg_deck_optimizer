// Tremor of íArathael — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If a card has been put into your banished zone this turn, Tremor of íArathael gains
// +2{p}."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var tremorOfIArathaelTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type TremorOfIArathaelRed struct{}

func (TremorOfIArathaelRed) ID() card.ID              { return card.TremorOfIArathaelRed }
func (TremorOfIArathaelRed) Name() string             { return "Tremor of íArathael" }
func (TremorOfIArathaelRed) Cost(*card.TurnState) int { return 1 }
func (TremorOfIArathaelRed) Pitch() int               { return 1 }
func (TremorOfIArathaelRed) Attack() int              { return 4 }
func (TremorOfIArathaelRed) Defense() int             { return 2 }
func (TremorOfIArathaelRed) Types() card.TypeSet      { return tremorOfIArathaelTypes }
func (TremorOfIArathaelRed) GoAgain() bool            { return false }

// not implemented: banished-zone +2{p} rider (banished-zone count not tracked)
func (TremorOfIArathaelRed) NotImplemented() {}
func (c TremorOfIArathaelRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type TremorOfIArathaelYellow struct{}

func (TremorOfIArathaelYellow) ID() card.ID              { return card.TremorOfIArathaelYellow }
func (TremorOfIArathaelYellow) Name() string             { return "Tremor of íArathael" }
func (TremorOfIArathaelYellow) Cost(*card.TurnState) int { return 1 }
func (TremorOfIArathaelYellow) Pitch() int               { return 2 }
func (TremorOfIArathaelYellow) Attack() int              { return 3 }
func (TremorOfIArathaelYellow) Defense() int             { return 2 }
func (TremorOfIArathaelYellow) Types() card.TypeSet      { return tremorOfIArathaelTypes }
func (TremorOfIArathaelYellow) GoAgain() bool            { return false }

// not implemented: banished-zone +2{p} rider (banished-zone count not tracked)
func (TremorOfIArathaelYellow) NotImplemented() {}
func (c TremorOfIArathaelYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type TremorOfIArathaelBlue struct{}

func (TremorOfIArathaelBlue) ID() card.ID              { return card.TremorOfIArathaelBlue }
func (TremorOfIArathaelBlue) Name() string             { return "Tremor of íArathael" }
func (TremorOfIArathaelBlue) Cost(*card.TurnState) int { return 1 }
func (TremorOfIArathaelBlue) Pitch() int               { return 3 }
func (TremorOfIArathaelBlue) Attack() int              { return 2 }
func (TremorOfIArathaelBlue) Defense() int             { return 2 }
func (TremorOfIArathaelBlue) Types() card.TypeSet      { return tremorOfIArathaelTypes }
func (TremorOfIArathaelBlue) GoAgain() bool            { return false }

// not implemented: banished-zone +2{p} rider (banished-zone count not tracked)
func (TremorOfIArathaelBlue) NotImplemented() {}
func (c TremorOfIArathaelBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
