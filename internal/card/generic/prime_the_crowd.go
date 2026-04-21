// Prime the Crowd — Generic Action. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next attack action card you play this turn gets +N{p}. **The crowd cheers** each
// Revered hero. **The crowd boos** each Reviled hero. **Go again**" (Red N=4, Yellow N=3, Blue
// N=2.)
//
// Simplification: Crowd cheers/boos keywords are dropped. Scans TurnState.CardsRemaining for the
// first matching attack action card and credits the bonus assuming it will be played; if none is
// scheduled after this card, the bonus fizzles.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var primeTheCrowdTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type PrimeTheCrowdRed struct{}

func (PrimeTheCrowdRed) ID() card.ID                 { return card.PrimeTheCrowdRed }
func (PrimeTheCrowdRed) Name() string                { return "Prime the Crowd (Red)" }
func (PrimeTheCrowdRed) Cost(*card.TurnState) int                   { return 2 }
func (PrimeTheCrowdRed) Pitch() int                  { return 1 }
func (PrimeTheCrowdRed) Attack() int                 { return 0 }
func (PrimeTheCrowdRed) Defense() int                { return 2 }
func (PrimeTheCrowdRed) Types() card.TypeSet         { return primeTheCrowdTypes }
func (PrimeTheCrowdRed) GoAgain() bool               { return true }
func (PrimeTheCrowdRed) Play(s *card.TurnState, _ *card.CardState) int { return nextAttackActionBonus(s, 4) }

type PrimeTheCrowdYellow struct{}

func (PrimeTheCrowdYellow) ID() card.ID                 { return card.PrimeTheCrowdYellow }
func (PrimeTheCrowdYellow) Name() string                { return "Prime the Crowd (Yellow)" }
func (PrimeTheCrowdYellow) Cost(*card.TurnState) int                   { return 2 }
func (PrimeTheCrowdYellow) Pitch() int                  { return 2 }
func (PrimeTheCrowdYellow) Attack() int                 { return 0 }
func (PrimeTheCrowdYellow) Defense() int                { return 2 }
func (PrimeTheCrowdYellow) Types() card.TypeSet         { return primeTheCrowdTypes }
func (PrimeTheCrowdYellow) GoAgain() bool               { return true }
func (PrimeTheCrowdYellow) Play(s *card.TurnState, _ *card.CardState) int { return nextAttackActionBonus(s, 3) }

type PrimeTheCrowdBlue struct{}

func (PrimeTheCrowdBlue) ID() card.ID                 { return card.PrimeTheCrowdBlue }
func (PrimeTheCrowdBlue) Name() string                { return "Prime the Crowd (Blue)" }
func (PrimeTheCrowdBlue) Cost(*card.TurnState) int                   { return 2 }
func (PrimeTheCrowdBlue) Pitch() int                  { return 3 }
func (PrimeTheCrowdBlue) Attack() int                 { return 0 }
func (PrimeTheCrowdBlue) Defense() int                { return 2 }
func (PrimeTheCrowdBlue) Types() card.TypeSet         { return primeTheCrowdTypes }
func (PrimeTheCrowdBlue) GoAgain() bool               { return true }
func (PrimeTheCrowdBlue) Play(s *card.TurnState, _ *card.CardState) int { return nextAttackActionBonus(s, 2) }
