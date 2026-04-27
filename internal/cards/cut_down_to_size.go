// Cut Down to Size — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, if they have 4 or more cards in hand, they discard a card."

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var cutDownToSizeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type CutDownToSizeRed struct{}

func (CutDownToSizeRed) ID() card.ID              { return card.CutDownToSizeRed }
func (CutDownToSizeRed) Name() string             { return "Cut Down to Size" }
func (CutDownToSizeRed) Cost(*card.TurnState) int { return 2 }
func (CutDownToSizeRed) Pitch() int               { return 1 }
func (CutDownToSizeRed) Attack() int              { return 6 }
func (CutDownToSizeRed) Defense() int             { return 2 }
func (CutDownToSizeRed) Types() card.TypeSet      { return cutDownToSizeTypes }
func (CutDownToSizeRed) GoAgain() bool            { return false }

// not implemented: on-hit opponent discard (conditional on hand size)
func (CutDownToSizeRed) NotImplemented() {}
func (CutDownToSizeRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type CutDownToSizeYellow struct{}

func (CutDownToSizeYellow) ID() card.ID              { return card.CutDownToSizeYellow }
func (CutDownToSizeYellow) Name() string             { return "Cut Down to Size" }
func (CutDownToSizeYellow) Cost(*card.TurnState) int { return 2 }
func (CutDownToSizeYellow) Pitch() int               { return 2 }
func (CutDownToSizeYellow) Attack() int              { return 5 }
func (CutDownToSizeYellow) Defense() int             { return 2 }
func (CutDownToSizeYellow) Types() card.TypeSet      { return cutDownToSizeTypes }
func (CutDownToSizeYellow) GoAgain() bool            { return false }

// not implemented: on-hit opponent discard (conditional on hand size)
func (CutDownToSizeYellow) NotImplemented() {}
func (CutDownToSizeYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type CutDownToSizeBlue struct{}

func (CutDownToSizeBlue) ID() card.ID              { return card.CutDownToSizeBlue }
func (CutDownToSizeBlue) Name() string             { return "Cut Down to Size" }
func (CutDownToSizeBlue) Cost(*card.TurnState) int { return 2 }
func (CutDownToSizeBlue) Pitch() int               { return 3 }
func (CutDownToSizeBlue) Attack() int              { return 4 }
func (CutDownToSizeBlue) Defense() int             { return 2 }
func (CutDownToSizeBlue) Types() card.TypeSet      { return cutDownToSizeTypes }
func (CutDownToSizeBlue) GoAgain() bool            { return false }

// not implemented: on-hit opponent discard (conditional on hand size)
func (CutDownToSizeBlue) NotImplemented() {}
func (CutDownToSizeBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
