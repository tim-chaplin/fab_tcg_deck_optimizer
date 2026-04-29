// Spring Load — Generic Action - Attack. Cost 1. Printed power: Red 2, Yellow 2, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, if you have no cards in hand, it gains +3{p}."
//
// "No cards in hand" reads `len(s.Hand) == 0`: s.Hand at chain-resolution time holds the Held
// cards (post-pitch, post-attacker, post-defender), so the rider fires precisely when every
// hand card is committed to the turn — pitched, played, or defending — and nothing is sitting
// in hand for arsenal-promotion at end of turn.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var springLoadTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// springLoadPlay applies the +3{p} 'no cards in hand' rider before crediting the attack so
// the bump folds into EffectiveAttack (and downstream LikelyToHit windows for any rider that
// reads the post-buff power).
func springLoadPlay(s *sim.TurnState, self *sim.CardState) {
	if len(s.Hand) == 0 {
		self.BonusAttack += 3
	}
	s.ApplyAndLogEffectiveAttack(self)
}

type SpringLoadRed struct{}

func (SpringLoadRed) ID() ids.CardID          { return ids.SpringLoadRed }
func (SpringLoadRed) Name() string            { return "Spring Load" }
func (SpringLoadRed) Cost(*sim.TurnState) int { return 1 }
func (SpringLoadRed) Pitch() int              { return 1 }
func (SpringLoadRed) Attack() int             { return 2 }
func (SpringLoadRed) Defense() int            { return 2 }
func (SpringLoadRed) Types() card.TypeSet     { return springLoadTypes }
func (SpringLoadRed) GoAgain() bool           { return false }
func (SpringLoadRed) Play(s *sim.TurnState, self *sim.CardState) {
	springLoadPlay(s, self)
}

type SpringLoadYellow struct{}

func (SpringLoadYellow) ID() ids.CardID          { return ids.SpringLoadYellow }
func (SpringLoadYellow) Name() string            { return "Spring Load" }
func (SpringLoadYellow) Cost(*sim.TurnState) int { return 1 }
func (SpringLoadYellow) Pitch() int              { return 2 }
func (SpringLoadYellow) Attack() int             { return 2 }
func (SpringLoadYellow) Defense() int            { return 2 }
func (SpringLoadYellow) Types() card.TypeSet     { return springLoadTypes }
func (SpringLoadYellow) GoAgain() bool           { return false }
func (SpringLoadYellow) Play(s *sim.TurnState, self *sim.CardState) {
	springLoadPlay(s, self)
}

type SpringLoadBlue struct{}

func (SpringLoadBlue) ID() ids.CardID          { return ids.SpringLoadBlue }
func (SpringLoadBlue) Name() string            { return "Spring Load" }
func (SpringLoadBlue) Cost(*sim.TurnState) int { return 1 }
func (SpringLoadBlue) Pitch() int              { return 3 }
func (SpringLoadBlue) Attack() int             { return 2 }
func (SpringLoadBlue) Defense() int            { return 2 }
func (SpringLoadBlue) Types() card.TypeSet     { return springLoadTypes }
func (SpringLoadBlue) GoAgain() bool           { return false }
func (SpringLoadBlue) Play(s *sim.TurnState, self *sim.CardState) {
	springLoadPlay(s, self)
}
