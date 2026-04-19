// Amplify the Arknight — Runeblade Action - Attack. Printed cost 3, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 6, Yellow 5, Blue 4.
// Text: "Amplify the Arknight costs {r} less to play for each Runechant you control."
//
// Variable cost: Cost reads s.Runechants to return max(0, printed - Runechants) at play time.
// Implements card.VariableCost so the solver can pre-screen with MinCost / MaxCost bounds.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var amplifyTheArknightTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

const amplifyTheArknightPrintedCost = 3

func amplifyTheArknightCost(s *card.TurnState) int {
	eff := amplifyTheArknightPrintedCost - s.Runechants
	if eff < 0 {
		return 0
	}
	return eff
}

type AmplifyTheArknightRed struct{}

func (AmplifyTheArknightRed) ID() card.ID                { return card.AmplifyTheArknightRed }
func (AmplifyTheArknightRed) Name() string               { return "Amplify the Arknight (Red)" }
func (AmplifyTheArknightRed) Cost(s *card.TurnState) int { return amplifyTheArknightCost(s) }
func (AmplifyTheArknightRed) MinCost() int               { return 0 }
func (AmplifyTheArknightRed) MaxCost() int               { return amplifyTheArknightPrintedCost }
func (AmplifyTheArknightRed) Pitch() int                 { return 1 }
func (AmplifyTheArknightRed) Attack() int                { return 6 }
func (AmplifyTheArknightRed) Defense() int               { return 3 }
func (AmplifyTheArknightRed) Types() card.TypeSet        { return amplifyTheArknightTypes }
func (AmplifyTheArknightRed) GoAgain() bool              { return false }
func (c AmplifyTheArknightRed) Play(*card.TurnState) int { return c.Attack() }

type AmplifyTheArknightYellow struct{}

func (AmplifyTheArknightYellow) ID() card.ID                { return card.AmplifyTheArknightYellow }
func (AmplifyTheArknightYellow) Name() string               { return "Amplify the Arknight (Yellow)" }
func (AmplifyTheArknightYellow) Cost(s *card.TurnState) int { return amplifyTheArknightCost(s) }
func (AmplifyTheArknightYellow) MinCost() int               { return 0 }
func (AmplifyTheArknightYellow) MaxCost() int               { return amplifyTheArknightPrintedCost }
func (AmplifyTheArknightYellow) Pitch() int                 { return 2 }
func (AmplifyTheArknightYellow) Attack() int                { return 5 }
func (AmplifyTheArknightYellow) Defense() int               { return 3 }
func (AmplifyTheArknightYellow) Types() card.TypeSet        { return amplifyTheArknightTypes }
func (AmplifyTheArknightYellow) GoAgain() bool              { return false }
func (c AmplifyTheArknightYellow) Play(*card.TurnState) int { return c.Attack() }

type AmplifyTheArknightBlue struct{}

func (AmplifyTheArknightBlue) ID() card.ID                { return card.AmplifyTheArknightBlue }
func (AmplifyTheArknightBlue) Name() string               { return "Amplify the Arknight (Blue)" }
func (AmplifyTheArknightBlue) Cost(s *card.TurnState) int { return amplifyTheArknightCost(s) }
func (AmplifyTheArknightBlue) MinCost() int               { return 0 }
func (AmplifyTheArknightBlue) MaxCost() int               { return amplifyTheArknightPrintedCost }
func (AmplifyTheArknightBlue) Pitch() int                 { return 3 }
func (AmplifyTheArknightBlue) Attack() int                { return 4 }
func (AmplifyTheArknightBlue) Defense() int               { return 3 }
func (AmplifyTheArknightBlue) Types() card.TypeSet        { return amplifyTheArknightTypes }
func (AmplifyTheArknightBlue) GoAgain() bool              { return false }
func (c AmplifyTheArknightBlue) Play(*card.TurnState) int { return c.Attack() }
