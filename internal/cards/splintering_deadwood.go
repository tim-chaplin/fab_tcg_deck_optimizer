// Splintering Deadwood — Runeblade Action - Attack. Cost 3, Defense 3.
// Printed power: Red 7, Yellow 6, Blue 5.
// Text: "When this attacks or hits, you may destroy an aura you control. If you do, create a
// Runechant token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var splinteringDeadwoodTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type SplinteringDeadwoodRed struct{}

func (SplinteringDeadwoodRed) ID() ids.CardID           { return ids.SplinteringDeadwoodRed }
func (SplinteringDeadwoodRed) Name() string             { return "Splintering Deadwood" }
func (SplinteringDeadwoodRed) Cost(*card.TurnState) int { return 3 }
func (SplinteringDeadwoodRed) Pitch() int               { return 1 }
func (SplinteringDeadwoodRed) Attack() int              { return 7 }
func (SplinteringDeadwoodRed) Defense() int             { return 3 }
func (SplinteringDeadwoodRed) Types() card.TypeSet      { return splinteringDeadwoodTypes }
func (SplinteringDeadwoodRed) GoAgain() bool            { return false }

// not implemented: aura-swap rider modelled net-zero; no tempo credit for trading a weak aura
// for a Runechant
func (SplinteringDeadwoodRed) NotImplemented() {}
func (SplinteringDeadwoodRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type SplinteringDeadwoodYellow struct{}

func (SplinteringDeadwoodYellow) ID() ids.CardID           { return ids.SplinteringDeadwoodYellow }
func (SplinteringDeadwoodYellow) Name() string             { return "Splintering Deadwood" }
func (SplinteringDeadwoodYellow) Cost(*card.TurnState) int { return 3 }
func (SplinteringDeadwoodYellow) Pitch() int               { return 2 }
func (SplinteringDeadwoodYellow) Attack() int              { return 6 }
func (SplinteringDeadwoodYellow) Defense() int             { return 3 }
func (SplinteringDeadwoodYellow) Types() card.TypeSet      { return splinteringDeadwoodTypes }
func (SplinteringDeadwoodYellow) GoAgain() bool            { return false }

// not implemented: aura-swap rider modelled net-zero; no tempo credit for trading a weak aura
// for a Runechant
func (SplinteringDeadwoodYellow) NotImplemented() {}
func (SplinteringDeadwoodYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type SplinteringDeadwoodBlue struct{}

func (SplinteringDeadwoodBlue) ID() ids.CardID           { return ids.SplinteringDeadwoodBlue }
func (SplinteringDeadwoodBlue) Name() string             { return "Splintering Deadwood" }
func (SplinteringDeadwoodBlue) Cost(*card.TurnState) int { return 3 }
func (SplinteringDeadwoodBlue) Pitch() int               { return 3 }
func (SplinteringDeadwoodBlue) Attack() int              { return 5 }
func (SplinteringDeadwoodBlue) Defense() int             { return 3 }
func (SplinteringDeadwoodBlue) Types() card.TypeSet      { return splinteringDeadwoodTypes }
func (SplinteringDeadwoodBlue) GoAgain() bool            { return false }

// not implemented: aura-swap rider modelled net-zero; no tempo credit for trading a weak aura
// for a Runechant
func (SplinteringDeadwoodBlue) NotImplemented() {}
func (SplinteringDeadwoodBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
