package sim

// newSequenceContextForTest builds a sequenceContext wired to a fresh attackBufs sized for
// the given chain length. Tests use this instead of hand-rolling the context fields so the
// common shape is centralised. Lives in a sim test file (rather than sim_test) so
// exports_test.go's NewSequenceContextForTest wrapper can reach it.
func newSequenceContextForTest(h Hero, pitched, deck []Card, resourceBudget, runechantCarryover, chainLen int) *sequenceContext {
	bufs := newAttackBufs(chainLen, 0, nil)
	return &sequenceContext{
		hero:               h,
		pitched:            pitched,
		deck:               deck,
		bufs:               bufs,
		resourceBudget:     resourceBudget,
		runechantCarryover: runechantCarryover,
		carryWinner:        &bufs.carryWinnerScratch,
	}
}
