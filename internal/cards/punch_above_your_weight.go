// Punch Above Your Weight — Generic Action - Attack. Cost 0. Printed power: Red 2, Yellow 2, Blue
// 2. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, you may pay {r}{r}{r}. If you do, this gets +5{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var punchAboveYourWeightTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type PunchAboveYourWeightRed struct{}

func (PunchAboveYourWeightRed) ID() ids.CardID           { return ids.PunchAboveYourWeightRed }
func (PunchAboveYourWeightRed) Name() string             { return "Punch Above Your Weight" }
func (PunchAboveYourWeightRed) Cost(*card.TurnState) int { return 0 }
func (PunchAboveYourWeightRed) Pitch() int               { return 1 }
func (PunchAboveYourWeightRed) Attack() int              { return 2 }
func (PunchAboveYourWeightRed) Defense() int             { return 2 }
func (PunchAboveYourWeightRed) Types() card.TypeSet      { return punchAboveYourWeightTypes }
func (PunchAboveYourWeightRed) GoAgain() bool            { return false }

// not implemented: pay-{r}{r}{r}-for-+5{p} mode
func (PunchAboveYourWeightRed) NotImplemented() {}
func (c PunchAboveYourWeightRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type PunchAboveYourWeightYellow struct{}

func (PunchAboveYourWeightYellow) ID() ids.CardID           { return ids.PunchAboveYourWeightYellow }
func (PunchAboveYourWeightYellow) Name() string             { return "Punch Above Your Weight" }
func (PunchAboveYourWeightYellow) Cost(*card.TurnState) int { return 0 }
func (PunchAboveYourWeightYellow) Pitch() int               { return 2 }
func (PunchAboveYourWeightYellow) Attack() int              { return 2 }
func (PunchAboveYourWeightYellow) Defense() int             { return 2 }
func (PunchAboveYourWeightYellow) Types() card.TypeSet      { return punchAboveYourWeightTypes }
func (PunchAboveYourWeightYellow) GoAgain() bool            { return false }

// not implemented: pay-{r}{r}{r}-for-+5{p} mode
func (PunchAboveYourWeightYellow) NotImplemented() {}
func (c PunchAboveYourWeightYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type PunchAboveYourWeightBlue struct{}

func (PunchAboveYourWeightBlue) ID() ids.CardID           { return ids.PunchAboveYourWeightBlue }
func (PunchAboveYourWeightBlue) Name() string             { return "Punch Above Your Weight" }
func (PunchAboveYourWeightBlue) Cost(*card.TurnState) int { return 0 }
func (PunchAboveYourWeightBlue) Pitch() int               { return 3 }
func (PunchAboveYourWeightBlue) Attack() int              { return 2 }
func (PunchAboveYourWeightBlue) Defense() int             { return 2 }
func (PunchAboveYourWeightBlue) Types() card.TypeSet      { return punchAboveYourWeightTypes }
func (PunchAboveYourWeightBlue) GoAgain() bool            { return false }

// not implemented: pay-{r}{r}{r}-for-+5{p} mode
func (PunchAboveYourWeightBlue) NotImplemented() {}
func (c PunchAboveYourWeightBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
