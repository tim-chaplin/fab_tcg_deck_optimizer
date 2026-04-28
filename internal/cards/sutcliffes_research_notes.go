// Sutcliffe's Research Notes — Runeblade Action. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Reveal the top N cards of your deck. Create a Runechant token for each Runeblade attack
// action card revealed this way, then put the cards on top of your deck in any order." (N = 3
// Red / 2 Yellow / 1 Blue.)
//
// Scan the top N cards of s.Deck; credit +1 per Runeblade attack action card revealed. The
// re-ordering clause isn't modelled.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var sutcliffesResearchNotesTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction)

// sutcliffesResearchNotesPlay scans the top revealCount cards of s.Deck and creates one
// runechant per Runeblade attack action card found, emitting the rider sub-line under self
// when any are created.
func sutcliffesResearchNotesPlay(s *sim.TurnState, self *sim.CardState, revealCount int) {
	s.ApplyAndLogEffectiveAttack(self)
	n := revealCount
	if n > len(s.Deck) {
		n = len(s.Deck)
	}
	count := 0
	for i := 0; i < n; i++ {
		t := s.Deck[i].Types()
		if t.Has(card.TypeRuneblade) && t.IsAttackAction() {
			count++
		}
	}
	s.CreateAndLogRunechantsOnPlay(self, count)
}

type SutcliffesResearchNotesRed struct{}

func (SutcliffesResearchNotesRed) ID() ids.CardID          { return ids.SutcliffesResearchNotesRed }
func (SutcliffesResearchNotesRed) Name() string            { return "Sutcliffe's Research Notes" }
func (SutcliffesResearchNotesRed) Cost(*sim.TurnState) int { return 1 }
func (SutcliffesResearchNotesRed) Pitch() int              { return 1 }
func (SutcliffesResearchNotesRed) Attack() int             { return 0 }
func (SutcliffesResearchNotesRed) Defense() int            { return 2 }
func (SutcliffesResearchNotesRed) Types() card.TypeSet     { return sutcliffesResearchNotesTypes }
func (SutcliffesResearchNotesRed) GoAgain() bool           { return true }

// not implemented: top-of-deck reordering clause
func (SutcliffesResearchNotesRed) NotImplemented() {}
func (SutcliffesResearchNotesRed) Play(s *sim.TurnState, self *sim.CardState) {
	sutcliffesResearchNotesPlay(s, self, 3)
}

type SutcliffesResearchNotesYellow struct{}

func (SutcliffesResearchNotesYellow) ID() ids.CardID          { return ids.SutcliffesResearchNotesYellow }
func (SutcliffesResearchNotesYellow) Name() string            { return "Sutcliffe's Research Notes" }
func (SutcliffesResearchNotesYellow) Cost(*sim.TurnState) int { return 1 }
func (SutcliffesResearchNotesYellow) Pitch() int              { return 2 }
func (SutcliffesResearchNotesYellow) Attack() int             { return 0 }
func (SutcliffesResearchNotesYellow) Defense() int            { return 2 }
func (SutcliffesResearchNotesYellow) Types() card.TypeSet     { return sutcliffesResearchNotesTypes }
func (SutcliffesResearchNotesYellow) GoAgain() bool           { return true }

// not implemented: top-of-deck reordering clause
func (SutcliffesResearchNotesYellow) NotImplemented() {}
func (SutcliffesResearchNotesYellow) Play(s *sim.TurnState, self *sim.CardState) {
	sutcliffesResearchNotesPlay(s, self, 2)
}

type SutcliffesResearchNotesBlue struct{}

func (SutcliffesResearchNotesBlue) ID() ids.CardID          { return ids.SutcliffesResearchNotesBlue }
func (SutcliffesResearchNotesBlue) Name() string            { return "Sutcliffe's Research Notes" }
func (SutcliffesResearchNotesBlue) Cost(*sim.TurnState) int { return 1 }
func (SutcliffesResearchNotesBlue) Pitch() int              { return 3 }
func (SutcliffesResearchNotesBlue) Attack() int             { return 0 }
func (SutcliffesResearchNotesBlue) Defense() int            { return 2 }
func (SutcliffesResearchNotesBlue) Types() card.TypeSet     { return sutcliffesResearchNotesTypes }
func (SutcliffesResearchNotesBlue) GoAgain() bool           { return true }

// not implemented: top-of-deck reordering clause
func (SutcliffesResearchNotesBlue) NotImplemented() {}
func (SutcliffesResearchNotesBlue) Play(s *sim.TurnState, self *sim.CardState) {
	sutcliffesResearchNotesPlay(s, self, 1)
}
