// Read the Runes — Runeblade Action. Cost 0, Defense 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Create N Runechant tokens." (Red N=3, Yellow N=2, Blue N=1.)
//
// Play returns N and sets AuraCreated so later cards this turn see the Runechants.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var readTheRunesTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction)

type ReadTheRunesRed struct{}

func (ReadTheRunesRed) ID() card.ID                 { return card.ReadTheRunesRed }
func (ReadTheRunesRed) Name() string               { return "Read the Runes (Red)" }
func (ReadTheRunesRed) Cost(*card.TurnState) int                  { return 0 }
func (ReadTheRunesRed) Pitch() int                 { return 1 }
func (ReadTheRunesRed) Attack() int                { return 0 }
func (ReadTheRunesRed) Defense() int               { return 2 }
func (ReadTheRunesRed) Types() card.TypeSet        { return readTheRunesTypes }
func (ReadTheRunesRed) GoAgain() bool              { return false }
func (ReadTheRunesRed) Play(s *card.TurnState, _ *card.CardState) int { return s.CreateRunechants(3) }

type ReadTheRunesYellow struct{}

func (ReadTheRunesYellow) ID() card.ID                 { return card.ReadTheRunesYellow }
func (ReadTheRunesYellow) Name() string               { return "Read the Runes (Yellow)" }
func (ReadTheRunesYellow) Cost(*card.TurnState) int                  { return 0 }
func (ReadTheRunesYellow) Pitch() int                 { return 2 }
func (ReadTheRunesYellow) Attack() int                { return 0 }
func (ReadTheRunesYellow) Defense() int               { return 2 }
func (ReadTheRunesYellow) Types() card.TypeSet        { return readTheRunesTypes }
func (ReadTheRunesYellow) GoAgain() bool              { return false }
func (ReadTheRunesYellow) Play(s *card.TurnState, _ *card.CardState) int { return s.CreateRunechants(2) }

type ReadTheRunesBlue struct{}

func (ReadTheRunesBlue) ID() card.ID                 { return card.ReadTheRunesBlue }
func (ReadTheRunesBlue) Name() string               { return "Read the Runes (Blue)" }
func (ReadTheRunesBlue) Cost(*card.TurnState) int                  { return 0 }
func (ReadTheRunesBlue) Pitch() int                 { return 3 }
func (ReadTheRunesBlue) Attack() int                { return 0 }
func (ReadTheRunesBlue) Defense() int               { return 2 }
func (ReadTheRunesBlue) Types() card.TypeSet        { return readTheRunesTypes }
func (ReadTheRunesBlue) GoAgain() bool              { return false }
func (ReadTheRunesBlue) Play(s *card.TurnState, _ *card.CardState) int { return s.CreateRunechants(1) }
