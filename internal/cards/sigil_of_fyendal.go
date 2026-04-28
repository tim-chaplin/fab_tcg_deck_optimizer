// Sigil of Fyendal — Generic Action - Aura. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "**Go again** At the beginning of your action phase, destroy this. When this leaves the
// arena, gain 1{h}."
//
// Handler credits +1 next turn for the 1{h} gain (valued 1-to-1 with damage).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var sigilOfFyendalTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAura)

type SigilOfFyendalBlue struct{}

func (SigilOfFyendalBlue) ID() ids.CardID          { return ids.SigilOfFyendalBlue }
func (SigilOfFyendalBlue) Name() string            { return "Sigil of Fyendal" }
func (SigilOfFyendalBlue) Cost(*sim.TurnState) int { return 0 }
func (SigilOfFyendalBlue) Pitch() int              { return 3 }
func (SigilOfFyendalBlue) Attack() int             { return 0 }
func (SigilOfFyendalBlue) Defense() int            { return 2 }
func (SigilOfFyendalBlue) Types() card.TypeSet     { return sigilOfFyendalTypes }
func (SigilOfFyendalBlue) GoAgain() bool           { return true }
func (SigilOfFyendalBlue) AddsFutureValue()        {}
func (c SigilOfFyendalBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.RegisterStartOfTurn(c, 1, "Gained 1 health", func(*sim.TurnState) int { return 1 })
	s.LogPlay(self)
}
