// Looking for a Scrap — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Looking for a Scrap, you may banish a card with 1{p} from
// your graveyard. When you do, this gains +1{p} and **go again**."
//
// Simplification: Graveyard-banish additional cost and bonus rider aren't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var lookingForAScrapTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type LookingForAScrapRed struct{}

func (LookingForAScrapRed) ID() card.ID                 { return card.LookingForAScrapRed }
func (LookingForAScrapRed) Name() string                { return "Looking for a Scrap (Red)" }
func (LookingForAScrapRed) Cost(*card.TurnState) int                   { return 1 }
func (LookingForAScrapRed) Pitch() int                  { return 1 }
func (LookingForAScrapRed) Attack() int                 { return 4 }
func (LookingForAScrapRed) Defense() int                { return 2 }
func (LookingForAScrapRed) Types() card.TypeSet         { return lookingForAScrapTypes }
func (LookingForAScrapRed) GoAgain() bool               { return false }
func (c LookingForAScrapRed) Play(s *card.TurnState) int { return c.Attack() }

type LookingForAScrapYellow struct{}

func (LookingForAScrapYellow) ID() card.ID                 { return card.LookingForAScrapYellow }
func (LookingForAScrapYellow) Name() string                { return "Looking for a Scrap (Yellow)" }
func (LookingForAScrapYellow) Cost(*card.TurnState) int                   { return 1 }
func (LookingForAScrapYellow) Pitch() int                  { return 2 }
func (LookingForAScrapYellow) Attack() int                 { return 3 }
func (LookingForAScrapYellow) Defense() int                { return 2 }
func (LookingForAScrapYellow) Types() card.TypeSet         { return lookingForAScrapTypes }
func (LookingForAScrapYellow) GoAgain() bool               { return false }
func (c LookingForAScrapYellow) Play(s *card.TurnState) int { return c.Attack() }

type LookingForAScrapBlue struct{}

func (LookingForAScrapBlue) ID() card.ID                 { return card.LookingForAScrapBlue }
func (LookingForAScrapBlue) Name() string                { return "Looking for a Scrap (Blue)" }
func (LookingForAScrapBlue) Cost(*card.TurnState) int                   { return 1 }
func (LookingForAScrapBlue) Pitch() int                  { return 3 }
func (LookingForAScrapBlue) Attack() int                 { return 2 }
func (LookingForAScrapBlue) Defense() int                { return 2 }
func (LookingForAScrapBlue) Types() card.TypeSet         { return lookingForAScrapTypes }
func (LookingForAScrapBlue) GoAgain() bool               { return false }
func (c LookingForAScrapBlue) Play(s *card.TurnState) int { return c.Attack() }
