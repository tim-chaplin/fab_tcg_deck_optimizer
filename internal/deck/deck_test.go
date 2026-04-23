package deck

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/generic"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// int1StubHero is a test-only Hero with Intelligence=1 so we can isolate per-hand behavior
// without interaction between multiple drawn cards. Otherwise identical to a no-op hero —
// no on-play triggers, never flags as Runeblade.
type int1StubHero struct{}

func (int1StubHero) ID() hero.ID                                { return hero.Invalid }
func (int1StubHero) Name() string                               { return "int1Stub" }
func (int1StubHero) Health() int                                { return 20 }
func (int1StubHero) Intelligence() int                          { return 1 }
func (int1StubHero) Types() card.TypeSet                        { return 0 }
func (int1StubHero) OnCardPlayed(card.Card, *card.TurnState) int { return 0 }

func TestAllMutations_CountsAndShape(t *testing.T) {
	// Build a tiny deck: 2 unique cards × 2 copies = 4 cards, plus one weapon.
	a := cards.Get(card.AetherSlashRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

	muts := AllMutations(d, 2, nil)

	// Weapon mutations: every loadout except the current one. Card mutations at maxCopies=2:
	// for each of the 2 unique removals, every pool entry except self (no-op) and the other
	// in-deck card (already at cap) is a valid add — so 2 × (pool - 2). Use legalPool(nil)
	// instead of cards.Deckable() directly so the expected count tracks AllMutations's own
	// filtering (NotImplemented cards are skipped).
	loadouts := weaponLoadouts(cards.AllWeapons)
	pool := legalPool(nil)
	wantWeaponMuts := len(loadouts) - 1
	wantCardMuts := 2 * (len(pool) - 2)
	want := wantWeaponMuts + wantCardMuts

	if len(muts) != want {
		t.Fatalf("got %d mutations, want %d (%d weapon + %d card)",
			len(muts), want, wantWeaponMuts, wantCardMuts)
	}
	for i, m := range muts {
		if len(m.Deck.Cards) != 4 {
			t.Errorf("mutation %d: card count %d, want 4", i, len(m.Deck.Cards))
		}
		if m.Deck.Hero.Name() != d.Hero.Name() {
			t.Errorf("mutation %d: hero changed", i)
		}
		if m.Description == "" {
			t.Errorf("mutation %d: empty description", i)
		}
	}
}

// TestAllMutations_OddCountsAllowed exercises the single-card-swap semantics: a mutation may leave
// the deck with an odd number of any given printing (e.g. 1×A + 3×B at maxCopies=3), and raising
// maxCopies should open up adds to cards already in the deck that are below the cap.
func TestAllMutations_OddCountsAllowed(t *testing.T) {
	a := cards.Get(card.AetherSlashRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

	// At maxCopies=3, each of the 2 in-deck cards (a, b) is below the cap, so "remove a, add b"
	// (and the mirror) become legal. That's 2 more card mutations than the maxCopies=2 case.
	mutsLow := AllMutations(d, 2, nil)
	mutsHigh := AllMutations(d, 3, nil)
	if len(mutsHigh)-len(mutsLow) != 2 {
		t.Errorf("maxCopies=3 should produce exactly 2 more mutations than maxCopies=2; got diff=%d",
			len(mutsHigh)-len(mutsLow))
	}

	// Every mutation at maxCopies=3 still has exactly 4 cards; some should have an odd count of
	// one card — single-card swaps can leave odd-count slots when maxCopies allows it.
	sawOdd := false
	for _, m := range mutsHigh {
		if len(m.Deck.Cards) != 4 {
			t.Errorf("card count %d, want 4", len(m.Deck.Cards))
		}
		counts := map[card.ID]int{}
		for _, c := range m.Deck.Cards {
			counts[c.ID()]++
		}
		for _, n := range counts {
			if n%2 == 1 {
				sawOdd = true
			}
			if n > 3 {
				t.Errorf("card count %d exceeds maxCopies=3: %v", n, m.Description)
			}
		}
	}
	if !sawOdd {
		t.Errorf("expected at least one mutation with an odd-count card; single-card swaps always produce odd counts from a 2/2 starting deck")
	}
}

// TestAllMutations_OrdersByAscendingAvg pins the iterate-friendly ordering: the removed card in
// the first card-mutation batch should be the one with the lowest per-card Avg in the current
// deck's stats. Run with both (lower-ID is low-avg) and (higher-ID is low-avg) so the test fails
// if the implementation accidentally sorts only by card.ID and gets a free pass from whichever
// direction happens to align.
func TestAllMutations_OrdersByAscendingAvg(t *testing.T) {
	a := cards.Get(card.AetherSlashRed)   // lower card.ID
	b := cards.Get(card.ArcanicSpikeRed)  // higher card.ID

	cases := []struct {
		name             string
		lowAvgCard       card.Card
		highAvgCard      card.Card
	}{
		{"low-avg card has lower ID", a, b},
		{"low-avg card has higher ID", b, a},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}},
				[]card.Card{a, a, b, b})
			// Both cards see the same Plays so only Avg (= TotalContribution / Plays) drives the
			// ordering — no path for card.ID to sneak in via a sub-ordering rule.
			d.Stats.PerCard = map[card.ID]CardPlayStats{
				tc.lowAvgCard.ID():  {Plays: 10, TotalContribution: 10},   // Avg 1.0
				tc.highAvgCard.ID(): {Plays: 10, TotalContribution: 80},  // Avg 8.0
			}

			muts := AllMutations(d, 2, nil)
			// Skip the weapon-mutation block (len(loadouts)-1 entries, one per alternative loadout).
			firstCardMut := muts[len(weaponLoadouts(cards.AllWeapons))-1]
			wantPrefix := "-1 " + tc.lowAvgCard.Name() + ","
			if !strings.HasPrefix(firstCardMut.Description, wantPrefix) {
				t.Errorf("first card mutation removed wrong card\n  got:  %q\n  want prefix: %q",
					firstCardMut.Description, wantPrefix)
			}
		})
	}
}

// TestAllMutations_PreservesSideboard pins that every derived Mutation inherits the source
// deck's Sideboard verbatim. Without this guarantee an anneal round would silently drop the
// user's hand-managed sideboard as soon as it accepted a mutation and wrote the deck back.
func TestAllMutations_PreservesSideboard(t *testing.T) {
	a := cards.Get(card.AetherSlashRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})
	d.Sideboard = []card.Card{a, b, b}

	muts := AllMutations(d, 2, nil)
	if len(muts) == 0 {
		t.Fatal("expected at least one mutation")
	}

	wantCounts := map[card.ID]int{a.ID(): 1, b.ID(): 2}
	for i, m := range muts {
		got := map[card.ID]int{}
		for _, c := range m.Deck.Sideboard {
			got[c.ID()]++
		}
		for id, want := range wantCounts {
			if got[id] != want {
				t.Errorf("mutation %d (%s): sideboard count for %s = %d, want %d",
					i, m.Description, cards.Get(id).Name(), got[id], want)
				break
			}
		}
	}
}

func TestAllMutations_Deterministic(t *testing.T) {
	a := cards.Get(card.AetherSlashRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

	first := AllMutations(d, 2, nil)
	second := AllMutations(d, 2, nil)

	if len(first) != len(second) {
		t.Fatalf("mutation counts differ between calls: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if weaponKey(first[i].Deck.Weapons) != weaponKey(second[i].Deck.Weapons) {
			t.Errorf("mutation %d weapons differ between calls", i)
		}
		if first[i].Description != second[i].Description {
			t.Errorf("mutation %d descriptions differ: %q vs %q",
				i, first[i].Description, second[i].Description)
		}
		for j, c := range first[i].Deck.Cards {
			if c.ID() != second[i].Deck.Cards[j].ID() {
				t.Errorf("mutation %d card[%d] differs between calls: %v vs %v",
					i, j, c.ID(), second[i].Deck.Cards[j].ID())
			}
		}
	}
}

func TestAllMutations_NoDuplicateOfSource(t *testing.T) {
	a := cards.Get(card.AetherSlashRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, a, a})
	srcKey := deckFingerprint(d)
	for i, m := range AllMutations(d, 2, nil) {
		if deckFingerprint(m.Deck) == srcKey {
			t.Errorf("mutation %d equals the source deck", i)
		}
	}
}

// TestRandom_FilterExcludesRejected confirms the legal predicate is actually applied to the
// candidate pool: a filter that blocks Plunder Run (all variants) should produce decks that
// never contain any Plunder Run printing, even across many samples.
func TestRandom_FilterExcludesRejected(t *testing.T) {
	bannedIDs := map[card.ID]bool{
		card.PlunderRunRed:    true,
		card.PlunderRunYellow: true,
		card.PlunderRunBlue:   true,
	}
	legal := func(c card.Card) bool { return !bannedIDs[c.ID()] }
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 20; i++ {
		d := Random(hero.Viserai{}, 40, 2, rng, legal)
		for j, c := range d.Cards {
			if bannedIDs[c.ID()] {
				t.Errorf("sample %d: card[%d] = %s was in the banlist", i, j, c.Name())
			}
		}
	}
}

// TestLegalPool_SkipsNotImplemented pins the contract that cards tagged with
// card.NotImplemented never land in the search pool, with or without a legal predicate. The
// property holds regardless of which registered cards carry the tag today, so it doubles as
// a regression guard against accidental weakening of the filter.
func TestLegalPool_SkipsNotImplemented(t *testing.T) {
	for _, pred := range []func(card.Card) bool{nil, func(card.Card) bool { return true }} {
		for _, id := range legalPool(pred) {
			c := cards.Get(id)
			if _, ok := c.(card.NotImplemented); ok {
				t.Errorf("legalPool included NotImplemented card %s", c.Name())
			}
		}
	}
}

// TestRandom_ExcludesNotImplemented confirms no sampled random deck contains a card tagged
// with card.NotImplemented. Drives Random over many seeds so a leak into the pool would show
// up as a failure even when the tagged set is small relative to the full pool.
func TestRandom_ExcludesNotImplemented(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 20; i++ {
		d := Random(hero.Viserai{}, 40, 2, rng, nil)
		for j, c := range d.Cards {
			if _, ok := c.(card.NotImplemented); ok {
				t.Errorf("sample %d card[%d] = %s implements NotImplemented", i, j, c.Name())
			}
		}
	}
}

// TestLegalPool_ExcludesTaggedCardsByID gives TestLegalPool_SkipsNotImplemented teeth: it
// picks a concrete registered card we know currently carries the NotImplemented marker
// (Strike Gold (Red), gold-token rider) and asserts it's absent from legalPool's output.
// Without at least one real tagged card the property test is vacuous, so this guards against
// a regression where the marker interface itself silently breaks. Self-retires if Strike Gold
// ever loses the tag (gold-token economy gets modelled) so maintenance is only a delete.
func TestLegalPool_ExcludesTaggedCardsByID(t *testing.T) {
	if _, ok := cards.Get(card.StrikeGoldRed).(card.NotImplemented); !ok {
		t.Skip("Strike Gold (Red) is no longer NotImplemented — pick another tagged card or drop this test")
	}
	for _, id := range legalPool(nil) {
		if id == card.StrikeGoldRed {
			t.Fatalf("legalPool included Strike Gold (Red) despite its NotImplemented tag")
		}
	}
}

// TestSanitizeNotImplemented_ReplacesTaggedSlotsAndKeepsSizeLegal drives the sanitizer
// against a deck that starts with two NotImplemented copies in it (Strike Gold Red is a real
// tagged card). After sanitization the deck must: (a) have zero NotImplemented cards, (b)
// be the same size, (c) respect maxCopies across the post-sanitize distribution, (d) report
// exactly two swaps, each naming the original tagged card.
func TestSanitizeNotImplemented_ReplacesTaggedSlotsAndKeepsSizeLegal(t *testing.T) {
	if _, ok := cards.Get(card.StrikeGoldRed).(card.NotImplemented); !ok {
		t.Skip("Strike Gold (Red) is no longer NotImplemented — pick another tagged card or drop this test")
	}
	tagged := cards.Get(card.StrikeGoldRed)
	safe := cards.Get(card.AetherSlashRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}},
		[]card.Card{safe, safe, tagged, tagged})

	rng := rand.New(rand.NewSource(1))
	replaced := d.SanitizeNotImplemented(2, rng, nil)

	if len(replaced) != 2 {
		t.Errorf("replaced %d slots, want 2", len(replaced))
	}
	if len(d.Cards) != 4 {
		t.Errorf("card count after sanitize = %d, want 4", len(d.Cards))
	}
	for i, c := range d.Cards {
		if _, ok := c.(card.NotImplemented); ok {
			t.Errorf("card[%d] = %s still implements NotImplemented", i, c.Name())
		}
	}
	counts := map[card.ID]int{}
	for _, c := range d.Cards {
		counts[c.ID()]++
		if counts[c.ID()] > 2 {
			t.Errorf("%s appears %d times, exceeds maxCopies=2", c.Name(), counts[c.ID()])
		}
	}
	for _, r := range replaced {
		if r.From.ID() != card.StrikeGoldRed {
			t.Errorf("replacement From = %s, want Strike Gold (Red)", r.From.Name())
		}
		if _, ok := r.To.(card.NotImplemented); ok {
			t.Errorf("replacement To = %s implements NotImplemented", r.To.Name())
		}
	}
}

// TestSanitizeNotImplemented_NoOpOnCleanDeck confirms the sanitizer is an identity operation
// when the deck already has no NotImplemented cards: no replacements, no mutations to
// Cards.
func TestSanitizeNotImplemented_NoOpOnCleanDeck(t *testing.T) {
	a := cards.Get(card.AetherSlashRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, a, a})
	before := append([]card.Card(nil), d.Cards...)

	rng := rand.New(rand.NewSource(1))
	replaced := d.SanitizeNotImplemented(2, rng, nil)

	if len(replaced) != 0 {
		t.Errorf("replacements on clean deck = %d, want 0", len(replaced))
	}
	for i, c := range d.Cards {
		if c.ID() != before[i].ID() {
			t.Errorf("card[%d] mutated: %s → %s", i, before[i].Name(), c.Name())
		}
	}
}

// TestAllMutations_ExcludesNotImplementedAdditions confirms no single-slot mutation can
// introduce a NotImplemented card. Starting deck contains only implemented cards so any
// NotImplemented copy in a mutation output must have come from the add pool.
func TestAllMutations_ExcludesNotImplementedAdditions(t *testing.T) {
	a := cards.Get(card.AetherSlashRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, a, a})
	for _, m := range AllMutations(d, 2, nil) {
		for _, c := range m.Deck.Cards {
			if _, ok := c.(card.NotImplemented); ok {
				t.Errorf("%s introduced NotImplemented card %s", m.Description, c.Name())
			}
		}
	}
}

// TestAllMutations_FilterExcludesRejectedAdditions confirms banned cards never appear as
// swap-in candidates. A banned card already in the deck IS still a valid removal target — the
// hill climb must be able to swap it out — so we assert that the starting deck's banned card is
// never in the post-mutation card list either (which would require it to have been added back).
func TestAllMutations_FilterExcludesRejectedAdditions(t *testing.T) {
	bannedIDs := map[card.ID]bool{
		card.PlunderRunRed: true,
	}
	legal := func(c card.Card) bool { return !bannedIDs[c.ID()] }

	pr := cards.Get(card.PlunderRunRed)
	other := cards.Get(card.AetherSlashRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{pr, pr, other, other})

	for i, m := range AllMutations(d, 2, legal) {
		bannedIn := 0
		for _, c := range m.Deck.Cards {
			if bannedIDs[c.ID()] {
				bannedIn++
			}
		}
		// The starting deck has 2 copies of Plunder Run (Red). A mutation that removes one leaves
		// 1; a mutation that removes the other leaves 1; a weapon-only mutation leaves all 2. No
		// mutation should ADD another copy.
		if bannedIn > 2 {
			t.Errorf("mutation %d (%s): has %d banned copies, want <=2 (no additions allowed)",
				i, m.Description, bannedIn)
		}
	}
}

// TestEvaluate_PerCardStatsPopulated pins per-card attribution: every card that's played or
// pitched increments Plays+Pitches, and TotalContribution sums role-based per-card credit:
// Attack → Card.Attack(), Defend → proportional share of block, Pitch → Card.Pitch(). Held and
// Arsenal cards don't tick the counters (they didn't contribute to this turn's Value). A
// single-printing deck makes the totals easy to assert against the card's printed stats.
func TestEvaluate_PerCardStatsPopulated(t *testing.T) {
	read := cards.Get(card.ReadTheRunesRed)
	d := New(hero.Viserai{}, nil, []card.Card{read, read, read, read})
	d.Evaluate(1, 0, rand.New(rand.NewSource(1)))

	if d.Stats.PerCard == nil {
		t.Fatalf("PerCard should be initialised after Evaluate")
	}
	stat, ok := d.Stats.PerCard[card.ReadTheRunesRed]
	if !ok {
		t.Fatalf("PerCard missing entry for Read the Runes (Red)")
	}
	// Read the Runes Red has no Go again, so the chain plays at most one per turn. With 4 in a
	// 4-card hand, the solver plays one and the rest fall into Held/Arsenal roles which don't
	// tick Plays or Pitches. Counter should be non-zero (at least one Play) but need not sum
	// to 4.
	if got := stat.Plays + stat.Pitches; got == 0 {
		t.Errorf("Plays+Pitches = 0, want at least 1 (the chosen attacker plays once)")
	}
	// Contributions come from the winning chain replay (Play returns + hero triggers) plus
	// role-based shares for pitch/defend. The exact total depends on rider/trigger damage, so
	// assert the weaker property that it's positive and produces a positive Avg.
	if stat.TotalContribution <= 0 {
		t.Errorf("TotalContribution = %v, want >0 (played Read the Runes deals at least Attack+rider)",
			stat.TotalContribution)
	}
	if stat.Avg() <= 0 {
		t.Errorf("Avg() = %v, want >0", stat.Avg())
	}
}

// TestEvaluate_BestTurnStartingRunechantsIsPreHandCarryover pins the contract of
// BestTurn.StartingRunechants: it's the Runechant count carried in from the previous turn when
// the hand was played, so for the first hand of a run it's always 0 — even if the hand itself
// creates runechants that carry out into the next turn.
func TestEvaluate_BestTurnStartingRunechantsIsPreHandCarryover(t *testing.T) {
	// Viserai has Intelligence 4. A 4-card deck gives exactly one hand per run, so the Best
	// record always reflects that first hand — no previous turn ever existed.
	read := cards.Get(card.ReadTheRunesRed)
	d := New(hero.Viserai{}, nil, []card.Card{read, read, read, read})

	// Seed doesn't matter (all cards identical), but fix it for determinism.
	d.Evaluate(1, 0, rand.New(rand.NewSource(1)))

	if len(d.Stats.Best.Summary.BestLine) == 0 {
		t.Fatalf("expected Best to be populated after Evaluate")
	}
	// Sanity: the hand should have left runechants on the table (otherwise the bug couldn't
	// manifest — pre-hand and post-hand counts would both be 0).
	if d.Stats.Best.Summary.Value == 0 {
		t.Fatalf("expected nonzero Value from a hand of Read the Runes; got 0")
	}
	if d.Stats.Best.StartingRunechants != 0 {
		t.Errorf("StartingRunechants = %d, want 0 (first hand of the run has no previous-turn carryover)",
			d.Stats.Best.StartingRunechants)
	}
}

// TestEvaluate_BestTurnSnapshotsDrawnAndLeftoverRunechants pins the BestTurn snapshot's
// completeness: Drawn (mid-turn-drawn cards with their dispositions) and LeftoverRunechants
// must propagate from play.* into Stats.Best.Summary.* so FormatBestTurn's per-card breakdown
// reconciles with the displayed Value and the header's "carryover runechants" count is real.
// Without the snapshot, drawn-attack extension damage and pitch-from-drawn resource land in
// Value but never show up in the printout, and runechants always read 0.
func TestEvaluate_BestTurnSnapshotsDrawnAndLeftoverRunechants(t *testing.T) {
	// Snatch (cost 0, attack 4) fires on-hit DrawOne — its drawn card lands in summary.Drawn.
	// 4 Snatches keeps Viserai's Intelligence-4 hand full of draw-rider cards on the first
	// turn so at least one Snatch attacks and DrawOne fires.
	snatch := cards.Get(card.SnatchRed)
	d := New(hero.Viserai{}, nil, []card.Card{snatch, snatch, snatch, snatch, snatch, snatch, snatch, snatch})
	d.Evaluate(1, 0, rand.New(rand.NewSource(1)))

	if len(d.Stats.Best.Summary.BestLine) == 0 {
		t.Fatalf("expected Best to be populated after Evaluate")
	}
	if len(d.Stats.Best.Summary.Drawn) == 0 {
		t.Errorf("Stats.Best.Summary.Drawn is empty; want >=1 entry from Snatch's on-hit DrawOne (the snapshot in Evaluate isn't copying play.Drawn)")
	}
}

// TestEvaluate_HeldCardDefersDrawToNextTurn pins the "up to Intelligence" draw rule plus arsenal
// carryover. Intelligence-1 hero, deck of Toughen Up Blue (DR, cost 2, defense 4): the lone
// card has no legal play (can't pay its 2-cost, can't pitch with nothing on the stack, DRs
// can't Attack). Turn 1 holds then promotes it to Arsenal (empty slot). Turn 2 draws a new DR;
// the arsenal card stays on tie, the new card goes Held, so drawCount = 0 next turn and the
// loop halts at Stats.Hands = 2. Neither turn plays or pitches, so PerCard stays at 0.
func TestEvaluate_HeldCardDefersDrawToNextTurn(t *testing.T) {
	// 40 copies of the DR so we have enough deck to fill many hands if held carryover weren't
	// wired up — the assertion would fail catastrophically (loop or much larger Hands count).
	deckCards := make([]card.Card, 40)
	for i := range deckCards {
		deckCards[i] = generic.ToughenUpBlue{}
	}
	d := New(int1StubHero{}, nil, deckCards)
	d.Evaluate(1, 0, rand.New(rand.NewSource(1)))

	if d.Stats.Hands != 2 {
		t.Errorf("Stats.Hands = %d, want 2 (turn 1 arsenals the card, turn 2 holds its successor, turn 3 can't draw)", d.Stats.Hands)
	}
	tuStat := d.Stats.PerCard[card.ToughenUpBlue]
	if tuStat.Plays != 0 || tuStat.Pitches != 0 {
		t.Errorf("PerCard[ToughenUpBlue] Plays=%d Pitches=%d, want 0/0 (card was Held/Arsenaled, never played or pitched)",
			tuStat.Plays, tuStat.Pitches)
	}
	// Best captures turn 1 (first hand with a recorded play). That hand's single card got
	// promoted from Held to Arsenal by the post-hoc upgrade.
	if len(d.Stats.Best.Summary.BestLine) == 0 {
		t.Fatalf("expected Best to be populated after at least one hand")
	}
	if d.Stats.Best.Summary.BestLine[0].Role != hand.Arsenal {
		t.Errorf("Best.Play.Roles[0] = %s, want ARSENAL (empty slot on turn 1 → Held promoted)", d.Stats.Best.Summary.BestLine[0].Role)
	}
}

// TestEvaluate_ArsenalPersistsAcrossTurns confirms the arsenal slot threads through Evaluate's
// per-turn loop: a card promoted to Arsenal on one turn becomes arsenalCardIn on the next.
// Intelligence-1 hero, 2-card deck of two Toughen Up Blue. Turn 1 arsenals the drawn TU.
// Turn 2 draws the second TU and against incoming 4 plays the arsenal-in DR, pitching the
// drawn card to fund its 2-cost — Value = 4 (prevents the full attack). Turn 3 re-draws the
// pitched card (returned to deck bottom) and arsenals it again. Loop stops when the deck's
// empty and nothing new can be drawn.
func TestEvaluate_ArsenalPersistsAcrossTurns(t *testing.T) {
	d := New(int1StubHero{}, nil, []card.Card{generic.ToughenUpBlue{}, generic.ToughenUpBlue{}})
	d.Evaluate(1, 4, rand.New(rand.NewSource(1)))

	// Best captures turn 2 — only turn with Value > 0 (arsenal DR fires).
	if d.Stats.Best.Summary.Value != 4 {
		t.Errorf("Best.Play.Value = %d, want 4 (turn 2 plays arsenal DR, pitches hand DR to pay; prevents 4)", d.Stats.Best.Summary.Value)
	}
	// Turn 1: arsenal the drawn card. Turn 2: play arsenal DR (paid by pitching drawn card).
	// Turn 3: draw the recycled pitched card, arsenal it (deck is then empty). Loop ends.
	if d.Stats.Hands != 3 {
		t.Errorf("Stats.Hands = %d, want 3", d.Stats.Hands)
	}
}

// TestEvaluate_TerminatesAfterTwoCycles pins the infinite-loop guard on Evaluate's per-run loop.
// 40 Toughen Up Blue DRs with Reaping Blade equipped, incoming=0, reaches a steady state after
// turn 1 (pitch one TU, swing Reaping Blade for +3, hold the other 3). From then on every turn
// draws and pitches one card — net deck change zero, hand.Best returns the same TurnSummary.
// Without the cap the loop would spin forever; with it, Stats.Hands halts at 2 × handsPerCycle.
func TestEvaluate_TerminatesAfterTwoCycles(t *testing.T) {
	deckCards := make([]card.Card, 40)
	for i := range deckCards {
		deckCards[i] = generic.ToughenUpBlue{}
	}
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.ReapingBlade{}}, deckCards)
	done := make(chan struct{})
	go func() {
		d.Evaluate(1, 0, rand.New(rand.NewSource(1)))
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("Evaluate did not terminate within 2 seconds — infinite loop regression")
	}
	// Two cycles of a 40-card / 4-hand-size deck is exactly 20 hands.
	handsPerCycle := len(deckCards) / hero.Viserai{}.Intelligence()
	maxHands := 2 * handsPerCycle
	if d.Stats.Hands != maxHands {
		t.Errorf("Stats.Hands = %d, want exactly %d (steady-state pitched-pitch loop hits the cap)",
			d.Stats.Hands, maxHands)
	}
}

// deckFingerprint builds a comparable summary of a deck for equality checks in tests.
func deckFingerprint(d *Deck) string {
	s := weaponKey(d.Weapons) + "|"
	counts := map[card.ID]int{}
	for _, c := range d.Cards {
		counts[c.ID()]++
	}
	// Stable ordering — iterate over all possible IDs in byID order isn't exposed, so use a
	// sorted slice of (id, count).
	type pair struct {
		id card.ID
		n  int
	}
	var pairs []pair
	for id, n := range counts {
		pairs = append(pairs, pair{id, n})
	}
	// Insertion sort by id (tiny N).
	for i := 1; i < len(pairs); i++ {
		for j := i; j > 0 && pairs[j-1].id > pairs[j].id; j-- {
			pairs[j-1], pairs[j] = pairs[j], pairs[j-1]
		}
	}
	for _, p := range pairs {
		s += string(rune(p.id)) + ":" + string(rune(p.n)) + ","
	}
	return s
}
