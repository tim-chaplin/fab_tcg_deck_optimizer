// Tip-Off — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4, Blue 3. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Instant** - Discard this: **Mark** target opposing hero."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var tipOffTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type TipOffRed struct{}

func (TipOffRed) ID() ids.CardID          { return ids.TipOffRed }
func (TipOffRed) Name() string            { return "Tip-Off" }
func (TipOffRed) Cost(*sim.TurnState) int { return 1 }
func (TipOffRed) Pitch() int              { return 1 }
func (TipOffRed) Attack() int             { return 5 }
func (TipOffRed) Defense() int            { return 2 }
func (TipOffRed) Types() card.TypeSet     { return tipOffTypes }
func (TipOffRed) GoAgain() bool           { return false }

// not implemented: instant discard-to-mark activation
func (TipOffRed) NotImplemented() {}
func (c TipOffRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type TipOffYellow struct{}

func (TipOffYellow) ID() ids.CardID          { return ids.TipOffYellow }
func (TipOffYellow) Name() string            { return "Tip-Off" }
func (TipOffYellow) Cost(*sim.TurnState) int { return 1 }
func (TipOffYellow) Pitch() int              { return 2 }
func (TipOffYellow) Attack() int             { return 4 }
func (TipOffYellow) Defense() int            { return 2 }
func (TipOffYellow) Types() card.TypeSet     { return tipOffTypes }
func (TipOffYellow) GoAgain() bool           { return false }

// not implemented: instant discard-to-mark activation
func (TipOffYellow) NotImplemented() {}
func (c TipOffYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type TipOffBlue struct{}

func (TipOffBlue) ID() ids.CardID          { return ids.TipOffBlue }
func (TipOffBlue) Name() string            { return "Tip-Off" }
func (TipOffBlue) Cost(*sim.TurnState) int { return 1 }
func (TipOffBlue) Pitch() int              { return 3 }
func (TipOffBlue) Attack() int             { return 3 }
func (TipOffBlue) Defense() int            { return 2 }
func (TipOffBlue) Types() card.TypeSet     { return tipOffTypes }
func (TipOffBlue) GoAgain() bool           { return false }

// not implemented: instant discard-to-mark activation
func (TipOffBlue) NotImplemented() {}
func (c TipOffBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
