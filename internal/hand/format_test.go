package hand

import (
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/generic"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// TestRole_String pins the human-readable labels for each Role value so display strings stay
// stable.
func TestRole_String(t *testing.T) {
	cases := []struct {
		r    Role
		want string
	}{
		{Pitch, "PITCH"},
		{Attack, "ATTACK"},
		{Defend, "DEFEND"},
		{Held, "HELD"},
		{Arsenal, "ARSENAL"},
		{Role(99), "UNKNOWN"},
	}
	for _, c := range cases {
		if got := c.r.String(); got != c.want {
			t.Errorf("Role(%d).String() = %q, want %q", c.r, got, c.want)
		}
	}
}

// TestFormatBestLine_Compact is the one-line compact formatter used in test error messages —
// just a comma-separated "card: ROLE" list with a " (from arsenal)" tag on arsenal-in entries.
func TestFormatBestLine_Compact(t *testing.T) {
	line := []CardAssignment{
		{Card: fake.RedAttack{}, Role: Pitch},
		{Card: fake.RedAttack{}, Role: Attack},
		{Card: generic.ToughenUpBlue{}, Role: Defend, FromArsenal: true},
	}
	got := FormatBestLine(line)
	want := "cardtest.RedAttack: PITCH, cardtest.RedAttack: ATTACK, Toughen Up (Blue) (from arsenal): DEFEND"
	if got != want {
		t.Errorf("FormatBestLine = %q\n  want = %q", got, want)
	}
}

// TestFormatBestTurn_AttackAndPitch verifies the basic numbered-list shape: pitches come first
// as "PITCH (my turn)" (no damage tag — resource isn't damage), attacks come from AttackChain
// with their Play damage. Hand: 2 Red Attacks + 2 Blue (one pitched to pay, one Held).
func TestFormatBestTurn_AttackAndPitch(t *testing.T) {
	h := []card.Card{fake.BlueAttack{}, fake.BlueAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	got := Best(stubHero{}, nil, h, 0, nil, 0, nil)
	out := FormatBestTurn(got)
	// There should be exactly one "PITCH (my turn)" line — a single Blue (pitch 3) covers the
	// 1-cost Red attacks' combined cost of 2. The other Blue ends up Held.
	if n := strings.Count(out, "PITCH (my turn)"); n != 1 {
		t.Errorf("want 1 'PITCH (my turn)' line, got %d in:\n%s", n, out)
	}
	// No defense phase → no "PITCH (opponent's turn)" line.
	if strings.Contains(out, "opponent's turn") {
		t.Errorf("didn't expect defense-phase pitch in:\n%s", out)
	}
	// Both Red attacks show up with their Attack() damage of 3.
	if n := strings.Count(out, ": ATTACK (+3)"); n != 2 {
		t.Errorf("want 2 Red attacks at +3 damage, got %d in:\n%s", n, out)
	}
}

// TestFormatBestTurn_HeroTriggerAttribution exercises the explicit hero-trigger line on an
// attack that fires OnCardPlayed. Viserai creates a Runechant on the 2nd+ non-attack action
// played; with a non-attack action first in chain, the next attack's card slot gets a
// "(+M hero trigger)" suffix rather than silently folding M into the attack's damage number.
func TestFormatBestTurn_HeroTriggerAttribution(t *testing.T) {
	h := []card.Card{runeblade.MauvrionSkiesRed{}, runeblade.ShrillOfSkullformRed{}, runeblade.MaleficIncantationBlue{}}
	got := Best(hero.Viserai{}, nil, h, 0, nil, 0, nil)
	out := FormatBestTurn(got)
	// The winning line plays Mauvrion (non-attack action) → Shrill (attack); Viserai triggers on
	// Shrill for +1 runechant and the display should split that off into its own tag.
	if !strings.Contains(out, "hero trigger") {
		t.Errorf("want a hero trigger tag on the chain, got:\n%s", out)
	}
}

// TestFormatBestTurn_ArsenalInPlayedAsDR checks the combined "arsenal-in played from the slot"
// + "defense reaction prevented" rendering. Hand: one Malefic Blue (pitch 3). Arsenal-in:
// Toughen Up Blue (DR cost 2). Malefic pitches to fund the DR, Toughen Up blocks 4 of 4 incoming.
// Display should put the pitch under "opponent's turn" and tag Toughen Up with "(from arsenal)"
// on the DEFENSE REACTION line.
func TestFormatBestTurn_ArsenalInPlayedAsDR(t *testing.T) {
	h := []card.Card{runeblade.MaleficIncantationBlue{}}
	got := Best(stubHero{}, nil, h, 4, nil, 0, generic.ToughenUpBlue{})
	out := FormatBestTurn(got)
	if !strings.Contains(out, "PITCH (opponent's turn)") {
		t.Errorf("want a defense-phase pitch line, got:\n%s", out)
	}
	if !strings.Contains(out, "Toughen Up (Blue) (from arsenal): DEFENSE REACTION") {
		t.Errorf("want the DR tagged as 'from arsenal', got:\n%s", out)
	}
	if !strings.Contains(out, "+4 prevented") {
		t.Errorf("want '+4 prevented' (4 incoming fully blocked by defense 4), got:\n%s", out)
	}
}

// TestFormatBestTurn_WeaponSwingInChain makes sure a swung weapon shows up in the chain with a
// WEAPON ATTACK label and its damage.
func TestFormatBestTurn_WeaponSwingInChain(t *testing.T) {
	h := []card.Card{fake.RedAttack{}}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(stubHero{}, weapons, h, 0, nil, 0, nil)
	out := FormatBestTurn(got)
	// Reaping Blade attack is 3.
	if !strings.Contains(out, "Reaping Blade: WEAPON ATTACK (+3)") {
		t.Errorf("want the weapon in the chain, got:\n%s", out)
	}
}

// TestFormatBestTurn_HeldAndArsenalFooter covers the trailing "(held: …)" / "(arsenal: …)"
// bookkeeping. A lone DR (no way to pay its cost, no incoming) is Held in the partition but
// then promoted to Arsenal post-hoc, so the output shows an arsenal footer (not a held footer).
func TestFormatBestTurn_HeldAndArsenalFooter(t *testing.T) {
	h := []card.Card{generic.ToughenUpBlue{}}
	got := Best(stubHero{}, nil, h, 4, nil, 0, nil)
	out := FormatBestTurn(got)
	if !strings.Contains(out, "(arsenal: Toughen Up (Blue) (new))") {
		t.Errorf("want an arsenal footer tagged '(new)', got:\n%s", out)
	}
}

// TestFormatBestTurn_StayingArsenalFooter tags the carrying-over arsenal card with "(stayed)"
// rather than "(new)" — useful for the reader to see the slot wasn't swapped this turn.
func TestFormatBestTurn_StayingArsenalFooter(t *testing.T) {
	// Hand with no attacks / no pitches to pay for the arsenal DR at incoming=0 (defense is
	// wasted anyway). Arsenal-in Toughen Up sits.
	h := []card.Card{generic.ToughenUpBlue{}}
	got := Best(stubHero{}, nil, h, 0, nil, 0, generic.ToughenUpBlue{})
	out := FormatBestTurn(got)
	if !strings.Contains(out, "(stayed)") {
		t.Errorf("want the arsenal-in card tagged '(stayed)', got:\n%s", out)
	}
}

// TestFormatBestTurn_EmptyBestLine covers the degenerate path — zero cards produces no output
// lines. Exercised by plugging an empty summary directly into the formatter.
func TestFormatBestTurn_EmptyBestLine(t *testing.T) {
	if got := FormatBestTurn(TurnSummary{}); got != "" {
		t.Errorf("empty summary should render as empty string, got %q", got)
	}
}

// TestFormatBestTurn_DelayedFromLastTurnLine surfaces cross-turn PlayNextTurn contributions on
// their own "(from previous turn)" line so the reader sees where the damage-equivalent came
// from (the Value would otherwise not reconcile with the on-turn per-card breakdown).
func TestFormatBestTurn_DelayedFromLastTurnLine(t *testing.T) {
	t.Run("damage credit", func(t *testing.T) {
		summary := TurnSummary{
			DelayedFromLastTurn: []DelayedContribution{
				{Card: fake.RedAttack{}, Damage: 3},
			},
		}
		out := FormatBestTurn(summary)
		want := "1. cardtest.RedAttack (from previous turn): START OF ACTION PHASE (+3)"
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	})
	t.Run("reveal into hand", func(t *testing.T) {
		summary := TurnSummary{
			DelayedFromLastTurn: []DelayedContribution{
				{Card: fake.RedAttack{}, ToHand: fake.BlueAttack{}},
			},
		}
		out := FormatBestTurn(summary)
		want := "1. cardtest.RedAttack (from previous turn): REVEALED cardtest.BlueAttack INTO HAND"
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	})
}

// TestFormatBestTurn_DrawnCardsRendered pins each role a drawn card can take to a tagged line
// in the printout, so per-card attribution reconciles with the turn's Value rather than
// silently folding extension damage / pitch-from-drawn into the total. The summary is
// hand-built (rather than going through Best) so the test exercises only the formatter.
func TestFormatBestTurn_DrawnCardsRendered(t *testing.T) {
	t.Run("attack extension shows on numbered timeline", func(t *testing.T) {
		summary := TurnSummary{
			Drawn: []CardAssignment{
				{Card: fake.RedAttack{}, Role: Attack, Contribution: 3},
			},
		}
		out := FormatBestTurn(summary)
		want := "1. cardtest.RedAttack (drawn): ATTACK (+3)"
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	})

	t.Run("pitch shows on numbered timeline", func(t *testing.T) {
		summary := TurnSummary{
			Drawn: []CardAssignment{
				{Card: fake.BlueAttack{}, Role: Pitch, Contribution: 3},
			},
		}
		out := FormatBestTurn(summary)
		want := "1. cardtest.BlueAttack (drawn): PITCH (my turn)"
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	})

	t.Run("held lands in the footer", func(t *testing.T) {
		summary := TurnSummary{
			Drawn: []CardAssignment{
				{Card: fake.RedAttack{}, Role: Held},
			},
		}
		out := FormatBestTurn(summary)
		want := "(held: cardtest.RedAttack (drawn))"
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	})

	t.Run("arsenal lands in the footer", func(t *testing.T) {
		summary := TurnSummary{
			Drawn: []CardAssignment{
				{Card: fake.RedAttack{}, Role: Arsenal},
			},
		}
		out := FormatBestTurn(summary)
		want := "(arsenal: cardtest.RedAttack (drawn))"
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	})
}

// TestFormatContribution_IntegerVsFractional covers the small helper that chooses between
// integer and single-decimal rendering. Defense-share contributions can be fractional (e.g. 2
// blockers splitting 3 incoming → 1.5 each).
func TestFormatContribution_IntegerVsFractional(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, "0"},
		{3, "3"},
		{1.5, "1.5"},
		{0.5, "0.5"},
	}
	for _, c := range cases {
		if got := formatContribution(c.in); got != c.want {
			t.Errorf("formatContribution(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}
