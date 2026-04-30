// Zealous Belting — Generic Action - Attack. Cost 2. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "While there is a card in your pitch zone with {p} greater than Zealous Belting's base {p},
// Zealous Belting has **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var zealousBeltingTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// zealousBeltingPlay grants go again when any pitched card this turn has base power greater
// than the card's own base power, then emits the chain step.
func zealousBeltingPlay(s *sim.TurnState, self *sim.CardState) {
	base := self.Card.Attack()
	for _, p := range s.Pitched {
		if p.Attack() > base {
			self.GrantedGoAgain = true
			break
		}
	}
	s.ApplyAndLogEffectiveAttack(self)
}

type ZealousBeltingRed struct{}

func (ZealousBeltingRed) ID() ids.CardID          { return ids.ZealousBeltingRed }
func (ZealousBeltingRed) Name() string            { return "Zealous Belting" }
func (ZealousBeltingRed) Cost(*sim.TurnState) int { return 2 }
func (ZealousBeltingRed) Pitch() int              { return 1 }
func (ZealousBeltingRed) Attack() int             { return 5 }
func (ZealousBeltingRed) Defense() int            { return 2 }
func (ZealousBeltingRed) Types() card.TypeSet     { return zealousBeltingTypes }
func (ZealousBeltingRed) GoAgain() bool           { return false }
func (ZealousBeltingRed) Play(s *sim.TurnState, self *sim.CardState) {
	zealousBeltingPlay(s, self)
}

type ZealousBeltingYellow struct{}

func (ZealousBeltingYellow) ID() ids.CardID          { return ids.ZealousBeltingYellow }
func (ZealousBeltingYellow) Name() string            { return "Zealous Belting" }
func (ZealousBeltingYellow) Cost(*sim.TurnState) int { return 2 }
func (ZealousBeltingYellow) Pitch() int              { return 2 }
func (ZealousBeltingYellow) Attack() int             { return 4 }
func (ZealousBeltingYellow) Defense() int            { return 2 }
func (ZealousBeltingYellow) Types() card.TypeSet     { return zealousBeltingTypes }
func (ZealousBeltingYellow) GoAgain() bool           { return false }
func (ZealousBeltingYellow) Play(s *sim.TurnState, self *sim.CardState) {
	zealousBeltingPlay(s, self)
}

type ZealousBeltingBlue struct{}

func (ZealousBeltingBlue) ID() ids.CardID          { return ids.ZealousBeltingBlue }
func (ZealousBeltingBlue) Name() string            { return "Zealous Belting" }
func (ZealousBeltingBlue) Cost(*sim.TurnState) int { return 2 }
func (ZealousBeltingBlue) Pitch() int              { return 3 }
func (ZealousBeltingBlue) Attack() int             { return 3 }
func (ZealousBeltingBlue) Defense() int            { return 2 }
func (ZealousBeltingBlue) Types() card.TypeSet     { return zealousBeltingTypes }
func (ZealousBeltingBlue) GoAgain() bool           { return false }
func (ZealousBeltingBlue) Play(s *sim.TurnState, self *sim.CardState) {
	zealousBeltingPlay(s, self)
}

func (ZealousBeltingRed) ConditionalGoAgain()    {}
func (ZealousBeltingYellow) ConditionalGoAgain() {}
func (ZealousBeltingBlue) ConditionalGoAgain()   {}
