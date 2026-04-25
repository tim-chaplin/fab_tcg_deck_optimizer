// Come to Fight — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 3.
//
// Text: "The next attack action card you play this turn gains +N{p}. **Go again**" (Red N=3,
// Yellow N=2, Blue N=1.)

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var comeToFightTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type ComeToFightRed struct{}

func (ComeToFightRed) ID() card.ID                 { return card.ComeToFightRed }
func (ComeToFightRed) Name() string                { return "Come to Fight (Red)" }
func (ComeToFightRed) Cost(*card.TurnState) int                   { return 1 }
func (ComeToFightRed) Pitch() int                  { return 1 }
func (ComeToFightRed) Attack() int                 { return 0 }
func (ComeToFightRed) Defense() int                { return 3 }
func (ComeToFightRed) Types() card.TypeSet         { return comeToFightTypes }
func (ComeToFightRed) GoAgain() bool               { return true }
func (ComeToFightRed) Play(s *card.TurnState, _ *card.CardState) int { return grantNextAttackActionBonus(s, 3) }

type ComeToFightYellow struct{}

func (ComeToFightYellow) ID() card.ID                 { return card.ComeToFightYellow }
func (ComeToFightYellow) Name() string                { return "Come to Fight (Yellow)" }
func (ComeToFightYellow) Cost(*card.TurnState) int                   { return 1 }
func (ComeToFightYellow) Pitch() int                  { return 2 }
func (ComeToFightYellow) Attack() int                 { return 0 }
func (ComeToFightYellow) Defense() int                { return 3 }
func (ComeToFightYellow) Types() card.TypeSet         { return comeToFightTypes }
func (ComeToFightYellow) GoAgain() bool               { return true }
func (ComeToFightYellow) Play(s *card.TurnState, _ *card.CardState) int { return grantNextAttackActionBonus(s, 2) }

type ComeToFightBlue struct{}

func (ComeToFightBlue) ID() card.ID                 { return card.ComeToFightBlue }
func (ComeToFightBlue) Name() string                { return "Come to Fight (Blue)" }
func (ComeToFightBlue) Cost(*card.TurnState) int                   { return 1 }
func (ComeToFightBlue) Pitch() int                  { return 3 }
func (ComeToFightBlue) Attack() int                 { return 0 }
func (ComeToFightBlue) Defense() int                { return 3 }
func (ComeToFightBlue) Types() card.TypeSet         { return comeToFightTypes }
func (ComeToFightBlue) GoAgain() bool               { return true }
func (ComeToFightBlue) Play(s *card.TurnState, _ *card.CardState) int { return grantNextAttackActionBonus(s, 1) }
