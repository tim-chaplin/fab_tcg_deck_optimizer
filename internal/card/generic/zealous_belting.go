// Zealous Belting — Generic Action - Attack. Cost 2. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "While there is a card in your pitch zone with {p} greater than Zealous Belting's base {p},
// Zealous Belting has **go again**."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var zealousBeltingTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// zealousBeltingPlay grants go again when any pitched card this turn has base power greater
// than the card's own base power.
func zealousBeltingPlay(base int, s *card.TurnState, self *card.CardState) int {
	for _, p := range s.Pitched {
		if p.Attack() > base {
			self.GrantedGoAgain = true
			break
		}
	}
	return base
}

type ZealousBeltingRed struct{}

func (ZealousBeltingRed) ID() card.ID                 { return card.ZealousBeltingRed }
func (ZealousBeltingRed) Name() string                { return "Zealous Belting (Red)" }
func (ZealousBeltingRed) Cost(*card.TurnState) int                   { return 2 }
func (ZealousBeltingRed) Pitch() int                  { return 1 }
func (ZealousBeltingRed) Attack() int                 { return 5 }
func (ZealousBeltingRed) Defense() int                { return 2 }
func (ZealousBeltingRed) Types() card.TypeSet         { return zealousBeltingTypes }
func (ZealousBeltingRed) GoAgain() bool               { return false }
func (c ZealousBeltingRed) Play(s *card.TurnState, self *card.CardState) int { return zealousBeltingPlay(c.Attack(), s, self) }

type ZealousBeltingYellow struct{}

func (ZealousBeltingYellow) ID() card.ID                 { return card.ZealousBeltingYellow }
func (ZealousBeltingYellow) Name() string                { return "Zealous Belting (Yellow)" }
func (ZealousBeltingYellow) Cost(*card.TurnState) int                   { return 2 }
func (ZealousBeltingYellow) Pitch() int                  { return 2 }
func (ZealousBeltingYellow) Attack() int                 { return 4 }
func (ZealousBeltingYellow) Defense() int                { return 2 }
func (ZealousBeltingYellow) Types() card.TypeSet         { return zealousBeltingTypes }
func (ZealousBeltingYellow) GoAgain() bool               { return false }
func (c ZealousBeltingYellow) Play(s *card.TurnState, self *card.CardState) int { return zealousBeltingPlay(c.Attack(), s, self) }

type ZealousBeltingBlue struct{}

func (ZealousBeltingBlue) ID() card.ID                 { return card.ZealousBeltingBlue }
func (ZealousBeltingBlue) Name() string                { return "Zealous Belting (Blue)" }
func (ZealousBeltingBlue) Cost(*card.TurnState) int                   { return 2 }
func (ZealousBeltingBlue) Pitch() int                  { return 3 }
func (ZealousBeltingBlue) Attack() int                 { return 3 }
func (ZealousBeltingBlue) Defense() int                { return 2 }
func (ZealousBeltingBlue) Types() card.TypeSet         { return zealousBeltingTypes }
func (ZealousBeltingBlue) GoAgain() bool               { return false }
func (c ZealousBeltingBlue) Play(s *card.TurnState, self *card.CardState) int { return zealousBeltingPlay(c.Attack(), s, self) }
