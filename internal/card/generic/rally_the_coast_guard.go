// Rally the Coast Guard — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Once per Turn Instant** - Discard a card: This gets +3{d}. Activate this only while this
// card is defending."
//
// Simplification: Defense-time instant activated ability isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var rallyTheCoastGuardTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type RallyTheCoastGuardRed struct{}

func (RallyTheCoastGuardRed) ID() card.ID                 { return card.RallyTheCoastGuardRed }
func (RallyTheCoastGuardRed) Name() string                { return "Rally the Coast Guard (Red)" }
func (RallyTheCoastGuardRed) Cost(*card.TurnState) int                   { return 3 }
func (RallyTheCoastGuardRed) Pitch() int                  { return 1 }
func (RallyTheCoastGuardRed) Attack() int                 { return 7 }
func (RallyTheCoastGuardRed) Defense() int                { return 2 }
func (RallyTheCoastGuardRed) Types() card.TypeSet         { return rallyTheCoastGuardTypes }
func (RallyTheCoastGuardRed) GoAgain() bool               { return false }
func (c RallyTheCoastGuardRed) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }

type RallyTheCoastGuardYellow struct{}

func (RallyTheCoastGuardYellow) ID() card.ID                 { return card.RallyTheCoastGuardYellow }
func (RallyTheCoastGuardYellow) Name() string                { return "Rally the Coast Guard (Yellow)" }
func (RallyTheCoastGuardYellow) Cost(*card.TurnState) int                   { return 3 }
func (RallyTheCoastGuardYellow) Pitch() int                  { return 2 }
func (RallyTheCoastGuardYellow) Attack() int                 { return 6 }
func (RallyTheCoastGuardYellow) Defense() int                { return 2 }
func (RallyTheCoastGuardYellow) Types() card.TypeSet         { return rallyTheCoastGuardTypes }
func (RallyTheCoastGuardYellow) GoAgain() bool               { return false }
func (c RallyTheCoastGuardYellow) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }

type RallyTheCoastGuardBlue struct{}

func (RallyTheCoastGuardBlue) ID() card.ID                 { return card.RallyTheCoastGuardBlue }
func (RallyTheCoastGuardBlue) Name() string                { return "Rally the Coast Guard (Blue)" }
func (RallyTheCoastGuardBlue) Cost(*card.TurnState) int                   { return 3 }
func (RallyTheCoastGuardBlue) Pitch() int                  { return 3 }
func (RallyTheCoastGuardBlue) Attack() int                 { return 5 }
func (RallyTheCoastGuardBlue) Defense() int                { return 2 }
func (RallyTheCoastGuardBlue) Types() card.TypeSet         { return rallyTheCoastGuardTypes }
func (RallyTheCoastGuardBlue) GoAgain() bool               { return false }
func (c RallyTheCoastGuardBlue) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }
