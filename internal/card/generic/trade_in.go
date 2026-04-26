// Trade In — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, you may discard a card. If you do, draw a card. If this was played from
// arsenal, it gains **go again**."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var tradeInTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type TradeInRed struct{}

func (TradeInRed) ID() card.ID                 { return card.TradeInRed }
func (TradeInRed) Name() string                { return "Trade In" }
func (TradeInRed) Cost(*card.TurnState) int                   { return 0 }
func (TradeInRed) Pitch() int                  { return 1 }
func (TradeInRed) Attack() int                 { return 3 }
func (TradeInRed) Defense() int                { return 2 }
func (TradeInRed) Types() card.TypeSet         { return tradeInTypes }
func (TradeInRed) GoAgain() bool               { return false }
// not implemented: discard-to-draw rider, arsenal-conditional go again
func (TradeInRed) NotImplemented()             {}
func (c TradeInRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type TradeInYellow struct{}

func (TradeInYellow) ID() card.ID                 { return card.TradeInYellow }
func (TradeInYellow) Name() string                { return "Trade In" }
func (TradeInYellow) Cost(*card.TurnState) int                   { return 0 }
func (TradeInYellow) Pitch() int                  { return 2 }
func (TradeInYellow) Attack() int                 { return 2 }
func (TradeInYellow) Defense() int                { return 2 }
func (TradeInYellow) Types() card.TypeSet         { return tradeInTypes }
func (TradeInYellow) GoAgain() bool               { return false }
// not implemented: discard-to-draw rider, arsenal-conditional go again
func (TradeInYellow) NotImplemented()             {}
func (c TradeInYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type TradeInBlue struct{}

func (TradeInBlue) ID() card.ID                 { return card.TradeInBlue }
func (TradeInBlue) Name() string                { return "Trade In" }
func (TradeInBlue) Cost(*card.TurnState) int                   { return 0 }
func (TradeInBlue) Pitch() int                  { return 3 }
func (TradeInBlue) Attack() int                 { return 1 }
func (TradeInBlue) Defense() int                { return 2 }
func (TradeInBlue) Types() card.TypeSet         { return tradeInTypes }
func (TradeInBlue) GoAgain() bool               { return false }
// not implemented: discard-to-draw rider, arsenal-conditional go again
func (TradeInBlue) NotImplemented()             {}
func (c TradeInBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
