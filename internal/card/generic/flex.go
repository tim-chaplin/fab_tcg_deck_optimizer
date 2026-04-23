// Flex — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When you attack or defend with Flex, you may pay {r}{r}. If you do, it gains +2{p}."
//
// Simplification: Pay-{r}{r}-for-+2{p} rider isn't modelled.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var flexTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type FlexRed struct{}

func (FlexRed) ID() card.ID                 { return card.FlexRed }
func (FlexRed) Name() string                { return "Flex (Red)" }
func (FlexRed) Cost(*card.TurnState) int                   { return 0 }
func (FlexRed) Pitch() int                  { return 1 }
func (FlexRed) Attack() int                 { return 4 }
func (FlexRed) Defense() int                { return 2 }
func (FlexRed) Types() card.TypeSet         { return flexTypes }
func (FlexRed) GoAgain() bool               { return false }
func (c FlexRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type FlexYellow struct{}

func (FlexYellow) ID() card.ID                 { return card.FlexYellow }
func (FlexYellow) Name() string                { return "Flex (Yellow)" }
func (FlexYellow) Cost(*card.TurnState) int                   { return 0 }
func (FlexYellow) Pitch() int                  { return 2 }
func (FlexYellow) Attack() int                 { return 3 }
func (FlexYellow) Defense() int                { return 2 }
func (FlexYellow) Types() card.TypeSet         { return flexTypes }
func (FlexYellow) GoAgain() bool               { return false }
func (c FlexYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type FlexBlue struct{}

func (FlexBlue) ID() card.ID                 { return card.FlexBlue }
func (FlexBlue) Name() string                { return "Flex (Blue)" }
func (FlexBlue) Cost(*card.TurnState) int                   { return 0 }
func (FlexBlue) Pitch() int                  { return 3 }
func (FlexBlue) Attack() int                 { return 2 }
func (FlexBlue) Defense() int                { return 2 }
func (FlexBlue) Types() card.TypeSet         { return flexTypes }
func (FlexBlue) GoAgain() bool               { return false }
func (c FlexBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
