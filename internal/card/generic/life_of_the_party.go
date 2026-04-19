// Life of the Party — Generic Action - Attack. Cost 2. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "You may discard or destroy a card you control named Crazy Brew rather than pay Life of the
// Party's {r} cost. If you do, choose all modes, otherwise choose 1 at random; - This gets "When
// this hits, gain life 2{h}." - This gets +2{p}. - This gets **go again**."
//
// Simplification: Crazy Brew substitute and random-mode selection aren't modelled, so all three
// modes (+2{p}, on-hit life, go again) default off — including go again. Returning true from
// GoAgain() would make the chain-legality check always pass, over-crediting sequences against
// the baseline where the random mode rolled a different option.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var lifeOfThePartyTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type LifeOfThePartyRed struct{}

func (LifeOfThePartyRed) ID() card.ID                 { return card.LifeOfThePartyRed }
func (LifeOfThePartyRed) Name() string                { return "Life of the Party (Red)" }
func (LifeOfThePartyRed) Cost(*card.TurnState) int                   { return 2 }
func (LifeOfThePartyRed) Pitch() int                  { return 1 }
func (LifeOfThePartyRed) Attack() int                 { return 4 }
func (LifeOfThePartyRed) Defense() int                { return 2 }
func (LifeOfThePartyRed) Types() card.TypeSet         { return lifeOfThePartyTypes }
func (LifeOfThePartyRed) GoAgain() bool               { return false }
func (c LifeOfThePartyRed) Play(s *card.TurnState) int { return c.Attack() }

type LifeOfThePartyYellow struct{}

func (LifeOfThePartyYellow) ID() card.ID                 { return card.LifeOfThePartyYellow }
func (LifeOfThePartyYellow) Name() string                { return "Life of the Party (Yellow)" }
func (LifeOfThePartyYellow) Cost(*card.TurnState) int                   { return 2 }
func (LifeOfThePartyYellow) Pitch() int                  { return 2 }
func (LifeOfThePartyYellow) Attack() int                 { return 3 }
func (LifeOfThePartyYellow) Defense() int                { return 2 }
func (LifeOfThePartyYellow) Types() card.TypeSet         { return lifeOfThePartyTypes }
func (LifeOfThePartyYellow) GoAgain() bool               { return false }
func (c LifeOfThePartyYellow) Play(s *card.TurnState) int { return c.Attack() }

type LifeOfThePartyBlue struct{}

func (LifeOfThePartyBlue) ID() card.ID                 { return card.LifeOfThePartyBlue }
func (LifeOfThePartyBlue) Name() string                { return "Life of the Party (Blue)" }
func (LifeOfThePartyBlue) Cost(*card.TurnState) int                   { return 2 }
func (LifeOfThePartyBlue) Pitch() int                  { return 3 }
func (LifeOfThePartyBlue) Attack() int                 { return 2 }
func (LifeOfThePartyBlue) Defense() int                { return 2 }
func (LifeOfThePartyBlue) Types() card.TypeSet         { return lifeOfThePartyTypes }
func (LifeOfThePartyBlue) GoAgain() bool               { return false }
func (c LifeOfThePartyBlue) Play(s *card.TurnState) int { return c.Attack() }
