// Come to Fight — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 3.
//
// Text: "The next attack action card you play this turn gains +N{p}. **Go again**" (Red N=3,
// Yellow N=2, Blue N=1.)
//
// Simplification: Scans TurnState.CardsRemaining for the first matching attack action card and
// credits the bonus assuming it will be played; if none is scheduled after this card, the bonus
// fizzles.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var comeToFightTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// comeToFightPlay returns n when a matching attack action card is scheduled later this turn.
func comeToFightPlay(s *card.TurnState, n int) int {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if !t.Has(card.TypeAttack) || !t.Has(card.TypeAction) {
			continue
		}
		return n
	}
	return 0
}

type ComeToFightRed struct{}

func (ComeToFightRed) ID() card.ID                 { return card.ComeToFightRed }
func (ComeToFightRed) Name() string                { return "Come to Fight (Red)" }
func (ComeToFightRed) Cost() int                   { return 1 }
func (ComeToFightRed) Pitch() int                  { return 1 }
func (ComeToFightRed) Attack() int                 { return 0 }
func (ComeToFightRed) Defense() int                { return 3 }
func (ComeToFightRed) Types() card.TypeSet         { return comeToFightTypes }
func (ComeToFightRed) GoAgain() bool               { return true }
func (ComeToFightRed) Play(s *card.TurnState) int { return comeToFightPlay(s, 3) }

type ComeToFightYellow struct{}

func (ComeToFightYellow) ID() card.ID                 { return card.ComeToFightYellow }
func (ComeToFightYellow) Name() string                { return "Come to Fight (Yellow)" }
func (ComeToFightYellow) Cost() int                   { return 1 }
func (ComeToFightYellow) Pitch() int                  { return 2 }
func (ComeToFightYellow) Attack() int                 { return 0 }
func (ComeToFightYellow) Defense() int                { return 3 }
func (ComeToFightYellow) Types() card.TypeSet         { return comeToFightTypes }
func (ComeToFightYellow) GoAgain() bool               { return true }
func (ComeToFightYellow) Play(s *card.TurnState) int { return comeToFightPlay(s, 2) }

type ComeToFightBlue struct{}

func (ComeToFightBlue) ID() card.ID                 { return card.ComeToFightBlue }
func (ComeToFightBlue) Name() string                { return "Come to Fight (Blue)" }
func (ComeToFightBlue) Cost() int                   { return 1 }
func (ComeToFightBlue) Pitch() int                  { return 3 }
func (ComeToFightBlue) Attack() int                 { return 0 }
func (ComeToFightBlue) Defense() int                { return 3 }
func (ComeToFightBlue) Types() card.TypeSet         { return comeToFightTypes }
func (ComeToFightBlue) GoAgain() bool               { return true }
func (ComeToFightBlue) Play(s *card.TurnState) int { return comeToFightPlay(s, 1) }
