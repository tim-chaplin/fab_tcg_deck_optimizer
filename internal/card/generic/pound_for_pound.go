// Pound for Pound — Generic Action - Attack. Cost 3. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When you play Pound for Pound, if you have less {h} than an opposing hero, it gains
// **dominate**."
//
// Modelling: the "less {h} than an opposing hero" clause is treated as a hero attribute — the
// Dominate grant fires for heroes that implement card.LowerHealthWanter (via
// simstate.HeroWantsLowerHealth) and never fires otherwise, a coarse proxy that skips per-turn
// life tracking. Standard self.GrantedDominate wiring (docs/dev-standards.md).

package generic

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/simstate"
)

var poundForPoundTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// poundForPoundPlay grants self Dominate when the current hero opts into LowerHealthWanter,
// then emits the chain step.
func poundForPoundPlay(s *card.TurnState, self *card.CardState) {
	if simstate.HeroWantsLowerHealth() {
		self.GrantedDominate = true
	}
	s.ApplyAndLogEffectiveAttack(self)
}

type PoundForPoundRed struct{}

func (PoundForPoundRed) ID() card.ID              { return card.PoundForPoundRed }
func (PoundForPoundRed) Name() string             { return "Pound for Pound" }
func (PoundForPoundRed) Cost(*card.TurnState) int { return 3 }
func (PoundForPoundRed) Pitch() int               { return 1 }
func (PoundForPoundRed) Attack() int              { return 6 }
func (PoundForPoundRed) Defense() int             { return 2 }
func (PoundForPoundRed) Types() card.TypeSet      { return poundForPoundTypes }
func (PoundForPoundRed) GoAgain() bool            { return false }
func (PoundForPoundRed) Play(s *card.TurnState, self *card.CardState) {
	poundForPoundPlay(s, self)
}

type PoundForPoundYellow struct{}

func (PoundForPoundYellow) ID() card.ID              { return card.PoundForPoundYellow }
func (PoundForPoundYellow) Name() string             { return "Pound for Pound" }
func (PoundForPoundYellow) Cost(*card.TurnState) int { return 3 }
func (PoundForPoundYellow) Pitch() int               { return 2 }
func (PoundForPoundYellow) Attack() int              { return 5 }
func (PoundForPoundYellow) Defense() int             { return 2 }
func (PoundForPoundYellow) Types() card.TypeSet      { return poundForPoundTypes }
func (PoundForPoundYellow) GoAgain() bool            { return false }
func (PoundForPoundYellow) Play(s *card.TurnState, self *card.CardState) {
	poundForPoundPlay(s, self)
}

type PoundForPoundBlue struct{}

func (PoundForPoundBlue) ID() card.ID              { return card.PoundForPoundBlue }
func (PoundForPoundBlue) Name() string             { return "Pound for Pound" }
func (PoundForPoundBlue) Cost(*card.TurnState) int { return 3 }
func (PoundForPoundBlue) Pitch() int               { return 3 }
func (PoundForPoundBlue) Attack() int              { return 4 }
func (PoundForPoundBlue) Defense() int             { return 2 }
func (PoundForPoundBlue) Types() card.TypeSet      { return poundForPoundTypes }
func (PoundForPoundBlue) GoAgain() bool            { return false }
func (PoundForPoundBlue) Play(s *card.TurnState, self *card.CardState) {
	poundForPoundPlay(s, self)
}
