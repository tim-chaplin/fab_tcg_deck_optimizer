// Punch Above Your Weight — Generic Action - Attack. Cost 0. Printed power: Red 2, Yellow 2, Blue
// 2. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, you may pay {r}{r}{r}. If you do, this gets +5{p}."
//
// Simplification: Pay-{r}{r}{r}-for-+5{p} rider isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var punchAboveYourWeightTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type PunchAboveYourWeightRed struct{}

func (PunchAboveYourWeightRed) ID() card.ID                 { return card.PunchAboveYourWeightRed }
func (PunchAboveYourWeightRed) Name() string                { return "Punch Above Your Weight (Red)" }
func (PunchAboveYourWeightRed) Cost(*card.TurnState) int                   { return 0 }
func (PunchAboveYourWeightRed) Pitch() int                  { return 1 }
func (PunchAboveYourWeightRed) Attack() int                 { return 2 }
func (PunchAboveYourWeightRed) Defense() int                { return 2 }
func (PunchAboveYourWeightRed) Types() card.TypeSet         { return punchAboveYourWeightTypes }
func (PunchAboveYourWeightRed) GoAgain() bool               { return false }
func (c PunchAboveYourWeightRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type PunchAboveYourWeightYellow struct{}

func (PunchAboveYourWeightYellow) ID() card.ID                 { return card.PunchAboveYourWeightYellow }
func (PunchAboveYourWeightYellow) Name() string                { return "Punch Above Your Weight (Yellow)" }
func (PunchAboveYourWeightYellow) Cost(*card.TurnState) int                   { return 0 }
func (PunchAboveYourWeightYellow) Pitch() int                  { return 2 }
func (PunchAboveYourWeightYellow) Attack() int                 { return 2 }
func (PunchAboveYourWeightYellow) Defense() int                { return 2 }
func (PunchAboveYourWeightYellow) Types() card.TypeSet         { return punchAboveYourWeightTypes }
func (PunchAboveYourWeightYellow) GoAgain() bool               { return false }
func (c PunchAboveYourWeightYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type PunchAboveYourWeightBlue struct{}

func (PunchAboveYourWeightBlue) ID() card.ID                 { return card.PunchAboveYourWeightBlue }
func (PunchAboveYourWeightBlue) Name() string                { return "Punch Above Your Weight (Blue)" }
func (PunchAboveYourWeightBlue) Cost(*card.TurnState) int                   { return 0 }
func (PunchAboveYourWeightBlue) Pitch() int                  { return 3 }
func (PunchAboveYourWeightBlue) Attack() int                 { return 2 }
func (PunchAboveYourWeightBlue) Defense() int                { return 2 }
func (PunchAboveYourWeightBlue) Types() card.TypeSet         { return punchAboveYourWeightTypes }
func (PunchAboveYourWeightBlue) GoAgain() bool               { return false }
func (c PunchAboveYourWeightBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
