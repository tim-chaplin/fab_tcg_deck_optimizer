// Sky Fire Lanterns — Runeblade Action. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Reveal the top card of your deck. If it's <same color as this variant>, create a
// Runechant token."
//
// Peek s.Deck[0] and compare its pitch to this variant's pitch (color). On match, create
// one Runechant.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// skyFireLanternsPlay emits the chain step then writes a runechant rider sub-line under
// self when the deck-top card matches this variant's pitch (color). Reads the deck top via
// s.Deck() so the cacheable bit flips — whether the rider fires depends on shuffle order.
func skyFireLanternsPlay(s *sim.TurnState, self *sim.CardState, selfPitch int) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	deck := s.Deck()
	if len(deck) == 0 || deck[0].Pitch() != selfPitch {
		return
	}
	s.CreateRunechants(1)
	s.LogRider(self, 1, "Created a runechant")
}

func (c SkyFireLanternsRed) Play(s *sim.TurnState, self *sim.CardState) {
	skyFireLanternsPlay(s, self, c.Pitch())
}

func (c SkyFireLanternsYellow) Play(s *sim.TurnState, self *sim.CardState) {
	skyFireLanternsPlay(s, self, c.Pitch())
}

func (c SkyFireLanternsBlue) Play(s *sim.TurnState, self *sim.CardState) {
	skyFireLanternsPlay(s, self, c.Pitch())
}
