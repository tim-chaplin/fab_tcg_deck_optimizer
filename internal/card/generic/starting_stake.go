// Starting Stake — Generic Action. Cost 0, Pitch 2, Defense 3. Only printed in Yellow.
//
// Text: "If you control no Gold tokens, create a Gold token."
//
// Simplification: Gold-token economy isn't tracked.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var startingStakeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type StartingStakeYellow struct{}

func (StartingStakeYellow) ID() card.ID                 { return card.StartingStakeYellow }
func (StartingStakeYellow) Name() string                { return "Starting Stake (Yellow)" }
func (StartingStakeYellow) Cost() int                   { return 0 }
func (StartingStakeYellow) Pitch() int                  { return 2 }
func (StartingStakeYellow) Attack() int                 { return 0 }
func (StartingStakeYellow) Defense() int                { return 3 }
func (StartingStakeYellow) Types() card.TypeSet         { return startingStakeTypes }
func (StartingStakeYellow) GoAgain() bool               { return false }
func (StartingStakeYellow) Play(s *card.TurnState) int { return 0 }
