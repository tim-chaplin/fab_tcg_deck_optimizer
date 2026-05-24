package sim

// End-to-end tests for TurnSummary.Cacheable propagation through hand.Best. The bit
// reports whether the chain that produced the summary depended on hidden state (deck
// order, prior-turn graveyard contents). The structural-enforcement principle: cards
// only reach deck / graveyard via accessor methods on TurnState (the underlying fields
// are package-private), so the cacheable signal is sound by the language's visibility
// rules — no cardlint backstop required.

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestBest_CacheableEmptyHand: no chain ran, nothing read hidden state — Cacheable=true.
// Pins the no-feasible-line fallback's seed default so an empty hand doesn't accidentally
// report uncacheable just because the search visited zero leaves.
func TestBest_CacheableEmptyHand(t *testing.T) {
	got := Best(nil, nil, nil, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	if !got.Cacheable {
		t.Errorf("empty hand: Cacheable = false, want true (no chain, no hidden read)")
	}
}

// TestBest_CacheablePlainAttackers: a hand of plain attack-action fakes whose Play does no
// deck / graveyard reads. The chain's output depends only on the inputs (hand + cards
// played), so Cacheable=true.
func TestBest_CacheablePlainAttackers(t *testing.T) {
	h := []card.Card{testutils.FakeRedAttack(), testutils.FakeRedAttack(), testutils.FakeRedAttack(), testutils.FakeRedAttack()}
	got := Best(nil, h, nil, gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetIncomingDamage(4).
		Build())
	if !got.Cacheable {
		t.Errorf("plain attackers: Cacheable = false, want true (no card touches hidden state)")
	}
}

// TestBest_UncacheableSkyFireLanterns: Sky Fire Lanterns peeks the deck top via
// ge.PeekDeck() to gate its runechant rider. Even when the rider doesn't fire, the read
// happened — Cacheable must report false.
func TestBest_UncacheableSkyFireLanterns(t *testing.T) {
	h := []card.Card{cards.SkyFireLanternsRed{}}
	deck := DeckOf(testutils.FakeRedAttack())
	got := Best(nil, h, deck, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	if got.Cacheable {
		t.Errorf("Sky Fire Lanterns hand: Cacheable = true, want false (Play reads s.PeekDeck())")
	}
}

// Tests that Sutcliffe's Research Notes pins Cacheable=false because its top-N scan
// reads ge.PeekTopN().
func TestBest_UncacheableSutcliffesResearchNotes(t *testing.T) {
	h := []card.Card{cards.SutcliffesResearchNotesRed{}, testutils.FakeBlueAttack()}
	deck := DeckOf(testutils.FakeRedAttack().WithTypes(card.TypeRuneblade))
	got := Best(nil, h, deck, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	if got.Cacheable {
		t.Errorf("Sutcliffe's hand: Cacheable = true, want false (Play scans s.PeekTopN())")
	}
}

// TestBest_UncacheableMoonWishTutor: Moon Wish's on-hit branch tutors via TutorFromDeck —
// flips the bit through the verb. Pair it with a hand card so the alt cost fires (still
// flips via PrependToDeck) and a deck card to tutor.
func TestBest_UncacheableMoonWishTutor(t *testing.T) {
	h := []card.Card{cards.MoonWishRed{}, testutils.FakeRedAttack()}
	deck := DeckOf(cards.SunKissRed{})
	got := Best(nil, h, deck, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	if got.Cacheable {
		t.Errorf("Moon Wish hand: Cacheable = true, want false (TutorFromDeck flips)")
	}
}

// TestBest_UncacheableRavenousRabble: the on-attack -X{p} debuff reads the deck top via
// ge.PeekDeck() — Cacheable=false even though the card "only" peeks.
func TestBest_UncacheableRavenousRabble(t *testing.T) {
	h := []card.Card{cards.RavenousRabbleRed{}}
	got := Best(nil, h, DeckOf(testutils.FakeRedAttack()), gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	if got.Cacheable {
		t.Errorf("Ravenous Rabble hand: Cacheable = true, want false (Play reads s.PeekDeck())")
	}
}

// Tests that a Snatch hit's on-hit DrawOne leaves the chain Cacheable.
func TestBest_CacheableSnatchHitDrawsViaDrawOne(t *testing.T) {
	h := []card.Card{cards.SnatchRed{}}
	deck := DeckOf(testutils.FakeRedAttack())
	got := Best(nil, h, deck, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	if !got.Cacheable {
		t.Errorf("Snatch [R] hand: Cacheable = false, want true (DrawOne no longer flips cacheable)")
	}
}

// Tests that Test of Strength's Clash flips Cacheable via ge.Clash().
func TestBest_UncacheableTestOfStrengthClash(t *testing.T) {
	h := []card.Card{cards.TestOfStrengthRed{}}
	got := Best(nil, h, DeckOf(testutils.FakeRedAttack()), gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetIncomingDamage(1).
		Build())
	if got.Cacheable {
		t.Errorf("Test of Strength hand: Cacheable = true, want false (Clash reads deck top)")
	}
}

// Tests that Weeping Battleground's BanishFromGraveyard call propagates the Cacheable flip
// through defendersDamage.
func TestBest_UncacheableWeepingBattlegroundDR(t *testing.T) {
	h := []card.Card{cards.WeepingBattlegroundRed{}, cards.MaleficIncantationBlue{}}
	got := Best(nil, h, nil, gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetIncomingDamage(3).
		Build())
	if got.Cacheable {
		t.Errorf("Weeping Battleground hand: Cacheable = true, want false (BanishFromGraveyard flips)")
	}
}

// Tests that any sibling partition's hidden-state read pins Cacheable=false even when the
// deck-reading card is only pitched in the winning line.
func TestBest_AggregationDeckReaderInHandPoisonsResultEvenWhenPitched(t *testing.T) {
	h := []card.Card{cards.SkyFireLanternsBlue{}, testutils.FakeRedAttack(), testutils.FakeRedAttack(), testutils.FakeRedAttack()}
	deck := DeckOf(testutils.FakeRedAttack())
	got := Best(nil, h, deck, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	if got.Cacheable {
		t.Errorf("hand with pitched deck-reader: Cacheable = true, want false (sibling-leaf reads poison)")
	}
}

// Tests that an Evaluator's uncacheable bit resets between Best calls so a prior uncacheable
// hand can't poison a later cacheable one.
func TestBest_ResetBetweenCallsClearsCacheableState(t *testing.T) {
	ev := NewEvaluator()
	deck := DeckOf(testutils.FakeRedAttack())

	// First call: Sky Fire Lanterns reads the deck top → expected Cacheable=false.
	first := ev.Best(nil, []card.Card{cards.SkyFireLanternsRed{}}, deck, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	if first.Cacheable {
		t.Fatalf("first call: Cacheable = true, want false (Sky Fire reads deck)")
	}

	// Second call: plain attackers, no hidden read → cacheable. If the bit leaked across
	// calls this assertion would fail.
	clean := ev.Best(nil, []card.Card{testutils.FakeRedAttack()}, deck, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	if !clean.Cacheable {
		t.Errorf("second call: Cacheable = false, want true (no hidden read; bit must reset between calls)")
	}
}

// Tests that cacheability tracks ACTUAL reads, not card identity: Snatch Yellow misses
// ge.LikelyToHit so DrawOne never fires and the result stays cacheable.
func TestBest_RuntimeGatedNonFlipSnatchYellowMisses(t *testing.T) {
	h := []card.Card{cards.SnatchYellow{}}
	deck := DeckOf(testutils.FakeRedAttack())
	got := Best(nil, h, deck, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	if !got.Cacheable {
		t.Errorf("Snatch [Y] alone: Cacheable = false, want true (s.LikelyToHit miss skips DrawOne)")
	}
}

// Tests that pre-Play cost rejection stops a would-be-uncacheable card from running Play.
func TestBest_RuntimeGatedNonFlipMoonWishBlockedAtCostCheck(t *testing.T) {
	h := []card.Card{cards.MoonWishRed{}}
	deck := DeckOf(cards.SunKissRed{})
	got := Best(nil, h, deck, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	if !got.Cacheable {
		t.Errorf("solo Moon Wish [R]: Cacheable = false, want true (cost check rejects before Play)")
	}
}
