// Rifting — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Rifting hits, you may play your next 'non-attack' action card this turn as though it
// were an instant."
//
// Simplification: On-hit instant-casting rider isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var riftingTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type RiftingRed struct{}

func (RiftingRed) ID() card.ID                 { return card.RiftingRed }
func (RiftingRed) Name() string                { return "Rifting (Red)" }
func (RiftingRed) Cost(*card.TurnState) int                   { return 2 }
func (RiftingRed) Pitch() int                  { return 1 }
func (RiftingRed) Attack() int                 { return 6 }
func (RiftingRed) Defense() int                { return 2 }
func (RiftingRed) Types() card.TypeSet         { return riftingTypes }
func (RiftingRed) GoAgain() bool               { return false }
func (c RiftingRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type RiftingYellow struct{}

func (RiftingYellow) ID() card.ID                 { return card.RiftingYellow }
func (RiftingYellow) Name() string                { return "Rifting (Yellow)" }
func (RiftingYellow) Cost(*card.TurnState) int                   { return 2 }
func (RiftingYellow) Pitch() int                  { return 2 }
func (RiftingYellow) Attack() int                 { return 5 }
func (RiftingYellow) Defense() int                { return 2 }
func (RiftingYellow) Types() card.TypeSet         { return riftingTypes }
func (RiftingYellow) GoAgain() bool               { return false }
func (c RiftingYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type RiftingBlue struct{}

func (RiftingBlue) ID() card.ID                 { return card.RiftingBlue }
func (RiftingBlue) Name() string                { return "Rifting (Blue)" }
func (RiftingBlue) Cost(*card.TurnState) int                   { return 2 }
func (RiftingBlue) Pitch() int                  { return 3 }
func (RiftingBlue) Attack() int                 { return 4 }
func (RiftingBlue) Defense() int                { return 2 }
func (RiftingBlue) Types() card.TypeSet         { return riftingTypes }
func (RiftingBlue) GoAgain() bool               { return false }
func (c RiftingBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
