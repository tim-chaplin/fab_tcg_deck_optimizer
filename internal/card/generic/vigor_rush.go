// Vigor Rush — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If you have played a 'non-attack' action card this turn, Vigor Rush gains **go again**."
//
// Go again is conditional on a prior non-attack action, so GoAgain() returns false and
// vigorRushPlay sets SelfGoAgain when the condition fires. Returning true from GoAgain()
// unconditionally would over-credit sequences with no non-attack action played.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var vigorRushTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// vigorRushPlay grants go again when any non-attack Action has been played earlier this turn.
func vigorRushPlay(base int, s *card.TurnState, self *card.PlayedCard) int {
	for _, pl := range s.CardsPlayed {
		if pl.Types().IsNonAttackAction() {
			self.GrantedGoAgain = true
			break
		}
	}
	return base
}

type VigorRushRed struct{}

func (VigorRushRed) ID() card.ID                 { return card.VigorRushRed }
func (VigorRushRed) Name() string                { return "Vigor Rush (Red)" }
func (VigorRushRed) Cost(*card.TurnState) int                   { return 1 }
func (VigorRushRed) Pitch() int                  { return 1 }
func (VigorRushRed) Attack() int                 { return 4 }
func (VigorRushRed) Defense() int                { return 2 }
func (VigorRushRed) Types() card.TypeSet         { return vigorRushTypes }
func (VigorRushRed) GoAgain() bool               { return false }
func (c VigorRushRed) Play(s *card.TurnState, self *card.PlayedCard) int { return vigorRushPlay(c.Attack(), s, self) }

type VigorRushYellow struct{}

func (VigorRushYellow) ID() card.ID                 { return card.VigorRushYellow }
func (VigorRushYellow) Name() string                { return "Vigor Rush (Yellow)" }
func (VigorRushYellow) Cost(*card.TurnState) int                   { return 1 }
func (VigorRushYellow) Pitch() int                  { return 2 }
func (VigorRushYellow) Attack() int                 { return 3 }
func (VigorRushYellow) Defense() int                { return 2 }
func (VigorRushYellow) Types() card.TypeSet         { return vigorRushTypes }
func (VigorRushYellow) GoAgain() bool               { return false }
func (c VigorRushYellow) Play(s *card.TurnState, self *card.PlayedCard) int { return vigorRushPlay(c.Attack(), s, self) }

type VigorRushBlue struct{}

func (VigorRushBlue) ID() card.ID                 { return card.VigorRushBlue }
func (VigorRushBlue) Name() string                { return "Vigor Rush (Blue)" }
func (VigorRushBlue) Cost(*card.TurnState) int                   { return 1 }
func (VigorRushBlue) Pitch() int                  { return 3 }
func (VigorRushBlue) Attack() int                 { return 2 }
func (VigorRushBlue) Defense() int                { return 2 }
func (VigorRushBlue) Types() card.TypeSet         { return vigorRushTypes }
func (VigorRushBlue) GoAgain() bool               { return false }
func (c VigorRushBlue) Play(s *card.TurnState, self *card.PlayedCard) int { return vigorRushPlay(c.Attack(), s, self) }
