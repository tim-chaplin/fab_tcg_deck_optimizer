// Runic Fellingsong — Runeblade Action - Attack. Cost 3, Defense 3, printed 1 arcane.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 7, Yellow 6, Blue 5.
// Text: "When this attacks, you may banish an aura from your graveyard. If you do, deal 1 arcane
// damage to target hero."
//
// Simplification: credit a flat +1 arcane on top of printed power, covering the printed 1 arcane
// OR the banish rider. Not both — the rider requires a banishable aura in the graveyard, which we
// don't track.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var runicFellingsongTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type RunicFellingsongRed struct{}

func (RunicFellingsongRed) Name() string               { return "Runic Fellingsong (Red)" }
func (RunicFellingsongRed) Cost() int                  { return 3 }
func (RunicFellingsongRed) Pitch() int                 { return 1 }
func (RunicFellingsongRed) Attack() int                { return 7 }
func (RunicFellingsongRed) Defense() int               { return 3 }
func (RunicFellingsongRed) Types() card.TypeSet        { return runicFellingsongTypes }
func (RunicFellingsongRed) GoAgain() bool              { return false }
func (c RunicFellingsongRed) Play(*card.TurnState) int { return c.Attack() + 1 }

type RunicFellingsongYellow struct{}

func (RunicFellingsongYellow) Name() string               { return "Runic Fellingsong (Yellow)" }
func (RunicFellingsongYellow) Cost() int                  { return 3 }
func (RunicFellingsongYellow) Pitch() int                 { return 2 }
func (RunicFellingsongYellow) Attack() int                { return 6 }
func (RunicFellingsongYellow) Defense() int               { return 3 }
func (RunicFellingsongYellow) Types() card.TypeSet        { return runicFellingsongTypes }
func (RunicFellingsongYellow) GoAgain() bool              { return false }
func (c RunicFellingsongYellow) Play(*card.TurnState) int { return c.Attack() + 1 }

type RunicFellingsongBlue struct{}

func (RunicFellingsongBlue) Name() string               { return "Runic Fellingsong (Blue)" }
func (RunicFellingsongBlue) Cost() int                  { return 3 }
func (RunicFellingsongBlue) Pitch() int                 { return 3 }
func (RunicFellingsongBlue) Attack() int                { return 5 }
func (RunicFellingsongBlue) Defense() int               { return 3 }
func (RunicFellingsongBlue) Types() card.TypeSet        { return runicFellingsongTypes }
func (RunicFellingsongBlue) GoAgain() bool              { return false }
func (c RunicFellingsongBlue) Play(*card.TurnState) int { return c.Attack() + 1 }
