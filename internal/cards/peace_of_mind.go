// Peace of Mind — Generic Instant. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text (Red): "The next time you would be dealt {p} damage, prevent 4 of that damage.
// Create a Ponder token." Yellow caps prevention at 3, Blue at 2. The Ponder token rider
// is dropped — token economies aren't modelled.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var peaceOfMindTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

func peaceOfMindPlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}

type PeaceOfMindRed struct{}

func (PeaceOfMindRed) ID() ids.CardID          { return ids.PeaceOfMindRed }
func (PeaceOfMindRed) Name() string            { return "Peace of Mind" }
func (PeaceOfMindRed) Cost(*sim.TurnState) int { return 2 }
func (PeaceOfMindRed) Pitch() int              { return 1 }
func (PeaceOfMindRed) Attack() int             { return 0 }
func (PeaceOfMindRed) Defense() int            { return 4 }
func (PeaceOfMindRed) Types() card.TypeSet     { return peaceOfMindTypes }
func (PeaceOfMindRed) GoAgain() bool           { return false }
func (PeaceOfMindRed) DefensiveInstant()       {}
func (PeaceOfMindRed) Play(s *sim.TurnState, self *sim.CardState) {
	peaceOfMindPlay(s, self)
}

type PeaceOfMindYellow struct{}

func (PeaceOfMindYellow) ID() ids.CardID          { return ids.PeaceOfMindYellow }
func (PeaceOfMindYellow) Name() string            { return "Peace of Mind" }
func (PeaceOfMindYellow) Cost(*sim.TurnState) int { return 2 }
func (PeaceOfMindYellow) Pitch() int              { return 2 }
func (PeaceOfMindYellow) Attack() int             { return 0 }
func (PeaceOfMindYellow) Defense() int            { return 3 }
func (PeaceOfMindYellow) Types() card.TypeSet     { return peaceOfMindTypes }
func (PeaceOfMindYellow) GoAgain() bool           { return false }
func (PeaceOfMindYellow) DefensiveInstant()       {}
func (PeaceOfMindYellow) Play(s *sim.TurnState, self *sim.CardState) {
	peaceOfMindPlay(s, self)
}

type PeaceOfMindBlue struct{}

func (PeaceOfMindBlue) ID() ids.CardID          { return ids.PeaceOfMindBlue }
func (PeaceOfMindBlue) Name() string            { return "Peace of Mind" }
func (PeaceOfMindBlue) Cost(*sim.TurnState) int { return 2 }
func (PeaceOfMindBlue) Pitch() int              { return 3 }
func (PeaceOfMindBlue) Attack() int             { return 0 }
func (PeaceOfMindBlue) Defense() int            { return 2 }
func (PeaceOfMindBlue) Types() card.TypeSet     { return peaceOfMindTypes }
func (PeaceOfMindBlue) GoAgain() bool           { return false }
func (PeaceOfMindBlue) DefensiveInstant()       {}
func (PeaceOfMindBlue) Play(s *sim.TurnState, self *sim.CardState) {
	peaceOfMindPlay(s, self)
}
