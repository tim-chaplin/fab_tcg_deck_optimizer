package sim

import (
	"math/rand"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// Test-only exports. Visible to package sim_test files in this directory only.

// Best re-exports the package-private best for sim_test consumers.
func Best(hero gameengine.Hero, weapons []Weapon, hand []card.Card, mp Matchup, d *deck.Deck, prior gameengine.Spec) TurnSummary {
	return best(hero, weapons, hand, mp, d, prior)
}

// DeckOf builds a *deck.Deck from a list of cards.
func DeckOf(cards ...card.Card) *deck.Deck {
	dc := make([]deck.Card, len(cards))
	for i, c := range cards {
		dc[i] = c
	}
	return deck.New(nil, nil, dc)
}

// SequenceContextForTest wraps *sequenceContext so sim_test files can drive playSequence /
// bestSequence without touching the unexported type directly.
type SequenceContextForTest struct{ ctx *sequenceContext }

// NewSequenceContextForTest builds a sequenceContext with the same shape as the in-package
// newSequenceContextForTest helper.
func NewSequenceContextForTest(h gameengine.Hero, pitched, deck []card.Card, resourceBudget, runechantCarryover, chainLen int) *SequenceContextForTest {
	return &SequenceContextForTest{ctx: newSequenceContextForTest(h, pitched, deck, resourceBudget, runechantCarryover, chainLen)}
}

// PlaySequence wraps (*sequenceContext).playSequence.
func (s *SequenceContextForTest) PlaySequence(order []card.Card) (damage int, futureValue int, residualBudget int, legal bool) {
	return s.ctx.playSequence(order)
}

// BestSequence wraps (*sequenceContext).bestSequence.
func (s *SequenceContextForTest) BestSequence(attackers []card.Card) (int, int, bool) {
	return s.ctx.bestSequence(attackers)
}

// FireEndOfTurn re-exports the engine's end-of-turn fire for sim_test consumers.
func FireEndOfTurn(state *gameengine.GameEngine) { state.FireEndOfTurn() }

// PromoteRandomHandCardToArsenal re-exports promoteRandomHandCardToArsenal.
func PromoteRandomHandCardToArsenal(best *TurnSummary, startingHand []card.Card, arsenalCardIn card.Card) {
	promoteRandomHandCardToArsenal(best, startingHand, arsenalCardIn)
}

// BeatsBest exposes findBest's partition-tiebreaker policy for direct testing.
func BeatsBest(v, futureValuePlayed int, willOccupyArsenal bool, best TurnSummary, bestFutureValuePlayed int, bestWillOccupyArsenal bool) bool {
	if cmp := chainScoreCmp(v, 0, futureValuePlayed, best.Value, 0, bestFutureValuePlayed); cmp != 0 {
		return cmp > 0
	}
	return willOccupyArsenal && !bestWillOccupyArsenal
}

// AppendGroupedChainEntries re-exports appendGroupedChainEntries.
func AppendGroupedChainEntries(out []string, log []turnlogger.LogEntry) []string {
	return appendGroupedChainEntries(out, log)
}

// DefendersDamage re-exports defendersDamage with an unbounded block budget.
func DefendersDamage(defenders, pitched []card.Card, d *deck.Deck, state *gameengine.GameEngine, gravBuf []card.Card, cs *card.CardState, incomingDamage, arsenalDefenderIdx int) (int, []card.Card) {
	total, gravBuf, _ := defendersDamage(defenders, pitched, d, state, gravBuf, cs, incomingDamage, noBlockBudgetCap, arsenalDefenderIdx)
	return total, gravBuf
}

// DefendersDamageWithBudget is the budget-aware export.
func DefendersDamageWithBudget(defenders, pitched []card.Card, d *deck.Deck, state *gameengine.GameEngine, gravBuf []card.Card, cs *card.CardState, incomingDamage, blockBudget, arsenalDefenderIdx int) (int, []card.Card) {
	total, gravBuf, _ := defendersDamage(defenders, pitched, d, state, gravBuf, cs, incomingDamage, blockBudget, arsenalDefenderIdx)
	return total, gravBuf
}

// FormatContribution re-exports formatContribution.
func FormatContribution(v float64) string { return formatContribution(v) }

// AttackBufs is the exported alias of attackBufs.
type AttackBufs = attackBufs

// NewAttackBufs re-exports newAttackBufs.
func NewAttackBufs(handSize, weaponCount int, weapons []Weapon) *AttackBufs {
	return newAttackBufs(handSize, weaponCount, weapons)
}

// Bufs returns the wrapped sequenceContext's pooled scratch buffers.
func (s *SequenceContextForTest) Bufs() *AttackBufs { return s.ctx.bufs }

// State / DefenseGravScratch / DRCardStateScratch expose the unexported attackBufs fields.
func (b *attackBufs) State() *gameengine.GameEngine       { return b.state }
func (b *attackBufs) DefenseGravScratch() []card.Card     { return b.defenseGravScratch }
func (b *attackBufs) DRCardStateScratch() *card.CardState { return &b.drCardStateScratch }

// EvaluateImplForTest re-exports the unexported (*Evaluator).evaluateImpl.
func (ev *Evaluator) EvaluateImplForTest(d *deck.Deck, maxRuns int, mp Matchup, rng *rand.Rand, stop func(stats *DeckStats, runs int) bool) DeckStats {
	return ev.evaluateImpl(d, maxRuns, mp, rng, stop)
}

// AdaptiveCheckInterval re-exports the adaptive-shuffle check interval constant.
const AdaptiveCheckInterval = adaptiveCheckInterval

// AdaptiveShufflesCap re-exports the adaptive-shuffle ceiling.
const AdaptiveShufflesCap = adaptiveShufflesCap

// MakeAdaptiveStop re-exports makeAdaptiveStop.
func MakeAdaptiveStop(targetSE float64) func(stats *DeckStats, runs int) bool {
	return makeAdaptiveStop(targetSE)
}

// MeanStandardError re-exports meanStandardError.
func MeanStandardError(stats *DeckStats) float64 { return meanStandardError(stats) }

// ProcessAurasAtStartOfTurn re-exports processAurasAtStartOfTurn.
func ProcessAurasAtStartOfTurn(queued []gameengine.Aura, d *deck.Deck) (
	survivors []gameengine.Aura,
	contribs []TriggerContribution,
	damage int,
	revealed []card.Card,
	graveyarded []card.Card,
) {
	return processAurasAtStartOfTurn(queued, d)
}
