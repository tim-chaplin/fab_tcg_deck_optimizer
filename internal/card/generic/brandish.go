// Brandish — Generic Action - Attack. Cost 1. Printed power: Red 3, Yellow 2, Blue 1. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Brandish hits, your next weapon attack this turn gains +1{p}. **Go again**"

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var brandishTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BrandishRed struct{}

func (BrandishRed) ID() card.ID                 { return card.BrandishRed }
func (BrandishRed) Name() string                { return "Brandish (Red)" }
func (BrandishRed) Cost(*card.TurnState) int                   { return 1 }
func (BrandishRed) Pitch() int                  { return 1 }
func (BrandishRed) Attack() int                 { return 3 }
func (BrandishRed) Defense() int                { return 2 }
func (BrandishRed) Types() card.TypeSet         { return brandishTypes }
func (BrandishRed) GoAgain() bool               { return true }
// not implemented: next-weapon-attack +1{p} grant (weapon chain not scanned)
func (BrandishRed) NotImplemented()             {}
func (c BrandishRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type BrandishYellow struct{}

func (BrandishYellow) ID() card.ID                 { return card.BrandishYellow }
func (BrandishYellow) Name() string                { return "Brandish (Yellow)" }
func (BrandishYellow) Cost(*card.TurnState) int                   { return 1 }
func (BrandishYellow) Pitch() int                  { return 2 }
func (BrandishYellow) Attack() int                 { return 2 }
func (BrandishYellow) Defense() int                { return 2 }
func (BrandishYellow) Types() card.TypeSet         { return brandishTypes }
func (BrandishYellow) GoAgain() bool               { return true }
// not implemented: next-weapon-attack +1{p} grant (weapon chain not scanned)
func (BrandishYellow) NotImplemented()             {}
func (c BrandishYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type BrandishBlue struct{}

func (BrandishBlue) ID() card.ID                 { return card.BrandishBlue }
func (BrandishBlue) Name() string                { return "Brandish (Blue)" }
func (BrandishBlue) Cost(*card.TurnState) int                   { return 1 }
func (BrandishBlue) Pitch() int                  { return 3 }
func (BrandishBlue) Attack() int                 { return 1 }
func (BrandishBlue) Defense() int                { return 2 }
func (BrandishBlue) Types() card.TypeSet         { return brandishTypes }
func (BrandishBlue) GoAgain() bool               { return true }
// not implemented: next-weapon-attack +1{p} grant (weapon chain not scanned)
func (BrandishBlue) NotImplemented()             {}
func (c BrandishBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
