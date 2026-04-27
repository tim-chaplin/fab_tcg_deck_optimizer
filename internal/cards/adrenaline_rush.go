// Adrenaline Rush — Generic Action - Attack. Cost 2. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When you play this, if you have less {h} than an opposing hero, this gets +3{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/simstate"
)

var adrenalineRushTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// adrenalineRushBonus returns the +3{p} rider when the current hero opts into
// LowerHealthWanter, else 0.
func adrenalineRushBonus() int {
	if simstate.HeroWantsLowerHealth() {
		return 3
	}
	return 0
}

type AdrenalineRushRed struct{}

func (AdrenalineRushRed) ID() card.ID              { return card.AdrenalineRushRed }
func (AdrenalineRushRed) Name() string             { return "Adrenaline Rush" }
func (AdrenalineRushRed) Cost(*card.TurnState) int { return 2 }
func (AdrenalineRushRed) Pitch() int               { return 1 }
func (AdrenalineRushRed) Attack() int              { return 4 }
func (AdrenalineRushRed) Defense() int             { return 2 }
func (AdrenalineRushRed) Types() card.TypeSet      { return adrenalineRushTypes }
func (AdrenalineRushRed) GoAgain() bool            { return false }
func (AdrenalineRushRed) Play(s *card.TurnState, self *card.CardState) {
	self.BonusAttack += adrenalineRushBonus()
	s.ApplyAndLogEffectiveAttack(self)
}

type AdrenalineRushYellow struct{}

func (AdrenalineRushYellow) ID() card.ID              { return card.AdrenalineRushYellow }
func (AdrenalineRushYellow) Name() string             { return "Adrenaline Rush" }
func (AdrenalineRushYellow) Cost(*card.TurnState) int { return 2 }
func (AdrenalineRushYellow) Pitch() int               { return 2 }
func (AdrenalineRushYellow) Attack() int              { return 3 }
func (AdrenalineRushYellow) Defense() int             { return 2 }
func (AdrenalineRushYellow) Types() card.TypeSet      { return adrenalineRushTypes }
func (AdrenalineRushYellow) GoAgain() bool            { return false }
func (AdrenalineRushYellow) Play(s *card.TurnState, self *card.CardState) {
	self.BonusAttack += adrenalineRushBonus()
	s.ApplyAndLogEffectiveAttack(self)
}

type AdrenalineRushBlue struct{}

func (AdrenalineRushBlue) ID() card.ID              { return card.AdrenalineRushBlue }
func (AdrenalineRushBlue) Name() string             { return "Adrenaline Rush" }
func (AdrenalineRushBlue) Cost(*card.TurnState) int { return 2 }
func (AdrenalineRushBlue) Pitch() int               { return 3 }
func (AdrenalineRushBlue) Attack() int              { return 2 }
func (AdrenalineRushBlue) Defense() int             { return 2 }
func (AdrenalineRushBlue) Types() card.TypeSet      { return adrenalineRushTypes }
func (AdrenalineRushBlue) GoAgain() bool            { return false }
func (AdrenalineRushBlue) Play(s *card.TurnState, self *card.CardState) {
	self.BonusAttack += adrenalineRushBonus()
	s.ApplyAndLogEffectiveAttack(self)
}
