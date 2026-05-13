package sim

import (
	"math/rand"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/token"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// Test-only exports. Visible to package sim_test files in this directory only.

// Best re-exports the package-private best for sim_test consumers.
func Best(weapons []Weapon, hand []card.Card, mp Matchup, d *deck.Deck, prior Prior) TurnSummary {
	return best(weapons, hand, mp, d, prior)
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
func NewSequenceContextForTest(h hero.Hero, pitched, deck []card.Card, resourceBudget, runechantCarryover, chainLen int) *SequenceContextForTest {
	return &SequenceContextForTest{ctx: newSequenceContextForTest(h, pitched, deck, resourceBudget, runechantCarryover, chainLen)}
}

// PlaySequence wraps (*sequenceContext).playSequence.
func (s *SequenceContextForTest) PlaySequence(order []card.Card) (damage int, futureValue int, residualBudget int, legal bool) {
	return s.ctx.playSequence(order)
}

// PermEngine returns a *GameEngine wrapping the *GameState the most recent
// PlaySequence call ran the chain against. Tests assert state via this engine
// (Graveyard, Hand, …) after PlaySequence.
func (s *SequenceContextForTest) PermEngine() *gameengine.GameEngine {
	if s.ctx.permState == nil {
		return nil
	}
	return s.ctx.permState.Engine()
}

// BestSequence wraps (*sequenceContext).bestSequence. Drops the returned winning engine
// pointer — sim_test consumers care only about the damage / future-value / legal triplet.
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

// AppendGroupedChainEntries re-exports appendGroupedChainEntries.
func AppendGroupedChainEntries(out []string, log []turnlogger.LogEntry) []string {
	return appendGroupedChainEntries(out, log)
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

// DefenseGravScratch / DRCardStateScratch expose the unexported attackBufs fields. The
// engine itself is no longer pooled on attackBufs (per-permutation Copy means each chain
// run owns its own); tests that need a chain-runner engine call State() — each call
// returns a fresh engine.
func (b *attackBufs) DefenseGravScratch() []card.Card     { return b.defenseGravScratch }
func (b *attackBufs) DRCardStateScratch() *card.CardState { return &b.drCardStateScratch }
func (b *attackBufs) State() *gameengine.GameEngine {
	return gameengine.New()
}

// EngineWithHand returns a fresh GameState seeded with hand h. Tests that build a
// TurnSummary by hand use this to populate the State *GameState without going through
// the full chain runner.
func EngineWithHand(h []card.Card) *gameengine.GameState {
	gs := gameengine.GameStateBuilder().Build()
	gs.SetHand(h)
	return gs
}

// EngineWithItems returns a fresh GameState with the supplied items installed.
func EngineWithItems(items []*token.Item) *gameengine.GameState {
	gs := gameengine.GameStateBuilder().Build()
	for _, it := range items {
		gs.CreateItem(it)
	}
	return gs
}

// EngineWith returns a fresh GameState with hand, items, and log entries installed.
// log can be nil to skip log seeding.
func EngineWith(h []card.Card, items []*token.Item, log []turnlogger.LogEntry) *gameengine.GameState {
	gs := gameengine.GameStateBuilder().Build()
	gs.SetHand(h)
	for _, it := range items {
		gs.CreateItem(it)
	}
	if len(log) > 0 {
		for _, e := range log {
			switch e.Kind {
			case turnlogger.LogEntryChainStep:
				gs.Logger().AppendChainStep(e.Text, e.N)
			case turnlogger.LogEntryPostTrigger:
				gs.Logger().AppendPostTrigger(e.Source, e.Text, e.N)
			case turnlogger.LogEntryPreTrigger:
				gs.Logger().AppendPreTrigger(e.Source, e.Text, e.N)
			}
		}
	}
	return gs
}

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

// ProcessAurasAtStartOfTurn re-exports processAurasAtStartOfTurn. Takes / returns
// gameengine.Aura at the boundary so sim_test callers (in package sim_test, outside
// sim's concrete-type namespace) keep working through the engine interface.
func ProcessAurasAtStartOfTurn(queued []gameengine.Aura, d *deck.Deck) (
	survivors []gameengine.Aura,
	contribs []TriggerContribution,
	damage int,
	revealed []card.Card,
	graveyarded []card.Card,
) {
	inAuras := make([]*aura.Aura, 0, len(queued))
	for _, a := range queued {
		inAuras = append(inAuras, a.(*aura.Aura))
	}
	s, c, dmg, rev, gv := processAurasAtStartOfTurn(inAuras, d)
	out := make([]gameengine.Aura, 0, len(s))
	for _, a := range s {
		out = append(out, a)
	}
	return out, c, dmg, rev, gv
}
