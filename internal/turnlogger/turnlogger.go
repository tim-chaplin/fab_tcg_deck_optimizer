// Package turnlogger collects the per-turn LogEntry stream the sim builds during chain
// resolution and the format layer renders into the printout. The package owns the data
// types (LogEntry, LogEntryKind) and the TurnLogger primitive that accumulates entries;
// it knows nothing about cards, decks, or any other sim concept. Callers feed it raw
// (text, source, n) tuples — the CardState->ChainStepText conversion lives in the
// caller, keeping turnlogger free of sim-package coupling.
//
// Nil-safety is the core ergonomic: every method on *TurnLogger handles a nil receiver
// silently, so the eval loop's "find best" pass passes a nil pointer and skips log
// materialization without per-call branching at the call site. The "replay best turn"
// pass swaps in a real *TurnLogger and the same calls write into its buffer.
package turnlogger

import "fmt"

// LogEntryKind classifies a LogEntry. Triggers come in two flavours because they fire on
// opposite sides of their parent in the FaB stack — the format layer needs to know which
// side a given entry sits on so two cards with the same display name in the same chain
// don't steal each other's triggers during grouping.
type LogEntryKind int8

const (
	// LogEntryChainStep is the sim's "<Card>: <VERB>" line. Stands alone in the printout
	// and acts as the parent that triggers attach to.
	LogEntryChainStep LogEntryKind = iota
	// LogEntryPreTrigger is a trigger that fires before its parent chain entry resolves
	// (hero / aura attack-action triggers). The format layer attaches it to the next
	// chain entry whose name matches Source.
	LogEntryPreTrigger
	// LogEntryPostTrigger is a trigger that fires after its parent chain entry resolves
	// (OnHit riders, AR buffs). The format layer attaches it to the previous chain entry
	// whose name matches Source.
	LogEntryPostTrigger
)

// LogEntry is one chain-event entry in TurnLogger's stream. Text is the freeform display
// string the producer authored ("Viserai created a runechant", "Consuming Volition [R]:
// ATTACK") — the format layer renders it verbatim. Kind tags the entry as a chain step
// or as a pre/post trigger so the grouping algorithm can attribute triggers correctly
// even when sibling chain entries share a name. Source names the card whose play caused
// the trigger and is matched against chain-entry names to pick the parent; only
// meaningful for triggers. N is the damage-equivalent credited when the entry was added.
type LogEntry struct {
	Text   string
	Source string
	Kind   LogEntryKind
	N      int
}

// TurnLogger accumulates the per-turn LogEntry stream. A nil *TurnLogger is the silent
// "find best" mode: every method returns without recording, so callers don't branch at
// the call site. A non-nil *TurnLogger appends each call to its entries slice; SetBuffer
// rebinds the backing storage at permutation boundaries so the sim's pre-sized scratch
// buffer is reused across permutations without allocation.
type TurnLogger struct {
	entries []LogEntry
}

// New returns a fresh empty *TurnLogger. Use SetBuffer afterwards to bind a pre-sized
// backing slice when the caller owns one.
func New() *TurnLogger { return &TurnLogger{} }

// SetBuffer rebinds the entries slice to buf (truncated to zero length). Used by the
// per-permutation reset to point the logger at the chain-runner's pre-sized scratch.
func (l *TurnLogger) SetBuffer(buf []LogEntry) {
	if l == nil {
		return
	}
	l.entries = buf[:0]
}

// Reset truncates the entries slice to zero length, preserving backing storage.
func (l *TurnLogger) Reset() {
	if l == nil {
		return
	}
	l.entries = l.entries[:0]
}

// Entries returns the accumulated log entries. Nil receiver returns nil.
func (l *TurnLogger) Entries() []LogEntry {
	if l == nil {
		return nil
	}
	return l.entries
}

// AppendChainStep appends a chain-step entry with text and damage-equivalent n. n is
// clamped at zero. No-op on a nil receiver.
func (l *TurnLogger) AppendChainStep(text string, n int) {
	if l == nil {
		return
	}
	if n < 0 {
		n = 0
	}
	l.entries = append(l.entries, LogEntry{Kind: LogEntryChainStep, Text: text, N: n})
}

// AppendChainStepf appends a chain-step entry whose text is fmt.Sprintf(format, args...).
// fmt.Sprintf only runs on a non-nil receiver — callers still pay variadic boxing at the
// call site regardless.
func (l *TurnLogger) AppendChainStepf(n int, format string, args ...any) {
	if l == nil {
		return
	}
	if n < 0 {
		n = 0
	}
	l.entries = append(l.entries, LogEntry{Kind: LogEntryChainStep, Text: fmt.Sprintf(format, args...), N: n})
}

// AppendPostTrigger appends a post-trigger entry attributed to source.
func (l *TurnLogger) AppendPostTrigger(source, text string, n int) {
	if l == nil {
		return
	}
	if n < 0 {
		n = 0
	}
	l.entries = append(l.entries, LogEntry{Kind: LogEntryPostTrigger, Source: source, Text: text, N: n})
}

// AppendPostTriggerf is the format variant of AppendPostTrigger.
func (l *TurnLogger) AppendPostTriggerf(source string, n int, format string, args ...any) {
	if l == nil {
		return
	}
	if n < 0 {
		n = 0
	}
	l.entries = append(l.entries, LogEntry{Kind: LogEntryPostTrigger, Source: source, Text: fmt.Sprintf(format, args...), N: n})
}

// AppendPreTrigger appends a pre-trigger entry attributed to source.
func (l *TurnLogger) AppendPreTrigger(source, text string, n int) {
	if l == nil {
		return
	}
	if n < 0 {
		n = 0
	}
	l.entries = append(l.entries, LogEntry{Kind: LogEntryPreTrigger, Source: source, Text: text, N: n})
}

// AppendPreTriggerf is the format variant of AppendPreTrigger.
func (l *TurnLogger) AppendPreTriggerf(source string, n int, format string, args ...any) {
	if l == nil {
		return
	}
	if n < 0 {
		n = 0
	}
	l.entries = append(l.entries, LogEntry{Kind: LogEntryPreTrigger, Source: source, Text: fmt.Sprintf(format, args...), N: n})
}

// AmendLastChainStepN adds n to the most recent ChainStep entry's N field. ARs use this
// to fold their +{p} buff into the buffed attack's display delta. No-op on n == 0 and
// when no chain-step entry exists yet.
func (l *TurnLogger) AmendLastChainStepN(n int) {
	if l == nil || n == 0 {
		return
	}
	for i := len(l.entries) - 1; i >= 0; i-- {
		if l.entries[i].Kind == LogEntryChainStep {
			l.entries[i].N += n
			return
		}
	}
}
