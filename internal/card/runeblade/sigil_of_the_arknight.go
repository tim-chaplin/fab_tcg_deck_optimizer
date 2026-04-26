// Sigil of the Arknight — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2. Go again.
// Only printed in Blue.
//
// Text: "At the beginning of your action phase, destroy this. When this leaves the arena,
// reveal the top card of your deck. If it's an attack action card, put it into your hand."
//
// Handler fires next turn on the post-draw deck: if s.Deck[0] is an attack action, append it
// to s.Revealed (the deck loop moves revealed cards into the hand) and pop it off s.Deck;
// non-attack reveals leave the top untouched. Damage is 0 either way — the tempo is the
// extra card, not a flat credit.

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
func (c SigilOfTheArknightBlue) Play(s *card.TurnState, _ *card.CardState) int {
	s.AddAuraTrigger(card.AuraTrigger{
		Self:    c,
		Type:    card.TriggerStartOfTurn,
		Count:   1,
		Handler: sigilOfTheArknightReveal,
	})
	return 0
}

// sigilOfTheArknightReveal implements the handler described in the file docstring.
func sigilOfTheArknightReveal(s *card.TurnState) int {
	if len(s.Deck) == 0 {
		return 0
	}
	top := s.Deck[0]
	if top.Types().IsAttackAction() {
		s.Revealed = append(s.Revealed, top)
		s.Deck = s.Deck[1:]
	}
	return 0
}
