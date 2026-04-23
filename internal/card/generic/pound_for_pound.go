// Pound for Pound — Generic Action - Attack. Cost 3. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When you play Pound for Pound, if you have less {h} than an opposing hero, it gains
// **dominate**."
//
// Modelling: the "less {h} than an opposing hero" clause is modelled as a hero attribute — the
// Dominate grant fires for heroes that implement card.LowerHealthWanter (via
// simstate.HeroWantsLowerHealth) and never fires otherwise, the same coarse proxy Adrenaline
// Rush / Blow for a Blow / Wounded Bull use to skip per-turn life tracking. Play flips
// self.GrantedDominate so downstream scanners and any future on-hit rider see the grant.

package generic

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/simstate"
)

var poundForPoundTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// poundForPoundPlay grants self Dominate when the current hero opts into LowerHealthWanter,
// then returns the base attack.
func poundForPoundPlay(attack int, self *card.CardState) int {
	if simstate.HeroWantsLowerHealth() {
		self.GrantedDominate = true
	}
	return attack
}

type PoundForPoundRed struct{}

func (PoundForPoundRed) ID() card.ID                 { return card.PoundForPoundRed }
func (PoundForPoundRed) Name() string                { return "Pound for Pound (Red)" }
func (PoundForPoundRed) Cost(*card.TurnState) int    { return 3 }
func (PoundForPoundRed) Pitch() int                  { return 1 }
func (PoundForPoundRed) Attack() int                 { return 6 }
func (PoundForPoundRed) Defense() int                { return 2 }
func (PoundForPoundRed) Types() card.TypeSet         { return poundForPoundTypes }
func (PoundForPoundRed) GoAgain() bool               { return false }
func (c PoundForPoundRed) Play(_ *card.TurnState, self *card.CardState) int {
	return poundForPoundPlay(c.Attack(), self)
}

type PoundForPoundYellow struct{}

func (PoundForPoundYellow) ID() card.ID                 { return card.PoundForPoundYellow }
func (PoundForPoundYellow) Name() string                { return "Pound for Pound (Yellow)" }
func (PoundForPoundYellow) Cost(*card.TurnState) int    { return 3 }
func (PoundForPoundYellow) Pitch() int                  { return 2 }
func (PoundForPoundYellow) Attack() int                 { return 5 }
func (PoundForPoundYellow) Defense() int                { return 2 }
func (PoundForPoundYellow) Types() card.TypeSet         { return poundForPoundTypes }
func (PoundForPoundYellow) GoAgain() bool               { return false }
func (c PoundForPoundYellow) Play(_ *card.TurnState, self *card.CardState) int {
	return poundForPoundPlay(c.Attack(), self)
}

type PoundForPoundBlue struct{}

func (PoundForPoundBlue) ID() card.ID                 { return card.PoundForPoundBlue }
func (PoundForPoundBlue) Name() string                { return "Pound for Pound (Blue)" }
func (PoundForPoundBlue) Cost(*card.TurnState) int    { return 3 }
func (PoundForPoundBlue) Pitch() int                  { return 3 }
func (PoundForPoundBlue) Attack() int                 { return 4 }
func (PoundForPoundBlue) Defense() int                { return 2 }
func (PoundForPoundBlue) Types() card.TypeSet         { return poundForPoundTypes }
func (PoundForPoundBlue) GoAgain() bool               { return false }
func (c PoundForPoundBlue) Play(_ *card.TurnState, self *card.CardState) int {
	return poundForPoundPlay(c.Attack(), self)
}
