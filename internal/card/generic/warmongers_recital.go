// Warmonger's Recital — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next attack action card you play this turn gains +3{p} and "When this hits, put it on
// the bottom of its owner's deck." **Go again**"
//
// Simplification: The 'bottom of deck' rider is dropped (just credit the +3). Scans
// TurnState.CardsRemaining for the first matching attack action card and credits the bonus assuming
// it will be played; if none is scheduled after this card, the bonus fizzles.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var warmongersRecitalTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// warmongersRecitalPlay returns 3 when a matching attack action card is scheduled later this turn.
func warmongersRecitalPlay(s *card.TurnState) int {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if !t.Has(card.TypeAttack) || !t.Has(card.TypeAction) {
			continue
		}
		return 3
	}
	return 0
}

type WarmongersRecitalRed struct{}

func (WarmongersRecitalRed) ID() card.ID                 { return card.WarmongersRecitalRed }
func (WarmongersRecitalRed) Name() string                { return "Warmonger's Recital (Red)" }
func (WarmongersRecitalRed) Cost() int                   { return 1 }
func (WarmongersRecitalRed) Pitch() int                  { return 1 }
func (WarmongersRecitalRed) Attack() int                 { return 0 }
func (WarmongersRecitalRed) Defense() int                { return 2 }
func (WarmongersRecitalRed) Types() card.TypeSet         { return warmongersRecitalTypes }
func (WarmongersRecitalRed) GoAgain() bool               { return true }
func (WarmongersRecitalRed) Play(s *card.TurnState) int { return warmongersRecitalPlay(s) }

type WarmongersRecitalYellow struct{}

func (WarmongersRecitalYellow) ID() card.ID                 { return card.WarmongersRecitalYellow }
func (WarmongersRecitalYellow) Name() string                { return "Warmonger's Recital (Yellow)" }
func (WarmongersRecitalYellow) Cost() int                   { return 1 }
func (WarmongersRecitalYellow) Pitch() int                  { return 2 }
func (WarmongersRecitalYellow) Attack() int                 { return 0 }
func (WarmongersRecitalYellow) Defense() int                { return 2 }
func (WarmongersRecitalYellow) Types() card.TypeSet         { return warmongersRecitalTypes }
func (WarmongersRecitalYellow) GoAgain() bool               { return true }
func (WarmongersRecitalYellow) Play(s *card.TurnState) int { return warmongersRecitalPlay(s) }

type WarmongersRecitalBlue struct{}

func (WarmongersRecitalBlue) ID() card.ID                 { return card.WarmongersRecitalBlue }
func (WarmongersRecitalBlue) Name() string                { return "Warmonger's Recital (Blue)" }
func (WarmongersRecitalBlue) Cost() int                   { return 1 }
func (WarmongersRecitalBlue) Pitch() int                  { return 3 }
func (WarmongersRecitalBlue) Attack() int                 { return 0 }
func (WarmongersRecitalBlue) Defense() int                { return 2 }
func (WarmongersRecitalBlue) Types() card.TypeSet         { return warmongersRecitalTypes }
func (WarmongersRecitalBlue) GoAgain() bool               { return true }
func (WarmongersRecitalBlue) Play(s *card.TurnState) int { return warmongersRecitalPlay(s) }
