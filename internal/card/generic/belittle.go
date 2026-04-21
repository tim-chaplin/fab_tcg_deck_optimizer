// Belittle — Generic Action - Attack. Cost 1. Printed power: Red 3, Yellow 2, Blue 1. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Belittle, you may reveal an attack action card with 3 or
// less base {p} from your hand. If you do, search your deck for a card named Minnowism, reveal it,
// put it into your hand, then shuffle your deck. **Go again**"
//
// Simplification: Additional-cost reveal and deck search for Minnowism aren't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var belittleTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BelittleRed struct{}

func (BelittleRed) ID() card.ID                 { return card.BelittleRed }
func (BelittleRed) Name() string                { return "Belittle (Red)" }
func (BelittleRed) Cost(*card.TurnState) int                   { return 1 }
func (BelittleRed) Pitch() int                  { return 1 }
func (BelittleRed) Attack() int                 { return 3 }
func (BelittleRed) Defense() int                { return 2 }
func (BelittleRed) Types() card.TypeSet         { return belittleTypes }
func (BelittleRed) GoAgain() bool               { return true }
func (BelittleRed) NotSilverAgeLegal()           {}
func (c BelittleRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type BelittleYellow struct{}

func (BelittleYellow) ID() card.ID                 { return card.BelittleYellow }
func (BelittleYellow) Name() string                { return "Belittle (Yellow)" }
func (BelittleYellow) Cost(*card.TurnState) int                   { return 1 }
func (BelittleYellow) Pitch() int                  { return 2 }
func (BelittleYellow) Attack() int                 { return 2 }
func (BelittleYellow) Defense() int                { return 2 }
func (BelittleYellow) Types() card.TypeSet         { return belittleTypes }
func (BelittleYellow) GoAgain() bool               { return true }
func (BelittleYellow) NotSilverAgeLegal()           {}
func (c BelittleYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type BelittleBlue struct{}

func (BelittleBlue) ID() card.ID                 { return card.BelittleBlue }
func (BelittleBlue) Name() string                { return "Belittle (Blue)" }
func (BelittleBlue) Cost(*card.TurnState) int                   { return 1 }
func (BelittleBlue) Pitch() int                  { return 3 }
func (BelittleBlue) Attack() int                 { return 1 }
func (BelittleBlue) Defense() int                { return 2 }
func (BelittleBlue) Types() card.TypeSet         { return belittleTypes }
func (BelittleBlue) GoAgain() bool               { return true }
func (BelittleBlue) NotSilverAgeLegal()           {}
func (c BelittleBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
