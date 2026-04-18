// Scout the Periphery — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Look at the top card of target hero's deck. The next attack action card you play from
// arsenal this turn gains +N{p}. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: The 'play from arsenal' gate is ignored (arsenal not modelled). Scans
// TurnState.CardsRemaining for the first matching attack action card and credits the bonus assuming
// it will be played; if none is scheduled after this card, the bonus fizzles.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var scoutThePeripheryTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// scoutThePeripheryPlay returns n when a matching attack action card is scheduled later this turn.
func scoutThePeripheryPlay(s *card.TurnState, n int) int {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if !t.Has(card.TypeAttack) || !t.Has(card.TypeAction) {
			continue
		}
		return n
	}
	return 0
}

type ScoutThePeripheryRed struct{}

func (ScoutThePeripheryRed) ID() card.ID                 { return card.ScoutThePeripheryRed }
func (ScoutThePeripheryRed) Name() string                { return "Scout the Periphery (Red)" }
func (ScoutThePeripheryRed) Cost() int                   { return 0 }
func (ScoutThePeripheryRed) Pitch() int                  { return 1 }
func (ScoutThePeripheryRed) Attack() int                 { return 0 }
func (ScoutThePeripheryRed) Defense() int                { return 2 }
func (ScoutThePeripheryRed) Types() card.TypeSet         { return scoutThePeripheryTypes }
func (ScoutThePeripheryRed) GoAgain() bool               { return true }
func (ScoutThePeripheryRed) Play(s *card.TurnState) int { return scoutThePeripheryPlay(s, 3) }

type ScoutThePeripheryYellow struct{}

func (ScoutThePeripheryYellow) ID() card.ID                 { return card.ScoutThePeripheryYellow }
func (ScoutThePeripheryYellow) Name() string                { return "Scout the Periphery (Yellow)" }
func (ScoutThePeripheryYellow) Cost() int                   { return 0 }
func (ScoutThePeripheryYellow) Pitch() int                  { return 2 }
func (ScoutThePeripheryYellow) Attack() int                 { return 0 }
func (ScoutThePeripheryYellow) Defense() int                { return 2 }
func (ScoutThePeripheryYellow) Types() card.TypeSet         { return scoutThePeripheryTypes }
func (ScoutThePeripheryYellow) GoAgain() bool               { return true }
func (ScoutThePeripheryYellow) Play(s *card.TurnState) int { return scoutThePeripheryPlay(s, 2) }

type ScoutThePeripheryBlue struct{}

func (ScoutThePeripheryBlue) ID() card.ID                 { return card.ScoutThePeripheryBlue }
func (ScoutThePeripheryBlue) Name() string                { return "Scout the Periphery (Blue)" }
func (ScoutThePeripheryBlue) Cost() int                   { return 0 }
func (ScoutThePeripheryBlue) Pitch() int                  { return 3 }
func (ScoutThePeripheryBlue) Attack() int                 { return 0 }
func (ScoutThePeripheryBlue) Defense() int                { return 2 }
func (ScoutThePeripheryBlue) Types() card.TypeSet         { return scoutThePeripheryTypes }
func (ScoutThePeripheryBlue) GoAgain() bool               { return true }
func (ScoutThePeripheryBlue) Play(s *card.TurnState) int { return scoutThePeripheryPlay(s, 1) }
