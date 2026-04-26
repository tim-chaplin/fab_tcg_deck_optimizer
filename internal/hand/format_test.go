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

// TestFormatBestTurn_AttackAndPitch verifies the basic numbered-list shape: pitches and the
// attack chain both fall under the "My turn:" section header, pitches render as ": PITCH"
// (no damage tag — resource isn't damage, and the section header disambiguates phase),
// attacks come from AttackChain with their Play damage. Hand: 2 Red Attacks + 2 Blue (one
// pitched to pay, one Held).
func TestFormatBestTurn_AttackAndPitch(t *testing.T) {
	h := []card.Card{fake.BlueAttack{}, fake.BlueAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	got := Best(stubHero, nil, h, 0, nil, 0, nil)
	out := FormatBestTurn(got, 0)
	if !strings.Contains(out, "  My turn:") {
		t.Errorf("want 'My turn:' section header, got:\n%s", out)
	}
	// Exactly one PITCH line — a single Blue (pitch 3) covers the 1-cost Red attacks'
	// combined cost of 2; the other Blue ends up Held.
	if n := strings.Count(out, ": PITCH"); n != 1 {
		t.Errorf("want 1 ': PITCH' line, got %d in:\n%s", n, out)
	}
	// No defense phase → no "Opponent's turn:" section at all.
	if strings.Contains(out, "Opponent's turn:") {
		t.Errorf("didn't expect defense-phase section in:\n%s", out)
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
	out := FormatBestTurn(got, 0)
	// The winning line plays Mauvrion (non-attack action) → Shrill (attack); Viserai triggers on
	// Shrill for +1 runechant and the display should split that off into its own tag.
	if !strings.Contains(out, "hero trigger") {
		t.Errorf("want a hero trigger tag on the chain, got:\n%s", out)
	}
}

// TestFormatBestTurn_NonAttackCardUsesPlayLabel pins the chain label to "PLAY" for cards that
// aren't attacks (e.g. Mauvrion Skies, a non-attack action). Attack cards keep the "ATTACK"
// label so the reader can distinguish damage-dealing chain steps from resource/setup plays.
func TestFormatBestTurn_NonAttackCardUsesPlayLabel(t *testing.T) {
	h := []card.Card{runeblade.MauvrionSkiesRed{}, runeblade.ShrillOfSkullformRed{}, runeblade.MaleficIncantationBlue{}}
	got := Best(hero.Viserai{}, nil, h, 0, nil, 0, nil)
	out := FormatBestTurn(got, 0)
	if !strings.Contains(out, "Mauvrion Skies (Red): PLAY") {
		t.Errorf("want Mauvrion (non-attack action) labelled PLAY, got:\n%s", out)
	}
	if !strings.Contains(out, "Shrill of Skullform (Red): ATTACK") {
		t.Errorf("want Shrill (attack action) labelled ATTACK, got:\n%s", out)
	}
}

// TestFormatBestTurn_AuraTriggerLabelledSeparately pins the chain display to attribute hero
// OnCardPlayed damage and mid-chain AuraTrigger damage as distinct comma-separated items
// inside a single trigger parenthesised group. When a Runeblade attack action triggers both
// Viserai (hero ability, +1 Runechant) and a carryover Malefic Incantation
// (TriggerAttackAction aura, +1 Runechant), the card line reads "(+1 hero trigger, +1 aura
// trigger)" — one parenthesised group so the eye tracks two related items together.
func TestFormatBestTurn_AuraTriggerLabelledSeparately(t *testing.T) {
	// Hand: a Generic non-attack action (Nimblism — sets NonAttackActionPlayed without being
	// Runeblade), a Runeblade non-attack action (Mauvrion — Viserai fires on it), and a
	// Runeblade attack action (Consuming Volition — Viserai fires AND the carryover Malefic
	// trigger fires). Viserai's contribution is +1; the carryover aura's is +1; the display
	// must attribute them to distinct items inside one combined group.
	h := []card.Card{generic.NimblismRed{}, runeblade.MauvrionSkiesRed{}, runeblade.ConsumingVolitionRed{}}
	prior := []card.AuraTrigger{{
		Self:        runeblade.MaleficIncantationRed{},
		Type:        card.TriggerAttackAction,
		Count:       2,
		OncePerTurn: true,
		Handler:     func(s *card.TurnState) int { return s.CreateRunechants(1) },
	}}
	got := BestWithTriggers(hero.Viserai{}, nil, h, 0, nil, 0, nil, prior)
	out := FormatBestTurn(got, 0)
	if !strings.Contains(out, "(+1 hero trigger, +1 aura trigger)") {
		t.Errorf("want combined '(+1 hero trigger, +1 aura trigger)' tag, got:\n%s", out)
	}
	if strings.Contains(out, "(+2 hero trigger)") {
		t.Errorf("the +1 aura damage must not appear under the hero tag; got:\n%s", out)
	}
}

// TestFormatBestTurn_ArsenalInPlayedAsDR checks the combined "arsenal-in played from the slot"
// + "defense reaction prevented" rendering. Hand: one Malefic Blue (pitch 3). Arsenal-in:
// Toughen Up Blue (DR cost 2). Malefic pitches to fund the DR, Toughen Up blocks 4 of 4 incoming.
// Display puts the pitch and DR lines under the "Opponent's turn:" section; the role label
// reads "DEFENSE REACTION from arsenal" since Toughen Up came out of the arsenal slot.
func TestFormatBestTurn_ArsenalInPlayedAsDR(t *testing.T) {
	h := []card.Card{runeblade.MaleficIncantationBlue{}}
	got := Best(stubHero, nil, h, 4, nil, 0, generic.ToughenUpBlue{})
	out := FormatBestTurn(got, 0)
	if !strings.Contains(out, "  Opponent's turn:") {
		t.Errorf("want 'Opponent's turn:' section header, got:\n%s", out)
	}
	if !strings.Contains(out, ": PITCH") {
		t.Errorf("want a defense-phase pitch line, got:\n%s", out)
	}
	if !strings.Contains(out, "Toughen Up (Blue): DEFENSE REACTION from arsenal") {
		t.Errorf("want 'DEFENSE REACTION from arsenal' on the role label, got:\n%s", out)
	}
	if !strings.Contains(out, "+4 prevented") {
		t.Errorf("want '+4 prevented' (4 incoming fully blocked by defense 4), got:\n%s", out)
	}
}

// TestFormatBestTurn_ArsenalInPlayedOnChain checks the role-label tag for an arsenal-in
// card played as part of the my-turn chain. Hand: one BlueAttack (pitch 3, cost 1).
// Arsenal-in: RedAttack (cost 1, attack 3). The solver pitches the Blue to pay the Red's
// cost and attacks from arsenal for 3; the chain line reads "cardtest.RedAttack: ATTACK
// from arsenal" — tag on the role, not on the card name.
func TestFormatBestTurn_ArsenalInPlayedOnChain(t *testing.T) {
	h := []card.Card{fake.BlueAttack{}}
	got := Best(stubHero, nil, h, 0, nil, 0, fake.RedAttack{})
	out := FormatBestTurn(got, 0)
	if !strings.Contains(out, "  My turn:") {
		t.Errorf("want 'My turn:' section header, got:\n%s", out)
	}
	if !strings.Contains(out, "cardtest.RedAttack: ATTACK from arsenal (+3)") {
		t.Errorf("want 'ATTACK from arsenal' on the role label, got:\n%s", out)
	}
	// The arsenal tag must hang off the role label, not the card name.
	if strings.Contains(out, "cardtest.RedAttack (from arsenal)") {
		t.Errorf("arsenal tag should live on the role label, not the card name; got:\n%s", out)
	}
}

// TestFormatBestTurn_WeaponSwingInChain makes sure a swung weapon shows up in the chain with a
// WEAPON ATTACK label and its damage.
func TestFormatBestTurn_WeaponSwingInChain(t *testing.T) {
	h := []card.Card{fake.RedAttack{}}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(stubHero, weapons, h, 0, nil, 0, nil)
	out := FormatBestTurn(got, 0)
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
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	out := FormatBestTurn(got, 0)
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
	got := Best(stubHero, nil, h, 0, nil, 0, generic.ToughenUpBlue{})
	out := FormatBestTurn(got, 0)
	if !strings.Contains(out, "(stayed)") {
		t.Errorf("want the arsenal-in card tagged '(stayed)', got:\n%s", out)
	}
}

// TestFormatBestTurn_EmptyBestLine covers the degenerate path — zero cards produces no output
// lines. Exercised by plugging an empty summary directly into the formatter.
func TestFormatBestTurn_EmptyBestLine(t *testing.T) {
	if got := FormatBestTurn(TurnSummary{}, 0); got != "" {
		t.Errorf("empty summary should render as empty string, got %q", got)
	}
}

// TestFormatBestTurn_TriggersFromLastTurnLine surfaces cross-turn AuraTrigger contributions
// under the "My turn:" section as the first numbered entries — the reveal / damage fires at
// the top of the action phase before pitches resolve. The aura name stands alone (the
// "Auras in play at start of turn" header already names the source, so a "(from previous
// turn)" suffix would just repeat).
func TestFormatBestTurn_TriggersFromLastTurnLine(t *testing.T) {
	summary := TurnSummary{
		TriggersFromLastTurn: []TriggerContribution{
			{Card: fake.RedAttack{}, Damage: 3},
		},
	}
	out := FormatBestTurn(summary, 0)
	want := "1. cardtest.RedAttack: START OF ACTION PHASE (+3)"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// TestFormatBestTurn_TriggersFromLastTurnRevealedLine surfaces the card a trigger handler
// revealed into the hand. Sigil of the Arknight fires at start of action phase with
// Damage=0 but reveals the deck top; the printout names the card it drew instead of a
// bare "(+0)".
func TestFormatBestTurn_TriggersFromLastTurnRevealedLine(t *testing.T) {
	summary := TurnSummary{
		TriggersFromLastTurn: []TriggerContribution{
			{Card: runeblade.SigilOfTheArknightBlue{}, Revealed: runeblade.MauvrionSkiesRed{}},
		},
	}
	out := FormatBestTurn(summary, 0)
	want := "1. Sigil of the Arknight (Blue): drew Mauvrion Skies (Red) into hand"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// TestFormatBestTurn_TriggersFromLastTurnZeroEffectDropped suppresses lines for carryover
// triggers that did nothing visible this turn (zero damage, no reveal). Output has no
// numbered entries at all — the My turn section is empty so its header elides too.
func TestFormatBestTurn_TriggersFromLastTurnZeroEffectDropped(t *testing.T) {
	summary := TurnSummary{
		TriggersFromLastTurn: []TriggerContribution{
			{Card: runeblade.SigilOfTheArknightBlue{}},
		},
	}
	out := FormatBestTurn(summary, 0)
	if out != "" {
		t.Errorf("zero-effect trigger with no other content should render empty; got:\n%s", out)
	}
}

// TestFormatBestTurn_StartOfTurnAurasHeader pins the header line that lists the aura cards
// in play at the top of the turn. Names sort alphabetically for determinism, and duplicates
// are preserved (two copies of the same aura render twice).
func TestFormatBestTurn_StartOfTurnAurasHeader(t *testing.T) {
	summary := TurnSummary{
		StartOfTurnAuras: []card.Card{
			runeblade.MaleficIncantationRed{},
			runeblade.MaleficIncantationRed{},
			runeblade.SigilOfTheArknightBlue{},
		},
	}
	out := FormatBestTurn(summary, 0)
	want := "Auras in play at start of turn: Malefic Incantation (Red), Malefic Incantation (Red), Sigil of the Arknight (Blue)"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// TestFormatBestTurn_StartOfTurnAurasHeaderSuppressedWhenEmpty pins the omission of the header
// line when no auras were in play and no starting runechants carry in — the empty state
// shouldn't render a dangling label.
func TestFormatBestTurn_StartOfTurnAurasHeaderSuppressedWhenEmpty(t *testing.T) {
	summary := TurnSummary{BestLine: []CardAssignment{{Card: fake.RedAttack{}, Role: Attack}}}
	out := FormatBestTurn(summary, 0)
	if strings.Contains(out, "Auras in play at start of turn") {
		t.Errorf("unexpected header in output:\n%s", out)
	}
}

// TestFormatBestTurn_StartOfTurnHeaderWithRunechants folds a non-zero starting Runechant
// carry into the "Auras in play at start of turn" line as the trailing item — a Runeblade
// hero carrying tokens from the previous turn sees them alongside any auras, as one
// combined start-of-turn state readout.
func TestFormatBestTurn_StartOfTurnHeaderWithRunechants(t *testing.T) {
	summary := TurnSummary{
		StartOfTurnAuras: []card.Card{runeblade.MaleficIncantationRed{}},
	}
	out := FormatBestTurn(summary, 3)
	want := "Auras in play at start of turn: Malefic Incantation (Red), 3 Runechants"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// TestFormatBestTurn_StartOfTurnHeaderRunechantsOnly folds a non-zero starting Runechant
// carry into the header even when no auras are in play, using singular "Runechant" when the
// count is 1.
func TestFormatBestTurn_StartOfTurnHeaderRunechantsOnly(t *testing.T) {
	out := FormatBestTurn(TurnSummary{}, 1)
	want := "Auras in play at start of turn: 1 Runechant"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
	// Plural noun when count > 1.
	out2 := FormatBestTurn(TurnSummary{}, 2)
	if !strings.Contains(out2, "2 Runechants") {
		t.Errorf("want plural 'Runechants' at count 2, got:\n%s", out2)
	}
}

// TestFormatBestTurn_HandHeldRenderedInFooter pins the held-footer rendering — every card in
// State.Hand surfaces as one "(held: NAME)" trailing line, regardless of whether it started
// the turn in hand or got drawn / tutored mid-chain (both flavours land in State.Hand).
func TestFormatBestTurn_HandHeldRenderedInFooter(t *testing.T) {
	summary := TurnSummary{
		State: CarryState{Hand: []card.Card{fake.RedAttack{}}},
	}
	out := FormatBestTurn(summary, 0)
	want := "(held: cardtest.RedAttack)"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
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
