package lint

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// TestTurntests_NewTestsUseEvalNTurnsForTesting rejects new turntests/ test files that
// don't drive a turn through the public sim.EvalOneTurnForTesting /
// sim.EvalTwoTurnsForTesting entry points. Per docs/dev-standards.md "Test layout",
// turntests/ is for public-entry-point tests only — anything that drives the engine
// directly (ResolveAttackStep, ge.Discard, ge.FireTriggers, …) belongs in a same-package
// unit test under the card's home directory.
//
// grandfatheredTurntestFiles is the allow-list of existing offenders; remove an entry to
// lock in its migration to either a same-package unit test or the public Eval entry
// points.
func TestTurntests_NewTestsUseEvalNTurnsForTesting(t *testing.T) {
	root := RepoRoot(t)
	turntestsDir := filepath.Join(root, "turntests")

	allowed := make(map[string]bool, len(grandfatheredTurntestFiles))
	for _, p := range grandfatheredTurntestFiles {
		allowed[p] = true
	}

	var offenders []string
	err := filepath.WalkDir(turntestsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		src := string(body)
		// A passing file must call EvalOneTurnForTesting or EvalTwoTurnsForTesting (its
		// purpose is to exercise the public turn-eval API). Files that drive engine
		// internals directly (ResolveAttackStep, ge.Discard, ge.FireTriggers, …) are
		// offenders unless grandfathered.
		if strings.Contains(src, "EvalOneTurnForTesting") || strings.Contains(src, "EvalTwoTurnsForTesting") {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		if !allowed[rel] {
			offenders = append(offenders, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}

	sort.Strings(offenders)
	for _, o := range offenders {
		t.Errorf("%s: doesn't drive a turn through sim.EvalOneTurnForTesting / "+
			"sim.EvalTwoTurnsForTesting. turntests/ is for public-entry-point tests only "+
			"(docs/dev-standards.md \"Test layout\"). Move card-specific unit tests to "+
			"internal/card/cards/<card>_test.go (same package).",
			o)
	}
}

// grandfatheredTurntestFiles is the allow-list of existing turntests/ files that don't
// drive a turn through sim.EvalOneTurnForTesting / sim.EvalTwoTurnsForTesting — they call
// ge.ResolveAttackStep or other engine internals directly. They need migrating to
// same-package unit tests (under internal/card/cards/<card>_test.go) or rewriting against
// the public Eval entry points. To migrate, remove the entry — the lint then rejects
// reintroduction. See TODO.md "Tech debt" for the migration plan.
var grandfatheredTurntestFiles = []string{
	"turntests/aether_slash_test.go",
	"turntests/arcane_polarity_test.go",
	"turntests/battlefront_bastion_test.go",
	"turntests/bloodspill_invocation_test.go",
	"turntests/brothers_in_arms_test.go",
	"turntests/cadaverous_contraband_test.go",
	"turntests/captains_call_test.go",
	"turntests/card_high_striker_test.go",
	"turntests/card_moon_wish_test.go",
	"turntests/card_performance_bonus_test.go",
	"turntests/card_plunder_run_test.go",
	"turntests/card_relentless_pursuit_test.go",
	"turntests/card_starting_stake_test.go",
	"turntests/card_tremor_of_iarathael_test.go",
	"turntests/card_warmongers_recital_test.go",
	"turntests/card_weeping_battleground_test.go",
	"turntests/card_yinti_yanti_test.go",
	"turntests/come_to_fight_test.go",
	"turntests/consuming_volition_test.go",
	"turntests/deathly_duet_test.go",
	"turntests/destructive_deliberation_test.go",
	"turntests/drawn_to_the_dark_dimension_test.go",
	"turntests/drowning_dire_test.go",
	"turntests/exposed_test.go",
	"turntests/fate_foreseen_test.go",
	"turntests/fervent_forerunner_test.go",
	"turntests/flying_high_test.go",
	"turntests/force_sight_test.go",
	"turntests/from_arsenal_go_again_test.go",
	"turntests/from_arsenal_next_attack_bonus_test.go",
	"turntests/jack_be_nimble_test.go",
	"turntests/jack_be_quick_test.go",
	"turntests/lifegain_variants_test.go",
	"turntests/looking_for_a_scrap_test.go",
	"turntests/lower_health_wanter_test.go",
	"turntests/meat_and_greet_test.go",
	"turntests/memorial_ground_test.go",
	"turntests/minnowism_test.go",
	"turntests/nimble_strike_test.go",
	"turntests/nimblism_test.go",
	"turntests/oath_of_the_arknight_test.go",
	"turntests/overload_test.go",
	"turntests/peace_of_mind_test.go",
	"turntests/prime_the_crowd_test.go",
	"turntests/public_bounty_test.go",
	"turntests/punch_above_your_weight_test.go",
	"turntests/pursue_to_the_pits_of_despair_test.go",
	"turntests/ravenous_rabble_test.go",
	"turntests/reduce_to_runechant_test.go",
	"turntests/reek_of_corruption_test.go",
	"turntests/regurgitating_slog_test.go",
	"turntests/right_behind_you_test.go",
	"turntests/rise_above_test.go",
	"turntests/runerager_swarm_test.go",
	"turntests/runic_reaping_test.go",
	"turntests/scout_the_periphery_test.go",
	"turntests/seek_horizon_test.go",
	"turntests/shrill_of_skullform_test.go",
	"turntests/sigil_of_protection_test.go",
	"turntests/sky_fire_lanterns_test.go",
	"turntests/sloggism_test.go",
	"turntests/smashing_good_time_test.go",
	"turntests/snatch_test.go",
	"turntests/spring_load_test.go",
	"turntests/strategic_planning_test.go",
	"turntests/sun_kiss_test.go",
	"turntests/sutcliffes_research_notes_test.go",
	"turntests/test_of_strength_test.go",
	"turntests/trot_along_test.go",
	"turntests/vantage_point_test.go",
	"turntests/vigor_rush_test.go",
	"turntests/water_the_seeds_test.go",
	"turntests/zealous_belting_test.go",
}
