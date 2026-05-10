package sim

// ChainStepText returns the "<DisplayName>: <VERB>[ from arsenal]" prefix the chain-step
// log line is built from. VERB picks WEAPON ATTACK for weapon activated-ability cards
// (Weapon + Attack), ATTACK for attack-action cards, DEFENSE REACTION for Defense
// Reactions, and PLAY for everything else; the "from arsenal" suffix tags entries played
// out of the arsenal slot. Declared as a var so optimizations can swap in a memoised
// per-(CardID, FromArsenal) implementation at init.
var ChainStepText = func(self *CardState) string {
	types := self.Card.Types()
	var verb string
	switch {
	case types.IsWeaponAttack():
		verb = "WEAPON ATTACK"
	case types.IsAttackAction():
		verb = "ATTACK"
	case types.IsDefenseReaction():
		verb = "DEFENSE REACTION"
	default:
		verb = "PLAY"
	}
	if self.FromArsenal {
		verb += " from arsenal"
	}
	return self.Card.DisplayName() + ": " + verb
}
