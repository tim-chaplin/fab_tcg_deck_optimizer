// Back Alley Breakline — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If an activated ability or action card effect puts Back Alley Breakline face up into a
// zone from your deck, gain 1 action point."

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var backAlleyBreaklineTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BackAlleyBreaklineRed struct{}

func (BackAlleyBreaklineRed) ID() card.ID              { return card.BackAlleyBreaklineRed }
func (BackAlleyBreaklineRed) Name() string             { return "Back Alley Breakline" }
func (BackAlleyBreaklineRed) Cost(*card.TurnState) int { return 1 }
func (BackAlleyBreaklineRed) Pitch() int               { return 1 }
func (BackAlleyBreaklineRed) Attack() int              { return 5 }
func (BackAlleyBreaklineRed) Defense() int             { return 2 }
func (BackAlleyBreaklineRed) Types() card.TypeSet      { return backAlleyBreaklineTypes }
func (BackAlleyBreaklineRed) GoAgain() bool            { return false }

// not implemented: face-up-from-deck action point grant
func (BackAlleyBreaklineRed) NotImplemented() {}
func (c BackAlleyBreaklineRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type BackAlleyBreaklineYellow struct{}

func (BackAlleyBreaklineYellow) ID() card.ID              { return card.BackAlleyBreaklineYellow }
func (BackAlleyBreaklineYellow) Name() string             { return "Back Alley Breakline" }
func (BackAlleyBreaklineYellow) Cost(*card.TurnState) int { return 1 }
func (BackAlleyBreaklineYellow) Pitch() int               { return 2 }
func (BackAlleyBreaklineYellow) Attack() int              { return 4 }
func (BackAlleyBreaklineYellow) Defense() int             { return 2 }
func (BackAlleyBreaklineYellow) Types() card.TypeSet      { return backAlleyBreaklineTypes }
func (BackAlleyBreaklineYellow) GoAgain() bool            { return false }

// not implemented: face-up-from-deck action point grant
func (BackAlleyBreaklineYellow) NotImplemented() {}
func (c BackAlleyBreaklineYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type BackAlleyBreaklineBlue struct{}

func (BackAlleyBreaklineBlue) ID() card.ID              { return card.BackAlleyBreaklineBlue }
func (BackAlleyBreaklineBlue) Name() string             { return "Back Alley Breakline" }
func (BackAlleyBreaklineBlue) Cost(*card.TurnState) int { return 1 }
func (BackAlleyBreaklineBlue) Pitch() int               { return 3 }
func (BackAlleyBreaklineBlue) Attack() int              { return 3 }
func (BackAlleyBreaklineBlue) Defense() int             { return 2 }
func (BackAlleyBreaklineBlue) Types() card.TypeSet      { return backAlleyBreaklineTypes }
func (BackAlleyBreaklineBlue) GoAgain() bool            { return false }

// not implemented: face-up-from-deck action point grant
func (BackAlleyBreaklineBlue) NotImplemented() {}
func (c BackAlleyBreaklineBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
