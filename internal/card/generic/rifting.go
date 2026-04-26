// Rifting — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Rifting hits, you may play your next 'non-attack' action card this turn as though it
// were an instant."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var riftingTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type RiftingRed struct{}

func (RiftingRed) ID() card.ID              { return card.RiftingRed }
func (RiftingRed) Name() string             { return "Rifting" }
func (RiftingRed) Cost(*card.TurnState) int { return 2 }
func (RiftingRed) Pitch() int               { return 1 }
func (RiftingRed) Attack() int              { return 6 }
func (RiftingRed) Defense() int             { return 2 }
func (RiftingRed) Types() card.TypeSet      { return riftingTypes }
func (RiftingRed) GoAgain() bool            { return false }

// not implemented: on-hit instant-casting grant
func (RiftingRed) NotImplemented()                                {}
func (c RiftingRed) Play(s *card.TurnState, self *card.CardState) { s.ApplyAndLogEffectiveAttack(self) }

type RiftingYellow struct{}

func (RiftingYellow) ID() card.ID              { return card.RiftingYellow }
func (RiftingYellow) Name() string             { return "Rifting" }
func (RiftingYellow) Cost(*card.TurnState) int { return 2 }
func (RiftingYellow) Pitch() int               { return 2 }
func (RiftingYellow) Attack() int              { return 5 }
func (RiftingYellow) Defense() int             { return 2 }
func (RiftingYellow) Types() card.TypeSet      { return riftingTypes }
func (RiftingYellow) GoAgain() bool            { return false }

// not implemented: on-hit instant-casting grant
func (RiftingYellow) NotImplemented() {}
func (c RiftingYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type RiftingBlue struct{}

func (RiftingBlue) ID() card.ID              { return card.RiftingBlue }
func (RiftingBlue) Name() string             { return "Rifting" }
func (RiftingBlue) Cost(*card.TurnState) int { return 2 }
func (RiftingBlue) Pitch() int               { return 3 }
func (RiftingBlue) Attack() int              { return 4 }
func (RiftingBlue) Defense() int             { return 2 }
func (RiftingBlue) Types() card.TypeSet      { return riftingTypes }
func (RiftingBlue) GoAgain() bool            { return false }

// not implemented: on-hit instant-casting grant
func (RiftingBlue) NotImplemented() {}
func (c RiftingBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
