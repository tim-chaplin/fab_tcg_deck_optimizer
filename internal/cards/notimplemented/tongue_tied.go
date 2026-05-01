// Tongue Tied — Generic Action - Attack. Cost 3, Pitch 1, Power 7, Defense 2. Only printed in Red.
//
// Text: "When this hits a hero, turn a card in their arsenal face-up, then banish an instant card
// from their arsenal."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var tongueTiedTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type TongueTiedRed struct{}

func (TongueTiedRed) ID() ids.CardID          { return ids.TongueTiedRed }
func (TongueTiedRed) Name() string            { return "Tongue Tied" }
func (TongueTiedRed) Cost(*sim.TurnState) int { return 3 }
func (TongueTiedRed) Pitch() int              { return 1 }
func (TongueTiedRed) Attack() int             { return 7 }
func (TongueTiedRed) Defense() int            { return 2 }
func (TongueTiedRed) Types() card.TypeSet     { return tongueTiedTypes }
func (TongueTiedRed) GoAgain() bool           { return false }

// not implemented: on-hit opponent-arsenal manipulation rider
func (TongueTiedRed) NotImplemented() {}
func (TongueTiedRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
