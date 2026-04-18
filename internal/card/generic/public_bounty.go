// Public Bounty — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 2.
//
// Text: "**Mark** target opposing hero. The next time you attack a **marked** hero this turn, the
// attack gets +N{p}. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: Mark isn't modelled; the +N rider is credited unconditionally. Scans
// TurnState.CardsRemaining for the first matching attack action card and credits the bonus assuming
// it will be played; if none is scheduled after this card, the bonus fizzles.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var publicBountyTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// publicBountyPlay returns n when a matching attack action card is scheduled later this turn.
func publicBountyPlay(s *card.TurnState, n int) int {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if !t.Has(card.TypeAttack) || !t.Has(card.TypeAction) {
			continue
		}
		return n
	}
	return 0
}

type PublicBountyRed struct{}

func (PublicBountyRed) ID() card.ID                 { return card.PublicBountyRed }
func (PublicBountyRed) Name() string                { return "Public Bounty (Red)" }
func (PublicBountyRed) Cost() int                   { return 1 }
func (PublicBountyRed) Pitch() int                  { return 1 }
func (PublicBountyRed) Attack() int                 { return 0 }
func (PublicBountyRed) Defense() int                { return 2 }
func (PublicBountyRed) Types() card.TypeSet         { return publicBountyTypes }
func (PublicBountyRed) GoAgain() bool               { return true }
func (PublicBountyRed) Play(s *card.TurnState) int { return publicBountyPlay(s, 3) }

type PublicBountyYellow struct{}

func (PublicBountyYellow) ID() card.ID                 { return card.PublicBountyYellow }
func (PublicBountyYellow) Name() string                { return "Public Bounty (Yellow)" }
func (PublicBountyYellow) Cost() int                   { return 1 }
func (PublicBountyYellow) Pitch() int                  { return 2 }
func (PublicBountyYellow) Attack() int                 { return 0 }
func (PublicBountyYellow) Defense() int                { return 2 }
func (PublicBountyYellow) Types() card.TypeSet         { return publicBountyTypes }
func (PublicBountyYellow) GoAgain() bool               { return true }
func (PublicBountyYellow) Play(s *card.TurnState) int { return publicBountyPlay(s, 2) }

type PublicBountyBlue struct{}

func (PublicBountyBlue) ID() card.ID                 { return card.PublicBountyBlue }
func (PublicBountyBlue) Name() string                { return "Public Bounty (Blue)" }
func (PublicBountyBlue) Cost() int                   { return 1 }
func (PublicBountyBlue) Pitch() int                  { return 3 }
func (PublicBountyBlue) Attack() int                 { return 0 }
func (PublicBountyBlue) Defense() int                { return 2 }
func (PublicBountyBlue) Types() card.TypeSet         { return publicBountyTypes }
func (PublicBountyBlue) GoAgain() bool               { return true }
func (PublicBountyBlue) Play(s *card.TurnState) int { return publicBountyPlay(s, 1) }
