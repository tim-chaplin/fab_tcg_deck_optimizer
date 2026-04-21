// Tip-Off — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4, Blue 3. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Instant** - Discard this: **Mark** target opposing hero."
//
// Simplification: Instant discard activation isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var tipOffTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type TipOffRed struct{}

func (TipOffRed) ID() card.ID                 { return card.TipOffRed }
func (TipOffRed) Name() string                { return "Tip-Off (Red)" }
func (TipOffRed) Cost(*card.TurnState) int                   { return 1 }
func (TipOffRed) Pitch() int                  { return 1 }
func (TipOffRed) Attack() int                 { return 5 }
func (TipOffRed) Defense() int                { return 2 }
func (TipOffRed) Types() card.TypeSet         { return tipOffTypes }
func (TipOffRed) GoAgain() bool               { return false }
func (c TipOffRed) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }

type TipOffYellow struct{}

func (TipOffYellow) ID() card.ID                 { return card.TipOffYellow }
func (TipOffYellow) Name() string                { return "Tip-Off (Yellow)" }
func (TipOffYellow) Cost(*card.TurnState) int                   { return 1 }
func (TipOffYellow) Pitch() int                  { return 2 }
func (TipOffYellow) Attack() int                 { return 4 }
func (TipOffYellow) Defense() int                { return 2 }
func (TipOffYellow) Types() card.TypeSet         { return tipOffTypes }
func (TipOffYellow) GoAgain() bool               { return false }
func (c TipOffYellow) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }

type TipOffBlue struct{}

func (TipOffBlue) ID() card.ID                 { return card.TipOffBlue }
func (TipOffBlue) Name() string                { return "Tip-Off (Blue)" }
func (TipOffBlue) Cost(*card.TurnState) int                   { return 1 }
func (TipOffBlue) Pitch() int                  { return 3 }
func (TipOffBlue) Attack() int                 { return 3 }
func (TipOffBlue) Defense() int                { return 2 }
func (TipOffBlue) Types() card.TypeSet         { return tipOffTypes }
func (TipOffBlue) GoAgain() bool               { return false }
func (c TipOffBlue) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }
