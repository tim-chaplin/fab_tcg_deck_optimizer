// Sutcliffe's Research Notes — Runeblade Action. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Reveal the top N cards of your deck. Create a Runechant token for each Runeblade attack
// action card revealed this way, then put the cards on top of your deck in any order." (N = 3
// Red / 2 Yellow / 1 Blue.)
//
// Scan the top N cards of s.Deck; credit +1 per Runeblade attack action card revealed. The
// post-reveal reorder isn't modelled.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// sutcliffesResearchNotesPlay scans the top revealCount cards of the deck and creates one
// runechant per Runeblade attack action card found, emitting the rider sub-line under self
// when any are created. Reads the top via PeekTopN so the cacheable bit flips — the
// runechant count produced depends on shuffle order.
func sutcliffesResearchNotesPlay(s *sim.TurnState, self *sim.CardState, revealCount int) {
	dmg := self.DealEffectiveAttack(s)
	s.Log(self, dmg)
	count := 0
	for _, c := range s.PeekTopN(revealCount) {
		t := c.Types()
		if t.Has(card.TypeRuneblade) && t.IsAttackAction() {
			count++
		}
	}
	s.CreateRunechants(count)
	s.LogRiderf(self, count, "Created %d runechants", count)
}

func (SutcliffesResearchNotesRed) Play(s *sim.TurnState, self *sim.CardState) {
	sutcliffesResearchNotesPlay(s, self, 3)
}

func (SutcliffesResearchNotesYellow) Play(s *sim.TurnState, self *sim.CardState) {
	sutcliffesResearchNotesPlay(s, self, 2)
}

func (SutcliffesResearchNotesBlue) Play(s *sim.TurnState, self *sim.CardState) {
	sutcliffesResearchNotesPlay(s, self, 1)
}
