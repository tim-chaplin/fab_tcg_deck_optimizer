// Frontline Scout — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "You may look at the defending hero's hand. If Frontline Scout is played from arsenal, it
// gains **go again**."
//
// Modelling: hand-peek isn't modelled. Standard played-from-arsenal go-again
// (docs/dev-standards.md).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var frontlineScoutTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// frontlineScoutPlay grants self Go again when this copy was played from arsenal, then
// emits the chain step.
func frontlineScoutPlay(s *card.TurnState, self *card.CardState) {
	if self.FromArsenal {
		self.GrantedGoAgain = true
	}
	s.ApplyAndLogEffectiveAttack(self)
}

type FrontlineScoutRed struct{}

func (FrontlineScoutRed) ID() card.ID              { return card.FrontlineScoutRed }
func (FrontlineScoutRed) Name() string             { return "Frontline Scout" }
func (FrontlineScoutRed) Cost(*card.TurnState) int { return 0 }
func (FrontlineScoutRed) Pitch() int               { return 1 }
func (FrontlineScoutRed) Attack() int              { return 3 }
func (FrontlineScoutRed) Defense() int             { return 2 }
func (FrontlineScoutRed) Types() card.TypeSet      { return frontlineScoutTypes }
func (FrontlineScoutRed) GoAgain() bool            { return false }

// not implemented: opposing-hand-peek rider
func (FrontlineScoutRed) NotImplemented() {}
func (FrontlineScoutRed) Play(s *card.TurnState, self *card.CardState) {
	frontlineScoutPlay(s, self)
}

type FrontlineScoutYellow struct{}

func (FrontlineScoutYellow) ID() card.ID              { return card.FrontlineScoutYellow }
func (FrontlineScoutYellow) Name() string             { return "Frontline Scout" }
func (FrontlineScoutYellow) Cost(*card.TurnState) int { return 0 }
func (FrontlineScoutYellow) Pitch() int               { return 2 }
func (FrontlineScoutYellow) Attack() int              { return 2 }
func (FrontlineScoutYellow) Defense() int             { return 2 }
func (FrontlineScoutYellow) Types() card.TypeSet      { return frontlineScoutTypes }
func (FrontlineScoutYellow) GoAgain() bool            { return false }

// not implemented: opposing-hand-peek rider
func (FrontlineScoutYellow) NotImplemented() {}
func (FrontlineScoutYellow) Play(s *card.TurnState, self *card.CardState) {
	frontlineScoutPlay(s, self)
}

type FrontlineScoutBlue struct{}

func (FrontlineScoutBlue) ID() card.ID              { return card.FrontlineScoutBlue }
func (FrontlineScoutBlue) Name() string             { return "Frontline Scout" }
func (FrontlineScoutBlue) Cost(*card.TurnState) int { return 0 }
func (FrontlineScoutBlue) Pitch() int               { return 3 }
func (FrontlineScoutBlue) Attack() int              { return 1 }
func (FrontlineScoutBlue) Defense() int             { return 2 }
func (FrontlineScoutBlue) Types() card.TypeSet      { return frontlineScoutTypes }
func (FrontlineScoutBlue) GoAgain() bool            { return false }

// not implemented: opposing-hand-peek rider
func (FrontlineScoutBlue) NotImplemented() {}
func (FrontlineScoutBlue) Play(s *card.TurnState, self *card.CardState) {
	frontlineScoutPlay(s, self)
}
