// Battlefront Bastion — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this defends alone, prevent the next 1 damage that would be dealt to you this turn."

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var battlefrontBastionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BattlefrontBastionRed struct{}

func (BattlefrontBastionRed) ID() card.ID              { return card.BattlefrontBastionRed }
func (BattlefrontBastionRed) Name() string             { return "Battlefront Bastion" }
func (BattlefrontBastionRed) Cost(*card.TurnState) int { return 3 }
func (BattlefrontBastionRed) Pitch() int               { return 1 }
func (BattlefrontBastionRed) Attack() int              { return 7 }
func (BattlefrontBastionRed) Defense() int             { return 2 }
func (BattlefrontBastionRed) Types() card.TypeSet      { return battlefrontBastionTypes }
func (BattlefrontBastionRed) GoAgain() bool            { return false }

// not implemented: defend-alone damage prevention rider
func (BattlefrontBastionRed) NotImplemented() {}
func (c BattlefrontBastionRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type BattlefrontBastionYellow struct{}

func (BattlefrontBastionYellow) ID() card.ID              { return card.BattlefrontBastionYellow }
func (BattlefrontBastionYellow) Name() string             { return "Battlefront Bastion" }
func (BattlefrontBastionYellow) Cost(*card.TurnState) int { return 3 }
func (BattlefrontBastionYellow) Pitch() int               { return 2 }
func (BattlefrontBastionYellow) Attack() int              { return 6 }
func (BattlefrontBastionYellow) Defense() int             { return 2 }
func (BattlefrontBastionYellow) Types() card.TypeSet      { return battlefrontBastionTypes }
func (BattlefrontBastionYellow) GoAgain() bool            { return false }

// not implemented: defend-alone damage prevention rider
func (BattlefrontBastionYellow) NotImplemented() {}
func (c BattlefrontBastionYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type BattlefrontBastionBlue struct{}

func (BattlefrontBastionBlue) ID() card.ID              { return card.BattlefrontBastionBlue }
func (BattlefrontBastionBlue) Name() string             { return "Battlefront Bastion" }
func (BattlefrontBastionBlue) Cost(*card.TurnState) int { return 3 }
func (BattlefrontBastionBlue) Pitch() int               { return 3 }
func (BattlefrontBastionBlue) Attack() int              { return 5 }
func (BattlefrontBastionBlue) Defense() int             { return 2 }
func (BattlefrontBastionBlue) Types() card.TypeSet      { return battlefrontBastionTypes }
func (BattlefrontBastionBlue) GoAgain() bool            { return false }

// not implemented: defend-alone damage prevention rider
func (BattlefrontBastionBlue) NotImplemented() {}
func (c BattlefrontBastionBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
