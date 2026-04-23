// Brutal Assault — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 3.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var brutalAssaultTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BrutalAssaultRed struct{}

func (BrutalAssaultRed) ID() card.ID                 { return card.BrutalAssaultRed }
func (BrutalAssaultRed) Name() string                { return "Brutal Assault (Red)" }
func (BrutalAssaultRed) Cost(*card.TurnState) int                   { return 2 }
func (BrutalAssaultRed) Pitch() int                  { return 1 }
func (BrutalAssaultRed) Attack() int                 { return 6 }
func (BrutalAssaultRed) Defense() int                { return 3 }
func (BrutalAssaultRed) Types() card.TypeSet         { return brutalAssaultTypes }
func (BrutalAssaultRed) GoAgain() bool               { return false }
func (c BrutalAssaultRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type BrutalAssaultYellow struct{}

func (BrutalAssaultYellow) ID() card.ID                 { return card.BrutalAssaultYellow }
func (BrutalAssaultYellow) Name() string                { return "Brutal Assault (Yellow)" }
func (BrutalAssaultYellow) Cost(*card.TurnState) int                   { return 2 }
func (BrutalAssaultYellow) Pitch() int                  { return 2 }
func (BrutalAssaultYellow) Attack() int                 { return 5 }
func (BrutalAssaultYellow) Defense() int                { return 3 }
func (BrutalAssaultYellow) Types() card.TypeSet         { return brutalAssaultTypes }
func (BrutalAssaultYellow) GoAgain() bool               { return false }
func (c BrutalAssaultYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type BrutalAssaultBlue struct{}

func (BrutalAssaultBlue) ID() card.ID                 { return card.BrutalAssaultBlue }
func (BrutalAssaultBlue) Name() string                { return "Brutal Assault (Blue)" }
func (BrutalAssaultBlue) Cost(*card.TurnState) int                   { return 2 }
func (BrutalAssaultBlue) Pitch() int                  { return 3 }
func (BrutalAssaultBlue) Attack() int                 { return 4 }
func (BrutalAssaultBlue) Defense() int                { return 3 }
func (BrutalAssaultBlue) Types() card.TypeSet         { return brutalAssaultTypes }
func (BrutalAssaultBlue) GoAgain() bool               { return false }
func (c BrutalAssaultBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
