// Amulet of Oblation — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** **Instant** - Destroy Amulet of Oblation: Until end of turn, target attack
// action gains "If this would be put into a graveyard, instead put it on the bottom of its owner's
// deck." Activate this ability only if a card has entered a graveyard this turn."
//
// Marked sim.Unplayable: the card itself is too weak to want in a deck. Best-case output is
// recycling one attack action back to the bottom of the deck — saves the card from the
// graveyard, but a deck-bottom recycle is worth ~1-2 future-turn value at most, and you've
// spent a card slot to get it. Even fully modelled the EV doesn't beat the slot cost; the
// per-turn evaluation caveat is a secondary modelling concern but not the deciding factor.

package cards

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
