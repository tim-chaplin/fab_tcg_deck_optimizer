package gameengine

import (
	"fmt"
	"io"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// NoopLogger satisfies card.Logger with every method as a no-op. Value-type so it costs
// nothing to embed in a GameState field — every attack step / trigger emission compiles to
// "call into iface, jump to empty function, return." The eval hot path installs this so
// no log work happens anywhere off the print path.
type NoopLogger struct{}

func (NoopLogger) AppendAttackStep(string, int)                   {}
func (NoopLogger) AppendAttackStepf(int, string, ...any)          {}
func (NoopLogger) AppendPostTrigger(string, string, int)          {}
func (NoopLogger) AppendPostTriggerf(string, int, string, ...any) {}
func (NoopLogger) AppendPreTrigger(string, string, int)           {}
func (NoopLogger) AppendPreTriggerf(string, int, string, ...any)  {}
func (NoopLogger) AmendLastAttackStepN(int)                       {}

// StreamLogger satisfies card.Logger by formatting each emission directly to w. The print
// path installs this on the replay state so attack steps and trigger riders land in the
// writer as they happen — no accumulation, no buffering. AmendLastAttackStepN is a no-op:
// the prior attack-step line has already been written, so AR buffs surface as their own
// trigger line rather than mutating the previous line in place (a "less pretty" trade-off
// the design explicitly accepts).
type StreamLogger struct {
	w io.Writer
}

// NewStreamLogger returns a StreamLogger that writes to w.
func NewStreamLogger(w io.Writer) *StreamLogger { return &StreamLogger{w: w} }

func (s *StreamLogger) AppendAttackStep(text string, n int) {
	if n != 0 {
		fmt.Fprintf(s.w, "%s (+%d)\n", text, n)
		return
	}
	fmt.Fprintln(s.w, text)
}

func (s *StreamLogger) AppendAttackStepf(n int, format string, args ...any) {
	s.AppendAttackStep(fmt.Sprintf(format, args...), n)
}

func (s *StreamLogger) AppendPostTrigger(source, text string, n int) {
	s.appendTrigger(source, text, n)
}

func (s *StreamLogger) AppendPostTriggerf(source string, n int, format string, args ...any) {
	s.appendTrigger(source, fmt.Sprintf(format, args...), n)
}

func (s *StreamLogger) AppendPreTrigger(source, text string, n int) {
	s.appendTrigger(source, text, n)
}

func (s *StreamLogger) AppendPreTriggerf(source string, n int, format string, args ...any) {
	s.appendTrigger(source, fmt.Sprintf(format, args...), n)
}

func (s *StreamLogger) AmendLastAttackStepN(int) {}

// appendTrigger renders a trigger line indented to sit beneath its attack step. The source
// tag is omitted from the wire form because the indentation already conveys "child of the
// prior attack step"; orphan triggers (no matching parent) would lose attribution under
// this rule, but the attack-turn runner pairs them tightly so orphans don't occur in practice.
func (s *StreamLogger) appendTrigger(_ string, text string, n int) {
	if n != 0 {
		fmt.Fprintf(s.w, "     %s (+%d)\n", text, n)
		return
	}
	fmt.Fprintf(s.w, "     %s\n", text)
}

// Compile-time check that both loggers satisfy card.Logger.
var (
	_ card.Logger = NoopLogger{}
	_ card.Logger = (*StreamLogger)(nil)
)
