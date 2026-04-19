// Seek Horizon — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Seek Horizon, you may put a card from your hand on top of
// your deck. If you do, Seek Horizon gains **go again**."
//
// Simplification: Hand-on-top additional cost and conditional go-again aren't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var seekHorizonTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type SeekHorizonRed struct{}

func (SeekHorizonRed) ID() card.ID                 { return card.SeekHorizonRed }
func (SeekHorizonRed) Name() string                { return "Seek Horizon (Red)" }
func (SeekHorizonRed) Cost(*card.TurnState) int                   { return 0 }
func (SeekHorizonRed) Pitch() int                  { return 1 }
func (SeekHorizonRed) Attack() int                 { return 4 }
func (SeekHorizonRed) Defense() int                { return 2 }
func (SeekHorizonRed) Types() card.TypeSet         { return seekHorizonTypes }
func (SeekHorizonRed) GoAgain() bool               { return false }
func (c SeekHorizonRed) Play(s *card.TurnState) int { return c.Attack() }

type SeekHorizonYellow struct{}

func (SeekHorizonYellow) ID() card.ID                 { return card.SeekHorizonYellow }
func (SeekHorizonYellow) Name() string                { return "Seek Horizon (Yellow)" }
func (SeekHorizonYellow) Cost(*card.TurnState) int                   { return 0 }
func (SeekHorizonYellow) Pitch() int                  { return 2 }
func (SeekHorizonYellow) Attack() int                 { return 3 }
func (SeekHorizonYellow) Defense() int                { return 2 }
func (SeekHorizonYellow) Types() card.TypeSet         { return seekHorizonTypes }
func (SeekHorizonYellow) GoAgain() bool               { return false }
func (c SeekHorizonYellow) Play(s *card.TurnState) int { return c.Attack() }

type SeekHorizonBlue struct{}

func (SeekHorizonBlue) ID() card.ID                 { return card.SeekHorizonBlue }
func (SeekHorizonBlue) Name() string                { return "Seek Horizon (Blue)" }
func (SeekHorizonBlue) Cost(*card.TurnState) int                   { return 0 }
func (SeekHorizonBlue) Pitch() int                  { return 3 }
func (SeekHorizonBlue) Attack() int                 { return 2 }
func (SeekHorizonBlue) Defense() int                { return 2 }
func (SeekHorizonBlue) Types() card.TypeSet         { return seekHorizonTypes }
func (SeekHorizonBlue) GoAgain() bool               { return false }
func (c SeekHorizonBlue) Play(s *card.TurnState) int { return c.Attack() }
