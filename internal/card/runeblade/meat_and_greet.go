// Meat and Greet — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When this hits, create a Runechant token. If you've dealt arcane damage to an opposing
// hero this turn, this gets go again."
//
// Modelling:
//   - The on-hit Runechant fires only when the attack's printed power satisfies
//     card.LikelyToHit; when it does, the token is credited as +1 damage via CreateRunechant
//     (which also sets AuraCreated). Blockable-power variants drop the rider.
//   - The go-again rider reads TurnState.ArcaneDamageDealt; when live, sets SelfGoAgain so
//     the solver's chain-legality check sees the conditional. The card's own Runechant fires
//     on a future turn, so it can't satisfy its own rider.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var meatAndGreetTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// meatAndGreetPlay is the shared Play implementation. Go again goes on SelfGoAgain (not the
// printed GoAgain) so the rider stays conditional on ArcaneDamageDealt. The on-hit Runechant
// rider is gated on card.LikelyToHit, mirroring how other on-hit rider cards treat blockable
// attacks.
func meatAndGreetPlay(c card.Card, s *card.TurnState) int {
	if s.ArcaneDamageDealt {
		s.SelfGoAgain = true
	}
	if card.LikelyToHit(c.Attack()) {
		return c.Attack() + s.CreateRunechant()
	}
	return c.Attack()
}

type MeatAndGreetRed struct{}

func (MeatAndGreetRed) ID() card.ID                   { return card.MeatAndGreetRed }
func (MeatAndGreetRed) Name() string                  { return "Meat and Greet (Red)" }
func (MeatAndGreetRed) Cost(*card.TurnState) int                     { return 1 }
func (MeatAndGreetRed) Pitch() int                    { return 1 }
func (MeatAndGreetRed) Attack() int                   { return 4 }
func (MeatAndGreetRed) Defense() int                  { return 3 }
func (MeatAndGreetRed) Types() card.TypeSet           { return meatAndGreetTypes }
func (MeatAndGreetRed) GoAgain() bool                 { return false }
func (c MeatAndGreetRed) Play(s *card.TurnState) int  { return meatAndGreetPlay(c, s) }

type MeatAndGreetYellow struct{}

func (MeatAndGreetYellow) ID() card.ID                   { return card.MeatAndGreetYellow }
func (MeatAndGreetYellow) Name() string                  { return "Meat and Greet (Yellow)" }
func (MeatAndGreetYellow) Cost(*card.TurnState) int                     { return 1 }
func (MeatAndGreetYellow) Pitch() int                    { return 2 }
func (MeatAndGreetYellow) Attack() int                   { return 3 }
func (MeatAndGreetYellow) Defense() int                  { return 3 }
func (MeatAndGreetYellow) Types() card.TypeSet           { return meatAndGreetTypes }
func (MeatAndGreetYellow) GoAgain() bool                 { return false }
func (c MeatAndGreetYellow) Play(s *card.TurnState) int  { return meatAndGreetPlay(c, s) }

type MeatAndGreetBlue struct{}

func (MeatAndGreetBlue) ID() card.ID                   { return card.MeatAndGreetBlue }
func (MeatAndGreetBlue) Name() string                  { return "Meat and Greet (Blue)" }
func (MeatAndGreetBlue) Cost(*card.TurnState) int                     { return 1 }
func (MeatAndGreetBlue) Pitch() int                    { return 3 }
func (MeatAndGreetBlue) Attack() int                   { return 2 }
func (MeatAndGreetBlue) Defense() int                  { return 3 }
func (MeatAndGreetBlue) Types() card.TypeSet           { return meatAndGreetTypes }
func (MeatAndGreetBlue) GoAgain() bool                 { return false }
func (c MeatAndGreetBlue) Play(s *card.TurnState) int  { return meatAndGreetPlay(c, s) }
