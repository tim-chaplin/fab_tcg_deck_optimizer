// Humble — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, they lose all hero card abilities until the end of their next
// turn."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var humbleTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type HumbleRed struct{}

func (HumbleRed) ID() ids.CardID          { return ids.HumbleRed }
func (HumbleRed) Name() string            { return "Humble" }
func (HumbleRed) Cost(*sim.TurnState) int { return 2 }
func (HumbleRed) Pitch() int              { return 1 }
func (HumbleRed) Attack() int             { return 6 }
func (HumbleRed) Defense() int            { return 2 }
func (HumbleRed) Types() card.TypeSet     { return humbleTypes }
func (HumbleRed) GoAgain() bool           { return false }

// not implemented: hero-ability suppression rider
func (HumbleRed) NotImplemented() {}
func (HumbleRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type HumbleYellow struct{}

func (HumbleYellow) ID() ids.CardID          { return ids.HumbleYellow }
func (HumbleYellow) Name() string            { return "Humble" }
func (HumbleYellow) Cost(*sim.TurnState) int { return 2 }
func (HumbleYellow) Pitch() int              { return 2 }
func (HumbleYellow) Attack() int             { return 5 }
func (HumbleYellow) Defense() int            { return 2 }
func (HumbleYellow) Types() card.TypeSet     { return humbleTypes }
func (HumbleYellow) GoAgain() bool           { return false }

// not implemented: hero-ability suppression rider
func (HumbleYellow) NotImplemented() {}
func (HumbleYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type HumbleBlue struct{}

func (HumbleBlue) ID() ids.CardID          { return ids.HumbleBlue }
func (HumbleBlue) Name() string            { return "Humble" }
func (HumbleBlue) Cost(*sim.TurnState) int { return 2 }
func (HumbleBlue) Pitch() int              { return 3 }
func (HumbleBlue) Attack() int             { return 4 }
func (HumbleBlue) Defense() int            { return 2 }
func (HumbleBlue) Types() card.TypeSet     { return humbleTypes }
func (HumbleBlue) GoAgain() bool           { return false }

// not implemented: hero-ability suppression rider
func (HumbleBlue) NotImplemented() {}
func (HumbleBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
