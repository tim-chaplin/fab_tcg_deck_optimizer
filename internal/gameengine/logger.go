package gameengine

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// NoopLogger satisfies card.Logger with every method as a no-op. Value-type so it costs
// nothing to embed in a GameState field — every chain step / trigger emission compiles to
// "call into iface, jump to empty function, return."
type NoopLogger struct{}

func (NoopLogger) AppendChainStep(string, int)                    {}
func (NoopLogger) AppendChainStepf(int, string, ...any)           {}
func (NoopLogger) AppendPostTrigger(string, string, int)          {}
func (NoopLogger) AppendPostTriggerf(string, int, string, ...any) {}
func (NoopLogger) AppendPreTrigger(string, string, int)           {}
func (NoopLogger) AppendPreTriggerf(string, int, string, ...any)  {}
func (NoopLogger) AmendLastChainStepN(int)                        {}

// Compile-time check that NoopLogger satisfies card.Logger.
var _ card.Logger = NoopLogger{}
