package sim

// Attack Reaction support. See docs/dev-standards.md "Attack Reactions".

// AttackReaction is implemented by every Attack Reaction card. ARTargetAllowed reports
// whether c is a legal target for this AR's chosen mode. Non-modal ARs ignore the mode
// parameter (it's always 0). Modal ARs (sim.ModalCard) dispatch on it: each mode's printed
// target text becomes its own predicate leg, and the chain runner rejects the permutation
// when the chosen mode doesn't accept the active attack.
type AttackReaction interface {
	ARTargetAllowed(c Card, mode int8) bool
}

// GrantAttackReactionBuff buffs the active attack target by n: adds to BonusAttack, credits
// s.Value, amends the target's chain-step delta, and logs the rider under the target's
// entry. Cards call this from Play; the chain runner has already validated the target.
func GrantAttackReactionBuff(s *TurnState, self *CardState, n int) {
	target := s.AttackReactionTarget()
	if target == nil {
		return
	}
	target.BonusAttack += n
	s.AddValue(n)
	s.AmendLastChainStepN(n)
	// N=0: the +n delta is folded into the parent chain step via AmendLastChainStepN.
	s.LogPostTriggerf(DisplayName(target.Card), 0, "%s buffed +%d{p}", DisplayName(self.Card), n)
}
