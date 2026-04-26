// Fervent Forerunner — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Fervent Forerunner hits, **opt 2**. If Fervent Forerunner is played from arsenal, it
// gains **go again**."
//
// Modelling: on-hit Opt 2 isn't modelled. The played-from-arsenal go-again fires via
// self.GrantedGoAgain when self.FromArsenal reports this copy came from the arsenal slot.
// GoAgain() stays false so hand-played copies don't get the grant.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var ferventForerunnerTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// ferventForerunnerPlay grants self Go again when this copy was played from arsenal,
// then emits the chain step.
func ferventForerunnerPlay(s *card.TurnState, self *card.CardState) {
	if self.FromArsenal {
		self.GrantedGoAgain = true
	}
	s.ApplyAndLogEffectiveAttack(self)
}

type FerventForerunnerRed struct{}

func (FerventForerunnerRed) ID() card.ID              { return card.FerventForerunnerRed }
func (FerventForerunnerRed) Name() string             { return "Fervent Forerunner" }
func (FerventForerunnerRed) Cost(*card.TurnState) int { return 0 }
func (FerventForerunnerRed) Pitch() int               { return 1 }
func (FerventForerunnerRed) Attack() int              { return 3 }
func (FerventForerunnerRed) Defense() int             { return 2 }
func (FerventForerunnerRed) Types() card.TypeSet      { return ferventForerunnerTypes }
func (FerventForerunnerRed) GoAgain() bool            { return false }

// not implemented: on-hit Opt 2 rider
func (FerventForerunnerRed) NotImplemented() {}
func (FerventForerunnerRed) Play(s *card.TurnState, self *card.CardState) {
	ferventForerunnerPlay(s, self)
}

type FerventForerunnerYellow struct{}

func (FerventForerunnerYellow) ID() card.ID              { return card.FerventForerunnerYellow }
func (FerventForerunnerYellow) Name() string             { return "Fervent Forerunner" }
func (FerventForerunnerYellow) Cost(*card.TurnState) int { return 0 }
func (FerventForerunnerYellow) Pitch() int               { return 2 }
func (FerventForerunnerYellow) Attack() int              { return 2 }
func (FerventForerunnerYellow) Defense() int             { return 2 }
func (FerventForerunnerYellow) Types() card.TypeSet      { return ferventForerunnerTypes }
func (FerventForerunnerYellow) GoAgain() bool            { return false }

// not implemented: on-hit Opt 2 rider
func (FerventForerunnerYellow) NotImplemented() {}
func (FerventForerunnerYellow) Play(s *card.TurnState, self *card.CardState) {
	ferventForerunnerPlay(s, self)
}

type FerventForerunnerBlue struct{}

func (FerventForerunnerBlue) ID() card.ID              { return card.FerventForerunnerBlue }
func (FerventForerunnerBlue) Name() string             { return "Fervent Forerunner" }
func (FerventForerunnerBlue) Cost(*card.TurnState) int { return 0 }
func (FerventForerunnerBlue) Pitch() int               { return 3 }
func (FerventForerunnerBlue) Attack() int              { return 1 }
func (FerventForerunnerBlue) Defense() int             { return 2 }
func (FerventForerunnerBlue) Types() card.TypeSet      { return ferventForerunnerTypes }
func (FerventForerunnerBlue) GoAgain() bool            { return false }

// not implemented: on-hit Opt 2 rider
func (FerventForerunnerBlue) NotImplemented() {}
func (FerventForerunnerBlue) Play(s *card.TurnState, self *card.CardState) {
	ferventForerunnerPlay(s, self)
}
