// Ransack and Raze — Generic Action. Cost X, Pitch 3, Defense 3. Only printed in Blue.
//
// Text: "Destroy target landmark with cost X. Create X Gold tokens. **Go again**"
//
// X cost treated as 0.

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var ransackAndRazeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type RansackAndRazeBlue struct{}

func (RansackAndRazeBlue) ID() card.ID              { return card.RansackAndRazeBlue }
func (RansackAndRazeBlue) Name() string             { return "Ransack and Raze" }
func (RansackAndRazeBlue) Cost(*card.TurnState) int { return 0 }
func (RansackAndRazeBlue) Pitch() int               { return 3 }
func (RansackAndRazeBlue) Attack() int              { return 0 }
func (RansackAndRazeBlue) Defense() int             { return 3 }
func (RansackAndRazeBlue) Types() card.TypeSet      { return ransackAndRazeTypes }
func (RansackAndRazeBlue) GoAgain() bool            { return true }

// not implemented: gold tokens, landmarks
func (RansackAndRazeBlue) NotImplemented()                              {}
func (RansackAndRazeBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
