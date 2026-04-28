// Emissary of Moon — Generic Action - Attack. Cost 0, Pitch 1, Power 4, Defense 2. Only printed in
// Red.
//
// Text: "When this attacks, you may put a card from your hand on the bottom of your deck. If you
// do, draw a card."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var emissaryOfMoonTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type EmissaryOfMoonRed struct{}

func (EmissaryOfMoonRed) ID() ids.CardID          { return ids.EmissaryOfMoonRed }
func (EmissaryOfMoonRed) Name() string            { return "Emissary of Moon" }
func (EmissaryOfMoonRed) Cost(*sim.TurnState) int { return 0 }
func (EmissaryOfMoonRed) Pitch() int              { return 1 }
func (EmissaryOfMoonRed) Attack() int             { return 4 }
func (EmissaryOfMoonRed) Defense() int            { return 2 }
func (EmissaryOfMoonRed) Types() card.TypeSet     { return emissaryOfMoonTypes }
func (EmissaryOfMoonRed) GoAgain() bool           { return false }

// not implemented: hand-cycle draw rider
func (EmissaryOfMoonRed) NotImplemented() {}
func (c EmissaryOfMoonRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
