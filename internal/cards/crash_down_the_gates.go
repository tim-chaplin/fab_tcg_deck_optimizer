// Crash Down the Gates — Generic Action - Attack. Cost 3. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks a hero, they reveal the top card of their deck. If this has {p} greater
// than the revealed card, this gets +2{p}. When this hits a hero, destroy the top card of their
// deck."

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var crashDownTheGatesTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type CrashDownTheGatesRed struct{}

func (CrashDownTheGatesRed) ID() card.ID              { return card.CrashDownTheGatesRed }
func (CrashDownTheGatesRed) Name() string             { return "Crash Down the Gates" }
func (CrashDownTheGatesRed) Cost(*card.TurnState) int { return 3 }
func (CrashDownTheGatesRed) Pitch() int               { return 1 }
func (CrashDownTheGatesRed) Attack() int              { return 6 }
func (CrashDownTheGatesRed) Defense() int             { return 2 }
func (CrashDownTheGatesRed) Types() card.TypeSet      { return crashDownTheGatesTypes }
func (CrashDownTheGatesRed) GoAgain() bool            { return false }

// not implemented: deck-reveal comparison + on-hit deck-top destruction
func (CrashDownTheGatesRed) NotImplemented() {}
func (CrashDownTheGatesRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type CrashDownTheGatesYellow struct{}

func (CrashDownTheGatesYellow) ID() card.ID              { return card.CrashDownTheGatesYellow }
func (CrashDownTheGatesYellow) Name() string             { return "Crash Down the Gates" }
func (CrashDownTheGatesYellow) Cost(*card.TurnState) int { return 3 }
func (CrashDownTheGatesYellow) Pitch() int               { return 2 }
func (CrashDownTheGatesYellow) Attack() int              { return 5 }
func (CrashDownTheGatesYellow) Defense() int             { return 2 }
func (CrashDownTheGatesYellow) Types() card.TypeSet      { return crashDownTheGatesTypes }
func (CrashDownTheGatesYellow) GoAgain() bool            { return false }

// not implemented: deck-reveal comparison + on-hit deck-top destruction
func (CrashDownTheGatesYellow) NotImplemented() {}
func (CrashDownTheGatesYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type CrashDownTheGatesBlue struct{}

func (CrashDownTheGatesBlue) ID() card.ID              { return card.CrashDownTheGatesBlue }
func (CrashDownTheGatesBlue) Name() string             { return "Crash Down the Gates" }
func (CrashDownTheGatesBlue) Cost(*card.TurnState) int { return 3 }
func (CrashDownTheGatesBlue) Pitch() int               { return 3 }
func (CrashDownTheGatesBlue) Attack() int              { return 4 }
func (CrashDownTheGatesBlue) Defense() int             { return 2 }
func (CrashDownTheGatesBlue) Types() card.TypeSet      { return crashDownTheGatesTypes }
func (CrashDownTheGatesBlue) GoAgain() bool            { return false }

// not implemented: deck-reveal comparison + on-hit deck-top destruction
func (CrashDownTheGatesBlue) NotImplemented() {}
func (CrashDownTheGatesBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
