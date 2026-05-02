package sim

// Matchup bundles the opponent-profile parameters that stay constant for the duration of
// an Evaluator's lifetime / a deck eval session. Adding a new matchup-level parameter
// (Mark / Tap state, hero life total, etc.) adds a field here rather than threading a new
// int through every Best / Evaluate / IterateParallel signature.
//
// The chain runner copies these fields onto TurnState in resetStateForPermutation; cards
// read the per-turn copy (s.IncomingDamage, s.ArcaneIncomingDamage) rather than reaching
// back into the Matchup so per-card hot paths stay one struct field deep.
type Matchup struct {
	// IncomingDamage is the opponent damage per turn.
	IncomingDamage int
	// ArcaneIncomingDamage is the opponent's arcane damage per turn.
	ArcaneIncomingDamage int
}
