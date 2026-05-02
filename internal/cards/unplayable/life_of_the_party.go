// Life of the Party — Generic Action - Attack. Cost 2. Printed power: Red 4, Yellow 3,
// Blue 2. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "You may discard or destroy a card you control named Crazy Brew rather than pay
// Life of the Party's {r} cost. If you do, choose all modes, otherwise choose 1 at random;
// - This gets 'When this hits, gain life 2{h}.'
// - This gets +2{p}.
// - This gets **go again**."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var lifeOfThePartyTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type LifeOfThePartyRed struct{}

func (LifeOfThePartyRed) ID() ids.CardID          { return ids.LifeOfThePartyRed }
func (LifeOfThePartyRed) Name() string            { return "Life of the Party" }
func (LifeOfThePartyRed) Cost(*sim.TurnState) int { return 2 }
func (LifeOfThePartyRed) Pitch() int              { return 1 }
func (LifeOfThePartyRed) Attack() int             { return 4 }
func (LifeOfThePartyRed) Defense() int            { return 2 }
func (LifeOfThePartyRed) Types() card.TypeSet     { return lifeOfThePartyTypes }
func (LifeOfThePartyRed) GoAgain() bool           { return false }
func (LifeOfThePartyRed) Unplayable()             {}
func (LifeOfThePartyRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type LifeOfThePartyYellow struct{}

func (LifeOfThePartyYellow) ID() ids.CardID          { return ids.LifeOfThePartyYellow }
func (LifeOfThePartyYellow) Name() string            { return "Life of the Party" }
func (LifeOfThePartyYellow) Cost(*sim.TurnState) int { return 2 }
func (LifeOfThePartyYellow) Pitch() int              { return 2 }
func (LifeOfThePartyYellow) Attack() int             { return 3 }
func (LifeOfThePartyYellow) Defense() int            { return 2 }
func (LifeOfThePartyYellow) Types() card.TypeSet     { return lifeOfThePartyTypes }
func (LifeOfThePartyYellow) GoAgain() bool           { return false }
func (LifeOfThePartyYellow) Unplayable()             {}
func (LifeOfThePartyYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type LifeOfThePartyBlue struct{}

func (LifeOfThePartyBlue) ID() ids.CardID          { return ids.LifeOfThePartyBlue }
func (LifeOfThePartyBlue) Name() string            { return "Life of the Party" }
func (LifeOfThePartyBlue) Cost(*sim.TurnState) int { return 2 }
func (LifeOfThePartyBlue) Pitch() int              { return 3 }
func (LifeOfThePartyBlue) Attack() int             { return 2 }
func (LifeOfThePartyBlue) Defense() int            { return 2 }
func (LifeOfThePartyBlue) Types() card.TypeSet     { return lifeOfThePartyTypes }
func (LifeOfThePartyBlue) GoAgain() bool           { return false }
func (LifeOfThePartyBlue) Unplayable()             {}
func (LifeOfThePartyBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
