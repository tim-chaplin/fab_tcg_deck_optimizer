// Sutcliffe's Research Notes — Runeblade Action. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Reveal the top N cards of your deck. Create a Runechant token for each Runeblade attack
// action card revealed this way, then put the cards on top of your deck in any order." (N = 3
// Red / 2 Yellow / 1 Blue.)
//
// Simplification: scan the top N cards of the remaining deck; credit +1 per Runeblade attack
// action card revealed and set AuraCreated if at least one Runechant is made. Opts out of the
// hand-evaluation memo because the result depends on deck composition. The re-ordering clause is
// ignored — we don't model future-turn draw order.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sutcliffesResearchNotesTypes = map[string]bool{"Runeblade": true, "Action": true}

func sutcliffesResearchNotesPlay(revealCount int, s *card.TurnState) int {
	n := revealCount
	if n > len(s.Deck) {
		n = len(s.Deck)
	}
	count := 0
	for i := 0; i < n; i++ {
		t := s.Deck[i].Types()
		if t["Runeblade"] && t["Attack"] && t["Action"] {
			count++
		}
	}
	if count > 0 {
		s.AuraCreated = true
	}
	return count
}

type SutcliffesResearchNotesRed struct{}

func (SutcliffesResearchNotesRed) Name() string                 { return "Sutcliffe's Research Notes (Red)" }
func (SutcliffesResearchNotesRed) Cost() int                    { return 1 }
func (SutcliffesResearchNotesRed) Pitch() int                   { return 1 }
func (SutcliffesResearchNotesRed) Attack() int                  { return 0 }
func (SutcliffesResearchNotesRed) Defense() int                 { return 2 }
func (SutcliffesResearchNotesRed) Types() map[string]bool       { return sutcliffesResearchNotesTypes }
func (SutcliffesResearchNotesRed) GoAgain() bool                { return true }
func (SutcliffesResearchNotesRed) NoMemo()                      {}
func (SutcliffesResearchNotesRed) Play(s *card.TurnState) int   { return sutcliffesResearchNotesPlay(3, s) }

type SutcliffesResearchNotesYellow struct{}

func (SutcliffesResearchNotesYellow) Name() string                 { return "Sutcliffe's Research Notes (Yellow)" }
func (SutcliffesResearchNotesYellow) Cost() int                    { return 1 }
func (SutcliffesResearchNotesYellow) Pitch() int                   { return 2 }
func (SutcliffesResearchNotesYellow) Attack() int                  { return 0 }
func (SutcliffesResearchNotesYellow) Defense() int                 { return 2 }
func (SutcliffesResearchNotesYellow) Types() map[string]bool       { return sutcliffesResearchNotesTypes }
func (SutcliffesResearchNotesYellow) GoAgain() bool                { return true }
func (SutcliffesResearchNotesYellow) NoMemo()                      {}
func (SutcliffesResearchNotesYellow) Play(s *card.TurnState) int   { return sutcliffesResearchNotesPlay(2, s) }

type SutcliffesResearchNotesBlue struct{}

func (SutcliffesResearchNotesBlue) Name() string                 { return "Sutcliffe's Research Notes (Blue)" }
func (SutcliffesResearchNotesBlue) Cost() int                    { return 1 }
func (SutcliffesResearchNotesBlue) Pitch() int                   { return 3 }
func (SutcliffesResearchNotesBlue) Attack() int                  { return 0 }
func (SutcliffesResearchNotesBlue) Defense() int                 { return 2 }
func (SutcliffesResearchNotesBlue) Types() map[string]bool       { return sutcliffesResearchNotesTypes }
func (SutcliffesResearchNotesBlue) GoAgain() bool                { return true }
func (SutcliffesResearchNotesBlue) NoMemo()                      {}
func (SutcliffesResearchNotesBlue) Play(s *card.TurnState) int   { return sutcliffesResearchNotesPlay(1, s) }
