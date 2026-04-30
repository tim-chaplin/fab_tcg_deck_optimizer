// Fate Foreseen — Generic Defense Reaction. Cost 0.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "Opt 1"
//
// Opt 1 is credited at sim.OptValue on top of the printed defense.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// fateForeseenPlay applies the printed defense and credits the Opt 1 rider as a sub-line.
func fateForeseenPlay(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
	s.ApplyAndLogRiderOnPlay(self, "Opt 1", sim.OptValue)
}

type FateForeseenRed struct{}

func (FateForeseenRed) ID() ids.CardID          { return ids.FateForeseenRed }
func (FateForeseenRed) Name() string            { return "Fate Foreseen" }
func (FateForeseenRed) Cost(*sim.TurnState) int { return 0 }
func (FateForeseenRed) Pitch() int              { return 1 }
func (FateForeseenRed) Attack() int             { return 0 }
func (FateForeseenRed) Defense() int            { return 4 }
func (FateForeseenRed) Types() card.TypeSet     { return defenseReactionTypes }
func (FateForeseenRed) GoAgain() bool           { return false }
func (FateForeseenRed) NotSilverAgeLegal()      {}
func (FateForeseenRed) Play(s *sim.TurnState, self *sim.CardState) {
	fateForeseenPlay(s, self)
}

type FateForeseenYellow struct{}

func (FateForeseenYellow) ID() ids.CardID          { return ids.FateForeseenYellow }
func (FateForeseenYellow) Name() string            { return "Fate Foreseen" }
func (FateForeseenYellow) Cost(*sim.TurnState) int { return 0 }
func (FateForeseenYellow) Pitch() int              { return 2 }
func (FateForeseenYellow) Attack() int             { return 0 }
func (FateForeseenYellow) Defense() int            { return 3 }
func (FateForeseenYellow) Types() card.TypeSet     { return defenseReactionTypes }
func (FateForeseenYellow) GoAgain() bool           { return false }
func (FateForeseenYellow) NotSilverAgeLegal()      {}
func (FateForeseenYellow) Play(s *sim.TurnState, self *sim.CardState) {
	fateForeseenPlay(s, self)
}

type FateForeseenBlue struct{}

func (FateForeseenBlue) ID() ids.CardID          { return ids.FateForeseenBlue }
func (FateForeseenBlue) Name() string            { return "Fate Foreseen" }
func (FateForeseenBlue) Cost(*sim.TurnState) int { return 0 }
func (FateForeseenBlue) Pitch() int              { return 3 }
func (FateForeseenBlue) Attack() int             { return 0 }
func (FateForeseenBlue) Defense() int            { return 2 }
func (FateForeseenBlue) Types() card.TypeSet     { return defenseReactionTypes }
func (FateForeseenBlue) GoAgain() bool           { return false }
func (FateForeseenBlue) NotSilverAgeLegal()      {}
func (FateForeseenBlue) Play(s *sim.TurnState, self *sim.CardState) {
	fateForeseenPlay(s, self)
}
