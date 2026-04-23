// Sigil of the Arknight — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2. Go again.
// Only printed in Blue.
//
// Text: "At the beginning of your action phase, destroy this. When this leaves the arena,
// reveal the top card of your deck. If it's an attack action card, put it into your hand."
//
// Modelling: Play flips AuraCreated and registers a start-of-turn AuraTrigger with Count=1.
// Next turn the sim fires the trigger on a TurnState whose Deck is the post-draw deck — the
// handler peeks s.Deck[0] and, when it's an attack action, appends that card to s.Revealed
// and shrinks s.Deck by one so the deck loop pops the top and adds it to this turn's hand.
// Non-attack reveals leave the top untouched. Count hits zero, so the sim graveyards Self.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfTheArknightTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type SigilOfTheArknightBlue struct{}

func (SigilOfTheArknightBlue) ID() card.ID              { return card.SigilOfTheArknightBlue }
func (SigilOfTheArknightBlue) Name() string             { return "Sigil of the Arknight (Blue)" }
func (SigilOfTheArknightBlue) Cost(*card.TurnState) int { return 0 }
func (SigilOfTheArknightBlue) Pitch() int               { return 3 }
func (SigilOfTheArknightBlue) Attack() int              { return 0 }
func (SigilOfTheArknightBlue) Defense() int             { return 2 }
func (SigilOfTheArknightBlue) Types() card.TypeSet      { return sigilOfTheArknightTypes }
func (SigilOfTheArknightBlue) GoAgain() bool            { return true }
func (SigilOfTheArknightBlue) AddsFutureValue()         {}
func (c SigilOfTheArknightBlue) Play(s *card.TurnState, _ *card.CardState) int {
	s.AuraCreated = true
	s.AddAuraTrigger(card.AuraTrigger{
		Self:    c,
		Type:    card.TriggerStartOfTurn,
		Count:   1,
		Handler: sigilOfTheArknightReveal,
	})
	return 0
}

// sigilOfTheArknightReveal peeks the top of the post-draw deck. If it's an attack action,
// append it to s.Revealed and pop it off s.Deck so the deck loop moves it into that turn's
// hand; otherwise the top stays in place. Damage is 0 either way — the tempo is captured by
// the extra card, not by a flat credit.
func sigilOfTheArknightReveal(s *card.TurnState) int {
	if len(s.Deck) == 0 {
		return 0
	}
	top := s.Deck[0]
	t := top.Types()
	if t.Has(card.TypeAttack) && t.Has(card.TypeAction) {
		s.Revealed = append(s.Revealed, top)
		s.Deck = s.Deck[1:]
	}
	return 0
}
