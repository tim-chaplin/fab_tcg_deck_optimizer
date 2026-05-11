// Sigil of the Arknight — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2. Go again.
// Only printed in Blue.
//
// Text: "At the beginning of your action phase, destroy this. When this leaves the arena,
// reveal the top card of your deck. If it's an attack action card, put it into your hand."
//
// Handler fires next turn on the post-draw deck: peek the top, and on an attack-action
// hit draw it into the hand. Whiffs leave the deck untouched. Damage is 0 either way —
// the tempo is the extra card, not a flat credit. The handler always logs (hit or whiff)
// so the printout names the card revealed in both cases.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (c SigilOfTheArknightBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	s.AddAura(sim.Aura{
		Trigger: sim.Trigger{TriggerType: sim.TriggerStartOfTurn, Handler: sigilOfTheArknightReveal},
		Self:    sim.CardOrTokenType{Card: c},
		Count:   1,
	})
}

// sigilOfTheArknightReveal implements the handler described in the file docstring. Logs
// the outcome on every fire — "drew X into hand" on a hit or "revealed X but didn't draw
// it" on a whiff — so the printout makes the random reveal visible either way. Empty deck
// is the silent edge case (no card to name). PeekDeck flips the cacheable bit either way
// since the reveal outcome depends on shuffle order.
func sigilOfTheArknightReveal(s *sim.TurnState, l sim.Logger, _ *sim.Trigger, a *sim.Aura) {
	s.DestroyAura(a, true)
	top, ok := s.PeekDeck()
	if !ok {
		return
	}
	self := SigilOfTheArknightBlue{}.DisplayName()
	if top.Types().IsAttackAction() {
		s.DrawOne()
		l.AppendPostTriggerf(self, 0, "%s drew %s into hand", self, top.DisplayName())
		return
	}
	l.AppendPostTriggerf(self, 0, "%s revealed %s but didn't draw it", self, top.DisplayName())
}
