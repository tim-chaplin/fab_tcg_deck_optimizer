package sim

// ResolveChainStep runs card.Play and then applies the standard chain-step resolution:
// for an attack-action or weapon-attack, credit self.EffectiveAttack() to s.value; for
// a defense-reaction (or DefensiveInstant), credit the EffectiveDefense capped at
// s.incomingDamage and decrement IncomingDamage; for everything else, log (+0).
// The "<DisplayName>: <VERB> (+N)" chain-step entry is appended after Play returns so
// any self-buffs Play applied (e.g. modal +2{p} riders flipping self.BonusAttack) are
// reflected in the displayed delta.
//
// Cards' Play body owns card-specific behaviour: riders that emit rider log lines,
// OnHit registration, conditional self-buffs, sub-card plays. The standard
// printed-attack-deals-damage / DR-blocks-incoming mechanic is the sim's job; cards
// don't reach for DealEffectiveAttack / DealEffectiveDefense or emit the chain step
// themselves.
//
// The cards-facing rider lines (post-triggers emitted during Play via
// l.AppendPostTrigger or via DealArcaneDamage's internal LogRider) appear before the
// auto-logged chain step in the log stream. The format pass
// (appendGroupedChainEntries) buffers post-triggers and attaches them under the
// matching chain step regardless of stream order, so the rendered printout still reads
// "chain step, then indented riders" as before.
func ResolveChainStep(s *TurnState, l Logger, self *CardState) {
	self.Card.Play(s, l, self)
	n := chainStepDelta(s, self)
	l.AppendChainStep(ChainStepText(self), n)
}

// chainStepDelta computes the chain step's display delta and applies the standard
// damage / block side effects. Returns the (+N) value for the log line.
func chainStepDelta(s *TurnState, self *CardState) int {
	types := self.Card.Types()
	switch {
	case types.IsAttackAction() || types.IsWeaponAttack():
		n := self.EffectiveAttack()
		s.value += n
		return n
	case types.IsDefenseReaction() || isDefensiveInstant(self.Card):
		n := self.EffectiveDefense()
		if n > s.incomingDamage {
			n = s.incomingDamage
		}
		if n < 0 {
			n = 0
		}
		s.incomingDamage -= n
		s.value += n
		return n
	}
	return 0
}

// isDefensiveInstant reports whether c opts into the DR resolution path via the
// DefensiveInstant marker. Centralised here so ResolveChainStep doesn't repeat the
// type-assertion shape.
func isDefensiveInstant(c Card) bool {
	_, ok := c.(DefensiveInstant)
	return ok
}
