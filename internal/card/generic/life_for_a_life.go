// Life for a Life — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this is played, if you have less {h} than an opposing hero, it gets **go again**.
// When this hits, gain 1{h}."
//
// Simplification: Health comparison for go-again and on-hit 1{h} gain aren't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var lifeForALifeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type LifeForALifeRed struct{}

func (LifeForALifeRed) ID() card.ID                 { return card.LifeForALifeRed }
func (LifeForALifeRed) Name() string                { return "Life for a Life (Red)" }
func (LifeForALifeRed) Cost() int                   { return 1 }
func (LifeForALifeRed) Pitch() int                  { return 1 }
func (LifeForALifeRed) Attack() int                 { return 4 }
func (LifeForALifeRed) Defense() int                { return 2 }
func (LifeForALifeRed) Types() card.TypeSet         { return lifeForALifeTypes }
func (LifeForALifeRed) GoAgain() bool               { return false }
func (c LifeForALifeRed) Play(s *card.TurnState) int { return c.Attack() }

type LifeForALifeYellow struct{}

func (LifeForALifeYellow) ID() card.ID                 { return card.LifeForALifeYellow }
func (LifeForALifeYellow) Name() string                { return "Life for a Life (Yellow)" }
func (LifeForALifeYellow) Cost() int                   { return 1 }
func (LifeForALifeYellow) Pitch() int                  { return 2 }
func (LifeForALifeYellow) Attack() int                 { return 3 }
func (LifeForALifeYellow) Defense() int                { return 2 }
func (LifeForALifeYellow) Types() card.TypeSet         { return lifeForALifeTypes }
func (LifeForALifeYellow) GoAgain() bool               { return false }
func (c LifeForALifeYellow) Play(s *card.TurnState) int { return c.Attack() }

type LifeForALifeBlue struct{}

func (LifeForALifeBlue) ID() card.ID                 { return card.LifeForALifeBlue }
func (LifeForALifeBlue) Name() string                { return "Life for a Life (Blue)" }
func (LifeForALifeBlue) Cost() int                   { return 1 }
func (LifeForALifeBlue) Pitch() int                  { return 3 }
func (LifeForALifeBlue) Attack() int                 { return 2 }
func (LifeForALifeBlue) Defense() int                { return 2 }
func (LifeForALifeBlue) Types() card.TypeSet         { return lifeForALifeTypes }
func (LifeForALifeBlue) GoAgain() bool               { return false }
func (c LifeForALifeBlue) Play(s *card.TurnState) int { return c.Attack() }
