// Barraging Brawnhide — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "While Barraging Brawnhide is defended by less than 2 non-equipment cards, it has +1{p}."
//
// Simplification: Defended-by-<2-non-equipment condition isn't modelled; the +1{p} rider never
// applies.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var barragingBrawnhideTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BarragingBrawnhideRed struct{}

func (BarragingBrawnhideRed) ID() card.ID                 { return card.BarragingBrawnhideRed }
func (BarragingBrawnhideRed) Name() string                { return "Barraging Brawnhide (Red)" }
func (BarragingBrawnhideRed) Cost(*card.TurnState) int                   { return 3 }
func (BarragingBrawnhideRed) Pitch() int                  { return 1 }
func (BarragingBrawnhideRed) Attack() int                 { return 7 }
func (BarragingBrawnhideRed) Defense() int                { return 2 }
func (BarragingBrawnhideRed) Types() card.TypeSet         { return barragingBrawnhideTypes }
func (BarragingBrawnhideRed) GoAgain() bool               { return false }
func (c BarragingBrawnhideRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type BarragingBrawnhideYellow struct{}

func (BarragingBrawnhideYellow) ID() card.ID                 { return card.BarragingBrawnhideYellow }
func (BarragingBrawnhideYellow) Name() string                { return "Barraging Brawnhide (Yellow)" }
func (BarragingBrawnhideYellow) Cost(*card.TurnState) int                   { return 3 }
func (BarragingBrawnhideYellow) Pitch() int                  { return 2 }
func (BarragingBrawnhideYellow) Attack() int                 { return 6 }
func (BarragingBrawnhideYellow) Defense() int                { return 2 }
func (BarragingBrawnhideYellow) Types() card.TypeSet         { return barragingBrawnhideTypes }
func (BarragingBrawnhideYellow) GoAgain() bool               { return false }
func (c BarragingBrawnhideYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type BarragingBrawnhideBlue struct{}

func (BarragingBrawnhideBlue) ID() card.ID                 { return card.BarragingBrawnhideBlue }
func (BarragingBrawnhideBlue) Name() string                { return "Barraging Brawnhide (Blue)" }
func (BarragingBrawnhideBlue) Cost(*card.TurnState) int                   { return 3 }
func (BarragingBrawnhideBlue) Pitch() int                  { return 3 }
func (BarragingBrawnhideBlue) Attack() int                 { return 5 }
func (BarragingBrawnhideBlue) Defense() int                { return 2 }
func (BarragingBrawnhideBlue) Types() card.TypeSet         { return barragingBrawnhideTypes }
func (BarragingBrawnhideBlue) GoAgain() bool               { return false }
func (c BarragingBrawnhideBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
