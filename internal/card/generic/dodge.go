// Dodge — Generic Defense Reaction. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var defenseReactionTypes = map[string]bool{"Generic": true, "Defense Reaction": true}

type DodgeBlue struct{}

func (DodgeBlue) Name() string             { return "Dodge (Blue)" }
func (DodgeBlue) Cost() int                { return 0 }
func (DodgeBlue) Pitch() int               { return 3 }
func (DodgeBlue) Attack() int              { return 0 }
func (DodgeBlue) Defense() int             { return 2 }
func (DodgeBlue) Types() map[string]bool   { return defenseReactionTypes }
func (DodgeBlue) GoAgain() bool            { return false }
func (DodgeBlue) Play(*card.TurnState) int { return 0 }
