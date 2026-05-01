// Amulet of Oblation — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** **Instant** - Destroy Amulet of Oblation: Until end of turn, target attack
// action gains "If this would be put into a graveyard, instead put it on the bottom of its owner's
// deck." Activate this ability only if a card has entered a graveyard this turn."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var amuletOfOblationTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type AmuletOfOblationBlue struct{}

func (AmuletOfOblationBlue) ID() ids.CardID          { return ids.AmuletOfOblationBlue }
func (AmuletOfOblationBlue) Name() string            { return "Amulet of Oblation" }
func (AmuletOfOblationBlue) Cost(*sim.TurnState) int { return 0 }
func (AmuletOfOblationBlue) Pitch() int              { return 3 }
func (AmuletOfOblationBlue) Attack() int             { return 0 }
func (AmuletOfOblationBlue) Defense() int            { return 0 }
func (AmuletOfOblationBlue) Types() card.TypeSet     { return amuletOfOblationTypes }
func (AmuletOfOblationBlue) GoAgain() bool           { return true }

func (AmuletOfOblationBlue) Unplayable()                                {}
func (AmuletOfOblationBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
