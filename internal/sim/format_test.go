package sim_test

import (
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"reflect"
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
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

// TestFormatBestTurn_AttackAndPitch verifies the basic numbered-list shape: pitches and the
// attack chain both fall under the "My turn:" section header. Hand: 2 Red Attacks + 2 Blues.
// One Blue pitches for 3 resource, funding the 3-cost chain (Blue + Red + Red, all cost 1
// each, all go-again).
func TestFormatBestTurn_AttackAndPitch(t *testing.T) {
	h := []Card{testutils.BlueAttack{}, testutils.BlueAttack{}, testutils.RedAttack{}, testutils.RedAttack{}}
	got := Best(StubHero, nil, h, 0, nil, 0, nil)
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
	h := []Card{cards.MauvrionSkiesRed{}, cards.ShrillOfSkullformRed{}, cards.MaleficIncantationBlue{}}
	got := Best(heroes.Viserai{}, nil, h, 0, nil, 0, nil)
	out := FormatBestTurn(got, 0)
	if !strings.Contains(out, "Mauvrion Skies [R]: PLAY") {
		t.Errorf("want Mauvrion (non-attack action) labelled PLAY, got:\n%s", out)
	}
	if !strings.Contains(out, "Shrill of Skullform [R]: ATTACK") {
		t.Errorf("want Shrill (attack action) labelled ATTACK, got:\n%s", out)
	}
}

// TestFormatBestTurn_LogAttributesEachTriggerSeparately pins phase-2 logging: each event in
// the chain (card Play, hero trigger, mid-chain aura trigger, ephemeral attack trigger) gets
// its own attributed line, grouped under the chain entry of the card that triggered it.
// Hand: Nimblism (Generic non-attack action, sets NonAttackActionPlayed) + Mauvrion
// (Runeblade non-attack action, registers an ephemeral "if hits, +3" trigger) + Consuming
// Volition (Runeblade attack action). Prior aura: a Malefic Incantation TriggerAttackAction.
// The chain log should hit:
//   - Mauvrion: PLAY (+0)             — non-attack action, no own damage
//   - Consuming Volition: ATTACK (+N) — printed power, with the three triggers it fires
//     attached underneath as indented children. Each trigger handler authors its own log
//     line (Viserai / Malefic / Mauvrion all create runechants). The "(from Consuming
//     Volition)" suffix is dropped because the visual grouping makes the source obvious.
func TestFormatBestTurn_LogAttributesEachTriggerSeparately(t *testing.T) {
	h := []Card{cards.NimblismRed{}, cards.MauvrionSkiesRed{}, cards.ConsumingVolitionRed{}}
	// Use the real Malefic Incantation card's Play to register the prior trigger so the
	// handler matches production exactly (logs via AddPreTriggerLogEntry, sources from
	// state.TriggeringCard).
	var bootstrap TurnState
	cards.MaleficIncantationRed{}.Play(&bootstrap, &CardState{Card: cards.MaleficIncantationRed{}})
	prior := bootstrap.AuraTriggers
	got := BestWithTriggers(heroes.Viserai{}, nil, h, 0, nil, 0, nil, prior)
	out := FormatBestTurn(got, 0)
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

// TestFormatBestTurn_LogSuppressesZeroTriggers pins that handlers returning 0 don't add a
// "(+0)" line for the trigger — log entries for triggers gate on positive contribution to
// avoid noise. Card-Play lines render unconditionally because the card resolved (the line
// proves the chain step happened); trigger lines only render when the trigger actually
// credited damage.
func TestFormatBestTurn_LogSuppressesZeroTriggers(t *testing.T) {
	// Hand: a single Red attack with no go-again. Viserai's OnCardPlayed contributes nothing
	// (the gate needs another non-attack action played first), no priors, no ephemerals — so
	// the chain log should be exactly one card-Play line, no trigger spam.
	h := []Card{testutils.RedAttack{}}
	got := Best(heroes.Viserai{}, nil, h, 0, nil, 0, nil)
	if strings.Contains(FormatBestTurn(got, 0), "Viserai created") {
		t.Errorf("hero trigger line shouldn't render when Viserai contributed 0; got:\n%s",
			FormatBestTurn(got, 0))
	}
}

// TestFormatBestTurn_MoonWishTutorAndPlayLogsAsPostTrigger: the go-again branch tutors
// Sun Kiss and immediately plays it. Moon Wish's chain step shows only its printed
// attack; the tutor narration line renders as a post-trigger child grouped beneath Moon
// Wish, Sun Kiss authors its own "PLAY" chain entry, and the heal lands as a "Gained 3
// health (+3)" sub-line under Sun Kiss. The third hand card is the alt-cost target so
// Flying High [R] stays in the chain to grant go-again.
func TestFormatBestTurn_MoonWishTutorAndPlayLogsAsPostTrigger(t *testing.T) {
	h := []Card{cards.FlyingHighRed{}, cards.MoonWishYellow{}, testutils.BlueAttack{}}
	deck := []Card{cards.SunKissRed{}}
	got := Best(StubHero, nil, h, 0, deck, 0, nil)
	out := FormatBestTurn(got, 0)
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

// TestFormatBestTurn_MoonWishTutorOnlyLogsAsPostTrigger: the no-go-again branch tutors
// Sun Kiss but doesn't play it (it lands in hand for next turn). The post-trigger line
// "Moon Wish [Y] tutored Sun Kiss [R]" renders without a (+N) since no damage credits.
// The Blue attack is held to satisfy Moon Wish's alt cost.
func TestFormatBestTurn_MoonWishTutorOnlyLogsAsPostTrigger(t *testing.T) {
	h := []Card{cards.MoonWishYellow{}, testutils.BlueAttack{}}
	deck := []Card{cards.SunKissRed{}}
	got := Best(StubHero, nil, h, 0, deck, 0, nil)
	out := FormatBestTurn(got, 0)
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

// TestFormatBestTurn_ArsenalInPlayedAsDR checks the combined "arsenal-in played from the slot"
// + "defense reaction prevented" rendering. Hand: one Malefic Blue (pitch 3). Arsenal-in:
// Toughen Up Blue (DR cost 2). Malefic pitches to fund the DR, Toughen Up blocks 4 of 4 incoming.
// Display puts the pitch and DR lines under the "Opponent's turn:" section; the role label
// reads "DEFENSE REACTION from arsenal" since Toughen Up came out of the arsenal slot.
func TestFormatBestTurn_ArsenalInPlayedAsDR(t *testing.T) {
	h := []Card{cards.MaleficIncantationBlue{}}
	got := Best(StubHero, nil, h, 4, nil, 0, cards.ToughenUpBlue{})
	out := FormatBestTurn(got, 0)
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

// TestFormatBestTurn_DefenseReactionLinesAndRiders pins the per-DR rendering: the chain step's
// "(+N)" folds in BonusDefense (the +1{d} bonus is rolled in just like BonusAttack feeds the
// attack chain step), and separable riders like the arcane ping each land as their own indented
// sub-line under the parent. Sigil of Suffering Red against incoming 4 has the Sigil block 4
// (printed 3 + 1 from the arcane-conditional bonus) and deal 1 arcane on a sub-line. Dodge has
// no riders or bonuses, so it renders as a single chain step with "(+2)" — the printed Defense.
func TestFormatBestTurn_DefenseReactionLinesAndRiders(t *testing.T) {
	cases := []struct {
		name     string
		hand     []Card
		incoming int
		wants    []string
	}{
		{
			name:     "Sigil of Suffering folds bonus into chain step + arcane sub-line",
			hand:     []Card{cards.SigilOfSufferingRed{}, testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}},
			incoming: 4,
			wants: []string{
				"Sigil of Suffering [R]: DEFENSE REACTION (+4)",
				"Dealt 1 arcane damage (+1)",
			},
		},
		{
			name:     "Dodge has no riders, single chain line",
			hand:     []Card{cards.DodgeBlue{}},
			incoming: 2,
			wants:    []string{"Dodge [B]: DEFENSE REACTION (+2)"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Best(StubHero, nil, tc.hand, tc.incoming, nil, 0, nil)
			out := FormatBestTurn(got, 0)
			for _, w := range tc.wants {
				if !strings.Contains(out, w) {
					t.Errorf("want %q in:\n%s", w, out)
				}
			}
		})
	}
}

// TestFormatBestTurn_ArsenalInPlayedOnChain checks the role-label tag for an arsenal-in
// card played as part of the my-turn chain. Hand: one BlueAttack (pitch 3, cost 1).
// Arsenal-in: RedAttack (cost 1, attack 3). The solver pitches the Blue to pay the Red's
// cost and attacks from arsenal for 3; the chain line reads "cardtest.RedAttack [R]: ATTACK
// from arsenal" — tag on the role, not on the card name.
func TestFormatBestTurn_ArsenalInPlayedOnChain(t *testing.T) {
	h := []Card{testutils.BlueAttack{}}
	got := Best(StubHero, nil, h, 0, nil, 0, testutils.RedAttack{})
	out := FormatBestTurn(got, 0)
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

// TestFormatBestTurn_WeaponSwingInChain makes sure a swung weapon shows up in the chain with
// a WEAPON ATTACK label, sourced from the dispatcher's Log entry. The State.Log assertion
// pins the dispatcher → log → format pipeline for weapons; FormatBestTurn reads weapon
// swings from State.Log rather than SwungWeapons.
func TestFormatBestTurn_WeaponSwingInChain(t *testing.T) {
	h := []Card{testutils.RedAttack{}}
	weapons := []Weapon{weapons.ReapingBlade{}}
	got := Best(StubHero, weapons, h, 0, nil, 0, nil)
	out := FormatBestTurn(got, 0)
	if !strings.Contains(out, "Reaping Blade: WEAPON ATTACK") {
		t.Errorf("want the weapon in the chain, got:\n%s", out)
	}
	var sawWeaponLog bool
	for _, e := range got.State.Log {
		if strings.Contains(e.Text, "Reaping Blade: WEAPON ATTACK") {
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
// Arsenal, so End of turn shows "Arsenal: Toughen Up [B] (new)".
func TestFormatBestTurn_EndOfTurnArsenalNew(t *testing.T) {
	h := []Card{cards.ToughenUpBlue{}}
	got := Best(StubHero, nil, h, 4, nil, 0, nil)
	out := FormatBestTurn(got, 0)
	if !strings.Contains(out, "Arsenal: Toughen Up [B] (new)") {
		t.Errorf("want an end-of-turn arsenal entry tagged '(new)', got:\n%s", out)
	}
}

// TestFormatBestTurn_EndOfTurnArsenalStayed tags the carrying-over arsenal card with
// "(stayed)" rather than "(new)" — useful for the reader to see the slot wasn't swapped.
func TestFormatBestTurn_EndOfTurnArsenalStayed(t *testing.T) {
	// Hand with no attacks / no pitches to pay for the arsenal DR at incoming=0 (defense is
	// wasted anyway). Arsenal-in Toughen Up sits.
	h := []Card{cards.ToughenUpBlue{}}
	got := Best(StubHero, nil, h, 0, nil, 0, cards.ToughenUpBlue{})
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
// at the top of the "My turn:" section as numbered entries — the reveal / damage fires at
// the top of the action phase before the chain runs, but it's an action, not pre-existing
// state, so it belongs to the action-phase numbering rather than the unnumbered Start of
// turn block.
func TestFormatBestTurn_TriggersFromLastTurnLine(t *testing.T) {
	summary := TurnSummary{
		TriggersFromLastTurn: []TriggerContribution{
			{Card: testutils.RedAttack{}, Damage: 3},
		},
	}
	out := FormatBestTurn(summary, 0)
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

// TestFormatBestTurn_StartOfTurnHandReadsDealtHand: the printout's "Start of turn → Hand:"
// line must reflect the cards dealt at end of last turn — NOT the augmented hand a Sigil
// reveal produces by appending its drawn card. The deck loop snapshots TurnSummary.DealtHand
// before processTriggersAtStartOfTurn modifies the working hand, and the formatter reads
// only that snapshot. The Sigil-revealed card appears under MyTurn (where it actually
// resolves), never in the start-of-turn hand line.
func TestFormatBestTurn_StartOfTurnHandReadsDealtHand(t *testing.T) {
	summary := TurnSummary{
		DealtHand: []Card{testutils.RedAttack{}},
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
	out := FormatBestTurn(summary, 0)
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

// TestFormatBestTurn_TriggersFromLastTurnRevealedLine surfaces the card a trigger handler
// revealed into the hand. Sigil of the Arknight fires at start of action phase with
// Damage=0 but reveals the deck top; the My turn section's first numbered entry names the
// card it drew.
func TestFormatBestTurn_TriggersFromLastTurnRevealedLine(t *testing.T) {
	summary := TurnSummary{
		TriggersFromLastTurn: []TriggerContribution{
			{Card: cards.SigilOfTheArknightBlue{}, Revealed: cards.MauvrionSkiesRed{}},
		},
	}
	out := FormatBestTurn(summary, 0)
	want := "1. Sigil of the Arknight [B]: drew Mauvrion Skies [R] into hand"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// TestFormatBestTurn_TriggersFromLastTurnHandlerAuthoredText: when the handler authors a
// custom log line via state.AddPostTriggerLogEntry, the format layer renders Text verbatim
// and skips the synthesised "Aura Name: drew X into hand" / "(+N)" suffix. Cards keep
// full ownership of their printout wording.
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
	out := FormatBestTurn(summary, 0)
	want := "1. Sigil of the Arknight [B] revealed Mauvrion Skies [R] but didn't draw it"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
	// The synthesised "drew X into hand" must NOT appear — Text takes precedence.
	if strings.Contains(out, "drew Mauvrion Skies [R] into hand") {
		t.Errorf("synthesised suffix leaked despite Text override; got:\n%s", out)
	}
}

// TestFormatBestTurn_TriggersFromLastTurnZeroEffectDropped suppresses lines for carryover
// triggers that did nothing visible this turn (zero damage, no reveal). Output has no
// numbered entries at all — the My turn section is empty so its header elides too.
func TestFormatBestTurn_TriggersFromLastTurnZeroEffectDropped(t *testing.T) {
	summary := TurnSummary{
		TriggersFromLastTurn: []TriggerContribution{
			{Card: cards.SigilOfTheArknightBlue{}},
		},
	}
	out := FormatBestTurn(summary, 0)
	if out != "" {
		t.Errorf("zero-effect trigger with no other content should render empty; got:\n%s", out)
	}
}

// TestAppendGroupedChainEntries_ClustersTriggersUnderTheirParent drives the grouping
// helper directly with a synthesised LogEntry slice — covers the three placement cases
// (pre-trigger before its parent, post-trigger after, trigger from card B interleaved
// before card B's chain entry) without needing a full Best invocation. Trigger Texts are
// freeform card-authored phrases ("Viserai created a runechant"); the grouping helper
// matches each trigger's Source field against the chain entry's "<Name>:" prefix.
func TestAppendGroupedChainEntries_ClustersTriggersUnderTheirParent(t *testing.T) {
	log := []LogEntry{
		// Card A's pre-trigger fires from a hero/aura before A's chain entry resolves.
		{Text: "Viserai created a runechant", Source: "Card A", Kind: LogEntryPreTrigger, N: 1},
		// Card A resolves.
		{Text: "Card A: ATTACK", N: 5},
		// Card A's ephemeral attack trigger fires after the hit.
		{Text: "Aura created 3 runechants on hit", Source: "Card A", Kind: LogEntryPostTrigger, N: 3},
		// Card B's pre-trigger queues for B.
		{Text: "Viserai created a runechant", Source: "Card B", Kind: LogEntryPreTrigger, N: 1},
		// Card B resolves; its post-trigger follows.
		{Text: "Card B: PLAY", N: 0},
		{Text: "Aura created 2 runechants on hit", Source: "Card B", Kind: LogEntryPostTrigger, N: 2},
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

// TestAppendGroupedChainEntries_PreTriggerAttachesToNextSameNameParent guards the
// duplicate-name disambiguation: when two chain entries share a display name (e.g. two
// copies of Mauvrion Skies in one chain), a pre-trigger between them belongs to the
// SECOND parent — it fires before the second card's Play. Source-name match alone is
// ambiguous in this case; Kind==LogEntryPreTrigger is what tells the grouping algorithm
// to skip the first parent's post-trigger lookforward and let the entry fall through to
// the next matching chain step.
func TestAppendGroupedChainEntries_PreTriggerAttachesToNextSameNameParent(t *testing.T) {
	log := []LogEntry{
		// First Mauvrion plays (no triggers).
		{Text: "Mauvrion Skies [R]: PLAY", N: 0},
		// Second Mauvrion's hero pre-trigger fires (Viserai now sees a non-attack
		// action played) — Source matches the first chain entry's name too, but it
		// belongs to the second.
		{Text: "Viserai created a runechant", Source: "Mauvrion Skies [R]", Kind: LogEntryPreTrigger, N: 1},
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

// TestAppendGroupedChainEntries_OrphanTriggerSurfacesAtTopLevel guards the defensive
// fallback: a trigger whose Source matches no chain entry shouldn't be silently dropped.
// Currently impossible in practice (playSequenceWithMeta emits triggers immediately around
// their parent) but the fallback keeps the data visible if that invariant ever loosens.
func TestAppendGroupedChainEntries_OrphanTriggerSurfacesAtTopLevel(t *testing.T) {
	log := []LogEntry{
		{Text: "Card A: ATTACK", N: 5},
		{Text: "Aura created 2 runechants on hit", Source: "Card Z", Kind: LogEntryPostTrigger, N: 2},
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

// TestFormatBestTurn_BlockLineCarriesDefenseValue pins the "(+N)" suffix on plain BLOCK
// lines. Each block line shows the defender's effective Defense so the reader can sum the
// wall against the incoming attack without re-checking each card. testutils.RedAttack has
// printed Defense=1; a synthesised BestLine drives the renderer directly.
func TestFormatBestTurn_BlockLineCarriesDefenseValue(t *testing.T) {
	summary := TurnSummary{
		BestLine: []CardAssignment{
			{Card: testutils.RedAttack{}, Role: Defend},
		},
	}
	out := FormatBestTurn(summary, 0)
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
		StartOfTurnAuras: []Card{
			cards.MaleficIncantationRed{},
			cards.MaleficIncantationRed{},
			cards.SigilOfTheArknightBlue{},
		},
	}
	out := FormatBestTurn(summary, 0)
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
		StartOfTurnAuras: []Card{cards.MaleficIncantationRed{}},
	}
	out := FormatBestTurn(summary, 3)
	want := "Auras: Malefic Incantation [R], 3 Runechants"
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
		State: CarryState{Hand: []Card{testutils.RedAttack{}}},
	}
	out := FormatBestTurn(summary, 0)
	want := "Hand: cardtest.RedAttack [R]"
	if !strings.Contains(out, want) {
		t.Errorf("missing %q in:\n%s", want, out)
	}
}

// TestFormatBestTurn_EndOfTurnAurasWithRunechants pins the End of turn "Auras: ..." entry —
// surviving AuraTriggers + the live Runechant count render as one comma-separated line,
// mirroring the start-of-turn formatting.
func TestFormatBestTurn_EndOfTurnAurasWithRunechants(t *testing.T) {
	summary := TurnSummary{
		State: CarryState{
			AuraTriggers: []AuraTrigger{
				{Self: cards.MaleficIncantationRed{}},
			},
			Runechants: 2,
		},
	}
	out := FormatBestTurn(summary, 0)
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
	out := FormatBestTurn(summary, 0)
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
