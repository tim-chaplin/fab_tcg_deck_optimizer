// Smashing Good Time — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next time an attack action card hits a hero this turn, you may destroy an item they
// control with cost 2 or less. If Smashing Good Time is played from arsenal, the next attack action
// card you play this turn gains +N{p}. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: Item-destruction rider ignored; arsenal-only +N credited unconditionally. Scans
// TurnState.CardsRemaining for the first matching attack action card and credits the bonus assuming
// it will be played; if none is scheduled after this card, the bonus fizzles.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var smashingGoodTimeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type SmashingGoodTimeRed struct{}

func (SmashingGoodTimeRed) ID() card.ID                 { return card.SmashingGoodTimeRed }
func (SmashingGoodTimeRed) Name() string                { return "Smashing Good Time (Red)" }
func (SmashingGoodTimeRed) Cost(*card.TurnState) int                   { return 0 }
func (SmashingGoodTimeRed) Pitch() int                  { return 1 }
func (SmashingGoodTimeRed) Attack() int                 { return 0 }
func (SmashingGoodTimeRed) Defense() int                { return 2 }
func (SmashingGoodTimeRed) Types() card.TypeSet         { return smashingGoodTimeTypes }
func (SmashingGoodTimeRed) GoAgain() bool               { return true }
func (SmashingGoodTimeRed) Play(s *card.TurnState) int { return nextAttackActionBonus(s, 3) }

type SmashingGoodTimeYellow struct{}

func (SmashingGoodTimeYellow) ID() card.ID                 { return card.SmashingGoodTimeYellow }
func (SmashingGoodTimeYellow) Name() string                { return "Smashing Good Time (Yellow)" }
func (SmashingGoodTimeYellow) Cost(*card.TurnState) int                   { return 0 }
func (SmashingGoodTimeYellow) Pitch() int                  { return 2 }
func (SmashingGoodTimeYellow) Attack() int                 { return 0 }
func (SmashingGoodTimeYellow) Defense() int                { return 2 }
func (SmashingGoodTimeYellow) Types() card.TypeSet         { return smashingGoodTimeTypes }
func (SmashingGoodTimeYellow) GoAgain() bool               { return true }
func (SmashingGoodTimeYellow) Play(s *card.TurnState) int { return nextAttackActionBonus(s, 2) }

type SmashingGoodTimeBlue struct{}

func (SmashingGoodTimeBlue) ID() card.ID                 { return card.SmashingGoodTimeBlue }
func (SmashingGoodTimeBlue) Name() string                { return "Smashing Good Time (Blue)" }
func (SmashingGoodTimeBlue) Cost(*card.TurnState) int                   { return 0 }
func (SmashingGoodTimeBlue) Pitch() int                  { return 3 }
func (SmashingGoodTimeBlue) Attack() int                 { return 0 }
func (SmashingGoodTimeBlue) Defense() int                { return 2 }
func (SmashingGoodTimeBlue) Types() card.TypeSet         { return smashingGoodTimeTypes }
func (SmashingGoodTimeBlue) GoAgain() bool               { return true }
func (SmashingGoodTimeBlue) Play(s *card.TurnState) int { return nextAttackActionBonus(s, 1) }
