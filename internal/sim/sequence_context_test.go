package sim

// newSequenceContextForTest builds a sequenceContext wired to a fresh attackBufs sized for
// the given chain length. Tests use this instead of hand-rolling the context fields so the
// common shape is centralised. Lives in a sim test file (rather than sim_test) so
// exports_test.go's NewSequenceContextForTest wrapper can reach it.
func newSequenceContextForTest(h Hero, pitched, deck []Card, resourceBudget, runechantCarryover, chainLen int) *sequenceContext {
	return &sequenceContext{
		hero:               h,
		pitched:            pitched,
		deck:               deck,
		bufs:               newAttackBufs(chainLen, 0, nil),
		resourceBudget:     resourceBudget,
		runechantCarryover: runechantCarryover,
	}
}
