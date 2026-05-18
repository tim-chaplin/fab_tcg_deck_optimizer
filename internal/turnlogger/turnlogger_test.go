package turnlogger

import "testing"

// Append* on a non-nil TurnLogger records each entry in order; Entries() returns the
// stream the format layer reads. Smoke test for the basic happy path.
func TestTurnLogger_AppendsEntriesInOrder(t *testing.T) {
	l := New()
	l.AppendChainStep("attack", 3)
	l.AppendPostTrigger("attack", "rider", 1)
	l.AppendPreTrigger("aura", "pre", 2)

	got := l.Entries()
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	if got[0].Kind != LogEntryChainStep || got[0].Text != "attack" || got[0].N != 3 {
		t.Errorf("entry 0 = %+v, want chain step", got[0])
	}
	if got[1].Kind != LogEntryPostTrigger || got[1].Source != "attack" || got[1].N != 1 {
		t.Errorf("entry 1 = %+v, want post trigger", got[1])
	}
	if got[2].Kind != LogEntryPreTrigger || got[2].Source != "aura" || got[2].N != 2 {
		t.Errorf("entry 2 = %+v, want pre trigger", got[2])
	}
}

// Every method handles a nil receiver as a silent no-op so the eval loop can pass nil
// during the find-best pass without per-call branching.
func TestTurnLogger_NilReceiverIsNoop(t *testing.T) {
	var l *TurnLogger
	l.AppendChainStep("x", 1)
	l.AppendChainStepf(0, "%d", 7)
	l.AppendPostTrigger("a", "b", 0)
	l.AppendPostTriggerf("a", 0, "%s", "c")
	l.AppendPreTrigger("a", "b", 0)
	l.AppendPreTriggerf("a", 0, "%s", "c")
	l.AmendLastChainStepN(5)
	l.SetBuffer(make([]LogEntry, 0, 4))
	l.Reset()
	if got := l.Entries(); got != nil {
		t.Errorf("nil-receiver Entries = %v, want nil", got)
	}
}

// AmendLastChainStepN walks back over post/pre triggers to land on the most recent
// chain-step entry. Mirrors the AR buff path: a ChainStep is pushed first, then a
// rider, then the AR buff amends the original chain step's N.
func TestTurnLogger_AmendLastChainStepN_SkipsTriggers(t *testing.T) {
	l := New()
	l.AppendChainStep("first", 2)
	l.AppendPostTrigger("first", "rider", 0)
	l.AmendLastChainStepN(5)

	got := l.Entries()
	if got[0].N != 7 {
		t.Errorf("chain step N = %d, want 7", got[0].N)
	}
	if got[1].N != 0 {
		t.Errorf("post trigger N = %d, want 0", got[1].N)
	}
}

// AmendLastChainStepN with no chain-step entry recorded yet is a no-op rather than a
// panic.
func TestTurnLogger_AmendLastChainStepN_NoChainStepIsNoop(t *testing.T) {
	l := New()
	l.AppendPostTrigger("a", "b", 1)
	l.AmendLastChainStepN(3) // no chain step exists; should not panic
	if l.Entries()[0].N != 1 {
		t.Errorf("post trigger N = %d, want 1 (untouched)", l.Entries()[0].N)
	}
}

// SetBuffer rebinds the entries slice so the chain runner's per-permutation scratch
// backs the logger without per-permutation allocation.
func TestTurnLogger_SetBuffer_RebindsBacking(t *testing.T) {
	l := New()
	buf := make([]LogEntry, 0, 8)
	l.SetBuffer(buf)
	for i := 0; i < 4; i++ {
		l.AppendChainStep("x", i)
	}
	if cap(l.Entries()) != 8 {
		t.Errorf("cap = %d, want 8 (SetBuffer should retain pre-sized backing)", cap(l.Entries()))
	}
	if len(l.Entries()) != 4 {
		t.Errorf("len = %d, want 4", len(l.Entries()))
	}
}

// Negative N is clamped to zero by every Append* method — mirrors TurnState.Log's old
// behavior for the Value-credit pairing convention.
func TestTurnLogger_NegativeNClampsToZero(t *testing.T) {
	l := New()
	l.AppendChainStep("a", -3)
	l.AppendPostTrigger("a", "b", -1)
	got := l.Entries()
	if got[0].N != 0 {
		t.Errorf("chain step N = %d, want 0", got[0].N)
	}
	if got[1].N != 0 {
		t.Errorf("post trigger N = %d, want 0", got[1].N)
	}
}
