// Overload — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Dominate** If Overload hits, it gains **go again**."
//
// Simplification: On-hit go-again and Dominate aren't modelled (keyword held but solver ignores
// hit-gated grants).
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var overloadTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type OverloadRed struct{}

func (OverloadRed) ID() card.ID                 { return card.OverloadRed }
func (OverloadRed) Name() string                { return "Overload (Red)" }
func (OverloadRed) Cost() int                   { return 0 }
func (OverloadRed) Pitch() int                  { return 1 }
func (OverloadRed) Attack() int                 { return 3 }
func (OverloadRed) Defense() int                { return 2 }
func (OverloadRed) Types() card.TypeSet         { return overloadTypes }
func (OverloadRed) GoAgain() bool               { return false }
func (c OverloadRed) Play(s *card.TurnState) int { return c.Attack() }

type OverloadYellow struct{}

func (OverloadYellow) ID() card.ID                 { return card.OverloadYellow }
func (OverloadYellow) Name() string                { return "Overload (Yellow)" }
func (OverloadYellow) Cost() int                   { return 0 }
func (OverloadYellow) Pitch() int                  { return 2 }
func (OverloadYellow) Attack() int                 { return 2 }
func (OverloadYellow) Defense() int                { return 2 }
func (OverloadYellow) Types() card.TypeSet         { return overloadTypes }
func (OverloadYellow) GoAgain() bool               { return false }
func (c OverloadYellow) Play(s *card.TurnState) int { return c.Attack() }

type OverloadBlue struct{}

func (OverloadBlue) ID() card.ID                 { return card.OverloadBlue }
func (OverloadBlue) Name() string                { return "Overload (Blue)" }
func (OverloadBlue) Cost() int                   { return 0 }
func (OverloadBlue) Pitch() int                  { return 3 }
func (OverloadBlue) Attack() int                 { return 1 }
func (OverloadBlue) Defense() int                { return 2 }
func (OverloadBlue) Types() card.TypeSet         { return overloadTypes }
func (OverloadBlue) GoAgain() bool               { return false }
func (c OverloadBlue) Play(s *card.TurnState) int { return c.Attack() }
