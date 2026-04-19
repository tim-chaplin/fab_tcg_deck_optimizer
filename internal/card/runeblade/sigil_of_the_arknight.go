// Sigil of the Arknight — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2. Go again. Only
// printed in Blue.
// Text: "At the beginning of your action phase, destroy this. When this leaves the arena, reveal
// the top card of your deck. If it's an attack action card, put it into your hand."
//
// Simplification: assume the Sigil enters then leaves at the start of next turn. By then the
// hero has drawn Intelligence cards into their next hand, so the revealed card is at index
// Intelligence in the remaining deck. If it's an attack action, credit +3 value (the draw-a-card
// damage-equivalent); otherwise 0. Cross-turn aura-persistence isn't modelled — the value is
// collapsed into the immediate Play return. Deck too short to reach that index → 0.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import (
	"log"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/simstate"
)

var sigilOfTheArknightTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type SigilOfTheArknightBlue struct{}

func (SigilOfTheArknightBlue) ID() card.ID                 { return card.SigilOfTheArknightBlue }
func (SigilOfTheArknightBlue) Name() string           { return "Sigil of the Arknight (Blue)" }
func (SigilOfTheArknightBlue) Cost(*card.TurnState) int              { return 0 }
func (SigilOfTheArknightBlue) Pitch() int             { return 3 }
func (SigilOfTheArknightBlue) Attack() int            { return 0 }
func (SigilOfTheArknightBlue) Defense() int           { return 2 }
func (SigilOfTheArknightBlue) Types() card.TypeSet    { return sigilOfTheArknightTypes }
func (SigilOfTheArknightBlue) GoAgain() bool          { return true }
func (SigilOfTheArknightBlue) NoMemo()                {}  // value depends on the state of the deck
func (SigilOfTheArknightBlue) Play(s *card.TurnState) int {
	s.AuraCreated = true
	// Next turn's draw takes the first Intelligence cards; the reveal then peeks at the next one.
	if simstate.CurrentHero == nil {
		log.Fatal("Sigil of the Arknight played with simstate.CurrentHero unset")
	}
	idx := simstate.CurrentHero.Intelligence()
	if idx >= len(s.Deck) {
		return 0
	}
	t := s.Deck[idx].Types()
	if t.Has(card.TypeAttack) && t.Has(card.TypeAction) {
		return card.DrawValue
	}
	return 0
}
