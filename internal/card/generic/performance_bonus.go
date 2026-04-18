// Performance Bonus — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits, create a Gold token. If this was played from arsenal, it gets **Go
// again**."
//
// Simplification: Gold-token creation and arsenal-only go-again aren't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var performanceBonusTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type PerformanceBonusRed struct{}

func (PerformanceBonusRed) ID() card.ID                 { return card.PerformanceBonusRed }
func (PerformanceBonusRed) Name() string                { return "Performance Bonus (Red)" }
func (PerformanceBonusRed) Cost() int                   { return 0 }
func (PerformanceBonusRed) Pitch() int                  { return 1 }
func (PerformanceBonusRed) Attack() int                 { return 3 }
func (PerformanceBonusRed) Defense() int                { return 2 }
func (PerformanceBonusRed) Types() card.TypeSet         { return performanceBonusTypes }
func (PerformanceBonusRed) GoAgain() bool               { return true }
func (c PerformanceBonusRed) Play(s *card.TurnState) int { return c.Attack() }

type PerformanceBonusYellow struct{}

func (PerformanceBonusYellow) ID() card.ID                 { return card.PerformanceBonusYellow }
func (PerformanceBonusYellow) Name() string                { return "Performance Bonus (Yellow)" }
func (PerformanceBonusYellow) Cost() int                   { return 0 }
func (PerformanceBonusYellow) Pitch() int                  { return 2 }
func (PerformanceBonusYellow) Attack() int                 { return 2 }
func (PerformanceBonusYellow) Defense() int                { return 2 }
func (PerformanceBonusYellow) Types() card.TypeSet         { return performanceBonusTypes }
func (PerformanceBonusYellow) GoAgain() bool               { return true }
func (c PerformanceBonusYellow) Play(s *card.TurnState) int { return c.Attack() }

type PerformanceBonusBlue struct{}

func (PerformanceBonusBlue) ID() card.ID                 { return card.PerformanceBonusBlue }
func (PerformanceBonusBlue) Name() string                { return "Performance Bonus (Blue)" }
func (PerformanceBonusBlue) Cost() int                   { return 0 }
func (PerformanceBonusBlue) Pitch() int                  { return 3 }
func (PerformanceBonusBlue) Attack() int                 { return 1 }
func (PerformanceBonusBlue) Defense() int                { return 2 }
func (PerformanceBonusBlue) Types() card.TypeSet         { return performanceBonusTypes }
func (PerformanceBonusBlue) GoAgain() bool               { return true }
func (c PerformanceBonusBlue) Play(s *card.TurnState) int { return c.Attack() }
