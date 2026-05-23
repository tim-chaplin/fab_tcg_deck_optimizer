# cmd/cardgen

The CLI wrapper around `internal/cardgen`. It regenerates the `<card>_gen.go` files and the
registry's `cardsByID` map from the `<card>.yaml` data files and writes the output to disk.
The generation logic and YAML schema live in `internal/cardgen`.

## Usage

```
go run ./cmd/cardgen [-registry <path>] <dir>...
```

- `<dir>...` — one or more directories holding `<card>.yaml` files. A `<card>_gen.go` is
  written next to each YAML.
- `-registry <path>` — when set, additionally emits the registry `cardsByID` file at that
  path, built from the directory whose package is `cards`.

You normally do not invoke this directly. `go generate ./internal/card/cards/...` runs it via
the `//go:generate` directive in `internal/card/cards/generate.go`, which already passes the
correct directory and `-registry` path.

## Important files

- `main.go` — flag parsing, calls `cardgen.Generate`, writes each returned file to disk.

## Gotchas

- The command exits non-zero on a generation error and writes nothing partial-by-design only
  per file; prefer running it through `go generate` so paths stay correct.
- `internal/lint`'s staleness test fails the build when committed `_gen.go` / `cards_gen.go`
  files drift from the YAML — rerun `go generate` after any YAML edit.
