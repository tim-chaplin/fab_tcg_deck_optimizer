// Regain Composure — Generic Action. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "Your next attack this turn gets +1{p} and "When this hits, {u} your hero." **Go again**"
//
// Simplification: The on-hit unfreeze rider is dropped. Scans TurnState.CardsRemaining for the
// first matching attack action card and credits the bonus assuming it will be played; if none is
// scheduled after this card, the bonus fizzles.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var regainComposureTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// regainComposurePlay returns 1 when a matching attack action card is scheduled later this turn.
func regainComposurePlay(s *card.TurnState) int {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if !t.Has(card.TypeAttack) || !t.Has(card.TypeAction) {
			continue
		}
		return 1
	}
	return 0
}

type RegainComposureBlue struct{}

func (RegainComposureBlue) ID() card.ID                 { return card.RegainComposureBlue }
func (RegainComposureBlue) Name() string                { return "Regain Composure (Blue)" }
func (RegainComposureBlue) Cost() int                   { return 0 }
func (RegainComposureBlue) Pitch() int                  { return 3 }
func (RegainComposureBlue) Attack() int                 { return 0 }
func (RegainComposureBlue) Defense() int                { return 2 }
func (RegainComposureBlue) Types() card.TypeSet         { return regainComposureTypes }
func (RegainComposureBlue) GoAgain() bool               { return true }
func (RegainComposureBlue) Play(s *card.TurnState) int { return regainComposurePlay(s) }
