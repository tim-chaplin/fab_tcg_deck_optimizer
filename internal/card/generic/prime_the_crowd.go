// Prime the Crowd — Generic Action. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next attack action card you play this turn gets +4{p}. **The crowd cheers** each
// Revered hero. **The crowd boos** each Reviled hero. **Go again**"
//
// Simplification: Crowd cheers/boos keywords are dropped. Scans TurnState.CardsRemaining for the
// first matching attack action card and credits the bonus assuming it will be played; if none is
// scheduled after this card, the bonus fizzles.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var primeTheCrowdTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// primeTheCrowdPlay returns 4 when a matching attack action card is scheduled later this turn.
func primeTheCrowdPlay(s *card.TurnState) int {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if !t.Has(card.TypeAttack) || !t.Has(card.TypeAction) {
			continue
		}
		return 4
	}
	return 0
}

type PrimeTheCrowdRed struct{}

func (PrimeTheCrowdRed) ID() card.ID                 { return card.PrimeTheCrowdRed }
func (PrimeTheCrowdRed) Name() string                { return "Prime the Crowd (Red)" }
func (PrimeTheCrowdRed) Cost() int                   { return 2 }
func (PrimeTheCrowdRed) Pitch() int                  { return 1 }
func (PrimeTheCrowdRed) Attack() int                 { return 0 }
func (PrimeTheCrowdRed) Defense() int                { return 2 }
func (PrimeTheCrowdRed) Types() card.TypeSet         { return primeTheCrowdTypes }
func (PrimeTheCrowdRed) GoAgain() bool               { return true }
func (PrimeTheCrowdRed) Play(s *card.TurnState) int { return primeTheCrowdPlay(s) }

type PrimeTheCrowdYellow struct{}

func (PrimeTheCrowdYellow) ID() card.ID                 { return card.PrimeTheCrowdYellow }
func (PrimeTheCrowdYellow) Name() string                { return "Prime the Crowd (Yellow)" }
func (PrimeTheCrowdYellow) Cost() int                   { return 2 }
func (PrimeTheCrowdYellow) Pitch() int                  { return 2 }
func (PrimeTheCrowdYellow) Attack() int                 { return 0 }
func (PrimeTheCrowdYellow) Defense() int                { return 2 }
func (PrimeTheCrowdYellow) Types() card.TypeSet         { return primeTheCrowdTypes }
func (PrimeTheCrowdYellow) GoAgain() bool               { return true }
func (PrimeTheCrowdYellow) Play(s *card.TurnState) int { return primeTheCrowdPlay(s) }

type PrimeTheCrowdBlue struct{}

func (PrimeTheCrowdBlue) ID() card.ID                 { return card.PrimeTheCrowdBlue }
func (PrimeTheCrowdBlue) Name() string                { return "Prime the Crowd (Blue)" }
func (PrimeTheCrowdBlue) Cost() int                   { return 2 }
func (PrimeTheCrowdBlue) Pitch() int                  { return 3 }
func (PrimeTheCrowdBlue) Attack() int                 { return 0 }
func (PrimeTheCrowdBlue) Defense() int                { return 2 }
func (PrimeTheCrowdBlue) Types() card.TypeSet         { return primeTheCrowdTypes }
func (PrimeTheCrowdBlue) GoAgain() bool               { return true }
func (PrimeTheCrowdBlue) Play(s *card.TurnState) int { return primeTheCrowdPlay(s) }
