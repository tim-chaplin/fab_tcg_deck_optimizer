package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/generic"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// TestBest_EmptyArsenalClaimsHeldCard confirms the post-hoc Arsenal promotion fires when the
// slot is empty and the winning partition has Held cards. A hand that can't play Toughen Up as
// DR (no other card to pitch for the 2-cost) leaves the DR Held; with arsenalCardIn=nil the
// slot is empty so the DR becomes Arsenal and rides into next turn as Play.ArsenalCard.
func TestBest_EmptyArsenalClaimsHeldCard(t *testing.T) {
	h := []card.Card{generic.ToughenUpBlue{}}
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	if got.BestLine[0].Role != Arsenal {
		t.Errorf("Roles[0] = %s, want ARSENAL", got.BestLine[0].Role)
	}
	if got.ArsenalCard == nil || got.ArsenalCard.ID() != card.ToughenUpBlue {
		t.Errorf("ArsenalCard = %v, want Toughen Up Blue", got.ArsenalCard)
	}
}

// TestBest_ArsenalInPlayDR covers the "arsenal card played as DR" branch. Previous turn left a
// Toughen Up Blue in arsenal; this turn we draw a Blue Malefic (pitch 3, cost 0). The pitched
// Malefic funds Toughen Up's 2-cost defense out of the arsenal, preventing 4 damage. Value = 4.
// Play.ArsenalCard is nil because the slot was vacated and no hand card ends up Held.
func TestBest_ArsenalInPlayDR(t *testing.T) {
	h := []card.Card{runeblade.MaleficIncantationBlue{}}
	got := Best(stubHero, nil, h, 4, nil, 0, generic.ToughenUpBlue{})
	if got.Value != 4 {
		t.Fatalf("Value = %d, want 4 (Malefic pitches to pay arsenal DR, prevents 4). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
	if got.ArsenalCard != nil {
		t.Errorf("ArsenalCard = %v, want nil (slot was vacated, no Held card to promote)", got.ArsenalCard)
	}
	// ArsenalIn surfaces the arsenal-in assignment so callers (the best-hand printout) can flag
	// that this card wasn't in hand this turn.
	ai, hasArsenal := got.ArsenalIn()
	if !hasArsenal || ai.Card.ID() != card.ToughenUpBlue {
		t.Errorf("ArsenalIn = %v, want Toughen Up Blue", ai)
	}
	if ai.Role != Defend {
		t.Errorf("ArsenalIn role = %s, want DEFEND", ai.Role)
	}
}

// TestBest_ArsenalInStayBlocksNewArsenal locks in that while the arsenal slot is occupied, a
// hand card that would otherwise be promoted to Arsenal (because it's Held) stays Held instead —
// one arsenal slot, no replacement until the old card is played. A lone DR in hand is Held;
// the arsenal-in Toughen Up Blue stays (incoming=0 makes defending pointless, and the hand
// can't fund a DR anyway); post-hoc the slot is occupied so no promotion happens.
func TestBest_ArsenalInStayBlocksNewArsenal(t *testing.T) {
	h := []card.Card{generic.ToughenUpBlue{}}
	got := Best(stubHero, nil, h, 0, nil, 0, generic.ToughenUpBlue{})
	if got.BestLine[0].Role != Held {
		t.Errorf("Roles[0] = %s, want HELD (slot occupied by arsenal-in, can't promote)", got.BestLine[0].Role)
	}
	if got.ArsenalCard == nil || got.ArsenalCard.ID() != card.ToughenUpBlue {
		t.Errorf("ArsenalCard = %v, want Toughen Up Blue (the staying arsenal-in card)", got.ArsenalCard)
	}
}

// TestBest_ArsenalInPlayAttack covers the "arsenal card played as attack" branch. A Red attack
// sits in arsenal from a previous turn; this turn we draw a single Red Attack which pitches
// (pitch 1) to fund both the hand Red's 1-cost and the arsenal Red's 1-cost... wait, one pitch
// can't pay two costs. Instead, the winning line plays the arsenal Red (funded by pitching the
// hand Red) and leaves the hand slot consumed. Value = 3 (arsenal Red's attack). With the
// arsenal slot now empty and no Held cards, ArsenalCard is nil.
func TestBest_ArsenalInPlayAttack(t *testing.T) {
	h := []card.Card{fake.RedAttack{}}
	got := Best(stubHero, nil, h, 0, nil, 0, fake.RedAttack{})
	if got.Value != 3 {
		t.Fatalf("Value = %d, want 3 (arsenal Red played, hand Red pitched to fund it). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
	if got.ArsenalCard != nil {
		t.Errorf("ArsenalCard = %v, want nil (slot vacated, no Held to promote)", got.ArsenalCard)
	}
}

// TestBest_ArsenalInNonAttackActionPlays covers the "arsenal card isn't tagged Attack but can
// still be played on your turn" rule — non-attack actions (auras, item cards, etc.) are playable
// from arsenal. Hand: Malefic Incantation Blue (cost 0, pitch 3, Play returns 1 flat with no
// follow-up attack). Arsenal: Arcane Cussing Red (cost 1, pitch 1, Play returns 3 when we
// block all incoming). The winning line pitches Malefic to fund Cussing's 1-cost and plays
// Cussing from arsenal for a flat 3.
func TestBest_ArsenalInNonAttackActionPlays(t *testing.T) {
	h := []card.Card{runeblade.MaleficIncantationBlue{}}
	got := Best(stubHero, nil, h, 0, nil, 0, runeblade.ArcaneCussingRed{})
	if got.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Malefic pitched, arsenal Cussing played for 3). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
	if got.ArsenalCard != nil {
		t.Errorf("ArsenalCard = %v, want nil (Cussing played out of arsenal)", got.ArsenalCard)
	}
}

// TestBest_ArsenalInUnmovableGrantsDefenseBonus pins the DR-from-arsenal +N{d} rider:
// Unmovable Red printed Defense() is 7 and grants +1{d} when played from arsenal. Hand: Blue
// Malefic (pitch 3, cost 0). Arsenal: Unmovable Red. Pitched Malefic funds Unmovable's 3-cost
// defense; effective defense is 7 + 1 (from-arsenal) = 8, fully blocking 8 incoming. Value = 8.
// If the rider didn't fire, prevented would cap at 7.
func TestBest_ArsenalInUnmovableGrantsDefenseBonus(t *testing.T) {
	h := []card.Card{runeblade.MaleficIncantationBlue{}}
	got := Best(stubHero, nil, h, 8, nil, 0, generic.UnmovableRed{})
	if got.Value != 8 {
		t.Fatalf("Value = %d, want 8 (Unmovable from arsenal blocks 7+1). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBest_HandUnmovableNoDefenseBonus confirms the +1{d} rider does NOT fire when Unmovable
// is played from hand. Hand: Blue Malefic + Unmovable Red, no arsenal. Pitched Malefic funds
// Unmovable's 3-cost; effective defense stays at printed 7, so 8 incoming caps prevented at 7.
// If the rider mistakenly fired from hand, prevented would be 8.
func TestBest_HandUnmovableNoDefenseBonus(t *testing.T) {
	h := []card.Card{runeblade.MaleficIncantationBlue{}, generic.UnmovableRed{}}
	got := Best(stubHero, nil, h, 8, nil, 0, nil)
	if got.Value != 7 {
		t.Fatalf("Value = %d, want 7 (hand-played Unmovable: no rider). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestBest_ArsenalInSmashingGoodTimeGatesOnlyArsenalCopy pins the from-arsenal gate: only the
// SGT that came from the arsenal grants its +3 rider; the hand copy returns 0. Hero = Viserai.
// Arsenal: SGT Red. Hand: SGT Red + Hocus Pocus Red. Best line plays both SGTs (non-attack
// actions, go again) ahead of Hocus Pocus. Arsenal SGT's Play scans CardsRemaining, finds Hocus
// (attack action) and credits +3; hand SGT's Play fails the FromArsenal check and returns 0.
// Hocus Pocus contributes 3 base + 1 from its own Runechant + 1 from Viserai's hero ability
// (fires because two non-attack actions were already played). Value = 3 + 0 + 3 + 1 + 1 = 8.
// If the from-arsenal gate weren't enforced, both SGTs would grant their rider and value would
// be 11.
func TestBest_ArsenalInSmashingGoodTimeGatesOnlyArsenalCopy(t *testing.T) {
	h := []card.Card{
		generic.SmashingGoodTimeRed{},
		runeblade.HocusPocusRed{},
	}
	got := Best(hero.Viserai{}, nil, h, 0, nil, 0, generic.SmashingGoodTimeRed{})
	if got.Value != 8 {
		t.Fatalf("Value = %d, want 8 (only arsenal SGT grants +3). Roles=[%s]",
			got.Value, FormatBestLine(got.BestLine))
	}
}

// TestPromoteRandomHeldToArsenal_SpreadsAcrossHands pins the post-hoc Held→Arsenal promotion's
// anti-bias property: the selection hashes the sorted hand IDs so different hands land on
// different Held positions rather than always picking slot 0 (which, with IDs sorted, would
// systematically prefer low-ID cards). Drives the helper directly with synthesised BestLine
// entries — all slots Held, all equivalent in value — so only the hash-based index selection
// is under test.
func TestPromoteRandomHeldToArsenal_SpreadsAcrossHands(t *testing.T) {
	// 20 different 4-card hands using Wounding Blow Red/Yellow/Blue as "arbitrary cards with
	// distinct IDs". Varying which card sits in which slot is enough to exercise the hash across
	// different inputs.
	wbR := generic.WoundingBlowRed{}
	wbY := generic.WoundingBlowYellow{}
	wbB := generic.WoundingBlowBlue{}
	hands := [][]card.Card{
		{wbR, wbR, wbR, wbY}, {wbR, wbR, wbY, wbY}, {wbR, wbR, wbY, wbB}, {wbR, wbY, wbY, wbB},
		{wbR, wbY, wbB, wbB}, {wbR, wbR, wbR, wbB}, {wbR, wbR, wbB, wbB}, {wbY, wbY, wbY, wbB},
		{wbY, wbY, wbB, wbB}, {wbY, wbB, wbB, wbB}, {wbR, wbR, wbR, wbR}, {wbY, wbY, wbY, wbY},
		{wbB, wbB, wbB, wbB}, {wbR, wbY, wbY, wbY}, {wbR, wbR, wbY, wbR}, {wbB, wbR, wbY, wbR},
		{wbB, wbY, wbR, wbR}, {wbR, wbB, wbB, wbY}, {wbR, wbR, wbB, wbY}, {wbY, wbR, wbB, wbB},
	}
	slots := map[int]int{}
	for _, h := range hands {
		// Sort to match bestUncached's canonical order. Build a BestLine where every slot is Held.
		handCopy := append([]card.Card(nil), h...)
		ids := make([]card.ID, len(handCopy))
		for i, c := range handCopy {
			ids[i] = c.ID()
		}
		sortHandByID(handCopy, ids, len(handCopy))
		line := make([]CardAssignment, len(handCopy))
		for i, c := range handCopy {
			line[i] = CardAssignment{Card: c, Role: Held}
		}
		best := TurnSummary{BestLine: line}
		promoteRandomHeldToArsenal(&best, handCopy, len(handCopy), nil)
		idx := -1
		for i, a := range best.BestLine {
			if a.Role == Arsenal {
				idx = i
				break
			}
		}
		if idx < 0 {
			t.Fatalf("hand %v: no Arsenal-role slot found in BestLine=%s", h, FormatBestLine(best.BestLine))
		}
		slots[idx]++
	}
	if len(slots) < 2 {
		t.Errorf("arsenal promotion only ever landed on slot %v across %d hands; expected spread across multiple slots", slots, len(hands))
	}
}

// TestPromoteRandomHeldToArsenal_DeterministicPerHand pins the other half of the contract: a
// given hand produces the SAME picked slot every call, so the memo cache doesn't drift between
// hits and repeated simulations of the same deck stay reproducible.
func TestPromoteRandomHeldToArsenal_DeterministicPerHand(t *testing.T) {
	hand := []card.Card{
		generic.WoundingBlowRed{}, generic.WoundingBlowYellow{},
		generic.WoundingBlowBlue{}, generic.WoundingBlowBlue{},
	}
	var firstIdx int
	for run := 0; run < 5; run++ {
		line := []CardAssignment{
			{Card: hand[0], Role: Held},
			{Card: hand[1], Role: Held},
			{Card: hand[2], Role: Held},
			{Card: hand[3], Role: Held},
		}
		best := TurnSummary{BestLine: line}
		promoteRandomHeldToArsenal(&best, hand, len(hand), nil)
		idx := -1
		for i, a := range best.BestLine {
			if a.Role == Arsenal {
				idx = i
				break
			}
		}
		if run == 0 {
			firstIdx = idx
			continue
		}
		if idx != firstIdx {
			t.Errorf("run %d: Arsenal at slot %d, want %d (deterministic per-hand)", run, idx, firstIdx)
		}
	}
}

// TestPromoteRandomHeldToArsenal_SingleHeldAlwaysPicked covers the n=1 edge of the hash-modulo
// selection: with exactly one Held index the modulo is deterministic (always 0), so the only
// candidate gets promoted regardless of the hash value.
func TestPromoteRandomHeldToArsenal_SingleHeldAlwaysPicked(t *testing.T) {
	// Hand of two cards where one is Attack-role and one is Held-role — exactly one candidate
	// for post-hoc promotion.
	hand := []card.Card{generic.WoundingBlowRed{}, generic.WoundingBlowBlue{}}
	line := []CardAssignment{
		{Card: hand[0], Role: Attack},
		{Card: hand[1], Role: Held},
	}
	best := TurnSummary{BestLine: line}
	promoteRandomHeldToArsenal(&best, hand, len(hand), nil)
	if best.BestLine[1].Role != Arsenal {
		t.Errorf("Role[1] = %s, want Arsenal (only Held candidate)", best.BestLine[1].Role)
	}
	if best.ArsenalCard == nil || best.ArsenalCard.ID() != hand[1].ID() {
		t.Errorf("ArsenalCard = %v, want %s", best.ArsenalCard, hand[1].Name())
	}
}

// TestPromoteRandomHeldToArsenal_NoHeldIsNoop covers the other end: a partition where every
// hand card plays/pitches/defends leaves zero Held candidates, so the promotion is a no-op and
// the arsenal slot stays empty.
func TestPromoteRandomHeldToArsenal_NoHeldIsNoop(t *testing.T) {
	hand := []card.Card{generic.WoundingBlowRed{}, generic.WoundingBlowBlue{}}
	line := []CardAssignment{
		{Card: hand[0], Role: Attack},
		{Card: hand[1], Role: Pitch},
	}
	best := TurnSummary{BestLine: line}
	promoteRandomHeldToArsenal(&best, hand, len(hand), nil)
	for i, a := range best.BestLine {
		if a.Role == Arsenal {
			t.Errorf("BestLine[%d].Role = Arsenal, want unchanged (no Held candidates)", i)
		}
	}
	if best.ArsenalCard != nil {
		t.Errorf("ArsenalCard = %v, want nil (no promotion possible)", best.ArsenalCard)
	}
}

// TestBeatsBest_ArsenalOccupancyTiebreaker pins the tiebreaker contract used by the partition
// enumerator: when two candidates tie on Value and LeftoverRunechants, the one that will end
// the turn with the arsenal slot occupied (either via arsenal-in staying OR a post-hoc Held →
// Arsenal promotion) beats the one that won't. Exercised directly so a comparison-inversion
// regression can't hide behind enumeration order at the Best() level.
func TestBeatsBest_ArsenalOccupancyTiebreaker(t *testing.T) {
	// Seed best: Value=10, Leftover=0, arsenal NOT occupied, no future-value plays.
	best := TurnSummary{Value: 10, LeftoverRunechants: 0}
	// Candidate with equal V/L/future-value but arsenal WILL be occupied — should beat.
	if !beatsBest(10, 0, 0, true, best, 0, false) {
		t.Error("willOccupy=true should beat a tied best with willOccupy=false")
	}
	// Candidate with equal V/L and arsenal NOT occupied — same as best, should NOT beat.
	if beatsBest(10, 0, 0, false, best, 0, false) {
		t.Error("willOccupy=false should not beat a tied best with willOccupy=false")
	}
	// Best already occupies; candidate also occupies — no advantage, should NOT beat.
	if beatsBest(10, 0, 0, true, best, 0, true) {
		t.Error("willOccupy=true should not beat a tied best that also has willOccupy=true")
	}
	// Strict-wins on Value still takes precedence over the occupancy tiebreaker.
	if !beatsBest(11, 0, 0, false, best, 0, true) {
		t.Error("higher Value should beat even when the candidate has no occupancy advantage")
	}
	// Strict-loses on Value — can't be rescued by occupancy.
	if beatsBest(9, 0, 0, true, best, 0, false) {
		t.Error("lower Value should lose regardless of occupancy advantage")
	}
	// Strict-wins on leftover takes precedence over occupancy.
	if !beatsBest(10, 1, 0, false, best, 0, true) {
		t.Error("higher LeftoverRunechants should beat even without occupancy advantage")
	}
}

// TestBeatsBest_FutureValueTiebreaker pins the future-value bias: at equal Value and
// LeftoverRunechants, a partition that plays more card.AddsFutureValue cards wins over one
// that plays fewer, regardless of arsenal occupancy. This corrects for the hidden later-turn
// value those cards carry — without the bias, a lone sigil ends up Held → promoted to
// arsenal because same-turn Value is 0 and arsenal occupancy wins the fallback tiebreak.
func TestBeatsBest_FutureValueTiebreaker(t *testing.T) {
	best := TurnSummary{Value: 5, LeftoverRunechants: 0}
	// Candidate plays 1 future-value card, best plays 0 — candidate wins even though arsenal
	// occupancy favours the best.
	if !beatsBest(5, 0, 1, false, best, 0, true) {
		t.Error("more future-value cards should beat a tied best with occupancy advantage")
	}
	// Reverse: best plays 1 future-value, candidate plays 0 — candidate loses even when
	// candidate has the occupancy advantage.
	if beatsBest(5, 0, 0, true, best, 1, false) {
		t.Error("fewer future-value cards should lose even with occupancy advantage")
	}
	// Strict-wins on Value still takes precedence over the future-value tiebreaker.
	if !beatsBest(6, 0, 0, false, best, 5, false) {
		t.Error("higher Value should beat even when the candidate plays fewer future-value cards")
	}
	// Strict-wins on LeftoverRunechants still takes precedence over future-value.
	if !beatsBest(5, 1, 0, false, best, 5, false) {
		t.Error("higher LeftoverRunechants should beat even when the candidate plays fewer future-value cards")
	}
}
