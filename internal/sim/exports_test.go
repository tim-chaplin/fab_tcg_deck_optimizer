package sim

import (
	"math/rand"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/item"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Test-only exports.

// Best re-exports the package-private best. state is the carryover *GameState — pass nil
// to start from a clean state.
func Best(weapons []weapon.Weapon, hand []card.Card, mp Matchup, d *deck.Deck, state *gameengine.GameState) TurnSummary {
	return best(weapons, hand, mp, d, state)
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
// bestSequence without the unexported type.
type SequenceContextForTest struct{ ctx *sequenceContext }

// NewSequenceContextForTest builds a sequenceContext via newSequenceContextForTest.
func NewSequenceContextForTest(h hero.Hero, pitched, deck []card.Card, resourceBudget, runechantCarryover, chainLen int) *SequenceContextForTest {
	return &SequenceContextForTest{ctx: newSequenceContextForTest(h, pitched, deck, resourceBudget, runechantCarryover, chainLen)}
}

// PlaySequence wraps (*sequenceContext).playSequence.
func (s *SequenceContextForTest) PlaySequence(order []card.Card) (damage int, futureValue int, residualBudget int, legal bool) {
	return s.ctx.playSequence(order)
}

// PermEngine returns a *GameEngine wrapping the *GameState the most recent PlaySequence
// call ran the chain against.
func (s *SequenceContextForTest) PermEngine() *gameengine.GameEngine {
	if s.ctx.permState == nil {
		return nil
	}
	return s.ctx.permState.Engine()
}

// BestSequence wraps (*sequenceContext).bestSequence and returns only the
// damage / future-value / legal triplet.
func (s *SequenceContextForTest) BestSequence(attackers []card.Card) (int, int, bool) {
	d, fv, _, ok := s.ctx.bestSequence(attackers)
	return d, fv, ok
}

// FireEndOfTurn re-exports the engine's end-of-turn fire for sim_test consumers.
func FireEndOfTurn(ge *gameengine.GameEngine) { ge.FireEndOfTurn() }

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

// DefendersDamage re-exports defendersDamage with an unbounded block budget.
func DefendersDamage(defenders, pitched []card.Card, d *deck.Deck, ge *gameengine.GameEngine, gravBuf []card.Card, cs *card.CardState, incomingDamage, arsenalDefenderIdx int) (int, []card.Card) {
	total, gravBuf, _ := defendersDamage(defenders, pitched, d, ge, gravBuf, cs, incomingDamage, noBlockBudgetCap, arsenalDefenderIdx)
	return total, gravBuf
}

// DefendersDamageWithBudget is the budget-aware export.
func DefendersDamageWithBudget(defenders, pitched []card.Card, d *deck.Deck, ge *gameengine.GameEngine, gravBuf []card.Card, cs *card.CardState, incomingDamage, blockBudget, arsenalDefenderIdx int) (int, []card.Card) {
	total, gravBuf, _ := defendersDamage(defenders, pitched, d, ge, gravBuf, cs, incomingDamage, blockBudget, arsenalDefenderIdx)
	return total, gravBuf
}

// AttackBufs is the exported alias of attackBufs.
type AttackBufs = attackBufs

// NewAttackBufs re-exports newAttackBufs.
func NewAttackBufs(handSize, weaponCount int, weapons []weapon.Weapon) *AttackBufs {
	return newAttackBufs(handSize, weaponCount, weapons)
}

// Bufs returns the wrapped sequenceContext's pooled scratch buffers.
func (s *SequenceContextForTest) Bufs() *AttackBufs { return s.ctx.bufs }

// DefenseGravScratch / DRCardStateScratch expose unexported attackBufs fields. State()
// returns a fresh engine per call; the chain runner uses per-permutation Copy so attackBufs
// doesn't pool one.
func (b *attackBufs) DefenseGravScratch() []card.Card     { return b.defenseGravScratch }
func (b *attackBufs) DRCardStateScratch() *card.CardState { return &b.drCardStateScratch }
func (b *attackBufs) State() *gameengine.GameEngine {
	return gameengine.New()
}

// EngineWithHand returns a fresh GameState seeded with hand h.
func EngineWithHand(h []card.Card) *gameengine.GameState {
	gs := gameengine.GameStateBuilder().Build()
	gs.SetHand(h)
	return gs
}

// EngineWithItems returns a fresh GameState with the supplied items installed.
func EngineWithItems(items []*item.Item) *gameengine.GameState {
	gs := gameengine.GameStateBuilder().Build()
	for _, it := range items {
		gs.CreateItem(it)
	}
	return gs
}

// EvaluateImplForTest re-exports the unexported (*Evaluator).evaluateImpl.
func (ev *Evaluator) EvaluateImplForTest(d *deck.Deck, maxRuns int, mp Matchup, rng *rand.Rand, stop func(stats *deck.Stats, runs int) bool) deck.Stats {
	return ev.evaluateImpl(d, maxRuns, mp, rng, stop)
}

// AdaptiveCheckInterval re-exports the adaptive-shuffle check interval constant.
const AdaptiveCheckInterval = adaptiveCheckInterval

// AdaptiveShufflesCap re-exports the adaptive-shuffle ceiling.
const AdaptiveShufflesCap = adaptiveShufflesCap

// MakeAdaptiveStop re-exports makeAdaptiveStop.
func MakeAdaptiveStop(targetSE float64) func(stats *deck.Stats, runs int) bool {
	return makeAdaptiveStop(targetSE)
}

// MeanStandardError re-exports meanStandardError.
func MeanStandardError(stats *deck.Stats) float64 { return meanStandardError(stats) }

// ProcessAurasAtStartOfTurnForTest drives processAurasAtStartOfTurn against an arbitrary
// aura queue and returns the post-tick (survivors, damage, revealedCards, graveyardedCards)
// tuple.
func ProcessAurasAtStartOfTurnForTest(queued []gameengine.Aura, d *deck.Deck) (
	survivors []gameengine.Aura,
	damage int,
	revealed []card.Card,
	graveyarded []card.Card,
) {
	gs := gameengine.GameStateBuilder().Build()
	for _, a := range queued {
		gs.CreateAura(a)
	}
	preGrav := len(gs.Graveyard())
	var revealedBuf []card.Card
	damage = processAurasAtStartOfTurn(gs, d, &revealedBuf)
	if newGrav := gs.Graveyard(); len(newGrav) > preGrav {
		graveyarded = append([]card.Card(nil), newGrav[preGrav:]...)
	}
	survivors = append([]gameengine.Aura(nil), gs.Auras()...)
	return survivors, damage, revealedBuf, graveyarded
}
