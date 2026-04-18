// Blow for a Blow — Generic Action - Attack. Cost 2, Pitch 1, Power 4, Defense 2. Only printed in
// Red.
//
// Text: "When this is played, if you have less {h} than an opposing hero, it gets **go again**.
// When this hits, deal 1 damage to any target."
//
// Simplification: Health comparison for go-again and on-hit 1 damage aren't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var blowForABlowTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BlowForABlowRed struct{}

func (BlowForABlowRed) ID() card.ID                 { return card.BlowForABlowRed }
func (BlowForABlowRed) Name() string                { return "Blow for a Blow (Red)" }
func (BlowForABlowRed) Cost() int                   { return 2 }
func (BlowForABlowRed) Pitch() int                  { return 1 }
func (BlowForABlowRed) Attack() int                 { return 4 }
func (BlowForABlowRed) Defense() int                { return 2 }
func (BlowForABlowRed) Types() card.TypeSet         { return blowForABlowTypes }
func (BlowForABlowRed) GoAgain() bool               { return false }
func (c BlowForABlowRed) Play(s *card.TurnState) int { return c.Attack() }
