// Sloggism — Generic Action. Cost 3. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next attack action card with cost 2 or greater you play this turn gains +6{p}. **Go
// again**"
//
// Simplification: Scans TurnState.CardsRemaining for the first matching attack action card and
// credits the bonus assuming it will be played; if none is scheduled after this card, the bonus
// fizzles.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sloggismTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// sloggismPlay returns 6 when a matching attack action card is scheduled later this turn.
func sloggismPlay(s *card.TurnState) int {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if !t.Has(card.TypeAttack) || !t.Has(card.TypeAction) {
			continue
		}
		if pc.Card.Cost() >= 2 {
			return 6
		}
		continue
	}
	return 0
}

type SloggismRed struct{}

func (SloggismRed) ID() card.ID                 { return card.SloggismRed }
func (SloggismRed) Name() string                { return "Sloggism (Red)" }
func (SloggismRed) Cost() int                   { return 3 }
func (SloggismRed) Pitch() int                  { return 1 }
func (SloggismRed) Attack() int                 { return 0 }
func (SloggismRed) Defense() int                { return 2 }
func (SloggismRed) Types() card.TypeSet         { return sloggismTypes }
func (SloggismRed) GoAgain() bool               { return true }
func (SloggismRed) Play(s *card.TurnState) int { return sloggismPlay(s) }

type SloggismYellow struct{}

func (SloggismYellow) ID() card.ID                 { return card.SloggismYellow }
func (SloggismYellow) Name() string                { return "Sloggism (Yellow)" }
func (SloggismYellow) Cost() int                   { return 3 }
func (SloggismYellow) Pitch() int                  { return 2 }
func (SloggismYellow) Attack() int                 { return 0 }
func (SloggismYellow) Defense() int                { return 2 }
func (SloggismYellow) Types() card.TypeSet         { return sloggismTypes }
func (SloggismYellow) GoAgain() bool               { return true }
func (SloggismYellow) Play(s *card.TurnState) int { return sloggismPlay(s) }

type SloggismBlue struct{}

func (SloggismBlue) ID() card.ID                 { return card.SloggismBlue }
func (SloggismBlue) Name() string                { return "Sloggism (Blue)" }
func (SloggismBlue) Cost() int                   { return 3 }
func (SloggismBlue) Pitch() int                  { return 3 }
func (SloggismBlue) Attack() int                 { return 0 }
func (SloggismBlue) Defense() int                { return 2 }
func (SloggismBlue) Types() card.TypeSet         { return sloggismTypes }
func (SloggismBlue) GoAgain() bool               { return true }
func (SloggismBlue) Play(s *card.TurnState) int { return sloggismPlay(s) }
