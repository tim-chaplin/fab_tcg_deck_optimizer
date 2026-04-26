// Nimby — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, you may search your deck for a Nimblism, reveal it, put it into your
// hand, then shuffle."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var nimbyTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type NimbyRed struct{}

func (NimbyRed) ID() card.ID              { return card.NimbyRed }
func (NimbyRed) Name() string             { return "Nimby" }
func (NimbyRed) Cost(*card.TurnState) int { return 0 }
func (NimbyRed) Pitch() int               { return 1 }
func (NimbyRed) Attack() int              { return 3 }
func (NimbyRed) Defense() int             { return 2 }
func (NimbyRed) Types() card.TypeSet      { return nimbyTypes }
func (NimbyRed) GoAgain() bool            { return false }
func (NimbyRed) NotSilverAgeLegal()       {}

// not implemented: Nimblism deck-search tutor
func (NimbyRed) NotImplemented()                                {}
func (c NimbyRed) Play(s *card.TurnState, self *card.CardState) { s.ApplyAndLogEffectiveAttack(self) }

type NimbyYellow struct{}

func (NimbyYellow) ID() card.ID              { return card.NimbyYellow }
func (NimbyYellow) Name() string             { return "Nimby" }
func (NimbyYellow) Cost(*card.TurnState) int { return 0 }
func (NimbyYellow) Pitch() int               { return 2 }
func (NimbyYellow) Attack() int              { return 2 }
func (NimbyYellow) Defense() int             { return 2 }
func (NimbyYellow) Types() card.TypeSet      { return nimbyTypes }
func (NimbyYellow) GoAgain() bool            { return false }
func (NimbyYellow) NotSilverAgeLegal()       {}

// not implemented: Nimblism deck-search tutor
func (NimbyYellow) NotImplemented() {}
func (c NimbyYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type NimbyBlue struct{}

func (NimbyBlue) ID() card.ID              { return card.NimbyBlue }
func (NimbyBlue) Name() string             { return "Nimby" }
func (NimbyBlue) Cost(*card.TurnState) int { return 0 }
func (NimbyBlue) Pitch() int               { return 3 }
func (NimbyBlue) Attack() int              { return 1 }
func (NimbyBlue) Defense() int             { return 2 }
func (NimbyBlue) Types() card.TypeSet      { return nimbyTypes }
func (NimbyBlue) GoAgain() bool            { return false }
func (NimbyBlue) NotSilverAgeLegal()       {}

// not implemented: Nimblism deck-search tutor
func (NimbyBlue) NotImplemented()                                {}
func (c NimbyBlue) Play(s *card.TurnState, self *card.CardState) { s.ApplyAndLogEffectiveAttack(self) }
