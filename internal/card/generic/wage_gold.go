// Wage Gold — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Universal** When this attacks a hero, you may **wager** a Gold token with them."
//
// Simplification: Universal keyword and Gold-token wager aren't modelled.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var wageGoldTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type WageGoldRed struct{}

func (WageGoldRed) ID() card.ID                 { return card.WageGoldRed }
func (WageGoldRed) Name() string                { return "Wage Gold (Red)" }
func (WageGoldRed) Cost(*card.TurnState) int                   { return 3 }
func (WageGoldRed) Pitch() int                  { return 1 }
func (WageGoldRed) Attack() int                 { return 7 }
func (WageGoldRed) Defense() int                { return 2 }
func (WageGoldRed) Types() card.TypeSet         { return wageGoldTypes }
func (WageGoldRed) GoAgain() bool               { return false }
// not implemented: gold tokens, universal keyword
func (WageGoldRed) NotImplemented()              {}
func (c WageGoldRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type WageGoldYellow struct{}

func (WageGoldYellow) ID() card.ID                 { return card.WageGoldYellow }
func (WageGoldYellow) Name() string                { return "Wage Gold (Yellow)" }
func (WageGoldYellow) Cost(*card.TurnState) int                   { return 3 }
func (WageGoldYellow) Pitch() int                  { return 2 }
func (WageGoldYellow) Attack() int                 { return 6 }
func (WageGoldYellow) Defense() int                { return 2 }
func (WageGoldYellow) Types() card.TypeSet         { return wageGoldTypes }
func (WageGoldYellow) GoAgain() bool               { return false }
// not implemented: gold tokens, universal keyword
func (WageGoldYellow) NotImplemented()              {}
func (c WageGoldYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type WageGoldBlue struct{}

func (WageGoldBlue) ID() card.ID                 { return card.WageGoldBlue }
func (WageGoldBlue) Name() string                { return "Wage Gold (Blue)" }
func (WageGoldBlue) Cost(*card.TurnState) int                   { return 3 }
func (WageGoldBlue) Pitch() int                  { return 3 }
func (WageGoldBlue) Attack() int                 { return 5 }
func (WageGoldBlue) Defense() int                { return 2 }
func (WageGoldBlue) Types() card.TypeSet         { return wageGoldTypes }
func (WageGoldBlue) GoAgain() bool               { return false }
// not implemented: gold tokens, universal keyword
func (WageGoldBlue) NotImplemented()              {}
func (c WageGoldBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
