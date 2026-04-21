// Mauvrion Skies — Runeblade Action. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "The next Runeblade attack action card you play this turn gains go again and 'If this
// hits, create N Runechant tokens.'"
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelling: the go-again grant is published via GrantedGoAgain on the first matching
// CardState in CardsRemaining. The N Runechants are credited only when that target's printed
// Attack() satisfies card.LikelyToHit — i.e. the "if this hits" clause is treated as firing
// when the opponent would rather eat the damage than over-block. Targets whose printed power
// lands in the blockable range drop the Runechants but still keep go-again. "Attack action
// card" excludes weapons.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var mauvrionSkiesTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction)

type MauvrionSkiesRed struct{}

func (MauvrionSkiesRed) ID() card.ID                 { return card.MauvrionSkiesRed }
func (MauvrionSkiesRed) Name() string               { return "Mauvrion Skies (Red)" }
func (MauvrionSkiesRed) Cost(*card.TurnState) int                  { return 0 }
func (MauvrionSkiesRed) Pitch() int                 { return 1 }
func (MauvrionSkiesRed) Attack() int                { return 0 }
func (MauvrionSkiesRed) Defense() int               { return 2 }
func (MauvrionSkiesRed) Types() card.TypeSet        { return mauvrionSkiesTypes }
func (MauvrionSkiesRed) GoAgain() bool              { return true }
func (MauvrionSkiesRed) Play(s *card.TurnState, _ *card.CardState) int { return mauvrionSkiesPlay(s, 3) }

type MauvrionSkiesYellow struct{}

func (MauvrionSkiesYellow) ID() card.ID                 { return card.MauvrionSkiesYellow }
func (MauvrionSkiesYellow) Name() string               { return "Mauvrion Skies (Yellow)" }
func (MauvrionSkiesYellow) Cost(*card.TurnState) int                  { return 0 }
func (MauvrionSkiesYellow) Pitch() int                 { return 2 }
func (MauvrionSkiesYellow) Attack() int                { return 0 }
func (MauvrionSkiesYellow) Defense() int               { return 2 }
func (MauvrionSkiesYellow) Types() card.TypeSet        { return mauvrionSkiesTypes }
func (MauvrionSkiesYellow) GoAgain() bool              { return true }
func (MauvrionSkiesYellow) Play(s *card.TurnState, _ *card.CardState) int { return mauvrionSkiesPlay(s, 2) }

type MauvrionSkiesBlue struct{}

func (MauvrionSkiesBlue) ID() card.ID                 { return card.MauvrionSkiesBlue }
func (MauvrionSkiesBlue) Name() string               { return "Mauvrion Skies (Blue)" }
func (MauvrionSkiesBlue) Cost(*card.TurnState) int                  { return 0 }
func (MauvrionSkiesBlue) Pitch() int                 { return 3 }
func (MauvrionSkiesBlue) Attack() int                { return 0 }
func (MauvrionSkiesBlue) Defense() int               { return 2 }
func (MauvrionSkiesBlue) Types() card.TypeSet        { return mauvrionSkiesTypes }
func (MauvrionSkiesBlue) GoAgain() bool              { return true }
func (MauvrionSkiesBlue) Play(s *card.TurnState, _ *card.CardState) int { return mauvrionSkiesPlay(s, 1) }

func mauvrionSkiesPlay(s *card.TurnState, n int) int {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if t.Has(card.TypeRuneblade) && t.Has(card.TypeAction) && t.Has(card.TypeAttack) {
			pc.GrantedGoAgain = true
			if card.LikelyToHit(pc.Card.Attack()) {
				return s.CreateRunechants(n)
			}
			// Target is blockable — the "if hits" clause doesn't fire, so no Runechants.
			return 0
		}
	}
	// No qualifying target — both the go-again grant and the runechant rider fizzle.
	return 0
}
