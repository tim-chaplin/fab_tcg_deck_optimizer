// Sigil of the Arknight — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2. Go again. Only
// printed in Blue.
//
// Text: "At the beginning of your action phase, destroy this. When this leaves the arena, reveal
// the top card of your deck. If it's an attack action card, put it into your hand."
//
// Modelling: Play only flips the aura-created flag — the aura enters the arena at end of the
// turn it's played and its destroy/reveal triggers fire at the start of the NEXT action phase.
// card.DelayedPlay routes PlayNextTurn through the deck loop; the callback peeks s.Deck[0]
// (the actual card about to be revealed) and, when it's an attack action, returns the card in
// ToHand so the deck loop pops it off the deck and appends it to that turn's hand. Non-attack
// reveals leave the top card untouched.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfTheArknightTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type SigilOfTheArknightBlue struct{}

func (SigilOfTheArknightBlue) ID() card.ID         { return card.SigilOfTheArknightBlue }
func (SigilOfTheArknightBlue) Name() string        { return "Sigil of the Arknight (Blue)" }
func (SigilOfTheArknightBlue) Cost(*card.TurnState) int { return 0 }
func (SigilOfTheArknightBlue) Pitch() int          { return 3 }
func (SigilOfTheArknightBlue) Attack() int         { return 0 }
func (SigilOfTheArknightBlue) Defense() int        { return 2 }
func (SigilOfTheArknightBlue) Types() card.TypeSet { return sigilOfTheArknightTypes }
func (SigilOfTheArknightBlue) GoAgain() bool       { return true }
func (SigilOfTheArknightBlue) AddsFutureValue()    {}
func (SigilOfTheArknightBlue) Play(s *card.TurnState, _ *card.CardState) int {
	s.AuraCreated = true
	return 0
}

// PlayNextTurn fires at the start of the action phase after this was played: destroy the aura
// and reveal the top card of the deck (Deck[0] — the deck slice passed in is already post-draw).
// If it's an attack action, return it in ToHand so the deck loop actually moves the card into
// the hand for this turn's best-line search rather than collapsing the tempo into a flat
// damage-equivalent.
func (SigilOfTheArknightBlue) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	if len(s.Deck) == 0 {
		return card.DelayedPlayResult{}
	}
	top := s.Deck[0]
	t := top.Types()
	if t.Has(card.TypeAttack) && t.Has(card.TypeAction) {
		return card.DelayedPlayResult{ToHand: top}
	}
	return card.DelayedPlayResult{}
}
