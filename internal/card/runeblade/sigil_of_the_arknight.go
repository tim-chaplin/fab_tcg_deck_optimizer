// Sigil of the Arknight — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2. Go again. Only
// printed in Blue.
// Text: "At the beginning of your action phase, destroy this. When this leaves the arena, reveal
// the top card of your deck. If it's an attack action card, put it into your hand."
//
// Simplification: assume the Sigil enters then leaves next turn; the expected value of the
// reveal is (fraction of attack action cards remaining in the deck) × 3, where 3 is the value we
// ascribe to a drawn card (matching Drawn to the Dark Dimension's draw-a-card valuation). Cross-
// turn aura-persistence isn't modelled; the value is collapsed to an immediate return on play.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfTheArknightTypes = map[string]bool{"Runeblade": true, "Action": true, "Aura": true}

type SigilOfTheArknightBlue struct{}

func (SigilOfTheArknightBlue) Name() string           { return "Sigil of the Arknight (Blue)" }
func (SigilOfTheArknightBlue) Cost() int              { return 0 }
func (SigilOfTheArknightBlue) Pitch() int             { return 3 }
func (SigilOfTheArknightBlue) Attack() int            { return 0 }
func (SigilOfTheArknightBlue) Defense() int           { return 2 }
func (SigilOfTheArknightBlue) Types() map[string]bool { return sigilOfTheArknightTypes }
func (SigilOfTheArknightBlue) GoAgain() bool          { return true }
func (SigilOfTheArknightBlue) NoMemo()                {}  // value depends on the state of the deck
func (SigilOfTheArknightBlue) Play(s *card.TurnState) int {
	s.AuraCreated = true
	if len(s.Deck) == 0 {
		return 0
	}
	attackActions := 0
	for _, c := range s.Deck {
		t := c.Types()
		if t["Attack"] && t["Action"] {
			attackActions++
		}
	}
	// Expected value: P(top card is attack action) × 3.
	return (attackActions * 3) / len(s.Deck)
}
