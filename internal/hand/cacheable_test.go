package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/generic"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
)

// End-to-end tests for TurnSummary.Cacheable: assert that running hand.Best with various
// inputs produces the right cacheability bit on the returned summary. Pins the propagation
// from per-permutation TurnState.IsCacheable() up through bestSequence → bestAttackWithWeapons
// → findBest → TurnSummary, plus aggregation across all partition leaves (any chain that
// reads deck or graveyard poisons the whole hand-eval result).

// TestBest_CacheableEmptyHand: no attackers, no chain runs at all. The "nothing was played"
// fallback partition produces Value=0 and never enters any Card.Play, so the hand-eval is
// trivially cacheable.
func TestBest_CacheableEmptyHand(t *testing.T) {
	got := Best(stubHero, nil, nil, 0, nil, 0, nil)
	if !got.Cacheable {
		t.Errorf("empty hand should be cacheable; got Cacheable=false")
	}
}

// TestBest_CacheablePlainAttackers: a hand of plain RedAttack stubs whose Play does no deck
// or graveyard reads. Every permutation runs without flipping IsCacheable, so the summary
// reports Cacheable=true.
func TestBest_CacheablePlainAttackers(t *testing.T) {
	h := []card.Card{fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	if !got.Cacheable {
		t.Errorf("plain attackers should be cacheable; got Cacheable=false")
	}
}

// TestBest_CacheableMixedNonReaders: a hand mixing pitched + attackers + a defense-reaction
// blocker, none of which read deck or graveyard, must remain cacheable.
func TestBest_CacheableMixedNonReaders(t *testing.T) {
	h := []card.Card{
		runeblade.MaleficIncantationBlue{}, // pitch-3, no deck/graveyard reads
		generic.ToughenUpBlue{},            // DR with no deck/graveyard reads
		fake.RedAttack{},
	}
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	if !got.Cacheable {
		t.Errorf("non-reader mix should be cacheable; got Cacheable=false")
	}
}

// TestBest_UncacheableSkyFireLanterns: Sky Fire Lanterns' Play peeks the deck top via
// s.Deck() to compare its pitch against the variant's color. The partition enumerator
// tries the Attack-role line for SkyFire, which fires its Play, which flips
// IsCacheable=false on at least one permutation — propagating up to summary.Cacheable.
func TestBest_UncacheableSkyFireLanterns(t *testing.T) {
	h := []card.Card{
		runeblade.SkyFireLanternsRed{},
		runeblade.MaleficIncantationBlue{}, // pitch fuel so SkyFire can be played
		fake.RedAttack{},
	}
	deck := []card.Card{fake.RedAttack{}}
	got := Best(stubHero, nil, h, 4, deck, 0, nil)
	if got.Cacheable {
		t.Errorf("Sky Fire Lanterns reads deck top; expected Cacheable=false")
	}
}

// TestBest_UncacheableSutcliffes: Sutcliffe's Research Notes' Play scans the top N cards of
// the deck via s.Deck(). Same propagation path as SkyFire.
func TestBest_UncacheableSutcliffes(t *testing.T) {
	h := []card.Card{
		runeblade.SutcliffesResearchNotesRed{},
		runeblade.MaleficIncantationBlue{},
		runeblade.MaleficIncantationBlue{},
	}
	deck := []card.Card{fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	got := Best(stubHero, nil, h, 4, deck, 0, nil)
	if got.Cacheable {
		t.Errorf("Sutcliffe's reads top-N of deck; expected Cacheable=false")
	}
}

// TestBest_UncacheableMoonWish: Moon Wish reads the deck via s.Deck() to tutor a Sun Kiss.
// The partition enumerator's Attack-role exploration fires Play, flipping IsCacheable.
func TestBest_UncacheableMoonWish(t *testing.T) {
	h := []card.Card{
		generic.MoonWishYellow{},
		runeblade.MaleficIncantationBlue{},
	}
	deck := []card.Card{generic.SunKissRed{}, fake.RedAttack{}}
	got := Best(stubHero, nil, h, 4, deck, 0, nil)
	if got.Cacheable {
		t.Errorf("Moon Wish tutors via deck scan; expected Cacheable=false")
	}
}

// TestBest_UncacheableRavenousRabble: Ravenous Rabble peeks deck top via s.Deck() to
// subtract its pitch from the printed power.
func TestBest_UncacheableRavenousRabble(t *testing.T) {
	h := []card.Card{
		generic.RavenousRabbleRed{},
		runeblade.MaleficIncantationBlue{},
		fake.RedAttack{},
	}
	deck := []card.Card{fake.RedAttack{}}
	got := Best(stubHero, nil, h, 4, deck, 0, nil)
	if got.Cacheable {
		t.Errorf("Ravenous Rabble reads deck top; expected Cacheable=false")
	}
}

// TestBest_UncacheableSnatchViaDrawOne: Snatch's Play (on hit) calls s.DrawOne, which
// reads s.Deck() through the framework helper. Cards inheriting the flip via DrawOne /
// ClashValue don't need any per-card change — that's the structural-enforcement guarantee.
func TestBest_UncacheableSnatchViaDrawOne(t *testing.T) {
	h := []card.Card{
		generic.SnatchRed{},
		runeblade.MaleficIncantationBlue{}, // pitch fuel
	}
	deck := []card.Card{fake.RedAttack{}}
	got := Best(stubHero, nil, h, 4, deck, 0, nil)
	if got.Cacheable {
		t.Errorf("Snatch (DrawOne on hit) reads deck top via framework helper; expected Cacheable=false")
	}
}

// TestBest_UncacheableTestOfStrengthViaClashValue: Test of Strength clashes via
// card.ClashValue, which reads s.Deck() to compare against the opponent's projected top.
func TestBest_UncacheableTestOfStrengthViaClashValue(t *testing.T) {
	h := []card.Card{
		generic.TestOfStrengthRed{},
		runeblade.MaleficIncantationBlue{},
	}
	deck := []card.Card{fake.RedAttack{}}
	got := Best(stubHero, nil, h, 4, deck, 0, nil)
	if got.Cacheable {
		t.Errorf("Test of Strength clashes via ClashValue; expected Cacheable=false")
	}
}

// TestBest_UncacheableWeepingBattlegroundViaPriorGraveyard: Weeping Battleground's Play
// scans s.Graveyard() for an aura to banish. The framework's per-played-card append to
// the graveyard doesn't seed the start-of-chain graveyard (that's empty); the scan still
// reads the live slice, flipping IsCacheable. Even when the scan finds nothing, the read
// itself is what flips — semantics match aura_banish iterating the slice.
func TestBest_UncacheableWeepingBattlegroundViaPriorGraveyard(t *testing.T) {
	h := []card.Card{
		runeblade.WeepingBattlegroundRed{},
		runeblade.MaleficIncantationBlue{}, // pitch fuel
		fake.RedAttack{},
	}
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	if got.Cacheable {
		t.Errorf("Weeping Battleground scans graveyard; expected Cacheable=false")
	}
}

// TestBest_CacheabilityIsAggregatedAcrossPartitions: a hand with one deck-reader plus
// several non-readers must surface Cacheable=false even when the WINNING partition is one
// where the deck-reader was pitched (not played). The aggregation in findBest poisons the
// summary the moment ANY explored partition flips IsCacheable. This protects a future cache
// from storing one partition's cacheable result against a key that another shuffle could
// hit and resolve to a different (uncacheable) winner.
func TestBest_CacheabilityIsAggregatedAcrossPartitions(t *testing.T) {
	h := []card.Card{
		runeblade.SkyFireLanternsRed{},
		fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{},
	}
	deck := []card.Card{fake.RedAttack{}}
	got := Best(stubHero, nil, h, 4, deck, 0, nil)
	if got.Cacheable {
		t.Errorf("SkyFire in hand should mark uncacheable even if pitched in winning line; got Cacheable=true (BestLine=%s)",
			FormatBestLine(got.BestLine))
	}
}

// TestBest_CacheableSnatchYellowDoesntHit: Yellow Snatch (attack 3) is in the
// uncacheableCards family — Snatch's Play calls DrawOne when LikelyToHit. But
// LikelyDamageHits(n, false) is true only for n ∈ {1, 4, 7}, so attack 3 doesn't hit and
// the DrawOne branch is skipped at runtime. With Snatch alone in hand, every partition
// path either runs Play-without-DrawOne (Attack role) or runs no Play at all (Pitch /
// Defend-as-plain-block / Held). Result must be cacheable.
//
// This pins the structural-enforcement principle's payoff: cacheability tracks ACTUAL
// reads, not card identity. A card "in the uncacheableCards family" doesn't poison the
// chain unless its Play actually fires the read on this run.
func TestBest_CacheableSnatchYellowDoesntHit(t *testing.T) {
	h := []card.Card{generic.SnatchYellow{}}
	got := Best(stubHero, nil, h, 0, nil, 0, nil)
	if !got.Cacheable {
		t.Errorf("Yellow Snatch (attack 3) skips DrawOne via LikelyToHit gate; expected Cacheable=true (BestLine=%s)",
			FormatBestLine(got.BestLine))
	}
}

// TestBest_CacheableSnatchBlueWithRedAttacker: Blue Snatch (attack 2) misses on every
// partition where it plays as Attack, and never runs Play when pitched / blocked / held.
// The companion Red attacker (fake.RedAttack, attack 3) is a non-reader stub. All paths
// cacheable.
func TestBest_CacheableSnatchBlueWithRedAttacker(t *testing.T) {
	h := []card.Card{generic.SnatchBlue{}, fake.RedAttack{}}
	got := Best(stubHero, nil, h, 1, nil, 0, nil)
	if !got.Cacheable {
		t.Errorf("Blue Snatch (attack 2, never hits) + plain Red attacker should be cacheable; got Cacheable=false (BestLine=%s)",
			FormatBestLine(got.BestLine))
	}
}

// TestBest_CacheableMoonWishBlockedAtCostCheck: Red Moon Wish alone in hand can never
// resolve as Attack — its 2-cost can't be paid (no other card to pitch, no other card to
// alt-cost), and the playSequenceWithMeta cost check rejects the permutation BEFORE Play
// runs. So no deck read fires from the tutor branch. The other partition roles never run
// Play either: Pitch (no Play), plain Defend against incoming > 0 (Moon Wish isn't a DR;
// plain blocks credit Defense() without Play), Held / Arsenal (no Play). All paths
// cacheable.
//
// Pins that the pre-Play resource gate is sufficient to keep would-be-uncacheable cards
// from poisoning a chain when partition enumeration would otherwise try them.
func TestBest_CacheableMoonWishBlockedAtCostCheck(t *testing.T) {
	h := []card.Card{generic.MoonWishRed{}}
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	if !got.Cacheable {
		t.Errorf("Red Moon Wish alone (cost 2, no payment source) should be cacheable; got Cacheable=false (BestLine=%s)",
			FormatBestLine(got.BestLine))
	}
}

// TestBest_CacheabilityResetsAcrossCalls: a single Evaluator reused across two calls — first
// with a deck-reader, then without — must report the second call as cacheable. The
// per-permutation TurnState reset zeroes the bit for free; if the reset path leaked the
// previous call's poisoned bit, the second call would falsely report uncacheable.
func TestBest_CacheabilityResetsAcrossCalls(t *testing.T) {
	e := NewEvaluator()

	// First call: deck-reader present → uncacheable.
	h1 := []card.Card{
		runeblade.SkyFireLanternsRed{},
		runeblade.MaleficIncantationBlue{},
		fake.RedAttack{},
	}
	deck := []card.Card{fake.RedAttack{}}
	got1 := e.BestWithTriggers(stubHero, nil, h1, 4, deck, 0, nil, nil)
	if got1.Cacheable {
		t.Fatalf("first call: expected Cacheable=false")
	}

	// Second call: no readers → cacheable. If the bit leaked from call 1 we'd see false.
	h2 := []card.Card{fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	got2 := e.BestWithTriggers(stubHero, nil, h2, 4, nil, 0, nil, nil)
	if !got2.Cacheable {
		t.Errorf("second call: expected Cacheable=true; cacheable bit leaked from prior call")
	}
}
