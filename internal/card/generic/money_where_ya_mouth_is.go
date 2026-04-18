// Money Where Ya Mouth Is — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue
// 3. Defense 2.
//
// Text: "Your next attack this turn gets +N{p} and "When this attacks a hero, you may **wager** a
// Gold token with them."" (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: Wager Gold token rider is dropped. Scans TurnState.CardsRemaining for the first
// matching attack action card and credits the bonus assuming it will be played; if none is
// scheduled after this card, the bonus fizzles.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var moneyWhereYaMouthIsTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// moneyWhereYaMouthIsPlay returns n when a matching attack action card is scheduled later this turn.
func moneyWhereYaMouthIsPlay(s *card.TurnState, n int) int {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if !t.Has(card.TypeAttack) || !t.Has(card.TypeAction) {
			continue
		}
		return n
	}
	return 0
}

type MoneyWhereYaMouthIsRed struct{}

func (MoneyWhereYaMouthIsRed) ID() card.ID                 { return card.MoneyWhereYaMouthIsRed }
func (MoneyWhereYaMouthIsRed) Name() string                { return "Money Where Ya Mouth Is (Red)" }
func (MoneyWhereYaMouthIsRed) Cost() int                   { return 1 }
func (MoneyWhereYaMouthIsRed) Pitch() int                  { return 1 }
func (MoneyWhereYaMouthIsRed) Attack() int                 { return 0 }
func (MoneyWhereYaMouthIsRed) Defense() int                { return 2 }
func (MoneyWhereYaMouthIsRed) Types() card.TypeSet         { return moneyWhereYaMouthIsTypes }
func (MoneyWhereYaMouthIsRed) GoAgain() bool               { return true }
func (MoneyWhereYaMouthIsRed) Play(s *card.TurnState) int { return moneyWhereYaMouthIsPlay(s, 3) }

type MoneyWhereYaMouthIsYellow struct{}

func (MoneyWhereYaMouthIsYellow) ID() card.ID                 { return card.MoneyWhereYaMouthIsYellow }
func (MoneyWhereYaMouthIsYellow) Name() string                { return "Money Where Ya Mouth Is (Yellow)" }
func (MoneyWhereYaMouthIsYellow) Cost() int                   { return 1 }
func (MoneyWhereYaMouthIsYellow) Pitch() int                  { return 2 }
func (MoneyWhereYaMouthIsYellow) Attack() int                 { return 0 }
func (MoneyWhereYaMouthIsYellow) Defense() int                { return 2 }
func (MoneyWhereYaMouthIsYellow) Types() card.TypeSet         { return moneyWhereYaMouthIsTypes }
func (MoneyWhereYaMouthIsYellow) GoAgain() bool               { return true }
func (MoneyWhereYaMouthIsYellow) Play(s *card.TurnState) int { return moneyWhereYaMouthIsPlay(s, 2) }

type MoneyWhereYaMouthIsBlue struct{}

func (MoneyWhereYaMouthIsBlue) ID() card.ID                 { return card.MoneyWhereYaMouthIsBlue }
func (MoneyWhereYaMouthIsBlue) Name() string                { return "Money Where Ya Mouth Is (Blue)" }
func (MoneyWhereYaMouthIsBlue) Cost() int                   { return 1 }
func (MoneyWhereYaMouthIsBlue) Pitch() int                  { return 3 }
func (MoneyWhereYaMouthIsBlue) Attack() int                 { return 0 }
func (MoneyWhereYaMouthIsBlue) Defense() int                { return 2 }
func (MoneyWhereYaMouthIsBlue) Types() card.TypeSet         { return moneyWhereYaMouthIsTypes }
func (MoneyWhereYaMouthIsBlue) GoAgain() bool               { return true }
func (MoneyWhereYaMouthIsBlue) Play(s *card.TurnState) int { return moneyWhereYaMouthIsPlay(s, 1) }
