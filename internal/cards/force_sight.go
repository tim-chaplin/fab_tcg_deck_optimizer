// Force Sight — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next attack action card you play this turn gains +N{p}. If Force Sight is played from
// arsenal, **opt 2**. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var forceSightTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type ForceSightRed struct{}

func (ForceSightRed) ID() ids.CardID          { return ids.ForceSightRed }
func (ForceSightRed) Name() string            { return "Force Sight" }
func (ForceSightRed) Cost(*sim.TurnState) int { return 1 }
func (ForceSightRed) Pitch() int              { return 1 }
func (ForceSightRed) Attack() int             { return 0 }
func (ForceSightRed) Defense() int            { return 2 }
func (ForceSightRed) Types() card.TypeSet     { return forceSightTypes }
func (ForceSightRed) GoAgain() bool           { return true }

// not implemented: arsenal-gated Opt 2
func (ForceSightRed) NotImplemented() {}
func (ForceSightRed) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextAttackActionBonus(s, 3)
	s.ApplyAndLogEffectiveAttack(self)
}

type ForceSightYellow struct{}

func (ForceSightYellow) ID() ids.CardID          { return ids.ForceSightYellow }
func (ForceSightYellow) Name() string            { return "Force Sight" }
func (ForceSightYellow) Cost(*sim.TurnState) int { return 1 }
func (ForceSightYellow) Pitch() int              { return 2 }
func (ForceSightYellow) Attack() int             { return 0 }
func (ForceSightYellow) Defense() int            { return 2 }
func (ForceSightYellow) Types() card.TypeSet     { return forceSightTypes }
func (ForceSightYellow) GoAgain() bool           { return true }

// not implemented: arsenal-gated Opt 2
func (ForceSightYellow) NotImplemented() {}
func (ForceSightYellow) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextAttackActionBonus(s, 2)
	s.ApplyAndLogEffectiveAttack(self)
}

type ForceSightBlue struct{}

func (ForceSightBlue) ID() ids.CardID          { return ids.ForceSightBlue }
func (ForceSightBlue) Name() string            { return "Force Sight" }
func (ForceSightBlue) Cost(*sim.TurnState) int { return 1 }
func (ForceSightBlue) Pitch() int              { return 3 }
func (ForceSightBlue) Attack() int             { return 0 }
func (ForceSightBlue) Defense() int            { return 2 }
func (ForceSightBlue) Types() card.TypeSet     { return forceSightTypes }
func (ForceSightBlue) GoAgain() bool           { return true }

// not implemented: arsenal-gated Opt 2
func (ForceSightBlue) NotImplemented() {}
func (ForceSightBlue) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextAttackActionBonus(s, 1)
	s.ApplyAndLogEffectiveAttack(self)
}
