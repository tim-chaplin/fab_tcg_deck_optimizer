package sim

// Attack Reaction support. See docs/dev-standards.md "Attack Reactions".

// AttackReaction is implemented by every Attack Reaction card. The predicate matches the
// printed target text.
type AttackReaction interface {
	ARTargetAllowed(c Card) bool
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
