// Force Sight — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next attack action card you play this turn gains +N{p}. If Force Sight is played from
// arsenal, **opt 2**. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)
//
// The Opt 2 fires only when this copy was played from arsenal (self.FromArsenal).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var forceSightTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// forceSightPlay grants the next attack action +bonus{p}, logs the chain step (Force
// Sight is a non-attack action — no Attack() to apply), and resolves the arsenal-gated
// Opt 2.
func forceSightPlay(s *sim.TurnState, self *sim.CardState, bonus int) {
	grantNextAttackActionBonus(s, bonus)
	s.Log(self, 0)
	if self.FromArsenal {
		s.Opt(2)
	}
}

type ForceSightRed struct{}

func (ForceSightRed) ID() ids.CardID          { return ids.ForceSightRed }
func (ForceSightRed) Name() string            { return "Force Sight" }
func (ForceSightRed) Cost(*sim.TurnState) int { return 1 }
func (ForceSightRed) Pitch() int              { return 1 }
func (ForceSightRed) Attack() int             { return 0 }
func (ForceSightRed) Defense() int            { return 2 }
func (ForceSightRed) Types() card.TypeSet     { return forceSightTypes }
func (ForceSightRed) GoAgain() bool           { return true }
func (ForceSightRed) Play(s *sim.TurnState, self *sim.CardState) {
	forceSightPlay(s, self, 3)
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
func (ForceSightYellow) Play(s *sim.TurnState, self *sim.CardState) {
	forceSightPlay(s, self, 2)
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
func (ForceSightBlue) Play(s *sim.TurnState, self *sim.CardState) {
	forceSightPlay(s, self, 1)
}
