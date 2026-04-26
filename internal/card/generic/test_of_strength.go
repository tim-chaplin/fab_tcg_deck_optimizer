// Test of Strength — Generic Block. Cost 0, Pitch 1, Defense 4. Only printed in Red.
//
// Text: "When this defends, **clash** with the attacking hero. The winner creates a Gold token."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

type TestOfStrengthRed struct{}

func (TestOfStrengthRed) ID() card.ID                { return card.TestOfStrengthRed }
func (TestOfStrengthRed) Name() string               { return "Test of Strength" }
func (TestOfStrengthRed) Cost(*card.TurnState) int   { return 0 }
func (TestOfStrengthRed) Pitch() int                 { return 1 }
func (TestOfStrengthRed) Attack() int                { return 0 }
func (TestOfStrengthRed) Defense() int               { return 4 }
func (TestOfStrengthRed) Types() card.TypeSet        { return defenseReactionTypes }
func (TestOfStrengthRed) GoAgain() bool              { return false }
// not implemented: gold tokens
func (TestOfStrengthRed) NotImplemented()            {}
func (TestOfStrengthRed) Play(s *card.TurnState, _ *card.CardState) { card.ClashValue(s, card.GoldTokenValue) }