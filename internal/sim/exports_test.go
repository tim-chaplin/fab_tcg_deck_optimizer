package sim

import (
	"math/rand"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
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

// SequenceContextForTest wraps *sequenceContext so sim_test files can drive
// playSequence / bestSequence without touching the unexported type directly. Production
// callers go through Best / BestWithTriggers instead.
type SequenceContextForTest struct{ ctx *sequenceContext }

// NewSequenceContextForTest builds a sequenceContext with the same shape as the
// in-package newSequenceContextForTest helper.
func NewSequenceContextForTest(h Hero, pitched, deck []Card, resourceBudget, runechantCarryover, chainLen int) *SequenceContextForTest {
	return &SequenceContextForTest{ctx: newSequenceContextForTest(h, pitched, deck, resourceBudget, runechantCarryover, chainLen)}
}

// PlaySequence wraps (*sequenceContext).playSequence.
func (s *SequenceContextForTest) PlaySequence(order []Card) (damage int, leftoverRunechants int, residualBudget int, legal bool) {
	return s.ctx.playSequence(order)
}

// BestSequence wraps (*sequenceContext).bestSequence.
func (s *SequenceContextForTest) BestSequence(attackers []Card) (int, int, bool) {
	return s.ctx.bestSequence(attackers)
}

// FireAttackActionTriggers re-exports fireAttackActionTriggers for sim_test consumers.
func FireAttackActionTriggers(state *TurnState, triggeringCard Card) {
	fireAttackActionTriggers(state, triggeringCard)
}

// FireEphemeralAttackTriggers re-exports fireEphemeralAttackTriggers for sim_test consumers.
func FireEphemeralAttackTriggers(state *TurnState, target *CardState) {
	fireEphemeralAttackTriggers(state, target)
}

// PromoteRandomHandCardToArsenal re-exports promoteRandomHandCardToArsenal for sim_test
// consumers exercising the post-hoc arsenal-promotion path in isolation.
func PromoteRandomHandCardToArsenal(best *TurnSummary, startingHand []Card, arsenalCardIn Card) {
	promoteRandomHandCardToArsenal(best, startingHand, arsenalCardIn)
}

// BeatsBest is the test-only entry point for the partition-tiebreaker policy. The shim
// builds a runningCarry seeded with the comparison's "current best" stats and asks
// whether the candidate would displace it; that's the same path findBest's recurse
// takes, so the tiebreak order under test is the production order.
//
// hasHeld is synthesised from bestWillOccupyArsenal: any non-nil arsenal already
// satisfies the willOccupy predicate, so passing hasHeld=true is enough to make the
// running winner's willOccupy=true regardless of the (test-ignored) arsenal field.
func BeatsBest(v, leftoverRunechants, futureValuePlayed int, willOccupyArsenal bool, best TurnSummary, bestFutureValuePlayed int, bestWillOccupyArsenal bool) bool {
	r := runningCarry{
		seen:               true,
		value:              best.Value,
		leftoverRunechants: best.State.Runechants,
		futureValuePlayed:  bestFutureValuePlayed,
		hasHeld:            bestWillOccupyArsenal,
	}
	return r.Beats(v, leftoverRunechants, futureValuePlayed, willOccupyArsenal)
}

// PairSwapMutations re-exports pairSwapMutations for sim_test consumers.
func PairSwapMutations(d *Deck, legal func(Card) bool) []Mutation {
	return pairSwapMutations(d, legal)
}

// CardMultisetKey re-exports cardMultisetKey for sim_test consumers.
func CardMultisetKey(cs []Card) string { return cardMultisetKey(cs) }

// CardPairs re-exports the cardPairs registry for sim_test consumers.
var CardPairs = cardPairs

// FilterMaxCopiesViolations re-exports filterMaxCopiesViolations for sim_test consumers.
func FilterMaxCopiesViolations(muts []Mutation, maxCopies int) []Mutation {
	return filterMaxCopiesViolations(muts, maxCopies)
}

// RespectsMaxCopies re-exports respectsMaxCopies for sim_test consumers.
func RespectsMaxCopies(cs []Card, maxCopies int) bool { return respectsMaxCopies(cs, maxCopies) }

// AppendGroupedChainEntries re-exports appendGroupedChainEntries for sim_test consumers.
func AppendGroupedChainEntries(out []string, log []LogEntry) []string {
	return appendGroupedChainEntries(out, log)
}

// DefendersDamage re-exports defendersDamage for sim_test consumers. Drops the trailing
// cacheable bool so call sites that don't care about cacheable propagation stay terse.
func DefendersDamage(defenders, pitched, deck []Card, state *TurnState, gravBuf []Card, cs *CardState, incomingDamage, arsenalDefenderIdx int) (int, []Card) {
	total, gravBuf, _ := defendersDamage(defenders, pitched, deck, state, gravBuf, cs, incomingDamage, arsenalDefenderIdx)
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
func (b *attackBufs) State() *TurnState                    { return b.state }
func (b *attackBufs) DefenseGravScratch() []Card           { return b.defenseGravScratch }
func (b *attackBufs) SetDefenseGravScratch(scratch []Card) { b.defenseGravScratch = scratch }
func (b *attackBufs) DRCardStateScratch() *CardState       { return &b.drCardStateScratch }

// LegalPool re-exports legalPool for sim_test consumers.
func LegalPool(legal func(Card) bool) []ids.CardID { return legalPool(legal) }

// LegalWeapons re-exports legalWeapons for sim_test consumers.
func LegalWeapons() []Weapon { return legalWeapons() }

// WeaponLoadouts re-exports weaponLoadouts for sim_test consumers.
func WeaponLoadouts(ws []Weapon) [][]Weapon { return weaponLoadouts(ws) }

// WeaponKey re-exports weaponKey for sim_test consumers.
func WeaponKey(ws []Weapon) string { return weaponKey(ws) }

// SortedIDPair re-exports sortedIDPair for sim_test consumers.
func SortedIDPair(a, b ids.CardID) (ids.CardID, ids.CardID) { return sortedIDPair(a, b) }

// DeckFingerprint re-exports deckFingerprint, the deck-equality helper used by sim_test
// files. The underlying helper is in package sim because it reads the unexported weaponKey.
func DeckFingerprint(d *Deck) string { return deckFingerprint(d) }

// PairAddAllowed re-exports pairAddAllowed for sim_test consumers.
func PairAddAllowed(c Card, legal func(Card) bool) bool { return pairAddAllowed(c, legal) }

// EvaluateImplForTest re-exports the unexported (*Deck).evaluateImpl as an exported method
// for sim_test consumers exercising the eval-with-stop-condition path directly.
func (d *Deck) EvaluateImplForTest(maxRuns int, incomingDamage int, rng *rand.Rand, ev *Evaluator, stop func(stats *Stats, runs int) bool) Stats {
	return d.evaluateImpl(maxRuns, incomingDamage, rng, ev, stop)
}

// DefaultEquipment re-exports defaultEquipment for sim_test consumers.
var DefaultEquipment = defaultEquipment

// AdaptiveCheckInterval re-exports the adaptive-shuffle check interval constant.
const AdaptiveCheckInterval = adaptiveCheckInterval

// AdaptiveShufflesCap re-exports the adaptive-shuffle ceiling.
const AdaptiveShufflesCap = adaptiveShufflesCap

// MakeAdaptiveStop re-exports makeAdaptiveStop for sim_test consumers.
func MakeAdaptiveStop(targetSE float64) func(stats *Stats, runs int) bool {
	return makeAdaptiveStop(targetSE)
}

// MeanStandardError re-exports meanStandardError for sim_test consumers.
func MeanStandardError(stats *Stats) float64 { return meanStandardError(stats) }

// ProcessTriggersAtStartOfTurn re-exports processTriggersAtStartOfTurn for sim_test
// consumers exercising the start-of-turn aura-trigger pipeline in isolation.
func ProcessTriggersAtStartOfTurn(queued []AuraTrigger, postDrawDeck []Card) (
	survivors []AuraTrigger,
	contribs []TriggerContribution,
	damage int,
	runes int,
	revealed []Card,
	graveyarded []Card,
) {
	return processTriggersAtStartOfTurn(queued, postDrawDeck)
}
