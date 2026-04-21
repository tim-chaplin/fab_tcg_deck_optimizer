// Dodge — Generic Defense Reaction. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var defenseReactionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeDefenseReaction)

type DodgeBlue struct{}

func (DodgeBlue) ID() card.ID                 { return card.DodgeBlue }
func (DodgeBlue) Name() string             { return "Dodge (Blue)" }
func (DodgeBlue) Cost(*card.TurnState) int                { return 0 }
func (DodgeBlue) Pitch() int               { return 3 }
func (DodgeBlue) Attack() int              { return 0 }
func (DodgeBlue) Defense() int             { return 2 }
func (DodgeBlue) Types() card.TypeSet      { return defenseReactionTypes }
func (DodgeBlue) GoAgain() bool            { return false }
func (DodgeBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
