// Package simstate holds process-wide simulation state that card effects read. Leaf package
// (only depends on internal/card) so any card implementation can import it without cycling
// through the hand/deck/cards stack.
package simstate

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// CurrentHero is the hero playing the current simulation. Set once at the start of a run; card
// effects read profile info like Intelligence without plumbing it through TurnState.
var CurrentHero card.Hero
