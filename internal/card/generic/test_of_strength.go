// Test of Strength — Generic Block. Cost 0, Pitch 1, Defense 4. Only printed in Red.
//
// Text: "When this defends, **clash** with the attacking hero. The winner creates a Gold token."
//
// Clash is modelled via card.ClashValue with a fixed-opponent heuristic: we win (+GoldTokenValue)
// when our deck's top-card attack is 6 or 7, tie (0) at 5, lose (-GoldTokenValue — the Gold
// token accrues to the opponent instead) at 4 or below. Block is typed as a Defense Reaction so
// the solver invokes Play during the defensive chain, making state.Deck available for the peek.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

type TestOfStrengthRed struct{}

func (TestOfStrengthRed) ID() card.ID                { return card.TestOfStrengthRed }
func (TestOfStrengthRed) Name() string               { return "Test of Strength (Red)" }
func (TestOfStrengthRed) Cost() int                  { return 0 }
func (TestOfStrengthRed) Pitch() int                 { return 1 }
func (TestOfStrengthRed) Attack() int                { return 0 }
func (TestOfStrengthRed) Defense() int               { return 4 }
func (TestOfStrengthRed) Types() card.TypeSet        { return defenseReactionTypes }
func (TestOfStrengthRed) GoAgain() bool              { return false }
func (TestOfStrengthRed) Play(s *card.TurnState) int { return card.ClashValue(s, card.GoldTokenValue) }
