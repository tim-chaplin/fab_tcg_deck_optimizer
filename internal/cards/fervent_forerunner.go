// Fervent Forerunner — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Fervent Forerunner hits, **opt 2**. If Fervent Forerunner is played from arsenal, it
// gains **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var ferventForerunnerTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// ferventForerunnerPlay grants self Go again when this copy was played from arsenal,
// emits the chain step, and resolves the on-hit Opt 2 (gated on LikelyToHit).
func ferventForerunnerPlay(s *sim.TurnState, self *sim.CardState) {
	if self.FromArsenal {
		self.GrantedGoAgain = true
	}
	s.ApplyAndLogEffectiveAttack(self)
	if sim.LikelyToHit(self) {
		s.Opt(2)
	}
}

type FerventForerunnerRed struct{}

func (FerventForerunnerRed) ID() ids.CardID          { return ids.FerventForerunnerRed }
func (FerventForerunnerRed) Name() string            { return "Fervent Forerunner" }
func (FerventForerunnerRed) Cost(*sim.TurnState) int { return 0 }
func (FerventForerunnerRed) Pitch() int              { return 1 }
func (FerventForerunnerRed) Attack() int             { return 3 }
func (FerventForerunnerRed) Defense() int            { return 2 }
func (FerventForerunnerRed) Types() card.TypeSet     { return ferventForerunnerTypes }
func (FerventForerunnerRed) GoAgain() bool           { return false }
func (FerventForerunnerRed) Play(s *sim.TurnState, self *sim.CardState) {
	ferventForerunnerPlay(s, self)
}

type FerventForerunnerYellow struct{}

func (FerventForerunnerYellow) ID() ids.CardID          { return ids.FerventForerunnerYellow }
func (FerventForerunnerYellow) Name() string            { return "Fervent Forerunner" }
func (FerventForerunnerYellow) Cost(*sim.TurnState) int { return 0 }
func (FerventForerunnerYellow) Pitch() int              { return 2 }
func (FerventForerunnerYellow) Attack() int             { return 2 }
func (FerventForerunnerYellow) Defense() int            { return 2 }
func (FerventForerunnerYellow) Types() card.TypeSet     { return ferventForerunnerTypes }
func (FerventForerunnerYellow) GoAgain() bool           { return false }
func (FerventForerunnerYellow) Play(s *sim.TurnState, self *sim.CardState) {
	ferventForerunnerPlay(s, self)
}

type FerventForerunnerBlue struct{}

func (FerventForerunnerBlue) ID() ids.CardID          { return ids.FerventForerunnerBlue }
func (FerventForerunnerBlue) Name() string            { return "Fervent Forerunner" }
func (FerventForerunnerBlue) Cost(*sim.TurnState) int { return 0 }
func (FerventForerunnerBlue) Pitch() int              { return 3 }
func (FerventForerunnerBlue) Attack() int             { return 1 }
func (FerventForerunnerBlue) Defense() int            { return 2 }
func (FerventForerunnerBlue) Types() card.TypeSet     { return ferventForerunnerTypes }
func (FerventForerunnerBlue) GoAgain() bool           { return false }
func (FerventForerunnerBlue) Play(s *sim.TurnState, self *sim.CardState) {
	ferventForerunnerPlay(s, self)
}
