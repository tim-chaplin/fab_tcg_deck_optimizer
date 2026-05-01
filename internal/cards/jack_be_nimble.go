// Jack Be Nimble — Generic Action - Attack. Cost 0, Pitch 1, Power 3, Defense 3. Only printed in
// Red.
//
// Text: "When this attacks, you may banish a Nimblism from your graveyard. If you do, this gets
// +1{p} and **go again**. When this hits a hero, steal an item they control until the end of this
// action phase."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var jackBeNimbleTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type JackBeNimbleRed struct{}

func (JackBeNimbleRed) ID() ids.CardID          { return ids.JackBeNimbleRed }
func (JackBeNimbleRed) Name() string            { return "Jack Be Nimble" }
func (JackBeNimbleRed) Cost(*sim.TurnState) int { return 0 }
func (JackBeNimbleRed) Pitch() int              { return 1 }
func (JackBeNimbleRed) Attack() int             { return 3 }
func (JackBeNimbleRed) Defense() int            { return 3 }
func (JackBeNimbleRed) Types() card.TypeSet     { return jackBeNimbleTypes }
func (JackBeNimbleRed) GoAgain() bool           { return false }

// not implemented: graveyard-banish cost + on-hit item steal
func (JackBeNimbleRed) NotImplemented() {}
func (JackBeNimbleRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}
