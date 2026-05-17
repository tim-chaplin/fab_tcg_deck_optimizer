// Package-level go:generate directive for the cards package. Run `go generate ./...` from
// the repo root to regenerate every <card>_gen.go file from cards.yaml. Hand-written
// <card>.go files (Play, riders, modal helpers, variable Cost, AR predicates, …) are
// untouched.
//
// New cards: edit cards.yaml, run `go generate`, then add the hand-written companion file
// supplying Play().

//go:generate go run ../../cmd/cardgen .

package cards
