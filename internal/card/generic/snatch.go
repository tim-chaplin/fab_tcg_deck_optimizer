// Snatch — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits, draw a card."
//
// The on-hit draw fires when card.LikelyToHit approves the printed attack.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var snatchTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// snatchPlay fires the on-hit draw when the attack is likely to land and returns its damage.
func snatchPlay(attack int, s *card.TurnState, self *card.CardState) int {
	if card.LikelyToHit(self) {
		s.DrawOne()
	}
	return attack
}

type SnatchRed struct{}

func (SnatchRed) ID() card.ID                  { return card.SnatchRed }
func (SnatchRed) Name() string                 { return "Snatch (Red)" }
func (SnatchRed) Cost(*card.TurnState) int                    { return 0 }
func (SnatchRed) Pitch() int                   { return 1 }
func (SnatchRed) Attack() int                  { return 4 }
func (SnatchRed) Defense() int                 { return 2 }
func (SnatchRed) Types() card.TypeSet          { return snatchTypes }
func (SnatchRed) GoAgain() bool                { return false }
func (SnatchRed) NoMemo()                      {} // on-hit DrawOne depends on top of deck
func (c SnatchRed) Play(s *card.TurnState, self *card.CardState) int { return snatchPlay(c.Attack(), s, self) }

type SnatchYellow struct{}

func (SnatchYellow) ID() card.ID                  { return card.SnatchYellow }
func (SnatchYellow) Name() string                 { return "Snatch (Yellow)" }
func (SnatchYellow) Cost(*card.TurnState) int                    { return 0 }
func (SnatchYellow) Pitch() int                   { return 2 }
func (SnatchYellow) Attack() int                  { return 3 }
func (SnatchYellow) Defense() int                 { return 2 }
func (SnatchYellow) Types() card.TypeSet          { return snatchTypes }
func (SnatchYellow) GoAgain() bool                { return false }
func (SnatchYellow) NoMemo()                      {}
func (c SnatchYellow) Play(s *card.TurnState, self *card.CardState) int { return snatchPlay(c.Attack(), s, self) }

type SnatchBlue struct{}

func (SnatchBlue) ID() card.ID                  { return card.SnatchBlue }
func (SnatchBlue) Name() string                 { return "Snatch (Blue)" }
func (SnatchBlue) Cost(*card.TurnState) int                    { return 0 }
func (SnatchBlue) Pitch() int                   { return 3 }
func (SnatchBlue) Attack() int                  { return 2 }
func (SnatchBlue) Defense() int                 { return 2 }
func (SnatchBlue) Types() card.TypeSet          { return snatchTypes }
func (SnatchBlue) GoAgain() bool                { return false }
func (SnatchBlue) NoMemo()                      {}
func (c SnatchBlue) Play(s *card.TurnState, self *card.CardState) int { return snatchPlay(c.Attack(), s, self) }
