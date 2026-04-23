// Trot Along — Generic Action. Cost 0, Pitch 3, Defense 3. Only printed in Blue.
//
// Text: "Your next attack with 3 or less base {p} this turn gets **go again**. **Go again**"
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var trotAlongTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// trotAlongPlay grants go again to the next qualifying attack action card scheduled later this turn.
func trotAlongPlay(s *card.TurnState) int {
	for _, pc := range s.CardsRemaining {
		if !pc.Card.Types().IsAttackAction() {
			continue
		}
		if pc.Card.Attack() <= 3 {
			pc.GrantedGoAgain = true
			return 0
		}
	}
	return 0
}

type TrotAlongBlue struct{}

func (TrotAlongBlue) ID() card.ID                 { return card.TrotAlongBlue }
func (TrotAlongBlue) Name() string                { return "Trot Along (Blue)" }
func (TrotAlongBlue) Cost(*card.TurnState) int                   { return 0 }
func (TrotAlongBlue) Pitch() int                  { return 3 }
func (TrotAlongBlue) Attack() int                 { return 0 }
func (TrotAlongBlue) Defense() int                { return 3 }
func (TrotAlongBlue) Types() card.TypeSet         { return trotAlongTypes }
func (TrotAlongBlue) GoAgain() bool               { return true }
func (TrotAlongBlue) Play(s *card.TurnState, _ *card.CardState) int { return trotAlongPlay(s) }
