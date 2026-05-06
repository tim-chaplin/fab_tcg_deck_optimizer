// Jack Be Nimble — Generic Action - Attack. Cost 0, Pitch 1, Power 3, Defense 3. Only printed in
// Red.
//
// Text: "When this attacks, you may banish a Nimblism from your graveyard. If you do, this gets
// +1{p} and **go again**. When this hits a hero, steal an item they control until the end of this
// action phase."
//
// The on-hit item-steal rider is not modelled: it's a sideboard-time consideration against
// item-heavy matchups, not a property of the card-vs-deck Value the simulator optimises.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var jackBeNimbleTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func jackBeNimblePlay(s *sim.TurnState, self *sim.CardState) {
	if _, ok := s.BanishFromGraveyard(isNimblism); ok {
		self.BonusAttack++
		self.GrantedGoAgain = true
		s.LogRider(self, 1, "Banished a Nimblism, +1{p} and go again")
	}
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type JackBeNimbleRed struct{}

func (JackBeNimbleRed) ID() ids.CardID          { return ids.JackBeNimbleRed }
func (JackBeNimbleRed) Name() string            { return "Jack Be Nimble" }
func (JackBeNimbleRed) Cost(*sim.TurnState) int { return 0 }
func (JackBeNimbleRed) Pitch() int              { return 1 }
func (JackBeNimbleRed) Attack() int             { return 3 }
func (JackBeNimbleRed) Defense() int            { return 3 }
func (JackBeNimbleRed) Types() card.TypeSet     { return jackBeNimbleTypes }
func (JackBeNimbleRed) GoAgain() bool           { return false }
func (JackBeNimbleRed) Play(s *sim.TurnState, self *sim.CardState) {
	jackBeNimblePlay(s, self)
}
