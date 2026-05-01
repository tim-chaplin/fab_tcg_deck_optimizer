// Sigil of the Arknight — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2. Go again.
// Only printed in Blue.
//
// Text: "At the beginning of your action phase, destroy this. When this leaves the arena,
// reveal the top card of your deck. If it's an attack action card, put it into your hand."
//
// Handler fires next turn on the post-draw deck: if s.Deck[0] is an attack action, append
// it to s.Revealed (the deck loop moves revealed cards into the hand) and pop it off
// s.Deck; non-attack reveals leave the top untouched. Damage is 0 either way — the tempo
// is the extra card, not a flat credit. The handler always logs (hit or whiff) via
// state.AddPostTriggerLogEntry so the printout names the card revealed in both cases.

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
	s.RegisterStartOfTurn(c, 1, "", sigilOfTheArknightReveal)
	s.LogChain(self, 0)
}

// sigilOfTheArknightReveal implements the handler described in the file docstring. Logs
// the outcome on every fire — "drew X into hand" on a hit or "revealed X but didn't draw
// it" on a whiff — so the printout makes the random reveal visible either way. Empty deck
// is the silent edge case (no card to name). Pops the top via PopDeckTop on a hit; on a
// whiff puts the card back via PrependToDeck so the deck order is preserved (both verbs
// flip the cacheable bit, which is what we want — the reveal outcome depends on shuffle).
func sigilOfTheArknightReveal(s *sim.TurnState, _ *sim.AuraTrigger) int {
	top, ok := s.PopDeckTop()
	if !ok {
		return 0
	}
	self := sim.DisplayName(SigilOfTheArknightBlue{})
	if top.Types().IsAttackAction() {
		s.Revealed = append(s.Revealed, top)
		s.LogPostTriggerf(self, 0, "%s drew %s into hand", self, sim.DisplayName(top))
		return 0
	}
	// Whiff — restore the deck top so non-attack reveals leave deck order untouched.
	s.PrependToDeck(top)
	s.LogPostTriggerf(self, 0, "%s revealed %s but didn't draw it", self, sim.DisplayName(top))
	return 0
}
