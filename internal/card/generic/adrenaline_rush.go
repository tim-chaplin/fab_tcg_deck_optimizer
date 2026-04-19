// Adrenaline Rush — Generic Action - Attack. Cost 2. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When you play this, if you have less {h} than an opposing hero, this gets +3{p}."
//
// Simplification: 'Less life than opposing hero' health comparison isn't modelled; the +3{p} rider
// never fires.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var adrenalineRushTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type AdrenalineRushRed struct{}

func (AdrenalineRushRed) ID() card.ID                 { return card.AdrenalineRushRed }
func (AdrenalineRushRed) Name() string                { return "Adrenaline Rush (Red)" }
func (AdrenalineRushRed) Cost(*card.TurnState) int                   { return 2 }
func (AdrenalineRushRed) Pitch() int                  { return 1 }
func (AdrenalineRushRed) Attack() int                 { return 4 }
func (AdrenalineRushRed) Defense() int                { return 2 }
func (AdrenalineRushRed) Types() card.TypeSet         { return adrenalineRushTypes }
func (AdrenalineRushRed) GoAgain() bool               { return false }
func (c AdrenalineRushRed) Play(s *card.TurnState) int { return c.Attack() }

type AdrenalineRushYellow struct{}

func (AdrenalineRushYellow) ID() card.ID                 { return card.AdrenalineRushYellow }
func (AdrenalineRushYellow) Name() string                { return "Adrenaline Rush (Yellow)" }
func (AdrenalineRushYellow) Cost(*card.TurnState) int                   { return 2 }
func (AdrenalineRushYellow) Pitch() int                  { return 2 }
func (AdrenalineRushYellow) Attack() int                 { return 3 }
func (AdrenalineRushYellow) Defense() int                { return 2 }
func (AdrenalineRushYellow) Types() card.TypeSet         { return adrenalineRushTypes }
func (AdrenalineRushYellow) GoAgain() bool               { return false }
func (c AdrenalineRushYellow) Play(s *card.TurnState) int { return c.Attack() }

type AdrenalineRushBlue struct{}

func (AdrenalineRushBlue) ID() card.ID                 { return card.AdrenalineRushBlue }
func (AdrenalineRushBlue) Name() string                { return "Adrenaline Rush (Blue)" }
func (AdrenalineRushBlue) Cost(*card.TurnState) int                   { return 2 }
func (AdrenalineRushBlue) Pitch() int                  { return 3 }
func (AdrenalineRushBlue) Attack() int                 { return 2 }
func (AdrenalineRushBlue) Defense() int                { return 2 }
func (AdrenalineRushBlue) Types() card.TypeSet         { return adrenalineRushTypes }
func (AdrenalineRushBlue) GoAgain() bool               { return false }
func (c AdrenalineRushBlue) Play(s *card.TurnState) int { return c.Attack() }
