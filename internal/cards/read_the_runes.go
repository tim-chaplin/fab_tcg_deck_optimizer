// Read the Runes — Runeblade Action. Cost 0, Defense 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Create N Runechant tokens." (Red N=3, Yellow N=2, Blue N=1.)
//
// Play returns N and sets AuraCreated so later cards this turn see the Runechants.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var readTheRunesTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction)

type ReadTheRunesRed struct{}

func (ReadTheRunesRed) ID() ids.CardID           { return ids.ReadTheRunesRed }
func (ReadTheRunesRed) Name() string             { return "Read the Runes" }
func (ReadTheRunesRed) Cost(*card.TurnState) int { return 0 }
func (ReadTheRunesRed) Pitch() int               { return 1 }
func (ReadTheRunesRed) Attack() int              { return 0 }
func (ReadTheRunesRed) Defense() int             { return 2 }
func (ReadTheRunesRed) Types() card.TypeSet      { return readTheRunesTypes }
func (ReadTheRunesRed) GoAgain() bool            { return false }
func (ReadTheRunesRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	s.CreateAndLogRunechantsOnPlay(self, 3)
}

type ReadTheRunesYellow struct{}

func (ReadTheRunesYellow) ID() ids.CardID           { return ids.ReadTheRunesYellow }
func (ReadTheRunesYellow) Name() string             { return "Read the Runes" }
func (ReadTheRunesYellow) Cost(*card.TurnState) int { return 0 }
func (ReadTheRunesYellow) Pitch() int               { return 2 }
func (ReadTheRunesYellow) Attack() int              { return 0 }
func (ReadTheRunesYellow) Defense() int             { return 2 }
func (ReadTheRunesYellow) Types() card.TypeSet      { return readTheRunesTypes }
func (ReadTheRunesYellow) GoAgain() bool            { return false }
func (ReadTheRunesYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	s.CreateAndLogRunechantsOnPlay(self, 2)
}

type ReadTheRunesBlue struct{}

func (ReadTheRunesBlue) ID() ids.CardID           { return ids.ReadTheRunesBlue }
func (ReadTheRunesBlue) Name() string             { return "Read the Runes" }
func (ReadTheRunesBlue) Cost(*card.TurnState) int { return 0 }
func (ReadTheRunesBlue) Pitch() int               { return 3 }
func (ReadTheRunesBlue) Attack() int              { return 0 }
func (ReadTheRunesBlue) Defense() int             { return 2 }
func (ReadTheRunesBlue) Types() card.TypeSet      { return readTheRunesTypes }
func (ReadTheRunesBlue) GoAgain() bool            { return false }
func (ReadTheRunesBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	s.CreateAndLogRunechantsOnPlay(self, 1)
}
