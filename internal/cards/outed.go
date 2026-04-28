// Outed — Generic Action - Attack. Cost 0, Pitch 1, Power 3, Defense 0. Only printed in Red.
//
// Text: "If you are **marked**, you can't play this. If the defending hero is **marked**, this gets
// +1{p}. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var outedTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type OutedRed struct{}

func (OutedRed) ID() ids.CardID          { return ids.OutedRed }
func (OutedRed) Name() string            { return "Outed" }
func (OutedRed) Cost(*sim.TurnState) int { return 0 }
func (OutedRed) Pitch() int              { return 1 }
func (OutedRed) Attack() int             { return 3 }
func (OutedRed) Defense() int            { return 0 }
func (OutedRed) Types() card.TypeSet     { return outedTypes }
func (OutedRed) GoAgain() bool           { return true }

// not implemented: marked-hero state not tracked; +1{p}-vs-marked-defender rider never fires
func (OutedRed) NotImplemented()                              {}
func (c OutedRed) Play(s *sim.TurnState, self *sim.CardState) { s.ApplyAndLogEffectiveAttack(self) }
