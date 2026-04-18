// Water the Seeds — Generic Action - Attack. Cost 1. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, your next attack this combat chain with 1 or less base {p} gets +1{p}.
// **Go again**"
//
// Simplification: Chain-bonus rider for next low-power attack isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var waterTheSeedsTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type WaterTheSeedsRed struct{}

func (WaterTheSeedsRed) ID() card.ID                 { return card.WaterTheSeedsRed }
func (WaterTheSeedsRed) Name() string                { return "Water the Seeds (Red)" }
func (WaterTheSeedsRed) Cost() int                   { return 1 }
func (WaterTheSeedsRed) Pitch() int                  { return 1 }
func (WaterTheSeedsRed) Attack() int                 { return 3 }
func (WaterTheSeedsRed) Defense() int                { return 2 }
func (WaterTheSeedsRed) Types() card.TypeSet         { return waterTheSeedsTypes }
func (WaterTheSeedsRed) GoAgain() bool               { return true }
func (c WaterTheSeedsRed) Play(s *card.TurnState) int { return c.Attack() }

type WaterTheSeedsYellow struct{}

func (WaterTheSeedsYellow) ID() card.ID                 { return card.WaterTheSeedsYellow }
func (WaterTheSeedsYellow) Name() string                { return "Water the Seeds (Yellow)" }
func (WaterTheSeedsYellow) Cost() int                   { return 1 }
func (WaterTheSeedsYellow) Pitch() int                  { return 2 }
func (WaterTheSeedsYellow) Attack() int                 { return 2 }
func (WaterTheSeedsYellow) Defense() int                { return 2 }
func (WaterTheSeedsYellow) Types() card.TypeSet         { return waterTheSeedsTypes }
func (WaterTheSeedsYellow) GoAgain() bool               { return true }
func (c WaterTheSeedsYellow) Play(s *card.TurnState) int { return c.Attack() }

type WaterTheSeedsBlue struct{}

func (WaterTheSeedsBlue) ID() card.ID                 { return card.WaterTheSeedsBlue }
func (WaterTheSeedsBlue) Name() string                { return "Water the Seeds (Blue)" }
func (WaterTheSeedsBlue) Cost() int                   { return 1 }
func (WaterTheSeedsBlue) Pitch() int                  { return 3 }
func (WaterTheSeedsBlue) Attack() int                 { return 1 }
func (WaterTheSeedsBlue) Defense() int                { return 2 }
func (WaterTheSeedsBlue) Types() card.TypeSet         { return waterTheSeedsTypes }
func (WaterTheSeedsBlue) GoAgain() bool               { return true }
func (c WaterTheSeedsBlue) Play(s *card.TurnState) int { return c.Attack() }
