// Flying High — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "Your next attack this turn gets **go again**. If it's red, it gets +1{p}. **Go again**"
//
// Simplification: The '+1{p} if red' rider is modelled: when the granted target has pitch 1, +1 is
// also credited.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var flyingHighTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// flyingHighPlay grants go again to the next attack action card scheduled later this turn. For Flying High, if
// that target is red (pitch 1) we also credit +1 power as a bonus.
func flyingHighPlay(s *card.TurnState) int {
	for _, pc := range s.CardsRemaining {
		t := pc.Card.Types()
		if !t.Has(card.TypeAttack) || !t.Has(card.TypeAction) {
			continue
		}
		bonus := 0
		if pc.Card.Pitch() == 1 {
			bonus = 1
		}
		pc.GrantedGoAgain = true
		return bonus
	}
	return 0
}

type FlyingHighRed struct{}

func (FlyingHighRed) ID() card.ID                 { return card.FlyingHighRed }
func (FlyingHighRed) Name() string                { return "Flying High (Red)" }
func (FlyingHighRed) Cost() int                   { return 0 }
func (FlyingHighRed) Pitch() int                  { return 1 }
func (FlyingHighRed) Attack() int                 { return 0 }
func (FlyingHighRed) Defense() int                { return 2 }
func (FlyingHighRed) Types() card.TypeSet         { return flyingHighTypes }
func (FlyingHighRed) GoAgain() bool               { return true }
func (FlyingHighRed) GrantsGoAgain() bool         { return true }
func (FlyingHighRed) Play(s *card.TurnState) int { return flyingHighPlay(s) }

type FlyingHighYellow struct{}

func (FlyingHighYellow) ID() card.ID                 { return card.FlyingHighYellow }
func (FlyingHighYellow) Name() string                { return "Flying High (Yellow)" }
func (FlyingHighYellow) Cost() int                   { return 0 }
func (FlyingHighYellow) Pitch() int                  { return 2 }
func (FlyingHighYellow) Attack() int                 { return 0 }
func (FlyingHighYellow) Defense() int                { return 2 }
func (FlyingHighYellow) Types() card.TypeSet         { return flyingHighTypes }
func (FlyingHighYellow) GoAgain() bool               { return true }
func (FlyingHighYellow) GrantsGoAgain() bool         { return true }
func (FlyingHighYellow) Play(s *card.TurnState) int { return flyingHighPlay(s) }

type FlyingHighBlue struct{}

func (FlyingHighBlue) ID() card.ID                 { return card.FlyingHighBlue }
func (FlyingHighBlue) Name() string                { return "Flying High (Blue)" }
func (FlyingHighBlue) Cost() int                   { return 0 }
func (FlyingHighBlue) Pitch() int                  { return 3 }
func (FlyingHighBlue) Attack() int                 { return 0 }
func (FlyingHighBlue) Defense() int                { return 2 }
func (FlyingHighBlue) Types() card.TypeSet         { return flyingHighTypes }
func (FlyingHighBlue) GoAgain() bool               { return true }
func (FlyingHighBlue) GrantsGoAgain() bool         { return true }
func (FlyingHighBlue) Play(s *card.TurnState) int { return flyingHighPlay(s) }
