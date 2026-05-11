package card

// GrantAttackReactionBuff buffs the active attack target by n: adds to BonusAttack, credits
// g's value, amends the target's chain-step delta, and logs the rider under the target's
// entry. Cards call this from Play; the chain runner has already validated the target.
func GrantAttackReactionBuff(g GameEngine, l Logger, self *CardState, n int) {
	target := g.AttackReactionTarget()
	if target == nil {
		return
	}
	target.BonusAttack += n
	g.AddValue(n)
	l.AmendLastChainStepN(n)
	// N=0: the +n delta is folded into the parent chain step via AmendLastChainStepN.
	l.AppendPostTriggerf(target.Card.DisplayName(), 0, "%s buffed +%d{p}", self.Card.DisplayName(), n)
}
