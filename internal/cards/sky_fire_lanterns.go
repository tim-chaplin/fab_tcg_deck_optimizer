// Sky Fire Lanterns — Runeblade Action. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Reveal the top card of your deck. If it's <same color as this variant>, create a
// Runechant token."
//
// Peek s.Deck[0] and compare its pitch to this variant's pitch (color). On match, create
// one Runechant.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var skyFireLanternsTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction)

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
	s.AddValue(s.CreateRunechants(1))
	s.LogRider(self, 1, "Created a runechant")
}

type SkyFireLanternsRed struct{}

func (SkyFireLanternsRed) ID() ids.CardID          { return ids.SkyFireLanternsRed }
func (SkyFireLanternsRed) Name() string            { return "Sky Fire Lanterns" }
func (SkyFireLanternsRed) Cost(*sim.TurnState) int { return 0 }
func (SkyFireLanternsRed) Pitch() int              { return 1 }
func (SkyFireLanternsRed) Attack() int             { return 0 }
func (SkyFireLanternsRed) Defense() int            { return 2 }
func (SkyFireLanternsRed) Types() card.TypeSet     { return skyFireLanternsTypes }
func (SkyFireLanternsRed) GoAgain() bool           { return true }
func (c SkyFireLanternsRed) Play(s *sim.TurnState, self *sim.CardState) {
	skyFireLanternsPlay(s, self, c.Pitch())
}

type SkyFireLanternsYellow struct{}

func (SkyFireLanternsYellow) ID() ids.CardID          { return ids.SkyFireLanternsYellow }
func (SkyFireLanternsYellow) Name() string            { return "Sky Fire Lanterns" }
func (SkyFireLanternsYellow) Cost(*sim.TurnState) int { return 0 }
func (SkyFireLanternsYellow) Pitch() int              { return 2 }
func (SkyFireLanternsYellow) Attack() int             { return 0 }
func (SkyFireLanternsYellow) Defense() int            { return 2 }
func (SkyFireLanternsYellow) Types() card.TypeSet     { return skyFireLanternsTypes }
func (SkyFireLanternsYellow) GoAgain() bool           { return true }
func (c SkyFireLanternsYellow) Play(s *sim.TurnState, self *sim.CardState) {
	skyFireLanternsPlay(s, self, c.Pitch())
}

type SkyFireLanternsBlue struct{}

func (SkyFireLanternsBlue) ID() ids.CardID          { return ids.SkyFireLanternsBlue }
func (SkyFireLanternsBlue) Name() string            { return "Sky Fire Lanterns" }
func (SkyFireLanternsBlue) Cost(*sim.TurnState) int { return 0 }
func (SkyFireLanternsBlue) Pitch() int              { return 3 }
func (SkyFireLanternsBlue) Attack() int             { return 0 }
func (SkyFireLanternsBlue) Defense() int            { return 2 }
func (SkyFireLanternsBlue) Types() card.TypeSet     { return skyFireLanternsTypes }
func (SkyFireLanternsBlue) GoAgain() bool           { return true }
func (c SkyFireLanternsBlue) Play(s *sim.TurnState, self *sim.CardState) {
	skyFireLanternsPlay(s, self, c.Pitch())
}
