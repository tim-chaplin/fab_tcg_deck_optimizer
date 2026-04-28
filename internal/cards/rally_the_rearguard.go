// Rally the Rearguard — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Once per Turn Instant** - Discard a card: This gets +3{d}. Activate this ability only
// while this is defending."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var rallyTheRearguardTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type RallyTheRearguardRed struct{}

func (RallyTheRearguardRed) ID() ids.CardID           { return ids.RallyTheRearguardRed }
func (RallyTheRearguardRed) Name() string             { return "Rally the Rearguard" }
func (RallyTheRearguardRed) Cost(*card.TurnState) int { return 2 }
func (RallyTheRearguardRed) Pitch() int               { return 1 }
func (RallyTheRearguardRed) Attack() int              { return 6 }
func (RallyTheRearguardRed) Defense() int             { return 2 }
func (RallyTheRearguardRed) Types() card.TypeSet      { return rallyTheRearguardTypes }
func (RallyTheRearguardRed) GoAgain() bool            { return false }

// not implemented: defense-time instant activated ability
func (RallyTheRearguardRed) NotImplemented() {}
func (c RallyTheRearguardRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type RallyTheRearguardYellow struct{}

func (RallyTheRearguardYellow) ID() ids.CardID           { return ids.RallyTheRearguardYellow }
func (RallyTheRearguardYellow) Name() string             { return "Rally the Rearguard" }
func (RallyTheRearguardYellow) Cost(*card.TurnState) int { return 2 }
func (RallyTheRearguardYellow) Pitch() int               { return 2 }
func (RallyTheRearguardYellow) Attack() int              { return 5 }
func (RallyTheRearguardYellow) Defense() int             { return 2 }
func (RallyTheRearguardYellow) Types() card.TypeSet      { return rallyTheRearguardTypes }
func (RallyTheRearguardYellow) GoAgain() bool            { return false }

// not implemented: defense-time instant activated ability
func (RallyTheRearguardYellow) NotImplemented() {}
func (c RallyTheRearguardYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type RallyTheRearguardBlue struct{}

func (RallyTheRearguardBlue) ID() ids.CardID           { return ids.RallyTheRearguardBlue }
func (RallyTheRearguardBlue) Name() string             { return "Rally the Rearguard" }
func (RallyTheRearguardBlue) Cost(*card.TurnState) int { return 2 }
func (RallyTheRearguardBlue) Pitch() int               { return 3 }
func (RallyTheRearguardBlue) Attack() int              { return 4 }
func (RallyTheRearguardBlue) Defense() int             { return 2 }
func (RallyTheRearguardBlue) Types() card.TypeSet      { return rallyTheRearguardTypes }
func (RallyTheRearguardBlue) GoAgain() bool            { return false }

// not implemented: defense-time instant activated ability
func (RallyTheRearguardBlue) NotImplemented() {}
func (c RallyTheRearguardBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
