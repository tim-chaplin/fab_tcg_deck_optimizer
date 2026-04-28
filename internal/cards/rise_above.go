// Rise Above — Generic Defense Reaction. Cost 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "You may put a card from your hand on top of your deck rather than pay Rise Above's {r}
// cost."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

type RiseAboveRed struct{}

func (RiseAboveRed) ID() ids.CardID           { return ids.RiseAboveRed }
func (RiseAboveRed) Name() string             { return "Rise Above" }
func (RiseAboveRed) Cost(*card.TurnState) int { return 2 }
func (RiseAboveRed) Pitch() int               { return 1 }
func (RiseAboveRed) Attack() int              { return 0 }
func (RiseAboveRed) Defense() int             { return 4 }
func (RiseAboveRed) Types() card.TypeSet      { return defenseReactionTypes }
func (RiseAboveRed) GoAgain() bool            { return false }

// not implemented: hand-as-cost alt cost not modelled; card fails when printed cost can't be paid
func (RiseAboveRed) NotImplemented() {}
func (RiseAboveRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}

type RiseAboveYellow struct{}

func (RiseAboveYellow) ID() ids.CardID           { return ids.RiseAboveYellow }
func (RiseAboveYellow) Name() string             { return "Rise Above" }
func (RiseAboveYellow) Cost(*card.TurnState) int { return 2 }
func (RiseAboveYellow) Pitch() int               { return 2 }
func (RiseAboveYellow) Attack() int              { return 0 }
func (RiseAboveYellow) Defense() int             { return 3 }
func (RiseAboveYellow) Types() card.TypeSet      { return defenseReactionTypes }
func (RiseAboveYellow) GoAgain() bool            { return false }

// not implemented: hand-as-cost alt cost not modelled; card fails when printed cost can't be paid
func (RiseAboveYellow) NotImplemented() {}
func (RiseAboveYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}

type RiseAboveBlue struct{}

func (RiseAboveBlue) ID() ids.CardID           { return ids.RiseAboveBlue }
func (RiseAboveBlue) Name() string             { return "Rise Above" }
func (RiseAboveBlue) Cost(*card.TurnState) int { return 2 }
func (RiseAboveBlue) Pitch() int               { return 3 }
func (RiseAboveBlue) Attack() int              { return 0 }
func (RiseAboveBlue) Defense() int             { return 2 }
func (RiseAboveBlue) Types() card.TypeSet      { return defenseReactionTypes }
func (RiseAboveBlue) GoAgain() bool            { return false }

// not implemented: hand-as-cost alt cost not modelled; card fails when printed cost can't be paid
func (RiseAboveBlue) NotImplemented() {}
func (RiseAboveBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}
