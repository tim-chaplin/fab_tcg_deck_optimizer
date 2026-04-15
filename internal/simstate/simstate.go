// Package simstate holds process-wide simulation state that card effects read. It's a leaf
// package (depends only on internal/card) so any card implementation can import it without
// triggering cycles through the hand/deck/cards stack.
package simstate

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// CurrentHero is the hero playing the current simulation. Set once at the start of a run (the
// hero doesn't change mid-program) so card effects can read profile info like Intelligence
// without plumbing it through TurnState.
var CurrentHero card.Hero
