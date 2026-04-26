// Stony Woottonhog — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "While Stony Woottonhog is defended by less than 2 non-equipment cards, it has +1{p}."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var stonyWoottonhogTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type StonyWoottonhogRed struct{}

func (StonyWoottonhogRed) ID() card.ID                 { return card.StonyWoottonhogRed }
func (StonyWoottonhogRed) Name() string                { return "Stony Woottonhog (Red)" }
func (StonyWoottonhogRed) Cost(*card.TurnState) int                   { return 2 }
func (StonyWoottonhogRed) Pitch() int                  { return 1 }
func (StonyWoottonhogRed) Attack() int                 { return 6 }
func (StonyWoottonhogRed) Defense() int                { return 2 }
func (StonyWoottonhogRed) Types() card.TypeSet         { return stonyWoottonhogTypes }
func (StonyWoottonhogRed) GoAgain() bool               { return false }
// not implemented: defended-by-<2-non-equipment +1{p} rider
func (StonyWoottonhogRed) NotImplemented()             {}
func (c StonyWoottonhogRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type StonyWoottonhogYellow struct{}

func (StonyWoottonhogYellow) ID() card.ID                 { return card.StonyWoottonhogYellow }
func (StonyWoottonhogYellow) Name() string                { return "Stony Woottonhog (Yellow)" }
func (StonyWoottonhogYellow) Cost(*card.TurnState) int                   { return 2 }
func (StonyWoottonhogYellow) Pitch() int                  { return 2 }
func (StonyWoottonhogYellow) Attack() int                 { return 5 }
func (StonyWoottonhogYellow) Defense() int                { return 2 }
func (StonyWoottonhogYellow) Types() card.TypeSet         { return stonyWoottonhogTypes }
func (StonyWoottonhogYellow) GoAgain() bool               { return false }
// not implemented: defended-by-<2-non-equipment +1{p} rider
func (StonyWoottonhogYellow) NotImplemented()             {}
func (c StonyWoottonhogYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type StonyWoottonhogBlue struct{}

func (StonyWoottonhogBlue) ID() card.ID                 { return card.StonyWoottonhogBlue }
func (StonyWoottonhogBlue) Name() string                { return "Stony Woottonhog (Blue)" }
func (StonyWoottonhogBlue) Cost(*card.TurnState) int                   { return 2 }
func (StonyWoottonhogBlue) Pitch() int                  { return 3 }
func (StonyWoottonhogBlue) Attack() int                 { return 4 }
func (StonyWoottonhogBlue) Defense() int                { return 2 }
func (StonyWoottonhogBlue) Types() card.TypeSet         { return stonyWoottonhogTypes }
func (StonyWoottonhogBlue) GoAgain() bool               { return false }
// not implemented: defended-by-<2-non-equipment +1{p} rider
func (StonyWoottonhogBlue) NotImplemented()             {}
func (c StonyWoottonhogBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
