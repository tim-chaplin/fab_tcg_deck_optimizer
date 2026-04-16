// Package fake provides generic stub Card implementations used by tests in multiple packages (hand,
// sim). These are not real FaB cards — they're deliberately simple attack actions with known stat
// lines so partition/ordering tests have predictable optimal values.
package fake

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var genericAttackTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// BlueAttack is a generic blue attack action: pitches 3, defends 3, attacks 1, costs 1.
type BlueAttack struct{}

func (BlueAttack) ID() card.ID                { return card.FakeBlueAttack }
func (BlueAttack) Name() string               { return "cardtest.BlueAttack" }
func (BlueAttack) Cost() int                  { return 1 }
func (BlueAttack) Pitch() int                 { return 3 }
func (BlueAttack) Attack() int                { return 1 }
func (BlueAttack) Defense() int               { return 3 }
func (BlueAttack) Types() card.TypeSet        { return genericAttackTypes }
func (BlueAttack) GoAgain() bool              { return true }
func (c BlueAttack) Play(*card.TurnState) int { return c.Attack() }

// RedAttack is a generic red attack action: pitches 1, defends 1, attacks 3, costs 1.
type RedAttack struct{}

func (RedAttack) ID() card.ID                { return card.FakeRedAttack }
func (RedAttack) Name() string               { return "cardtest.RedAttack" }
func (RedAttack) Cost() int                  { return 1 }
func (RedAttack) Pitch() int                 { return 1 }
func (RedAttack) Attack() int                { return 3 }
func (RedAttack) Defense() int               { return 1 }
func (RedAttack) Types() card.TypeSet        { return genericAttackTypes }
func (RedAttack) GoAgain() bool              { return true }
func (c RedAttack) Play(*card.TurnState) int { return c.Attack() }
