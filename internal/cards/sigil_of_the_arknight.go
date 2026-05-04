// Sigil of the Arknight — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2. Go again.
// Only printed in Blue.
//
// Text: "At the beginning of your action phase, destroy this. When this leaves the arena,
// reveal the top card of your deck. If it's an attack action card, put it into your hand."
//
// Handler fires next turn on the post-draw deck: peek the top, and on an attack-action
// hit pop it into s.Revealed (the deck loop moves revealed cards into the hand). Whiffs
// leave the deck untouched. Damage is 0 either way — the tempo is the extra card, not a
// flat credit. The handler always logs (hit or whiff) so the printout names the card
// revealed in both cases.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var sigilOfTheArknightTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type SigilOfTheArknightBlue struct{}

func (SigilOfTheArknightBlue) ID() ids.CardID          { return ids.SigilOfTheArknightBlue }
func (SigilOfTheArknightBlue) Name() string            { return "Sigil of the Arknight" }
func (SigilOfTheArknightBlue) Cost(*sim.TurnState) int { return 0 }
func (SigilOfTheArknightBlue) Pitch() int              { return 3 }
func (SigilOfTheArknightBlue) Attack() int             { return 0 }
func (SigilOfTheArknightBlue) Defense() int            { return 2 }
func (SigilOfTheArknightBlue) Types() card.TypeSet     { return sigilOfTheArknightTypes }
func (SigilOfTheArknightBlue) GoAgain() bool           { return true }
func (SigilOfTheArknightBlue) AddsFutureValue()        {}
func (c SigilOfTheArknightBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.AddAura(sim.Aura{
		Self:        c,
		TriggerType: sim.TriggerStartOfTurn,
		Count:       1,
		Handler:     sigilOfTheArknightReveal,
	})
	s.Log(self, 0)
}

// sigilOfTheArknightReveal implements the handler described in the file docstring. Logs
// the outcome on every fire — "drew X into hand" on a hit or "revealed X but didn't draw
// it" on a whiff — so the printout makes the random reveal visible either way. Empty deck
// is the silent edge case (no card to name). PeekDeck flips the cacheable bit either way
// since the reveal outcome depends on shuffle order.
func sigilOfTheArknightReveal(s *sim.TurnState, t *sim.Aura) {
	s.DestroyAura(t, true)
	top, ok := s.PeekDeck()
	if !ok {
		return
	}
	self := sim.DisplayName(SigilOfTheArknightBlue{})
	if top.Types().IsAttackAction() {
		s.PopDeckTop()
		s.Revealed = append(s.Revealed, top)
		s.LogPostTriggerf(self, 0, "%s drew %s into hand", self, sim.DisplayName(top))
		return
	}
	s.LogPostTriggerf(self, 0, "%s revealed %s but didn't draw it", self, sim.DisplayName(top))
}
