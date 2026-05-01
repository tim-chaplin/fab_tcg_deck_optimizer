// Relentless Pursuit — Generic Action. Cost 0, Pitch 3, Defense 3. Only printed in Blue.
//
// Text: "**Mark** target opposing hero. If you've attacked them this turn, put this on the bottom
// of its owner's deck. **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var relentlessPursuitTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type RelentlessPursuitBlue struct{}

func (RelentlessPursuitBlue) ID() ids.CardID          { return ids.RelentlessPursuitBlue }
func (RelentlessPursuitBlue) Name() string            { return "Relentless Pursuit" }
func (RelentlessPursuitBlue) Cost(*sim.TurnState) int { return 0 }
func (RelentlessPursuitBlue) Pitch() int              { return 3 }
func (RelentlessPursuitBlue) Attack() int             { return 0 }
func (RelentlessPursuitBlue) Defense() int            { return 3 }
func (RelentlessPursuitBlue) Types() card.TypeSet     { return relentlessPursuitTypes }
func (RelentlessPursuitBlue) GoAgain() bool           { return true }

// not implemented: marked-target gate + 'attacked them this turn' chain rider
func (RelentlessPursuitBlue) NotImplemented()                            {}
func (RelentlessPursuitBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
