package deck

import (
	"math/rand"
	"strings"
	"testing"

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

	muts := AllMutations(d, 2)

	// Weapon mutations: every loadout except the current one. Card mutations at maxCopies=2:
	// for each of the 2 unique removals, every pool entry except self (no-op) and the other
	// in-deck card (already at cap) is a valid add — so 2 × (pool - 2).
	loadouts := weaponLoadouts(cards.AllWeapons)
	pool := cards.Deckable()
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
	mutsLow := AllMutations(d, 2)
	mutsHigh := AllMutations(d, 3)
	if len(mutsHigh)-len(mutsLow) != 2 {
		t.Errorf("maxCopies=3 should produce exactly 2 more mutations than maxCopies=2; got diff=%d",
			len(mutsHigh)-len(mutsLow))
	}

	// Every mutation at maxCopies=3 still has exactly 4 cards; some should have an odd count of
	// one card (that's the whole point — no longer forced to keep pairs).
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

			muts := AllMutations(d, 2)
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

func TestAllMutations_Deterministic(t *testing.T) {
	a := cards.Get(card.AetherSlashRed)
	b := cards.Get(card.ArcanicSpikeRed)
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.NebulaBlade{}}, []card.Card{a, a, b, b})

	first := AllMutations(d, 2)
	second := AllMutations(d, 2)

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
	for i, m := range AllMutations(d, 2) {
		if deckFingerprint(m.Deck) == srcKey {
			t.Errorf("mutation %d equals the source deck", i)
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

// TestEvaluate_BestHandStartingRunechantsIsPreHandCarryover pins down a subtle bug: Evaluate
// used to write the post-hand LeftoverRunechants into BestHand.StartingRunechants, so the field
// surfaced the wrong turn's count. The field is documented as "the Runechant count carried in
// from the previous turn when this hand was played", which for the first hand of a run is
// always 0 — even if the hand itself creates runechants that leftover into the next turn.
//
// Without the snapshot fix this test fails: StartingRunechants equals LeftoverRunechants (nonzero)
// instead of 0.
func TestEvaluate_BestHandStartingRunechantsIsPreHandCarryover(t *testing.T) {
	// Viserai has Intelligence 4. A 4-card deck gives exactly one hand per run, so the Best
	// record always reflects that first hand — no previous turn ever existed.
	read := cards.Get(card.ReadTheRunesRed)
	d := New(hero.Viserai{}, nil, []card.Card{read, read, read, read})

	// Seed doesn't matter (all cards identical), but fix it for determinism.
	d.Evaluate(1, 0, rand.New(rand.NewSource(1)))

	if d.Stats.Best.Hand == nil {
		t.Fatalf("expected Best to be populated after Evaluate")
	}
	// Sanity: the hand should have left runechants on the table (otherwise the bug couldn't
	// manifest — pre-hand and post-hand counts would both be 0).
	if d.Stats.Best.Play.Value == 0 {
		t.Fatalf("expected nonzero Value from a hand of Read the Runes; got 0")
	}
	if d.Stats.Best.StartingRunechants != 0 {
		t.Errorf("StartingRunechants = %d, want 0 (first hand of the run has no previous-turn carryover)",
			d.Stats.Best.StartingRunechants)
	}
}

// TestEvaluate_HeldCardDefersDrawToNextTurn pins down the "up to Intelligence" draw rule and
// the arsenal carryover that sits on top of it. An Intelligence-1 hero with a deck of Toughen
// Up Blue (DR, cost 2, defense 4) has no legal play: the lone card can't pay its own 2-cost to
// fire as a DR, can't be pitched (nothing unpaid on the stack), and can't Attack (DRs can't).
//
// Turn 1: the drawn card is Held by the partition; the arsenal slot is empty so the post-hoc
// upgrade promotes it to Arsenal. Turn 2: a fresh DR is drawn; the arsenal card (from turn 1)
// can stay or play as DR, and staying ties on Value, so it sits. The newly-drawn card goes
// Held and next turn's drawCount is handSize - held = 0, halting the loop. Stats.Hands = 2.
// Neither turn plays or pitches the card, so PerCard counters stay at 0.
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
	if d.Stats.Best.Hand == nil || len(d.Stats.Best.Play.Roles) == 0 {
		t.Fatalf("expected Best to be populated after at least one hand")
	}
	if d.Stats.Best.Play.Roles[0] != hand.Arsenal {
		t.Errorf("Best.Play.Roles[0] = %s, want ARSENAL (empty slot on turn 1 → Held promoted)", d.Stats.Best.Play.Roles[0])
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
	if d.Stats.Best.Play.Value != 4 {
		t.Errorf("Best.Play.Value = %d, want 4 (turn 2 plays arsenal DR, pitches hand DR to pay; prevents 4)", d.Stats.Best.Play.Value)
	}
	// Turn 1: arsenal the drawn card. Turn 2: play arsenal DR (paid by pitching drawn card).
	// Turn 3: draw the recycled pitched card, arsenal it (deck is then empty). Loop ends.
	if d.Stats.Hands != 3 {
		t.Errorf("Stats.Hands = %d, want 3", d.Stats.Hands)
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
