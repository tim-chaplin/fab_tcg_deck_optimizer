// Captain's Call — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 2.
//
// Text: "Choose 1; The next attack action card with cost N or less you play this turn gains +2{p}.
// The next attack action card with cost N or less you play this turn gains **go again**. **Go
// again**" (Red N=2, Yellow N=1, Blue N=0.)
//
// Simplification: Modal: we pick the +2 power mode; the alternative 'go again' mode is dropped.
// Scans TurnState.CardsRemaining for the first matching attack action card and credits the bonus
// assuming it will be played; if none is scheduled after this card, the bonus fizzles.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var captainsCallTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// captainsCallPlay returns 2 when a matching attack action card with cost <= maxCost is scheduled
// later this turn.
func captainsCallPlay(s *card.TurnState, maxCost int) int {
	for _, pc := range s.CardsRemaining {
		if !pc.Card.Types().IsAttackAction() {
			continue
		}
		if pc.Card.Cost(s) <= maxCost {
			return 2
		}
	}
	return 0
}

type CaptainsCallRed struct{}

func (CaptainsCallRed) ID() card.ID                 { return card.CaptainsCallRed }
func (CaptainsCallRed) Name() string                { return "Captain's Call (Red)" }
func (CaptainsCallRed) Cost(*card.TurnState) int                   { return 0 }
func (CaptainsCallRed) Pitch() int                  { return 1 }
func (CaptainsCallRed) Attack() int                 { return 0 }
func (CaptainsCallRed) Defense() int                { return 2 }
func (CaptainsCallRed) Types() card.TypeSet         { return captainsCallTypes }
func (CaptainsCallRed) GoAgain() bool               { return true }
func (CaptainsCallRed) Play(s *card.TurnState, _ *card.CardState) int { return captainsCallPlay(s, 2) }

type CaptainsCallYellow struct{}

func (CaptainsCallYellow) ID() card.ID                 { return card.CaptainsCallYellow }
func (CaptainsCallYellow) Name() string                { return "Captain's Call (Yellow)" }
func (CaptainsCallYellow) Cost(*card.TurnState) int                   { return 0 }
func (CaptainsCallYellow) Pitch() int                  { return 2 }
func (CaptainsCallYellow) Attack() int                 { return 0 }
func (CaptainsCallYellow) Defense() int                { return 2 }
func (CaptainsCallYellow) Types() card.TypeSet         { return captainsCallTypes }
func (CaptainsCallYellow) GoAgain() bool               { return true }
func (CaptainsCallYellow) Play(s *card.TurnState, _ *card.CardState) int { return captainsCallPlay(s, 1) }

type CaptainsCallBlue struct{}

func (CaptainsCallBlue) ID() card.ID                 { return card.CaptainsCallBlue }
func (CaptainsCallBlue) Name() string                { return "Captain's Call (Blue)" }
func (CaptainsCallBlue) Cost(*card.TurnState) int                   { return 0 }
func (CaptainsCallBlue) Pitch() int                  { return 3 }
func (CaptainsCallBlue) Attack() int                 { return 0 }
func (CaptainsCallBlue) Defense() int                { return 2 }
func (CaptainsCallBlue) Types() card.TypeSet         { return captainsCallTypes }
func (CaptainsCallBlue) GoAgain() bool               { return true }
func (CaptainsCallBlue) Play(s *card.TurnState, _ *card.CardState) int { return captainsCallPlay(s, 0) }
