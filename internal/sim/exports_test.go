package sim

import (
	"math/rand"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Test-only exports.

// Best re-exports the package-private best. state carries the matchup figures (set via
// SetIncomingDamage) and any other carryover.
func Best(weapons []weapon.Weapon, hand []card.Card, d *deck.Deck, state *gameengine.GameState) TurnSummary {
	return best(weapons, hand, d, state)
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
func (s *SequenceContextForTest) PlaySequence(order []card.Card) (damage int, totalCounters int, residualBudget int, legal bool) {
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

// BestSequence wraps (*sequenceContext).bestSequence and returns the winning leaf's
// damage / totalCounters / legal triplet.
func (s *SequenceContextForTest) BestSequence(attackers []card.Card) (int, int, bool) {
	score, _, ok := s.ctx.bestSequence(attackers)
	return score.value, score.totalCounters, ok
}

// FireEndOfTurn re-exports the engine's end-of-turn fire for sim_test consumers.
func FireEndOfTurn(ge *gameengine.GameEngine) { ge.FireTriggers(triggertype.EndOfTurn, nil) }

// PromoteRandomHandCardToArsenal re-exports promoteRandomHandCardToArsenal.
func PromoteRandomHandCardToArsenal(best *TurnSummary, startingHand []card.Card, arsenalCardIn card.Card) {
	promoteRandomHandCardToArsenal(best, startingHand, arsenalCardIn)
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
// tuple. damage is gs.Value() post-call, since handlers write directly to it. Cards added
// to state.Hand during firing (e.g. reveal-into-hand handlers) come back as `revealed`.
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
	preHand := len(gs.Hand())
	processAurasAtStartOfTurn(gs, d)
	damage = gs.Value()
	if newGrav := gs.Graveyard(); len(newGrav) > preGrav {
		graveyarded = append([]card.Card(nil), newGrav[preGrav:]...)
	}
	if newHand := gs.Hand(); len(newHand) > preHand {
		revealed = append([]card.Card(nil), newHand[preHand:]...)
	}
	survivors = append([]gameengine.Aura(nil), gs.Auras()...)
	return survivors, damage, revealed, graveyarded
}
