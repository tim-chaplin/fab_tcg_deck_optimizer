// Sutcliffe's Research Notes — Runeblade Action. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Reveal the top N cards of your deck. Create a Runechant token for each Runeblade attack
// action card revealed this way, then put the cards on top of your deck in any order." (N = 3
// Red / 2 Yellow / 1 Blue.)
//
// Scan the top N cards of ge.Deck; credit +1 per Runeblade attack action card revealed. The
// post-reveal reorder isn't modelled.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// sutcliffesResearchNotesPlay scans the top revealCount cards of the deck and creates one
// runechant per Runeblade attack action card found.
func sutcliffesResearchNotesPlay(ge card.GameEngine, l card.Logger, self *card.CardState, revealCount int) {
	count := 0
	for _, c := range ge.PeekTopN(revealCount) {
		t := c.Types(nil)
		if t.Has(card.TypeRuneblade) && t.IsAttackAction() {
			count++
		}
	}
	createRunechantsAndLog(ge, l, self.Card.DisplayName(), count)
}

func (SutcliffesResearchNotesRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	sutcliffesResearchNotesPlay(ge, l, self, 3)
}

func (SutcliffesResearchNotesYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	sutcliffesResearchNotesPlay(ge, l, self, 2)
}

func (SutcliffesResearchNotesBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	sutcliffesResearchNotesPlay(ge, l, self, 1)
}
