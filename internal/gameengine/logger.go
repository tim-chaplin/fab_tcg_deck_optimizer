package gameengine

import (
	"fmt"
	"io"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// NoopLogger satisfies card.Logger with every method as a no-op. Value-type so it costs
// nothing to embed in a GameState field. Used as the default logger during evaluation
// where no rendered output is wanted — every chain step / trigger emission compiles to
// "call into iface, jump to empty function, return."
type NoopLogger struct{}

func (NoopLogger) AppendChainStep(string, int)                       {}
func (NoopLogger) AppendChainStepf(int, string, ...any)              {}
func (NoopLogger) AppendPostTrigger(string, string, int)             {}
func (NoopLogger) AppendPostTriggerf(string, int, string, ...any)    {}
func (NoopLogger) AppendPreTrigger(string, string, int)              {}
func (NoopLogger) AppendPreTriggerf(string, int, string, ...any)     {}
func (NoopLogger) AmendLastChainStepN(int)                           {}

// StreamLogger satisfies card.Logger by formatting each emission directly to w. Used at
// print time to render a turn's chain log straight to stdout (or any io.Writer) as the
// chain runner replays it — nothing accumulates.
//
// Chain steps render as "text (N)" or "text" when N is zero; triggers indent under their
// source. AmendLastChainStepN is a no-op (streaming can't go back and edit a written
// line) — Attack-Reaction buffs surface as their own trigger line instead, so readers
// see the buff as a discrete event rather than baked into the prior chain step's N.
type StreamLogger struct {
	w io.Writer
}

// NewStreamLogger returns a StreamLogger that writes to w.
func NewStreamLogger(w io.Writer) *StreamLogger { return &StreamLogger{w: w} }

func (s *StreamLogger) AppendChainStep(text string, n int) {
	if n != 0 {
		fmt.Fprintf(s.w, "%s (%d)\n", text, n)
	} else {
		fmt.Fprintln(s.w, text)
	}
}

func (s *StreamLogger) AppendChainStepf(n int, format string, args ...any) {
	s.AppendChainStep(fmt.Sprintf(format, args...), n)
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

func (s *StreamLogger) AmendLastChainStepN(int) {}

func (s *StreamLogger) appendTrigger(source, text string, n int) {
	switch {
	case source == "" && n != 0:
		fmt.Fprintf(s.w, "  %s (%d)\n", text, n)
	case source == "":
		fmt.Fprintf(s.w, "  %s\n", text)
	case n != 0:
		fmt.Fprintf(s.w, "  %s: %s (%d)\n", source, text, n)
	default:
		fmt.Fprintf(s.w, "  %s: %s\n", source, text)
	}
}

// Compile-time check that both loggers satisfy card.Logger.
var (
	_ card.Logger = NoopLogger{}
	_ card.Logger = (*StreamLogger)(nil)
)
