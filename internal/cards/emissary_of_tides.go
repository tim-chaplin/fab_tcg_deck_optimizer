// Emissary of Tides — Generic Action - Attack. Cost 0, Pitch 1, Power 4, Defense 2. Only printed in
// Red.
//
// Text: "When this attacks, you may put a card from your hand on the bottom of your deck. If you
// do, this gets +2{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var emissaryOfTidesTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type EmissaryOfTidesRed struct{}

func (EmissaryOfTidesRed) ID() ids.CardID          { return ids.EmissaryOfTidesRed }
func (EmissaryOfTidesRed) Name() string            { return "Emissary of Tides" }
func (EmissaryOfTidesRed) Cost(*sim.TurnState) int { return 0 }
func (EmissaryOfTidesRed) Pitch() int              { return 1 }
func (EmissaryOfTidesRed) Attack() int             { return 4 }
func (EmissaryOfTidesRed) Defense() int            { return 2 }
func (EmissaryOfTidesRed) Types() card.TypeSet     { return emissaryOfTidesTypes }
func (EmissaryOfTidesRed) GoAgain() bool           { return false }

// not implemented: hand-cycle-for-+2{p} rider
func (EmissaryOfTidesRed) NotImplemented() {}
func (c EmissaryOfTidesRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
