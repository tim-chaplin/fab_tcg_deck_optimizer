package sim

import (
	"math/rand"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
)

// Test-only exports. Visible to package sim_test files in this directory only — these
// aliases never appear in the production binary because the file is _test.go. Cards' Play
// methods take *sim.TurnState / *sim.CardState, which means cards imports sim, which means
// sim test files that want to instantiate cards.X structs would form a cycle if they were
// in package sim. The fix is to keep those tests in package sim_test (a separate test
// package); this file re-exports the unexported helpers they rely on.

// IsExcludedFromPool re-exports isExcludedFromPool for sim_test consumers exercising the
// pool-exclusion marker contract directly.
func IsExcludedFromPool(c Card) bool { return isExcludedFromPool(c) }

// Best re-exports the package-private best for sim_test consumers. External-package
// (turntests) callers don't see this — only sim_test files in the same directory do.
func Best(hero Hero, weapons []Weapon, hand []Card, mp Matchup, deck []Card, prior TurnState) TurnSummary {
	return best(hero, weapons, hand, mp, deck, prior)
}

// SequenceContextForTest wraps *sequenceContext so sim_test files can drive
// playSequence / bestSequence without touching the unexported type directly. Production
// callers go through best / Evaluator.Best instead.
type SequenceContextForTest struct{ ctx *sequenceContext }

// NewSequenceContextForTest builds a sequenceContext with the same shape as the
// in-package newSequenceContextForTest helper.
func NewSequenceContextForTest(h Hero, pitched, deck []Card, resourceBudget, runechantCarryover, chainLen int) *SequenceContextForTest {
	return &SequenceContextForTest{ctx: newSequenceContextForTest(h, pitched, deck, resourceBudget, runechantCarryover, chainLen)}
}

// PlaySequence wraps (*sequenceContext).playSequence.
func (s *SequenceContextForTest) PlaySequence(order []Card) (damage int, futureValue int, residualBudget int, legal bool) {
	return s.ctx.playSequence(order)
}

// BestSequence wraps (*sequenceContext).bestSequence.
func (s *SequenceContextForTest) BestSequence(attackers []Card) (int, int, bool) {
	return s.ctx.bestSequence(attackers)
}

// FireAttackActionAuras re-exports fireAttackActionAuras for sim_test consumers.
func FireAttackActionAuras(state *TurnState, triggeringCard Card) {
	fireAttackActionAuras(state, triggeringCard)
}

// FireEndOfTurn re-exports fireEndOfTurn for sim_test consumers.
func FireEndOfTurn(state *TurnState) { fireEndOfTurn(state) }

// PromoteRandomHandCardToArsenal re-exports promoteRandomHandCardToArsenal for sim_test
// consumers exercising the post-hoc arsenal-promotion path in isolation.
func PromoteRandomHandCardToArsenal(best *TurnSummary, startingHand []Card, arsenalCardIn Card) {
	promoteRandomHandCardToArsenal(best, startingHand, arsenalCardIn)
}

// BeatsBest exposes findBest's partition-tiebreaker policy for direct testing.
func BeatsBest(v, futureValuePlayed int, willOccupyArsenal bool, best TurnSummary, bestFutureValuePlayed int, bestWillOccupyArsenal bool) bool {
	if cmp := chainScoreCmp(v, 0, futureValuePlayed, best.Value, 0, bestFutureValuePlayed); cmp != 0 {
		return cmp > 0
	}
	return willOccupyArsenal && !bestWillOccupyArsenal
}

// PairSwapMutations re-exports pairSwapMutations for sim_test consumers.
func PairSwapMutations(d *deck.Deck, legal func(deck.Card) bool) []Mutation {
	return pairSwapMutations(d, legal)
}

// CardMultisetKey re-exports cardMultisetKey for sim_test consumers. Takes []deck.Card so
// callers pass *deck.Deck.Cards directly; cardMultisetKey only reads .ID() so the narrow
// interface suffices.
func CardMultisetKey(cs []deck.Card) string { return cardMultisetKey(cs) }

// CardPairs re-exports the cardPairs registry for sim_test consumers.
var CardPairs = cardPairs

// FilterMaxCopiesViolations re-exports filterMaxCopiesViolations for sim_test consumers.
func FilterMaxCopiesViolations(muts []Mutation, maxCopies int) []Mutation {
	return filterMaxCopiesViolations(muts, maxCopies)
}

// RespectsMaxCopies re-exports respectsMaxCopies for sim_test consumers.
func RespectsMaxCopies(cs []deck.Card, maxCopies int) bool {
	return respectsMaxCopies(cs, maxCopies)
}

// AppendGroupedChainEntries re-exports appendGroupedChainEntries for sim_test consumers.
func AppendGroupedChainEntries(out []string, log []LogEntry) []string {
	return appendGroupedChainEntries(out, log)
}

// DefendersDamage re-exports defendersDamage for sim_test consumers with an unbounded
// block budget (no modal-blocker pressure). Drops the trailing cacheable bool so call
// sites that don't care about cacheable propagation stay terse.
func DefendersDamage(defenders, pitched, deck []Card, state *TurnState, gravBuf []Card, cs *CardState, incomingDamage, arsenalDefenderIdx int) (int, []Card) {
	total, gravBuf, _ := defendersDamage(defenders, pitched, deck, state, gravBuf, cs, incomingDamage, noBlockBudgetCap, arsenalDefenderIdx)
	return total, gravBuf
}

// DefendersDamageWithBudget is the budget-aware export — sim_test consumers exercising
// modal blockers (Brothers in Arms, …) pass blockBudget directly to gate the mode pick.
func DefendersDamageWithBudget(defenders, pitched, deck []Card, state *TurnState, gravBuf []Card, cs *CardState, incomingDamage, blockBudget, arsenalDefenderIdx int) (int, []Card) {
	total, gravBuf, _ := defendersDamage(defenders, pitched, deck, state, gravBuf, cs, incomingDamage, blockBudget, arsenalDefenderIdx)
	return total, gravBuf
}

// FormatContribution re-exports formatContribution for sim_test consumers.
func FormatContribution(v float64) string { return formatContribution(v) }

// AttackBufs is the exported alias of attackBufs for sim_test consumers. Type alias, not
// a fresh type — *AttackBufs and *attackBufs are identical, so wrappers thread cleanly.
type AttackBufs = attackBufs

// NewAttackBufs re-exports newAttackBufs.
func NewAttackBufs(handSize, weaponCount int, weapons []Weapon) *AttackBufs {
	return newAttackBufs(handSize, weaponCount, weapons)
}

// Bufs returns the wrapped sequenceContext's pooled scratch buffers. Lets sim_test files
// reach into per-permutation state mid-flight (e.g. assert ctx.Bufs().State().Graveyard
// after running a sequence).
func (s *SequenceContextForTest) Bufs() *AttackBufs { return s.ctx.bufs }

// State / DefenseGravScratch / DRCardStateScratch expose the unexported attackBufs fields
// the sim_test files probe for chain-state assertions.
func (b *attackBufs) State() *TurnState              { return b.state }
func (b *attackBufs) DefenseGravScratch() []Card     { return b.defenseGravScratch }
func (b *attackBufs) DRCardStateScratch() *CardState { return &b.drCardStateScratch }

// WeaponLoadouts re-exports weaponLoadouts for sim_test consumers.
func WeaponLoadouts(ws []deck.Weapon) [][]deck.Weapon { return weaponLoadouts(ws) }

// WeaponKey re-exports weaponKey for sim_test consumers.
func WeaponKey(ws []deck.Weapon) string { return weaponKey(ws) }

// SortedIDPair re-exports sortedIDPair for sim_test consumers.
func SortedIDPair(a, b ids.CardID) (ids.CardID, ids.CardID) { return sortedIDPair(a, b) }

// DeckFingerprint re-exports deckFingerprint, the deck-equality helper used by sim_test
// files. The underlying helper is in package sim because it reads the unexported weaponKey.
func DeckFingerprint(d *deck.Deck) string { return deckFingerprint(d) }

// PairAddAllowed re-exports pairAddAllowed for sim_test consumers.
func PairAddAllowed(c Card, legal func(deck.Card) bool) bool { return pairAddAllowed(c, legal) }

// EvaluateImplForTest re-exports the unexported (*Evaluator).evaluateImpl as an exported
// method for sim_test consumers exercising the eval-with-stop-condition path directly.
func (ev *Evaluator) EvaluateImplForTest(d *deck.Deck, maxRuns int, mp Matchup, rng *rand.Rand, stop func(stats *DeckStats, runs int) bool) DeckStats {
	return ev.evaluateImpl(d, maxRuns, mp, rng, stop)
}

// AdaptiveCheckInterval re-exports the adaptive-shuffle check interval constant.
const AdaptiveCheckInterval = adaptiveCheckInterval

// AdaptiveShufflesCap re-exports the adaptive-shuffle ceiling.
const AdaptiveShufflesCap = adaptiveShufflesCap

// MakeAdaptiveStop re-exports makeAdaptiveStop for sim_test consumers.
func MakeAdaptiveStop(targetSE float64) func(stats *DeckStats, runs int) bool {
	return makeAdaptiveStop(targetSE)
}

// MeanStandardError re-exports meanStandardError for sim_test consumers.
func MeanStandardError(stats *DeckStats) float64 { return meanStandardError(stats) }

// ProcessAurasAtStartOfTurn re-exports processAurasAtStartOfTurn for sim_test
// consumers exercising the start-of-turn aura-trigger pipeline in isolation.
func ProcessAurasAtStartOfTurn(queued []Aura, postDrawDeck []Card) (
	survivors []Aura,
	contribs []TriggerContribution,
	damage int,
	revealed []Card,
	graveyarded []Card,
) {
	return processAurasAtStartOfTurn(queued, postDrawDeck)
}
