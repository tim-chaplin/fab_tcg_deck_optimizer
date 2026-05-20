package lint

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// TestTurntests_NewTestsUseEvalNTurnsForTesting rejects new turntests/ test files that
// drive the chain via ge.ResolveChainStep(...) directly. Per docs/dev-standards.md "Test
// layout", turntests/ is for public-entry-point tests only: sim.EvalOneTurnForTesting and
// sim.EvalTwoTurnsForTesting.
//
// The grandfatheredResolveChainStepFiles list freezes the v2-migration backlog: every
// file currently using ResolveChainStep is allow-listed. Migrating one of those files to
// EvalOneTurnForTesting (or moving it to a same-package unit test) means removing it
// from the list — the lint will then reject re-introduction. New files using
// ResolveChainStep fail immediately.
func TestTurntests_NewTestsUseEvalNTurnsForTesting(t *testing.T) {
	root := RepoRoot(t)
	turntestsDir := filepath.Join(root, "turntests")

	allowed := make(map[string]bool, len(grandfatheredResolveChainStepFiles))
	for _, p := range grandfatheredResolveChainStepFiles {
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
		if !strings.Contains(string(body), "ResolveChainStep") {
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
		t.Errorf("%s: drives the chain via ResolveChainStep. turntests/ tests must use "+
			"sim.EvalOneTurnForTesting / sim.EvalTwoTurnsForTesting (the public entry points) "+
			"per docs/dev-standards.md \"Test layout\". Either rewrite the test against the "+
			"public API or move it to a same-package unit test under the card's home directory.",
			o)
	}
}

// grandfatheredResolveChainStepFiles freezes the v2-migration backlog. Files were dropped
// into turntests/ during the v2 reorganisation but still drive ResolveChainStep directly;
// they need a future home (same-package unit tests, or rewritten against the public Eval
// entry points). See TODO.md "Tech debt" for the migration plan. To migrate a file, remove
// its entry from this list — the lint will then reject reintroduction.
var grandfatheredResolveChainStepFiles = []string{
	"turntests/aether_slash_test.go",
	"turntests/arcane_polarity_test.go",
	"turntests/arcanic_spike_test.go",
	"turntests/blessing_of_occult_test.go",
	"turntests/bloodspill_invocation_test.go",
	"turntests/blow_for_a_blow_test.go",
	"turntests/bluster_buff_test.go",
	"turntests/brush_off_test.go",
	"turntests/cadaverous_contraband_test.go",
	"turntests/calming_breeze_test.go",
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
	"turntests/clearwater_elixir_test.go",
	"turntests/come_to_fight_test.go",
	"turntests/condemn_to_slaughter_test.go",
	"turntests/consuming_volition_test.go",
	"turntests/deathly_duet_test.go",
	"turntests/demolition_crew_test.go",
	"turntests/destructive_deliberation_test.go",
	"turntests/drowning_dire_test.go",
	"turntests/enchanting_melody_test.go",
	"turntests/exposed_test.go",
	"turntests/fate_foreseen_test.go",
	"turntests/fervent_forerunner_test.go",
	"turntests/flex_test.go",
	"turntests/flying_high_test.go",
	"turntests/force_sight_test.go",
	"turntests/freewheeling_renegades_test.go",
	"turntests/from_arsenal_go_again_test.go",
	"turntests/from_arsenal_next_attack_bonus_test.go",
	"turntests/hit_the_high_notes_test.go",
	"turntests/jack_be_nimble_test.go",
	"turntests/jack_be_quick_test.go",
	"turntests/life_for_a_life_test.go",
	"turntests/lifegain_variants_test.go",
	"turntests/looking_for_a_scrap_test.go",
	"turntests/lower_health_wanter_test.go",
	"turntests/malefic_incantation_test.go",
	"turntests/meat_and_greet_test.go",
	"turntests/memorial_ground_test.go",
	"turntests/minnowism_test.go",
	"turntests/nimble_strike_test.go",
	"turntests/nimblism_test.go",
	"turntests/oasis_respite_test.go",
	"turntests/oath_of_the_arknight_test.go",
	"turntests/outed_test.go",
	"turntests/overload_test.go",
	"turntests/peace_of_mind_test.go",
	"turntests/prime_the_crowd_test.go",
	"turntests/public_bounty_test.go",
	"turntests/punch_above_your_weight_test.go",
	"turntests/pursue_to_the_edge_of_oblivion_test.go",
	"turntests/pursue_to_the_pits_of_despair_test.go",
	"turntests/push_the_point_test.go",
	"turntests/ravenous_rabble_test.go",
	"turntests/reduce_to_runechant_test.go",
	"turntests/reek_of_corruption_test.go",
	"turntests/regain_composure_test.go",
	"turntests/regurgitating_slog_test.go",
	"turntests/restvine_elixir_test.go",
	"turntests/rise_above_test.go",
	"turntests/runechant_on_play_test.go",
	"turntests/runerager_swarm_test.go",
	"turntests/runic_fellingsong_test.go",
	"turntests/runic_reaping_test.go",
	"turntests/sapwood_elixir_test.go",
	"turntests/scout_the_periphery_test.go",
	"turntests/seek_horizon_test.go",
	"turntests/shrill_of_skullform_test.go",
	"turntests/sigil_of_fyendal_test.go",
	"turntests/sigil_of_protection_test.go",
	"turntests/sigil_of_silphidae_test.go",
	"turntests/sigil_of_suffering_test.go",
	"turntests/sigil_of_the_arknight_test.go",
	"turntests/sky_fire_lanterns_test.go",
	"turntests/sloggism_test.go",
	"turntests/smashing_good_time_test.go",
	"turntests/snatch_test.go",
	"turntests/spring_load_test.go",
	"turntests/starting_stake_test.go",
	"turntests/strategic_planning_test.go",
	"turntests/sun_kiss_test.go",
	"turntests/sutcliffes_research_notes_test.go",
	"turntests/test_of_strength_test.go",
	"turntests/trot_along_test.go",
	"turntests/vantage_point_test.go",
	"turntests/vigor_rush_test.go",
	"turntests/water_the_seeds_test.go",
	"turntests/whisper_of_the_oracle_test.go",
	"turntests/zealous_belting_test.go",
}
