// Wounded Bull — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When you play this, if you have less {h} than an opposing hero, this gains +1{p}."
//
// Simplification: Health comparison isn't modelled; +1{p} rider never fires.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var woundedBullTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type WoundedBullRed struct{}

func (WoundedBullRed) ID() card.ID                 { return card.WoundedBullRed }
func (WoundedBullRed) Name() string                { return "Wounded Bull (Red)" }
func (WoundedBullRed) Cost(*card.TurnState) int                   { return 3 }
func (WoundedBullRed) Pitch() int                  { return 1 }
func (WoundedBullRed) Attack() int                 { return 7 }
func (WoundedBullRed) Defense() int                { return 2 }
func (WoundedBullRed) Types() card.TypeSet         { return woundedBullTypes }
func (WoundedBullRed) GoAgain() bool               { return false }
func (c WoundedBullRed) Play(s *card.TurnState) int { return c.Attack() }

type WoundedBullYellow struct{}

func (WoundedBullYellow) ID() card.ID                 { return card.WoundedBullYellow }
func (WoundedBullYellow) Name() string                { return "Wounded Bull (Yellow)" }
func (WoundedBullYellow) Cost(*card.TurnState) int                   { return 3 }
func (WoundedBullYellow) Pitch() int                  { return 2 }
func (WoundedBullYellow) Attack() int                 { return 6 }
func (WoundedBullYellow) Defense() int                { return 2 }
func (WoundedBullYellow) Types() card.TypeSet         { return woundedBullTypes }
func (WoundedBullYellow) GoAgain() bool               { return false }
func (c WoundedBullYellow) Play(s *card.TurnState) int { return c.Attack() }

type WoundedBullBlue struct{}

func (WoundedBullBlue) ID() card.ID                 { return card.WoundedBullBlue }
func (WoundedBullBlue) Name() string                { return "Wounded Bull (Blue)" }
func (WoundedBullBlue) Cost(*card.TurnState) int                   { return 3 }
func (WoundedBullBlue) Pitch() int                  { return 3 }
func (WoundedBullBlue) Attack() int                 { return 5 }
func (WoundedBullBlue) Defense() int                { return 2 }
func (WoundedBullBlue) Types() card.TypeSet         { return woundedBullTypes }
func (WoundedBullBlue) GoAgain() bool               { return false }
func (c WoundedBullBlue) Play(s *card.TurnState) int { return c.Attack() }
