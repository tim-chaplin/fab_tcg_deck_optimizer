// On the Horizon — Generic Block. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Printed
// defense: Red 4, Yellow 3, Blue 2.
//
// Text: "When this defends, look at the top card of your deck."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

type OnTheHorizonRed struct{}

func (OnTheHorizonRed) ID() card.ID              { return card.OnTheHorizonRed }
func (OnTheHorizonRed) Name() string             { return "On the Horizon" }
func (OnTheHorizonRed) Cost(*card.TurnState) int { return 0 }
func (OnTheHorizonRed) Pitch() int               { return 1 }
func (OnTheHorizonRed) Attack() int              { return 0 }
func (OnTheHorizonRed) Defense() int             { return 4 }
func (OnTheHorizonRed) Types() card.TypeSet      { return defenseReactionTypes }
func (OnTheHorizonRed) GoAgain() bool            { return false }

// not implemented: deck-peek defend trigger
func (OnTheHorizonRed) NotImplemented()                              {}
func (OnTheHorizonRed) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }

type OnTheHorizonYellow struct{}

func (OnTheHorizonYellow) ID() card.ID              { return card.OnTheHorizonYellow }
func (OnTheHorizonYellow) Name() string             { return "On the Horizon" }
func (OnTheHorizonYellow) Cost(*card.TurnState) int { return 0 }
func (OnTheHorizonYellow) Pitch() int               { return 2 }
func (OnTheHorizonYellow) Attack() int              { return 0 }
func (OnTheHorizonYellow) Defense() int             { return 3 }
func (OnTheHorizonYellow) Types() card.TypeSet      { return defenseReactionTypes }
func (OnTheHorizonYellow) GoAgain() bool            { return false }

// not implemented: deck-peek defend trigger
func (OnTheHorizonYellow) NotImplemented()                              {}
func (OnTheHorizonYellow) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }

type OnTheHorizonBlue struct{}

func (OnTheHorizonBlue) ID() card.ID              { return card.OnTheHorizonBlue }
func (OnTheHorizonBlue) Name() string             { return "On the Horizon" }
func (OnTheHorizonBlue) Cost(*card.TurnState) int { return 0 }
func (OnTheHorizonBlue) Pitch() int               { return 3 }
func (OnTheHorizonBlue) Attack() int              { return 0 }
func (OnTheHorizonBlue) Defense() int             { return 2 }
func (OnTheHorizonBlue) Types() card.TypeSet      { return defenseReactionTypes }
func (OnTheHorizonBlue) GoAgain() bool            { return false }

// not implemented: deck-peek defend trigger
func (OnTheHorizonBlue) NotImplemented()                              {}
func (OnTheHorizonBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
