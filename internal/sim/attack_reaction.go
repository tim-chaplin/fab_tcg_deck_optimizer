package sim

// Attack Reaction support. See docs/dev-standards.md "Attack Reactions" for the wiring
// contract; this file is the framework half (validator + buff helper).

// AttackReaction is implemented by every Attack Reaction card. The predicate matches the
// printed target text. The validator excludes self-targets, so a predicate matching the AR's
// own type set is harmless.
type AttackReaction interface {
	ARTargetAllowed(c Card) bool
}

// partitionHasValidARTargets reports whether every Attack Reaction in chain has at least one
// non-self chain card matching its target predicate. The chain slice must include both
// partition attackers and the weapons selected by the active wmask. Vacuously true when no
// ARs are present.
func partitionHasValidARTargets(chain []Card) bool {
	for i, c := range chain {
		ar, ok := c.(AttackReaction)
		if !ok {
			continue
		}
		if !arHasAnyTarget(ar, chain, i) {
			return false
		}
	}
	return true
}

// arHasAnyTarget reports whether any chain card other than chain[selfIdx] satisfies ar's
// predicate. FaB's targeting rules require a distinct target, so self is always skipped.
func arHasAnyTarget(ar AttackReaction, chain []Card, selfIdx int) bool {
	for j, target := range chain {
		if j == selfIdx {
			continue
		}
		if ar.ARTargetAllowed(target) {
			return true
		}
	}
	return false
}

// GrantAttackReactionBuff adds n to BonusAttack on the first card in s.CardsRemaining
// matching predicate. Fizzles silently when no remaining card matches —
// partitionHasValidARTargets guarantees a target exists somewhere in the chain, but
// orderings where the AR plays after every legal target naturally contribute zero and other
// orderings recover the value.
func GrantAttackReactionBuff(s *TurnState, predicate func(Card) bool, n int) {
	for _, pc := range s.CardsRemaining {
		if predicate(pc.Card) {
			pc.BonusAttack += n
			return
		}
	}
}
