// Walk the Plank — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a Pirate hero, {t} them or an ally they control."
//
// Simplification: Pirate-specific target-freezing rider isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var walkThePlankTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type WalkThePlankRed struct{}

func (WalkThePlankRed) ID() card.ID                 { return card.WalkThePlankRed }
func (WalkThePlankRed) Name() string                { return "Walk the Plank (Red)" }
func (WalkThePlankRed) Cost() int                   { return 3 }
func (WalkThePlankRed) Pitch() int                  { return 1 }
func (WalkThePlankRed) Attack() int                 { return 7 }
func (WalkThePlankRed) Defense() int                { return 2 }
func (WalkThePlankRed) Types() card.TypeSet         { return walkThePlankTypes }
func (WalkThePlankRed) GoAgain() bool               { return false }
func (c WalkThePlankRed) Play(s *card.TurnState) int { return c.Attack() }

type WalkThePlankYellow struct{}

func (WalkThePlankYellow) ID() card.ID                 { return card.WalkThePlankYellow }
func (WalkThePlankYellow) Name() string                { return "Walk the Plank (Yellow)" }
func (WalkThePlankYellow) Cost() int                   { return 3 }
func (WalkThePlankYellow) Pitch() int                  { return 2 }
func (WalkThePlankYellow) Attack() int                 { return 6 }
func (WalkThePlankYellow) Defense() int                { return 2 }
func (WalkThePlankYellow) Types() card.TypeSet         { return walkThePlankTypes }
func (WalkThePlankYellow) GoAgain() bool               { return false }
func (c WalkThePlankYellow) Play(s *card.TurnState) int { return c.Attack() }

type WalkThePlankBlue struct{}

func (WalkThePlankBlue) ID() card.ID                 { return card.WalkThePlankBlue }
func (WalkThePlankBlue) Name() string                { return "Walk the Plank (Blue)" }
func (WalkThePlankBlue) Cost() int                   { return 3 }
func (WalkThePlankBlue) Pitch() int                  { return 3 }
func (WalkThePlankBlue) Attack() int                 { return 5 }
func (WalkThePlankBlue) Defense() int                { return 2 }
func (WalkThePlankBlue) Types() card.TypeSet         { return walkThePlankTypes }
func (WalkThePlankBlue) GoAgain() bool               { return false }
func (c WalkThePlankBlue) Play(s *card.TurnState) int { return c.Attack() }
