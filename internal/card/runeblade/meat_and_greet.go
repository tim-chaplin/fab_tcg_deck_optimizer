// Meat and Greet — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When this hits, create a Runechant token. If you've dealt arcane damage to an opposing
// hero this turn, this gets go again."
//
// Simplifications:
//   - The "on hit, create a Runechant" is baked in as +1 damage (the Runechant's future value)
//     and AuraCreated is set for any later-in-chain effects that care.
//   - Assume any attack triggers some Runechant arcane damage when played, so the "dealt arcane
//     damage this turn" clause is always satisfied and Go again is treated as printed.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var meatAndGreetTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type MeatAndGreetRed struct{}

func (MeatAndGreetRed) ID() card.ID                 { return card.MeatAndGreetRed }
func (MeatAndGreetRed) Name() string                 { return "Meat and Greet (Red)" }
func (MeatAndGreetRed) Cost() int                    { return 1 }
func (MeatAndGreetRed) Pitch() int                   { return 1 }
func (MeatAndGreetRed) Attack() int                  { return 4 }
func (MeatAndGreetRed) Defense() int                 { return 3 }
func (MeatAndGreetRed) Types() card.TypeSet          { return meatAndGreetTypes }
func (MeatAndGreetRed) GoAgain() bool                { return true }
func (c MeatAndGreetRed) Play(s *card.TurnState) int { s.AuraCreated = true; return c.Attack() + 1 }

type MeatAndGreetYellow struct{}

func (MeatAndGreetYellow) ID() card.ID                 { return card.MeatAndGreetYellow }
func (MeatAndGreetYellow) Name() string                 { return "Meat and Greet (Yellow)" }
func (MeatAndGreetYellow) Cost() int                    { return 1 }
func (MeatAndGreetYellow) Pitch() int                   { return 2 }
func (MeatAndGreetYellow) Attack() int                  { return 3 }
func (MeatAndGreetYellow) Defense() int                 { return 3 }
func (MeatAndGreetYellow) Types() card.TypeSet          { return meatAndGreetTypes }
func (MeatAndGreetYellow) GoAgain() bool                { return true }
func (c MeatAndGreetYellow) Play(s *card.TurnState) int { s.AuraCreated = true; return c.Attack() + 1 }

type MeatAndGreetBlue struct{}

func (MeatAndGreetBlue) ID() card.ID                 { return card.MeatAndGreetBlue }
func (MeatAndGreetBlue) Name() string                 { return "Meat and Greet (Blue)" }
func (MeatAndGreetBlue) Cost() int                    { return 1 }
func (MeatAndGreetBlue) Pitch() int                   { return 3 }
func (MeatAndGreetBlue) Attack() int                  { return 2 }
func (MeatAndGreetBlue) Defense() int                 { return 3 }
func (MeatAndGreetBlue) Types() card.TypeSet          { return meatAndGreetTypes }
func (MeatAndGreetBlue) GoAgain() bool                { return true }
func (c MeatAndGreetBlue) Play(s *card.TurnState) int { s.AuraCreated = true; return c.Attack() + 1 }
