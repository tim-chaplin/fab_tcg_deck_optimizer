// Mauvrion Skies — Runeblade Action. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "The next Runeblade attack action card you play this turn gains go again and 'If this
// hits, create N Runechant tokens.'"
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: "if hits" is assumed, so the N Runechants are counted as N damage attributed to
// Mauvrion Skies (contingent on a qualifying next attack existing in the turn's ordering). The
// go-again grant is published by flipping GrantedGoAgain on the first matching PlayedCard in
// TurnState.CardsRemaining; the solver's chain-legality check ORs that with the card's printed
// Go again. Like Runic Reaping, "attack action card" excludes weapons.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var mauvrionSkiesTypes = map[string]bool{"Runeblade": true, "Action": true}

type MauvrionSkiesRed struct{}

func (MauvrionSkiesRed) Name() string               { return "Mauvrion Skies (Red)" }
func (MauvrionSkiesRed) Cost() int                  { return 0 }
func (MauvrionSkiesRed) Pitch() int                 { return 1 }
func (MauvrionSkiesRed) Attack() int                { return 0 }
func (MauvrionSkiesRed) Defense() int               { return 2 }
func (MauvrionSkiesRed) Types() map[string]bool     { return mauvrionSkiesTypes }
func (MauvrionSkiesRed) GoAgain() bool              { return true }
func (MauvrionSkiesRed) Play(s *card.TurnState) int { return mauvrionSkiesPlay(s, 3) }

type MauvrionSkiesYellow struct{}

func (MauvrionSkiesYellow) Name() string               { return "Mauvrion Skies (Yellow)" }
func (MauvrionSkiesYellow) Cost() int                  { return 0 }
func (MauvrionSkiesYellow) Pitch() int                 { return 2 }
func (MauvrionSkiesYellow) Attack() int                { return 0 }
func (MauvrionSkiesYellow) Defense() int               { return 2 }
func (MauvrionSkiesYellow) Types() map[string]bool     { return mauvrionSkiesTypes }
func (MauvrionSkiesYellow) GoAgain() bool              { return true }
func (MauvrionSkiesYellow) Play(s *card.TurnState) int { return mauvrionSkiesPlay(s, 2) }

type MauvrionSkiesBlue struct{}

func (MauvrionSkiesBlue) Name() string               { return "Mauvrion Skies (Blue)" }
func (MauvrionSkiesBlue) Cost() int                  { return 0 }
func (MauvrionSkiesBlue) Pitch() int                 { return 3 }
func (MauvrionSkiesBlue) Attack() int                { return 0 }
func (MauvrionSkiesBlue) Defense() int               { return 2 }
func (MauvrionSkiesBlue) Types() map[string]bool     { return mauvrionSkiesTypes }
func (MauvrionSkiesBlue) GoAgain() bool              { return true }
func (MauvrionSkiesBlue) Play(s *card.TurnState) int { return mauvrionSkiesPlay(s, 1) }

func mauvrionSkiesPlay(s *card.TurnState, n int) int {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if t["Runeblade"] && t["Action"] && t["Attack"] {
			pc.GrantedGoAgain = true
			s.AuraCreated = true
			return n
		}
	}
	// No qualifying target — both the go-again grant and the runechant rider fizzle.
	return 0
}
