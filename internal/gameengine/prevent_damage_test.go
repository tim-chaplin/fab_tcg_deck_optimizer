package gameengine

import "testing"

// Tests that PreventArcaneDamage banks into the per-turn arcaneDamageBlocked accumulator
// rather than mutating the constant matchup figure, so prevention on one turn can't leak
// into later turns. Regression for the bug where PreventArcaneDamage decremented
// incomingArcaneDamage itself — a carryover field — permanently shrinking arcane for every
// subsequent turn of the shuffle.
func TestPreventArcaneDamage_ResetsPerTurn(t *testing.T) {
	gs := GameStateBuilder().SetIncomingArcaneDamage(3).Build()
	ge := &GameEngine{GameState: gs}

	if got := ge.PreventArcaneDamage(2); got != 2 {
		t.Fatalf("PreventArcaneDamage(2) = %d, want 2", got)
	}
	if got := ge.IncomingArcaneDamage(); got != 3 {
		t.Errorf("IncomingArcaneDamage = %d after prevention, want 3 (matchup figure must stay constant)", got)
	}
	if got := ge.RemainingArcaneDamage(); got != 1 {
		t.Errorf("RemainingArcaneDamage = %d, want 1 (3 incoming - 2 prevented)", got)
	}

	// Turn boundary: the per-turn accumulator resets, so the full matchup arcane returns.
	gs.ResetEphemeralState()
	if got := ge.RemainingArcaneDamage(); got != 3 {
		t.Errorf("RemainingArcaneDamage = %d after turn reset, want 3 (prevention must not leak across turns)", got)
	}
}

// Tests that PreventPhysicalDamage and PreventArcaneDamage bank into independent
// accumulators — preventing one type doesn't consume the other's budget.
func TestPreventDamage_PhysicalAndArcaneIndependent(t *testing.T) {
	gs := GameStateBuilder().
		SetIncomingPhysicalDamage(4).
		SetIncomingArcaneDamage(2).
		Build()
	ge := &GameEngine{GameState: gs}

	if got := ge.PreventPhysicalDamage(3); got != 3 {
		t.Fatalf("PreventPhysicalDamage(3) = %d, want 3", got)
	}
	if got := ge.RemainingPhysicalDamage(); got != 1 {
		t.Errorf("RemainingPhysicalDamage = %d, want 1", got)
	}
	if got := ge.RemainingArcaneDamage(); got != 2 {
		t.Errorf("RemainingArcaneDamage = %d, want 2 (physical prevention must not touch arcane)", got)
	}
}
