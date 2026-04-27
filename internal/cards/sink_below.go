// Sink Below — Generic Defense Reaction. Cost 0.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "You may put a card from your hand on the bottom of your deck. If you do, draw a card."

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

type SinkBelowRed struct{}

func (SinkBelowRed) ID() card.ID              { return card.SinkBelowRed }
func (SinkBelowRed) Name() string             { return "Sink Below" }
func (SinkBelowRed) Cost(*card.TurnState) int { return 0 }
func (SinkBelowRed) Pitch() int               { return 1 }
func (SinkBelowRed) Attack() int              { return 0 }
func (SinkBelowRed) Defense() int             { return 4 }
func (SinkBelowRed) Types() card.TypeSet      { return defenseReactionTypes }
func (SinkBelowRed) GoAgain() bool            { return false }
func (SinkBelowRed) NotSilverAgeLegal()       {}

// not implemented: discard-to-cycle rider (hand cycling not modelled)
func (SinkBelowRed) NotImplemented() {}
func (SinkBelowRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}

type SinkBelowYellow struct{}

func (SinkBelowYellow) ID() card.ID              { return card.SinkBelowYellow }
func (SinkBelowYellow) Name() string             { return "Sink Below" }
func (SinkBelowYellow) Cost(*card.TurnState) int { return 0 }
func (SinkBelowYellow) Pitch() int               { return 2 }
func (SinkBelowYellow) Attack() int              { return 0 }
func (SinkBelowYellow) Defense() int             { return 3 }
func (SinkBelowYellow) Types() card.TypeSet      { return defenseReactionTypes }
func (SinkBelowYellow) GoAgain() bool            { return false }
func (SinkBelowYellow) NotSilverAgeLegal()       {}

// not implemented: discard-to-cycle rider (hand cycling not modelled)
func (SinkBelowYellow) NotImplemented() {}
func (SinkBelowYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}

type SinkBelowBlue struct{}

func (SinkBelowBlue) ID() card.ID              { return card.SinkBelowBlue }
func (SinkBelowBlue) Name() string             { return "Sink Below" }
func (SinkBelowBlue) Cost(*card.TurnState) int { return 0 }
func (SinkBelowBlue) Pitch() int               { return 3 }
func (SinkBelowBlue) Attack() int              { return 0 }
func (SinkBelowBlue) Defense() int             { return 2 }
func (SinkBelowBlue) Types() card.TypeSet      { return defenseReactionTypes }
func (SinkBelowBlue) GoAgain() bool            { return false }
func (SinkBelowBlue) NotSilverAgeLegal()       {}

// not implemented: discard-to-cycle rider (hand cycling not modelled)
func (SinkBelowBlue) NotImplemented() {}
func (SinkBelowBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}
