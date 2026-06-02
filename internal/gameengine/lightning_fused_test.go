package gameengine

import "testing"

// Tests the Lightning-fused flag: set/get round-trips and the turn reset clears it.
func TestLightningFused_SetGetResets(t *testing.T) {
	ge := New()
	if ge.HasLightningFused() {
		t.Fatal("HasLightningFused true on a fresh engine, want false")
	}
	ge.SetLightningFused(true)
	if !ge.HasLightningFused() {
		t.Error("HasLightningFused false after SetLightningFused(true)")
	}
	ge.ResetEphemeralState()
	if ge.HasLightningFused() {
		t.Error("HasLightningFused still true after ResetEphemeralState — it's per-turn and must clear")
	}
}
