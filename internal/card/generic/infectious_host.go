// Infectious Host — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks a hero, if you control a Frailty token, create a Frailty token under
// their control, then repeat for Inertia and Bloodrot Pox."
//
// Simplification: Frailty/Inertia/Bloodrot Pox token creation isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var infectiousHostTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type InfectiousHostRed struct{}

func (InfectiousHostRed) ID() card.ID                 { return card.InfectiousHostRed }
func (InfectiousHostRed) Name() string                { return "Infectious Host (Red)" }
func (InfectiousHostRed) Cost(*card.TurnState) int                   { return 0 }
func (InfectiousHostRed) Pitch() int                  { return 1 }
func (InfectiousHostRed) Attack() int                 { return 4 }
func (InfectiousHostRed) Defense() int                { return 2 }
func (InfectiousHostRed) Types() card.TypeSet         { return infectiousHostTypes }
func (InfectiousHostRed) GoAgain() bool               { return false }
func (c InfectiousHostRed) Play(s *card.TurnState) int { return c.Attack() }

type InfectiousHostYellow struct{}

func (InfectiousHostYellow) ID() card.ID                 { return card.InfectiousHostYellow }
func (InfectiousHostYellow) Name() string                { return "Infectious Host (Yellow)" }
func (InfectiousHostYellow) Cost(*card.TurnState) int                   { return 0 }
func (InfectiousHostYellow) Pitch() int                  { return 2 }
func (InfectiousHostYellow) Attack() int                 { return 3 }
func (InfectiousHostYellow) Defense() int                { return 2 }
func (InfectiousHostYellow) Types() card.TypeSet         { return infectiousHostTypes }
func (InfectiousHostYellow) GoAgain() bool               { return false }
func (c InfectiousHostYellow) Play(s *card.TurnState) int { return c.Attack() }

type InfectiousHostBlue struct{}

func (InfectiousHostBlue) ID() card.ID                 { return card.InfectiousHostBlue }
func (InfectiousHostBlue) Name() string                { return "Infectious Host (Blue)" }
func (InfectiousHostBlue) Cost(*card.TurnState) int                   { return 0 }
func (InfectiousHostBlue) Pitch() int                  { return 3 }
func (InfectiousHostBlue) Attack() int                 { return 2 }
func (InfectiousHostBlue) Defense() int                { return 2 }
func (InfectiousHostBlue) Types() card.TypeSet         { return infectiousHostTypes }
func (InfectiousHostBlue) GoAgain() bool               { return false }
func (c InfectiousHostBlue) Play(s *card.TurnState) int { return c.Attack() }
