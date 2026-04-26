// Nimble Strike — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Nimble Strike, you may banish a card named Nimblism from
// your graveyard. If you do, Nimble Strike gain +1{p} and **go again**."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var nimbleStrikeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type NimbleStrikeRed struct{}

func (NimbleStrikeRed) ID() card.ID                 { return card.NimbleStrikeRed }
func (NimbleStrikeRed) Name() string                { return "Nimble Strike (Red)" }
func (NimbleStrikeRed) Cost(*card.TurnState) int                   { return 1 }
func (NimbleStrikeRed) Pitch() int                  { return 1 }
func (NimbleStrikeRed) Attack() int                 { return 4 }
func (NimbleStrikeRed) Defense() int                { return 2 }
func (NimbleStrikeRed) Types() card.TypeSet         { return nimbleStrikeTypes }
func (NimbleStrikeRed) GoAgain() bool               { return false }
// not implemented: graveyard-banish-Nimblism additional cost and the +1{p}/go-again bonus rider
func (NimbleStrikeRed) NotImplemented()             {}
func (c NimbleStrikeRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type NimbleStrikeYellow struct{}

func (NimbleStrikeYellow) ID() card.ID                 { return card.NimbleStrikeYellow }
func (NimbleStrikeYellow) Name() string                { return "Nimble Strike (Yellow)" }
func (NimbleStrikeYellow) Cost(*card.TurnState) int                   { return 1 }
func (NimbleStrikeYellow) Pitch() int                  { return 2 }
func (NimbleStrikeYellow) Attack() int                 { return 3 }
func (NimbleStrikeYellow) Defense() int                { return 2 }
func (NimbleStrikeYellow) Types() card.TypeSet         { return nimbleStrikeTypes }
func (NimbleStrikeYellow) GoAgain() bool               { return false }
// not implemented: graveyard-banish-Nimblism additional cost and the +1{p}/go-again bonus rider
func (NimbleStrikeYellow) NotImplemented()             {}
func (c NimbleStrikeYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type NimbleStrikeBlue struct{}

func (NimbleStrikeBlue) ID() card.ID                 { return card.NimbleStrikeBlue }
func (NimbleStrikeBlue) Name() string                { return "Nimble Strike (Blue)" }
func (NimbleStrikeBlue) Cost(*card.TurnState) int                   { return 1 }
func (NimbleStrikeBlue) Pitch() int                  { return 3 }
func (NimbleStrikeBlue) Attack() int                 { return 2 }
func (NimbleStrikeBlue) Defense() int                { return 2 }
func (NimbleStrikeBlue) Types() card.TypeSet         { return nimbleStrikeTypes }
func (NimbleStrikeBlue) GoAgain() bool               { return false }
// not implemented: graveyard-banish-Nimblism additional cost and the +1{p}/go-again bonus rider
func (NimbleStrikeBlue) NotImplemented()             {}
func (c NimbleStrikeBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
