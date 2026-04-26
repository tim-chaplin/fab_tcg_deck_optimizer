// Hocus Pocus — Runeblade Action - Attack. Cost 0, Defense 3.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "When this attacks, create a Runechant token."

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var hocusPocusTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type HocusPocusRed struct{}

func (HocusPocusRed) ID() card.ID                 { return card.HocusPocusRed }
func (HocusPocusRed) Name() string               { return "Hocus Pocus" }
func (HocusPocusRed) Cost(*card.TurnState) int                  { return 0 }
func (HocusPocusRed) Pitch() int                 { return 1 }
func (HocusPocusRed) Attack() int                { return 3 }
func (HocusPocusRed) Defense() int               { return 3 }
func (HocusPocusRed) Types() card.TypeSet     { return hocusPocusTypes }
func (HocusPocusRed) GoAgain() bool              { return false }
func (c HocusPocusRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() + s.CreateRunechant() }

type HocusPocusYellow struct{}

func (HocusPocusYellow) ID() card.ID                 { return card.HocusPocusYellow }
func (HocusPocusYellow) Name() string               { return "Hocus Pocus" }
func (HocusPocusYellow) Cost(*card.TurnState) int                  { return 0 }
func (HocusPocusYellow) Pitch() int                 { return 2 }
func (HocusPocusYellow) Attack() int                { return 2 }
func (HocusPocusYellow) Defense() int               { return 3 }
func (HocusPocusYellow) Types() card.TypeSet     { return hocusPocusTypes }
func (HocusPocusYellow) GoAgain() bool              { return false }
func (c HocusPocusYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() + s.CreateRunechant() }

type HocusPocusBlue struct{}

func (HocusPocusBlue) ID() card.ID                 { return card.HocusPocusBlue }
func (HocusPocusBlue) Name() string               { return "Hocus Pocus" }
func (HocusPocusBlue) Cost(*card.TurnState) int                  { return 0 }
func (HocusPocusBlue) Pitch() int                 { return 3 }
func (HocusPocusBlue) Attack() int                { return 1 }
func (HocusPocusBlue) Defense() int               { return 3 }
func (HocusPocusBlue) Types() card.TypeSet     { return hocusPocusTypes }
func (HocusPocusBlue) GoAgain() bool              { return false }
func (c HocusPocusBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() + s.CreateRunechant() }
