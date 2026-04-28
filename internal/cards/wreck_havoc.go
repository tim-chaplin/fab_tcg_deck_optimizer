// Wreck Havoc — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "Defense reactions can't be played to this chain link. When this hits a hero, you may turn
// a card in their arsenal face up, then destroy a defense reaction in their arsenal."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var wreckHavocTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type WreckHavocRed struct{}

func (WreckHavocRed) ID() ids.CardID          { return ids.WreckHavocRed }
func (WreckHavocRed) Name() string            { return "Wreck Havoc" }
func (WreckHavocRed) Cost(*sim.TurnState) int { return 2 }
func (WreckHavocRed) Pitch() int              { return 1 }
func (WreckHavocRed) Attack() int             { return 6 }
func (WreckHavocRed) Defense() int            { return 2 }
func (WreckHavocRed) Types() card.TypeSet     { return wreckHavocTypes }
func (WreckHavocRed) GoAgain() bool           { return false }

// not implemented: defense-reaction lockout, on-hit arsenal banish
func (WreckHavocRed) NotImplemented() {}
func (WreckHavocRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type WreckHavocYellow struct{}

func (WreckHavocYellow) ID() ids.CardID          { return ids.WreckHavocYellow }
func (WreckHavocYellow) Name() string            { return "Wreck Havoc" }
func (WreckHavocYellow) Cost(*sim.TurnState) int { return 2 }
func (WreckHavocYellow) Pitch() int              { return 2 }
func (WreckHavocYellow) Attack() int             { return 5 }
func (WreckHavocYellow) Defense() int            { return 2 }
func (WreckHavocYellow) Types() card.TypeSet     { return wreckHavocTypes }
func (WreckHavocYellow) GoAgain() bool           { return false }

// not implemented: defense-reaction lockout, on-hit arsenal banish
func (WreckHavocYellow) NotImplemented() {}
func (WreckHavocYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type WreckHavocBlue struct{}

func (WreckHavocBlue) ID() ids.CardID          { return ids.WreckHavocBlue }
func (WreckHavocBlue) Name() string            { return "Wreck Havoc" }
func (WreckHavocBlue) Cost(*sim.TurnState) int { return 2 }
func (WreckHavocBlue) Pitch() int              { return 3 }
func (WreckHavocBlue) Attack() int             { return 4 }
func (WreckHavocBlue) Defense() int            { return 2 }
func (WreckHavocBlue) Types() card.TypeSet     { return wreckHavocTypes }
func (WreckHavocBlue) GoAgain() bool           { return false }

// not implemented: defense-reaction lockout, on-hit arsenal banish
func (WreckHavocBlue) NotImplemented() {}
func (WreckHavocBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
