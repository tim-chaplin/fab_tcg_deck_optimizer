package sim_test

import (
	"reflect"
	"strings"
	"testing"

	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/token"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
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
		{Card: testutils.RedAttack{}, Role: Pitch},
		{Card: testutils.RedAttack{}, Role: Attack},
		{Card: cards.ToughenUpBlue{}, Role: Defend, FromArsenal: true},
	}
	got := FormatBestLine(line)
	want := "cardtest.RedAttack [R]: PITCH, cardtest.RedAttack [R]: ATTACK, Toughen Up [B] (from arsenal): DEFEND"
	if got != want {
		t.Errorf("FormatBestLine = %q\n  want = %q", got, want)
	}
}

// Tests that pitches and the attack chain both render under the "My turn:" section header.
func TestFormatBestTurn_AttackAndPitch(t *testing.T) {
	h := []card.Card{testutils.BlueAttack{}, testutils.BlueAttack{}, testutils.RedAttack{}, testutils.RedAttack{}}
	got := Best(nil, h, Matchup{IncomingDamage: 0}, nil, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	out := FormatBestTurn(got, nil, nil)
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
	h := []card.Card{cards.MauvrionSkiesRed{}, cards.ShrillOfSkullformRed{}, cards.MaleficIncantationBlue{}}
	got := Best(nil, h, Matchup{IncomingDamage: 0}, nil, gameengine.GameStateBuilder().SetHero(heroes.Viserai{}).Build())
	out := FormatBestTurn(got, nil, nil)
	if !strings.Contains(out, "Mauvrion Skies [R]: PLAY") {
		t.Errorf("want Mauvrion (non-attack action) labelled PLAY, got:\n%s", out)
	}
	if !strings.Contains(out, "Shrill of Skullform [R]: ATTACK") {
		t.Errorf("want Shrill (attack action) labelled ATTACK, got:\n%s", out)
	}
}

// Tests log attribution: each chain event (Play, hero trigger, aura trigger, OnHit) gets its
// own line grouped under the triggering card; "(from ...)" suffix is dropped under indentation.
func TestFormatBestTurn_LogAttributesEachTriggerSeparately(t *testing.T) {
	h := []card.Card{cards.NimblismRed{}, cards.MauvrionSkiesRed{}, cards.ConsumingVolitionRed{}}
	// Use the real Malefic Incantation card's Play to register the prior trigger so the
	// handler matches production exactly (logs via AddPreTriggerLogEntry, sources from
	// state.TriggeringCard).
	bootstrap := gameengine.New()
	bootstrap.ResolveChainStep(bootstrap.Logger(), &card.CardState{Card: cards.MaleficIncantationRed{}})
	state := gameengine.GameStateBuilder().SetHero(heroes.Viserai{}).Build()
	for _, a := range bootstrap.Auras() {
		state.CreateAura(a)
	}
	got := Best(nil, h, Matchup{}, nil, state)
	out := FormatBestTurn(got, nil, nil)
	// Trigger lines render indented (9 spaces) with no "(from <source>)" suffix — the
	// indentation under the parent chain entry conveys attribution. Each line carries
	// the verb phrase the trigger handler authored.
	wants := []string{
		"Consuming Volition [R]: ATTACK",
		"         Viserai created a runechant (+1)",
		"         Malefic Incantation [R] created a runechant (+1)",
		"         Mauvrion Skies [R] created 3 runechants on hit (+3)",
	}
	for _, want := range wants {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	// "(from <source>)" attribution belongs only to orphan triggers; grouped triggers
	// drop it.
	if strings.Contains(out, "(from Consuming Volition") {
		t.Errorf("grouped trigger should not carry '(from <source>)' suffix; got:\n%s", out)
	}
	// Trigger lines are card-authored freeform text — generic dispatcher verbs like
	// "HERO/AURA/ATTACK TRIGGER" only appear when a centralised format leaks back in.
	for _, gone := range []string{"HERO TRIGGER", "AURA TRIGGER", "ATTACK TRIGGER"} {
		if strings.Contains(out, gone) {
			t.Errorf("trigger line still uses generic %q verb; got:\n%s", gone, out)
		}
	}
}

// Tests that trigger handlers returning 0 don't render a "(+0)" line; card-Play lines always
// render, trigger lines only render on positive contribution.
func TestFormatBestTurn_LogSuppressesZeroTriggers(t *testing.T) {
	h := []card.Card{testutils.RedAttack{}}
	got := Best(nil, h, Matchup{IncomingDamage: 0}, nil, gameengine.GameStateBuilder().SetHero(heroes.Viserai{}).Build())
	if strings.Contains(FormatBestTurn(got, nil, nil), "Viserai created") {
		t.Errorf("hero trigger line shouldn't render when Viserai contributed 0; got:\n%s",
			FormatBestTurn(got, nil, nil))
	}
}

// Tests Moon Wish go-again log shape: tutor line groups under Moon Wish, Sun Kiss authors
// its own PLAY chain entry, heal renders as a "Gained 3 health" child under Sun Kiss.
func TestFormatBestTurn_MoonWishTutorAndPlayLogsAsPostTrigger(t *testing.T) {
	h := []card.Card{cards.FlyingHighRed{}, cards.MoonWishYellow{}, testutils.BlueAttack{}}
	deck := DeckOf(cards.SunKissRed{})
	got := Best(nil, h, Matchup{IncomingDamage: 0}, deck, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	out := FormatBestTurn(got, nil, nil)
	wants := []string{
		"Moon Wish [Y]: ATTACK (+4)",
		"         Moon Wish [Y] tutored Sun Kiss [R] and played it",
		"Sun Kiss [R]: PLAY",
		"         Gained 3 health (+3)",
	}
	for _, want := range wants {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	// The chain entry must NOT bundle Sun Kiss's damage into Moon Wish's (+N), and Sun
	// Kiss's chain entry must NOT bundle the heal into its own (+N).
	if strings.Contains(out, "Moon Wish [Y]: ATTACK (+7)") {
		t.Errorf("chain entry bundled Sun Kiss damage; got:\n%s", out)
	}
	if strings.Contains(out, "Sun Kiss [R]: PLAY (+3)") {
		t.Errorf("Sun Kiss chain entry bundled heal into (+N) instead of using a child line\n%s", out)
	}
}

// Tests Moon Wish tutor-only log: post-trigger line renders without a (+N) when nothing
// is credited.
func TestFormatBestTurn_MoonWishTutorOnlyLogsAsPostTrigger(t *testing.T) {
	h := []card.Card{cards.MoonWishYellow{}, testutils.BlueAttack{}}
	deck := DeckOf(cards.SunKissRed{})
	got := Best(nil, h, Matchup{IncomingDamage: 0}, deck, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	out := FormatBestTurn(got, nil, nil)
	wants := []string{
		"Moon Wish [Y]: ATTACK (+4)",
		"         Moon Wish [Y] tutored Sun Kiss [R]",
	}
	for _, want := range wants {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "tutored Sun Kiss [R] (+") {
		t.Errorf("tutor-only line shouldn't carry a (+N) suffix; got:\n%s", out)
	}
}

// Tests rendering of an arsenal-in DR played: section "Opponent's turn:" with role label
// "DEFENSE REACTION from arsenal".
func TestFormatBestTurn_ArsenalInPlayedAsDR(t *testing.T) {
	h := []card.Card{cards.MaleficIncantationBlue{}}
	state := gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetArsenal(cards.ToughenUpBlue{}).
		Build()
	got := Best(nil, h, Matchup{IncomingDamage: 4}, nil, state)
	out := FormatBestTurn(got, nil, nil)
	if !strings.Contains(out, "  Opponent's turn:") {
		t.Errorf("want 'Opponent's turn:' section header, got:\n%s", out)
	}
	if !strings.Contains(out, ": PITCH") {
		t.Errorf("want a defense-phase pitch line, got:\n%s", out)
	}
	if !strings.Contains(out, "Toughen Up [B]: DEFENSE REACTION from arsenal") {
		t.Errorf("want 'DEFENSE REACTION from arsenal' on the role label, got:\n%s", out)
	}
}

// Tests DR rendering: chain step's (+N) folds in BonusDefense; separable riders render
// as indented sub-lines.
func TestFormatBestTurn_DefenseReactionLinesAndRiders(t *testing.T) {
	cases := []struct {
		name     string
		hand     []card.Card
		incoming int
		wants    []string
	}{
		{
			name:     "Sigil of Suffering folds bonus into chain step + arcane sub-line",
			hand:     []card.Card{cards.SigilOfSufferingRed{}, testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}},
			incoming: 4,
			wants: []string{
				"Sigil of Suffering [R]: DEFENSE REACTION (+4)",
				"Dealt 1 arcane damage (+1)",
			},
		},
		{
			name:     "Dodge has no riders, single chain line",
			hand:     []card.Card{cards.DodgeBlue{}},
			incoming: 2,
			wants:    []string{"Dodge [B]: DEFENSE REACTION (+2)"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Best(nil, tc.hand, Matchup{IncomingDamage: tc.incoming}, nil, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
			out := FormatBestTurn(got, nil, nil)
			for _, w := range tc.wants {
				if !strings.Contains(out, w) {
					t.Errorf("want %q in:\n%s", w, out)
				}
			}
		})
	}
}

// Tests that an arsenal-in card played on the my-turn chain renders as "ATTACK from arsenal"
// (role tag, not card-name tag).
func TestFormatBestTurn_ArsenalInPlayedOnChain(t *testing.T) {
	h := []card.Card{testutils.BlueAttack{}}
	state := gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetArsenal(testutils.RedAttack{}).
		Build()
	got := Best(nil, h, Matchup{IncomingDamage: 0}, nil, state)
	out := FormatBestTurn(got, nil, nil)
	if !strings.Contains(out, "  My turn:") {
		t.Errorf("want 'My turn:' section header, got:\n%s", out)
	}
	if !strings.Contains(out, "cardtest.RedAttack [R]: ATTACK from arsenal") {
		t.Errorf("want 'ATTACK from arsenal' on the role label, got:\n%s", out)
	}
	// The arsenal tag must hang off the role label, not the card name.
	if strings.Contains(out, "cardtest.RedAttack [R] (from arsenal)") {
		t.Errorf("arsenal tag should live on the role label, not the card name; got:\n%s", out)
	}
}

// Tests that a swung weapon shows up in the chain as "WEAPON ATTACK" sourced from
// State.Log (FormatBestTurn reads weapon swings from State.Log, not SwungWeapons).
func TestFormatBestTurn_WeaponSwingInChain(t *testing.T) {
	h := []card.Card{testutils.RedAttack{}}
	weapons := []Weapon{weapons.ReapingBlade{}}
	got := Best(weapons, h, Matchup{IncomingDamage: 0}, nil, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	out := FormatBestTurn(got, nil, nil)
	if !strings.Contains(out, "Reaping Blade: WEAPON ATTACK") {
		t.Errorf("want the weapon in the chain, got:\n%s", out)
	}
	var sawWeaponLog bool
	for _, e := range got.State.LogEntries() {
		if strings.Contains(e.Text, "Reaping Blade: WEAPON ATTACK") {
			sawWeaponLog = true
			break
		}
	}
	if !sawWeaponLog {
		t.Errorf("State.Log missing the weapon swing entry; format-layer match was a fluke. Log=%v", got.State.LogEntries())
	}
}

// Tests that the End of turn section tags a post-hoc-promoted Held → Arsenal card as "(new)".
func TestFormatBestTurn_EndOfTurnArsenalNew(t *testing.T) {
	h := []card.Card{cards.ToughenUpBlue{}}
	got := Best(nil, h, Matchup{IncomingDamage: 4}, nil, gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4}).Build())
	out := FormatBestTurn(got, nil, nil)
	if !strings.Contains(out, "Arsenal: Toughen Up [B] (new)") {
		t.Errorf("want an end-of-turn arsenal entry tagged '(new)', got:\n%s", out)
	}
}

// TestFormatBestTurn_EndOfTurnArsenalStayed tags the carrying-over arsenal card with
// "(stayed)" rather than "(new)" — useful for the reader to see the slot wasn't swapped.
func TestFormatBestTurn_EndOfTurnArsenalStayed(t *testing.T) {
	// Hand with no attacks / no pitches to pay for the arsenal DR at incoming=0 (defense is
	// wasted anyway). Arsenal-in Toughen Up sits.
	h := []card.Card{cards.ToughenUpBlue{}}
	state := gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetArsenal(cards.ToughenUpBlue{}).
		Build()
	got := Best(nil, h, Matchup{IncomingDamage: 0}, nil, state)
	out := FormatBestTurn(got, nil, nil)
	if !strings.Contains(out, "(stayed)") {
		t.Errorf("want the arsenal-in card tagged '(stayed)', got:\n%s", out)
	}
}

// TestFormatBestTurn_EmptyBestLine covers the degenerate path — zero cards produces no output
// lines. Exercised by plugging an empty summary directly into the formatter.
func TestFormatBestTurn_EmptyBestLine(t *testing.T) {
	if got := FormatBestTurn(TurnSummary{}, nil, nil); got != "" {
		t.Errorf("empty summary should render as empty string, got %q", got)
	}
}

// Tests that cross-turn Aura contributions render as numbered entries at the top of "My turn:",
// not in the unnumbered Start of turn block.
func TestFormatBestTurn_TriggersFromLastTurnLine(t *testing.T) {
	summary := TurnSummary{
		TriggersFromLastTurn: []TriggerContribution{
			{Card: testutils.RedAttack{}, Damage: 3},
		},
	}
	out := FormatBestTurn(summary, nil, nil)
	want := "1. cardtest.RedAttack [R]: START OF ACTION PHASE (+3)"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
	if !strings.Contains(out, "  My turn:") {
		t.Errorf("trigger line should sit under 'My turn:', got:\n%s", out)
	}
	if strings.Contains(out, "Start of turn:") {
		t.Errorf("Start of turn section shouldn't render for trigger-only summary; got:\n%s", out)
	}
}

// Tests that "Start of turn → Hand:" reads from DealtHand (not an augmented hand from a
// Sigil reveal); the revealed card appears under MyTurn instead.
func TestFormatBestTurn_StartOfTurnHandReadsDealtHand(t *testing.T) {
	summary := TurnSummary{
		DealtHand: []card.Card{testutils.RedAttack{}},
		BestLine: []CardAssignment{
			{Card: testutils.RedAttack{}, Role: Attack},
			// Mauvrion is in BestLine because the reveal augmented the hand the partition
			// saw, but it never appeared in DealtHand — so it must not show up in the
			// start-of-turn hand line.
			{Card: cards.MauvrionSkiesRed{}, Role: Held},
		},
		TriggersFromLastTurn: []TriggerContribution{
			{Card: cards.SigilOfTheArknightBlue{}, Revealed: cards.MauvrionSkiesRed{}},
		},
	}
	out := FormatBestTurn(summary, nil, nil)
	if !strings.Contains(out, "Hand: cardtest.RedAttack [R]\n") {
		t.Errorf("Hand line should list only DealtHand cards; got:\n%s", out)
	}
	if strings.Contains(out, "Hand: cardtest.RedAttack [R], Mauvrion Skies [R]") {
		t.Errorf("revealed Mauvrion must not appear in start-of-turn hand; got:\n%s", out)
	}
	if !strings.Contains(out, "Sigil of the Arknight [B]: drew Mauvrion Skies [R] into hand") {
		t.Errorf("MyTurn should still record the reveal; got:\n%s", out)
	}
}

// Tests that a trigger's Revealed card surfaces in the My turn section's first numbered entry.
func TestFormatBestTurn_TriggersFromLastTurnRevealedLine(t *testing.T) {
	summary := TurnSummary{
		TriggersFromLastTurn: []TriggerContribution{
			{Card: cards.SigilOfTheArknightBlue{}, Revealed: cards.MauvrionSkiesRed{}},
		},
	}
	out := FormatBestTurn(summary, nil, nil)
	want := "1. Sigil of the Arknight [B]: drew Mauvrion Skies [R] into hand"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// Tests that a TriggerContribution.Text override renders verbatim, skipping the synthesised
// "drew X into hand" suffix.
func TestFormatBestTurn_TriggersFromLastTurnHandlerAuthoredText(t *testing.T) {
	summary := TurnSummary{
		TriggersFromLastTurn: []TriggerContribution{
			{
				Card:     cards.SigilOfTheArknightBlue{},
				Revealed: cards.MauvrionSkiesRed{},
				Text:     "Sigil of the Arknight [B] revealed Mauvrion Skies [R] but didn't draw it",
			},
		},
	}
	out := FormatBestTurn(summary, nil, nil)
	want := "1. Sigil of the Arknight [B] revealed Mauvrion Skies [R] but didn't draw it"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
	// The synthesised "drew X into hand" must NOT appear — Text takes precedence.
	if strings.Contains(out, "drew Mauvrion Skies [R] into hand") {
		t.Errorf("synthesised suffix leaked despite Text override; got:\n%s", out)
	}
}

// Tests that zero-effect carryover triggers (no damage, no reveal) render no line at all.
func TestFormatBestTurn_TriggersFromLastTurnZeroEffectDropped(t *testing.T) {
	summary := TurnSummary{
		TriggersFromLastTurn: []TriggerContribution{
			{Card: cards.SigilOfTheArknightBlue{}},
		},
	}
	out := FormatBestTurn(summary, nil, nil)
	if out != "" {
		t.Errorf("zero-effect trigger with no other content should render empty; got:\n%s", out)
	}
}

// Tests AppendGroupedChainEntries clusters pre-/post-triggers under their Source-matched
// chain parent.
func TestAppendGroupedChainEntries_ClustersTriggersUnderTheirParent(t *testing.T) {
	log := []turnlogger.LogEntry{
		// Card A's pre-trigger fires from a hero/aura before A's chain entry resolves.
		{Text: "Viserai created a runechant", Source: "Card A", Kind: turnlogger.LogEntryPreTrigger, N: 1},
		// Card A resolves.
		{Text: "Card A: ATTACK", N: 5},
		// Card A's OnHit fires after the hit.
		{Text: "Aura created 3 runechants on hit", Source: "Card A", Kind: turnlogger.LogEntryPostTrigger, N: 3},
		// Card B's pre-trigger queues for B.
		{Text: "Viserai created a runechant", Source: "Card B", Kind: turnlogger.LogEntryPreTrigger, N: 1},
		// Card B resolves; its post-trigger follows.
		{Text: "Card B: PLAY", N: 0},
		{Text: "Aura created 2 runechants on hit", Source: "Card B", Kind: turnlogger.LogEntryPostTrigger, N: 2},
	}
	got := AppendGroupedChainEntries(nil, log)
	want := []string{
		"Card A: ATTACK (+5)",
		"  Viserai created a runechant (+1)",
		"  Aura created 3 runechants on hit (+3)",
		"Card B: PLAY",
		"  Viserai created a runechant (+1)",
		"  Aura created 2 runechants on hit (+2)",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("grouped output mismatch\n got: %#v\nwant: %#v", got, want)
	}
}

// Tests duplicate-name disambiguation: a pre-trigger between two same-name chain entries
// attaches to the SECOND (its real parent), not the first.
func TestAppendGroupedChainEntries_PreTriggerAttachesToNextSameNameParent(t *testing.T) {
	log := []turnlogger.LogEntry{
		// First Mauvrion plays (no triggers).
		{Text: "Mauvrion Skies [R]: PLAY", N: 0},
		// Second Mauvrion's hero pre-trigger fires (Viserai now sees a non-attack
		// action played) — Source matches the first chain entry's name too, but it
		// belongs to the second.
		{Text: "Viserai created a runechant", Source: "Mauvrion Skies [R]", Kind: turnlogger.LogEntryPreTrigger, N: 1},
		// Second Mauvrion's chain entry.
		{Text: "Mauvrion Skies [R]: PLAY", N: 0},
	}
	got := AppendGroupedChainEntries(nil, log)
	want := []string{
		"Mauvrion Skies [R]: PLAY",
		"Mauvrion Skies [R]: PLAY",
		"  Viserai created a runechant (+1)",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("pre-trigger attached to wrong parent\n got: %#v\nwant: %#v", got, want)
	}
}

// Tests the defensive fallback: a trigger whose Source matches no chain entry surfaces as
// a top-level line rather than being dropped.
func TestAppendGroupedChainEntries_OrphanTriggerSurfacesAtTopLevel(t *testing.T) {
	log := []turnlogger.LogEntry{
		{Text: "Card A: ATTACK", N: 5},
		{Text: "Aura created 2 runechants on hit", Source: "Card Z", Kind: turnlogger.LogEntryPostTrigger, N: 2},
	}
	got := AppendGroupedChainEntries(nil, log)
	want := []string{
		"Card A: ATTACK (+5)",
		"Aura created 2 runechants on hit (+2) (from Card Z)",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("orphan trigger should render as top-level line\n got: %#v\nwant: %#v",
			got, want)
	}
}

// Tests that plain BLOCK lines render the defender's effective Defense as a "(+N)" suffix.
func TestFormatBestTurn_BlockLineCarriesDefenseValue(t *testing.T) {
	summary := TurnSummary{
		BestLine: []CardAssignment{
			{Card: testutils.RedAttack{}, Role: Defend},
		},
	}
	out := FormatBestTurn(summary, nil, nil)
	want := "cardtest.RedAttack [R]: BLOCK (+1)"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// TestFormatBestTurn_StartOfTurnAurasLine pins the Start of turn section's "Auras: ..."
// entry. Names sort alphabetically for determinism, and duplicates are preserved (two copies
// of the same aura render twice).
func TestFormatBestTurn_StartOfTurnAurasLine(t *testing.T) {
	summary := TurnSummary{
		StartOfTurnAuras: []card.Card{
			cards.MaleficIncantationRed{},
			cards.MaleficIncantationRed{},
			cards.SigilOfTheArknightBlue{},
		},
	}
	out := FormatBestTurn(summary, nil, nil)
	want := "Auras: Malefic Incantation [R], Malefic Incantation [R], Sigil of the Arknight [B]"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// TestFormatBestTurn_StartOfTurnAurasSuppressedWhenEmpty pins the omission of the Auras
// entry when no auras were in play and no starting runechants carry in — the empty state
// shouldn't render a dangling "Auras:" label.
func TestFormatBestTurn_StartOfTurnAurasSuppressedWhenEmpty(t *testing.T) {
	summary := TurnSummary{BestLine: []CardAssignment{{Card: testutils.RedAttack{}, Role: Attack}}}
	out := FormatBestTurn(summary, nil, nil)
	if strings.Contains(out, "Auras: ") {
		t.Errorf("unexpected Auras line in output:\n%s", out)
	}
}

// TestFormatBestTurn_StartOfTurnAurasWithRunechants folds a non-zero starting Runechant
// carry into the "Auras:" entry as the trailing item — a Runeblade hero carrying tokens
// from the previous turn sees them alongside any auras as one combined readout.
func TestFormatBestTurn_StartOfTurnAurasWithRunechants(t *testing.T) {
	summary := TurnSummary{
		StartOfTurnAuras: []card.Card{cards.MaleficIncantationRed{}},
	}
	out := FormatBestTurn(summary, []*aura.Aura{cards.NewRunechant(3)}, nil)
	want := "Auras: Malefic Incantation [R], 3 Runechants"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// TestFormatBestTurn_StartOfTurnRunechantsOnly folds a non-zero starting Runechant carry
// into the Auras entry even when no auras are in play, using singular "Runechant" when the
// count is 1.
func TestFormatBestTurn_StartOfTurnRunechantsOnly(t *testing.T) {
	out := FormatBestTurn(TurnSummary{}, []*aura.Aura{cards.NewRunechant(1)}, nil)
	want := "Auras: 1 Runechant"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
	// Plural noun when count > 1.
	out2 := FormatBestTurn(TurnSummary{}, []*aura.Aura{cards.NewRunechant(2)}, nil)
	if !strings.Contains(out2, "2 Runechants") {
		t.Errorf("want plural 'Runechants' at count 2, got:\n%s", out2)
	}
}

// TestFormatBestTurn_StartOfTurnGoldItems surfaces a Gold token carryover as an
// "Items: N Gold" line in the Start of turn section.
func TestFormatBestTurn_StartOfTurnGoldItems(t *testing.T) {
	out := FormatBestTurn(TurnSummary{}, nil, []*token.Item{cards.NewGold(2)})
	want := "Items: 2 Gold"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// TestFormatBestTurn_EndOfTurnGoldItems surfaces Gold tokens surviving into next turn
// as an "Items: N Gold" line in the End of turn section.
func TestFormatBestTurn_EndOfTurnGoldItems(t *testing.T) {
	summary := TurnSummary{
		State: EngineWithItems([]*token.Item{cards.NewGold(1)}),
	}
	out := FormatBestTurn(summary, nil, nil)
	want := "Items: 1 Gold"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// TestFormatBestTurn_EndOfTurnHandLine pins the End of turn "Hand: ..." entry — every card
// in State.Hand surfaces as part of one comma-separated line, regardless of whether it
// started the turn in hand or got drawn / tutored mid-chain.
func TestFormatBestTurn_EndOfTurnHandLine(t *testing.T) {
	summary := TurnSummary{
		State: EngineWithHand([]card.Card{testutils.RedAttack{}}),
	}
	out := FormatBestTurn(summary, nil, nil)
	want := "Hand: cardtest.RedAttack [R]"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// TestFormatBestTurn_EndOfTurnAurasWithRunechants pins the End of turn "Auras: ..." entry —
// surviving Auras + the live Runechant count render as one comma-separated line,
// mirroring the start-of-turn formatting.
func TestFormatBestTurn_EndOfTurnAurasWithRunechants(t *testing.T) {
	gs := gameengine.GameStateBuilder().Build()
	gs.CreateAura(aura.NewCard(
		cards.MaleficIncantationRed{},
		triggertype.StartOfTurn,
		nil, 1, false,
	))
	gs.CreateAura(cards.NewRunechant(2))
	summary := TurnSummary{State: gs}
	out := FormatBestTurn(summary, nil, nil)
	want := "Auras: Malefic Incantation [R], 2 Runechants"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// TestFormatBestTurn_EndOfTurnArsenalStayedDirect drives endingArsenalLine's (stayed) branch
// directly via a synthesised BestLine. Pairs with TestFormatBestTurn_EndOfTurnArsenalNew's
// (new) branch coverage; the round-trip integration tests only exercise (new).
func TestFormatBestTurn_EndOfTurnArsenalStayedDirect(t *testing.T) {
	summary := TurnSummary{
		BestLine: []CardAssignment{
			{Card: cards.ToughenUpBlue{}, Role: Arsenal, FromArsenal: true},
		},
	}
	out := FormatBestTurn(summary, nil, nil)
	want := "Arsenal: Toughen Up [B] (stayed)"
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
		if got := FormatContribution(c.in); got != c.want {
			t.Errorf("FormatContribution(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}
