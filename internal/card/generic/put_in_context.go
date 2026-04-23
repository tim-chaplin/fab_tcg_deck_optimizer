// Put in Context — Generic Defense Reaction. Cost 0, Pitch 3, Defense 3. Only printed in Blue.
// Text: "This can only defend an attack with 3 or less base {p}."
// Simplification: the base-power cap on which attack this can block is ignored — we assume every
// attack qualifies and treat the defense as unconditional.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

type PutInContextBlue struct{}

func (PutInContextBlue) ID() card.ID                 { return card.PutInContextBlue }
func (PutInContextBlue) Name() string             { return "Put in Context (Blue)" }
func (PutInContextBlue) Cost(*card.TurnState) int                { return 0 }
func (PutInContextBlue) Pitch() int               { return 3 }
func (PutInContextBlue) Attack() int              { return 0 }
func (PutInContextBlue) Defense() int             { return 3 }
func (PutInContextBlue) Types() card.TypeSet      { return defenseReactionTypes }
func (PutInContextBlue) GoAgain() bool            { return false }
func (PutInContextBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
