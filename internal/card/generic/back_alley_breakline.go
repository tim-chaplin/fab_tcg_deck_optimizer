// Back Alley Breakline — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If an activated ability or action card effect puts Back Alley Breakline face up into a
// zone from your deck, gain 1 action point."
//
// Simplification: Face-up-from-deck action-point grant isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var backAlleyBreaklineTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BackAlleyBreaklineRed struct{}

func (BackAlleyBreaklineRed) ID() card.ID                 { return card.BackAlleyBreaklineRed }
func (BackAlleyBreaklineRed) Name() string                { return "Back Alley Breakline (Red)" }
func (BackAlleyBreaklineRed) Cost(*card.TurnState) int                   { return 1 }
func (BackAlleyBreaklineRed) Pitch() int                  { return 1 }
func (BackAlleyBreaklineRed) Attack() int                 { return 5 }
func (BackAlleyBreaklineRed) Defense() int                { return 2 }
func (BackAlleyBreaklineRed) Types() card.TypeSet         { return backAlleyBreaklineTypes }
func (BackAlleyBreaklineRed) GoAgain() bool               { return false }
func (c BackAlleyBreaklineRed) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }

type BackAlleyBreaklineYellow struct{}

func (BackAlleyBreaklineYellow) ID() card.ID                 { return card.BackAlleyBreaklineYellow }
func (BackAlleyBreaklineYellow) Name() string                { return "Back Alley Breakline (Yellow)" }
func (BackAlleyBreaklineYellow) Cost(*card.TurnState) int                   { return 1 }
func (BackAlleyBreaklineYellow) Pitch() int                  { return 2 }
func (BackAlleyBreaklineYellow) Attack() int                 { return 4 }
func (BackAlleyBreaklineYellow) Defense() int                { return 2 }
func (BackAlleyBreaklineYellow) Types() card.TypeSet         { return backAlleyBreaklineTypes }
func (BackAlleyBreaklineYellow) GoAgain() bool               { return false }
func (c BackAlleyBreaklineYellow) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }

type BackAlleyBreaklineBlue struct{}

func (BackAlleyBreaklineBlue) ID() card.ID                 { return card.BackAlleyBreaklineBlue }
func (BackAlleyBreaklineBlue) Name() string                { return "Back Alley Breakline (Blue)" }
func (BackAlleyBreaklineBlue) Cost(*card.TurnState) int                   { return 1 }
func (BackAlleyBreaklineBlue) Pitch() int                  { return 3 }
func (BackAlleyBreaklineBlue) Attack() int                 { return 3 }
func (BackAlleyBreaklineBlue) Defense() int                { return 2 }
func (BackAlleyBreaklineBlue) Types() card.TypeSet         { return backAlleyBreaklineTypes }
func (BackAlleyBreaklineBlue) GoAgain() bool               { return false }
func (c BackAlleyBreaklineBlue) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }
