// Promise of Plenty — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Promise of Plenty hits, each hero who doesn't have a card in their arsenal puts the top
// card of their deck face down into their arsenal. If Promise of Plenty is played from arsenal, it
// gains **go again**."
//
// Modelling: The arsenal-placement rider isn't modelled (arsenal/deck content tracking would
// be required). The played-from-arsenal go-again fires via self.GrantedGoAgain when
// self.FromArsenal reports this copy came from the arsenal slot.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var promiseOfPlentyTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// promiseOfPlentyPlay grants self Go again when this copy was played from arsenal.
func promiseOfPlentyPlay(c card.Card, self *card.CardState) int {
	if self.FromArsenal {
		self.GrantedGoAgain = true
	}
	return c.Attack()
}

type PromiseOfPlentyRed struct{}

func (PromiseOfPlentyRed) ID() card.ID                 { return card.PromiseOfPlentyRed }
func (PromiseOfPlentyRed) Name() string                { return "Promise of Plenty (Red)" }
func (PromiseOfPlentyRed) Cost(*card.TurnState) int                   { return 0 }
func (PromiseOfPlentyRed) Pitch() int                  { return 1 }
func (PromiseOfPlentyRed) Attack() int                 { return 3 }
func (PromiseOfPlentyRed) Defense() int                { return 2 }
func (PromiseOfPlentyRed) Types() card.TypeSet         { return promiseOfPlentyTypes }
func (PromiseOfPlentyRed) GoAgain() bool               { return false }
func (c PromiseOfPlentyRed) Play(_ *card.TurnState, self *card.CardState) int { return promiseOfPlentyPlay(c, self) }

type PromiseOfPlentyYellow struct{}

func (PromiseOfPlentyYellow) ID() card.ID                 { return card.PromiseOfPlentyYellow }
func (PromiseOfPlentyYellow) Name() string                { return "Promise of Plenty (Yellow)" }
func (PromiseOfPlentyYellow) Cost(*card.TurnState) int                   { return 0 }
func (PromiseOfPlentyYellow) Pitch() int                  { return 2 }
func (PromiseOfPlentyYellow) Attack() int                 { return 2 }
func (PromiseOfPlentyYellow) Defense() int                { return 2 }
func (PromiseOfPlentyYellow) Types() card.TypeSet         { return promiseOfPlentyTypes }
func (PromiseOfPlentyYellow) GoAgain() bool               { return false }
func (c PromiseOfPlentyYellow) Play(_ *card.TurnState, self *card.CardState) int { return promiseOfPlentyPlay(c, self) }

type PromiseOfPlentyBlue struct{}

func (PromiseOfPlentyBlue) ID() card.ID                 { return card.PromiseOfPlentyBlue }
func (PromiseOfPlentyBlue) Name() string                { return "Promise of Plenty (Blue)" }
func (PromiseOfPlentyBlue) Cost(*card.TurnState) int                   { return 0 }
func (PromiseOfPlentyBlue) Pitch() int                  { return 3 }
func (PromiseOfPlentyBlue) Attack() int                 { return 1 }
func (PromiseOfPlentyBlue) Defense() int                { return 2 }
func (PromiseOfPlentyBlue) Types() card.TypeSet         { return promiseOfPlentyTypes }
func (PromiseOfPlentyBlue) GoAgain() bool               { return false }
func (c PromiseOfPlentyBlue) Play(_ *card.TurnState, self *card.CardState) int { return promiseOfPlentyPlay(c, self) }
