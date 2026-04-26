// Unmovable — Generic Defense Reaction. Cost 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 7, Yellow 6, Blue 5.
// Text: "If Unmovable is played from arsenal, it gains +1{d}."
//
// Modelling: The +1{d} rider opts in via card.ArsenalDefenseBonus; CardState.EffectiveDefense
// folds it in only when this copy was the start-of-turn arsenal-in card.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

type UnmovableRed struct{}

func (UnmovableRed) ID() card.ID              { return card.UnmovableRed }
func (UnmovableRed) Name() string             { return "Unmovable" }
func (UnmovableRed) Cost(*card.TurnState) int { return 3 }
func (UnmovableRed) Pitch() int               { return 1 }
func (UnmovableRed) Attack() int              { return 0 }
func (UnmovableRed) Defense() int             { return 7 }
func (UnmovableRed) Types() card.TypeSet      { return defenseReactionTypes }
func (UnmovableRed) GoAgain() bool            { return false }
func (UnmovableRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}
func (UnmovableRed) ArsenalDefenseBonus() int { return 1 }

type UnmovableYellow struct{}

func (UnmovableYellow) ID() card.ID              { return card.UnmovableYellow }
func (UnmovableYellow) Name() string             { return "Unmovable" }
func (UnmovableYellow) Cost(*card.TurnState) int { return 3 }
func (UnmovableYellow) Pitch() int               { return 2 }
func (UnmovableYellow) Attack() int              { return 0 }
func (UnmovableYellow) Defense() int             { return 6 }
func (UnmovableYellow) Types() card.TypeSet      { return defenseReactionTypes }
func (UnmovableYellow) GoAgain() bool            { return false }
func (UnmovableYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}
func (UnmovableYellow) ArsenalDefenseBonus() int { return 1 }

type UnmovableBlue struct{}

func (UnmovableBlue) ID() card.ID              { return card.UnmovableBlue }
func (UnmovableBlue) Name() string             { return "Unmovable" }
func (UnmovableBlue) Cost(*card.TurnState) int { return 3 }
func (UnmovableBlue) Pitch() int               { return 3 }
func (UnmovableBlue) Attack() int              { return 0 }
func (UnmovableBlue) Defense() int             { return 5 }
func (UnmovableBlue) Types() card.TypeSet      { return defenseReactionTypes }
func (UnmovableBlue) GoAgain() bool            { return false }
func (UnmovableBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}
func (UnmovableBlue) ArsenalDefenseBonus() int { return 1 }
