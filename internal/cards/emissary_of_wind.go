// Emissary of Wind — Generic Action - Attack. Cost 0, Pitch 1, Power 4, Defense 2. Only printed in
// Red.
//
// Text: "When this attacks, you may put a card from your hand on the bottom of your deck. If you
// do, this gets **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var emissaryOfWindTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type EmissaryOfWindRed struct{}

func (EmissaryOfWindRed) ID() ids.CardID          { return ids.EmissaryOfWindRed }
func (EmissaryOfWindRed) Name() string            { return "Emissary of Wind" }
func (EmissaryOfWindRed) Cost(*sim.TurnState) int { return 0 }
func (EmissaryOfWindRed) Pitch() int              { return 1 }
func (EmissaryOfWindRed) Attack() int             { return 4 }
func (EmissaryOfWindRed) Defense() int            { return 2 }
func (EmissaryOfWindRed) Types() card.TypeSet     { return emissaryOfWindTypes }
func (EmissaryOfWindRed) GoAgain() bool           { return false }

// not implemented: hand-cycle-for-go-again rider
func (EmissaryOfWindRed) NotImplemented() {}
func (c EmissaryOfWindRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
