// Plunder Run — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next time an attack action card you control hits this turn, draw a card. If Plunder
// Run is played from arsenal, the next attack action card you play this turn gains +3{p}. **Go
// again**"
//
// Simplification: Draw rider on hit is dropped; the arsenal-only +3 is credited unconditionally.
// Scans TurnState.CardsRemaining for the first matching attack action card and credits the bonus
// assuming it will be played; if none is scheduled after this card, the bonus fizzles.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var plunderRunTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// plunderRunPlay returns 3 when a matching attack action card is scheduled later this turn.
func plunderRunPlay(s *card.TurnState) int {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if !t.Has(card.TypeAttack) || !t.Has(card.TypeAction) {
			continue
		}
		return 3
	}
	return 0
}

type PlunderRunRed struct{}

func (PlunderRunRed) ID() card.ID                 { return card.PlunderRunRed }
func (PlunderRunRed) Name() string                { return "Plunder Run (Red)" }
func (PlunderRunRed) Cost() int                   { return 0 }
func (PlunderRunRed) Pitch() int                  { return 1 }
func (PlunderRunRed) Attack() int                 { return 0 }
func (PlunderRunRed) Defense() int                { return 2 }
func (PlunderRunRed) Types() card.TypeSet         { return plunderRunTypes }
func (PlunderRunRed) GoAgain() bool               { return true }
func (PlunderRunRed) Play(s *card.TurnState) int { return plunderRunPlay(s) }

type PlunderRunYellow struct{}

func (PlunderRunYellow) ID() card.ID                 { return card.PlunderRunYellow }
func (PlunderRunYellow) Name() string                { return "Plunder Run (Yellow)" }
func (PlunderRunYellow) Cost() int                   { return 0 }
func (PlunderRunYellow) Pitch() int                  { return 2 }
func (PlunderRunYellow) Attack() int                 { return 0 }
func (PlunderRunYellow) Defense() int                { return 2 }
func (PlunderRunYellow) Types() card.TypeSet         { return plunderRunTypes }
func (PlunderRunYellow) GoAgain() bool               { return true }
func (PlunderRunYellow) Play(s *card.TurnState) int { return plunderRunPlay(s) }

type PlunderRunBlue struct{}

func (PlunderRunBlue) ID() card.ID                 { return card.PlunderRunBlue }
func (PlunderRunBlue) Name() string                { return "Plunder Run (Blue)" }
func (PlunderRunBlue) Cost() int                   { return 0 }
func (PlunderRunBlue) Pitch() int                  { return 3 }
func (PlunderRunBlue) Attack() int                 { return 0 }
func (PlunderRunBlue) Defense() int                { return 2 }
func (PlunderRunBlue) Types() card.TypeSet         { return plunderRunTypes }
func (PlunderRunBlue) GoAgain() bool               { return true }
func (PlunderRunBlue) Play(s *card.TurnState) int { return plunderRunPlay(s) }
