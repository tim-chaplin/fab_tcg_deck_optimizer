// Minnowism — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next attack action card with 3 or less base {p} you play this turn gains +N{p}. **Go
// again**" (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: Scans TurnState.CardsRemaining for the first matching attack action card and
// credits the bonus assuming it will be played; if none is scheduled after this card, the bonus
// fizzles.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var minnowismTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// minnowismPlay returns n when a matching attack action card is scheduled later this turn.
func minnowismPlay(s *card.TurnState, n int) int {
	for _, pc := range s.CardsRemaining {
		if !pc.Card.Types().IsAttackAction() {
			continue
		}
		if pc.Card.Attack() <= 3 {
			return n
		}
	}
	return 0
}

type MinnowismRed struct{}

func (MinnowismRed) ID() card.ID                 { return card.MinnowismRed }
func (MinnowismRed) Name() string                { return "Minnowism (Red)" }
func (MinnowismRed) Cost(*card.TurnState) int                   { return 0 }
func (MinnowismRed) Pitch() int                  { return 1 }
func (MinnowismRed) Attack() int                 { return 0 }
func (MinnowismRed) Defense() int                { return 2 }
func (MinnowismRed) Types() card.TypeSet         { return minnowismTypes }
func (MinnowismRed) GoAgain() bool               { return true }
func (MinnowismRed) Play(s *card.TurnState, _ *card.CardState) int { return minnowismPlay(s, 3) }

type MinnowismYellow struct{}

func (MinnowismYellow) ID() card.ID                 { return card.MinnowismYellow }
func (MinnowismYellow) Name() string                { return "Minnowism (Yellow)" }
func (MinnowismYellow) Cost(*card.TurnState) int                   { return 0 }
func (MinnowismYellow) Pitch() int                  { return 2 }
func (MinnowismYellow) Attack() int                 { return 0 }
func (MinnowismYellow) Defense() int                { return 2 }
func (MinnowismYellow) Types() card.TypeSet         { return minnowismTypes }
func (MinnowismYellow) GoAgain() bool               { return true }
func (MinnowismYellow) Play(s *card.TurnState, _ *card.CardState) int { return minnowismPlay(s, 2) }

type MinnowismBlue struct{}

func (MinnowismBlue) ID() card.ID                 { return card.MinnowismBlue }
func (MinnowismBlue) Name() string                { return "Minnowism (Blue)" }
func (MinnowismBlue) Cost(*card.TurnState) int                   { return 0 }
func (MinnowismBlue) Pitch() int                  { return 3 }
func (MinnowismBlue) Attack() int                 { return 0 }
func (MinnowismBlue) Defense() int                { return 2 }
func (MinnowismBlue) Types() card.TypeSet         { return minnowismTypes }
func (MinnowismBlue) GoAgain() bool               { return true }
func (MinnowismBlue) Play(s *card.TurnState, _ *card.CardState) int { return minnowismPlay(s, 1) }
