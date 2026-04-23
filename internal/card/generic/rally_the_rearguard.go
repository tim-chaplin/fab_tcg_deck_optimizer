// Rally the Rearguard — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Once per Turn Instant** - Discard a card: This gets +3{d}. Activate this ability only
// while this is defending."
//
// Simplification: Defense-time instant activated ability isn't modelled.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var rallyTheRearguardTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type RallyTheRearguardRed struct{}

func (RallyTheRearguardRed) ID() card.ID                 { return card.RallyTheRearguardRed }
func (RallyTheRearguardRed) Name() string                { return "Rally the Rearguard (Red)" }
func (RallyTheRearguardRed) Cost(*card.TurnState) int                   { return 2 }
func (RallyTheRearguardRed) Pitch() int                  { return 1 }
func (RallyTheRearguardRed) Attack() int                 { return 6 }
func (RallyTheRearguardRed) Defense() int                { return 2 }
func (RallyTheRearguardRed) Types() card.TypeSet         { return rallyTheRearguardTypes }
func (RallyTheRearguardRed) GoAgain() bool               { return false }
func (c RallyTheRearguardRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type RallyTheRearguardYellow struct{}

func (RallyTheRearguardYellow) ID() card.ID                 { return card.RallyTheRearguardYellow }
func (RallyTheRearguardYellow) Name() string                { return "Rally the Rearguard (Yellow)" }
func (RallyTheRearguardYellow) Cost(*card.TurnState) int                   { return 2 }
func (RallyTheRearguardYellow) Pitch() int                  { return 2 }
func (RallyTheRearguardYellow) Attack() int                 { return 5 }
func (RallyTheRearguardYellow) Defense() int                { return 2 }
func (RallyTheRearguardYellow) Types() card.TypeSet         { return rallyTheRearguardTypes }
func (RallyTheRearguardYellow) GoAgain() bool               { return false }
func (c RallyTheRearguardYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type RallyTheRearguardBlue struct{}

func (RallyTheRearguardBlue) ID() card.ID                 { return card.RallyTheRearguardBlue }
func (RallyTheRearguardBlue) Name() string                { return "Rally the Rearguard (Blue)" }
func (RallyTheRearguardBlue) Cost(*card.TurnState) int                   { return 2 }
func (RallyTheRearguardBlue) Pitch() int                  { return 3 }
func (RallyTheRearguardBlue) Attack() int                 { return 4 }
func (RallyTheRearguardBlue) Defense() int                { return 2 }
func (RallyTheRearguardBlue) Types() card.TypeSet         { return rallyTheRearguardTypes }
func (RallyTheRearguardBlue) GoAgain() bool               { return false }
func (c RallyTheRearguardBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
