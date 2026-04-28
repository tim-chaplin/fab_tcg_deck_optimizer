// Walk the Plank — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a Pirate hero, {t} them or an ally they control."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var walkThePlankTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type WalkThePlankRed struct{}

func (WalkThePlankRed) ID() ids.CardID           { return ids.WalkThePlankRed }
func (WalkThePlankRed) Name() string             { return "Walk the Plank" }
func (WalkThePlankRed) Cost(*card.TurnState) int { return 3 }
func (WalkThePlankRed) Pitch() int               { return 1 }
func (WalkThePlankRed) Attack() int              { return 7 }
func (WalkThePlankRed) Defense() int             { return 2 }
func (WalkThePlankRed) Types() card.TypeSet      { return walkThePlankTypes }
func (WalkThePlankRed) GoAgain() bool            { return false }

// not implemented: pirate-target freeze rider
func (WalkThePlankRed) NotImplemented() {}
func (WalkThePlankRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type WalkThePlankYellow struct{}

func (WalkThePlankYellow) ID() ids.CardID           { return ids.WalkThePlankYellow }
func (WalkThePlankYellow) Name() string             { return "Walk the Plank" }
func (WalkThePlankYellow) Cost(*card.TurnState) int { return 3 }
func (WalkThePlankYellow) Pitch() int               { return 2 }
func (WalkThePlankYellow) Attack() int              { return 6 }
func (WalkThePlankYellow) Defense() int             { return 2 }
func (WalkThePlankYellow) Types() card.TypeSet      { return walkThePlankTypes }
func (WalkThePlankYellow) GoAgain() bool            { return false }

// not implemented: pirate-target freeze rider
func (WalkThePlankYellow) NotImplemented() {}
func (WalkThePlankYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type WalkThePlankBlue struct{}

func (WalkThePlankBlue) ID() ids.CardID           { return ids.WalkThePlankBlue }
func (WalkThePlankBlue) Name() string             { return "Walk the Plank" }
func (WalkThePlankBlue) Cost(*card.TurnState) int { return 3 }
func (WalkThePlankBlue) Pitch() int               { return 3 }
func (WalkThePlankBlue) Attack() int              { return 5 }
func (WalkThePlankBlue) Defense() int             { return 2 }
func (WalkThePlankBlue) Types() card.TypeSet      { return walkThePlankTypes }
func (WalkThePlankBlue) GoAgain() bool            { return false }

// not implemented: pirate-target freeze rider
func (WalkThePlankBlue) NotImplemented() {}
func (WalkThePlankBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
