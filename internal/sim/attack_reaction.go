package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Attack Reaction support. See docs/dev-standards.md "Attack Reactions".

// AttackReaction is implemented by every Attack Reaction card. ARTargetAllowed reports
// whether c is a legal target for this AR's chosen mode. Non-modal ARs ignore the mode
// parameter (it's always 0). Modal ARs (sim.ModalCard) dispatch on it: each mode's printed
// target text becomes its own predicate leg, and the chain runner rejects the permutation
// when the chosen mode doesn't accept the active attack.
//
// card.GrantAttackReactionBuff (the helper most ARs call from Play) lives in v2/card —
// it's pure GameEngine / Logger / CardState plumbing.
type AttackReaction interface {
	ARTargetAllowed(c card.Card, mode int8) bool
}
