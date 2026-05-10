package sim

import "github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"

// simLogger is the framework-side concrete Logger implementation. Each method takes the
// CardState-aware shape cards expect and forwards to a *turnlogger.TurnLogger that owns
// the entries slice. The CardState->ChainStepText conversion lives here so v2/turnlogger
// stays free of sim-package coupling.
//
// Nil-safety: a typed-nil *simLogger is the find-best silent mode — every method
// returns at the entry guard so callers stay branch-free. The underlying *TurnLogger is
// nil-safe in its own right; the layered guards keep both paths defensive.
type simLogger struct {
	tl *turnlogger.TurnLogger
}

// newSimLogger wraps a *turnlogger.TurnLogger so it satisfies Logger. Allocated once per
// attackBufs alongside the TurnLogger and reused across permutations.
func newSimLogger(tl *turnlogger.TurnLogger) *simLogger {
	return &simLogger{tl: tl}
}

// Entries surfaces the underlying log stream for tests / format-layer readers. Returns
// nil for a nil receiver so the find-best path's "no logger" mode reads as "no entries".
func (l *simLogger) Entries() []LogEntry {
	if l == nil {
		return nil
	}
	return l.tl.Entries()
}

func (l *simLogger) Log(self *CardState, n int) {
	if l == nil {
		return
	}
	l.tl.AppendChainStep(ChainStepText(self), n)
}

func (l *simLogger) Logf(n int, format string, args ...any) {
	if l == nil {
		return
	}
	l.tl.AppendChainStepf(n, format, args...)
}

func (l *simLogger) LogRider(self *CardState, n int, text string) {
	if l == nil {
		return
	}
	l.tl.AppendPostTrigger(self.Card.DisplayName(), text, n)
}

func (l *simLogger) LogRiderf(self *CardState, n int, format string, args ...any) {
	if l == nil {
		return
	}
	l.tl.AppendPostTriggerf(self.Card.DisplayName(), n, format, args...)
}

func (l *simLogger) LogPreTrigger(source, text string, n int) {
	if l == nil {
		return
	}
	l.tl.AppendPreTrigger(source, text, n)
}

func (l *simLogger) LogPreTriggerf(source string, n int, format string, args ...any) {
	if l == nil {
		return
	}
	l.tl.AppendPreTriggerf(source, n, format, args...)
}

func (l *simLogger) LogPostTrigger(source, text string, n int) {
	if l == nil {
		return
	}
	l.tl.AppendPostTrigger(source, text, n)
}

func (l *simLogger) LogPostTriggerf(source string, n int, format string, args ...any) {
	if l == nil {
		return
	}
	l.tl.AppendPostTriggerf(source, n, format, args...)
}

func (l *simLogger) AmendLastChainStepN(n int) {
	if l == nil {
		return
	}
	l.tl.AmendLastChainStepN(n)
}
