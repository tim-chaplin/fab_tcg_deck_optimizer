// Put in Context — Generic Defense Reaction. Cost 0, Pitch 3, Defense 3. Only printed in Blue.
// Text: "This can only defend an attack with 3 or less base {p}."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

type PutInContextBlue struct{}

func (PutInContextBlue) ID() card.ID              { return card.PutInContextBlue }
func (PutInContextBlue) Name() string             { return "Put in Context" }
func (PutInContextBlue) Cost(*card.TurnState) int { return 0 }
func (PutInContextBlue) Pitch() int               { return 3 }
func (PutInContextBlue) Attack() int              { return 0 }
func (PutInContextBlue) Defense() int             { return 3 }
func (PutInContextBlue) Types() card.TypeSet      { return defenseReactionTypes }
func (PutInContextBlue) GoAgain() bool            { return false }

// not implemented: base-power cap on what this can defend is ignored; treated as legal vs every
// attack
func (PutInContextBlue) NotImplemented()                              {}
func (PutInContextBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
