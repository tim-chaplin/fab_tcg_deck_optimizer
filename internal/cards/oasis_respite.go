// Oasis Respite — Generic Instant. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text (Red): "Prevent the next 4 damage that would be dealt to target hero this turn by
// a source of your choice. If they have less life than each other hero, they may gain
// 1{h}." Yellow caps at 3, Blue at 2. The 1{h} life-gain rider fires for heroes opting
// into sim.LowerHealthWanter via sim.HeroWantsLowerHealth — life gain is credited to
// Value the same as damage prevention, so the rider lands on top of DealEffectiveDefense.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var oasisRespiteTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

func oasisRespitePlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	if sim.HeroWantsLowerHealth() {
		s.AddValue(1)
		n++
	}
	s.Log(self, n)
}

type OasisRespiteRed struct{}

func (OasisRespiteRed) ID() ids.CardID          { return ids.OasisRespiteRed }
func (OasisRespiteRed) Name() string            { return "Oasis Respite" }
func (OasisRespiteRed) Cost(*sim.TurnState) int { return 1 }
func (OasisRespiteRed) Pitch() int              { return 1 }
func (OasisRespiteRed) Attack() int             { return 0 }
func (OasisRespiteRed) Defense() int            { return 4 }
func (OasisRespiteRed) Types() card.TypeSet     { return oasisRespiteTypes }
func (OasisRespiteRed) GoAgain() bool           { return false }
func (OasisRespiteRed) DefensiveInstant()       {}
func (OasisRespiteRed) Play(s *sim.TurnState, self *sim.CardState) {
	oasisRespitePlay(s, self)
}

type OasisRespiteYellow struct{}

func (OasisRespiteYellow) ID() ids.CardID          { return ids.OasisRespiteYellow }
func (OasisRespiteYellow) Name() string            { return "Oasis Respite" }
func (OasisRespiteYellow) Cost(*sim.TurnState) int { return 1 }
func (OasisRespiteYellow) Pitch() int              { return 2 }
func (OasisRespiteYellow) Attack() int             { return 0 }
func (OasisRespiteYellow) Defense() int            { return 3 }
func (OasisRespiteYellow) Types() card.TypeSet     { return oasisRespiteTypes }
func (OasisRespiteYellow) GoAgain() bool           { return false }
func (OasisRespiteYellow) DefensiveInstant()       {}
func (OasisRespiteYellow) Play(s *sim.TurnState, self *sim.CardState) {
	oasisRespitePlay(s, self)
}

type OasisRespiteBlue struct{}

func (OasisRespiteBlue) ID() ids.CardID          { return ids.OasisRespiteBlue }
func (OasisRespiteBlue) Name() string            { return "Oasis Respite" }
func (OasisRespiteBlue) Cost(*sim.TurnState) int { return 1 }
func (OasisRespiteBlue) Pitch() int              { return 3 }
func (OasisRespiteBlue) Attack() int             { return 0 }
func (OasisRespiteBlue) Defense() int            { return 2 }
func (OasisRespiteBlue) Types() card.TypeSet     { return oasisRespiteTypes }
func (OasisRespiteBlue) GoAgain() bool           { return false }
func (OasisRespiteBlue) DefensiveInstant()       {}
func (OasisRespiteBlue) Play(s *sim.TurnState, self *sim.CardState) {
	oasisRespitePlay(s, self)
}
