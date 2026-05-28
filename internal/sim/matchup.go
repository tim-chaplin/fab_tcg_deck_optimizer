package sim

// Matchup bundles the opponent-profile parameters that stay constant for the duration of
// an Evaluator's lifetime / a deck eval session. Adding a new matchup-level parameter
// (Mark / Tap state, hero life total, etc.) adds a field here rather than threading a new
// int through every Best / Evaluate / RunMutationRound signature.
//
// The attack-turn runner copies these fields onto TurnState in resetStateForPermutation; cards
// read the per-turn copy (gs.incomingPhysicalDamage, gs.incomingArcaneDamage) rather than reaching
// back into the Matchup so per-card hot paths stay one struct field deep.
type Matchup struct {
	// IncomingPhysicalDamage is the opponent's physical damage per turn.
	IncomingPhysicalDamage int
	// IncomingArcaneDamage is the opponent's arcane damage per turn.
	IncomingArcaneDamage int
}
