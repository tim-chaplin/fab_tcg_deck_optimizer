// Sigil of the Arknight — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2. Go again.
// Only printed in Blue.
//
// Text: "At the beginning of your action phase, destroy this. When this leaves the arena,
// reveal the top card of your deck. If it's an attack action card, put it into your hand."
//
// Handler fires next turn on the post-draw deck: if s.Deck()[0] is an attack action, append
// it to s.Revealed (the deck loop moves revealed cards into the hand) and pop it off the
// deck via SetDeck; non-attack reveals leave the top untouched. Damage is 0 either way — the tempo
// is the extra card, not a flat credit. The handler always logs (hit or whiff) via
// state.AddPostTriggerLogEntry so the printout names the card revealed in both cases.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfTheArknightTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type SigilOfTheArknightBlue struct{}

func (SigilOfTheArknightBlue) ID() card.ID              { return card.SigilOfTheArknightBlue }
func (SigilOfTheArknightBlue) Name() string             { return "Sigil of the Arknight" }
func (SigilOfTheArknightBlue) Cost(*card.TurnState) int { return 0 }
func (SigilOfTheArknightBlue) Pitch() int               { return 3 }
func (SigilOfTheArknightBlue) Attack() int              { return 0 }
func (SigilOfTheArknightBlue) Defense() int             { return 2 }
func (SigilOfTheArknightBlue) Types() card.TypeSet      { return sigilOfTheArknightTypes }
func (SigilOfTheArknightBlue) GoAgain() bool            { return true }
func (SigilOfTheArknightBlue) AddsFutureValue()         {}
func (c SigilOfTheArknightBlue) Play(s *card.TurnState, self *card.CardState) {
	s.RegisterStartOfTurn(c, 1, "", sigilOfTheArknightReveal)
	s.LogPlay(self)
}

// sigilOfTheArknightReveal implements the handler described in the file docstring. Logs
// the outcome on every fire — "drew X into hand" on a hit or "revealed X but didn't draw
// it" on a whiff — so the printout makes the random reveal visible either way. Empty deck
// is the silent edge case (no card to name).
func sigilOfTheArknightReveal(s *card.TurnState) int {
	deck := s.Deck()
	if len(deck) == 0 {
		return 0
	}
	top := deck[0]
	self := card.DisplayName(SigilOfTheArknightBlue{})
	if top.Types().IsAttackAction() {
		s.Revealed = append(s.Revealed, top)
		s.SetDeck(deck[1:])
		s.AddPostTriggerLogEntry(self+" drew "+card.DisplayName(top)+" into hand", self, 0)
		return 0
	}
	s.AddPostTriggerLogEntry(self+" revealed "+card.DisplayName(top)+" but didn't draw it", self, 0)
	return 0
}
