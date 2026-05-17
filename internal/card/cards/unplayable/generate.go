// See ../generate.go for the codegen contract. Unplayable cards live here so the pool-
// exclusion marker isn't carried by anything in the main cards package.

//go:generate go run ../../../cmd/cardgen .

package unplayable
