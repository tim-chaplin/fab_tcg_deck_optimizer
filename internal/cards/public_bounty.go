// Public Bounty — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 2.
//
// Text: "**Mark** target opposing hero. The next time you attack a **marked** hero this turn, the
// attack gets +N{p}. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var publicBountyTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type PublicBountyRed struct{}

func (PublicBountyRed) ID() ids.CardID          { return ids.PublicBountyRed }
func (PublicBountyRed) Name() string            { return "Public Bounty" }
func (PublicBountyRed) Cost(*sim.TurnState) int { return 1 }
func (PublicBountyRed) Pitch() int              { return 1 }
func (PublicBountyRed) Attack() int             { return 0 }
func (PublicBountyRed) Defense() int            { return 2 }
func (PublicBountyRed) Types() card.TypeSet     { return publicBountyTypes }
func (PublicBountyRed) GoAgain() bool           { return true }

// not implemented: mark not tracked; +3{p} 'marked defender' rider fires unconditionally
func (PublicBountyRed) NotImplemented() {}
func (PublicBountyRed) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextAttackActionBonus(s, 3)
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}

type PublicBountyYellow struct{}

func (PublicBountyYellow) ID() ids.CardID          { return ids.PublicBountyYellow }
func (PublicBountyYellow) Name() string            { return "Public Bounty" }
func (PublicBountyYellow) Cost(*sim.TurnState) int { return 1 }
func (PublicBountyYellow) Pitch() int              { return 2 }
func (PublicBountyYellow) Attack() int             { return 0 }
func (PublicBountyYellow) Defense() int            { return 2 }
func (PublicBountyYellow) Types() card.TypeSet     { return publicBountyTypes }
func (PublicBountyYellow) GoAgain() bool           { return true }

// not implemented: mark not tracked; +3{p} 'marked defender' rider fires unconditionally
func (PublicBountyYellow) NotImplemented() {}
func (PublicBountyYellow) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextAttackActionBonus(s, 2)
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}

type PublicBountyBlue struct{}

func (PublicBountyBlue) ID() ids.CardID          { return ids.PublicBountyBlue }
func (PublicBountyBlue) Name() string            { return "Public Bounty" }
func (PublicBountyBlue) Cost(*sim.TurnState) int { return 1 }
func (PublicBountyBlue) Pitch() int              { return 3 }
func (PublicBountyBlue) Attack() int             { return 0 }
func (PublicBountyBlue) Defense() int            { return 2 }
func (PublicBountyBlue) Types() card.TypeSet     { return publicBountyTypes }
func (PublicBountyBlue) GoAgain() bool           { return true }

// not implemented: mark not tracked; +3{p} 'marked defender' rider fires unconditionally
func (PublicBountyBlue) NotImplemented() {}
func (PublicBountyBlue) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextAttackActionBonus(s, 1)
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}
