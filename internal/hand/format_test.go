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
// attack chain both fall under the "My turn:" section header. Hand: 2 Red Attacks + 2 Blues.
// One Blue pitches for 3 resource, funding the 3-cost chain (Blue + Red + Red, all cost 1
// each, all go-again).
func TestFormatBestTurn_AttackAndPitch(t *testing.T) {
	h := []card.Card{fake.BlueAttack{}, fake.BlueAttack{}, fake.RedAttack{}, fake.RedAttack{}}
	got := Best(stubHero, nil, h, 0, nil, 0, nil)
	out := FormatBestTurn(got, 0)
	if !strings.Contains(out, "  My turn:") {
		t.Errorf("want 'My turn:' section header, got:\n%s", out)
	}
	// Exactly one PITCH line — one Blue funds the 3-cost chain.
	if n := strings.Count(out, ": PITCH"); n != 1 {
		t.Errorf("want 1 ': PITCH' line, got %d in:\n%s", n, out)
	}
	// No defense phase → no "Opponent's turn:" section at all.
	if strings.Contains(out, "Opponent's turn:") {
		t.Errorf("didn't expect defense-phase section in:\n%s", out)
	}
	// Three ATTACK lines: 1 Blue + 2 Reds chain on go-again.
	if n := strings.Count(out, ": ATTACK"); n != 3 {
		t.Errorf("want 3 ': ATTACK' lines, got %d in:\n%s", n, out)
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

// TestFormatBestTurn_LogAttributesEachTriggerSeparately pins phase-2 logging: each event in
// the chain (card Play, hero trigger, mid-chain aura trigger, ephemeral attack trigger) gets
// its own attributed line. Hand: Nimblism (Generic non-attack action, sets
// NonAttackActionPlayed) + Mauvrion (Runeblade non-attack action, registers an ephemeral
// "if hits, +3" trigger) + Consuming Volition (Runeblade attack action). Prior aura: a
// Malefic Incantation TriggerAttackAction. The chain log should hit:
//   - Mauvrion: PLAY (+0)            — non-attack action, no own damage
//   - Consuming Volition: ATTACK (+N) — printed power
//   - Viserai: HERO TRIGGER (+1)     — fires on Volition (attack action)
//   - Malefic Incantation: AURA TRIGGER (+1) — carryover aura fires on Volition
//   - Mauvrion: ATTACK TRIGGER (+3)  — ephemeral fires on Volition's hit
func TestFormatBestTurn_LogAttributesEachTriggerSeparately(t *testing.T) {
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
	wants := []string{
		"Consuming Volition (Red): ATTACK",
		"Viserai: HERO TRIGGER (+1)",
		"Malefic Incantation (Red): AURA TRIGGER (+1)",
		"Mauvrion Skies (Red): ATTACK TRIGGER (+3)",
	}
	for _, want := range wants {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}

// TestFormatBestTurn_LogSuppressesZeroTriggers pins that handlers returning 0 don't add a
// "(+0)" line for the trigger — log entries for triggers gate on positive contribution to
// avoid noise. Card-Play lines render unconditionally because the card resolved (the line
// proves the chain step happened); trigger lines only render when the trigger actually
// credited damage.
func TestFormatBestTurn_LogSuppressesZeroTriggers(t *testing.T) {
	// Hand: a single Red attack with no go-again. Viserai's OnCardPlayed contributes nothing
	// (the gate needs another non-attack action played first), no priors, no ephemerals — so
	// the chain log should be exactly one card-Play line, no trigger spam.
	h := []card.Card{fake.RedAttack{}}
	got := Best(hero.Viserai{}, nil, h, 0, nil, 0, nil)
	if strings.Contains(FormatBestTurn(got, 0), "HERO TRIGGER") {
		t.Errorf("hero trigger line shouldn't render when Viserai contributed 0; got:\n%s",
			FormatBestTurn(got, 0))
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
	if !strings.Contains(out, "cardtest.RedAttack: ATTACK from arsenal") {
		t.Errorf("want 'ATTACK from arsenal' on the role label, got:\n%s", out)
	}
	// The arsenal tag must hang off the role label, not the card name.
	if strings.Contains(out, "cardtest.RedAttack (from arsenal)") {
		t.Errorf("arsenal tag should live on the role label, not the card name; got:\n%s", out)
	}
}

// TestFormatBestTurn_WeaponSwingInChain makes sure a swung weapon shows up in the chain with
// a WEAPON ATTACK label, sourced from the dispatcher's Log entry (chainVerbFor's TypeWeapon
// branch). The State.Log assertion pins the dispatcher → log → format pipeline for weapons
// since FormatBestTurn now reads weapon swings from State.Log rather than SwungWeapons.
func TestFormatBestTurn_WeaponSwingInChain(t *testing.T) {
	h := []card.Card{fake.RedAttack{}}
	weapons := []weapon.Weapon{weapon.ReapingBlade{}}
	got := Best(stubHero, weapons, h, 0, nil, 0, nil)
	out := FormatBestTurn(got, 0)
	if !strings.Contains(out, "Reaping Blade: WEAPON ATTACK") {
		t.Errorf("want the weapon in the chain, got:\n%s", out)
	}
	var sawWeaponLog bool
	for _, line := range got.State.Log {
		if strings.Contains(line, "Reaping Blade: WEAPON ATTACK") {
			sawWeaponLog = true
			break
		}
	}
	if !sawWeaponLog {
		t.Errorf("State.Log missing the weapon swing entry; format-layer match was a fluke. Log=%v", got.State.Log)
	}
}

// TestFormatBestTurn_EndOfTurnArsenalNew pins the End of turn section's arsenal entry tagged
// "(new)" when a Held hand card got promoted into an empty arsenal slot post-hoc. A lone DR
// (no way to pay its cost, no incoming) is Held in the partition but then promoted to
// Arsenal, so End of turn shows "Arsenal: Toughen Up (Blue) (new)".
func TestFormatBestTurn_EndOfTurnArsenalNew(t *testing.T) {
	h := []card.Card{generic.ToughenUpBlue{}}
	got := Best(stubHero, nil, h, 4, nil, 0, nil)
	out := FormatBestTurn(got, 0)
	if !strings.Contains(out, "Arsenal: Toughen Up (Blue) (new)") {
		t.Errorf("want an end-of-turn arsenal entry tagged '(new)', got:\n%s", out)
	}
}

// TestFormatBestTurn_EndOfTurnArsenalStayed tags the carrying-over arsenal card with
// "(stayed)" rather than "(new)" — useful for the reader to see the slot wasn't swapped.
func TestFormatBestTurn_EndOfTurnArsenalStayed(t *testing.T) {
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
// under the "Start of turn:" section as unnumbered entries — the reveal / damage fires at
// the top of the action phase before the chain runs.
func TestFormatBestTurn_TriggersFromLastTurnLine(t *testing.T) {
	summary := TurnSummary{
		TriggersFromLastTurn: []TriggerContribution{
			{Card: fake.RedAttack{}, Damage: 3},
		},
	}
	out := FormatBestTurn(summary, 0)
	want := "cardtest.RedAttack: START OF ACTION PHASE (+3)"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// TestFormatBestTurn_TriggersFromLastTurnRevealedLine surfaces the card a trigger handler
// revealed into the hand. Sigil of the Arknight fires at start of action phase with
// Damage=0 but reveals the deck top; the Start of turn section names the card it drew.
func TestFormatBestTurn_TriggersFromLastTurnRevealedLine(t *testing.T) {
	summary := TurnSummary{
		TriggersFromLastTurn: []TriggerContribution{
			{Card: runeblade.SigilOfTheArknightBlue{}, Revealed: runeblade.MauvrionSkiesRed{}},
		},
	}
	out := FormatBestTurn(summary, 0)
	want := "Sigil of the Arknight (Blue): drew Mauvrion Skies (Red) into hand"
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

// TestFormatBestTurn_StartOfTurnAurasLine pins the Start of turn section's "Auras: ..."
// entry. Names sort alphabetically for determinism, and duplicates are preserved (two copies
// of the same aura render twice).
func TestFormatBestTurn_StartOfTurnAurasLine(t *testing.T) {
	summary := TurnSummary{
		StartOfTurnAuras: []card.Card{
			runeblade.MaleficIncantationRed{},
			runeblade.MaleficIncantationRed{},
			runeblade.SigilOfTheArknightBlue{},
		},
	}
	out := FormatBestTurn(summary, 0)
	want := "Auras: Malefic Incantation (Red), Malefic Incantation (Red), Sigil of the Arknight (Blue)"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// TestFormatBestTurn_StartOfTurnAurasSuppressedWhenEmpty pins the omission of the Auras
// entry when no auras were in play and no starting runechants carry in — the empty state
// shouldn't render a dangling "Auras:" label.
func TestFormatBestTurn_StartOfTurnAurasSuppressedWhenEmpty(t *testing.T) {
	summary := TurnSummary{BestLine: []CardAssignment{{Card: fake.RedAttack{}, Role: Attack}}}
	out := FormatBestTurn(summary, 0)
	if strings.Contains(out, "Auras: ") {
		t.Errorf("unexpected Auras line in output:\n%s", out)
	}
}

// TestFormatBestTurn_StartOfTurnAurasWithRunechants folds a non-zero starting Runechant
// carry into the "Auras:" entry as the trailing item — a Runeblade hero carrying tokens
// from the previous turn sees them alongside any auras as one combined readout.
func TestFormatBestTurn_StartOfTurnAurasWithRunechants(t *testing.T) {
	summary := TurnSummary{
		StartOfTurnAuras: []card.Card{runeblade.MaleficIncantationRed{}},
	}
	out := FormatBestTurn(summary, 3)
	want := "Auras: Malefic Incantation (Red), 3 Runechants"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// TestFormatBestTurn_StartOfTurnRunechantsOnly folds a non-zero starting Runechant carry
// into the Auras entry even when no auras are in play, using singular "Runechant" when the
// count is 1.
func TestFormatBestTurn_StartOfTurnRunechantsOnly(t *testing.T) {
	out := FormatBestTurn(TurnSummary{}, 1)
	want := "Auras: 1 Runechant"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
	// Plural noun when count > 1.
	out2 := FormatBestTurn(TurnSummary{}, 2)
	if !strings.Contains(out2, "2 Runechants") {
		t.Errorf("want plural 'Runechants' at count 2, got:\n%s", out2)
	}
}

// TestFormatBestTurn_EndOfTurnHandLine pins the End of turn "Hand: ..." entry — every card
// in State.Hand surfaces as part of one comma-separated line, regardless of whether it
// started the turn in hand or got drawn / tutored mid-chain.
func TestFormatBestTurn_EndOfTurnHandLine(t *testing.T) {
	summary := TurnSummary{
		State: CarryState{Hand: []card.Card{fake.RedAttack{}}},
	}
	out := FormatBestTurn(summary, 0)
	want := "Hand: cardtest.RedAttack"
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
