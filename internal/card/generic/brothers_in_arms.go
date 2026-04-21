// Brothers in Arms — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this defends, you may pay {r}. If you do, it gets +2{d}."
//
// Simplification: Pay-to-buff-defence rider isn't modelled (defence-side costs aren't solved).
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var brothersInArmsTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BrothersInArmsRed struct{}

func (BrothersInArmsRed) ID() card.ID                 { return card.BrothersInArmsRed }
func (BrothersInArmsRed) Name() string                { return "Brothers in Arms (Red)" }
func (BrothersInArmsRed) Cost(*card.TurnState) int                   { return 2 }
func (BrothersInArmsRed) Pitch() int                  { return 1 }
func (BrothersInArmsRed) Attack() int                 { return 6 }
func (BrothersInArmsRed) Defense() int                { return 2 }
func (BrothersInArmsRed) Types() card.TypeSet         { return brothersInArmsTypes }
func (BrothersInArmsRed) GoAgain() bool               { return false }
func (c BrothersInArmsRed) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }

type BrothersInArmsYellow struct{}

func (BrothersInArmsYellow) ID() card.ID                 { return card.BrothersInArmsYellow }
func (BrothersInArmsYellow) Name() string                { return "Brothers in Arms (Yellow)" }
func (BrothersInArmsYellow) Cost(*card.TurnState) int                   { return 2 }
func (BrothersInArmsYellow) Pitch() int                  { return 2 }
func (BrothersInArmsYellow) Attack() int                 { return 5 }
func (BrothersInArmsYellow) Defense() int                { return 2 }
func (BrothersInArmsYellow) Types() card.TypeSet         { return brothersInArmsTypes }
func (BrothersInArmsYellow) GoAgain() bool               { return false }
func (c BrothersInArmsYellow) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }

type BrothersInArmsBlue struct{}

func (BrothersInArmsBlue) ID() card.ID                 { return card.BrothersInArmsBlue }
func (BrothersInArmsBlue) Name() string                { return "Brothers in Arms (Blue)" }
func (BrothersInArmsBlue) Cost(*card.TurnState) int                   { return 2 }
func (BrothersInArmsBlue) Pitch() int                  { return 3 }
func (BrothersInArmsBlue) Attack() int                 { return 4 }
func (BrothersInArmsBlue) Defense() int                { return 2 }
func (BrothersInArmsBlue) Types() card.TypeSet         { return brothersInArmsTypes }
func (BrothersInArmsBlue) GoAgain() bool               { return false }
func (c BrothersInArmsBlue) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }
