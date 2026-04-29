// Trot Along — Generic Action. Cost 0, Pitch 3, Defense 3. Only printed in Blue.
//
// Text: "Your next attack with 3 or less base {p} this turn gets **go again**. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var trotAlongTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// trotAlongApplySideEffect grants go again to the next qualifying attack scheduled later
// this turn — attack action card OR weapon swing per the "your next attack" wording —
// gated on base power 3 or less.
func trotAlongApplySideEffect(s *sim.TurnState) {
	for _, pc := range s.CardsRemaining {
		if !pc.Card.Types().IsAttack() {
			continue
		}
		if pc.Card.Attack() <= 3 {
			pc.GrantedGoAgain = true
			return
		}
	}
}

type TrotAlongBlue struct{}

func (TrotAlongBlue) ID() ids.CardID          { return ids.TrotAlongBlue }
func (TrotAlongBlue) Name() string            { return "Trot Along" }
func (TrotAlongBlue) Cost(*sim.TurnState) int { return 0 }
func (TrotAlongBlue) Pitch() int              { return 3 }
func (TrotAlongBlue) Attack() int             { return 0 }
func (TrotAlongBlue) Defense() int            { return 3 }
func (TrotAlongBlue) Types() card.TypeSet     { return trotAlongTypes }
func (TrotAlongBlue) GoAgain() bool           { return true }
func (TrotAlongBlue) Play(s *sim.TurnState, self *sim.CardState) {
	trotAlongApplySideEffect(s)
	s.ApplyAndLogEffectiveAttack(self)
}
