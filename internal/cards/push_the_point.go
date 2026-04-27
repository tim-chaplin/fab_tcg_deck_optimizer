// Push the Point — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If the last attack on this combat chain hit, Push the Point gains +2{p}."

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var pushThePointTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type PushThePointRed struct{}

func (PushThePointRed) ID() card.ID              { return card.PushThePointRed }
func (PushThePointRed) Name() string             { return "Push the Point" }
func (PushThePointRed) Cost(*card.TurnState) int { return 1 }
func (PushThePointRed) Pitch() int               { return 1 }
func (PushThePointRed) Attack() int              { return 4 }
func (PushThePointRed) Defense() int             { return 2 }
func (PushThePointRed) Types() card.TypeSet      { return pushThePointTypes }
func (PushThePointRed) GoAgain() bool            { return false }

// not implemented: chain-history +2{p} rider (in-chain history not readable from Play)
func (PushThePointRed) NotImplemented() {}
func (c PushThePointRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type PushThePointYellow struct{}

func (PushThePointYellow) ID() card.ID              { return card.PushThePointYellow }
func (PushThePointYellow) Name() string             { return "Push the Point" }
func (PushThePointYellow) Cost(*card.TurnState) int { return 1 }
func (PushThePointYellow) Pitch() int               { return 2 }
func (PushThePointYellow) Attack() int              { return 3 }
func (PushThePointYellow) Defense() int             { return 2 }
func (PushThePointYellow) Types() card.TypeSet      { return pushThePointTypes }
func (PushThePointYellow) GoAgain() bool            { return false }

// not implemented: chain-history +2{p} rider (in-chain history not readable from Play)
func (PushThePointYellow) NotImplemented() {}
func (c PushThePointYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type PushThePointBlue struct{}

func (PushThePointBlue) ID() card.ID              { return card.PushThePointBlue }
func (PushThePointBlue) Name() string             { return "Push the Point" }
func (PushThePointBlue) Cost(*card.TurnState) int { return 1 }
func (PushThePointBlue) Pitch() int               { return 3 }
func (PushThePointBlue) Attack() int              { return 2 }
func (PushThePointBlue) Defense() int             { return 2 }
func (PushThePointBlue) Types() card.TypeSet      { return pushThePointTypes }
func (PushThePointBlue) GoAgain() bool            { return false }

// not implemented: chain-history +2{p} rider (in-chain history not readable from Play)
func (PushThePointBlue) NotImplemented() {}
func (c PushThePointBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
