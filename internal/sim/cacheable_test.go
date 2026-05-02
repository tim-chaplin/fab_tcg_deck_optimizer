package sim_test

// End-to-end tests for TurnSummary.Cacheable propagation through hand.Best. The bit
// reports whether the chain that produced the summary depended on hidden state (deck
// order, prior-turn graveyard contents). The structural-enforcement principle: cards
// only reach deck / graveyard via accessor methods on TurnState (the underlying fields
// are package-private), so the cacheable signal is sound by the language's visibility
// rules — no cardlint backstop required.

import (
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	notimpl "github.com/tim-chaplin/fab-deck-optimizer/internal/cards/notimplemented"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestBest_CacheableEmptyHand: no chain ran, nothing read hidden state — Cacheable=true.
// Pins the no-feasible-line fallback's seed default so an empty hand doesn't accidentally
// report uncacheable just because the search visited zero leaves.
func TestBest_CacheableEmptyHand(t *testing.T) {
	got := Best(testutils.Hero{Intel: 4}, nil, nil, Matchup{IncomingDamage: 0}, nil, 0, nil)
	if !got.Cacheable {
		t.Errorf("empty hand: Cacheable = false, want true (no chain, no hidden read)")
	}
}

// TestBest_CacheablePlainAttackers: a hand of plain attack-action fakes whose Play does no
// deck / graveyard reads. The chain's output depends only on the inputs (hand + cards
// played), so Cacheable=true.
func TestBest_CacheablePlainAttackers(t *testing.T) {
	h := []Card{testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}}
	got := Best(testutils.Hero{Intel: 4}, nil, h, Matchup{IncomingDamage: 4}, nil, 0, nil)
	if !got.Cacheable {
		t.Errorf("plain attackers: Cacheable = false, want true (no card touches hidden state)")
	}
}

// TestBest_UncacheableSkyFireLanterns: Sky Fire Lanterns peeks the deck top via s.Deck() to
// gate its runechant rider. Even when the rider doesn't fire, the read happened — Cacheable
// must report false.
func TestBest_UncacheableSkyFireLanterns(t *testing.T) {
	h := []Card{cards.SkyFireLanternsRed{}}
	deck := []Card{testutils.RedAttack{}}
	got := Best(testutils.Hero{Intel: 4}, nil, h, Matchup{IncomingDamage: 0}, deck, 0, nil)
	if got.Cacheable {
		t.Errorf("Sky Fire Lanterns hand: Cacheable = true, want false (Play reads s.Deck())")
	}
}

// TestBest_UncacheableSutcliffesResearchNotes: the top-N runeblade-attack scan calls
// s.Deck() so Cacheable=false for any chain that includes a Sutcliffe's printing. Need
// a second pitchable card to fund Sutcliffe's cost-1 — solo it's never feasible to play
// and Play never fires.
func TestBest_UncacheableSutcliffesResearchNotes(t *testing.T) {
	h := []Card{cards.SutcliffesResearchNotesRed{}, testutils.BlueAttack{}}
	deck := []Card{testutils.RunebladeAttack{}}
	got := Best(testutils.Hero{Intel: 4}, nil, h, Matchup{IncomingDamage: 0}, deck, 0, nil)
	if got.Cacheable {
		t.Errorf("Sutcliffe's hand: Cacheable = true, want false (Play scans s.Deck())")
	}
}

// TestBest_UncacheableMoonWishTutor: Moon Wish's on-hit branch tutors via TutorFromDeck —
// flips the bit through the verb. Pair it with a hand card so the alt cost fires (still
// flips via PrependToDeck) and a deck card to tutor.
func TestBest_UncacheableMoonWishTutor(t *testing.T) {
	h := []Card{cards.MoonWishRed{}, testutils.RedAttack{}}
	deck := []Card{cards.SunKissRed{}}
	got := Best(testutils.Hero{Intel: 4}, nil, h, Matchup{IncomingDamage: 0}, deck, 0, nil)
	if got.Cacheable {
		t.Errorf("Moon Wish hand: Cacheable = true, want false (TutorFromDeck flips)")
	}
}

// TestBest_UncacheableRavenousRabble: the on-attack -X{p} debuff reads the deck top via
// s.Deck() — Cacheable=false even though the card "only" peeks.
func TestBest_UncacheableRavenousRabble(t *testing.T) {
	h := []Card{cards.RavenousRabbleRed{}}
	deck := []Card{testutils.GenericAttackPitch(0, 0, 1)}
	got := Best(testutils.Hero{Intel: 4}, nil, h, Matchup{IncomingDamage: 0}, deck, 0, nil)
	if got.Cacheable {
		t.Errorf("Ravenous Rabble hand: Cacheable = true, want false (Play reads s.Deck())")
	}
}

// TestBest_UncacheableSnatchHitDrawsViaDrawOne: Snatch Red's printed 4 attack lands in the
// hit window, so its on-hit DrawOne fires — that's a PopDeckTop call and the chain's output
// depends on what got drawn. Cacheable=false even though Snatch never names s.Deck()
// directly; the framework helper inherits the flip.
func TestBest_UncacheableSnatchHitDrawsViaDrawOne(t *testing.T) {
	h := []Card{cards.SnatchRed{}}
	deck := []Card{testutils.RedAttack{}}
	got := Best(testutils.Hero{Intel: 4}, nil, h, Matchup{IncomingDamage: 0}, deck, 0, nil)
	if got.Cacheable {
		t.Errorf("Snatch [R] hand: Cacheable = true, want false (DrawOne flips via PopDeckTop)")
	}
}

// TestBest_UncacheableTestOfStrengthClash: Test of Strength's clash reads the deck top via
// ClashValue, which inherits the flip from s.Deck(). Incoming = 1 to give the partition an
// actual defend step where the DR can fire (the partition skips Defend assignments at 0
// incoming since FaB has no defense step without an attack).
func TestBest_UncacheableTestOfStrengthClash(t *testing.T) {
	h := []Card{notimpl.TestOfStrengthRed{}}
	deck := []Card{testutils.GenericAttack(0, 7)}
	got := Best(testutils.Hero{Intel: 4}, nil, h, Matchup{IncomingDamage: 1}, deck, 0, nil)
	if got.Cacheable {
		t.Errorf("Test of Strength hand: Cacheable = true, want false (ClashValue flips via Deck())")
	}
}

// TestBest_UncacheableWeepingBattlegroundDR: Weeping Battleground is a Defense Reaction; its
// banish-from-graveyard rider routes through BanishFromGraveyard and reads the graveyard
// contents. The DR Play runs even when the partition's chain rejects, so the leaf-level
// uncacheable bit pings the propagation through defendersDamage.
func TestBest_UncacheableWeepingBattlegroundDR(t *testing.T) {
	// Pair Weeping Battleground with a Malefic Incantation Blue (cost 0 / pitch 3) so the DR
	// has a pitched card funding its 0 cost. Incoming damage > 0 forces the partition to
	// actually run defenders.
	h := []Card{cards.WeepingBattlegroundRed{}, cards.MaleficIncantationBlue{}}
	got := Best(testutils.Hero{Intel: 4}, nil, h, Matchup{IncomingDamage: 3}, nil, 0, nil)
	if got.Cacheable {
		t.Errorf("Weeping Battleground hand: Cacheable = true, want false (BanishFromGraveyard flips)")
	}
}

// TestBest_AggregationDeckReaderInHandPoisonsResultEvenWhenPitched: a deck-reading card
// pitched in the winning line still pins Cacheable=false because the search explored every
// partition and the partition that played the deck-reader (a sibling, not the winner) ran a
// chain that touched hidden state. The cacheable signal is "did any leaf the search visited
// touch hidden state", not "did the winner".
func TestBest_AggregationDeckReaderInHandPoisonsResultEvenWhenPitched(t *testing.T) {
	// Sky Fire Lanterns Blue (pitch 3) + 3 Red Attacks lets the partition pitch the Sky Fire
	// to fund 3 Red Attacks (3 res for 3× cost-1). The winner doesn't play Sky Fire — but a
	// sibling partition that did would've called s.Deck(), pinning the leaf uncacheable.
	h := []Card{cards.SkyFireLanternsBlue{}, testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}}
	deck := []Card{testutils.RedAttack{}}
	got := Best(testutils.Hero{Intel: 4}, nil, h, Matchup{IncomingDamage: 0}, deck, 0, nil)
	if got.Cacheable {
		t.Errorf("hand with pitched deck-reader: Cacheable = true, want false (sibling-leaf reads poison)")
	}
}

// TestBest_ResetBetweenCallsClearsCacheableState: a single Evaluator reused across two
// independent Best calls must report each call's cacheable state on its own merits — the
// uncacheable bit must reset per call. Otherwise a deck-eval loop's first uncacheable hand
// would falsely pin every subsequent hand.
func TestBest_ResetBetweenCallsClearsCacheableState(t *testing.T) {
	ev := NewEvaluator()
	deck := []Card{testutils.RedAttack{}}

	// First call: Sky Fire Lanterns reads the deck top → expected Cacheable=false.
	first := ev.Best(testutils.Hero{Intel: 4}, nil, []Card{cards.SkyFireLanternsRed{}}, Matchup{IncomingDamage: 0}, deck, 0, nil)
	if first.Cacheable {
		t.Fatalf("first call: Cacheable = true, want false (Sky Fire reads deck)")
	}

	// Second call: plain attackers, no hidden read → cacheable. If the bit leaked across
	// calls this assertion would fail.
	clean := ev.Best(testutils.Hero{Intel: 4}, nil, []Card{testutils.RedAttack{}}, Matchup{IncomingDamage: 0}, deck, 0, nil)
	if !clean.Cacheable {
		t.Errorf("second call: Cacheable = false, want true (no hidden read; bit must reset between calls)")
	}
}

// TestBest_RuntimeGatedNonFlipSnatchYellowMisses: Snatch Yellow's printed attack 3 isn't in
// the LikelyToHit set, so the on-hit DrawOne branch never fires. With Snatch Yellow alone
// in hand every partition path either runs Play-without-DrawOne (Attack) or runs no Play at
// all (Pitch / Held). The deck is never read — Cacheable=true even though Snatch is "in the
// uncacheableCards family". Pins that cacheability tracks ACTUAL reads, not card identity.
func TestBest_RuntimeGatedNonFlipSnatchYellowMisses(t *testing.T) {
	h := []Card{cards.SnatchYellow{}}
	deck := []Card{testutils.RedAttack{}}
	got := Best(testutils.Hero{Intel: 4}, nil, h, Matchup{IncomingDamage: 0}, deck, 0, nil)
	if !got.Cacheable {
		t.Errorf("Snatch [Y] alone: Cacheable = false, want true (LikelyToHit miss skips DrawOne)")
	}
}

// TestBest_RuntimeGatedNonFlipMoonWishBlockedAtCostCheck: Red Moon Wish's printed cost is 2
// and the alt cost requires another hand card. Solo in hand it fails the cost check before
// Play runs — the tutor's TutorFromDeck never fires. Cacheable=true. Pins that pre-Play
// resource gates protect cacheability when partition enumeration tries would-be-uncacheable
// cards in infeasible roles.
func TestBest_RuntimeGatedNonFlipMoonWishBlockedAtCostCheck(t *testing.T) {
	h := []Card{cards.MoonWishRed{}}
	deck := []Card{cards.SunKissRed{}}
	got := Best(testutils.Hero{Intel: 4}, nil, h, Matchup{IncomingDamage: 0}, deck, 0, nil)
	if !got.Cacheable {
		t.Errorf("solo Moon Wish [R]: Cacheable = false, want true (cost check rejects before Play)")
	}
}
