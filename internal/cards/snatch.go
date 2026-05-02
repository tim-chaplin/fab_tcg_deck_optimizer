// Snatch — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits, draw a card."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var snatchTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func snatchPlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	self.OnHit = append(self.OnHit, func(state *sim.TurnState) {
		state.DrawOne()
	})
}

type SnatchRed struct{}

func (SnatchRed) ID() ids.CardID          { return ids.SnatchRed }
func (SnatchRed) Name() string            { return "Snatch" }
func (SnatchRed) Cost(*sim.TurnState) int { return 0 }
func (SnatchRed) Pitch() int              { return 1 }
func (SnatchRed) Attack() int             { return 4 }
func (SnatchRed) Defense() int            { return 2 }
func (SnatchRed) Types() card.TypeSet     { return snatchTypes }
func (SnatchRed) GoAgain() bool           { return false }
func (SnatchRed) Play(s *sim.TurnState, self *sim.CardState) {
	snatchPlay(s, self)
}

type SnatchYellow struct{}

func (SnatchYellow) ID() ids.CardID          { return ids.SnatchYellow }
func (SnatchYellow) Name() string            { return "Snatch" }
func (SnatchYellow) Cost(*sim.TurnState) int { return 0 }
func (SnatchYellow) Pitch() int              { return 2 }
func (SnatchYellow) Attack() int             { return 3 }
func (SnatchYellow) Defense() int            { return 2 }
func (SnatchYellow) Types() card.TypeSet     { return snatchTypes }
func (SnatchYellow) GoAgain() bool           { return false }
func (SnatchYellow) Play(s *sim.TurnState, self *sim.CardState) {
	snatchPlay(s, self)
}

type SnatchBlue struct{}

func (SnatchBlue) ID() ids.CardID          { return ids.SnatchBlue }
func (SnatchBlue) Name() string            { return "Snatch" }
func (SnatchBlue) Cost(*sim.TurnState) int { return 0 }
func (SnatchBlue) Pitch() int              { return 3 }
func (SnatchBlue) Attack() int             { return 2 }
func (SnatchBlue) Defense() int            { return 2 }
func (SnatchBlue) Types() card.TypeSet     { return snatchTypes }
func (SnatchBlue) GoAgain() bool           { return false }
func (SnatchBlue) Play(s *sim.TurnState, self *sim.CardState) {
	snatchPlay(s, self)
}
