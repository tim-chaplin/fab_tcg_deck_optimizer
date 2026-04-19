// Nimblism — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next attack action card with cost 1 or less you play this turn gains +N{p}. **Go
// again**" (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: Scans TurnState.CardsRemaining for the first matching attack action card and
// credits the bonus assuming it will be played; if none is scheduled after this card, the bonus
// fizzles.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var nimblismTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// nimblismPlay returns n when a matching attack action card is scheduled later this turn.
func nimblismPlay(s *card.TurnState, n int) int {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if !t.Has(card.TypeAttack) || !t.Has(card.TypeAction) {
			continue
		}
		if pc.Card.Cost(s) <= 1 {
			return n
		}
		continue
	}
	return 0
}

type NimblismRed struct{}

func (NimblismRed) ID() card.ID                 { return card.NimblismRed }
func (NimblismRed) Name() string                { return "Nimblism (Red)" }
func (NimblismRed) Cost(*card.TurnState) int                   { return 0 }
func (NimblismRed) Pitch() int                  { return 1 }
func (NimblismRed) Attack() int                 { return 0 }
func (NimblismRed) Defense() int                { return 2 }
func (NimblismRed) Types() card.TypeSet         { return nimblismTypes }
func (NimblismRed) GoAgain() bool               { return true }
func (NimblismRed) Play(s *card.TurnState) int { return nimblismPlay(s, 3) }

type NimblismYellow struct{}

func (NimblismYellow) ID() card.ID                 { return card.NimblismYellow }
func (NimblismYellow) Name() string                { return "Nimblism (Yellow)" }
func (NimblismYellow) Cost(*card.TurnState) int                   { return 0 }
func (NimblismYellow) Pitch() int                  { return 2 }
func (NimblismYellow) Attack() int                 { return 0 }
func (NimblismYellow) Defense() int                { return 2 }
func (NimblismYellow) Types() card.TypeSet         { return nimblismTypes }
func (NimblismYellow) GoAgain() bool               { return true }
func (NimblismYellow) Play(s *card.TurnState) int { return nimblismPlay(s, 2) }

type NimblismBlue struct{}

func (NimblismBlue) ID() card.ID                 { return card.NimblismBlue }
func (NimblismBlue) Name() string                { return "Nimblism (Blue)" }
func (NimblismBlue) Cost(*card.TurnState) int                   { return 0 }
func (NimblismBlue) Pitch() int                  { return 3 }
func (NimblismBlue) Attack() int                 { return 0 }
func (NimblismBlue) Defense() int                { return 2 }
func (NimblismBlue) Types() card.TypeSet         { return nimblismTypes }
func (NimblismBlue) GoAgain() bool               { return true }
func (NimblismBlue) Play(s *card.TurnState) int { return nimblismPlay(s, 1) }
