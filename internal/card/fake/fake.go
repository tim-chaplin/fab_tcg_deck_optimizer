// Package fake provides generic stub Card implementations used by tests in multiple packages (hand,
// sim). These are not real FaB cards — they're deliberately simple attack actions with known stat
// lines so partition/ordering tests have predictable optimal values.
package fake

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var genericAttackTypes = map[string]bool{"Generic": true, "Action": true, "Attack": true}

// Blue is a generic blue action: pitches 3, defends 3, attacks 1, costs 1.
type Blue struct{}

func (Blue) Name() string                { return "cardtest.Blue" }
func (Blue) Cost() int                   { return 1 }
func (Blue) Pitch() int                  { return 3 }
func (Blue) Attack() int                 { return 1 }
func (Blue) Defense() int                { return 3 }
func (Blue) Types() map[string]bool      { return genericAttackTypes }
func (Blue) GoAgain() bool               { return true }
func (c Blue) Play(*card.TurnState) int  { return c.Attack() }

// Red is a generic red action: pitches 1, defends 1, attacks 3, costs 1.
type Red struct{}

func (Red) Name() string                { return "cardtest.Red" }
func (Red) Cost() int                   { return 1 }
func (Red) Pitch() int                  { return 1 }
func (Red) Attack() int                 { return 3 }
func (Red) Defense() int                { return 1 }
func (Red) Types() map[string]bool      { return genericAttackTypes }
func (Red) GoAgain() bool               { return true }
func (c Red) Play(*card.TurnState) int  { return c.Attack() }
