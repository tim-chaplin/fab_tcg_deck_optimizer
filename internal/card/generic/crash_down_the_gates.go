// Crash Down the Gates — Generic Action - Attack. Cost 3. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks a hero, they reveal the top card of their deck. If this has {p} greater
// than the revealed card, this gets +2{p}. When this hits a hero, destroy the top card of their
// deck."
//
// Simplification: Deck-reveal comparison and top-of-deck destruction aren't modelled.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var crashDownTheGatesTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type CrashDownTheGatesRed struct{}

func (CrashDownTheGatesRed) ID() card.ID                 { return card.CrashDownTheGatesRed }
func (CrashDownTheGatesRed) Name() string                { return "Crash Down the Gates (Red)" }
func (CrashDownTheGatesRed) Cost(*card.TurnState) int                   { return 3 }
func (CrashDownTheGatesRed) Pitch() int                  { return 1 }
func (CrashDownTheGatesRed) Attack() int                 { return 6 }
func (CrashDownTheGatesRed) Defense() int                { return 2 }
func (CrashDownTheGatesRed) Types() card.TypeSet         { return crashDownTheGatesTypes }
func (CrashDownTheGatesRed) GoAgain() bool               { return false }
func (c CrashDownTheGatesRed) Play(s *card.TurnState, self *card.CardState) int { return crashDownTheGatesDamage(c.Attack(), self) }

type CrashDownTheGatesYellow struct{}

func (CrashDownTheGatesYellow) ID() card.ID                 { return card.CrashDownTheGatesYellow }
func (CrashDownTheGatesYellow) Name() string                { return "Crash Down the Gates (Yellow)" }
func (CrashDownTheGatesYellow) Cost(*card.TurnState) int                   { return 3 }
func (CrashDownTheGatesYellow) Pitch() int                  { return 2 }
func (CrashDownTheGatesYellow) Attack() int                 { return 5 }
func (CrashDownTheGatesYellow) Defense() int                { return 2 }
func (CrashDownTheGatesYellow) Types() card.TypeSet         { return crashDownTheGatesTypes }
func (CrashDownTheGatesYellow) GoAgain() bool               { return false }
func (c CrashDownTheGatesYellow) Play(s *card.TurnState, self *card.CardState) int { return crashDownTheGatesDamage(c.Attack(), self) }

type CrashDownTheGatesBlue struct{}

func (CrashDownTheGatesBlue) ID() card.ID                 { return card.CrashDownTheGatesBlue }
func (CrashDownTheGatesBlue) Name() string                { return "Crash Down the Gates (Blue)" }
func (CrashDownTheGatesBlue) Cost(*card.TurnState) int                   { return 3 }
func (CrashDownTheGatesBlue) Pitch() int                  { return 3 }
func (CrashDownTheGatesBlue) Attack() int                 { return 4 }
func (CrashDownTheGatesBlue) Defense() int                { return 2 }
func (CrashDownTheGatesBlue) Types() card.TypeSet         { return crashDownTheGatesTypes }
func (CrashDownTheGatesBlue) GoAgain() bool               { return false }
func (c CrashDownTheGatesBlue) Play(s *card.TurnState, self *card.CardState) int { return crashDownTheGatesDamage(c.Attack(), self) }

// crashDownTheGatesDamage is a breadcrumb for the on-hit "destroy top of their deck" rider —
// not modelled yet (see TODO.md).
func crashDownTheGatesDamage(attack int, self *card.CardState) int {
	if card.LikelyToHit(attack, self.EffectiveDominate()) {
		// TODO: model on-hit deck-top destruction rider.
	}
	return attack
}
