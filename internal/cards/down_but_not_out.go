// Down But Not Out — Generic Action - Attack. Cost 3. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 3.
//
// Text: "When this attacks a hero, if you have less {h} and control fewer equipment and tokens than
// them, this gets +3{p}, **overpower**, and "When this hits, create an Agility, Might, and Vigor
// token.""

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var downButNotOutTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type DownButNotOutRed struct{}

func (DownButNotOutRed) ID() ids.CardID          { return ids.DownButNotOutRed }
func (DownButNotOutRed) Name() string            { return "Down But Not Out" }
func (DownButNotOutRed) Cost(*sim.TurnState) int { return 3 }
func (DownButNotOutRed) Pitch() int              { return 1 }
func (DownButNotOutRed) Attack() int             { return 5 }
func (DownButNotOutRed) Defense() int            { return 3 }
func (DownButNotOutRed) Types() card.TypeSet     { return downButNotOutTypes }
func (DownButNotOutRed) GoAgain() bool           { return false }

// not implemented: health/equipment/token comparison, agility/might/vigor tokens, overpower
func (DownButNotOutRed) NotImplemented() {}
func (DownButNotOutRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type DownButNotOutYellow struct{}

func (DownButNotOutYellow) ID() ids.CardID          { return ids.DownButNotOutYellow }
func (DownButNotOutYellow) Name() string            { return "Down But Not Out" }
func (DownButNotOutYellow) Cost(*sim.TurnState) int { return 3 }
func (DownButNotOutYellow) Pitch() int              { return 2 }
func (DownButNotOutYellow) Attack() int             { return 4 }
func (DownButNotOutYellow) Defense() int            { return 3 }
func (DownButNotOutYellow) Types() card.TypeSet     { return downButNotOutTypes }
func (DownButNotOutYellow) GoAgain() bool           { return false }

// not implemented: health/equipment/token comparison, agility/might/vigor tokens, overpower
func (DownButNotOutYellow) NotImplemented() {}
func (DownButNotOutYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type DownButNotOutBlue struct{}

func (DownButNotOutBlue) ID() ids.CardID          { return ids.DownButNotOutBlue }
func (DownButNotOutBlue) Name() string            { return "Down But Not Out" }
func (DownButNotOutBlue) Cost(*sim.TurnState) int { return 3 }
func (DownButNotOutBlue) Pitch() int              { return 3 }
func (DownButNotOutBlue) Attack() int             { return 3 }
func (DownButNotOutBlue) Defense() int            { return 3 }
func (DownButNotOutBlue) Types() card.TypeSet     { return downButNotOutTypes }
func (DownButNotOutBlue) GoAgain() bool           { return false }

// not implemented: health/equipment/token comparison, agility/might/vigor tokens, overpower
func (DownButNotOutBlue) NotImplemented() {}
func (DownButNotOutBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
