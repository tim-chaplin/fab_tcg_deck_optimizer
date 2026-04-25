// Fyendal's Fighting Spirit — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue
// 5. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks or defends, if you have less {h} than an opposing hero, gain 1{h}."

package generic

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/simstate"
)

var fyendalsFightingSpiritTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// fyendalsFightingSpiritDamage returns attack plus the 1{h} gain credit when the current hero opts
// into LowerHealthWanter.
func fyendalsFightingSpiritDamage(attack int) int {
	if simstate.HeroWantsLowerHealth() {
		return attack + 1
	}
	return attack
}

type FyendalsFightingSpiritRed struct{}

func (FyendalsFightingSpiritRed) ID() card.ID                 { return card.FyendalsFightingSpiritRed }
func (FyendalsFightingSpiritRed) Name() string                { return "Fyendal's Fighting Spirit (Red)" }
func (FyendalsFightingSpiritRed) Cost(*card.TurnState) int                   { return 3 }
func (FyendalsFightingSpiritRed) Pitch() int                  { return 1 }
func (FyendalsFightingSpiritRed) Attack() int                 { return 7 }
func (FyendalsFightingSpiritRed) Defense() int                { return 2 }
func (FyendalsFightingSpiritRed) Types() card.TypeSet         { return fyendalsFightingSpiritTypes }
func (FyendalsFightingSpiritRed) GoAgain() bool               { return false }
func (c FyendalsFightingSpiritRed) Play(s *card.TurnState, _ *card.CardState) int { return fyendalsFightingSpiritDamage(c.Attack()) }

type FyendalsFightingSpiritYellow struct{}

func (FyendalsFightingSpiritYellow) ID() card.ID                 { return card.FyendalsFightingSpiritYellow }
func (FyendalsFightingSpiritYellow) Name() string                { return "Fyendal's Fighting Spirit (Yellow)" }
func (FyendalsFightingSpiritYellow) Cost(*card.TurnState) int                   { return 3 }
func (FyendalsFightingSpiritYellow) Pitch() int                  { return 2 }
func (FyendalsFightingSpiritYellow) Attack() int                 { return 6 }
func (FyendalsFightingSpiritYellow) Defense() int                { return 2 }
func (FyendalsFightingSpiritYellow) Types() card.TypeSet         { return fyendalsFightingSpiritTypes }
func (FyendalsFightingSpiritYellow) GoAgain() bool               { return false }
func (c FyendalsFightingSpiritYellow) Play(s *card.TurnState, _ *card.CardState) int { return fyendalsFightingSpiritDamage(c.Attack()) }

type FyendalsFightingSpiritBlue struct{}

func (FyendalsFightingSpiritBlue) ID() card.ID                 { return card.FyendalsFightingSpiritBlue }
func (FyendalsFightingSpiritBlue) Name() string                { return "Fyendal's Fighting Spirit (Blue)" }
func (FyendalsFightingSpiritBlue) Cost(*card.TurnState) int                   { return 3 }
func (FyendalsFightingSpiritBlue) Pitch() int                  { return 3 }
func (FyendalsFightingSpiritBlue) Attack() int                 { return 5 }
func (FyendalsFightingSpiritBlue) Defense() int                { return 2 }
func (FyendalsFightingSpiritBlue) Types() card.TypeSet         { return fyendalsFightingSpiritTypes }
func (FyendalsFightingSpiritBlue) GoAgain() bool               { return false }
func (c FyendalsFightingSpiritBlue) Play(s *card.TurnState, _ *card.CardState) int { return fyendalsFightingSpiritDamage(c.Attack()) }
