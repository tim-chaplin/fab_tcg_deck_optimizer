// Meat and Greet — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When this hits, create a Runechant token. If you've dealt arcane damage to an opposing
// hero this turn, this gets go again."
//
// On-hit Runechant fires only when the attack's printed power satisfies card.LikelyToHit;
// blockable variants drop the rider. Go-again is conditional on TurnState.ArcaneDamageDealt.
// The card's own Runechant fires on a future turn, so it can't satisfy its own rider.

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var meatAndGreetTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// meatAndGreetPlay is the shared Play implementation. See the file docstring for rider
// modelling.
func meatAndGreetPlay(s *card.TurnState, self *card.CardState) {
	if s.ArcaneDamageDealt {
		self.GrantedGoAgain = true
	}
	s.ApplyAndLogEffectiveAttack(self)
	if card.LikelyToHit(self) {
		s.ApplyAndLogRiderOnPlay(self, "On-hit created a runechant", s.CreateRunechant())
	}
}

type MeatAndGreetRed struct{}

func (MeatAndGreetRed) ID() card.ID              { return card.MeatAndGreetRed }
func (MeatAndGreetRed) Name() string             { return "Meat and Greet" }
func (MeatAndGreetRed) Cost(*card.TurnState) int { return 1 }
func (MeatAndGreetRed) Pitch() int               { return 1 }
func (MeatAndGreetRed) Attack() int              { return 4 }
func (MeatAndGreetRed) Defense() int             { return 3 }
func (MeatAndGreetRed) Types() card.TypeSet      { return meatAndGreetTypes }
func (MeatAndGreetRed) GoAgain() bool            { return false }
func (MeatAndGreetRed) Play(s *card.TurnState, self *card.CardState) {
	meatAndGreetPlay(s, self)
}

type MeatAndGreetYellow struct{}

func (MeatAndGreetYellow) ID() card.ID              { return card.MeatAndGreetYellow }
func (MeatAndGreetYellow) Name() string             { return "Meat and Greet" }
func (MeatAndGreetYellow) Cost(*card.TurnState) int { return 1 }
func (MeatAndGreetYellow) Pitch() int               { return 2 }
func (MeatAndGreetYellow) Attack() int              { return 3 }
func (MeatAndGreetYellow) Defense() int             { return 3 }
func (MeatAndGreetYellow) Types() card.TypeSet      { return meatAndGreetTypes }
func (MeatAndGreetYellow) GoAgain() bool            { return false }
func (MeatAndGreetYellow) Play(s *card.TurnState, self *card.CardState) {
	meatAndGreetPlay(s, self)
}

type MeatAndGreetBlue struct{}

func (MeatAndGreetBlue) ID() card.ID              { return card.MeatAndGreetBlue }
func (MeatAndGreetBlue) Name() string             { return "Meat and Greet" }
func (MeatAndGreetBlue) Cost(*card.TurnState) int { return 1 }
func (MeatAndGreetBlue) Pitch() int               { return 3 }
func (MeatAndGreetBlue) Attack() int              { return 2 }
func (MeatAndGreetBlue) Defense() int             { return 3 }
func (MeatAndGreetBlue) Types() card.TypeSet      { return meatAndGreetTypes }
func (MeatAndGreetBlue) GoAgain() bool            { return false }
func (MeatAndGreetBlue) Play(s *card.TurnState, self *card.CardState) {
	meatAndGreetPlay(s, self)
}
