# Streaming-logger refactor — handoff

You're continuing a refactor that was started in a prior session and paused at the foundation step. This doc has everything you need to finish it.

## Branch state

- Branch: `refactor/streaming-logger` (pushed to origin, no open PR yet).
- Based off: `main`.
- One commit: `668d5e1e` — switched `gs.logger` field from concrete `*turnlogger.TurnLogger` to the existing `card.Logger` interface; added `NoopLogger` (no-op value type) and `StreamLogger` (writes to `io.Writer`) as new implementations in `internal/gameengine/logger.go`. Behaviorally identical, all tests green.
- A parallel open PR (#429, `refactor/shared-carryover-helpers`) will rebase on top of this branch once it lands. **Don't try to merge or rebase it now** — that comes after this PR is merged.

## Design mandate (load-bearing — the user is firm on this)

> "Logging is purely ephemeral, no state needs to be communicated between functions; when cards call any Log functions just send it straight to stdout, that's it. If the logs are a little less pretty that's fine, the logging isn't an important enough feature to drive any complexity in the sim."

Consequences:
- **No log accumulation anywhere in the sim.** No `[]LogEntry` buffers, no `gs.LogEntries()`, no `BuildTurnLog`, no `deck.TurnLog` struct, no JSON serialization of formatted log strings, no `TriggersFromLastTurn` / `StartOfTurnAuras` / `DealtHand` fields on `TurnSummary`.
- **No `contribs` plumbing.** `processAurasAtStartOfTurn` should return just `damage` (or nothing) and rely on the engine's logger for any per-aura logging. Today it returns `contribs []deck.TriggerContribution` purely so `BuildTurnLog` can format trigger lines later — that's exactly the antipattern we're killing.
- **Output is FLAT chronological.** No tree-style grouping of triggers under their parent attack. `AmendLastChainStepN` (Attack-Reaction buff folding) is a no-op in `StreamLogger` — the buff surfaces as its own trigger line instead of mutating a prior line. This is intentional and accepted.
- **Single replay at print time.** During eval, every turn uses `NoopLogger`. Only at the final stats printout do we install a `StreamLogger` and replay the best turn's chain so its log lines stream to stdout.

## Concrete work remaining

In rough order. Each step should leave tests green so you can commit/checkpoint.

### 1. Drop the per-turn pre-fetches in `runOneShuffle`

In `internal/sim/deck_eval.go`, the inner loop builds these every turn purely to feed `BuildTurnLog` later:
- `startingAuras := concreteAuras(master.Auras())`
- `startingItems := concreteItems(master.Items())`
- `startOfTurnAuras := snapshotStartOfTurnAuras(startingAuras)`
- `dealtHand := append([]card.Card(nil), h...)`
- `play.TriggersFromLastTurn = contribs`
- `play.StartOfTurnAuras = startOfTurnAuras`
- `play.DealtHand = dealtHand`

All of this disappears. The per-turn allocation savings (~22M allocs/eval on the anneal bench) is the primary perf win.

### 2. Drop `contribs` from `processAurasAtStartOfTurn`

New signature: `processAurasAtStartOfTurn(master *gameengine.GameState, d *deck.Deck) (damage int)` (or even `int` directly). Aura handlers already get a `card.Logger` via the `FireStartOfTurn` callback path — they log directly to it. During eval that logger is `NoopLogger` and the call no-ops; during print-time replay it's a `StreamLogger`.

The export wrapper in `internal/sim/exports_test.go` (`ProcessAurasAtStartOfTurn`) will need rework — easiest is to delete the test-export entirely and rewrite the few `internal/sim/deck_aura_test.go` tests that depend on it to assert on master state instead of returned tuples.

The `gameengine.FireStartOfTurn` callback's `newEntries` parameter exists only to feed `contribs`. With `contribs` gone, the callback shape simplifies — `newEntries` can be deleted from the signature and from the `engine.go` body (the type-assertion to `Entries()` I just added in the foundation commit becomes dead code).

### 3. Drop `replayBestForTurnWithLog` from `runOneShuffle`

The eager replay-with-log on every new-best disappears. New-best fires ~15 times per eval; today each fires one extra full chain run with log materialization. Gone in the new design — log replay happens only once at print time.

### 4. Add `bestSnapshot` capture

```go
type bestSnapshot struct {
    Master   *gameengine.GameState  // post-tick state best() saw
    Deck     *deck.Deck             // post-tick deck best() saw
    Hand     []card.Card            // post-tick hand (held + drawn + reveals)
    Weapons  []weapon.Weapon
    Mp       Matchup
    BestLine []deck.CardAssignment
}
```

Capture site: where `recordTurnStats` returns true (the new-best branch in `runOneShuffle`). Use `master.CopyPersistentState()` + `d.Copy()` — both are cheap and rarely called (log N times per eval). Bestline gets `append([]CardAssignment(nil), play.BestLine...)`.

Stash the snapshot somewhere the print path can find it. Two options:
- (a) Add a `Stats.PrintBest func(io.Writer)` closure that captures the snapshot + Evaluator pointer. Set after eval finishes. `json:"-"` so it doesn't serialize. Loaded-from-JSON `Stats` has nil PrintBest, so the print path no-ops the chain section.
- (b) Add a typed `*bestSnapshot` field on `deck.Stats` with `json:"-"`, and write `sim.PrintBestTurn(stats, w)` separately.

(a) is cleaner because `deck.Stats` stays unaware of `sim` types.

**Parallel merge concern:** `mergeStatsInto` in `deck_eval.go` swaps `dst.Best = src.Best` when src wins. The snapshot needs to ride along — extend the merge to also propagate the snapshot (or its closure).

### 5. Add `sim.PrintBestTurn(snap *bestSnapshot, w io.Writer)`

Orchestrates the print. Roughly:

```go
func PrintBestTurn(ev *Evaluator, snap *bestSnapshot, w io.Writer) {
    snap.Master.SetLogger(gameengine.NewStreamLogger(w))

    // Start-of-turn: hand, arsenal, auras (post-tick — the destroyed sigil
    // is gone; reader sees survivors only, which is the "less pretty" cost).
    fmt.Fprintln(w, "Start of turn:")
    fmt.Fprintf(w, "  Hand: %s\n", handDisplay(snap.Hand))
    if a := snap.Master.Arsenal(); a != nil {
        fmt.Fprintf(w, "  Arsenal: %s\n", a.DisplayName())
    }
    // auras / items lines...

    // My turn: pitches, then chain replay.
    fmt.Fprintln(w, "My turn:")
    for _, a := range snap.BestLine {
        if a.Role == deck.Pitch && !defenseSide(a) {
            fmt.Fprintf(w, "  %s: PITCH\n", a.Card.DisplayName())
        }
    }
    // Chain runs via ev.replayBest (cache-hit replay path) — its chain
    // steps log to the StreamLogger inline. Only ONE chain executes
    // (the captured BestLine), so output is clean — no interleaving
    // from concurrent leafStates.
    ev.replayBestFromSnapshot(snap)  // new helper that wraps replayBest

    // Opp turn: defense pitches, blocks, DR replays (each DR's Play emits
    // its own log lines via the StreamLogger).
    // ...

    // End of turn: final hand / arsenal / auras.
    // ...
}
```

The four-section structure mirrors today's `BuildTurnLog`. Most of `turnlog.go` can be cannibalized for the rendering helpers (`startingHandLine`, `endingArsenalLine`, etc.) — just rewire them to write to `w` instead of appending to a `[]string`.

### 6. Wire `cmd/fabsim/print.go`

Replace `printBestTurn(s deck.Stats)`'s body — it currently calls `sim.FormatTurnLog(s.Best.Log)`. New version calls `stats.PrintBest(os.Stdout)` (or `sim.PrintBestTurn(...)` directly), nil-checking PrintBest so JSON-loaded stats degrade gracefully.

### 7. Delete the dead code

- `deck.TurnLog` struct (`internal/deck/stats.go`) — gone.
- `deck.BestTurn.Log` field — gone.
- `TurnSummary.TriggersFromLastTurn`, `StartOfTurnAuras`, `DealtHand` (`internal/sim/hand_types.go`) — gone.
- `sim.BuildTurnLog`, `sim.FormatTurnLog` and their helpers (`internal/sim/turnlog.go`, `internal/sim/format.go`) — mostly gone; rendering helpers move into the `PrintBestTurn` flow.
- `processAurasAtStartOfTurn`'s `contribs` / `revealed` / `graveyarded` returns — gone (the engine mutates master directly).
- `gameengine.FireStartOfTurn`'s `newEntries` callback parameter — gone.
- `gs.LogEntries()` accessor — gone (no one reads it anymore; tests rewrite).
- The `AmendLastChainStepN` method on `card.Logger` — could stay as a no-op everywhere or be removed (and the one call site in `card/card_state.go` rewritten to emit a separate trigger line).
- Eventually the entire `internal/turnlogger` package can be deleted if no tests still need entry capture. Likely yes; check first.

### 8. JSON shape

`internal/textio/json_types.go::BestTurnJSON` currently has `{Value int, Log deck.TurnLog}`. Drop `Log`. JSON for a `BestTurn` becomes just `{Value, BestLine}` (BestLine isn't there today; add it if you want loaded evals to show the played line). Saved-and-loaded evals lose the chain log on reload — accepted regression per the design mandate.

### 9. Test cleanup

The big ones to plan for:
- `internal/sim/format_test.go` — 34 tests, all assert on output of `FormatBestTurn`. Delete the file entirely; `BuildTurnLog`/`FormatTurnLog` are gone.
- `internal/sim/attack_reaction_test.go` (lines 48-80) — uses `ge.LogEntries()` to assert that `AmendLastChainStepN` folded a buff into the prior chain step's N. Rewrite to assert on the actual buff effect via state, OR install a `*turnlogger.TurnLogger` explicitly and keep the entry-based assertion.
- `internal/sim/resolve_test.go` (lines 90-149) — multiple `ge.LogEntries()` assertions on chain-step entries. Same options: rewrite to assert on state, or install a TurnLogger.
- `internal/sim/opt_debug_test.go` (lines 53, 63) — uses `ge.Logger()` for Opt logging. Probably fine — just keep the default-builder behavior (which still installs `turnlogger.New()`).
- `internal/sim/turn_state_opt_test.go` (lines 172-219) — reads LogEntries. Same handling.
- `internal/sim/deck_aura_test.go` — tests around `ProcessAurasAtStartOfTurn` test export. Rewrite to assert on master state.

The simplest cleanup story: keep `turnlogger.TurnLogger` around as a test helper (it satisfies `card.Logger`), and tests that want entry capture install one explicitly via `ge.SetLogger(turnlogger.New())`. Production never installs it — production uses `NoopLogger` everywhere except print, which uses `StreamLogger`.

### 10. Benchmark + ship

Capture `BenchmarkEvaluate` numbers vs main baseline. Expected: the per-turn alloc pre-fetches dropping should reclaim the previous +1.01% regression from PR #429 and then some — order-of-magnitude estimate is 5-10% wall-clock improvement and ~40% allocs reduction.

Run the standard pre-PR pass: `pr-standards-reviewer` + `comment-simplifier` agents in parallel before pushing the PR per `docs/dev-standards.md` and the user's memory.

## Conversation context (so you understand the why)

The user and the previous session walked through this reasoning:

1. **The original problem**: `processAurasAtStartOfTurn` returns a 4-field struct (`damage`, `contribs`, `revealed`, `graveyarded`) that callers spread to four different sinks. Felt overwrought to the user.
2. **First proposed mitigation** (rejected): keep the struct, add a `wantContribs` flag to skip allocation. User said this was still too much — the goal is no plumbing at all.
3. **The actual mandate**: logging is purely ephemeral. Drop the data plumbing entirely. Cards log straight to stdout (via a writer interface). Eval uses a no-op writer; print uses stdout.
4. **Realized perf upside**: today's `runOneShuffle` pre-fetches `TriggersFromLastTurn` / `StartOfTurnAuras` / `DealtHand` / `startingAuras` / `startingItems` on EVERY turn just so `BuildTurnLog` can format them later, even though only the ~15 new-best turns ever use them. That's ~22M wasted allocs per 200k-shuffle eval. The cleanup is also a 5-10% perf win.
5. **Trade-off accepted**: trigger lines like "Sigil of the Arknight revealed Aether Slash" will disappear from print output since we capture post-tick state and can't reconstruct the tick's per-aura attribution. The reader will see the revealed card in the hand without explanation. "Less pretty is fine."
6. **Cache replay note** (originally a concern, dismissed): cache hits don't re-execute the chain, but that's fine — the print path uses `replayBest` against the captured `BestLine`, which always runs exactly one chain. No interleaving issue.
7. **PR #429 status**: the carryover-helpers refactor is open and will rebase onto this streaming branch once it lands. The two are orthogonal.

## Don't do

- Don't try to preserve the trigger-tree formatting (parent attack with children indented) — the user explicitly accepted flat output.
- Don't try to keep `BuildTurnLog` "for now and clean up later" — the user wants it gone in this PR.
- Don't add a `wantContribs` flag or any other "opt out of metadata" knob — that's the antipattern.
- Don't serialize the snapshot to JSON (you'd have to write a card.Card / hero / aura marshaler for everything; user accepted that loaded evals won't have chain logs).
- Don't open the PR until `pr-standards-reviewer` + `comment-simplifier` agents have run per the user's memory.
