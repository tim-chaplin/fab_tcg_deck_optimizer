package sim_test

// Blank-import the registry so its init() populates sim's forward-declared hooks
// (sim.GetCard, sim.DeckableCards, sim.AllWeapons) before any sim_test exercises a path
// that reaches them. Production callers pull the registry in transitively via the cards
// package; sim's own tests need this explicit import because the sim package itself
// can't import registry (registry → cards → sim cycle).

import _ "github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
