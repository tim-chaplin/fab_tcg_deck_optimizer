// Feisty Locals — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If this is defended by an action card, this gets +2{p}."
//
// Simplification: The 'defended by action card' +2{p} rider isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var feistyLocalsTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type FeistyLocalsRed struct{}

func (FeistyLocalsRed) ID() card.ID                 { return card.FeistyLocalsRed }
func (FeistyLocalsRed) Name() string                { return "Feisty Locals (Red)" }
func (FeistyLocalsRed) Cost(*card.TurnState) int                   { return 0 }
func (FeistyLocalsRed) Pitch() int                  { return 1 }
func (FeistyLocalsRed) Attack() int                 { return 3 }
func (FeistyLocalsRed) Defense() int                { return 2 }
func (FeistyLocalsRed) Types() card.TypeSet         { return feistyLocalsTypes }
func (FeistyLocalsRed) GoAgain() bool               { return false }
func (c FeistyLocalsRed) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }

type FeistyLocalsYellow struct{}

func (FeistyLocalsYellow) ID() card.ID                 { return card.FeistyLocalsYellow }
func (FeistyLocalsYellow) Name() string                { return "Feisty Locals (Yellow)" }
func (FeistyLocalsYellow) Cost(*card.TurnState) int                   { return 0 }
func (FeistyLocalsYellow) Pitch() int                  { return 2 }
func (FeistyLocalsYellow) Attack() int                 { return 2 }
func (FeistyLocalsYellow) Defense() int                { return 2 }
func (FeistyLocalsYellow) Types() card.TypeSet         { return feistyLocalsTypes }
func (FeistyLocalsYellow) GoAgain() bool               { return false }
func (c FeistyLocalsYellow) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }

type FeistyLocalsBlue struct{}

func (FeistyLocalsBlue) ID() card.ID                 { return card.FeistyLocalsBlue }
func (FeistyLocalsBlue) Name() string                { return "Feisty Locals (Blue)" }
func (FeistyLocalsBlue) Cost(*card.TurnState) int                   { return 0 }
func (FeistyLocalsBlue) Pitch() int                  { return 3 }
func (FeistyLocalsBlue) Attack() int                 { return 1 }
func (FeistyLocalsBlue) Defense() int                { return 2 }
func (FeistyLocalsBlue) Types() card.TypeSet         { return feistyLocalsTypes }
func (FeistyLocalsBlue) GoAgain() bool               { return false }
func (c FeistyLocalsBlue) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }
