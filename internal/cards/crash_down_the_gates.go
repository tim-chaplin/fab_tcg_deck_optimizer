// Crash Down the Gates — Generic Action - Attack. Cost 3. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks a hero, they reveal the top card of their deck. If this has {p} greater
// than the revealed card, this gets +2{p}. When this hits a hero, destroy the top card of their
// deck."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var crashDownTheGatesTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type CrashDownTheGatesRed struct{}

func (CrashDownTheGatesRed) ID() ids.CardID          { return ids.CrashDownTheGatesRed }
func (CrashDownTheGatesRed) Name() string            { return "Crash Down the Gates" }
func (CrashDownTheGatesRed) Cost(*sim.TurnState) int { return 3 }
func (CrashDownTheGatesRed) Pitch() int              { return 1 }
func (CrashDownTheGatesRed) Attack() int             { return 6 }
func (CrashDownTheGatesRed) Defense() int            { return 2 }
func (CrashDownTheGatesRed) Types() card.TypeSet     { return crashDownTheGatesTypes }
func (CrashDownTheGatesRed) GoAgain() bool           { return false }

// not implemented: deck-reveal comparison + on-hit deck-top destruction
func (CrashDownTheGatesRed) NotImplemented() {}
func (CrashDownTheGatesRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type CrashDownTheGatesYellow struct{}

func (CrashDownTheGatesYellow) ID() ids.CardID          { return ids.CrashDownTheGatesYellow }
func (CrashDownTheGatesYellow) Name() string            { return "Crash Down the Gates" }
func (CrashDownTheGatesYellow) Cost(*sim.TurnState) int { return 3 }
func (CrashDownTheGatesYellow) Pitch() int              { return 2 }
func (CrashDownTheGatesYellow) Attack() int             { return 5 }
func (CrashDownTheGatesYellow) Defense() int            { return 2 }
func (CrashDownTheGatesYellow) Types() card.TypeSet     { return crashDownTheGatesTypes }
func (CrashDownTheGatesYellow) GoAgain() bool           { return false }

// not implemented: deck-reveal comparison + on-hit deck-top destruction
func (CrashDownTheGatesYellow) NotImplemented() {}
func (CrashDownTheGatesYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type CrashDownTheGatesBlue struct{}

func (CrashDownTheGatesBlue) ID() ids.CardID          { return ids.CrashDownTheGatesBlue }
func (CrashDownTheGatesBlue) Name() string            { return "Crash Down the Gates" }
func (CrashDownTheGatesBlue) Cost(*sim.TurnState) int { return 3 }
func (CrashDownTheGatesBlue) Pitch() int              { return 3 }
func (CrashDownTheGatesBlue) Attack() int             { return 4 }
func (CrashDownTheGatesBlue) Defense() int            { return 2 }
func (CrashDownTheGatesBlue) Types() card.TypeSet     { return crashDownTheGatesTypes }
func (CrashDownTheGatesBlue) GoAgain() bool           { return false }

// not implemented: deck-reveal comparison + on-hit deck-top destruction
func (CrashDownTheGatesBlue) NotImplemented() {}
func (CrashDownTheGatesBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
