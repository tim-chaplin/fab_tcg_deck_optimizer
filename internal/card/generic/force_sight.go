// Force Sight — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next attack action card you play this turn gains +N{p}. If Force Sight is played from
// arsenal, **opt 2**. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var forceSightTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type ForceSightRed struct{}

func (ForceSightRed) ID() card.ID                 { return card.ForceSightRed }
func (ForceSightRed) Name() string                { return "Force Sight (Red)" }
func (ForceSightRed) Cost(*card.TurnState) int                   { return 1 }
func (ForceSightRed) Pitch() int                  { return 1 }
func (ForceSightRed) Attack() int                 { return 0 }
func (ForceSightRed) Defense() int                { return 2 }
func (ForceSightRed) Types() card.TypeSet         { return forceSightTypes }
func (ForceSightRed) GoAgain() bool               { return true }
// not implemented: arsenal-gated Opt 2
func (ForceSightRed) NotImplemented()             {}
func (ForceSightRed) Play(s *card.TurnState, _ *card.CardState) int { return grantNextAttackActionBonus(s, 3) }

type ForceSightYellow struct{}

func (ForceSightYellow) ID() card.ID                 { return card.ForceSightYellow }
func (ForceSightYellow) Name() string                { return "Force Sight (Yellow)" }
func (ForceSightYellow) Cost(*card.TurnState) int                   { return 1 }
func (ForceSightYellow) Pitch() int                  { return 2 }
func (ForceSightYellow) Attack() int                 { return 0 }
func (ForceSightYellow) Defense() int                { return 2 }
func (ForceSightYellow) Types() card.TypeSet         { return forceSightTypes }
func (ForceSightYellow) GoAgain() bool               { return true }
// not implemented: arsenal-gated Opt 2
func (ForceSightYellow) NotImplemented()             {}
func (ForceSightYellow) Play(s *card.TurnState, _ *card.CardState) int { return grantNextAttackActionBonus(s, 2) }

type ForceSightBlue struct{}

func (ForceSightBlue) ID() card.ID                 { return card.ForceSightBlue }
func (ForceSightBlue) Name() string                { return "Force Sight (Blue)" }
func (ForceSightBlue) Cost(*card.TurnState) int                   { return 1 }
func (ForceSightBlue) Pitch() int                  { return 3 }
func (ForceSightBlue) Attack() int                 { return 0 }
func (ForceSightBlue) Defense() int                { return 2 }
func (ForceSightBlue) Types() card.TypeSet         { return forceSightTypes }
func (ForceSightBlue) GoAgain() bool               { return true }
// not implemented: arsenal-gated Opt 2
func (ForceSightBlue) NotImplemented()             {}
func (ForceSightBlue) Play(s *card.TurnState, _ *card.CardState) int { return grantNextAttackActionBonus(s, 1) }
