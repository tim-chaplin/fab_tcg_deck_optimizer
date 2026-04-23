// Splintering Deadwood — Runeblade Action - Attack. Cost 3, Defense 3.
// Printed power: Red 7, Yellow 6, Blue 5.
// Text: "When this attacks or hits, you may destroy an aura you control. If you do, create a
// Runechant token."
//
// Simplification: the effect swaps an existing aura (worth a Runechant's value) for a new
// Runechant — net zero. Play returns base power.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var splinteringDeadwoodTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type SplinteringDeadwoodRed struct{}

func (SplinteringDeadwoodRed) ID() card.ID                 { return card.SplinteringDeadwoodRed }
func (SplinteringDeadwoodRed) Name() string               { return "Splintering Deadwood (Red)" }
func (SplinteringDeadwoodRed) Cost(*card.TurnState) int                  { return 3 }
func (SplinteringDeadwoodRed) Pitch() int                 { return 1 }
func (SplinteringDeadwoodRed) Attack() int                { return 7 }
func (SplinteringDeadwoodRed) Defense() int               { return 3 }
func (SplinteringDeadwoodRed) Types() card.TypeSet        { return splinteringDeadwoodTypes }
func (SplinteringDeadwoodRed) GoAgain() bool              { return false }
func (c SplinteringDeadwoodRed) Play(*card.TurnState, *card.CardState) int { return c.Attack() }

type SplinteringDeadwoodYellow struct{}

func (SplinteringDeadwoodYellow) ID() card.ID                 { return card.SplinteringDeadwoodYellow }
func (SplinteringDeadwoodYellow) Name() string               { return "Splintering Deadwood (Yellow)" }
func (SplinteringDeadwoodYellow) Cost(*card.TurnState) int                  { return 3 }
func (SplinteringDeadwoodYellow) Pitch() int                 { return 2 }
func (SplinteringDeadwoodYellow) Attack() int                { return 6 }
func (SplinteringDeadwoodYellow) Defense() int               { return 3 }
func (SplinteringDeadwoodYellow) Types() card.TypeSet        { return splinteringDeadwoodTypes }
func (SplinteringDeadwoodYellow) GoAgain() bool              { return false }
func (c SplinteringDeadwoodYellow) Play(*card.TurnState, *card.CardState) int { return c.Attack() }

type SplinteringDeadwoodBlue struct{}

func (SplinteringDeadwoodBlue) ID() card.ID                 { return card.SplinteringDeadwoodBlue }
func (SplinteringDeadwoodBlue) Name() string               { return "Splintering Deadwood (Blue)" }
func (SplinteringDeadwoodBlue) Cost(*card.TurnState) int                  { return 3 }
func (SplinteringDeadwoodBlue) Pitch() int                 { return 3 }
func (SplinteringDeadwoodBlue) Attack() int                { return 5 }
func (SplinteringDeadwoodBlue) Defense() int               { return 3 }
func (SplinteringDeadwoodBlue) Types() card.TypeSet        { return splinteringDeadwoodTypes }
func (SplinteringDeadwoodBlue) GoAgain() bool              { return false }
func (c SplinteringDeadwoodBlue) Play(*card.TurnState, *card.CardState) int { return c.Attack() }
