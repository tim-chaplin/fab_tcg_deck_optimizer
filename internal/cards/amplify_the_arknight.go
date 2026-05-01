// Amplify the Arknight — Runeblade Action - Attack. Printed cost 3, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 6, Yellow 5, Blue 4.
// Text: "Amplify the Arknight costs {r} less to play for each Runechant you control."
//
// Variable cost: Cost reads s.Runechants to return max(0, printed - Runechants).
// Standard sim.VariableCost wiring (docs/dev-standards.md).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var amplifyTheArknightTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

const amplifyTheArknightPrintedCost = 3

func amplifyTheArknightCost(s *sim.TurnState) int {
	eff := amplifyTheArknightPrintedCost - s.Runechants
	if eff < 0 {
		return 0
	}
	return eff
}

type AmplifyTheArknightRed struct{}

func (AmplifyTheArknightRed) ID() ids.CardID            { return ids.AmplifyTheArknightRed }
func (AmplifyTheArknightRed) Name() string              { return "Amplify the Arknight" }
func (AmplifyTheArknightRed) Cost(s *sim.TurnState) int { return amplifyTheArknightCost(s) }
func (AmplifyTheArknightRed) MinCost() int              { return 0 }
func (AmplifyTheArknightRed) MaxCost() int              { return amplifyTheArknightPrintedCost }
func (AmplifyTheArknightRed) Pitch() int                { return 1 }
func (AmplifyTheArknightRed) Attack() int               { return 6 }
func (AmplifyTheArknightRed) Defense() int              { return 3 }
func (AmplifyTheArknightRed) Types() card.TypeSet       { return amplifyTheArknightTypes }
func (AmplifyTheArknightRed) GoAgain() bool             { return false }
func (AmplifyTheArknightRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type AmplifyTheArknightYellow struct{}

func (AmplifyTheArknightYellow) ID() ids.CardID            { return ids.AmplifyTheArknightYellow }
func (AmplifyTheArknightYellow) Name() string              { return "Amplify the Arknight" }
func (AmplifyTheArknightYellow) Cost(s *sim.TurnState) int { return amplifyTheArknightCost(s) }
func (AmplifyTheArknightYellow) MinCost() int              { return 0 }
func (AmplifyTheArknightYellow) MaxCost() int              { return amplifyTheArknightPrintedCost }
func (AmplifyTheArknightYellow) Pitch() int                { return 2 }
func (AmplifyTheArknightYellow) Attack() int               { return 5 }
func (AmplifyTheArknightYellow) Defense() int              { return 3 }
func (AmplifyTheArknightYellow) Types() card.TypeSet       { return amplifyTheArknightTypes }
func (AmplifyTheArknightYellow) GoAgain() bool             { return false }
func (AmplifyTheArknightYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type AmplifyTheArknightBlue struct{}

func (AmplifyTheArknightBlue) ID() ids.CardID            { return ids.AmplifyTheArknightBlue }
func (AmplifyTheArknightBlue) Name() string              { return "Amplify the Arknight" }
func (AmplifyTheArknightBlue) Cost(s *sim.TurnState) int { return amplifyTheArknightCost(s) }
func (AmplifyTheArknightBlue) MinCost() int              { return 0 }
func (AmplifyTheArknightBlue) MaxCost() int              { return amplifyTheArknightPrintedCost }
func (AmplifyTheArknightBlue) Pitch() int                { return 3 }
func (AmplifyTheArknightBlue) Attack() int               { return 4 }
func (AmplifyTheArknightBlue) Defense() int              { return 3 }
func (AmplifyTheArknightBlue) Types() card.TypeSet       { return amplifyTheArknightTypes }
func (AmplifyTheArknightBlue) GoAgain() bool             { return false }
func (AmplifyTheArknightBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
