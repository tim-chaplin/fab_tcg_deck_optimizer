// Battlefront Bastion — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this defends alone, prevent the next 1 damage that would be dealt to you this turn."
//
// Block scans s.Defenders for a second plain blocker; if none is present (DRs alongside
// don't count), the +1 prevention fires by bumping self.BonusDefense.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var battlefrontBastionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func battlefrontBastionPlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

// battlefrontBastionBlock fires the +1 alone-bonus when this is the only plain blocker.
// Iterates Defenders and short-circuits on the second non-DR sighting so the typical
// partition pays at most a few comparisons.
func battlefrontBastionBlock(s *sim.TurnState, self *sim.CardState) {
	plainCount := 0
	for _, d := range s.Defenders {
		if d.Types().IsDefenseReaction() {
			continue
		}
		plainCount++
		if plainCount > 1 {
			return
		}
	}
	self.BonusDefense += 1
}

type BattlefrontBastionRed struct{}

func (BattlefrontBastionRed) ID() ids.CardID          { return ids.BattlefrontBastionRed }
func (BattlefrontBastionRed) Name() string            { return "Battlefront Bastion" }
func (BattlefrontBastionRed) Cost(*sim.TurnState) int { return 3 }
func (BattlefrontBastionRed) Pitch() int              { return 1 }
func (BattlefrontBastionRed) Attack() int             { return 7 }
func (BattlefrontBastionRed) Defense() int            { return 2 }
func (BattlefrontBastionRed) Types() card.TypeSet     { return battlefrontBastionTypes }
func (BattlefrontBastionRed) GoAgain() bool           { return false }
func (BattlefrontBastionRed) Block(s *sim.TurnState, self *sim.CardState) {
	battlefrontBastionBlock(s, self)
}
func (BattlefrontBastionRed) Play(s *sim.TurnState, self *sim.CardState) {
	battlefrontBastionPlay(s, self)
}

type BattlefrontBastionYellow struct{}

func (BattlefrontBastionYellow) ID() ids.CardID          { return ids.BattlefrontBastionYellow }
func (BattlefrontBastionYellow) Name() string            { return "Battlefront Bastion" }
func (BattlefrontBastionYellow) Cost(*sim.TurnState) int { return 3 }
func (BattlefrontBastionYellow) Pitch() int              { return 2 }
func (BattlefrontBastionYellow) Attack() int             { return 6 }
func (BattlefrontBastionYellow) Defense() int            { return 2 }
func (BattlefrontBastionYellow) Types() card.TypeSet     { return battlefrontBastionTypes }
func (BattlefrontBastionYellow) GoAgain() bool           { return false }
func (BattlefrontBastionYellow) Block(s *sim.TurnState, self *sim.CardState) {
	battlefrontBastionBlock(s, self)
}
func (BattlefrontBastionYellow) Play(s *sim.TurnState, self *sim.CardState) {
	battlefrontBastionPlay(s, self)
}

type BattlefrontBastionBlue struct{}

func (BattlefrontBastionBlue) ID() ids.CardID          { return ids.BattlefrontBastionBlue }
func (BattlefrontBastionBlue) Name() string            { return "Battlefront Bastion" }
func (BattlefrontBastionBlue) Cost(*sim.TurnState) int { return 3 }
func (BattlefrontBastionBlue) Pitch() int              { return 3 }
func (BattlefrontBastionBlue) Attack() int             { return 5 }
func (BattlefrontBastionBlue) Defense() int            { return 2 }
func (BattlefrontBastionBlue) Types() card.TypeSet     { return battlefrontBastionTypes }
func (BattlefrontBastionBlue) GoAgain() bool           { return false }
func (BattlefrontBastionBlue) Block(s *sim.TurnState, self *sim.CardState) {
	battlefrontBastionBlock(s, self)
}
func (BattlefrontBastionBlue) Play(s *sim.TurnState, self *sim.CardState) {
	battlefrontBastionPlay(s, self)
}
