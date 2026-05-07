// Jack Be Quick — Generic Action - Attack. Cost 0, Pitch 1, Power 3, Defense 3. Only printed in
// Red.
//
// Text: "When this attacks, you may banish a Nimblism from your graveyard. If you do, this gets
// +1{p} and **go again**. When this hits a hero, {u} an ally they control, then steal it until the
// end of this action phase."
//
// The on-hit ally-steal rider is not modelled: it's a sideboard-time consideration against
// ally-heavy matchups, not a property of the card-vs-deck Value the simulator optimises.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var jackBeQuickTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func jackBeQuickPlay(s *sim.TurnState, self *sim.CardState) {
	if _, ok := s.BanishFromGraveyard(isNimblism); ok {
		self.BonusAttack++
		self.GrantedGoAgain = true
		s.LogRider(self, 1, "Banished a Nimblism, +1{p} and go again")
	}
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type JackBeQuickRed struct{}

func (JackBeQuickRed) ID() ids.CardID          { return ids.JackBeQuickRed }
func (JackBeQuickRed) Name() string            { return "Jack Be Quick" }
func (JackBeQuickRed) Cost(*sim.TurnState) int { return 0 }
func (JackBeQuickRed) Pitch() int              { return 1 }
func (JackBeQuickRed) Attack() int             { return 3 }
func (JackBeQuickRed) Defense() int            { return 3 }
func (JackBeQuickRed) Types() card.TypeSet     { return jackBeQuickTypes }
func (JackBeQuickRed) GoAgain() bool           { return false }
func (JackBeQuickRed) Play(s *sim.TurnState, self *sim.CardState) {
	jackBeQuickPlay(s, self)
}
