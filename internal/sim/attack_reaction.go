package sim

// Attack Reaction support: ARs are 0-AP attack-step buff cards that target an attack on the
// stack. Optimizer plumbing is in two pieces:
//
//   1. Partition-level validity. An AR in the attack bag is illegal in FaB unless something
//      in the same chain (an attack action card or a swinging weapon, depending on the AR's
//      printed predicate) is a legal target. partitionHasValidARTargets is the bestSequence
//      pre-screen that rejects attack bags carrying an AR with no eligible target.
//
//   2. Buff application. ARs apply their +N{p} (or +1 GrantedGoAgain, +N{d}, …) to the
//      first matching CardState in TurnState.CardsRemaining when their Play resolves. Past
//      attacks have already credited damage and the buff can't retroactively change that;
//      orderings where the AR plays after every legal target naturally contribute zero,
//      and the permutation search's other orderings recover the buff's value.
//
// Cards opt in by implementing AttackReaction (the predicate the validator queries) and
// calling GrantAttackReactionBuff in their Play body to apply the actual rider.

// AttackReaction is implemented by every Attack Reaction card so the partition validator
// can ask "given this attack bag, does this AR have at least one legal target to react
// to?". The predicate matches the printed target text (e.g. Lunging Press's "target attack
// action card" rejects weapons; Thrust's "target sword attack" accepts both sword weapons
// and sword attack action cards). Returning true on a card that is itself the AR is fine —
// the validator excludes self-targets explicitly.
type AttackReaction interface {
	ARTargetAllowed(c Card) bool
}

// partitionHasValidARTargets reports whether every Attack Reaction in chain has at least
// one non-self chain card matching its target predicate. The chain slice must contain
// every card the chain will play this iteration — both partition attackers and the
// weapons selected by the active wmask — so that an AR which targets a weapon attack can
// see the weapon. Returns true when the chain has no ARs (the predicate is vacuously
// satisfied).
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
// target predicate. Self-targeting is excluded so an AR can't "target itself" — even if
// its own type set would technically match the predicate, FaB's targeting rules require a
// distinct target.
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

// GrantAttackReactionBuff finds the first card in s.CardsRemaining matching the predicate
// and adds n to its BonusAttack. Cards call this inside their Play body — see Lunging
// Press for the canonical shape. predicate is the AR's printed target text expressed as a
// Go function on Card; the helper iterates CardsRemaining in chain order so the buff
// lands on the next legal target, mirroring how a player would announce the buff target
// in a forward-resolving chain.
//
// When no remaining card matches, the buff fizzles silently. partitionHasValidARTargets
// guarantees a target exists somewhere in the chain, but the permutation search is free
// to place the AR after every legal target — those orderings contribute zero, and other
// orderings the search explores recover the value.
func GrantAttackReactionBuff(s *TurnState, predicate func(Card) bool, n int) {
	for _, pc := range s.CardsRemaining {
		if predicate(pc.Card) {
			pc.BonusAttack += n
			return
		}
	}
}
