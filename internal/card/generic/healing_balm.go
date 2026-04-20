// Healing Balm — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 2. Printed health-gain: Red 3{h}, Yellow 2{h}, Blue 1{h}.
//
// Text: "Gain N{h}." (N is the printed variant value above.)
//
// Modelling: health is valued 1-to-1 with damage, so Play returns +N damage-equivalent per
// variant.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var healingBalmTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type HealingBalmRed struct{}

func (HealingBalmRed) ID() card.ID                 { return card.HealingBalmRed }
func (HealingBalmRed) Name() string                { return "Healing Balm (Red)" }
func (HealingBalmRed) Cost(*card.TurnState) int                   { return 0 }
func (HealingBalmRed) Pitch() int                  { return 1 }
func (HealingBalmRed) Attack() int                 { return 0 }
func (HealingBalmRed) Defense() int                { return 2 }
func (HealingBalmRed) Types() card.TypeSet         { return healingBalmTypes }
func (HealingBalmRed) GoAgain() bool               { return false }
func (HealingBalmRed) Play(s *card.TurnState) int { return 3 }

type HealingBalmYellow struct{}

func (HealingBalmYellow) ID() card.ID                 { return card.HealingBalmYellow }
func (HealingBalmYellow) Name() string                { return "Healing Balm (Yellow)" }
func (HealingBalmYellow) Cost(*card.TurnState) int                   { return 0 }
func (HealingBalmYellow) Pitch() int                  { return 2 }
func (HealingBalmYellow) Attack() int                 { return 0 }
func (HealingBalmYellow) Defense() int                { return 2 }
func (HealingBalmYellow) Types() card.TypeSet         { return healingBalmTypes }
func (HealingBalmYellow) GoAgain() bool               { return false }
func (HealingBalmYellow) Play(s *card.TurnState) int { return 2 }

type HealingBalmBlue struct{}

func (HealingBalmBlue) ID() card.ID                 { return card.HealingBalmBlue }
func (HealingBalmBlue) Name() string                { return "Healing Balm (Blue)" }
func (HealingBalmBlue) Cost(*card.TurnState) int                   { return 0 }
func (HealingBalmBlue) Pitch() int                  { return 3 }
func (HealingBalmBlue) Attack() int                 { return 0 }
func (HealingBalmBlue) Defense() int                { return 2 }
func (HealingBalmBlue) Types() card.TypeSet         { return healingBalmTypes }
func (HealingBalmBlue) GoAgain() bool               { return false }
func (HealingBalmBlue) Play(s *card.TurnState) int { return 1 }
